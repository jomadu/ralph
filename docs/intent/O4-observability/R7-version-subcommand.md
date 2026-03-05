# R7: Version Subcommand

**Outcome:** O4 — Observability

## Requirement

The system provides a subcommand that reports the running Ralph version so the user and scripts can identify which binary is in use.

## Specification

Ralph provides a top-level subcommand `ralph version` (or equivalent; the exact command name is implementation-defined as long as it is discoverable via `ralph --help`).

**Behavior:**

- Invocation: `ralph version` (no arguments, no prompt alias or run context required).
- Output: A version string is written to stdout. The format is implementation-defined (e.g. `ralph 0.1.0` or `ralph version 0.1.0`). The string must allow a human or script to identify the Ralph release.
- Exit code: 0 on success. No config files are loaded; no loop is executed.
- Optional: Build or runtime metadata (e.g. Go version, build date) may be included at the implementer's discretion but is not required by this requirement.

**Interaction with other commands:**

- `ralph version` does not require a config file, prompt alias, or any other setup. It is safe to run in any directory and in scripts (e.g. for version checks before invoking `ralph run`).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| `ralph version` with no config present | Prints version to stdout, exits 0 |
| `ralph version` in a directory with ralph-config.yml | Prints version; config is not loaded |
| `ralph --help` or `ralph -h` | Help output includes the version subcommand so users can discover it |

### Examples

#### Basic invocation

**Input:**
`ralph version`

**Expected output:**
A line on stdout containing a version identifier (e.g. `ralph 0.1.0`).

**Verification:**
- Exit code 0
- stdout contains a recognizable version string

## Acceptance criteria

- [ ] Invoking the version subcommand prints a version string to stdout and exits 0
- [ ] The version subcommand does not require config, a prompt alias, or a run context
- [ ] The subcommand is listed in the main help output so users can discover it

## Dependencies

_None identified._
