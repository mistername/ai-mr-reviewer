package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/adlandh/ai-mr-reviewer/internal/domain"
)

type MiniMaxClient struct {
	http    *http.Client
	apiKey  string
	baseURL string
	model   string
}

type minimaxRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Stream      bool      `json:"stream"`
	Temperature float64   `json:"temperature"`
}

type minimaxResponse struct {
	Choices []choice `json:"choices"`
	Usage   usage    `json:"usage,omitempty"`
}

type usage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

func NewMiniMaxClient(apiKey, baseURL, model string) *MiniMaxClient {
	return &MiniMaxClient{apiKey: apiKey, baseURL: baseURL, model: model, http: &http.Client{}}
}

func (c *MiniMaxClient) ReviewCode(filePath, diff, language string) (string, error) {
	prompt := buildReviewPrompt(filePath, diff, language)
	reqBody := minimaxRequest{
		Model: c.model,
		Messages: []message{{
			Role:    "user",
			Content: prompt,
		}},
		Stream:      false,
		Temperature: 0,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	endpoint, _ := url.JoinPath(c.baseURL, "text", "chatcompletion_v2")

	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("call minimax: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("minimax returned status %d", resp.StatusCode)
	}

	var parsed minimaxResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return parsed.Choices[0].Message.Content, nil
}

var _ domain.AIProviderPort = (*MiniMaxClient)(nil)
