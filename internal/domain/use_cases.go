package domain

type MRProviderPort interface {
	GetMergeRequestChanges() ([]Diff, error)
	GetExistingComments() (map[string][]string, error)
	AddMergeRequestDiscussion(file string, line int, note string) error
}

type AIProviderPort interface {
	ReviewCode(filePath, diff, language string) (string, error)
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
	GetOllamaModel() string
	GetOpenAIAPIKey() string
	GetOpenAIBaseURL() string
	GetOpenAIModel() string
	GetAnthropicAuthToken() string
	GetAnthropicBaseURL() string
	GetAnthropicModel() string
}
