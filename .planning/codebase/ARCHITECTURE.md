# Architecture

**Analysis Date:** 2026-03-04

## Pattern Overview

**Overall:** Hexagonal Architecture (Ports & Adapters)

**Key Characteristics:**
- Clear separation between domain logic and external integrations
- Dependency injection via `go.uber.org/fx` for container management
- Interface-based design with ports defining contracts
- Concrete adapters implementing multiple provider integrations
- Single responsibility with focused layers

## Layers

**Domain (Core Business Logic):**
- Purpose: Contains pure business rules, interfaces (ports), and entities independent of frameworks
- Location: `internal/domain/`
- Contains: Port interfaces (`ConfigPort`, `MRProviderPort`, `AIProviderPort`), entity definitions (`MergeRequest`, `Diff`, `Review`, `ReviewIssue`), use case definitions, prompt templates
- Depends on: Nothing (framework-independent)
- Used by: Application layer, adapters for interface compliance verification

**Application (Use Cases & Orchestration):**
- Purpose: Implements business workflows, orchestrates ports, handles cross-cutting concerns
- Location: `internal/application/`
- Contains: `Reviewer` struct that coordinates the review process - fetches MR changes, filters diffs, calls AI provider, parses responses, posts comments
- Depends on: Domain interfaces (`ConfigPort`, `MRProviderPort`, `AIProviderPort`), Zap logger for structured logging
- Used by: Main entry point

**Adapters (External Integrations):**
- Purpose: Implement domain ports, handle external service communication
- Location: `internal/adapters/`
- Contains:
  - **Config adapter** (`adapters/config/`) - Parses environment variables via `caarlos0/env`, implements `ConfigPort`
  - **VCS adapters** (`adapters/github/`, `adapters/gitlab/`) - Implement `MRProviderPort` for different Git platforms
  - **AI adapters** (`adapters/ai/`) - Implement `AIProviderPort` for multiple LLM providers (Ollama, OpenAI, Anthropic, MiniMax, Copilot)
- Depends on: Domain interfaces, external SDKs (GitHub API, GitLab API, HTTP clients)
- Used by: Main factory functions for dependency injection

**Entry Point (Bootstrap):**
- Purpose: Initializes DI container, coordinates component construction, runs application
- Location: `main.go`
- Contains: `main()` entry point, factory functions (`newConfig`, `newVCSClient`, `newAIProvider`, `newReviewer`), `runReview()` orchestrator
- Depends on: All layers via fx.Provide declarations
- Uses: Zap for logging, fx for dependency injection and lifecycle management

## Data Flow

**Review Processing Pipeline:**

1. **Initialization** (`main.go`)
   - Create logger (Zap)
   - Build fx container with all dependencies
   - Load config from environment via `ConfigAdapter`
   - Select and instantiate VCS provider (GitHub or GitLab) based on config
   - Select and instantiate AI provider (Ollama, OpenAI, Anthropic, MiniMax, or Copilot) based on config

2. **Review Execution** (`Reviewer.Run()`)
   - Delete existing bot comments (optional, controlled by `DELETE_BOT_COMMENTS`)
   - Fetch existing comments from MR via `MRProviderPort.GetExistingComments()`
   - Fetch all diffs from MR via `MRProviderPort.GetMergeRequestChanges()`
   - Filter to only new diffs (files without existing bot comments)

3. **Diff Analysis** (`Reviewer.reviewDiffs()`)
   - Combine filtered diffs into single prompt with file paths, languages, and content
   - Send combined diff to AI provider via `AIProviderPort.ReviewCode()`
   - Parse JSON response from AI provider - extracts file, line, severity, message
   - Validate issues against known files and skip unknown file references

4. **Comment Posting** (per issue from review)
   - Format comment body with severity prefix
   - Post comment at specific file/line via `MRProviderPort.AddMergeRequestDiscussion()`
   - On GitHub: attempt review comment first, fallback to issue comment if unable to comment on line

**State Management:**
- **Config State:** Immutable, loaded once at startup via environment variables
- **MR State:** Fetched fresh during each run (not cached)
- **Comment State:** Fetched to avoid duplicate reviews, deleted at start if configured
- **Review State:** Parsed from AI response, validated, then posted immediately (no intermediate persistence)

## Key Abstractions

