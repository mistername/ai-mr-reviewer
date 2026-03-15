## Why

The codebase still carries MiniMax-specific provider wiring even though it is no longer a desired supported integration. Removing it now reduces configuration surface area, documentation noise, and maintenance cost for an AI provider that the project should no longer expose.

## What Changes

- Remove MiniMax as a supported AI provider from runtime selection and provider factory logic.
- Remove MiniMax-specific configuration fields, environment variables, and validation paths.
- Remove MiniMax references from user-facing documentation and setup examples.
- Remove MiniMax adapter code and any tests that exist solely for that provider.
- **BREAKING**: `AI_PROVIDER=minimax` and related MiniMax environment variables will no longer be supported.

## Capabilities

### New Capabilities
- `ai-provider-support`: Define the supported set of AI providers and the configuration/runtime behavior when unsupported providers are requested.

### Modified Capabilities

## Impact

- Affected code: `internal/adapters/ai`, `internal/adapters/config`, `internal/domain`, `main.go`, `README.md`, and related tests.
- Affected APIs: environment-based provider selection and configuration behavior.
- Affected systems: local runs, CI setups, and docs that currently mention MiniMax.
