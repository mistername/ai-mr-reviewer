## ADDED Requirements

### Requirement: Supported AI provider set
The system SHALL support only the AI providers explicitly implemented by the repository after this change: `ollama`, `openai`, `anthropic`, and `copilot`.

#### Scenario: Supported provider is selected
- **WHEN** configuration sets `AI_PROVIDER` to `ollama`, `openai`, `anthropic`, or `copilot`
- **THEN** the application SHALL resolve the matching provider implementation

### Requirement: Unsupported AI providers are rejected
The system MUST reject unsupported `AI_PROVIDER` values, including `minimax`, during application startup with an error that lists the currently supported providers.

#### Scenario: Removed MiniMax provider is configured
- **WHEN** configuration sets `AI_PROVIDER=minimax`
- **THEN** application startup SHALL fail
- **THEN** the error SHALL indicate that MiniMax is unsupported
- **THEN** the error SHALL list the supported provider values

### Requirement: Public documentation reflects supported providers
Repository documentation MUST list only currently supported AI providers and MUST NOT instruct users to configure MiniMax-specific environment variables or workflows.

#### Scenario: User reads setup documentation
- **WHEN** a user follows the repository README for AI provider setup
- **THEN** they SHALL see configuration and examples only for supported providers
