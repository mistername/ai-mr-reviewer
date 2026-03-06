package domain

import (
	"context"
	"time"
)

//go:generate go tool mockery

type MRProviderPort interface {
	GetMergeRequestChanges(ctx context.Context) ([]Diff, error)
	GetExistingComments(ctx context.Context) (map[string][]string, error)
	AddMergeRequestDiscussion(ctx context.Context, file string, line int, note string) error
	DeleteBotCommentsExceptResolved(ctx context.Context) error
}

type AIProviderPort interface {
	ReviewCode(ctx context.Context, diff string) (string, error)
}

type ConfigPort interface {
	GetVCSProvider() string
	GetGitLabURL() string
	GetGitLabToken() string
	GetProjectID() string
	GetMergeRequestIID() string
	GetCommitSHA() string
	GetMergeRequestDiffBaseSHA() string
	GetGitHubToken() string
	GetGitHubOwner() string
	GetGitHubRepo() string
	GetGitHubPRNumber() string
	GetAIProvider() string
	GetOllamaURL() string
	GetOllamaAPIKey() string
	GetOllamaModel() string
	GetOpenAIAPIKey() string
	GetOpenAIBaseURL() string
	GetOpenAIModel() string
	GetAnthropicAuthToken() string
	GetAnthropicBaseURL() string
	GetAnthropicModel() string
	GetMiniMaxAPIKey() string
	GetMiniMaxBaseURL() string
	GetMiniMaxModel() string
	GetCopilotBaseURL() string
	GetCopilotModel() string
	GetRunTimeout() time.Duration
	GetDeleteBotComments() bool
	GetCommentPrefix() string
}
