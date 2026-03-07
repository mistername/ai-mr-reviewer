package github

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	gogithub "github.com/google/go-github/v82/github"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}

func TestClientGetMergeRequestChanges(t *testing.T) {
	t.Parallel()

	apiClient := gogithub.NewClient(&http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodGet {
				t.Fatalf("unexpected method: %s", r.Method)
			}
			if r.URL.Path != "/repos/acme/repo/pulls/7/files" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}

			return jsonResponse(http.StatusOK, `[
			{"filename":"new.go","patch":"@@ -1 +1 @@","previous_filename":"old.go"},
			{"filename":"same.go","patch":"@@ -2 +2 @@"}
		]`), nil
		}),
	})

	baseURL, err := url.Parse("https://api.github.test/")
	if err != nil {
		t.Fatalf("parse base url: %v", err)
	}
	apiClient.BaseURL = baseURL

	client := &Client{
		client:   apiClient,
		owner:    "acme",
		repo:     "repo",
		prNumber: 7,
	}

	diffs, err := client.GetMergeRequestChanges(context.Background())
	if err != nil {
		t.Fatalf("GetMergeRequestChanges returned error: %v", err)
	}
	if len(diffs) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(diffs))
	}
	if diffs[0].NewPath != "new.go" || diffs[0].OldPath != "old.go" || diffs[0].Content != "@@ -1 +1 @@" {
		t.Fatalf("unexpected first diff: %+v", diffs[0])
	}
	if diffs[1].NewPath != "same.go" || diffs[1].OldPath != "" || diffs[1].Content != "@@ -2 +2 @@" {
		t.Fatalf("unexpected second diff: %+v", diffs[1])
	}
}

func TestClientAddMergeRequestDiscussionFallsBackToIssueComment(t *testing.T) {
	t.Parallel()

	var issueCommentBody string
	apiClient := gogithub.NewClient(&http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			switch {
			case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/repo/pulls/7/comments":
				return jsonResponse(http.StatusUnprocessableEntity, `{"message":"validation failed"}`), nil
			case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/repo/issues/7/comments":
				var payload struct {
					Body string `json:"body"`
				}
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Fatalf("decode issue comment payload: %v", err)
				}
				issueCommentBody = payload.Body
				return jsonResponse(http.StatusCreated, `{"id":1}`), nil
			default:
				t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
				return nil, nil
			}
		}),
	})

	baseURL, err := url.Parse("https://api.github.test/")
	if err != nil {
		t.Fatalf("parse base url: %v", err)
	}
	apiClient.BaseURL = baseURL

	client := &Client{
		client:        apiClient,
		commentPrefix: "ai-mr-reviewer",
		owner:         "acme",
		repo:          "repo",
		commitSHA:     "abc123",
		prNumber:      7,
	}

	err = client.AddMergeRequestDiscussion(context.Background(), "foo.go", 12, "please fix this")
	if err != nil {
		t.Fatalf("AddMergeRequestDiscussion returned error: %v", err)
	}

	want := "ai-mr-reviewer:**File: foo.go**\n\nplease fix this"
	if issueCommentBody != want {
		t.Fatalf("unexpected fallback issue comment body: %q", issueCommentBody)
	}
}
