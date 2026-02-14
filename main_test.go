package main

import (
	"testing"

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

func (g *mrProviderMockForTest) GetMergeRequestChanges() ([]domain.Diff, error) {
	return nil, nil
}
func (g *mrProviderMockForTest) GetExistingComments() (map[string][]string, error) {
	return nil, nil
}
func (g *mrProviderMockForTest) AddMergeRequestDiscussion(string, int, string) error { return nil }
func (g *mrProviderMockForTest) DeleteBotCommentsExceptResolved() error              { return nil }

var _ domain.MRProviderPort = (*mrProviderMockForTest)(nil)

type aiProviderMockForTest struct{}

func (o *aiProviderMockForTest) ReviewCode(string, string, string) (string, error) {
	return `{"issues":[]}`, nil
}

var _ domain.AIProviderPort = (*aiProviderMockForTest)(nil)

type configMockForTest struct{}

func (c *configMockForTest) GetVCSProvider() string             { return "gitlab" }
func (c *configMockForTest) GetGitLabURL() string               { return "https://gitlab.com" }
func (c *configMockForTest) GetGitLabToken() string             { return "token" }
func (c *configMockForTest) GetProjectID() string               { return "123" }
func (c *configMockForTest) GetMergeRequestIID() string         { return "1" }
func (c *configMockForTest) GetCommitSHA() string               { return "abc123" }
func (c *configMockForTest) GetMergeRequestDiffBaseSHA() string { return "def456" }
func (c *configMockForTest) GetGitHubToken() string             { return "" }
func (c *configMockForTest) GetGitHubOwner() string             { return "" }
func (c *configMockForTest) GetGitHubRepo() string              { return "" }
func (c *configMockForTest) GetGitHubPRNumber() string          { return "" }
func (c *configMockForTest) GetAIProvider() string              { return "ollama" }
func (c *configMockForTest) GetOllamaURL() string               { return "http://localhost:11434" }
func (c *configMockForTest) GetOllamaModel() string             { return "llama3.2" }
func (c *configMockForTest) GetOpenAIAPIKey() string            { return "" }
func (c *configMockForTest) GetOpenAIBaseURL() string           { return "https://api.openai.com/v1" }
func (c *configMockForTest) GetOpenAIModel() string             { return "gpt-4o" }
func (c *configMockForTest) GetAnthropicAuthToken() string      { return "" }
func (c *configMockForTest) GetAnthropicBaseURL() string        { return "https://api.anthropic.com/v1/" }
func (c *configMockForTest) GetAnthropicModel() string          { return "claude-sonnet-4-20250514" }
func (c *configMockForTest) GetDeleteBotComments() bool         { return false }

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
