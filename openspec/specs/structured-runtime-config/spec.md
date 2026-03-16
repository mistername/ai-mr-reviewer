### Requirement: Configuration is exposed as typed config structures
The system SHALL expose application configuration as typed concrete config structures grouped by concern rather than through a getter-based config interface.

#### Scenario: Startup wiring needs AI and VCS settings
- **WHEN** application startup constructs VCS and AI providers
- **THEN** it SHALL read from typed config structures grouped by concern
- **THEN** it SHALL NOT require a getter-based config port to access those values

### Requirement: Config is split by responsibility
The configuration model MUST separate shared runtime settings from provider-specific settings, including distinct groupings for AI-related settings and GitHub/GitLab-related settings.

#### Scenario: Code reads provider-specific configuration
- **WHEN** code needs GitHub, GitLab, or AI provider configuration
- **THEN** it SHALL access only the relevant config group for that concern

### Requirement: Tests can use concrete config fixtures
Tests using application configuration MUST be able to construct concrete config fixtures directly without implementing or mocking a broad getter-based config interface.

#### Scenario: Unit test config setup
- **WHEN** a test needs configuration values for a reviewer or factory
- **THEN** the test SHALL be able to provide those values through concrete config data structures
