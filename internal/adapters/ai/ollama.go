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
	model   string
}

type ollamaRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaResponse struct {
	Message message `json:"message"`
}

func NewOllamaClient(baseURL, model string) *OllamaClient {
	return &OllamaClient{baseURL: baseURL, model: model, http: &http.Client{}}
}

func (c *OllamaClient) ReviewCode(filePath, diff, language string) (string, error) {
	prompt := buildReviewPrompt(filePath, diff, language)
	reqBody := ollamaRequest{
		Model: c.model,
		Messages: []message{{
			Role:    "user",
			Content: prompt,
		}},
		Stream: false,
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

	resp, err := c.http.Do(httpReq)
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

func buildReviewPrompt(filePath, diff, language string) string {
	return fmt.Sprintf(`You are a strict code reviewer.
Review this Git diff and return issues tied to exact NEW file line numbers.
Focus on bugs, security, maintainability, and correctness.

File: %s
Language: %s

Diff:
%s

Return only valid JSON in this format:
{
  "issues": [
    {
      "line": 12,
      "severity": "error|warning|info",
      "message": "text"
    }
  ]
}

If no issues:
{"issues": []}
`, filePath, language, diff)
}

var _ domain.AIProviderPort = (*OllamaClient)(nil)
