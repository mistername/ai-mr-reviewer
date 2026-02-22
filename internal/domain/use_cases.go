package domain

import "time"

//go:generate go tool mockery

type MRProviderPort interface {
	GetMergeRequestChanges() ([]Diff, error)
	GetExistingComments() (map[string][]string, error)
	AddMergeRequestDiscussion(file string, line int, note string) error
	DeleteBotCommentsExceptResolved() error
}

type AIProviderPort interface {
	ReviewCode(diff string) (string, error)
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
