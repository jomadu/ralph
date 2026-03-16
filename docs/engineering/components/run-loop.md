# Run-loop

## Responsibility

The run-loop component executes the iteration loop for `ralph run`: it validates the AI command before starting, loads the prompt once and buffers it, then repeatedly invokes the backend with the assembled prompt, captures stdout, detects success or failure signals, and continues or exits based on iteration limit, failure threshold, or success. It produces run reports (completion message, iteration count, timing), distinct exit codes for success, failure-threshold exit, max-iterations exit, and interrupt, and supports dry-run (show assembled prompt without running the AI) and log level control.

Implements the requirements assigned to this component in the [engineering README](../README.md).

## Interfaces

**Consumes**

- Resolved config from the config component: prompt source (alias, file path, or stdin), loop settings (max iterations, failure threshold, timeout, success/failure signals, preamble, AI command/alias, streaming, log level, max output buffer). AI command/alias may come from config (root or per-prompt), environment, or CLI; the CLI passes the already-resolved value.
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
3. **Iterate** — For each iteration: invoke backend with assembled prompt (optionally with preamble) on stdin; capture stdout (capped by `max_output_buffer` when set—see config—so the last line is preserved within the cap). Scan **the last non-empty line** of stdout for configured success and failure signals. When both signals appear on that line, apply static precedence (success wins). If success: emit completion message, iteration count, timing; exit with success code. If failure: increment consecutive-failure count; if count ≥ failure threshold, report and exit with failure-threshold code. If max iterations reached: report and exit with max-iterations code. If timeout: treat as failure or defined behavior. On interrupt (e.g. SIGINT): exit with distinct interrupt code.
4. **Observability** — Emit iteration statistics (e.g. per-iteration timing) when configured; respect log level and quiet flag for what is printed.

### Exit code semantics (run command)

The run-loop is the authority for run exit codes. User and automation documentation must document the exact values. Semantics:

| Outcome | Exit code | When |
|---------|-----------|------|
| **Success** | 0 | Success signal detected; completion message, iteration count, timing. |
| **Error (pre-loop)** | 2 | Invalid or missing AI command, invalid config, or prompt source error; clear error message before loop starts (O001/R001, O004/R001). |
| **Failure threshold** | 4 | Consecutive failures reached; report and exit (O001/R005, O004/R003). |
| **Max iterations** | 3 | Iteration limit reached without success; report count and limit, then exit (O001/R007, O004/R004). |
| **Interrupt** | 130 | User interrupted (e.g. SIGINT/Ctrl+C, or SIGTERM on Unix); distinct code (O004/R005, T3.9). |

### Signal detection

- Success and failure signals are configured strings (or patterns). **Only the last non-empty line of the captured stdout is scanned** for the configured success and failure signals. No other part of stdout is used for signal detection. **Last non-empty line** is defined as: split stdout on newline (`\n`), trim each line (leading and trailing whitespace), take the last line that is non-empty after trim; if there is no non-empty line, the scanned region is treated as empty (no signal). Stdout capture may be capped by `max_output_buffer` (see config); when set, only the last N bytes are retained so the last line is preserved within that limit.
- **Static precedence (O001/R006):** When both success and failure signals appear **on that (last non-empty) line**, the iteration is classified by a defined rule so the outcome is never ambiguous. **Success is checked first** — if the success signal is present on that line, the iteration is treated as success regardless of the failure signal; only if the success signal is absent is the failure signal considered. So with static precedence, "success wins" when both are present on that line. This is the only supported behavior; it is documented for users and automation.

### Process exit without signal (O001/R009, T3.8)

When the AI process exits without emitting the configured success or failure signal (e.g. exit 0 with no signal in output, crash, kill, timeout, or invocation error), the iteration is treated as a **failure**: the run-loop increments the consecutive-failure count and continues or exits according to the same failure threshold as for failure-signal (R005). Thus no iteration is left undefined. The report when exiting due to threshold **distinguishes** "no signal" from "failure signal present": e.g. "Stopped after N consecutive iteration(s) without success or failure signal (threshold: T)" vs "Stopped after N consecutive failure(s) (threshold: T)", so the user can tell the two cases apart for debugging or tuning.

### Dry-run

When dry-run is enabled, the run-loop does not invoke the backend. It assembles the prompt (including preamble if configured) and outputs it to the user (stdout or logs per log level). Exit code and report semantics for "dry-run completed" are documented (e.g. 0 and a message that no run was performed).
