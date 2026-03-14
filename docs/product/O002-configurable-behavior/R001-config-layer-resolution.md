# R001: Config Layer Resolution

**Outcome:** O002 — Configurable Behavior

## Requirement

The system resolves configuration from defined layers (defaults, global file, workspace file, explicit file, environment, prompt-level overrides, command-line options) with a defined override order. Later layers override earlier ones for the same setting.

## Detail

Configuration is merged from multiple sources so the user can keep shared defaults, override per project or per prompt, and still override once at run time without editing files. Resolution is deterministic: for each setting, the effective value is the one from the highest-priority layer that supplies that setting.

**Layers (lowest to highest priority):**

1. **Defaults** — Built-in values so the tool works out of the box. No config file required.
2. **Global config file** — User-level; stored in the user's config directory (platform-specific; may be overridden by a documented environment variable). Optional: if missing, skipped.
3. **Workspace config file** — Project-level in the current working directory. Optional; if missing, skipped. Overrides global for the same setting when both exist.
4. **Explicit config file** — When the user points the tool at a specific config file via the documented config file option, that file is the only file-based source: global and workspace are not loaded. The file must exist or the system reports an error (see R005).
5. **Environment variables** — Override file-based config. A documented environment variable can also control where the tool looks for the user's global config file.
6. **Prompt-level overrides** — In config files, each prompt can specify its own loop settings; those apply when running or listing that prompt and override root loop settings, but are still overridden by environment and command-line options for that run.
7. **Command-line options** — Override all other layers for that run.

When the user specifies an explicit config file, layers 2 and 3 are not loaded; only that file (plus environment and command-line options) applies.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No config files exist | Defaults apply; resolution succeeds. |
| Global config missing | Skipped; workspace (if present) and defaults apply. |
| Workspace config missing | Skipped; global (if present) and defaults apply. |
| Explicit config path supplied | Only that file is used; global and workspace not loaded. Explicit file missing → error (R005). |
| Same setting in global and workspace | Workspace value wins for that setting. |
| Same setting in file and environment | Environment value wins. |
| Same setting in resolved config and command-line option | Command-line option wins. |
| Prompt override for a loop setting | When running that prompt, prompt override wins over root loop section; environment and command-line options still override for that run. |

### Examples

#### Resolution with global and workspace

**Input:** Global config sets the configured iteration limit to 10; workspace config sets it to 5. No explicit config file; no environment or command-line override.

**Expected output:** Effective iteration limit is 5 (workspace overrides global).

**Verification:** Run a loop and observe it stops after 5 iterations (or inspect resolved config in documentation or tooling).

#### Command-line overrides file

**Input:** Workspace config sets the failure threshold to 3. User runs with the documented option to set failure threshold to 1 for that run.

**Expected output:** Effective failure threshold for that run is 1.

**Verification:** Run until first failure; loop exits after 1 consecutive failure.

#### Explicit file excludes global and workspace

**Input:** User runs with the documented config file option pointing to a specific file (e.g. a user-chosen path). Global and workspace config files exist.

**Expected output:** Only settings from the specified file (plus defaults for unspecified settings, environment, and command-line options) apply. Global and workspace are not read.

**Verification:** Change global or workspace config; run again with the same explicit config file; behavior unchanged by global or workspace.

## Acceptance criteria

- [ ] For each configurable setting, the system uses a defined set of layers and override order (defaults → global → workspace → explicit file → environment → prompt overrides → command-line options).
- [ ] When no explicit config file is specified, global config (if present) and workspace config (if present) are merged; workspace overrides global for the same setting.
- [ ] When an explicit config file is specified, only that file is used for file-based config; global and workspace are not loaded.
- [ ] Environment variables override file-based values for the same setting.
- [ ] Command-line options override environment and file-based values for the same setting for that run.
- [ ] Prompt-level overrides in config apply when running or listing that prompt and override root loop settings; environment and command-line options still override for that run.
- [ ] Missing global or workspace config file does not cause an error; that layer is skipped.

## Dependencies

- None (foundational for other O002 requirements that consume resolved config).
