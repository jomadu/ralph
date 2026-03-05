# R5: Log Level Control

**Outcome:** O4 — Observability

## Requirement

The system supports configurable log verbosity levels, allowing the user to control how much operational output Ralph produces. Log levels govern Ralph's own operational messages (iteration progress, provenance, warnings, errors) and do not control AI output streaming, which is managed independently by the --verbose flag (see O4/R3).

## Specification

Ralph emits operational messages (e.g., iteration progress per R6, completion statistics per R2, errors, warnings, provenance). Log level controls which of these messages are emitted. There are four levels, from most to least verbose: **debug**, **info**, **warn**, **error**. A message at a given level is emitted only if the effective log level is at or below that level (i.e., more or equal verbosity). So: at **error**, only error messages; at **warn**, warn and error; at **info**, info, warn, and error; at **debug**, all four.

**Supported levels (order of decreasing verbosity):** debug, info, warn, error.

**Default:** When no flag or config overrides, the effective log level is **info**.

**How effective log level is determined (precedence, highest first):**

1. **`--log-level <level>`** — Explicit level. If present, it sets the effective log level. Valid values: `debug`, `info`, `warn`, `error` (case-sensitive or case-insensitive as specified by implementation; recommend case-insensitive).
2. **`--quiet` or `-q`** — Sets effective log level to **error**. Suppresses all non-error operational output. If `--log-level` is also present, `--log-level` wins (explicit takes precedence over quiet).
3. **`--verbose` or `-v`** — In addition to enabling AI output streaming (R3), sets effective log level to **debug** for Ralph's messages. If `--log-level` is also present, `--log-level` wins for log level; `-v` still enables AI streaming regardless. So: `-v --log-level warn` → log level warn, AI output streamed.
4. **Config / environment** — If no CLI flag sets log level, config (e.g., `log_level` in `ralph-config.yml`) or environment (e.g., `RALPH_LOG_LEVEL`) may set it; otherwise default is **info**.

**Output destination:** All Ralph operational log output (at any level) goes to **stderr**. This keeps stdout clean for piping and matches R2, R6 (statistics and progress to stderr).

**Scope:** Log level governs only Ralph's own messages. It does **not** enable or disable AI CLI output streaming; that is controlled solely by `--verbose` (R3).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No log level flag, no config | Default info: info, warn, error messages shown; debug hidden |
| `--log-level debug` | All operational messages (debug, info, warn, error) shown |
| `--log-level info` | info, warn, error shown; debug hidden |
| `--log-level warn` | warn, error shown; info and debug hidden (e.g., R6 progress suppressed) |
| `--log-level error` | Only error messages shown |
| `--quiet` or `-q` | Effective level error; only errors shown; same as `--log-level error` unless overridden |
| `--verbose` | Effective level debug (for Ralph messages); R3: AI output also streamed |
| `--verbose --log-level warn` | Log level warn (explicit wins); AI output still streamed |
| `--quiet --log-level info` | Log level info (explicit wins); non-error messages shown |
| Invalid `--log-level foo` | Error; invalid value; do not default silently (per O2 validation spirit) |
| All log output | Sent to stderr only |
| AI output streaming | Unchanged by log level; only `-v` enables it (R3) |

### Examples

#### Default — info level

**Input:**
`ralph run build` with no log level or verbose flags.

**Expected output:**
User sees info-level messages (e.g., iteration progress "Iteration 1/10" per R6, completion statistics per R2). User does not see debug-level messages (if any). Errors and warnings appear if they occur.

**Verification:**
- Progress messages visible
- No debug-only lines (e.g., internal state) visible

#### Quiet — scripting

**Input:**
`ralph run build -q`. Run completes with success on iteration 2.

**Expected output:**
No iteration progress lines (R6 is info; quiet = error). No completion statistics on stderr (if stats are at info level) or only errors if any. Exit code 0. Script can rely on exit code only.

**Verification:**
- Exit code 0
- stderr has minimal or no Ralph output (only errors if present)

#### Verbose plus log level override

**Input:**
`ralph run build -v --log-level warn`.

**Expected output:**
AI CLI output is streamed (R3: -v enables streaming). Ralph's operational messages are at warn level only: no info (e.g., no iteration progress), no debug; only warn and error.

**Verification:**
- AI output visible in real time
- "Iteration N/M" progress lines not visible (info suppressed)

## Acceptance criteria

- [ ] Supported log levels: debug, info, warn, error (in order of decreasing verbosity)
- [ ] Default log level is info
- [ ] --quiet sets the effective log level to error, suppressing all non-error output from Ralph
- [ ] --verbose sets the effective log level to debug (in addition to enabling AI output streaming per O4/R3)
- [ ] --log-level explicitly sets the log level and takes precedence over --quiet and --verbose for log verbosity, but does not affect AI output streaming
- [ ] All log output goes to stderr

## Dependencies

- R3 (verbose output streaming) — `--verbose` enables AI streaming and also sets log level to debug; log level itself does not affect AI streaming. R5 and R3 must agree on precedence when both -v and --log-level are used.
- R6 (iteration progress) — progress messages are at info level; they are suppressed when log level is warn or error (e.g., --quiet).
- R2 (iteration statistics) — statistics are Ralph operational output; typically emitted at info level so they are suppressed when --quiet.
