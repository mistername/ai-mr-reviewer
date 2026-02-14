package config

import (
	"fmt"

	"github.com/adlandh/ai-mr-reviewer/internal/domain"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	VCSProvider             string `env:"VCS_PROVIDER,notEmpty" envDefault:"gitlab"`
	GitLabURL               string `env:"CI_SERVER_URL" envDefault:"https://gitlab.com"`
	GitLabToken             string `env:"GITLAB_TOKEN"`
	ProjectID               string `env:"CI_PROJECT_ID"`
	MergeRequestIID         string `env:"CI_MERGE_REQUEST_IID"`
	CommitSHA               string `env:"CI_COMMIT_SHA"`
	MergeRequestDiffBaseSHA string `env:"CI_MERGE_REQUEST_DIFF_BASE_SHA"`
	GitHubToken             string `env:"GITHUB_TOKEN"`
	GitHubOwner             string `env:"GITHUB_OWNER"`
	GitHubRepo              string `env:"GITHUB_REPO"`
	GitHubPRNumber          string `env:"GITHUB_PR_NUMBER"`
	AIProvider              string `env:"AI_PROVIDER,notEmpty" envDefault:"ollama"`
	OllamaURL               string `env:"OLLAMA_URL" envDefault:"http://localhost:11434"`
	OllamaModel             string `env:"OLLAMA_MODEL" envDefault:"llama3.2"`
	OpenAIAPIKey            string `env:"OPENAI_API_KEY"`
	OpenAIBaseURL           string `env:"OPENAI_BASE_URL" envDefault:"https://api.openai.com/v1"`
	OpenAIModel             string `env:"OPENAI_MODEL" envDefault:"gpt-4o"`
	AnthropicAuthToken      string `env:"ANTHROPIC_AUTH_TOKEN"`
	AnthropicBaseURL        string `env:"ANTHROPIC_BASE_URL" envDefault:"https://api.anthropic.com/v1/"`
	AnthropicModel          string `env:"ANTHROPIC_MODEL" envDefault:"claude-sonnet-4-20250514"`
	DeleteBotComments       bool   `env:"DELETE_BOT_COMMENTS"`
}

func New() (*Config, error) {
	cfg, err := env.ParseAsWithOptions[Config](env.Options{RequiredIfNoDef: false})
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) GetGitLabURL() string {
	return c.GitLabURL
}

func (c *Config) GetGitLabToken() string {
	return c.GitLabToken
}

func (c *Config) GetProjectID() string {
	return c.ProjectID
}

func (c *Config) GetMergeRequestIID() string {
	return c.MergeRequestIID
}

func (c *Config) GetCommitSHA() string {
	return c.CommitSHA
}

func (c *Config) GetMergeRequestDiffBaseSHA() string {
	return c.MergeRequestDiffBaseSHA
}

func (c *Config) GetVCSProvider() string {
	return c.VCSProvider
}

func (c *Config) GetGitHubToken() string {
	return c.GitHubToken
}

func (c *Config) GetGitHubOwner() string {
	return c.GitHubOwner
}

func (c *Config) GetGitHubRepo() string {
	return c.GitHubRepo
}

func (c *Config) GetGitHubPRNumber() string {
	return c.GitHubPRNumber
}

func (c *Config) GetAIProvider() string {
	return c.AIProvider
}

func (c *Config) GetOllamaURL() string {
	return c.OllamaURL
}

func (c *Config) GetOllamaModel() string {
	return c.OllamaModel
}

func (c *Config) GetOpenAIAPIKey() string {
	return c.OpenAIAPIKey
}

func (c *Config) GetOpenAIBaseURL() string {
	return c.OpenAIBaseURL
}

func (c *Config) GetOpenAIModel() string {
	return c.OpenAIModel
}

func (c *Config) GetAnthropicAuthToken() string {
	return c.AnthropicAuthToken
}

func (c *Config) GetAnthropicBaseURL() string {
	return c.AnthropicBaseURL
}

func (c *Config) GetAnthropicModel() string {
	return c.AnthropicModel
}

func (c *Config) GetDeleteBotComments() bool {
	return c.DeleteBotComments
}

var _ domain.ConfigPort = (*Config)(nil)
