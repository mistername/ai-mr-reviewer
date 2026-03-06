package main

import (
	"context"
	"testing"
	"time"

	"github.com/adlandh/ai-mr-reviewer/internal/domain"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func TestMainFunctionExists(t *testing.T) {
	_ = main
}

func TestFxAppValidation(t *testing.T) {
	t.Setenv("VCS_PROVIDER", "gitlab")
	t.Setenv("GITLAB_TOKEN", "test-token")
	t.Setenv("CI_PROJECT_ID", "123")
	t.Setenv("CI_MERGE_REQUEST_IID", "1")
	t.Setenv("CI_COMMIT_SHA", "abc123")
	t.Setenv("CI_MERGE_REQUEST_DIFF_BASE_SHA", "def456")
	t.Setenv("AI_PROVIDER", "ollama")

	logger := zap.NewNop()

	err := fx.ValidateApp(
		fx.Supply(logger),
		fx.Provide(
			newConfig,
			newVCSClient,
			newAIProvider,
			newReviewer,
		),
		fx.NopLogger,
	)
	if err != nil {
		t.Fatalf("fx.ValidateApp failed: %v", err)
	}
}

func TestProvidersReturnCorrectTypes(t *testing.T) {
	t.Setenv("VCS_PROVIDER", "gitlab")
	t.Setenv("GITLAB_TOKEN", "test-token")
	t.Setenv("CI_PROJECT_ID", "123")
	t.Setenv("CI_MERGE_REQUEST_IID", "1")
	t.Setenv("CI_COMMIT_SHA", "abc123")
	t.Setenv("CI_MERGE_REQUEST_DIFF_BASE_SHA", "def456")
	t.Setenv("AI_PROVIDER", "ollama")

	cfg, err := newConfig()
	if err != nil {
		t.Fatalf("newConfig failed: %v", err)
	}

	mrProvider, err := newVCSClient(cfg)
	if err != nil {
		t.Fatalf("newVCSClient failed: %v", err)
	}

	aiProvider, err := newAIProvider(cfg)
	if err != nil {
		t.Fatalf("newAIProvider failed: %v", err)
	}

	logger := zap.NewNop()
	reviewer := newReviewer(cfg, mrProvider, aiProvider, logger)
	if reviewer == nil {
		t.Fatal("newReviewer returned nil")
	}
}

type mrProviderMockForTest struct{}

func (mrProviderMockForTest) GetMergeRequestChanges(context.Context) ([]domain.Diff, error) {
	return nil, nil
}
func (mrProviderMockForTest) GetExistingComments(context.Context) (map[string][]string, error) {
	return nil, nil
}
func (mrProviderMockForTest) AddMergeRequestDiscussion(context.Context, string, int, string) error {
	return nil
}
func (mrProviderMockForTest) DeleteBotCommentsExceptResolved(context.Context) error { return nil }

var _ domain.MRProviderPort = (*mrProviderMockForTest)(nil)

type aiProviderMockForTest struct{}

func (o *aiProviderMockForTest) ReviewCode(context.Context, string) (string, error) {
	return `{"issues":[]}`, nil
}

var _ domain.AIProviderPort = (*aiProviderMockForTest)(nil)

type configMockForTest struct{}

func (configMockForTest) GetVCSProvider() string             { return "gitlab" }
func (configMockForTest) GetGitLabURL() string               { return "https://gitlab.com" }
func (configMockForTest) GetGitLabToken() string             { return "token" }
func (configMockForTest) GetProjectID() string               { return "123" }
func (configMockForTest) GetMergeRequestIID() string         { return "1" }
func (configMockForTest) GetCommitSHA() string               { return "abc123" }
func (configMockForTest) GetMergeRequestDiffBaseSHA() string { return "def456" }
func (configMockForTest) GetGitHubToken() string             { return "" }
func (configMockForTest) GetGitHubOwner() string             { return "" }
func (configMockForTest) GetGitHubRepo() string              { return "" }
func (configMockForTest) GetGitHubPRNumber() string          { return "" }
func (configMockForTest) GetAIProvider() string              { return "ollama" }
func (configMockForTest) GetOllamaURL() string               { return "http://localhost:11434" }
func (configMockForTest) GetOllamaAPIKey() string            { return "" }
func (configMockForTest) GetOllamaModel() string             { return "llama3.2" }
func (configMockForTest) GetOpenAIAPIKey() string            { return "" }
func (configMockForTest) GetOpenAIBaseURL() string           { return "https://api.openai.com/v1" }
func (configMockForTest) GetOpenAIModel() string             { return "gpt-4o" }
func (configMockForTest) GetAnthropicAuthToken() string      { return "" }
func (configMockForTest) GetAnthropicBaseURL() string        { return "https://api.anthropic.com/v1/" }
func (configMockForTest) GetAnthropicModel() string          { return "claude-sonnet-4-20250514" }
func (configMockForTest) GetDeleteBotComments() bool         { return false }
func (configMockForTest) GetCommentPrefix() string           { return "ai-mr-reviewer" }
func (configMockForTest) GetMiniMaxAPIKey() string           { return "" }
func (configMockForTest) GetMiniMaxBaseURL() string          { return "https://api.minimax.chat/v1" }
func (configMockForTest) GetMiniMaxModel() string            { return "MiniMax-M2.5" }
func (configMockForTest) GetCopilotBaseURL() string          { return "https://models.github.ai/inference" }
func (configMockForTest) GetCopilotModel() string            { return "openai/gpt-4.1" }
func (configMockForTest) GetRunTimeout() time.Duration       { return 10 * time.Minute }

var _ domain.ConfigPort = (*configMockForTest)(nil)

func TestReviewerWithMocks(t *testing.T) {
	cfg := &configMockForTest{}
	mrProvider := &mrProviderMockForTest{}
	aiProvider := &aiProviderMockForTest{}
	logger := zap.NewNop()

	reviewer := newReviewer(cfg, mrProvider, aiProvider, logger)
	if reviewer == nil {
		t.Fatal("reviewer is nil")
	}
}
