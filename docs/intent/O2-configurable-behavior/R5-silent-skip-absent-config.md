# R5: Silent Skip for Absent Config Files

**Outcome:** O2 — Configurable Behavior

## Requirement

The system silently skips config files that do not exist at their default locations, applying only the layers that are present. This allows Ralph to run with no config files at all (using built-in defaults) and avoids requiring users to create config files before first use.

## Specification

Ralph loads config from two default file locations and one optional explicit path. The behavior differs between default locations (silent skip when absent) and the explicit path (error when absent).

**Default config file locations:**

| Layer | Path resolution | Fallback chain |
|-------|----------------|----------------|
| Global config | `$RALPH_CONFIG_HOME/ralph-config.yml` → `$XDG_CONFIG_HOME/ralph/ralph-config.yml` → `~/.config/ralph/ralph-config.yml` | First existing directory in the chain is used; if none exist, global config is skipped |
| Workspace config | `./ralph-config.yml` (relative to current working directory) | No fallback |

**Loading behavior for default locations:**

1. Resolve the global config path using the fallback chain above.
2. Attempt to open the file. If the file does not exist, skip it silently — no warning, no error, no log message at info level. A debug-level log is acceptable: `config: global config not found at <path>, skipping`.
3. If the file exists but is not readable (permission denied), emit a warning and skip: `warning: cannot read global config <path>: permission denied`.
4. If the file exists and is readable, parse it. YAML parse errors are fatal — exit with code 1 and the parse error message.
5. Repeat steps 2–4 for the workspace config path.

**Explicit config path (`--config <path>`):**

When the user provides `--config <path>`:

1. The explicit file replaces **both** file-based layers (global and workspace). Neither `~/.config/ralph/ralph-config.yml` nor `./ralph-config.yml` is loaded.
2. If the file at `<path>` does not exist, Ralph exits with error: `config file not found: <path>`. Exit code 1.
3. If the file exists but is not readable, Ralph exits with error: `config file not readable: <path>: permission denied`. Exit code 1.
4. If the file exists and is readable, parse it as the sole file-based config layer.

The distinction is intent: default locations are optional convenience; `--config` expresses an explicit expectation of a fully deterministic file-based config. This is critical for testing — an integration test that passes `--config ./testdata/test-config.yml` (or another path under the project's test fixture directory) is not polluted by a developer's global config in their home directory.

**Test fixture configs:** Config files used only for tests (e.g. integration or manual testing) live in the repository root directory `testdata/`. Example: `testdata/test-config.yml`. Tests that need a deterministic config source should use `--config testdata/<fixture>.yml` when run from the repository root.

**Global config directory resolution:**

The global config directory is resolved once at startup:

1. If `RALPH_CONFIG_HOME` is set and non-empty, use its value as the directory.
2. Else if `XDG_CONFIG_HOME` is set and non-empty, use `$XDG_CONFIG_HOME/ralph/` as the directory.
3. Else use `~/.config/ralph/` as the directory.

The config file name within that directory is always `ralph-config.yml`. Ralph does not create this directory or file — it only reads from it if it exists.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Neither global nor workspace config file exists | Ralph starts with built-in defaults. No warnings, no errors. |
| Global config exists but workspace config does not | Global config is loaded, workspace layer is skipped silently |
| Workspace config exists but global config does not | Workspace config is loaded, global layer is skipped silently |
| `--config ./custom.yml` and `./custom.yml` does not exist | Error: `config file not found: ./custom.yml`. Exit 1. |
| `--config ./custom.yml` and `./custom.yml` exists | File is used as the sole file-based config layer. Neither global config nor default `./ralph-config.yml` is loaded. |
| `--config ./custom.yml` and a global config exists at `~/.config/ralph/ralph-config.yml` | Global config is not loaded. Only the explicit file is used. |
| Global config file exists but contains invalid YAML | Fatal error: Ralph exits with the YAML parse error. Invalid syntax is never silently skipped. |
| Global config file is empty (0 bytes) | Valid — an empty YAML file parses as nil/empty. No values are loaded from this layer. No error. |
| `RALPH_CONFIG_HOME` points to a directory that exists but contains no `ralph-config.yml` | Global config is skipped silently — the file doesn't exist. |
| `RALPH_CONFIG_HOME` is set to an empty string | Treated as unset; falls through to `XDG_CONFIG_HOME` check |
| `~` in `RALPH_CONFIG_HOME` value | Ralph does not expand `~` in environment variables — the OS/shell must resolve it before Ralph sees the value |

### Examples

#### No config files at all

**Input:**
No `~/.config/ralph/ralph-config.yml`. No `./ralph-config.yml`. No `--config` flag. User runs `ralph run -f prompt.md`.

**Expected output:**
Ralph starts with built-in defaults. No errors or warnings related to config files.

**Verification:**
- Ralph starts successfully
- Debug log shows all values sourced from `default`
- No warning or error output about missing config files

#### Explicit config file missing

**Input:**
User runs `ralph run build --config ./my-config.yml`. File `./my-config.yml` does not exist.

**Expected output:**
```
error: config file not found: ./my-config.yml
```
Ralph exits with code 1.

**Verification:**
- Exit code is 1
- Error message names the explicit path
- Ralph does not fall back to `./ralph-config.yml`

#### Global config via XDG_CONFIG_HOME

**Input:**
`RALPH_CONFIG_HOME` is not set. `XDG_CONFIG_HOME` is set to `/home/user/.local/config`. File `/home/user/.local/config/ralph/ralph-config.yml` exists with `loop.failure_threshold: 10`.

**Expected output:**
Ralph loads the global config from `/home/user/.local/config/ralph/ralph-config.yml`. `failure_threshold` resolves to 10 with provenance `global`.

**Verification:**
- Debug log shows `loop.failure_threshold = 10 (source: global)`

## Acceptance criteria

- [ ] A missing global config file (~/.config/ralph/ralph-config.yml) does not produce an error or warning
- [ ] A missing workspace config file (./ralph-config.yml) does not produce an error or warning
- [ ] When no config files are present, built-in defaults are used for all values
- [ ] When --config specifies an explicit path that does not exist, Ralph exits with an error (the user explicitly named a file, so its absence is an error)
- [ ] When --config is provided, neither global config nor default workspace config is loaded — the explicit file is the sole file-based config source

## Dependencies

_None identified._
