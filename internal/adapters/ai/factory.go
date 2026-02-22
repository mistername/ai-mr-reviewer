package ai

import (
	"fmt"

	"github.com/adlandh/ai-mr-reviewer/internal/domain"
)

func NewAIProvider(config domain.ConfigPort) (domain.AIProviderPort, error) {
	provider := config.GetAIProvider()

	switch provider {
	case "ollama":
		return NewOllamaClient(config.GetOllamaURL(), config.GetOllamaAPIKey(), config.GetOllamaModel()), nil

	case "openai":
		if config.GetOpenAIAPIKey() == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY is required for openai provider")
		}

		return NewOpenAIClient(config.GetOpenAIAPIKey(), config.GetOpenAIBaseURL(), config.GetOpenAIModel()), nil

	case "anthropic":
		if config.GetAnthropicAuthToken() == "" {
			return nil, fmt.Errorf("ANTHROPIC_AUTH_TOKEN is required for anthropic provider")
		}

		return NewAnthropicClient(config.GetAnthropicAuthToken(), config.GetAnthropicBaseURL(), config.GetAnthropicModel()), nil

	case "minimax":
		if config.GetMiniMaxAPIKey() == "" {
			return nil, fmt.Errorf("MINIMAX_API_KEY is required for minimax provider")
		}

		return NewMiniMaxClient(config.GetMiniMaxAPIKey(), config.GetMiniMaxBaseURL(), config.GetMiniMaxModel()), nil

	case "copilot":
		if config.GetGitHubToken() == "" {
			return nil, fmt.Errorf("GITHUB_TOKEN is required for copilot provider")
		}

		return NewCopilotClient(config.GetGitHubToken(), config.GetCopilotBaseURL(), config.GetCopilotModel()), nil

	default:
		return nil, fmt.Errorf("unsupported AI provider: %s (supported: ollama, openai, anthropic, minimax, copilot)", provider)
	}
}
