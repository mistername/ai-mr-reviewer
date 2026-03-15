package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/adlandh/ai-mr-reviewer/internal/testutil/httpstub"
)

func TestOpenAIClientReviewCode(t *testing.T) {
	t.Parallel()

	client := NewOpenAIClient("openai-key", "https://api.openai.test/v1", "gpt-test")
	client.http = &http.Client{Transport: httpstub.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		assertRequestPath(t, req, "/v1/chat/completions")
		assertHeader(t, req, "Authorization", "Bearer openai-key")
		assertHeader(t, req, "Content-Type", "application/json")

		var body chatCompletionRequest
		decodeJSONBody(t, req, &body)
		if body.Model != "gpt-test" || body.Stream || body.Temperature != 0 {
			t.Fatalf("unexpected request body: %+v", body)
		}
		if len(body.Messages) != 1 || body.Messages[0].Role != "user" || body.Messages[0].Content == "" {
			t.Fatalf("unexpected messages: %+v", body.Messages)
		}

		return httpstub.JSONResponse(http.StatusOK, `{"choices":[{"message":{"content":"openai-ok"}}]}`), nil
	})}

	got, err := client.ReviewCode(context.Background(), "diff --git a/a.go b/a.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "openai-ok" {
		t.Fatalf("unexpected response: %q", got)
	}
}

func TestCopilotClientReviewCode(t *testing.T) {
	t.Parallel()

	client := NewCopilotClient("github-token", "https://models.github.ai/inference", "openai/gpt-4.1")
	client.http = &http.Client{Transport: httpstub.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		assertRequestPath(t, req, "/inference/chat/completions")
		assertHeader(t, req, "Authorization", "Bearer github-token")
		assertHeader(t, req, "Accept", "application/vnd.github+json")
		assertHeader(t, req, "X-GitHub-Api-Version", "2022-11-28")

		var body chatCompletionRequest
		decodeJSONBody(t, req, &body)
		if body.Model != "openai/gpt-4.1" {
			t.Fatalf("unexpected request body: %+v", body)
		}

		return httpstub.JSONResponse(http.StatusOK, `{"choices":[{"message":{"content":"copilot-ok"}}]}`), nil
	})}

	got, err := client.ReviewCode(context.Background(), "diff --git a/a.go b/a.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "copilot-ok" {
		t.Fatalf("unexpected response: %q", got)
	}
}

func TestAnthropicClientReviewCode(t *testing.T) {
	t.Parallel()

	client := NewAnthropicClient("anthropic-token", "https://api.anthropic.test/v1", "claude-test")
	client.http = &http.Client{Transport: httpstub.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		assertRequestPath(t, req, "/v1/messages")
		assertHeader(t, req, "Content-Type", "application/json")
		assertHeader(t, req, "x-api-key", "anthropic-token")
		assertHeader(t, req, "anthropic-version", "2023-06-01")

		var body anthropicRequest
		decodeJSONBody(t, req, &body)
		if body.Model != "claude-test" || body.MaxTokens != 4096 || body.Temperature != 0 {
			t.Fatalf("unexpected request body: %+v", body)
		}
		if len(body.Messages) != 1 || body.Messages[0].Content == "" {
			t.Fatalf("unexpected messages: %+v", body.Messages)
		}

		return httpstub.JSONResponse(http.StatusOK, `{"content":[{"text":"anthropic-ok"}]}`), nil
	})}

	got, err := client.ReviewCode(context.Background(), "diff --git a/a.go b/a.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "anthropic-ok" {
		t.Fatalf("unexpected response: %q", got)
	}
}

