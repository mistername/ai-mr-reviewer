package ai

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestSendChatCompletionRequestSetsGitHubModelsHeadersForCopilot(t *testing.T) {
	t.Parallel()

	var got *http.Request
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			got = req
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"choices":[{"message":{"content":"ok"}}]}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	_, err := sendChatCompletionRequest(
		context.Background(),
		client,
		"copilot",
		"https://models.github.ai/inference",
		"test-token",
		chatCompletionRequest{Model: "openai/gpt-4.1"},
		"chat",
		"completions",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got == nil {
		t.Fatal("expected request to be captured")
	}
	if auth := got.Header.Get("Authorization"); auth != "Bearer test-token" {
		t.Fatalf("unexpected authorization header: %q", auth)
	}
	if accept := got.Header.Get("Accept"); accept != "application/vnd.github+json" {
		t.Fatalf("unexpected accept header: %q", accept)
	}
	if version := got.Header.Get("X-GitHub-Api-Version"); version != "2022-11-28" {
		t.Fatalf("unexpected github api version: %q", version)
	}
}
