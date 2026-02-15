package application

import (
	"testing"

	"github.com/adlandh/ai-mr-reviewer/internal/domain"
	"go.uber.org/zap"
)

type configMock struct {
	iid string
}

func (c *configMock) GetVCSProvider() string             { return "gitlab" }
func (c *configMock) GetGitLabURL() string               { return "https://gitlab.com" }
func (c *configMock) GetGitLabToken() string             { return "token" }
func (c *configMock) GetProjectID() string               { return "123" }
func (c *configMock) GetMergeRequestIID() string         { return c.iid }
func (c *configMock) GetCommitSHA() string               { return "abc123" }
func (c *configMock) GetMergeRequestDiffBaseSHA() string { return "def456" }
func (c *configMock) GetGitHubToken() string             { return "" }
func (c *configMock) GetGitHubOwner() string             { return "" }
func (c *configMock) GetGitHubRepo() string              { return "" }
func (c *configMock) GetGitHubPRNumber() string          { return "" }
func (c *configMock) GetAIProvider() string              { return "ollama" }
func (c *configMock) GetOllamaURL() string               { return "http://localhost:11434" }
func (c *configMock) GetOllamaAPIKey() string            { return "" }
func (c *configMock) GetOllamaModel() string             { return "llama3.2" }
func (c *configMock) GetOpenAIAPIKey() string            { return "" }
func (c *configMock) GetOpenAIBaseURL() string           { return "https://api.openai.com/v1" }
func (c *configMock) GetOpenAIModel() string             { return "gpt-4" }
func (c *configMock) GetAnthropicAuthToken() string      { return "" }
func (c *configMock) GetAnthropicBaseURL() string        { return "https://api.anthropic.com/v1/" }
func (c *configMock) GetAnthropicModel() string          { return "claude-sonnet-4-20250514" }
func (c *configMock) GetDeleteBotComments() bool         { return false }
func (c *configMock) GetCommentPrefix() string           { return "ai-mr-reviewer" }
func (c *configMock) GetMiniMaxAPIKey() string           { return "" }
func (c *configMock) GetMiniMaxBaseURL() string          { return "https://api.minimax.chat/v1" }
func (c *configMock) GetMiniMaxModel() string            { return "MiniMax-M2.5" }

var _ domain.ConfigPort = (*configMock)(nil)

type mrProviderMock struct {
	comments         map[string][]string
	commentsErr      error
	diffs            []domain.Diff
	diffsErr         error
	addedDiscussions []addedDiscussion
}

type addedDiscussion struct {
	file string
	line int
	body string
}

func (m *mrProviderMock) GetMergeRequestChanges() ([]domain.Diff, error) {
	return m.diffs, m.diffsErr
}
func (m *mrProviderMock) GetExistingComments() (map[string][]string, error) {
	return m.comments, m.commentsErr
}
func (m *mrProviderMock) AddMergeRequestDiscussion(file string, line int, note string) error {
	m.addedDiscussions = append(m.addedDiscussions, addedDiscussion{file: file, line: line, body: note})
	return nil
}
func (m *mrProviderMock) DeleteBotCommentsExceptResolved() error { return nil }

var _ domain.MRProviderPort = (*mrProviderMock)(nil)

type ollamaMock struct {
	response string
	err      error
}

func (m *ollamaMock) ReviewCode(string) (string, error) { return m.response, m.err }

func TestParseReviewResponse(t *testing.T) {
	issues, err := parseReviewResponse("some text {\"issues\":[{\"file\":\"a.go\",\"line\":3,\"severity\":\"warning\",\"message\":\"x\"}]} tail")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 || issues[0].Line != 3 || issues[0].FilePath != "a.go" {
		t.Fatalf("unexpected issues: %+v", issues)
	}
}

func TestDetectLanguage(t *testing.T) {
	if got := detectLanguage("a.go"); got != "Go" {
		t.Fatalf("unexpected language: %s", got)
	}
	if got := detectLanguage("a.unknown"); got != "Unknown" {
		t.Fatalf("unexpected language: %s", got)
	}
}

func TestRunReviewsOnlyNewDiffs(t *testing.T) {
	c := &configMock{iid: "5"}
	g := &mrProviderMock{
		comments: map[string][]string{
			"already.go:1": {"ai-mr-reviewer:**WARNING**: fix it"},
		},
		diffs: []domain.Diff{
			{NewPath: "already.go", Content: "diff1"},
			{NewPath: "new.go", Content: "diff2"},
		},
	}
	o := &ollamaMock{response: `{"issues":[{"file":"new.go","line":10,"severity":"warning","message":"fix it"}]}`}
	r := NewReviewer(c, g, o, zap.NewNop())

	if err := r.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.addedDiscussions) != 1 {
		t.Fatalf("expected 1 discussion, got %d", len(g.addedDiscussions))
	}
	if g.addedDiscussions[0].file != "new.go" || g.addedDiscussions[0].line != 10 {
		t.Fatalf("unexpected discussion: %+v", g.addedDiscussions[0])
	}
}

func TestRunReviewsNewDiffsNoFilter(t *testing.T) {
	c := &configMock{iid: "5"}
	g := &mrProviderMock{
		diffs: []domain.Diff{
			{NewPath: "new.go", Content: "diff2"},
		},
	}
	o := &ollamaMock{response: `{"issues":[{"file":"new.go","line":10,"severity":"warning","message":"fix it"}]}`}
	r := NewReviewer(c, g, o, zap.NewNop())

	if err := r.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.addedDiscussions) != 1 {
		t.Fatalf("expected 1 discussion, got %d", len(g.addedDiscussions))
	}
}
