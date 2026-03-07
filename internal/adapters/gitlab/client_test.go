package gitlab

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/adlandh/ai-mr-reviewer/internal/domain/mocks"
	gogitlab "gitlab.com/gitlab-org/api/client-go"
)

const testGitLabBaseURL = "https://gitlab.example.com/api/v4"
const testGitLabMRPath = "/api/v4/projects/123/merge_requests/5"
const testGitLabNotesPathPrefix = testGitLabMRPath + "/notes/"

type discussionRequest struct {
	Body     string `json:"body"`
	Position struct {
		PositionType string `json:"position_type"`
		NewLine      int64  `json:"new_line"`
		HeadSHA      string `json:"head_sha"`
		BaseSHA      string `json:"base_sha"`
		StartSHA     string `json:"start_sha"`
		OldPath      string `json:"old_path"`
		NewPath      string `json:"new_path"`
	} `json:"position"`
}

func TestClientAddMergeRequestDiscussionIncludesPositionAndPrefix(t *testing.T) {
	t.Parallel()

	cfg := mocks.NewConfigPort(t)
	cfg.EXPECT().GetCommitSHA().Return("head-sha")
	cfg.EXPECT().GetMergeRequestDiffBaseSHA().Return("base-sha")
	cfg.EXPECT().GetCommentPrefix().Return("ai-mr-reviewer")

	var got discussionRequest

	client, err := NewClient(testGitLabBaseURL, "token", "123", 5, cfg)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	client.git, err = newStubGitLabClient(t, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != testGitLabMRPath+"/discussions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		return jsonResponse(http.StatusCreated, `{}`), nil
	}))
	if err != nil {
		t.Fatalf("create stub gitlab client: %v", err)
	}

	err = client.AddMergeRequestDiscussion(context.Background(), "foo.go", 42, "please fix this")
	if err != nil {
		t.Fatalf("AddMergeRequestDiscussion returned error: %v", err)
	}

	assertDiscussionRequest(t, got)
}

func TestClientDeleteBotCommentsExceptResolvedDeletesOnlyUnresolvedBotNotes(t *testing.T) {
	t.Parallel()

	cfg := mocks.NewConfigPort(t)
	cfg.On("GetCommentPrefix").Return("ai-mr-reviewer").Maybe()

	var deleted []string
	client, err := NewClient(testGitLabBaseURL, "token", "123", 5, cfg)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	client.git, err = gogitlab.NewClient(
		"token",
		gogitlab.WithBaseURL(testGitLabBaseURL),
		gogitlab.WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				switch {
				case r.Method == http.MethodGet && r.URL.Path == testGitLabMRPath+"/notes":
					return jsonResponse(http.StatusOK, `[
						{"id":1,"body":"ai-mr-reviewer: first","resolved":false,"system":false},
						{"id":2,"body":"ai-mr-reviewer: resolved","resolved":true,"system":false},
						{"id":3,"body":"system note","resolved":false,"system":true},
						{"id":4,"body":"human note","resolved":false,"system":false}
					]`), nil
				case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, testGitLabNotesPathPrefix):
					deleted = append(deleted, strings.TrimPrefix(r.URL.Path, testGitLabNotesPathPrefix))
					return &http.Response{
						StatusCode: http.StatusNoContent,
						Body:       io.NopCloser(bytes.NewBuffer(nil)),
						Header:     make(http.Header),
					}, nil
				default:
					t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
					return nil, nil
				}
			}),
		}),
	)
	if err != nil {
		t.Fatalf("create stub gitlab client: %v", err)
	}

	err = client.DeleteBotCommentsExceptResolved(context.Background())
	if err != nil {
		t.Fatalf("DeleteBotCommentsExceptResolved returned error: %v", err)
	}

	if len(deleted) != 1 || deleted[0] != "1" {
		t.Fatalf("unexpected deleted note ids: %v", deleted)
	}
}

func newStubGitLabClient(t *testing.T, transport roundTripFunc) (*gogitlab.Client, error) {
	t.Helper()

	return gogitlab.NewClient(
		"token",
		gogitlab.WithBaseURL(testGitLabBaseURL),
		gogitlab.WithHTTPClient(&http.Client{
			Transport: transport,
		}),
	)
}

func assertDiscussionRequest(t *testing.T, got discussionRequest) {
	t.Helper()

	if got.Body != "ai-mr-reviewer: please fix this" {
		t.Fatalf("unexpected discussion body: %q", got.Body)
	}
	if got.Position.PositionType != "line" {
		t.Fatalf("unexpected position type: %q", got.Position.PositionType)
	}
	if got.Position.NewLine != 42 {
		t.Fatalf("unexpected new line: %d", got.Position.NewLine)
	}
	if got.Position.HeadSHA != "head-sha" || got.Position.BaseSHA != "base-sha" || got.Position.StartSHA != "head-sha" {
		t.Fatalf("unexpected SHAs: %+v", got.Position)
	}
	if got.Position.OldPath != "foo.go" || got.Position.NewPath != "foo.go" {
		t.Fatalf("unexpected paths: %+v", got.Position)
	}
}

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
