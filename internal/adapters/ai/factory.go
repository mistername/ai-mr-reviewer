package ai

import (
	"fmt"

	"github.com/adlandh/ai-mr-reviewer/internal/domain"
)

func NewAIProvider(config domain.AIConfig) (domain.AIProviderPort, error) {
	provider := config.Provider

	switch provider {
	case "ollama":
		return NewOllamaClient(config.Ollama.URL, config.Ollama.APIKey, config.Ollama.Model), nil
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

func newOpenAIProvider(config domain.AIConfig) (domain.AIProviderPort, error) {
	if config.OpenAI.APIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required for openai provider")
	}

	return NewOpenAIClient(config.OpenAI.APIKey, config.OpenAI.BaseURL, config.OpenAI.Model), nil
}

func newAnthropicProvider(config domain.AIConfig) (domain.AIProviderPort, error) {
	if config.Anthropic.AuthToken == "" {
		return nil, fmt.Errorf("ANTHROPIC_AUTH_TOKEN is required for anthropic provider")
	}

	return NewAnthropicClient(config.Anthropic.AuthToken, config.Anthropic.BaseURL, config.Anthropic.Model), nil
}

func newCopilotProvider(config domain.AIConfig) (domain.AIProviderPort, error) {
	if config.Copilot.Token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN is required for copilot provider")
	}

	return NewCopilotClient(config.Copilot.Token, config.Copilot.BaseURL, config.Copilot.Model), nil
}
