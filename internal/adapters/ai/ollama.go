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

type OllamaClient struct {
	http    *http.Client
	baseURL string
	apiKey  string
	model   string
}

type ollamaRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool      `json:"stream"`
	Options  options   `json:"options"`
}

type options struct {
	Temperature float64 `json:"temperature"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaResponse struct {
	Message message `json:"message"`
}

func NewOllamaClient(baseURL, apiKey, model string) *OllamaClient {
	return &OllamaClient{baseURL: baseURL, apiKey: apiKey, model: model, http: &http.Client{}}
}

func (c *OllamaClient) ReviewCode(diff string) (string, error) {
	prompt := domain.BuildReviewPrompt(diff)
	reqBody := ollamaRequest{
		Model: c.model,
		Messages: []message{{
			Role:    "user",
			Content: prompt,
		}},
		Stream: false,
		Options: options{
			Temperature: 0,
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	endpoint, _ := url.JoinPath(c.baseURL, "api", "chat")

	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.http.Do(httpReq) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("call ollama: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var parsed ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return parsed.Message.Content, nil
}

var _ domain.AIProviderPort = (*OllamaClient)(nil)
