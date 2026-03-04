# Technology Stack

**Analysis Date:** 2026-03-04

## Languages

**Primary:**
- Go 1.26 - Main application language, compiled CLI tool for code review automation

## Runtime

**Environment:**
- Go 1.26+

**Package Manager:**
- Go Modules (go.mod/go.sum)
- Lockfile: Present (go.sum)

## Frameworks

**Core:**
- Uber FX v1.24.0 - Dependency injection framework used throughout the application architecture
- Zap v1.27.1 - Structured logging framework for all application logging

**API Clients:**
- google/go-github v82.0.0 - GitHub API v3 client for Pull Request operations
- gitlab.com/gitlab-org/api/client-go v1.46.0 - GitLab API client for Merge Request operations

**Configuration:**
- caarlos0/env v11.4.0 - Environment variable parsing and type conversion

**Testing:**
- stretchr/testify v1.11.1 - Testing assertions and mocking framework
- vektra/mockery v3.6.4 - Code generation tool for mock creation (used via `go.mod` tool directive)

**Build/Dev:**
- golangci-lint - Multi-linter runner (config: `.golangci.yml`, targets Go 1.23)
- mockery v3.6.4 - Test mock code generator (generates mock implementations in `internal/domain/mocks/`)

## Key Dependencies

**Critical:**
- github.com/google/go-github/v82 - Required for GitHub PR integration with OAuth token authentication
- gitlab.com/gitlab-org/api/client-go v1.46.0 - Required for GitLab MR integration with API token authentication
- go.uber.org/fx v1.24.0 - Dependency injection orchestration for application lifecycle management
- go.uber.org/zap v1.27.1 - Structured logging for all subsystems

**Infrastructure:**
- golang.org/x/oauth2 v0.35.0 - OAuth2 support for authentication flows
- net/http - Standard library HTTP client for all external API communication (Ollama, OpenAI, Anthropic, MiniMax, Copilot)

**Supporting:**
- golang.org/x/sync v0.19.0 - Synchronization primitives
- golang.org/x/time v0.14.0 - Time utilities
- gopkg.in/yaml.v3 v3.0.1 - YAML parsing (indirect, via configuration)

## Configuration

**Environment:**
- Environment variable driven configuration via `caarlos0/env/v11`
- Required env vars: `VCS_PROVIDER`, `AI_PROVIDER`, and provider-specific credentials
- All configuration lives in `internal/adapters/config/config.go` struct with `env` tags
- Configuration accessed via `domain.ConfigPort` interface

**Build:**
- `.golangci.yml` - Linter configuration for Go 1.23+
- `.mockery.yml` - Mock code generation configuration
- `go.mod` - Module declaration with Go 1.26 requirement
- `go.sum` - Dependency lock file

## Platform Requirements

**Development:**
- Go 1.26 or higher
- golangci-lint for linting
- mockery v3.6.4 for test mock generation

**Production:**
- Single binary executable compiled from Go source
- No external runtime dependencies beyond environment variables
- HTTP client connectivity to external services (GitHub, GitLab, AI providers)

**External Service Dependencies:**
- GitHub: github.com (HTTPS API endpoint)
- GitLab: Configurable URL (default: https://gitlab.com)
- Ollama: Configurable URL (default: http://localhost:11434)
- OpenAI: https://api.openai.com/v1
- Anthropic: https://api.anthropic.com/v1/
- MiniMax: https://api.minimax.io/v1
- GitHub Copilot: https://models.github.ai/inference

---

*Stack analysis: 2026-03-04*
