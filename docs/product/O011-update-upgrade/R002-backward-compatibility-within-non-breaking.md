# R002: Backward compatibility within non-breaking

**Outcome:** O011 — Update Upgrade

## Requirement

Within a non-breaking version range (e.g. same major), existing config and prompts continue to work without required migration; Ralph does not automatically rewrite user config.

## Detail

When a user updates or upgrades within the compatibility contract (e.g. same major version), their existing config files and prompt files must keep working. No mandatory migration step may block normal use. Ralph reads user config but does not automatically migrate or rewrite it to a new format when it finds old format; if a format or behavior change requires migration, the user is given documentation or optional tools they run themselves. This requirement applies to the scope Ralph owns: config schema, CLI options, and documented behavior. Third-party tools (e.g. AI CLIs) are out of scope.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Config from older patch/minor (same major) | Ralph accepts it and runs; no forced migration. |
| New optional config keys in a release | Old config without those keys still works; new keys are additive. |
| Deprecated option still in config | Continues to work within the non-breaking range, or deprecation is documented with migration path (R004). |
| Ralph finds config in an older format | Reads it; does not rewrite or migrate it automatically. User may run an optional migration tool or follow docs. |
| Major version upgrade | May introduce breaking changes; migration is documented (R003, R004). This requirement is about within non-breaking range. |
| Prompt files | Same principle: prompts that worked before continue to work; Ralph does not rewrite prompt files on upgrade. |

### Examples

#### Config works after update

**Input:** User has existing config (prompts, loop config, AI command) from an older release. User performs a non-breaking upgrade using the documented process (R001).

**Expected output:** The product starts and runs using the existing config without requiring the user to migrate or edit the file. No automatic rewrite of the config file.

**Verification:** Same config file; the run command (or equivalent) succeeds with the same invocation and exit code semantics as before the update.

#### No automatic config migration

**Input:** A future release introduces a new config schema; user's config is in the old format. User upgrades to that release.

**Expected output:** The product does not silently rewrite or migrate the config file. Either the old format is still supported within the non-breaking contract, or the user is informed (e.g. via release notes) and given documentation or an optional migration path they run themselves.

**Verification:** User's config file on disk is unchanged by Ralph unless the user explicitly runs a migration tool or edits the file.

## Acceptance criteria

- [ ] Within a non-breaking version range (e.g. same major), config that worked in an older version continues to work in the newer version without required migration.
- [ ] Within that range, prompt files and documented workflows continue to work without required migration.
- [ ] Ralph does not automatically rewrite or migrate user config files when it finds old format; the user initiates any migration via documented steps or optional tools.
- [ ] When format or behavior changes require migration, the user has documentation or optional tools they run themselves (not automatic rewrite by Ralph at startup or upgrade).

## Dependencies

- R003: Defines what "non-breaking" means for CLI contract; R002 focuses on config and prompts.
- O009: Predictability — not rewriting config aligns with Ralph not changing user content without explicit request.
