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

type OpenAIClient struct {
	http    *http.Client
	apiKey  string
	baseURL string
	model   string
}

type openaiRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type openaiResponse struct {
	Choices []choice `json:"choices"`
}

type choice struct {
	Message message `json:"message"`
}

func NewOpenAIClient(apiKey, baseURL, model string) *OpenAIClient {
	return &OpenAIClient{apiKey: apiKey, baseURL: baseURL, model: model, http: &http.Client{}}
}

func (c *OpenAIClient) ReviewCode(filePath, diff, language string) (string, error) {
	prompt := buildReviewPrompt(filePath, diff, language)
	reqBody := openaiRequest{
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

	endpoint, _ := url.JoinPath(c.baseURL, "chat", "completions")
	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("call openai: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openai returned status %d", resp.StatusCode)
	}

	var parsed openaiResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return parsed.Choices[0].Message.Content, nil
}

var _ domain.AIProviderPort = (*OpenAIClient)(nil)
