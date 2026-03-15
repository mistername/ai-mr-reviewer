## 1. Remove MiniMax runtime support

- [x] 1.1 Remove the MiniMax provider branch from AI provider selection and update unsupported-provider error messages to list only supported providers.
- [x] 1.2 Remove MiniMax-specific adapter code and any imports or wiring that exist only for that provider.
- [x] 1.3 Remove MiniMax configuration fields and getters that are no longer needed after provider removal.

## 2. Align tests and documentation

- [x] 2.1 Update or remove tests affected by the deleted MiniMax config and provider paths.
- [x] 2.2 Remove MiniMax references from `README.md` and any setup or environment examples.
- [x] 2.3 Run targeted validation for touched packages, then run the full lint and test suite to confirm the remaining providers still work.
