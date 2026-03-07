# TASK.md

Feature requirements and specifications for Ralph. See also `docs/intent/`, `building-intent.md`, and `AGENTS.md`.

---

## Output and verbosity: desired state

Desired behavior for **log output** (Ralph’s own messages) and **AI command output** (the child process stdout/stderr), and the flags that control them.

### Output destination (stdout vs stderr)

**Both Ralph’s operational messages and the AI command stream are the run’s log.** They both go to **stdout**. One stream, one place: `ralph run build > run.log` captures the full log (progress, stats, and AI output in order). **stderr** is reserved for output that is not part of the normal log — e.g. fatal startup errors (bad config, missing binary, usage errors). In the happy path, stderr may be unused.

### Two output elements (both on stdout)

1. **Log output** — Ralph’s operational messages: iteration progress, statistics, warnings, errors, debug lines. Destination: **stdout**. Controlled by **log level** (how much we show: error-only, warn, info, or debug).

2. **AI command output** — The stream from the AI CLI process (stdout/stderr). When “show” is on, we mirror it to the terminal in real time (in addition to capturing it for signal scanning). Destination: **stdout**. Controlled by **show AI command output** (on or off).

The two axes can be combined: e.g. verbose logs with no AI output (`-v --no-ai-cmd-output`), or quiet with no AI output (`-q`).

### Defaults

- **Log level:** info (show info, warn, error; hide debug).
- **Show AI command output:** **true** — by default we stream the AI command output to the terminal. So we only need a way to **turn it off**.

### Flags

| Flag | Effect on log level | Effect on show AI command output |
|------|---------------------|----------------------------------|
| **`-q` / `--quiet`** | Quiet: only error (or warn+error). Suppresses info/debug. | Set to **false** — do not stream AI command output. |
| **`-v` / `--verbose`** | Verbose: e.g. debug. Show all log levels. | Set to **true** (show), unless overridden by `--no-ai-cmd-output`. |
| **`--log-level <level>`** | Set log level explicitly (e.g. `debug`, `info`, `warn`, `error`). | No change. |
| **`--no-ai-cmd-output`** | No change. | Set to **false** — do not stream AI command output. |

So:

- **Quiet** sets log level to quiet and **also** suppresses AI command output.
- **Verbose** sets log level to verbose and turns on AI output unless `--no-ai-cmd-output` is set.
- **Log level** only sets the log level.
- **Show AI command output** is on by default; we only need one flag to **disable** it: `--no-ai-cmd-output`.

### Precedence and interaction

- **Log level** is determined by: `--log-level` (if present), else `-q` (quiet), else `-v` (verbose), else config/env, else default (info). (Explicit `--log-level` can override `-q`/`-v` for the log-level axis if we want; otherwise `-q`/`-v` override config.)

- **Show AI command output** is determined by:
  - If `--no-ai-cmd-output` is set → **false** (user explicitly hid it).
  - Else if `-q` is set → **false** (quiet suppresses AI output too).
  - Else if `-v` is set → **true** (verbose implies show).
  - Else config/env, else **default true**.

So: **If verbose is specified and [show AI output] is set to false, then the logging will be verbose, but the AI command output will not be included.** — The `--no-ai-cmd-output` flag wins; quiet also forces hide.

### Summary

- **Output destination:** Both Ralph operational messages and AI command stream go to **stdout** (the run’s log). stderr is for fatal/startup errors only.
- **Two axes:** log level (amount of Ralph log output) and show AI command output (whether we stream the AI CLI output).
- **Defaults:** log level = info, show AI command output = **true**.
- **Flags:** `-q` (quiet log level + suppress AI output), `-v` (verbose log level + show AI output), `--log-level <level>`, and `--no-ai-cmd-output` to disable AI output. Quiet suppresses both log verbosity and AI output; verbose increases log verbosity and turns on AI output unless `--no-ai-cmd-output` is set.
