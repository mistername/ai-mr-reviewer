# AI MR Reviewer

AI-powered code review bot for GitLab Merge Requests and GitHub Pull Requests supporting multiple AI providers (Ollama, OpenAI, Anthropic, MiniMax, Copilot).

## Overview

This tool automatically reviews code changes in merge requests and pull requests using AI. It analyzes the full MR/PR diff in a single request and posts inline review comments for issues found by the model.

## Quick Start

Local GitHub review with OpenAI:

```bash
export VCS_PROVIDER="github"
export GITHUB_TOKEN="your-github-token"
export GITHUB_OWNER="your-username"
export GITHUB_REPO="your-repo"
export GITHUB_PR_NUMBER="123"
export CI_COMMIT_SHA="abc123"
export AI_PROVIDER="openai"
export OPENAI_API_KEY="your-openai-key"

go run ./main.go
```

Local GitLab review with Ollama:

```bash
export VCS_PROVIDER="gitlab"
export GITLAB_TOKEN="your-gitlab-token"
export CI_PROJECT_ID="123"
export CI_MERGE_REQUEST_IID="1"
export CI_COMMIT_SHA="abc123"
export CI_MERGE_REQUEST_DIFF_BASE_SHA="def456"
export AI_PROVIDER="ollama"

go run ./main.go
```

## Features

- Supports both GitLab Merge Requests and GitHub Pull Requests
- AI-powered analysis (Ollama, OpenAI, Anthropic, MiniMax, Copilot)
- Inline comments with line numbers
- Configurable comment prefix and cleanup of previous bot comments
- Hexagonal architecture
- Uber FX dependency injection
- Zap structured logging

## Requirements

- Go 1.26+
- GitLab API token or GitHub token
- One of:
  - Ollama server with a model (e.g., llama3.2)
  - OpenAI API key
  - Anthropic API key
  - MiniMax API key
  - GitHub token (for Copilot/GitHub Models)

## Configuration

### VCS Provider Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `VCS_PROVIDER` | Version control system: `gitlab` or `github` | `gitlab` |

### GitLab Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `GITLAB_TOKEN` | GitLab API token with `api` and `write_discussion` scope | - |
| `CI_PROJECT_ID` | GitLab project ID | - |
| `CI_MERGE_REQUEST_IID` | Merge Request IID | - |
| `CI_COMMIT_SHA` | Commit SHA | - |
| `CI_MERGE_REQUEST_DIFF_BASE_SHA` | Base diff SHA | - |
| `CI_SERVER_URL` | GitLab server URL | `https://gitlab.com` |

### GitHub Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `GITHUB_TOKEN` | GitHub token used for PR access and, for `copilot`, GitHub Models auth | - |
| `GITHUB_OWNER` | GitHub repository owner | - |
| `GITHUB_REPO` | GitHub repository name | - |
| `GITHUB_PR_NUMBER` | Pull Request number | - |
| `CI_COMMIT_SHA` | Commit SHA | - |

### AI Provider Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `AI_PROVIDER` | AI provider: `ollama`, `openai`, `anthropic`, `minimax`, or `copilot` | `ollama` |

### Review Behavior

| Variable | Description | Default |
|----------|-------------|---------|
| `COMMENT_PREFIX` | Prefix used for every bot comment (`<prefix>:`) | `ai-mr-reviewer` |
| `DELETE_BOT_COMMENTS` | Delete previous unresolved bot comments before new run | `true` |
| `RUN_TIMEOUT` | Max duration of a single review run | `10m` |

### Ollama Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `OLLAMA_URL` | Ollama server URL | `http://localhost:11434` |
| `OLLAMA_API_KEY` | Ollama API key for protected/cloud endpoints | - |
| `OLLAMA_MODEL` | Ollama model name | `llama3.2` |

### OpenAI Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `OPENAI_API_KEY` | OpenAI API key | - |
| `OPENAI_BASE_URL` | OpenAI API base URL | `https://api.openai.com/v1` |
| `OPENAI_MODEL` | OpenAI model name | `GPT-5.2-Codex` |

### Anthropic Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `ANTHROPIC_AUTH_TOKEN` | Anthropic auth token | - |
| `ANTHROPIC_BASE_URL` | Anthropic API base URL | `https://api.anthropic.com/v1/` |
| `ANTHROPIC_MODEL` | Anthropic model name | `claude-sonnet-4-20250514` |

