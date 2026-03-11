# Backend

## Responsibility

The backend component invokes the user-chosen AI CLI command with the assembled prompt supplied on stdin and captures the command's stdout. It inherits the process environment and working directory so that the AI command runs in the same context as Ralph. It does not parse or interpret the AI output; it only delivers input and returns output. The run-loop and review components use the backend to run the AI; they are responsible for signal detection and review logic.

Implements the requirements assigned to this component in the [engineering README](../README.md).

## Interfaces

**Consumes**

- **Prompt** — Raw bytes (assembled prompt) to be written to the AI process stdin.
- **Command** — Resolved command line (executable + args) to run. No shell; the command is invoked directly (exec-style or documented spawn API).
- **Environment** — Process environment (inherit from Ralph's process unless overridden by config).
- **Working directory** — Current working directory (inherit from Ralph's cwd unless overridden by config).

**Produces**

- **Stdout** — Full stdout of the AI process, captured for the caller (run-loop for signal scanning, review for report generation).
- **Exit code** — The AI process exit code (for optional use by run-loop or review; Ralph's exit code is determined by run-loop or review logic, not necessarily the AI's exit code).
- **Error** — When the process could not be started (e.g. command not found, permission denied), or when the invocation fails in a way the caller must handle.

**Called by**

- Run-loop: once per iteration with the assembled prompt.
- Review: when the review flow invokes the AI to evaluate the prompt (e.g. for report and suggested revision).

## Implementation spec

### Built-in AI command aliases

Ralph ships with the following AI command aliases (see `internal/config/aliases.go`). Users can select one via config or CLI (e.g. `loop.ai_cmd_alias` or `--ai-cmd-alias`). User-defined aliases in config override a built-in alias with the same name.

| Alias | Command |
|-------|---------|
| `claude` | `claude -p --dangerously-skip-permissions` |
| `kiro` | `kiro-cli chat --no-interactive --trust-all-tools` |
| `copilot` | `copilot --yolo` |
| `cursor-agent` | `agent -p --force --output-format stream-json --stream-partial-output` |

The backend receives the **resolved** command string (after alias expansion by the config component); it does not resolve alias names itself.

### Invocation contract

- **Stdin** — The assembled prompt is written to the child process's stdin. The stream is closed after write so the AI receives EOF when it has consumed the prompt.
- **Stdout** — Captured in full and returned to the caller. Stderr may be captured, passed through to the user's terminal, or merged per configuration; exact behavior is implementation-defined and documented (e.g. for streaming mode).
- **No shell** — The AI command is invoked without a shell. Users who need shell features (pipes, redirects, expansion) must use a wrapper script or binary as the command.
- **Environment and cwd** — By default the child process inherits the parent's environment and current working directory. Overrides (if any) are config-defined and documented.

### Validation

The backend may be called only after the CLI or run-loop has validated that the command is present and resolvable. If the backend is invoked with a command that cannot be executed, it returns an error; the caller (run-loop or review) is responsible for reporting a clear error and using the documented failure exit code.

### Timeout

Per-iteration timeout (when configured) is applied by the run-loop or the backend; the run-loop may pass a timeout to the backend so that a single invocation is bounded. Exact placement (backend vs run-loop) is implementation-defined; the observable behavior is that the run does not hang beyond the configured timeout.
