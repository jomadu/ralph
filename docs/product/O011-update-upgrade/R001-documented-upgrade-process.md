# R001: Documented upgrade process

**Outcome:** O011 — Update Upgrade

## Requirement

Ralph documents how to upgrade to a chosen version and how to update within a non-breaking version (e.g. patch or minor on the same major).

## Detail

Users who already have Ralph installed need a clear, documented way to move to a specific version (upgrade) or to get the latest patch/minor within the same major line (update). The documentation covers the steps the user follows (e.g. re-run install or upgrade steps with a version identifier, or use a package manager). No automatic or unattended update is required; the user initiates the process. "Non-breaking" is defined by the project's compatibility contract (e.g. semantic versioning: same major = non-breaking).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User wants a specific version | Documented process describes how to install or upgrade to that version. |
| User wants latest within current major | Documented process describes how to update within a non-breaking range (e.g. patch/minor per compatibility contract). |
| User upgrades from source vs package manager | Documentation covers the supported methods; not every method is required. |
| First-time install | Out of scope; O006 covers install. This requirement is for upgrading or updating an existing install. |
| Uninstall then reinstall | May be one documented option; uninstall is O006. |

### Examples

#### Upgrade to a chosen version

**Input:** User has an older version installed and wants a specific newer version. User consults the documented upgrade process.

**Expected output:** Documentation describes how to upgrade to the chosen version (e.g. re-run install or upgrade steps with that version, or use a package manager). After following the steps, the new version is invocable and behaves as documented for that version.

**Verification:** User can complete the steps without guessing; the new version is invocable and behaves as documented for that version.

#### Update within non-breaking version

**Input:** User has a current install and wants the latest within the same major (a non-breaking upgrade). User follows the documented update process.

**Expected output:** Documentation explains how to get the latest patch or minor on the same major (per the compatibility contract). After update, config and workflows that worked before continue to work per R002.

**Verification:** User obtains a non-breaking upgrade; no breaking change is introduced by the process itself.

## Acceptance criteria

- [ ] Documentation describes how to upgrade to a chosen (specific) version of Ralph.
- [ ] Documentation describes how to update within a non-breaking version range (e.g. same major).
- [ ] The documented process is actionable: a user can follow it and achieve the upgrade or update.
- [ ] After completing the process, the new version is invocable and behaves as documented for that version.
- [ ] Documentation does not require Ralph to perform automatic or unattended updates; the user initiates upgrade or update.

## Dependencies

- O006 (Install/Uninstall): assumes an install path exists; upgrade/update builds on that.
