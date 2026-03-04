# Codebase Structure

**Analysis Date:** 2026-03-04

## Directory Layout

```
ai-mr-reviewer/
├── internal/                    # Internal packages (non-exported)
│   ├── domain/                 # Domain layer (core business logic)
│   │   ├── entity.go           # Data structures for review entities
│   │   ├── entity_test.go      # Tests for entities
│   │   ├── use_cases.go        # Port interfaces (contracts)
│   │   ├── prompt.go           # Review prompt template
│   │   └── mocks/              # Generated mocks for testing
│   │       ├── AIProviderPort.gen.go
│   │       ├── MRProviderPort.gen.go
│   │       └── ConfigPort.gen.go
│   ├── application/             # Application layer (use cases/orchestration)
│   │   ├── reviewer.go         # Review workflow coordinator
│   │   └── reviewer_test.go    # Tests for reviewer
│   └── adapters/                # Adapter layer (external integrations)
│       ├── config/             # Configuration adapter
│       │   ├── config.go       # Config struct & ConfigPort implementation
│       │   └── config_test.go  # Config parsing tests
│       ├── github/             # GitHub VCS adapter
│       │   └── client.go       # GitHub client & MRProviderPort implementation
│       ├── gitlab/             # GitLab VCS adapter
│       │   └── client.go       # GitLab client & MRProviderPort implementation
│       └── ai/                 # AI provider adapters
│           ├── factory.go      # Factory for creating AI providers
│           ├── chat_completions.go  # Shared HTTP logic for chat APIs
│           ├── ollama.go       # Ollama client
│           ├── openai.go       # OpenAI client
│           ├── anthropic.go    # Anthropic client
│           ├── minimax.go      # MiniMax client
│           └── copilot.go      # GitHub Copilot client
├── main.go                      # Entry point & dependency injection
├── main_test.go                 # Integration tests
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── README.md                    # Project documentation
├── LICENSE                      # MIT License
├── .github/                     # GitHub metadata (workflows, etc.)
├── .golangci.yml               # Linter configuration
├── .mockery.yml                # Mock generation configuration
├── .lefthook.yml               # Git hook configuration
└── .deepsource.toml            # Code quality service config
```

## Directory Purposes

**internal/domain:**
- Purpose: Encapsulate core business logic and define contracts between layers
- Contains: Entity types, port interfaces (contracts for external dependencies), prompt templates
- Key files: `use_cases.go` defines the three critical ports (ConfigPort, MRProviderPort, AIProviderPort)
- Not exported: Accessible only within this package and subdirectories
- Mocks: Generated via `mockery` tool in `mocks/` subdirectory for testing

**internal/application:**
- Purpose: Implement business workflows and orchestrate dependencies
- Contains: Reviewer orchestrator that coordinates the review workflow
- Key files: `reviewer.go` contains `Reviewer` struct with `Run()` method - main workflow coordinator
- Dependencies: Imports from `internal/domain` for interfaces and zap for logging
- Single file: Application logic concentrated in one `reviewer.go` file with parsing/formatting helpers

**internal/adapters/config:**
- Purpose: Provide runtime configuration via environment variables
- Contains: Config struct with tags for `caarlos0/env` parser, getter methods implementing ConfigPort
- Key files: `config.go` has Config struct with 30+ environment variable mappings
- Pattern: Implements ConfigPort interface - verified with `var _ domain.ConfigPort = (*Config)(nil)`

**internal/adapters/github:**
- Purpose: Implement MRProviderPort for GitHub Pull Requests
- Contains: Client struct wrapping `github.com/google/go-github` SDK
- Key methods: GetMergeRequestChanges, GetExistingComments, AddMergeRequestDiscussion, DeleteBotCommentsExceptResolved
- Special handling: Fallback logic for posting comments - tries review comment first, then issue comment

**internal/adapters/gitlab:**
- Purpose: Implement MRProviderPort for GitLab Merge Requests
- Contains: Client struct wrapping `gitlab.com/gitlab-org/api/client-go` SDK
- Key methods: GetMergeRequestChanges, GetExistingComments, AddMergeRequestDiscussion, DeleteBotCommentsExceptResolved
- Storage: Holds config reference for accessing commit SHA and diff base SHA needed for discussion positioning

**internal/adapters/ai:**
- Purpose: Provide multiple LLM provider implementations for code review
- Contains: Factory pattern selector, shared HTTP chat completion logic, concrete provider clients
- Key files:
  - `factory.go` - Selects provider based on config, validates required credentials
  - `chat_completions.go` - Shared HTTP utility for POST requests with Bearer auth
  - `ollama.go`, `openai.go`, `anthropic.go`, `minimax.go`, `copilot.go` - Provider-specific implementations
- Pattern: All providers implement AIProviderPort with ReviewCode(diff string) method