func TestAnthropicClientReviewCodeReturnsStatusError(t *testing.T) {
	t.Parallel()

	client := NewAnthropicClient("anthropic-token", "https://api.anthropic.test/v1", "claude-test")
	client.http = &http.Client{Transport: httpstub.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		return httpstub.JSONResponse(http.StatusBadGateway, `{"error":"bad gateway"}`), nil
	})}

	_, err := client.ReviewCode(context.Background(), "diff --git a/a.go b/a.go")
	if err == nil {
		t.Fatal("expected status error")
	}
	if err.Error() != "anthropic returned status 502" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnthropicClientReviewCodeReturnsDecodeError(t *testing.T) {
	t.Parallel()

	client := NewAnthropicClient("anthropic-token", "https://api.anthropic.test/v1", "claude-test")
	client.http = &http.Client{Transport: httpstub.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		return httpstub.JSONResponse(http.StatusOK, `{"content":`), nil
	})}

	_, err := client.ReviewCode(context.Background(), "diff --git a/a.go b/a.go")
	if err == nil {
		t.Fatal("expected decode error")
	}
	if !containsAll(err.Error(), "decode response", "unexpected EOF") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnthropicClientReviewCodeReturnsNoContentError(t *testing.T) {
	t.Parallel()

	client := NewAnthropicClient("anthropic-token", "https://api.anthropic.test/v1", "claude-test")
	client.http = &http.Client{Transport: httpstub.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		return httpstub.JSONResponse(http.StatusOK, `{"content":[]}`), nil
	})}

	_, err := client.ReviewCode(context.Background(), "diff --git a/a.go b/a.go")
	if err == nil {
		t.Fatal("expected no content error")
	}
	if err.Error() != "no content in response" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOllamaClientReviewCode(t *testing.T) {
	t.Parallel()

	client := NewOllamaClient("http://localhost:11434", "ollama-key", "llama3.2")
	client.http = &http.Client{Transport: httpstub.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		assertRequestPath(t, req, "/api/chat")
		assertHeader(t, req, "Content-Type", "application/json")
		assertHeader(t, req, "Authorization", "Bearer ollama-key")

		var body ollamaRequest
		decodeJSONBody(t, req, &body)
		if body.Model != "llama3.2" || body.Stream || body.Options.Temperature != 0 {
			t.Fatalf("unexpected request body: %+v", body)
		}
		if len(body.Messages) != 1 || body.Messages[0].Content == "" {
			t.Fatalf("unexpected messages: %+v", body.Messages)
		}

		return httpstub.JSONResponse(http.StatusOK, `{"message":{"content":"ollama-ok"}}`), nil
	})}

	got, err := client.ReviewCode(context.Background(), "diff --git a/a.go b/a.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ollama-ok" {
		t.Fatalf("unexpected response: %q", got)
	}
}

func TestOllamaClientReviewCodeWithoutAPIKey(t *testing.T) {
	t.Parallel()

	client := NewOllamaClient("http://localhost:11434", "", "llama3.2")
	client.http = &http.Client{Transport: httpstub.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		if got := req.Header.Get("Authorization"); got != "" {
			t.Fatalf("unexpected authorization header: %q", got)
		}

		return httpstub.JSONResponse(http.StatusOK, `{"message":{"content":"ollama-ok"}}`), nil
	})}

	_, err := client.ReviewCode(context.Background(), "diff --git a/a.go b/a.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOllamaClientReviewCodeReturnsStatusError(t *testing.T) {
	t.Parallel()

	client := NewOllamaClient("http://localhost:11434", "ollama-key", "llama3.2")
	client.http = &http.Client{Transport: httpstub.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		return httpstub.JSONResponse(http.StatusBadGateway, `{"error":"bad gateway"}`), nil
	})}

	_, err := client.ReviewCode(context.Background(), "diff --git a/a.go b/a.go")
	if err == nil {
		t.Fatal("expected status error")
	}
	if err.Error() != "ollama returned status 502" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOllamaClientReviewCodeReturnsDecodeError(t *testing.T) {
	t.Parallel()

	client := NewOllamaClient("http://localhost:11434", "ollama-key", "llama3.2")
	client.http = &http.Client{Transport: httpstub.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		return httpstub.JSONResponse(http.StatusOK, `{"message":`), nil
	})}

	_, err := client.ReviewCode(context.Background(), "diff --git a/a.go b/a.go")
	if err == nil {
		t.Fatal("expected decode error")
	}
	if !containsAll(err.Error(), "decode response", "unexpected EOF") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertRequestPath(t *testing.T, req *http.Request, want string) {
	t.Helper()

	if req.URL.Path != want {
		t.Fatalf("unexpected path: got %q want %q", req.URL.Path, want)
	}
}

func assertHeader(t *testing.T, req *http.Request, name, want string) {
	t.Helper()

	if got := req.Header.Get(name); got != want {
		t.Fatalf("unexpected %s header: got %q want %q", name, got, want)
	}
}

func decodeJSONBody(t *testing.T, req *http.Request, target any) {
	t.Helper()

	if err := json.NewDecoder(req.Body).Decode(target); err != nil {
		t.Fatalf("decode request body: %v", err)
	}
}

func containsAll(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if !strings.Contains(s, substr) {
			return false
		}
	}

	return true
}