**ConfigPort:** `internal/domain/use_cases.go`
- Purpose: Abstract configuration access from external sources
- Implementations: `internal/adapters/config/config.go` (Environment variables via caarlos0/env)
- Pattern: Interface with getters for each configuration parameter

**MRProviderPort:** `internal/domain/use_cases.go`
- Purpose: Abstract VCS operations (Merge Request/Pull Request handling)
- Implementations: 
  - `internal/adapters/github/client.go` (GitHub via google/go-github)
  - `internal/adapters/gitlab/client.go` (GitLab via gitlab.com/gitlab-org/api/client-go)
- Pattern: Interface defining 4 operations: get changes, get existing comments, add discussion, delete bot comments
- Key detail: GitHub adapter has fallback logic - if review comment on specific line fails, posts as general PR comment

**AIProviderPort:** `internal/domain/use_cases.go`
- Purpose: Abstract code review AI requests
- Implementations:
  - `internal/adapters/ai/ollama.go` (Local/remote Ollama instances)
  - `internal/adapters/ai/openai.go` (OpenAI API via chat completions)
  - `internal/adapters/ai/anthropic.go` (Anthropic Claude API)
  - `internal/adapters/ai/minimax.go` (MiniMax API)
  - `internal/adapters/ai/copilot.go` (GitHub Copilot Models API)
- Pattern: Factory pattern in `adapters/ai/factory.go` selects implementation, shared HTTP communication via `sendChatCompletionRequest()`

**Domain Entities:**
- `Diff` - Represents a single file's changes (OldPath, NewPath, Content)
- `ReviewIssue` - Represents identified problem (FilePath, Line, Severity, Message)
- `Review` - Groups issues by file with detected language
- `Comment` - Structure for posting (FilePath, Body, Line)

## Entry Points

**main():**
- Location: `main.go:19-57`
- Triggers: Program execution
- Responsibilities: 
  - Initialize Zap logger for structured logging
  - Create fx container with all providers and invokers
  - Start container (DI setup phase)
  - Stop container gracefully (cleanup phase)
  - Handle startup/shutdown timeouts (30 seconds each)

**runReview():**
- Location: `main.go:155-177`
- Triggers: Called by fx during container start
- Responsibilities:
  - Execute `Reviewer.Run()` in goroutine with timeout
  - Handle review errors and timeout scenarios
  - Return error to fx for container stop
  - Logs success/failure

**Reviewer.Run():**
- Location: `internal/application/reviewer.go:69-100`
- Triggers: Called from `runReview()`
- Responsibilities:
  - Orchestrate entire review workflow
  - Handle optional bot comment cleanup
  - Manage early exit if no new diffs found
  - Coordinate calls to MR provider and AI provider

## Error Handling

**Strategy:** Wrapping errors with context via `fmt.Errorf(...%w...)`

**Patterns:**
- All adapter constructors return `(*Type, error)` - fail-fast on initialization
- All main operations return `error` - propagate with context wrapping
- Application layer catches warnings (existing comments, bot comment deletion) but continues
- Validation errors in main.go factory functions during DI setup (fail before starting review)
- JSON parsing from AI response treats empty/invalid JSON as no issues (graceful degradation)

**Error Categories:**
1. Configuration errors - checked at startup in `newConfig()`, `newVCSClient()`, `newAIProvider()`
2. API errors - wrapped with context (e.g., "failed to get PR files: %w")
3. Parsing errors - JSON unmarshaling and line number parsing with context
4. Comments failures - logged as warnings, don't stop other issues from posting

## Cross-Cutting Concerns

**Logging:** 
- Framework: `go.uber.org/zap`
- Integration: Passed to `Reviewer`, used in `Reviewer.Run()` for warnings and errors
- Also passed to `runReview()` for info/error messages
- Key messages: "starting MR Reviewer", "review completed", warnings for deleted comments and failed posts

**Validation:**
- Config validation: In factory functions - ensures required fields present before creating clients
- Diff validation: In `Reviewer.reviewDiffs()` - validates issue file paths against known files
- JSON validation: `isJSON()` function checks structure before unmarshaling

**Authentication:**
- Pattern: Bearer token authorization in HTTP headers (OpenAI, Anthropic, MiniMax, Copilot)
- GitHub Copilot: Uses GitHub token via `client.WithAuthToken()`
- Stored in environment variables, passed via config port

---

*Architecture analysis: 2026-03-04*
