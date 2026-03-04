# Testing Patterns

**Analysis Date:** 2026-03-04

## Test Framework

**Runner:**
- Go standard testing: `testing` package (no external test runner like Jest or pytest)
- Configuration: none required; `go test` command used via build system

**Assertion Library:**
- Testify: `github.com/stretchr/testify` v1.11.1 (imported in go.mod line 8)
- Used for assertions and mocking: `testify/assert`, `testify/require`, `testify/mock`

**Run Commands:**
```bash
go test ./...              # Run all tests
go test -v ./...           # Verbose output
go test -race ./...        # Race condition detection
go test -cover ./...       # Coverage report
golangci-lint run ./...    # Linting (includes test linting)
```

**Mock Generation:**
- Tool: Mockery v3 (github.com/vektra/mockery/v3)
- Configuration: `.mockery.yml` at project root
- Command: `go generate ./...` to regenerate mocks based on //go:generate directives
- Output format: Single file per interface in `mocks/` subdirectory
- Example output files: `internal/domain/mocks/AIProviderPort.gen.go`, `internal/domain/mocks/MRProviderPort.gen.go`

## Test File Organization

**Location:**
- Co-located with source files in same package
- Test files named `*_test.go` in the same directory as code being tested
- Examples:
  - `main.go` with `main_test.go`
  - `config/config.go` with `config/config_test.go`
  - `application/reviewer.go` with `application/reviewer_test.go`
  - `domain/entity.go` with `domain/entity_test.go`

**Naming:**
- Test function names start with `Test`: `TestNewConfigWithDefaults()`, `TestParseReviewResponse()`, `TestRunReviewsOnlyNewDiffs()`
- Descriptive test names indicating what is tested: `TestProvidersReturnCorrectTypes()`, `TestEntitiesCanBeConstructed()`

**File Count:**
- 3 test files in codebase: `main_test.go`, `config_test.go`, `reviewer_test.go`, `entity_test.go`
- Test coverage focused on critical paths, not exhaustive

## Test Structure

**Suite Organization:**
```go
// Test function follows standard Go testing pattern
func TestSomething(t *testing.T) {
	// Setup
	t.Setenv("ENV_VAR", "value")
	
	// Execute
	result, err := someFunction()
	
	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != expected {
		t.Fatalf("unexpected result: %+v", result)
	}
}
```

**Patterns Observed:**

1. **Environment Variable Setup (Configuration Tests):**
   ```go
   func TestNewConfigWithDefaults(t *testing.T) {
       t.Setenv("VCS_PROVIDER", "gitlab")
       t.Setenv("GITLAB_TOKEN", "token")
       t.Setenv("CI_PROJECT_ID", "123")
       // ... more env vars
       
       cfg, err := New()
       if err != nil {
           t.Fatalf("expected config, got error: %v", err)
       }
   }
   ```
   From: `internal/adapters/config/config_test.go` lines 8-48

2. **Mock Creation and Assertion (Application Tests):**
   ```go
   func TestRunReviewsOnlyNewDiffs(t *testing.T) {
       c := &configMock{iid: "5"}
       g := &mrProviderMock{
           comments: map[string][]string{
               "already.go:1": {"ai-mr-reviewer:**WARNING**: fix it"},
           },
           diffs: []domain.Diff{
               {NewPath: "already.go", Content: "diff1"},
               {NewPath: "new.go", Content: "diff2"},
           },
       }
       o := &ollamaMock{response: `{"issues":[...]}`}
       r := NewReviewer(c, g, o, zap.NewNop())
       
       if err := r.Run(); err != nil {
           t.Fatalf("unexpected error: %v", err)
       }
       if len(g.addedDiscussions) != 1 {
           t.Fatalf("expected 1 discussion, got %d", len(g.addedDiscussions))
       }
   }
   ```
   From: `internal/application/reviewer_test.go` lines 101-124

3. **No Setup/Teardown Functions:**
   - No `SetupTest()` or `TeardownTest()` methods used
   - Setup done inline within test functions
   - Testify not used for assertion syntax; standard `if err != nil { t.Fatalf(...) }` pattern used

## Mocking

**Framework:** Manual struct implementation (not Testify mocks generated)

**Patterns:**
```go
// Manual mock struct in test file
type mrProviderMock struct {
    comments         map[string][]string
    commentsErr      error
    diffs            []domain.Diff
    diffsErr         error
    addedDiscussions []addedDiscussion
}

// Implement interface methods
func (m *mrProviderMock) GetMergeRequestChanges() ([]domain.Diff, error) {
    return m.diffs, m.diffsErr
}

func (m *mrProviderMock) GetExistingComments() (map[string][]string, error) {
    return m.comments, m.commentsErr
}

func (m *mrProviderMock) AddMergeRequestDiscussion(file string, line int, note string) error {
    m.addedDiscussions = append(m.addedDiscussions, addedDiscussion{file: file, line: line, body: note})
    return nil
}

// Verify interface implementation
var _ domain.MRProviderPort = (*mrProviderMock)(nil)
```

From: `internal/application/reviewer_test.go` lines 47-73

**What to Mock:**
- External dependencies (interfaces): `ConfigPort`, `MRProviderPort`, `AIProviderPort`
- Mocks allow injecting test data and error conditions
- Mocks capture calls for assertion: `addedDiscussions` slice captures `AddMergeRequestDiscussion()` calls

**What NOT to Mock:**
- Internal logic functions: `parseReviewResponse()`, `detectLanguage()`, `buildCombinedDiff()` tested directly
- Data structures/entities: `Config`, `Reviewer`, `Diff`, `ReviewIssue` not mocked
- Standard library types handled directly

