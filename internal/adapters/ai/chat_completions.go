package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type chatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Stream      bool      `json:"stream"`
	Temperature float64   `json:"temperature"`
}

type chatCompletionResponse struct {
	Choices []choice `json:"choices"`
}

type choice struct {
	Message message `json:"message"`
}

func sendChatCompletionRequest(
	ctx context.Context,
	httpClient *http.Client,
	providerName string,
	baseURL string,
	apiKey string,
	reqBody chatCompletionRequest,
	endpointParts ...string,
) (string, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	endpoint, _ := url.JoinPath(baseURL, endpointParts...)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := httpClient.Do(httpReq) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("call %s: %w", providerName, err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s returned status %d", providerName, resp.StatusCode)
	}

	var parsed chatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return parsed.Choices[0].Message.Content, nil
}
