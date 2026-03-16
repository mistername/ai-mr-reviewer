## Why

The current configuration API exposes one large getter-heavy port that mixes VCS settings, AI provider settings, and generic runtime behavior in a single shape. Splitting that configuration into focused structs now will reduce coupling, simplify call sites, and make future config changes easier to reason about and test.

## What Changes

- Split the current monolithic config into smaller config groups such as AI-related settings, GitHub/GitLab-related settings, and shared runtime settings.
- Replace getter-based config access with direct field access on typed config structs passed where needed.
- Update app wiring, adapters, and tests to depend on narrower config inputs instead of one large `ConfigPort`.
- Remove the old getter-based config interface and the redundant getter methods from the config adapter.
- **BREAKING**: internal code depending on `domain.ConfigPort` will need to migrate to the new config structure.

## Capabilities

### New Capabilities
- `structured-runtime-config`: Define how the application exposes typed runtime configuration without a getter-based config port.

### Modified Capabilities

## Impact

- Affected code: `internal/adapters/config`, `internal/domain/use_cases.go`, `main.go`, AI/VCS adapter factories, application wiring, and tests using config mocks.
- Affected interfaces: the current `domain.ConfigPort` contract and generated config mocks.
- Affected maintenance: future config additions should land in focused config groups instead of expanding a single interface.