**Generated Mocks:**
- Mockery generates mocks for domain interfaces via `//go:generate go tool mockery` directive
- Generated mocks go to `internal/domain/mocks/` directory
- Example: `internal/domain/mocks/AIProviderPort.gen.go` (generated, not manually written)
- Not actively used in provided test files; manual mocks preferred for clarity

## Fixtures and Factories

**Test Data:**
Test fixtures typically inlined in test functions rather than shared:

```go
// Inline fixture creation
func TestEntitiesCanBeConstructed(t *testing.T) {
    mr := MergeRequest{IID: 10, ProjectID: "42", Title: "MR", WebURL: "https://gitlab/mr/10"}
    diff := Diff{NewPath: "a.go", OldPath: "a.go", Content: "@@ -1 +1 @@"}
    issue := ReviewIssue{Line: 12, Severity: "warning", Message: "msg"}
    review := Review{FilePath: "a.go", Language: "Go", Issues: []ReviewIssue{issue}}
    comment := Comment{FilePath: "a.go", Line: 12, Body: "text"}
    
    if mr.IID != 10 || diff.NewPath == "" || review.Language != "Go" || comment.Line != 12 {
        t.Fatal("entity values are not set as expected")
    }
}
```

From: `internal/domain/entity_test.go` lines 7-17

**Factory Test Helpers:**
- Simple mock structs created fresh in each test
- No builder patterns or factory methods for test data
- Environment variable setup inline with `t.Setenv()`

## Test Types

**Unit Tests:**
- Scope: Individual functions and methods
- Approach: Isolated with mocked dependencies
- Examples:
  - `TestParseReviewResponse()` - tests JSON parsing logic
  - `TestDetectLanguage()` - tests language detection map lookup
  - `TestNewConfigWithDefaults()` - tests config initialization with defaults

**Integration Tests:**
- Scope: Multiple components working together
- Approach: Use mocked external services, test business logic flow
- Examples:
  - `TestRunReviewsOnlyNewDiffs()` - tests Reviewer.Run() with mocked providers
  - `TestProvidersReturnCorrectTypes()` - tests fx dependency injection setup
  - `TestFxAppValidation()` - validates fx app construction

**E2E Tests:**
- Not present in codebase
- Could be added as separate test suite in `e2e/` directory if needed

## Test-Related Configuration

**Environment Tests:**
- Use `t.Setenv()` to set environment variables in test scope
- Variables set before creating objects that read them
- Example from `config_test.go`:
  ```go
  func TestNewConfigWithDefaults(t *testing.T) {
      t.Setenv("VCS_PROVIDER", "gitlab")
      t.Setenv("GITLAB_TOKEN", "token")
      // ...
      cfg, err := New()
  }
  ```

**FX Testing:**
- `fx.ValidateApp()` used to validate dependency injection graph
- `fx.NopLogger` used in tests to suppress DI logging
- Example: `main_test.go` lines 16-40

## Common Patterns

**Async Testing:**
Not heavily used in this codebase. Review goroutines tested via channel semantics in main.go:
```go
// From main.go (not tested specifically)
result := make(chan error, 1)
go func() {
    result <- reviewer.Run()
}()

select {
case err := <-result:
    // handle result
case <-time.After(cfg.GetRunTimeout()):
    // handle timeout
}
```

**Error Testing:**
```go
func TestNewConfigMissingRequired(t *testing.T) {
    // Clear all required env vars
    t.Setenv("VCS_PROVIDER", "")
    // ... other env vars
    
    cfg, err := New()
    if err != nil {
        t.Fatalf("expected config without error, got: %v", err)
    }
    
    // Verify defaults applied
    if cfg.GetVCSProvider() != "gitlab" {
        t.Fatalf("unexpected VCS provider: %s", cfg.GetVCSProvider())
    }
}
```

From: `internal/adapters/config/config_test.go` lines 50-66

**Critical Path Testing:**
- Focus on workflows that matter most: configuration loading, provider selection, code review flow
- Example: `TestRunReviewsOnlyNewDiffs()` tests the core review filtering logic

## Coverage

**Requirements:** No explicit coverage target enforced

**View Coverage:**
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Gaps:**
- GitLab and GitHub client implementations not tested (integration-heavy)
- AI provider client implementations not tested (external API calls)
- Error cases in adapters not fully covered
- Main.go goroutine handling not directly tested

## Logging in Tests

**Framework:** zap (same as production)

**Pattern:**
- Test logger created with `zap.NewNop()` to suppress output during tests
- Example: `reviewer_test.go` line 113: `r := NewReviewer(c, g, o, zap.NewNop())`
- Allows testing code that requires logger without test output pollution

## Test Data Validation

**Approach:**
- Simple equality checks: `if got != expected { t.Fatalf(...) }`
- Custom struct validation: `if len(issues) != 1 || issues[0].Line != 3 { t.Fatalf(...) }`
- Map/slice inspection: `if len(g.addedDiscussions) != 1 { t.Fatalf(...) }`

## Generated Code Testing

**Mocks Directory:**
- `internal/domain/mocks/` contains generated mock implementations
- Files: `AIProviderPort.gen.go`, `ConfigPort.gen.go`, `MRProviderPort.gen.go`
- Regenerated via `go generate ./...`
- Testify template used for generation (configured in `.mockery.yml` line 11)

## Test Execution

**Standard Test Commands:**
- All tests in module: `go test ./...`
- Specific package: `go test ./internal/application/...`
- Single test: `go test -run TestRunReviewsOnlyNewDiffs ./internal/application/`
- With race detector: `go test -race ./...`

**CI/CD Integration:**
- Tests run via standard Go toolchain
- Linting via golangci-lint (`.golangci.yml` config)
- No custom test framework or CI-specific test commands visible

---

*Testing analysis: 2026-03-04*