### MiniMax Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `MINIMAX_API_KEY` | MiniMax API key | - |
| `MINIMAX_BASE_URL` | MiniMax API base URL | `https://api.minimax.io/v1` |
| `MINIMAX_MODEL` | MiniMax model name | `MiniMax-M2.5` |

### Copilot Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `GITHUB_TOKEN` | GitHub token used as a bearer token for GitHub Models | - |
| `COPILOT_BASE_URL` | GitHub Models-compatible base URL | `https://models.github.ai/inference` |
| `COPILOT_MODEL` | GitHub Models-compatible model name | `openai/gpt-4.1` |

## Usage

The sections below show full environment examples for the supported provider combinations.

### Running locally (GitLab + Ollama)

```bash
export VCS_PROVIDER="gitlab"
export GITLAB_TOKEN="your-gitlab-token"
export CI_PROJECT_ID="123"
export CI_MERGE_REQUEST_IID="1"
export CI_COMMIT_SHA="abc123"
export CI_MERGE_REQUEST_DIFF_BASE_SHA="def456"

go run ./main.go
```

### Running locally (GitHub + OpenAI)

```bash
export VCS_PROVIDER="github"
export GITHUB_TOKEN="your-github-token"
export GITHUB_OWNER="your-username"
export GITHUB_REPO="your-repo"
export GITHUB_PR_NUMBER="123"
export CI_COMMIT_SHA="abc123"
export AI_PROVIDER="openai"
export OPENAI_API_KEY="your-openai-key"

go run ./main.go
```

### Running locally (GitHub + MiniMax)

```bash
export VCS_PROVIDER="github"
export GITHUB_TOKEN="your-github-token"
export GITHUB_OWNER="your-username"
export GITHUB_REPO="your-repo"
export GITHUB_PR_NUMBER="123"
export CI_COMMIT_SHA="abc123"
export AI_PROVIDER="minimax"
export MINIMAX_API_KEY="your-minimax-key"

go run ./main.go
```

### Running locally (GitHub + Copilot)

```bash
export VCS_PROVIDER="github"
export GITHUB_TOKEN="your-github-token"
export GITHUB_OWNER="your-username"
export GITHUB_REPO="your-repo"
export GITHUB_PR_NUMBER="123"
export CI_COMMIT_SHA="abc123"
export AI_PROVIDER="copilot"
export COPILOT_BASE_URL="https://models.github.ai/inference" # optional
export COPILOT_MODEL="openai/gpt-4.1"                         # optional

go run ./main.go
```

### Running locally (GitLab + Anthropic)

```bash
export VCS_PROVIDER="gitlab"
export GITLAB_TOKEN="your-gitlab-token"
export CI_PROJECT_ID="123"
export CI_MERGE_REQUEST_IID="1"
export CI_COMMIT_SHA="abc123"
export CI_MERGE_REQUEST_DIFF_BASE_SHA="def456"
export AI_PROVIDER="anthropic"
export ANTHROPIC_AUTH_TOKEN="your-anthropic-token"

go run ./main.go
```

### Running in GitLab CI

```yaml
stages:
  - review

ai-code-review:
  stage: review
  image: golang:1.26-alpine
  services:
    - name: ollama/ollama:latest
      alias: ollama
  variables:
    VCS_PROVIDER: "gitlab"
    CI_SERVER_URL: $CI_SERVER_URL
    GITLAB_TOKEN: $GITLAB_TOKEN
    CI_PROJECT_ID: $CI_PROJECT_ID
    CI_MERGE_REQUEST_IID: $CI_MERGE_REQUEST_IID
    CI_COMMIT_SHA: $CI_COMMIT_SHA
    CI_MERGE_REQUEST_DIFF_BASE_SHA: $CI_MERGE_REQUEST_DIFF_BASE_SHA
    AI_PROVIDER: "ollama"
    OLLAMA_URL: "http://ollama:11434"
    OLLAMA_MODEL: "llama3.2"
  before_script:
    - go install github.com/adlandh/ai-mr-reviewer@latest
  script:
    - ai-mr-reviewer
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    - if: $CI_MERGE_REQUEST_IID
  allow_failure: true
```

### Running in GitHub Actions

