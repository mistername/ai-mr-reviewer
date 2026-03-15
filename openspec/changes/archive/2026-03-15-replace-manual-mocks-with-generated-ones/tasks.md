## 1. Migrate manual port mocks in existing tests

- [x] 1.1 Replace manual `domain.ConfigPort`, `domain.MRProviderPort`, and `domain.AIProviderPort` implementations in `main_test.go` with generated mocks from `internal/domain/mocks`.
- [x] 1.2 Replace manual port mock implementations in `internal/application/reviewer_test.go` with generated mocks while preserving current test behavior and assertions.
- [x] 1.3 Remove now-unused local mock structs and helper methods from the migrated test files.

## 2. Validate and stabilize the migration

- [x] 2.1 Regenerate mocks only if interface-driven changes require it, and keep generated files current.
- [x] 2.2 Run targeted tests for touched packages and fix any expectation or behavior regressions caused by the migration.
- [x] 2.3 Run the relevant lint/test suite to confirm the repository still passes with generated mocks as the default approach.
