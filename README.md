# AI MR Reviewer

AI-powered code review bot for GitLab Merge Requests and GitHub Pull Requests supporting multiple AI providers (Ollama, OpenAI, Anthropic, MiniMax, Copilot).

## Overview

This tool automatically reviews code changes in Merge Requests using AI. It analyzes the full MR/PR diff in a single request and adds inline comments with issues found by the AI model.

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
- GitLab API Token or GitHub Personal Access Token
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
| `GITHUB_TOKEN` | GitHub Personal Access Token | - |
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
| `GITHUB_TOKEN` | GitHub token used as bearer token for Copilot | - |
| `COPILOT_BASE_URL` | Copilot-compatible base URL | `https://models.github.ai/inference` |
| `COPILOT_MODEL` | Copilot-compatible model name | `openai/gpt-4.1` |

## Usage

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

OpenAI example (`env` additions):

```yaml
          AI_PROVIDER: openai
          OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
          OPENAI_MODEL: gpt-5-codex
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

## GitLab Token Setup

1. Go to GitLab → Settings → Access Tokens
2. Create a new token with:
   - `api` scope
   - `write_discussion` scope
3. Add the token as a CI/CD variable `GITLAB_TOKEN` (masked)

## GitHub Token Setup

1. Go to GitHub → Settings → Developer settings → Personal access tokens
2. Create a new token with:
   - `repo` scope (for private repos) or `public_repo` scope (for public repos)
3. Add the token as a secret in your repository

## Development

### Build

```bash
go install github.com/adlandh/ai-mr-reviewer@latest
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

## License

MIT
