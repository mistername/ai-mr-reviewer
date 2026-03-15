### Requirement: Domain port tests use generated mocks by default
Tests that need doubles for repository domain ports SHALL use generated mocks from `internal/domain/mocks` as the default mechanism instead of hand-written structs that reimplement those interfaces.

#### Scenario: Test needs a ConfigPort double
- **WHEN** a test depends on `domain.ConfigPort`
- **THEN** the test SHALL use the generated `ConfigPort` mock from `internal/domain/mocks`

### Requirement: Manual reimplementations of generated domain mocks are removed from touched tests
When a test file is updated for this migration, it MUST remove local mock structs that only exist to implement a domain port already covered by generated mocks.

#### Scenario: Touched test file contains local MRProviderPort implementation
- **WHEN** a migrated test file contains a local struct used only to satisfy `domain.MRProviderPort`
- **THEN** that local mock SHALL be replaced by the generated `MRProviderPort` mock or removed

### Requirement: Migrated tests preserve existing behavior
Replacing manual mocks with generated mocks MUST preserve the observable test behavior, including configured return values, captured interactions, and failure expectations.

#### Scenario: Existing test asserts interaction behavior
- **WHEN** a manual mock is replaced with a generated mock
- **THEN** the migrated test SHALL continue asserting the same behavior and outcomes as before
