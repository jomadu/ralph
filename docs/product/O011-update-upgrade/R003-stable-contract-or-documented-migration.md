# R003: Stable contract or documented migration

**Outcome:** O011 — Update Upgrade

## Requirement

Documented commands, options, and exit codes remain valid across non-breaking upgrades, or changes are documented in release notes with migration guidance.

## Detail

Scripts and documentation rely on the product's CLI contract: which commands exist, which options they accept, and what exit codes mean. Within a non-breaking version range (e.g. same major per compatibility contract), that contract stays stable so existing scripts and docs keep working. If the project intentionally changes the contract in a non-breaking release (e.g. deprecating an option while keeping it working), release notes (or equivalent) describe the change and migration guidance so users can adjust. Breaking contract changes belong in a major version (or as documented) with documented migration. The focus is user-facing contract stability and predictability, not internal implementation.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Non-breaking upgrade | Documented commands, options, and exit codes remain valid; same invocation and exit code semantics. |
| Intentional contract change in non-breaking release | Change and migration guidance are described in release notes (R004). |
| Deprecation (option or behavior) | Deprecated element continues to work within the non-breaking range; deprecation and replacement are documented. |
| Major version or documented breaking release | Migration path is documented; release notes describe the change and migration guidance. |
| Undocumented or internal behavior | Not part of the stable contract; only documented commands, options, and exit codes are in scope. |

### Examples

#### Scripts keep working after non-breaking upgrade

**Input:** User has scripts or docs that invoke the product using documented commands and options and rely on documented exit codes. User performs a non-breaking upgrade.

**Expected output:** The same invocations work; documented commands, options, and exit codes remain valid. No silent change to the contract.

**Verification:** Scripts and docs that depended on the documented contract continue to work without modification, or the user was given release notes (or equivalent) describing the change and migration guidance.

#### Contract change with migration guidance

**Input:** A release introduces a change to an option or exit code meaning. User looks for guidance.

**Expected output:** Release notes (or equivalent) describe the change and migration guidance so users can adjust config or scripts. Existing behavior is either preserved within the non-breaking range or the change is clearly communicated with a path to adapt.

**Verification:** User can find the change and migration guidance and update scripts or docs accordingly.

## Acceptance criteria

- [ ] Within a non-breaking version range (e.g. same major per compatibility contract), documented commands, options, and exit codes remain valid.
- [ ] When the project changes the contract in a way that could break scripts or docs, the change and migration guidance are documented (e.g. in release notes).
- [ ] Users can rely on the same invocation and exit code semantics across non-breaking upgrades, or learn about changes and how to adapt from release notes.

## Dependencies

- R004: Release notes are the vehicle for documenting changes and migration; R003 defines the stability expectation for the contract.
- O010: Automation relies on stable exit codes; R003 ensures they remain stable or are explicitly documented with migration guidance.
