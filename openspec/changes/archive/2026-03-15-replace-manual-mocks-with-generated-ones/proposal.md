## Why

The codebase currently mixes generated `mockery` mocks with hand-written test doubles for the same domain ports. Replacing the manual mocks now reduces duplicate maintenance, keeps tests aligned with the source interfaces, and makes future interface changes cheaper to update safely.

## What Changes

- Replace hand-written test mocks for domain ports with generated mocks from `internal/domain/mocks`.
- Update affected tests to express expectations through generated mock APIs instead of custom stub structs where practical.
- Remove redundant manual mock types and helper methods that exist only to satisfy domain interfaces in tests.
- Keep narrowly scoped non-port test helpers only where generated mocks are not the right fit.

## Capabilities

### New Capabilities
- `generated-test-mocks`: Define the expectation that tests relying on domain port doubles use generated mocks as the default mechanism.

### Modified Capabilities

## Impact

- Affected code: `main_test.go`, `internal/application/reviewer_test.go`, and any other tests with manual implementations of `domain.ConfigPort`, `domain.MRProviderPort`, or `domain.AIProviderPort`.
- Affected tooling: `go generate ./...` / `mockery` remains the source of truth for port mocks.
- Affected maintenance: interface changes should require fewer manual test updates after migration.
