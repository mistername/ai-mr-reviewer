package main

import (
	"context"
	"fmt"
	"time"

	aiadapter "github.com/mistername/ai-mr-reviewer/internal/adapters/ai"
	configadapter "github.com/mistername/ai-mr-reviewer/internal/adapters/config"
	githubadapter "github.com/mistername/ai-mr-reviewer/internal/adapters/github"
	gitlabadapter "github.com/mistername/ai-mr-reviewer/internal/adapters/gitlab"
	"github.com/mistername/ai-mr-reviewer/internal/application"
	"github.com/mistername/ai-mr-reviewer/internal/domain"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = logger.Sync()
	}()

	app := fx.New(
		fx.Supply(logger),
		fx.WithLogger(func(l *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: l}
		}),
		fx.Provide(
			newConfig,
			newVCSClient,
			newAIProvider,
			newReviewer,
		),
		fx.Invoke(runReview),
	)

	startCtx, cancelStart := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelStart()

	err = app.Start(startCtx)
	if err != nil {
		logger.Fatal("start app", zap.Error(err))
	}

	stopCtx, cancelStop := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelStop()

	if err := app.Stop(stopCtx); err != nil {
		logger.Fatal("stop app", zap.Error(err))
	}
}

func newConfig() (*domain.Config, error) {
	cfg, err := configadapter.New()
	if err != nil {
		return nil, fmt.Errorf("create config: %w", err)
	}

	return cfg, nil
}

func newVCSClient(cfg *domain.Config) (domain.MRProviderPort, error) {
	provider := cfg.VCS.Provider

	switch provider {
	case "github":
		return newGitHubClient(cfg)
	case "gitlab":
		return newGitLabClient(cfg)
	default:
		return nil, fmt.Errorf("unsupported VCS provider: %s (supported: gitlab, github)", provider)
	}
}

func newGitHubClient(cfg *domain.Config) (domain.MRProviderPort, error) {
	github := cfg.VCS.GitHub

	if github.Token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN is required for github provider")
	}

	if github.Owner == "" || github.Repo == "" {
		return nil, fmt.Errorf("GITHUB_OWNER and GITHUB_REPO are required for github provider")
	}

	if github.PRNumber == "" {
		return nil, fmt.Errorf("GITHUB_PR_NUMBER is required for github provider")
	}

	client, err := githubadapter.NewClient(
		github.Token,
		github.Owner,
		github.Repo,
		github.PRNumber,
		github.CommitSHA,
		cfg.Runtime.CommentPrefix,
	)
	if err != nil {
		return nil, fmt.Errorf("create GitHub client: %w", err)
	}

	return client, nil
}

func newGitLabClient(cfg *domain.Config) (domain.MRProviderPort, error) {
	gitlab := cfg.VCS.GitLab

	if gitlab.Token == "" {
		return nil, fmt.Errorf("GITLAB_TOKEN is required for gitlab provider")
	}

	if gitlab.ProjectID == "" {
		return nil, fmt.Errorf("CI_PROJECT_ID is required for gitlab provider")
	}

	iidStr := gitlab.MergeRequestIID
	if iidStr == "" {
		return nil, fmt.Errorf("CI_MERGE_REQUEST_IID is required for gitlab provider")
	}

	var iid int
	if _, err := fmt.Sscanf(iidStr, "%d", &iid); err != nil {
		return nil, fmt.Errorf("parse CI_MERGE_REQUEST_IID: %w", err)
	}

	client, err := gitlabadapter.NewClient(gitlab, cfg.Runtime, iid)
	if err != nil {
		return nil, fmt.Errorf("create GitLab client: %w", err)
	}

	return client, nil
}

func newAIProvider(cfg *domain.Config) (domain.AIProviderPort, error) {
	provider, err := aiadapter.NewAIProvider(cfg.AI)
	if err != nil {
		return nil, fmt.Errorf("create AI provider: %w", err)
	}

	return provider, nil
}

func newReviewer(cfg *domain.Config, mrProvider domain.MRProviderPort, aiProvider domain.AIProviderPort, logger *zap.Logger) *application.Reviewer {
	return application.NewReviewer(cfg.Runtime, mrProvider, aiProvider, logger)
}

func runReview(reviewer *application.Reviewer, logger *zap.Logger, cfg *domain.Config) error {
	logger.Info("starting MR Reviewer")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Runtime.RunTimeout)
	defer cancel()

	if err := reviewer.Run(ctx); err != nil {
		logger.Error("review failed", zap.Error(err))
		return fmt.Errorf("run review: %w", err)
	}

	logger.Info("review completed")

	return nil
}
