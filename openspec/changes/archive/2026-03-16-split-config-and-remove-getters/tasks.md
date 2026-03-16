## 1. Introduce structured config types

- [x] 1.1 Replace the flat getter-based config model in `internal/adapters/config` with typed config groups for runtime, VCS, and AI concerns while preserving current env variable behavior.
- [x] 1.2 Update configuration parsing and root config construction so startup code can consume concrete config values directly.
- [x] 1.3 Remove the old getter methods and eliminate `domain.ConfigPort` from `internal/domain/use_cases.go` if it is no longer needed.

## 2. Migrate wiring and consumers

- [x] 2.1 Update `main.go` factory functions and adapter constructors to accept the new concrete config structures instead of a getter-based config interface.
- [x] 2.2 Update AI and VCS provider initialization logic to read only the relevant config subgroup for each concern.
- [x] 2.3 Replace config mocks in tests with concrete config fixtures or narrower test helpers based on the new config structures.

## 3. Validate the refactor

- [x] 3.1 Regenerate mocks only if remaining domain interfaces changed during the refactor.
- [x] 3.2 Run targeted tests for config, application, and adapter packages touched by the migration.
- [x] 3.3 Run lint and the full test suite to confirm the config refactor is behavior-preserving.
