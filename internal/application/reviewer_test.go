package application

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/adlandh/ai-mr-reviewer/internal/domain"
	"github.com/adlandh/ai-mr-reviewer/internal/domain/mocks"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type addedDiscussion struct {
	file string
	line int
	body string
}

const (
	warningComment           = "ai-mr-reviewer:**WARNING**: fix it"
	expectedOneDiscussionFmt = "expected 1 discussion, got %d"
)

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
	c := newConfigMock(t, false)
	g := mocks.NewMRProviderPort(t)
	o := mocks.NewAIProviderPort(t)
	added := make([]addedDiscussion, 0, 1)

	g.EXPECT().GetExistingComments(mock.Anything).Return(map[string][]string{
		"already.go:1": {warningComment},
	}, nil)
	g.EXPECT().GetMergeRequestChanges(mock.Anything).Return([]domain.Diff{
		{NewPath: "already.go", Content: "diff1"},
		{NewPath: "new.go", Content: "diff2"},
	}, nil)
	o.EXPECT().ReviewCode(mock.Anything, mock.Anything).Return(`{"issues":[{"file":"new.go","line":10,"severity":"warning","message":"fix it"}]}`, nil)
	g.EXPECT().AddMergeRequestDiscussion(mock.Anything, "new.go", 10, warningComment).
		Run(func(_ context.Context, file string, line int, body string) {
			added = append(added, addedDiscussion{file: file, line: line, body: body})
		}).
		Return(nil)

	r := NewReviewer(c, g, o, zap.NewNop())

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(added) != 1 {
		t.Fatalf(expectedOneDiscussionFmt, len(added))
	}
	if added[0].file != "new.go" || added[0].line != 10 {
		t.Fatalf("unexpected discussion: %+v", added[0])
	}
}

func TestRunReviewsNewDiffsNoFilter(t *testing.T) {
	c := newConfigMock(t, false)
	g := mocks.NewMRProviderPort(t)
	o := mocks.NewAIProviderPort(t)
	callCount := 0

	g.EXPECT().GetExistingComments(mock.Anything).Return(map[string][]string{}, nil)
	g.EXPECT().GetMergeRequestChanges(mock.Anything).Return([]domain.Diff{{NewPath: "new.go", Content: "diff2"}}, nil)
	o.EXPECT().ReviewCode(mock.Anything, mock.Anything).Return(`{"issues":[{"file":"new.go","line":10,"severity":"warning","message":"fix it"}]}`, nil)
	g.EXPECT().AddMergeRequestDiscussion(mock.Anything, "new.go", 10, warningComment).
		Run(func(context.Context, string, int, string) { callCount++ }).
		Return(nil)

	r := NewReviewer(c, g, o, zap.NewNop())

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 1 {
		t.Fatalf(expectedOneDiscussionFmt, callCount)
	}
}

func TestRunContinuesWhenExistingCommentsFail(t *testing.T) {
	c := newConfigMock(t, false)
	g := mocks.NewMRProviderPort(t)
	o := mocks.NewAIProviderPort(t)
	callCount := 0

	g.EXPECT().GetExistingComments(mock.Anything).Return(nil, context.DeadlineExceeded)
	g.EXPECT().GetMergeRequestChanges(mock.Anything).Return([]domain.Diff{{NewPath: "new.go", Content: "diff2"}}, nil)
	o.EXPECT().ReviewCode(mock.Anything, mock.Anything).Return(`{"issues":[{"file":"new.go","line":10,"severity":"warning","message":"fix it"}]}`, nil)
	g.EXPECT().AddMergeRequestDiscussion(mock.Anything, "new.go", 10, warningComment).
		Run(func(context.Context, string, int, string) { callCount++ }).
		Return(nil)

	r := NewReviewer(c, g, o, zap.NewNop())

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 1 {
		t.Fatalf(expectedOneDiscussionFmt, callCount)
	}
}

