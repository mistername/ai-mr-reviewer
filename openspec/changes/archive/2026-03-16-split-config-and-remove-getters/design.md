## Context

Today the application treats configuration as one large `domain.ConfigPort` interface with dozens of getter methods. That shape is consumed across startup wiring, VCS adapter creation, AI adapter creation, and tests. The actual source config already lives in one concrete env-backed struct, so the interface mostly adds indirection while forcing tests and generated mocks to model a broad surface area they rarely need in full.

## Goals / Non-Goals

**Goals:**
- Replace the getter-heavy config port with typed config structs.
- Split config by responsibility, at minimum across VCS settings, AI settings, and shared runtime settings.
- Pass narrower config values into factories and adapters so call sites depend only on the fields they need.
- Reduce test setup complexity by eliminating the need to mock a wide config interface.

**Non-Goals:**
- Change environment variable names or user-facing configuration semantics unless required by the refactor.
- Redesign the reviewer workflow or provider behavior.
- Introduce a new external configuration framework.

## Decisions

### Use concrete typed config groups instead of a domain config port
- Decision: parse env into a root config struct composed of smaller typed sections such as `VCSConfig`, `AIConfig`, and `RuntimeConfig`, then pass those concrete values through startup wiring.
- Rationale: config is application wiring data, not a domain abstraction; concrete types are simpler and remove a large mock surface.
- Alternative considered: keep `ConfigPort` but split it into multiple smaller getter-based interfaces. Rejected because it preserves the same indirection pattern and still pushes tests toward interface mocks for plain data.

### Keep provider-specific settings grouped under the owning concern
- Decision: model AI config with provider-specific nested settings and model VCS config with platform-specific nested settings.
- Rationale: this matches how the app actually branches at startup and makes it clearer which fields matter for each provider.
- Alternative considered: keep one flat struct and just remove getters. Rejected because it reduces method count but not conceptual coupling.

### Migrate startup and tests in one coherent refactor
- Decision: update `main.go`, adapters, and tests in the same change so no compatibility layer is needed for the old config port.
- Rationale: a temporary compatibility layer would increase churn and duplicate config access paths.
- Alternative considered: add an adapter from the new config types back to `ConfigPort`. Rejected because the stated goal is to get rid of getters, not wrap them.

## Risks / Trade-offs

- [Large refactor touches many call sites] -> Keep the new config model small and move package-by-package, validating tests after each slice.
- [Tests may need substantial rewrites] -> Prefer simple concrete config fixtures once getters are removed, which should reduce long-term complexity.
- [Startup wiring may become temporarily inconsistent] -> Change constructors in a narrow order: config adapter, then VCS/AI factories, then reviewer/bootstrap tests.

## Migration Plan

1. Introduce the new concrete config types and env parsing.
2. Update startup wiring to pass concrete config groups into factory functions.
3. Remove `domain.ConfigPort` and migrate affected adapters/application code.
4. Replace generated config mocks in tests with concrete config fixtures where possible.
5. Run generation only if any remaining domain interfaces change, then run full lint/tests.

## Open Questions

- None. The desired direction is explicit: split config by concern and remove getter-based access.
