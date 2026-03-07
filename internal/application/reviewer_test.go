package application

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/adlandh/ai-mr-reviewer/internal/domain"
	"go.uber.org/zap"
)

type configMock struct {
	iid               string
	deleteBotComments bool
}

func (configMock) GetVCSProvider() string             { return "gitlab" }
func (configMock) GetGitLabURL() string               { return "https://gitlab.com" }
func (configMock) GetGitLabToken() string             { return "token" }
func (configMock) GetProjectID() string               { return "123" }
func (c configMock) GetMergeRequestIID() string       { return c.iid }
func (configMock) GetCommitSHA() string               { return "abc123" }
func (configMock) GetMergeRequestDiffBaseSHA() string { return "def456" }
func (configMock) GetGitHubToken() string             { return "" }
func (configMock) GetGitHubOwner() string             { return "" }
func (configMock) GetGitHubRepo() string              { return "" }
func (configMock) GetGitHubPRNumber() string          { return "" }
func (configMock) GetAIProvider() string              { return "ollama" }
func (configMock) GetOllamaURL() string               { return "http://localhost:11434" }
func (configMock) GetOllamaAPIKey() string            { return "" }
func (configMock) GetOllamaModel() string             { return "llama3.2" }
func (configMock) GetOpenAIAPIKey() string            { return "" }
func (configMock) GetOpenAIBaseURL() string           { return "https://api.openai.com/v1" }
func (configMock) GetOpenAIModel() string             { return "gpt-4" }
func (configMock) GetAnthropicAuthToken() string      { return "" }
func (configMock) GetAnthropicBaseURL() string        { return "https://api.anthropic.com/v1/" }
func (configMock) GetAnthropicModel() string          { return "claude-sonnet-4-20250514" }
func (c configMock) GetDeleteBotComments() bool       { return c.deleteBotComments }
func (configMock) GetCommentPrefix() string           { return "ai-mr-reviewer" }
func (configMock) GetMiniMaxAPIKey() string           { return "" }
func (configMock) GetMiniMaxBaseURL() string          { return "https://api.minimax.chat/v1" }
func (configMock) GetMiniMaxModel() string            { return "MiniMax-M2.5" }
func (configMock) GetCopilotBaseURL() string          { return "https://models.github.ai/inference" }
func (configMock) GetCopilotModel() string            { return "openai/gpt-4.1" }
func (configMock) GetRunTimeout() time.Duration       { return 10 * time.Minute }

var _ domain.ConfigPort = (*configMock)(nil)

type mrProviderMock struct {
	comments         map[string][]string
	commentsErr      error
	diffs            []domain.Diff
	diffsErr         error
	addedDiscussions []addedDiscussion
	deleteErr        error
	deleteCalls      int
	addErr           error
}

type addedDiscussion struct {
	file string
	line int
	body string
}

func (m *mrProviderMock) GetMergeRequestChanges(context.Context) ([]domain.Diff, error) {
	return m.diffs, m.diffsErr
}
func (m *mrProviderMock) GetExistingComments(context.Context) (map[string][]string, error) {
	return m.comments, m.commentsErr
}
func (m *mrProviderMock) AddMergeRequestDiscussion(_ context.Context, file string, line int, note string) error {
	m.addedDiscussions = append(m.addedDiscussions, addedDiscussion{file: file, line: line, body: note})
	return m.addErr
}
func (m *mrProviderMock) DeleteBotCommentsExceptResolved(context.Context) error {
	m.deleteCalls++
	return m.deleteErr
}

var _ domain.MRProviderPort = (*mrProviderMock)(nil)

type ollamaMock struct {
	response string
	err      error
}

func (m *ollamaMock) ReviewCode(context.Context, string) (string, error) { return m.response, m.err }

type blockingAIMock struct{}

func (blockingAIMock) ReviewCode(ctx context.Context, _ string) (string, error) {
	<-ctx.Done()
	return "", ctx.Err()
}

func TestParseReviewResponse(t *testing.T) {
	issues, err := parseReviewResponse("some text {\"issues\":[{\"file\":\"a.go\",\"line\":3,\"severity\":\"warning\",\"message\":\"x\"}]} tail")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 || issues[0].Line != 3 || issues[0].FilePath != "a.go" {
		t.Fatalf("unexpected issues: %+v", issues)
	}
}

