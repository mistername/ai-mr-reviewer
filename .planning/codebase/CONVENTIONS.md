# Coding Conventions

**Analysis Date:** 2026-03-04

## Naming Patterns

**Files:**
- Lowercase with underscores for multi-word files: `config.go`, `factory.go`, `reviewer.go`
- Test files suffixed with `_test.go`: `config_test.go`, `reviewer_test.go`, `entity_test.go`
- No prefixes or special naming conventions for different file types

**Packages:**
- Lowercase, single-word package names: `domain`, `application`, `config`, `ai`, `github`, `gitlab`
- Adapter packages grouped by integration: `adapters/ai/`, `adapters/github/`, `adapters/gitlab/`, `adapters/config/`
- Domain logic isolated in `internal/domain/`
- Application logic in `internal/application/`

**Functions:**
- PascalCase for exported functions: `New()`, `NewReviewer()`, `NewClient()`, `Run()`, `ReviewCode()`
- camelCase for unexported functions: `newConfig()`, `newVCSClient()`, `parseReviewResponse()`, `buildCombinedDiff()`
- Descriptive names indicating action: `GetMergeRequestChanges()`, `AddMergeRequestDiscussion()`, `DeleteBotCommentsExceptResolved()`
- Constructor functions named `New` or `NewTypeName`: `NewReviewer()`, `NewClient()`, `New()`

**Variables:**
- camelCase for local variables: `diffs`, `existing`, `reviewText`, `issues`, `filePath`
- Short names for loop counters: `i`, `d` (when type is clear from context)
- Descriptive names for struct fields: `ProjectID`, `MergeRequestIID`, `CommitSHA`
- Receiver names typically single letter: `r` for Reviewer, `c` for Config, `m` for mock types

**Types:**
- PascalCase for struct names: `Config`, `Reviewer`, `Diff`, `ReviewIssue`, `MergeRequest`, `Comment`
- PascalCase for interface names: `ConfigPort`, `MRProviderPort`, `AIProviderPort`
- Suffix interfaces with "Port" to indicate adapter/port pattern: `ConfigPort`, `MRProviderPort`, `AIProviderPort`
- Unexported structs in test files: `mrProviderMock`, `configMock`, `ollamaMock`

**Constants:**
- UPPERCASE with underscores for module-level constants: `promptFormat` (actually lowercase, not typical Go convention but used here)
- Map constants use lowercase: `languageMap`

## Code Style

**Formatting:**
- Standard Go formatting via `gofmt` with `goimports` for import organization
- Configured in `.golangci.yml` formatters section
- Line length follows Go standard (no explicit limit, but reasonable lines)
- Indentation: tabs (Go standard)

**Linting:**
- golangci-lint enabled with strict rules via `.golangci.yml`
- Active linters: asciicheck, bodyclose, copyloopvar, cyclop, dupl, dupword, errcheck, errorlint, forcetypeassert, funlen, gocritic, gosec, gosmopolitan, govet, ineffassign, loggercheck, makezero, mirror, noctx, prealloc, predeclared, staticcheck, thelper, unused, wastedassign, wrapcheck, wsl_v5
- Linter exclusions for generated code and third_party directories
- All issues reported: `max-issues-per-linter: 0`

**Imports:**
- Organized by goimports in three groups:
  1. Standard library: `context`, `encoding/json`, `fmt`, `time`, etc.
  2. Third-party: `github.com/...`, `go.uber.org/...`, `gitlab.com/...`
  3. Local project: `github.com/adlandh/ai-mr-reviewer/...`
- See `main.go` lines 3-17 for example organization

**Wrapping and Error Handling:**
- Consistent error wrapping with `fmt.Errorf()` using `%w` verb: `fmt.Errorf("create config: %w", err)`
- All errors wrapped to preserve context and stack trace
- Error messages lowercase, starting with action taken: "create config", "get MR changes", "parse review response"
- Error handling at call sites, not deferred to callers

## Struct Tags

**Environment Variables (Config):**
- Tag format: `` `env:"ENV_VAR_NAME" envDefault:"default_value"` ``
- Required fields marked with `env:"...,notEmpty"`
- Optional fields with just `env:"..."`
- See `internal/adapters/config/config.go` lines 11-41 for complete pattern

**JSON Serialization:**
- Tag format: `` `json:"field_name"` ``
- Fields in review response structs: `json:"issues"`, `json:"file"`, `json:"line"`, `json:"severity"`, `json:"message"`
- Lowercase JSON field names for external API compatibility

## Function Design

**Size:**
- Functions kept to single responsibility principle
- Typical functions 15-50 lines
- Complex operations like `parseReviewResponse()` (30 lines) and `extractJSON()` (30 lines) are acceptable when clearly focused
- Large functions flagged by `funlen` linter (enabled in `.golangci.yml`)

**Parameters:**
- Interface types preferred over concrete types for testability: `domain.ConfigPort`, `domain.MRProviderPort`
- Maximum 4-5 parameters typical
- Constructor functions pass required dependencies: `NewReviewer(config, mrProvider, aiProvider, logger)`
- Receiver method pattern common: `(r *Reviewer) Run() error`

