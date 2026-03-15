## Context

The repository already generates mocks for the key domain ports in `internal/domain/mocks` via `mockery`, but several tests still declare hand-written stub structs implementing the same interfaces. The clearest examples are in `main_test.go` and `internal/application/reviewer_test.go`, where local test types duplicate `ConfigPort`, `MRProviderPort`, and `AIProviderPort` behavior. This makes interface evolution noisier because changes must be mirrored in both generated mocks and handwritten test doubles.

## Goals / Non-Goals

**Goals:**
- Standardize tests on generated mocks for domain ports.
- Remove manual mock structs that only exist to satisfy those domain port interfaces.
- Keep the behavior and intent of existing tests intact while improving maintainability.
- Preserve simple test readability by using generated mock expectations consistently.

**Non-Goals:**
- Replace every helper struct in every test file regardless of purpose.
- Change production interfaces or `mockery` configuration as part of this change unless strictly required.
- Rewrite unrelated test logic or assertion style.

## Decisions

### Use generated mocks only for domain ports
- Decision: migrate tests that mock `domain.ConfigPort`, `domain.MRProviderPort`, and `domain.AIProviderPort` to `internal/domain/mocks`.
- Rationale: those interfaces are already declared as ports and already have generated mocks, so using a second mocking style is redundant.
- Alternative considered: keep handwritten mocks where they are “simple enough.” Rejected because the duplication cost grows each time the interface changes.

### Keep narrow non-port helpers if they are not replacing a generated mock
- Decision: allow tiny helper types only when they are not stand-ins for a generated domain port.
- Rationale: some tests may still need minimal local helpers for behavior that is more naturally expressed as a tiny concrete helper than a mock expectation chain.
- Alternative considered: ban all local test doubles. Rejected because the goal is consistency for port mocks, not dogmatic removal of all helpers.

### Migrate incrementally at the file level
- Decision: replace manual mocks in the currently affected files without broader repository-wide refactors in the same change.
- Rationale: file-scoped migration reduces risk and keeps the implementation reviewable.
- Alternative considered: migrate every test file that contains any helper in one pass. Rejected because it increases blast radius and review complexity.

## Risks / Trade-offs

- [Generated mock expectations can make some tests more verbose] -> Keep expectation setup close to each test and extract only small reusable helpers when repetition becomes distracting.
- [Some existing tests may rely on mutable stub state] -> Use `Run` callbacks or mock returns to preserve behavior rather than rewriting test intent.
- [Future contributors may reintroduce manual mocks] -> Add spec-level guidance and keep examples in touched tests aligned with generated mock usage.

## Migration Plan

1. Identify tests using local structs to implement domain ports.
2. Replace them with generated mocks from `internal/domain/mocks`.
3. Regenerate mocks only if interface changes are required during migration.
4. Run targeted tests for affected packages, then full lint/test validation.
5. If rollback is needed, restore the prior local test doubles in the touched files.

## Open Questions

- None. The current repo already has generated mocks for the relevant interfaces, so the migration path is straightforward.
