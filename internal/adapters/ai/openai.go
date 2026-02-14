package ai

import (
	"net/http"

	"github.com/adlandh/ai-mr-reviewer/internal/domain"
)

type OpenAIClient struct {
	http    *http.Client
	apiKey  string
	baseURL string
	model   string
}

func NewOpenAIClient(apiKey, baseURL, model string) *OpenAIClient {
	return &OpenAIClient{apiKey: apiKey, baseURL: baseURL, model: model, http: &http.Client{}}
}

func (c *OpenAIClient) ReviewCode(filePath, diff, language string) (string, error) {
	prompt := domain.BuildReviewPrompt(filePath, diff, language)

	reqBody := chatCompletionRequest{
		Model: c.model,
		Messages: []message{{
			Role:    "user",
			Content: prompt,
		}},
		Stream:      false,
		Temperature: 0,
	}

	return sendChatCompletionRequest(c.http, "openai", c.baseURL, c.apiKey, reqBody, "chat", "completions")
}

var _ domain.AIProviderPort = (*OpenAIClient)(nil)
