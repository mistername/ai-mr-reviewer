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

type AnthropicClient struct {
	http    *http.Client
	apiKey  string
	baseURL string
	model   string
}

type anthropicRequest struct {
	Model     string    `json:"model"`
	Messages  []message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type anthropicResponse struct {
	Content []contentBlock `json:"content"`
}

type contentBlock struct {
	Text string `json:"text"`
}

func NewAnthropicClient(apiKey, baseURL, model string) *AnthropicClient {
	return &AnthropicClient{apiKey: apiKey, baseURL: baseURL, model: model, http: &http.Client{}}
}

func (c *AnthropicClient) ReviewCode(filePath, diff, language string) (string, error) {
	prompt := buildReviewPrompt(filePath, diff, language)
	reqBody := anthropicRequest{
		Model: c.model,
		Messages: []message{{
			Role:    "user",
			Content: prompt,
		}},
		MaxTokens: 4096,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	endpoint, _ := url.JoinPath(c.baseURL, "messages")
	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("call anthropic: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("anthropic returned status %d", resp.StatusCode)
	}

	var parsed anthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(parsed.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return parsed.Content[0].Text, nil
}

var _ domain.AIProviderPort = (*AnthropicClient)(nil)
