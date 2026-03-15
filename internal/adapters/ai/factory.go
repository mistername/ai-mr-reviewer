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
		return newOpenAIProvider(config)
	case "anthropic":
		return newAnthropicProvider(config)
	case "copilot":
		return newCopilotProvider(config)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s (supported: ollama, openai, anthropic, copilot)", provider)
	}
}

func newOpenAIProvider(config domain.ConfigPort) (domain.AIProviderPort, error) {
	if config.GetOpenAIAPIKey() == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required for openai provider")
	}

	return NewOpenAIClient(config.GetOpenAIAPIKey(), config.GetOpenAIBaseURL(), config.GetOpenAIModel()), nil
}

func newAnthropicProvider(config domain.ConfigPort) (domain.AIProviderPort, error) {
	if config.GetAnthropicAuthToken() == "" {
		return nil, fmt.Errorf("ANTHROPIC_AUTH_TOKEN is required for anthropic provider")
	}

	return NewAnthropicClient(config.GetAnthropicAuthToken(), config.GetAnthropicBaseURL(), config.GetAnthropicModel()), nil
}

func newCopilotProvider(config domain.ConfigPort) (domain.AIProviderPort, error) {
	if config.GetGitHubToken() == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN is required for copilot provider")
	}

	return NewCopilotClient(config.GetGitHubToken(), config.GetCopilotBaseURL(), config.GetCopilotModel()), nil
}
