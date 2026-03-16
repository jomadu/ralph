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

- **Stdout** — Stdout of the AI process, captured for the caller (run-loop for signal scanning, review for report generation). When `max_output_buffer` is set, only the last N bytes are retained (sliding window) so the last line is preserved within the cap.
- **Exit code** — The AI process exit code (for optional use by run-loop or review; Ralph's exit code is determined by run-loop or review logic, not necessarily the AI's exit code).
- **Error** — When the process could not be started (e.g. command not found, permission denied), or when the invocation fails in a way the caller must handle.

**Called by**

- Run-loop: once per iteration with the assembled prompt.
- Review: when the review flow invokes the AI to evaluate the prompt The review component invokes the AI with a review prompt that includes the report directory path (interpolated from run options). The AI creates the five report files in that directory and may respond with a short confirmation. The review component does not parse stdout for report content; it reads result.json and (for apply) revision.md from the report directory after the invoke. See [review](review.md) for report directory layout and file formats.

**Implementation (T2.1)** — The interface is implemented in `internal/backend`: type `Invoker` with method `Invoke(command string, promptBytes []byte, cwd string, env []string) (stdout []byte, exitCode int, err error)`. Package function `Invoke` is the exec-style implementation (no shell; stdin receives prompt then stream closed; full stdout captured). Empty or whitespace command returns `ErrEmptyCommand`.

## Implementation spec

### Built-in AI command aliases

Ralph ships with the following AI command aliases (see `internal/config/resolve.go`). Users can select one via config or CLI (e.g. `loop.ai_cmd_alias` or `--ai-cmd-alias`). User-defined aliases in config override a built-in alias with the same name.

| Alias | Command |
|-------|---------|
| `claude` | `claude -p --dangerously-skip-permissions` |
| `kiro` | `kiro-cli chat --no-interactive --trust-all-tools` |
| `copilot` | `copilot --yolo` |
| `cursor-agent` | `agent -p --force --output-format stream-json --stream-partial-output` |

For agents that emit structured or noisy output, users can use a wrapper script that sends progress to stderr and assistant text to stdout so Ralph can scan for signals; see [Agent wrapper pattern](../../agent-wrapper-pattern.md). Cursor is one example; that doc links to [Cursor’s headless docs — Real-time progress tracking](https://cursor.com/docs/cli/headless#real-time-progress-tracking).

The backend receives the **resolved** command string (after alias expansion by the config component); it does not resolve alias names itself.

### Invocation contract

- **Stdin** — The assembled prompt is written to the child process's stdin. The stream is closed after write so the AI receives EOF when it has consumed the prompt.
- **Stdout** — Captured in full and returned to the caller. Stderr may be captured, passed through to the user's terminal, or merged per configuration; exact behavior is implementation-defined and documented (e.g. for streaming mode).
- **No shell** — The AI command is invoked without a shell. Users who need shell features (pipes, redirects, expansion) must use a wrapper script or binary as the command.
- **Environment and cwd** — The caller passes `cwd` and `env` into the backend. When `cwd` is empty, the backend does not set the child's working directory, so the child inherits the parent's current working directory. When `env` is nil or empty, the backend uses the parent's environment (e.g. `os.Environ()`), so the child sees the same environment as Ralph. When the caller passes a non-empty `cwd` or non-empty `env`, that overrides inheritance; such overrides are supplied by the config layer when configured (no default overrides in the initial implementation).

### Validation

The backend may be called only after the CLI or run-loop has validated that the command is present and resolvable. If the backend is invoked with a command that cannot be executed, it returns an error; the caller (run-loop or review) is responsible for reporting a clear error and using the documented failure exit code.

### Timeout (T2.4)

**Placement:** Per-iteration timeout is implemented in the **backend**. The run-loop (or review, when it invokes the AI) passes the effective `timeout_seconds` (from config) into the backend; the backend kills the process after N seconds when N > 0.

- **Interface:** `Invoker.Invoke(..., timeoutSec int)`. When `timeoutSec` is 0, no timeout is applied. When `timeoutSec` > 0, the backend uses `exec.CommandContext` with a deadline; when the context deadline is exceeded, the process is killed and the backend returns `ErrTimeout` and exit code -1.
- **Observable behavior:** A single invocation does not run longer than the configured timeout; the run does not hang. The caller (run-loop) receives `ErrTimeout` and can treat the iteration as a failure (e.g. increment failure count, respect failure threshold).
