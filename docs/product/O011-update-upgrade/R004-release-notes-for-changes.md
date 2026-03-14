# R004: Release notes for changes

**Outcome:** O011 — Update Upgrade

## Requirement

Ralph provides release notes (or equivalent) for each release describing intentional behavior changes and deprecations so users can adjust config or scripts.

## Detail

When a release includes intentional behavior changes, deprecations, or changes that could affect config or scripts, users need a place to read about them. Release notes (or an equivalent channel, e.g. changelog, tagged release description) are published for each release and describe such changes in enough detail for users to understand what changed and how to adapt. This supports R003 (migration guidance for contract changes) and helps users assess impact before or after upgrading. Patch releases that only fix bugs or add optional behavior may have minimal notes; the requirement is that meaningful changes are communicated.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Release with breaking or behavior-changing change | Release notes describe the change and, where applicable, migration or workaround. |
| Release with deprecation | Deprecation is noted; timeline or replacement (e.g. "use X instead; will be removed in 2.0") is stated. |
| Patch release with only bug fixes | Release notes may be brief; no requirement to list every internal fix if behavior is unchanged. |
| Pre-release (alpha/beta) | May have release notes; format may be lighter; intent is that users can learn about changes. |
| No release notes location yet | Project establishes a consistent place (e.g. GitHub releases, CHANGELOG, docs) that users can find from the main docs or repo. |
| Config or CLI contract unchanged | Release notes may still note other changes (e.g. performance); "behavior changes" focus on user-facing and contract-relevant items. |

### Examples

#### User reads about a new requirement or change

**Input:** A release adds a new required config key or changes the meaning of an exit code. User looks for release notes for that release.

**Expected output:** Release notes (or equivalent) describe the new key or exit code change and how to update config or scripts. User can adjust before or after upgrading.

**Verification:** User finds the release notes, understands the change, and can make the necessary adjustments.

#### Deprecation documented

**Input:** Project deprecates an option in favor of a replacement. User upgrades to that release.

**Expected output:** Release notes state that the deprecated option is deprecated, point to the replacement, and indicate when removal is planned (e.g. next major per compatibility contract). Existing scripts using the deprecated option still work within the non-breaking range but user is informed.

**Verification:** User can find the deprecation and migration path in release notes (or equivalent).

## Acceptance criteria

- [ ] Each release has release notes (or equivalent) that users can find (e.g. from docs or repo).
- [ ] Intentional behavior changes that affect config, CLI, or scripts are described in the notes for the release that introduces them.
- [ ] Deprecations are documented with replacement or timeline so users can adjust.
- [ ] Users can use release notes to understand what changed and how to adapt config or scripts after an upgrade or update.
- [ ] Release notes need not exhaustively list every internal change; focus is on user-facing and contract-relevant changes.

## Dependencies

- R003: Release notes provide the migration guidance for contract changes; R004 ensures the channel exists.
- R001: Users upgrading via the documented process can be pointed to release notes for the version they are moving to.
