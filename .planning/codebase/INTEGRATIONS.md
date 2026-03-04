# External Integrations

**Analysis Date:** 2026-03-04

## APIs & External Services

**Version Control Providers:**
- GitHub API - Pull Request review automation
  - SDK/Client: `github.com/google/go-github/v82`
  - Auth: `GITHUB_TOKEN` (Personal Access Token with `repo` or `public_repo` scopes)
  - Endpoints: Files (`ListFiles`), Comments (`ListComments`, `CreateCommentReply`), PR data
  - Implementation: `internal/adapters/github/client.go`

- GitLab API - Merge Request review automation
  - SDK/Client: `gitlab.com/gitlab-org/api/client-go v1.46.0`
  - Auth: `GITLAB_TOKEN` (API token with `api` and `write_discussion` scopes)
  - Endpoints: MR Diffs (`ListMergeRequestDiffs`), Notes (`ListMergeRequestNotes`, `CreateNote`)
  - Implementation: `internal/adapters/gitlab/client.go`

**AI Providers:**
- Ollama - Local/self-hosted LLM
  - SDK/Client: Direct HTTP client (`net/http`)
  - Auth: `OLLAMA_API_KEY` (optional, for cloud endpoints)
  - Base URL: `OLLAMA_URL` (default: `http://localhost:11434`)
  - Endpoint: `/api/chat`
  - Model: `OLLAMA_MODEL` (default: `llama3.2`)
  - Implementation: `internal/adapters/ai/ollama.go`

- OpenAI API - GPT models for code review
  - SDK/Client: Direct HTTP client via shared completion handler
  - Auth: `OPENAI_API_KEY`
  - Base URL: `OPENAI_BASE_URL` (default: `https://api.openai.com/v1`)
  - Endpoint: `/chat/completions`
  - Model: `OPENAI_MODEL` (default: `GPT-5.2-Codex`)
  - Implementation: `internal/adapters/ai/openai.go`

- Anthropic API - Claude models for code review
  - SDK/Client: Direct HTTP client
  - Auth: `ANTHROPIC_AUTH_TOKEN`
  - Base URL: `ANTHROPIC_BASE_URL` (default: `https://api.anthropic.com/v1/`)
  - Endpoint: `/messages`
  - Model: `ANTHROPIC_MODEL` (default: `claude-sonnet-4-20250514`)
  - Max Tokens: 4096
  - Implementation: `internal/adapters/ai/anthropic.go`

- MiniMax API - Chinese LLM provider
  - SDK/Client: Direct HTTP client via shared completion handler
  - Auth: `MINIMAX_API_KEY`
  - Base URL: `MINIMAX_BASE_URL` (default: `https://api.minimax.io/v1`)
  - Endpoint: `/chat/completions`
  - Model: `MINIMAX_MODEL` (default: `MiniMax-M2.5`)
  - Implementation: `internal/adapters/ai/minimax.go`

- GitHub Copilot (OpenAI Compatible) - GitHub Models API
  - SDK/Client: Direct HTTP client via shared completion handler
  - Auth: `GITHUB_TOKEN` (used as Bearer token)
  - Base URL: `COPILOT_BASE_URL` (default: `https://models.github.ai/inference`)
  - Endpoint: `/chat/completions`
  - Model: `COPILOT_MODEL` (default: `openai/gpt-4.1`)
  - Implementation: `internal/adapters/ai/copilot.go`

## Data Storage

**Databases:**
- Not applicable - This is a stateless code review bot

**File Storage:**
- Not used - All data is transient (read MR/PR diffs, analyze, write comments)

**Caching:**
- None - Each review run fetches fresh data from VCS providers

## Authentication & Identity

**Auth Approach:**
- Token-based authentication for all external services
- No database-backed user management
- Credentials passed via environment variables only

**VCS Authentication:**
- GitHub: Personal Access Token via `GITHUB_TOKEN` env var
  - Scopes required: `repo` (private) or `public_repo` (public only)
  - Passed to: `github.com/google/go-github/v82` client via `WithAuthToken()`