**Return Values:**
- Single return value for simple getters: `GetVCSProvider() string`
- Error as last return value: `ReviewCode(diff string) (string, error)`
- Multiple values when needed: `([]Diff, error)`, `(map[string][]string, error)`
- Error never ignored: explicit `_ = logger.Sync()` or error handling

## Interface Design

**Port Interfaces (Adapter Pattern):**
- Interfaces defined in domain package: `internal/domain/use_cases.go`
- Named with "Port" suffix to indicate adapter boundary: `ConfigPort`, `MRProviderPort`, `AIProviderPort`
- Minimal method sets: ConfigPort has 26 getter methods, AIProviderPort has 1 method, MRProviderPort has 4 methods
- Implementation verified with explicit interface check: `var _ domain.ConfigPort = (*Config)(nil)`

## Error Handling

**Patterns:**
- Errors wrapped at source with context: `fmt.Errorf("create config: %w", err)` in `newConfig()`
- Non-critical errors logged as warnings and execution continues: `r.logger.Warn("cannot delete bot comments", zap.Error(err))`
- Critical errors propagated: `return fmt.Errorf("review failed: %w", err)`
- Early returns for error conditions: see `newGitHubClient()` lines 82-88
- Validation errors clear and specific: "GITHUB_TOKEN is required for github provider"

**Logging on Error:**
```go
if err := r.mrProvider.DeleteBotCommentsExceptResolved(); err != nil {
	r.logger.Warn("cannot delete bot comments", zap.Error(err))
}
```

## Logging

**Framework:** uber/zap (structured logging)
- Imported as: `"go.uber.org/zap"`
- Logger provided via dependency injection in main.go line 30

**Patterns:**
- Info logs for milestones: `logger.Info("starting MR Reviewer")`
- Warn logs for non-critical failures: `logger.Warn("cannot delete bot comments", zap.Error(err))`
- Error logs for critical failures: `logger.Error("review failed", zap.Error(err))`
- Contextual fields with zap helpers: `zap.Error(err)`, `zap.String("file", path)`, `zap.Int("line", line)`
- See `internal/application/reviewer.go` lines 71-73 and 156-164 for examples

**Structured Logging:**
- No string concatenation for log messages; use structured fields
- Proper typing with zap helpers: `zap.String()`, `zap.Int()`, `zap.Error()`

## Comments

**When to Comment:**
- Interface contracts documented: `//go:generate go tool mockery` before interfaces
- Non-obvious logic explained: `// Focus on correctness, security, performance, maintainability...`
- Exported functions documented: Not strictly followed; minimal comments
- Complex algorithm steps commented: See `extractJSON()` function

**Go Generate Directives:**
- Used for mock generation: `//go:generate go tool mockery` in `internal/domain/use_cases.go` line 5
- Placed immediately before interface definition

**Explicit Interface Satisfaction:**
- Used at end of files to verify interface implementation: `var _ domain.ConfigPort = (*Config)(nil)`
- Compile-time check that type satisfies interface
- Example: `main_test.go` lines 84, 92, 126

## Module Design

**Exports:**
- Constructor functions always exported: `NewReviewer()`, `NewClient()`, `New()`
- Struct types exported when used externally: `Config`, `Reviewer`
- Internal helpers unexported: `parseReviewResponse()`, `buildCombinedDiff()`, `detectLanguage()`

**Barrel/Convenience Exports:**
- Not heavily used in this codebase
- `domain/use_cases.go` contains interface definitions (conceptual barrel)

## Dependency Injection

**Pattern:** Uber FX (dependency injection framework)
- Framework: `go.uber.org/fx` imported in main.go
- Provides: `fx.New()`, `fx.Supply()`, `fx.Provide()`, `fx.Invoke()`
- Constructor functions registered with fx: `newConfig`, `newVCSClient`, `newAIProvider`, `newReviewer`
- See `main.go` lines 29-41 for fx app setup

**Factory Functions:**
- Constructor functions act as providers: `newConfig()`, `newVCSClient()`, `newGitHubClient()`, `newAIProvider()`
- Each handles validation and returns domain Port interfaces or errors
- Factories in main.go for app-level composition; package-level in adapters

## Code Patterns

**Strategy Pattern (AI Provider Selection):**
- Factory selects implementation via interface: `internal/adapters/ai/factory.go` lines 9-26
- Implementations: OllamaClient, OpenAIClient, AnthropicClient, MiniMaxClient, CopilotClient
- All satisfy `domain.AIProviderPort` interface

**Adapter Pattern (VCS Provider):**
- Main.go selects VCS client: GitHubClient or GitLabClient via factory pattern
- Both satisfy `domain.MRProviderPort` interface
- Adapter layer isolates external API details: `internal/adapters/github/`, `internal/adapters/gitlab/`

**Repository Pattern (Conceptual):**
- `MRProviderPort` acts as repository for merge request data and operations
- Abstracts GitLab and GitHub API differences

---

*Convention analysis: 2026-03-04*