## Key File Locations

**Entry Points:**
- `main.go` - Main function creates fx container and runs application
- `internal/application/reviewer.go:Run()` - Executes review workflow

**Configuration:**
- `internal/adapters/config/config.go` - Loads environment variables, provides getters
- `.mockery.yml` - Specifies mock generation for ports
- `go.mod` - Lists all dependencies (fx, zap, github SDK, gitlab SDK, env parser)

**Core Logic:**
- `internal/application/reviewer.go` - Review orchestration, diff filtering, issue parsing, comment formatting
- `internal/domain/use_cases.go` - Port interface definitions
- `internal/domain/entity.go` - Domain entities (Diff, ReviewIssue, etc.)
- `internal/domain/prompt.go` - System prompt for AI review instructions

**Testing:**
- `main_test.go` - Integration tests for full workflow
- `internal/application/reviewer_test.go` - Tests for reviewer orchestration
- `internal/adapters/config/config_test.go` - Tests for config parsing
- `internal/domain/entity_test.go` - Tests for domain entities
- `internal/domain/mocks/*` - Generated mocks for testing

## Naming Conventions

**Files:**
- Adapter implementations: `client.go` (GitHub, GitLab adapters)
- Specific providers: `{provider}.go` (ollama.go, openai.go, anthropic.go, minimax.go, copilot.go)
- Configuration: `config.go`
- Shared utilities: `{feature}.go` (chat_completions.go for shared HTTP logic)
- Factory: `factory.go` (ai adapter factory)
- Tests: `{file}_test.go` (same directory as implementation)

**Directories:**
- Plural domain names: `adapters/`, `mocks/`
- Specific service names: `github/`, `gitlab/`, `ai/`, `config/`
- Layer organization: `domain/`, `application/`, `adapters/`

**Packages:**
- Match directory name: `package config`, `package github`, `package gitlab`, `package ai`
- Each package organized as separate directory under `internal/`

**Types:**
- Adapter clients: `Client` (in github, gitlab, ai provider packages)
- Config: `Config` struct
- Reviewer: `Reviewer` struct in application
- Domain objects: `MergeRequest`, `Diff`, `ReviewIssue`, `Review`, `Comment`

**Functions:**
- Factory constructors: `New{Type}()` or `new{Type}()` (exported vs private)
- Interface implementations: Methods named according to interface contracts
- Helpers in reviewer: `filterNewDiffs`, `reviewDiffs`, `parseReviewResponse`, `buildCombinedDiff`
- Parser utilities: `extractJSON`, `findJSONStart`, `isJSON`, `findMatchingBracket`, `detectLanguage`

## Where to Add New Code

**New VCS Provider (e.g., Gitea, Bitbucket):**
- Primary code: Create new directory `internal/adapters/{provider}/` with `client.go`
- Implementation: Implement `MRProviderPort` interface (4 methods)
- Registration: Add case in `main.go:newVCSClient()` switch statement
- Tests: Create `client_test.go` in same directory

**New AI Provider (e.g., Gemini, Claude Direct):**
- Primary code: Add `{provider}.go` to `internal/adapters/ai/`
- Implementation: Implement `AIProviderPort` interface (ReviewCode method)
- Configuration: Add env vars to `internal/adapters/config/config.go` Config struct
- Registration: Add case in `internal/adapters/ai/factory.go` NewAIProvider() switch
- Tests: Create `{provider}_test.go` next to implementation

**New Review Workflow Feature:**
- Primary code: Add logic to `internal/application/reviewer.go`
- If new concept needed: Add entity to `internal/domain/entity.go`
- If new port needed: Add interface to `internal/domain/use_cases.go`
- Tests: Add to `internal/application/reviewer_test.go`

**New Cross-Cutting Concern (Retry Logic, Rate Limiting, etc):**
- If generic: Consider adding to `internal/adapters/ai/chat_completions.go` (shared utilities)
- If domain-specific: Add to relevant adapter
- Wrap existing calls with new logic without changing interface contracts

## Special Directories

**internal/domain/mocks/:**
- Purpose: Generated mock implementations for interfaces
- Generated: Yes - created by `mockery` tool from port definitions
- Committed: Yes - included in git for reproducible tests
- Configuration: `go.mod` includes `github.com/vektra/mockery/v3` as tool dependency

**.github/:**
- Purpose: GitHub workflows and repository metadata
- Contains: CI/CD pipeline definitions
- Committed: Yes

**Tests (main_test.go, *_test.go):**
- Purpose: Unit and integration tests
- Executed: Via `go test ./...`
- Coverage: Critical paths - config parsing, reviewer orchestration
- Pattern: Each package can have _test.go files (not in separate test directory)

---

*Structure analysis: 2026-03-04*
