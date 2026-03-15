## Context

The application currently exposes MiniMax as a first-class AI provider alongside Ollama, OpenAI, Anthropic, and Copilot. That support spans several layers: environment-backed config fields, provider selection in the AI factory, provider-specific adapter code, documentation examples, and tests. The change is small in scope but cross-cutting, so a short design helps keep code, docs, and behavior aligned.

## Goals / Non-Goals

**Goals:**
- Remove MiniMax from the set of supported AI providers.
- Ensure unsupported-provider behavior is explicit and consistent after removal.
- Remove dead code, config surface area, and documentation references tied only to MiniMax.
- Keep the remaining provider flow unchanged for Ollama, OpenAI, Anthropic, and Copilot.

**Non-Goals:**
- Redesign the AI provider abstraction.
- Change request/response behavior for the remaining providers.
- Introduce a replacement provider in the same change.

## Decisions

### Remove MiniMax end-to-end rather than leaving partial compatibility
- Decision: delete MiniMax-specific adapter wiring, config getters/fields, factory branches, and docs.
- Rationale: partial deprecation would preserve maintenance burden and ambiguous runtime behavior.
- Alternative considered: keep config parsing but return a deprecation error for `minimax`. Rejected because the project goal is to fully remove the connection, not preserve a transitional compatibility layer.

### Keep unsupported-provider validation centralized in existing factory/startup paths
- Decision: rely on the current provider selection path to reject unsupported values once `minimax` is removed from supported lists and error messages.
- Rationale: this preserves the existing architecture and avoids duplicating validation logic in config parsing.
- Alternative considered: add explicit provider allowlists in config. Rejected because current validation already happens at composition time and does not require a broader config redesign.

### Update docs and examples in the same change
- Decision: remove MiniMax from README feature lists, configuration tables, and example workflows as part of the implementation.
- Rationale: docs are part of the user-facing contract; leaving them stale would create broken setup paths.
- Alternative considered: defer docs cleanup. Rejected because it would leave the repository internally inconsistent.

## Risks / Trade-offs

- [Users still configured for MiniMax will fail at startup] -> Update unsupported-provider errors and docs so failure is immediate and understandable.
- [Tests may silently rely on MiniMax config interfaces] -> Run targeted package tests plus full suite after removing config surface.
- [Removing config methods may touch multiple packages] -> Keep changes mechanical and scoped to references that exist today.

## Migration Plan

1. Remove MiniMax references from config, provider factory, adapter code, and tests.
2. Update README and any setup guidance to list only supported providers.
3. Validate that unsupported `AI_PROVIDER` values now reject MiniMax.
4. For rollback, restore the MiniMax adapter and config wiring in a single revert if needed.

## Open Questions

- None. The desired outcome is a straightforward removal rather than a phased deprecation.