Base workflow:

```yaml
name: AI Code Review

on:
  pull_request:

jobs:
  review:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      models: read
      pull-requests: write
    steps:
      - uses: actions/checkout@v6
      - name: Set up Go
        uses: actions/setup-go@v6
        with:
          go-version: '1.26'
      - name: Install and Run Review
        env:
          VCS_PROVIDER: github
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          GITHUB_OWNER: ${{ github.repository_owner }}
          GITHUB_REPO: ${{ github.event.repository.name }}
          GITHUB_PR_NUMBER: ${{ github.event.pull_request.number }}
          CI_COMMIT_SHA: ${{ github.sha }}
          DELETE_BOT_COMMENTS: true
        run: go install github.com/adlandh/ai-mr-reviewer@latest && ai-mr-reviewer
```

`models: read` is only required when `AI_PROVIDER=copilot`. It is safe to leave enabled in the base workflow if you switch providers.

OpenAI example (`env` additions):

```yaml
          AI_PROVIDER: openai
          OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
          OPENAI_MODEL: GPT-5.2-Codex
```

Ollama example (`env` additions):

```yaml
          AI_PROVIDER: ollama
          OLLAMA_URL: http://localhost:11434
          OLLAMA_API_KEY: ${{ secrets.OLLAMA_API_KEY }} # optional for cloud/protected endpoint
          OLLAMA_MODEL: llama3.2
```
(`steps` additions for local runner model start):
```yaml
        - name: Run model
          uses: ai-action/ollama-action@v2
          with:
            model: llama3.2
```
For Ollama cloud/remote endpoints, set `OLLAMA_URL` + `OLLAMA_API_KEY` and skip the `Run model` step.

Anthropic example (`env` additions):

```yaml
          AI_PROVIDER: anthropic
          ANTHROPIC_AUTH_TOKEN: ${{ secrets.ANTHROPIC_AUTH_TOKEN }}
          ANTHROPIC_MODEL: claude-sonnet-4-20250514
```

MiniMax example (`env` additions):

```yaml
          AI_PROVIDER: minimax
          MINIMAX_API_KEY: ${{ secrets.MINIMAX_API_KEY }}
          MINIMAX_MODEL: MiniMax-M2.5
```

Copilot example (`env` additions):

```yaml
          AI_PROVIDER: copilot
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          COPILOT_BASE_URL: https://models.github.ai/inference
          COPILOT_MODEL: openai/gpt-4.1
```

Copilot/GitHub Models workflow permission addition:

```yaml
    permissions:
      contents: read
      models: read
      pull-requests: write
```

## GitLab Token Setup

1. Go to GitLab → Settings → Access Tokens
2. Create a new token with:
   - `api` scope
   - `write_discussion` scope
3. Add the token as a CI/CD variable `GITLAB_TOKEN` (masked)

## GitHub Token Setup

For local runs with a personal token:

1. Go to GitHub → Settings → Developer settings → Personal access tokens.
2. Create a token with:
   - `repo` scope for private repositories, or `public_repo` for public repositories
   - access to GitHub Models if you use `AI_PROVIDER=copilot`
3. Export it as `GITHUB_TOKEN`.

For GitHub Actions:

1. Use `${{ secrets.GITHUB_TOKEN }}` for repository access.
2. If you use `AI_PROVIDER=copilot`, grant the workflow `models: read` so the workflow token can call the GitHub Models inference API.

## Development

### Build

```bash
go install .
```

### Test

```bash
go test ./...
```

### Lint

```bash
golangci-lint run
```

## Architecture

```
main.go              # Application entry point
internal/
  application/         # Use cases
    reviewer.go       # Review orchestration
  domain/
    entity.go        # Domain entities
    use_cases.go     # Port interfaces
  adapters/
    gitlab/         # GitLab API client
    github/         # GitHub API client
    ai/             # AI provider adapters (Ollama, OpenAI, Anthropic, MiniMax, Copilot)
    config/         # Configuration
```

## Notes

- Reviews are sent as a single combined diff to the configured model. Very large PRs or MRs can hit provider context limits.
- `RUN_TIMEOUT` bounds the full review run, including VCS API calls and model inference.
- On GitHub, inline comments may fall back to a general PR comment if GitHub rejects the exact line placement.

## License

MIT
