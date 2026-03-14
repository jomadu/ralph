# R005: Explicit Config File Only

**Outcome:** O002 — Configurable Behavior

## Requirement

When the user specifies an explicit config file, the system uses only that file and reports an error if it is missing. Global and workspace config are not loaded.

## Detail

The user can point the tool at a specific config file via the documented config file option (e.g. a user-chosen path). This mode is for tests, CI, or running with an alternate config without changing the current directory or user config.

**When explicit config is specified:**

- The system loads configuration only from that file (plus defaults for unspecified settings, environment variables, and command-line options). Global config file and workspace config file are not read.
- The specified file must exist and be readable. If it is missing or unreadable, the system reports an error and does not proceed with run or list using that config.
- The path may be relative to the current working directory or absolute, as documented.

**When explicit config is not specified:**

- Normal layer resolution applies: defaults, then global file (if present), then workspace file (if present), then environment, then prompt overrides, then command-line options (R001). Missing global or workspace is not an error.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Explicit path supplied, file exists | Only that file is used; global and workspace not loaded. |
| Explicit path supplied, file missing | Error reported; no fallback to global or workspace. |
| Explicit path supplied, file unreadable (permissions) | Error reported. |
| Explicit path is a directory | Defined behavior (e.g. error or documented convention such as a config file inside it). |
| Explicit path empty or invalid | Error reported. |

### Examples

#### Explicit file used alone

**Input:** User runs the run command with the documented config file option pointing to a specific file (e.g. a path such as testdata/ci.yml). Files exist at global and workspace config locations.

**Expected output:** Only the specified file is loaded for file-based config. Global and workspace are ignored. Run proceeds with that config (plus defaults, environment, command-line options).

**Verification:** Change global or workspace config; run again with the same explicit config file; behavior unchanged.

#### Explicit file missing

**Input:** User runs the run command with the documented config file option pointing to a path that does not exist.

**Expected output:** The system reports an error that the config file is missing (or unreadable) and does not start the run. No silent fallback to global or workspace.

**Verification:** Run with a missing path; observe error and no run.

## Acceptance criteria

- [ ] When the user specifies an explicit config file path via the documented option, the system uses only that file for file-based configuration; global and workspace config files are not loaded.
- [ ] If the specified file does not exist, the system reports an error and does not proceed with the run/list using that config.
- [ ] If the specified file exists but is not readable (e.g. permissions), the system reports an error.
- [ ] Defaults, environment variables, and command-line options still apply on top of the explicit file when it is used.
- [ ] When no explicit config is specified, behavior is unchanged: global and workspace are loaded when present (R001).

## Dependencies

- R001 — Config layer resolution (explicit file is a distinct mode that bypasses global and workspace).
