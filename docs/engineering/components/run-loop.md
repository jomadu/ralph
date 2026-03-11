# Run-loop

## Responsibility

The run-loop component executes the iteration loop for `ralph run`: it validates the AI command before starting, loads the prompt once and buffers it, then repeatedly invokes the backend with the assembled prompt, captures stdout, detects success or failure signals (or applies AI-interpreted precedence when configured), and continues or exits based on iteration limit, failure threshold, or success. It produces run reports (completion message, iteration count, timing), distinct exit codes for success, failure-threshold exit, max-iterations exit, and interrupt, and supports dry-run (show assembled prompt without running the AI) and log level control.

Implements the requirements assigned to this component in the [engineering README](../README.md).

## Interfaces

**Consumes**

- Resolved config from the config component: prompt source (alias, file path, or stdin), loop settings (max iterations, failure threshold, timeout, success/failure signals, signal precedence mode, preamble, AI command/alias, streaming, log level).
- Optional CLI overrides for the run (e.g. dry-run, log level, config file already applied by config).

**Produces**

- **Iteration outcome** — Success (success signal detected), failure-threshold reached, max-iterations reached, or interrupt. Each outcome has a documented exit code.
- **Run report** — Completion message, iteration count, timing (and on failure-threshold or max-iterations: report and exit code per observability requirements). On success: completion message, count, timing, and documented success exit code.
- **Dry-run output** — When dry-run is set, the assembled prompt is shown (e.g. to stdout or logs); no backend invocation.
- **Exit code** — Documented success code on success; distinct codes for failure threshold, max iterations, interrupt, and clear error (e.g. missing AI command). Exact numeric values are documented in user docs and automation docs; this component is the producer of those codes.

**Calls**

- Config: already resolved by CLI; run-loop receives effective config.
- Backend: once per iteration with the assembled prompt on stdin; captures stdout for signal detection.

## Implementation spec

### Loop algorithm

1. **Validate** — Resolve AI command (alias or direct); if missing or invalid, report clear error and exit with documented error code. Do not start the loop.
2. **Load prompt once** — Read prompt from the resolved source (alias → file, file path, or stdin); buffer in memory. Prompt is not re-read between iterations.
3. **Iterate** — For each iteration: invoke backend with assembled prompt (optionally with preamble) on stdin; capture stdout. Scan stdout for configured success and failure signals. Apply signal precedence (static: first match wins; or AI-interpreted when configured). If success: emit completion message, iteration count, timing; exit with success code. If failure: increment consecutive-failure count; if count ≥ failure threshold, report and exit with failure-threshold code. If max iterations reached: report and exit with max-iterations code. If timeout: treat as failure or defined behavior. On interrupt (e.g. SIGINT): exit with distinct interrupt code.
4. **Observability** — Emit iteration statistics (e.g. per-iteration timing) when configured; respect log level and quiet flag for what is printed.

### Exit code semantics (run command)

The run-loop is the authority for run exit codes. User and automation documentation must document the exact values. Semantics:

| Outcome | Exit code | When |
|---------|-----------|------|
| **Success** | 0 | Success signal detected; completion message, iteration count, timing. |
| **Error (pre-loop)** | 2 | Invalid or missing AI command, invalid config, or prompt source error; clear error message before loop starts (O001/R001, O004/R001). |
| **Failure threshold** | (TBD) | Consecutive failures reached; report and exit. |
| **Max iterations** | (TBD) | Iteration limit reached without success; report and exit. |
| **Interrupt** | (TBD) | User interrupted (e.g. SIGINT); distinct code. |

### Signal detection

- Success and failure signals are configured strings (or patterns). The run-loop scans the captured stdout for their presence.
- When both success and failure signals appear in the same output, precedence is either **static** (e.g. first match wins, or a defined order) or **AI-interpreted** (when configured): the AI output may be interpreted to decide outcome; the exact mechanism is implementation-defined and documented.

### Dry-run

When dry-run is enabled, the run-loop does not invoke the backend. It assembles the prompt (including preamble if configured) and outputs it to the user (stdout or logs per log level). Exit code and report semantics for "dry-run completed" are documented (e.g. 0 and a message that no run was performed).