func TestRunDeletesBotCommentsWhenEnabled(t *testing.T) {
	c := newConfigMock(t, true)
	g := mocks.NewMRProviderPort(t)
	o := mocks.NewAIProviderPort(t)
	deleteCalls := 0

	g.EXPECT().DeleteBotCommentsExceptResolved(mock.Anything).
		Run(func(context.Context) { deleteCalls++ }).
		Return(nil)
	g.EXPECT().GetExistingComments(mock.Anything).Return(map[string][]string{}, nil)
	g.EXPECT().GetMergeRequestChanges(mock.Anything).Return([]domain.Diff{{NewPath: "new.go", Content: "diff2"}}, nil)
	o.EXPECT().ReviewCode(mock.Anything, mock.Anything).Return(`{"issues":[]}`, nil)

	r := NewReviewer(c, g, o, zap.NewNop())

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deleteCalls != 1 {
		t.Fatalf("expected delete call, got %d", deleteCalls)
	}
}

func TestRunUsesOnlyKnownDiffPathWhenIssueFileIsEmpty(t *testing.T) {
	c := newConfigMock(t, false)
	g := mocks.NewMRProviderPort(t)
	o := mocks.NewAIProviderPort(t)
	added := make([]addedDiscussion, 0, 1)

	g.EXPECT().GetExistingComments(mock.Anything).Return(map[string][]string{}, nil)
	g.EXPECT().GetMergeRequestChanges(mock.Anything).Return([]domain.Diff{{NewPath: "new.go", Content: "diff2"}}, nil)
	o.EXPECT().ReviewCode(mock.Anything, mock.Anything).Return(`{"issues":[{"line":10,"severity":"warning","message":"fix it"}]}`, nil)
	g.EXPECT().AddMergeRequestDiscussion(mock.Anything, "new.go", 10, warningComment).
		Run(func(_ context.Context, file string, line int, body string) {
			added = append(added, addedDiscussion{file: file, line: line, body: body})
		}).
		Return(nil)

	r := NewReviewer(c, g, o, zap.NewNop())

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(added) != 1 || added[0].file != "new.go" {
		t.Fatalf("unexpected discussions: %+v", added)
	}
}

func TestRunSkipsUnknownFilesFromAIResponse(t *testing.T) {
	c := newConfigMock(t, false)
	g := mocks.NewMRProviderPort(t)
	o := mocks.NewAIProviderPort(t)

	g.EXPECT().GetExistingComments(mock.Anything).Return(map[string][]string{}, nil)
	g.EXPECT().GetMergeRequestChanges(mock.Anything).Return([]domain.Diff{{NewPath: "new.go", Content: "diff2"}}, nil)
	o.EXPECT().ReviewCode(mock.Anything, mock.Anything).Return(`{"issues":[{"file":"other.go","line":10,"severity":"warning","message":"fix it"}]}`, nil)

	r := NewReviewer(c, g, o, zap.NewNop())

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
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
	c := newConfigMock(t, false)
	g := mocks.NewMRProviderPort(t)
	ai := mocks.NewAIProviderPort(t)

	g.EXPECT().GetExistingComments(mock.Anything).Return(map[string][]string{}, nil)
	g.EXPECT().GetMergeRequestChanges(mock.Anything).Return([]domain.Diff{{NewPath: "new.go", Content: "diff2"}}, nil)
	ai.EXPECT().ReviewCode(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, _ string) (string, error) {
		<-ctx.Done()
		return "", ctx.Err()
	})

	r := NewReviewer(c, g, ai, zap.NewNop())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := r.Run(ctx)
	if err == nil {
		t.Fatal("expected cancellation error")
	}
	if !strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
		t.Fatalf("expected deadline exceeded error, got %v", err)
	}
}

func newConfigMock(t *testing.T, deleteBotComments bool) *mocks.ConfigPort {
	t.Helper()

	config := mocks.NewConfigPort(t)
	config.On("GetDeleteBotComments").Return(deleteBotComments)
	config.On("GetCommentPrefix").Return("ai-mr-reviewer")

	return config
}
