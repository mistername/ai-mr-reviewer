package main

import (
	"testing"

	"github.com/adlandh/ai-mr-reviewer/internal/domain/mocks"
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

func TestReviewerWithMocks(t *testing.T) {
	cfg := mocks.NewConfigPort(t)
	mrProvider := mocks.NewMRProviderPort(t)
	aiProvider := mocks.NewAIProviderPort(t)
	logger := zap.NewNop()

	reviewer := newReviewer(cfg, mrProvider, aiProvider, logger)
	if reviewer == nil {
		t.Fatal("reviewer is nil")
	}
}
