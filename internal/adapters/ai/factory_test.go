package ai

import (
	"strings"
	"testing"

	"github.com/adlandh/ai-mr-reviewer/internal/domain"
	"github.com/adlandh/ai-mr-reviewer/internal/domain/mocks"
)

func TestNewAIProviderCreatesSupportedProviders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		provider  string
		setupMock func(*mocks.ConfigPort)
		assert    func(*testing.T, domain.AIProviderPort)
	}{
		{
			name:     "ollama",
			provider: "ollama",
			setupMock: func(config *mocks.ConfigPort) {
				config.EXPECT().GetOllamaURL().Return("http://localhost:11434")
				config.EXPECT().GetOllamaAPIKey().Return("ollama-key")
				config.EXPECT().GetOllamaModel().Return("llama3.2")
			},
			assert: func(t *testing.T, provider domain.AIProviderPort) {
				t.Helper()

				client, ok := provider.(*OllamaClient)
				if !ok {
					t.Fatalf("expected *OllamaClient, got %T", provider)
				}

				if client.baseURL != "http://localhost:11434" || client.apiKey != "ollama-key" || client.model != "llama3.2" {
					t.Fatalf("unexpected ollama client: %+v", client)
				}
			},
		},
		{
			name:     "openai",
			provider: "openai",
			setupMock: func(config *mocks.ConfigPort) {
				config.EXPECT().GetOpenAIAPIKey().Return("openai-key").Times(2)
				config.EXPECT().GetOpenAIBaseURL().Return("https://api.openai.test/v1")
				config.EXPECT().GetOpenAIModel().Return("gpt-test")
			},
			assert: func(t *testing.T, provider domain.AIProviderPort) {
				t.Helper()

				client, ok := provider.(*OpenAIClient)
				if !ok {
					t.Fatalf("expected *OpenAIClient, got %T", provider)
				}

				if client.baseURL != "https://api.openai.test/v1" || client.apiKey != "openai-key" || client.model != "gpt-test" {
					t.Fatalf("unexpected openai client: %+v", client)
				}
			},
		},
		{
			name:     "anthropic",
			provider: "anthropic",
			setupMock: func(config *mocks.ConfigPort) {
				config.EXPECT().GetAnthropicAuthToken().Return("anthropic-token").Times(2)
				config.EXPECT().GetAnthropicBaseURL().Return("https://api.anthropic.test/v1")
				config.EXPECT().GetAnthropicModel().Return("claude-test")
			},
			assert: func(t *testing.T, provider domain.AIProviderPort) {
				t.Helper()

				client, ok := provider.(*AnthropicClient)
				if !ok {
					t.Fatalf("expected *AnthropicClient, got %T", provider)
				}

				if client.baseURL != "https://api.anthropic.test/v1" || client.apiKey != "anthropic-token" || client.model != "claude-test" {
					t.Fatalf("unexpected anthropic client: %+v", client)
				}
			},
		},
		{
			name:     "copilot",
			provider: "copilot",
			setupMock: func(config *mocks.ConfigPort) {
				config.EXPECT().GetGitHubToken().Return("github-token").Times(2)
				config.EXPECT().GetCopilotBaseURL().Return("https://models.github.ai/inference")
				config.EXPECT().GetCopilotModel().Return("openai/gpt-4.1")
			},
			assert: func(t *testing.T, provider domain.AIProviderPort) {
				t.Helper()

				client, ok := provider.(*CopilotClient)
				if !ok {
					t.Fatalf("expected *CopilotClient, got %T", provider)
				}

				if client.baseURL != "https://models.github.ai/inference" || client.apiKey != "github-token" || client.model != "openai/gpt-4.1" {
					t.Fatalf("unexpected copilot client: %+v", client)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config := mocks.NewConfigPort(t)
			config.EXPECT().GetAIProvider().Return(tt.provider)
			tt.setupMock(config)

			provider, err := NewAIProvider(config)
			if err != nil {
				t.Fatalf("NewAIProvider returned error: %v", err)
			}

			tt.assert(t, provider)
		})
	}
}

func TestNewAIProviderRejectsMiniMax(t *testing.T) {
	t.Parallel()

	config := mocks.NewConfigPort(t)
	config.EXPECT().GetAIProvider().Return("minimax")

	_, err := NewAIProvider(config)
	if err == nil {
		t.Fatal("expected error for removed MiniMax provider")
	}

	if !strings.Contains(err.Error(), "unsupported AI provider: minimax") {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(err.Error(), "ollama, openai, anthropic, copilot") {
		t.Fatalf("unexpected supported providers list: %v", err)
	}
}

func TestNewAIProviderRequiresProviderCredentials(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		provider    string
		setupMock   func(*mocks.ConfigPort)
		errorSubstr string
	}{
		{
			name:     "openai",
			provider: "openai",
			setupMock: func(config *mocks.ConfigPort) {
				config.EXPECT().GetOpenAIAPIKey().Return("")
			},
			errorSubstr: "OPENAI_API_KEY is required",
		},
		{
			name:     "anthropic",
			provider: "anthropic",
			setupMock: func(config *mocks.ConfigPort) {
				config.EXPECT().GetAnthropicAuthToken().Return("")
			},
			errorSubstr: "ANTHROPIC_AUTH_TOKEN is required",
		},
		{
			name:     "copilot",
			provider: "copilot",
			setupMock: func(config *mocks.ConfigPort) {
				config.EXPECT().GetGitHubToken().Return("")
			},
			errorSubstr: "GITHUB_TOKEN is required",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config := mocks.NewConfigPort(t)
			config.EXPECT().GetAIProvider().Return(tt.provider)
			tt.setupMock(config)

			_, err := NewAIProvider(config)
			if err == nil {
				t.Fatal("expected credential validation error")
			}

			if !strings.Contains(err.Error(), tt.errorSubstr) {
				t.Fatalf("expected error containing %q, got %v", tt.errorSubstr, err)
			}
		})
	}
}