func TestParseReviewResponseUsesPathFallback(t *testing.T) {
	issues, err := parseReviewResponse(`{"issues":[{"path":"a.go","line":4,"severity":"warning","message":"x"}]}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 || issues[0].FilePath != "a.go" || issues[0].Line != 4 {
		t.Fatalf("unexpected issues: %+v", issues)
	}
}

func TestParseReviewResponseExtractsJSONFromCodeFence(t *testing.T) {
	issues, err := parseReviewResponse("```json\n{\"issues\":[{\"file\":\"a.go\",\"line\":3,\"severity\":\"warning\",\"message\":\"x\"}]}\n```")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 || issues[0].FilePath != "a.go" {
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

	if err := r.Run(context.Background()); err != nil {
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

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.addedDiscussions) != 1 {
		t.Fatalf("expected 1 discussion, got %d", len(g.addedDiscussions))
	}
}

func TestRunContinuesWhenExistingCommentsFail(t *testing.T) {
	c := &configMock{iid: "5"}
	g := &mrProviderMock{
		commentsErr: context.DeadlineExceeded,
		diffs: []domain.Diff{
			{NewPath: "new.go", Content: "diff2"},
		},
	}
	o := &ollamaMock{response: `{"issues":[{"file":"new.go","line":10,"severity":"warning","message":"fix it"}]}`}
	r := NewReviewer(c, g, o, zap.NewNop())

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.addedDiscussions) != 1 {
		t.Fatalf("expected 1 discussion, got %d", len(g.addedDiscussions))
	}
}

func TestRunDeletesBotCommentsWhenEnabled(t *testing.T) {
	c := &configMock{iid: "5", deleteBotComments: true}
	g := &mrProviderMock{
		diffs: []domain.Diff{
			{NewPath: "new.go", Content: "diff2"},
		},
	}
	o := &ollamaMock{response: `{"issues":[]}`}
	r := NewReviewer(c, g, o, zap.NewNop())

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.deleteCalls != 1 {
		t.Fatalf("expected delete call, got %d", g.deleteCalls)
	}
}

func TestRunUsesOnlyKnownDiffPathWhenIssueFileIsEmpty(t *testing.T) {
	c := &configMock{iid: "5"}
	g := &mrProviderMock{
		diffs: []domain.Diff{
			{NewPath: "new.go", Content: "diff2"},
		},
	}
	o := &ollamaMock{response: `{"issues":[{"line":10,"severity":"warning","message":"fix it"}]}`}
	r := NewReviewer(c, g, o, zap.NewNop())

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.addedDiscussions) != 1 || g.addedDiscussions[0].file != "new.go" {
		t.Fatalf("unexpected discussions: %+v", g.addedDiscussions)
	}
}

func TestRunSkipsUnknownFilesFromAIResponse(t *testing.T) {
	c := &configMock{iid: "5"}
	g := &mrProviderMock{
		diffs: []domain.Diff{
			{NewPath: "new.go", Content: "diff2"},
		},
	}
	o := &ollamaMock{response: `{"issues":[{"file":"other.go","line":10,"severity":"warning","message":"fix it"}]}`}
	r := NewReviewer(c, g, o, zap.NewNop())

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.addedDiscussions) != 0 {
		t.Fatalf("expected no discussions, got %+v", g.addedDiscussions)
	}
}

func TestBuildCombinedDiffSortsPathsAndDetectsLanguages(t *testing.T) {
	got := buildCombinedDiff([]domain.Diff{
		{NewPath: "b.unknown", Content: "diff-b"},
		{NewPath: "a.go", Content: "diff-a"},
	})

	want := "File: a.go\nLanguage: Go\nDiff:\ndiff-a\n\nFile: b.unknown\nLanguage: Unknown\nDiff:\ndiff-b"
	if got != want {
		t.Fatalf("unexpected combined diff:\n%s", got)
	}
}

func TestRunCancelsInFlightReview(t *testing.T) {
	c := &configMock{iid: "5"}
	g := &mrProviderMock{
		diffs: []domain.Diff{
			{NewPath: "new.go", Content: "diff2"},
		},
	}
	r := NewReviewer(c, g, blockingAIMock{}, zap.NewNop())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := r.Run(ctx)
	if err == nil {
		t.Fatal("expected cancellation error")
	}
	if !strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
		t.Fatalf("expected deadline exceeded error, got %v", err)
	}
	if len(g.addedDiscussions) != 0 {
		t.Fatalf("expected no discussions after cancellation, got %d", len(g.addedDiscussions))
	}
}