- GitLab: API Token via `GITLAB_TOKEN` env var
  - Scopes required: `api` and `write_discussion`
  - Passed to: `gitlab.com/gitlab-org/api/client-go` client constructor

**AI Provider Authentication:**
- Ollama: Optional `OLLAMA_API_KEY` for protected endpoints (Bearer token)
- OpenAI: `OPENAI_API_KEY` (Bearer token in `Authorization` header)
- Anthropic: `ANTHROPIC_AUTH_TOKEN` (Bearer token in `x-api-key` header)
- MiniMax: `MINIMAX_API_KEY` (Bearer token in `Authorization` header)
- Copilot: `GITHUB_TOKEN` (Bearer token in `Authorization` header)

## Monitoring & Observability

**Error Tracking:**
- Not integrated - Errors logged via Zap structured logging to stderr

**Logs:**
- Zap v1.27.1 structured JSON logging
- Configuration: Production logger in `main.go` line 20
- Output: Logs written to stdout/stderr
- Levels used: Info, Error, Fatal
- Context included: Component names, error stack traces

## CI/CD & Deployment

**Hosting:**
- Self-hosted: Single compiled Go binary
- GitHub Actions: GitHub-hosted runners (see README examples)
- GitLab CI: GitLab CI/CD pipeline with docker runner
- Can run locally with environment variables

**CI Pipeline:**
- GitHub Actions: Standard ubuntu-latest runner with Go 1.26
- GitLab CI: golang:1.26-alpine image with optional Ollama service
- Both use environment variables for configuration

**Deployment Method:**
- Build: `go install github.com/adlandh/ai-mr-reviewer@latest` or `go run ./main.go`
- Package: Single binary executable
- No Docker image in official repository (but examples show using golang:1.26-alpine)

## Environment Configuration

**Required env vars (shared):**
- `CI_COMMIT_SHA` - Current commit SHA for code review
- `COMMENT_PREFIX` - Prefix for bot comments (default: `ai-mr-reviewer`)
- `VCS_PROVIDER` - `github` or `gitlab` (default: `gitlab`)
- `AI_PROVIDER` - `ollama`, `openai`, `anthropic`, `minimax`, or `copilot` (default: `ollama`)

**GitHub-specific required:**
- `GITHUB_TOKEN` - GitHub Personal Access Token
- `GITHUB_OWNER` - Repository owner
- `GITHUB_REPO` - Repository name
- `GITHUB_PR_NUMBER` - Pull Request number

**GitLab-specific required:**
- `GITLAB_TOKEN` - GitLab API token
- `CI_PROJECT_ID` - GitLab project ID
- `CI_MERGE_REQUEST_IID` - Merge Request IID
- `CI_MERGE_REQUEST_DIFF_BASE_SHA` - Base diff SHA

**Optional behavior vars:**
- `DELETE_BOT_COMMENTS` - Delete previous bot comments (default: `true`)
- `RUN_TIMEOUT` - Max review duration (default: `10m`)
- `CI_SERVER_URL` - GitLab server URL (default: `https://gitlab.com`)

**Secrets location:**
- GitHub Actions: Repository secrets (e.g., `GITHUB_TOKEN`, `OPENAI_API_KEY`)
- GitLab CI: Project CI/CD variables (masked as sensitive)
- Local development: `.env` file or shell export (NOT committed)

## Webhooks & Callbacks

**Incoming:**
- Not applicable - Application is invoked by CI/CD pipelines, not triggered by webhooks

**Outgoing:**
- GitHub: Creates PR review comments via GitHub API (`POST /repos/{owner}/{repo}/pulls/{pull_number}/comments`)
- GitLab: Creates MR discussion notes via GitLab API (`POST /projects/{id}/merge_requests/{merge_request_iid}/notes`)
- Both comment types include inline position information (file path and line number)

---

*Integration audit: 2026-03-04*
