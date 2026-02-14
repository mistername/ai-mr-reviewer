package ai

import (
	"net/http"

	"github.com/adlandh/ai-mr-reviewer/internal/domain"
)

type MiniMaxClient struct {
	http    *http.Client
	apiKey  string
	baseURL string
	model   string
}

func NewMiniMaxClient(apiKey, baseURL, model string) *MiniMaxClient {
	return &MiniMaxClient{apiKey: apiKey, baseURL: baseURL, model: model, http: &http.Client{}}
}

func (c *MiniMaxClient) ReviewCode(diff string) (string, error) {
	prompt := domain.BuildReviewPrompt(diff)

	reqBody := chatCompletionRequest{
		Model: c.model,
		Messages: []message{{
			Role:    "user",
			Content: prompt,
		}},
		Stream:      false,
		Temperature: 0,
	}

	return sendChatCompletionRequest(c.http, "minimax", c.baseURL, c.apiKey, reqBody, "text", "chatcompletion_v2")
}

var _ domain.AIProviderPort = (*MiniMaxClient)(nil)
