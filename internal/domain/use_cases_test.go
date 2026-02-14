package domain

import (
	"testing"
)

type mrProviderMock struct{}

func (mrProviderMock) GetMergeRequestChanges() ([]Diff, error) {
	return nil, nil
}
func (mrProviderMock) GetExistingComments() (map[string][]string, error) {
	return map[string][]string{}, nil
}
func (mrProviderMock) AddMergeRequestDiscussion(string, int, string) error { return nil }
func (mrProviderMock) DeleteBotCommentsExceptResolved() error              { return nil }

type ollamaMock struct{}

func (ollamaMock) ReviewCode(string, string, string) (string, error) { return "", nil }

type configMock struct{}

func (configMock) GetVCSProvider() string             { return "gitlab" }
func (configMock) GetGitLabURL() string               { return "" }
func (configMock) GetGitLabToken() string             { return "" }
func (configMock) GetProjectID() string               { return "" }
func (configMock) GetMergeRequestIID() string         { return "" }
func (configMock) GetCommitSHA() string               { return "" }
func (configMock) GetMergeRequestDiffBaseSHA() string { return "" }
func (configMock) GetGitHubToken() string             { return "" }
func (configMock) GetGitHubOwner() string             { return "" }
func (configMock) GetGitHubRepo() string              { return "" }
func (configMock) GetGitHubPRNumber() string          { return "" }
func (configMock) GetAIProvider() string              { return "ollama" }
func (configMock) GetOllamaURL() string               { return "" }
func (configMock) GetOllamaModel() string             { return "" }
func (configMock) GetOpenAIAPIKey() string            { return "" }
func (configMock) GetOpenAIBaseURL() string           { return "" }
func (configMock) GetOpenAIModel() string             { return "" }
func (configMock) GetAnthropicAuthToken() string      { return "" }
func (configMock) GetAnthropicBaseURL() string        { return "" }
func (configMock) GetAnthropicModel() string          { return "" }
func (configMock) GetDeleteBotComments() bool         { return false }
func (configMock) GetCommentPrefix() string           { return "ai-mr-reviewer" }

func TestPortContractsCompile(t *testing.T) {
	var _ MRProviderPort = mrProviderMock{}
	var _ AIProviderPort = ollamaMock{}
	var _ ConfigPort = configMock{}
}
