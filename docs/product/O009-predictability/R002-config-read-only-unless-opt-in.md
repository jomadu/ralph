# R002: Config read-only unless opt-in

**Outcome:** O009 — Predictability

## Requirement

Ralph reads user config but does not rewrite or migrate it unless the user invokes a documented opt-in flow that is described as modifying config.

## Detail

Ralph uses user configuration from config files (e.g. global, workspace, or explicitly specified file per O002). Reading and merging config is required for run, review, list, and other commands. Ralph must not write to those config files — no automatic normalization, no migration to a new schema, no comment stripping or reformatting — unless the user explicitly invokes a documented flow that is described as modifying config (e.g. "migrate config," "upgrade config," or "write config"). If such an opt-in flow exists, it must be clearly documented that it changes the user's config files. Normal operation (run, review, list, etc.) is read-only for config.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User runs any normal command (run, review, list) | Config files are read only; no write to global, workspace, or explicit config file. |
| Config file has non-canonical formatting or legacy keys | Ralph reads and uses the config; does not rewrite the file to normalize or remove unknown keys. |
| Config schema version could be "upgraded" | Ralph does not automatically migrate or rewrite config; user must invoke a documented opt-in migration flow if one exists. |
| User invokes a documented "migrate config" or "write config" command | Ralph may write to the config file(s) as described in the documentation for that command. |
| No opt-in config-modification flow exists | All config access is read-only. |

### Examples

#### Normal run does not modify config

**Input:** User runs the run command with workspace config (e.g. at ./.ralph/config.yml) that contains comments and non-standard key order.

**Expected output:** Ralph loads config and runs. The config file is unchanged on disk (same content, comments, order).

**Verification:** Checksum or content of config file before and after run is identical.

#### List prompts does not modify config

**Input:** User runs the list-prompts command (or equivalent). Config defines prompts and other options.

**Expected output:** Ralph reads config and lists prompts. No config file is written or modified.

**Verification:** Config file unchanged after the command.

#### Documented opt-in migration (if implemented)

**Input:** User runs a documented config-modification command (e.g. migrate config) that is described in docs as updating the config file to a new schema.

**Expected output:** Ralph may rewrite the target config file(s) as documented. The command is explicitly described as modifying config.

**Verification:** Documentation states that the command modifies config; user has opted in by invoking it.

## Acceptance criteria

- [ ] When the user runs run, review, list, or any other command that is not a documented config-modification flow, Ralph does not write to the user's config files (global, workspace, or explicitly specified).
- [ ] Ralph does not automatically normalize, migrate, or reformat config files (e.g. strip comments, reorder keys, remove unknown keys) as a side effect of reading config.
- [ ] If Ralph provides a flow that modifies config (e.g. migration, upgrade), that flow is documented as modifying config and is invoked only when the user explicitly runs it (opt-in).
- [ ] Config files are read for resolution and use; no silent writes to config for "fixes" or "upgrades" during normal operation.

## Dependencies

- O002 (Configurable Behavior) defines config sources and layer resolution; R002 constrains that Ralph does not write to those sources unless the user opts in via a documented flow.
