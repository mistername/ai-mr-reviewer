package ai

import (
	"context"
	"net/http"

	"github.com/mistername/ai-mr-reviewer/internal/domain"
)

type CopilotClient struct {
	http    *http.Client
	apiKey  string
	baseURL string
	model   string
}

func NewCopilotClient(apiKey, baseURL, model string) *CopilotClient {
	return &CopilotClient{apiKey: apiKey, baseURL: baseURL, model: model, http: &http.Client{}}
}

func (c *CopilotClient) ReviewCode(ctx context.Context, diff string) (string, error) {
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

	return sendChatCompletionRequest(ctx, c.http, "copilot", c.baseURL, c.apiKey, reqBody, "chat", "completions")
}

var _ domain.AIProviderPort = (*CopilotClient)(nil)
