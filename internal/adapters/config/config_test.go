package config

import (
	"testing"
	"time"
)

func TestNewConfigWithDefaults(t *testing.T) {
	t.Setenv("VCS_PROVIDER", "gitlab")
	t.Setenv("GITLAB_TOKEN", "token")
	t.Setenv("CI_PROJECT_ID", "123")
	t.Setenv("CI_MERGE_REQUEST_IID", "1")
	t.Setenv("CI_COMMIT_SHA", "abc123")
	t.Setenv("CI_MERGE_REQUEST_DIFF_BASE_SHA", "def456")
	t.Setenv("CI_SERVER_URL", "")
	t.Setenv("OLLAMA_URL", "")
	t.Setenv("OLLAMA_MODEL", "")

	cfg, err := New()
	if err != nil {
		t.Fatalf("expected config, got error: %v", err)
	}

	if cfg.GetVCSProvider() != "gitlab" {
		t.Fatalf("unexpected VCS provider: %s", cfg.GetVCSProvider())
	}
	if cfg.GetGitLabURL() != "https://gitlab.com" {
		t.Fatalf("unexpected GitLab URL: %s", cfg.GetGitLabURL())
	}
	if cfg.GetOllamaURL() != "http://localhost:11434" {
		t.Fatalf("unexpected Ollama URL: %s", cfg.GetOllamaURL())
	}
	if cfg.GetOllamaModel() != "llama3.2" {
		t.Fatalf("unexpected Ollama model: %s", cfg.GetOllamaModel())
	}
	if cfg.GetMergeRequestIID() != "1" {
		t.Fatalf("unexpected MergeRequestIID: %s", cfg.GetMergeRequestIID())
	}
	if cfg.GetCommitSHA() != "abc123" {
		t.Fatalf("unexpected CommitSHA: %s", cfg.GetCommitSHA())
	}
	if cfg.GetMergeRequestDiffBaseSHA() != "def456" {
		t.Fatalf("unexpected MergeRequestDiffBaseSHA: %s", cfg.GetMergeRequestDiffBaseSHA())
	}
	if cfg.GetRunTimeout() != 10*time.Minute {
		t.Fatalf("unexpected RunTimeout: %s", cfg.GetRunTimeout())
	}
}

func TestNewConfigMissingRequired(t *testing.T) {
	t.Setenv("VCS_PROVIDER", "")
	t.Setenv("GITLAB_TOKEN", "")
	t.Setenv("CI_PROJECT_ID", "")
	t.Setenv("CI_MERGE_REQUEST_IID", "")
	t.Setenv("CI_COMMIT_SHA", "")
	t.Setenv("CI_MERGE_REQUEST_DIFF_BASE_SHA", "")

	cfg, err := New()
	if err != nil {
		t.Fatalf("expected config without error, got: %v", err)
	}

	if cfg.GetVCSProvider() != "gitlab" {
		t.Fatalf("unexpected VCS provider: %s", cfg.GetVCSProvider())
	}
}
