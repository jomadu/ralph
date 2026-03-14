# O004: Observability

## Who

Users who run Ralph and need to know why a command or operation stopped, how it performed, and what to do when something fails — for diagnosis, tuning, and scripting.

## Statement

The user understands why the command or operation stopped, how it performed, and what to do when something fails.

## Why it matters

A loop runner that silently exits is a black box. The user re-runs the command, checks file diffs, reads AI output, and tries to reconstruct what happened. Did it succeed? Did it hit the failure threshold? Did it time out? Did the signal never appear? Without clear termination reporting and execution statistics, diagnosing loop behavior is guesswork and tuning config is trial and error. The same applies to the review command: the user must know whether the review succeeded, whether the prompt had errors, or whether the run failed.

## Output and log controls

Two separate controls govern what the user sees:

- **Log level** — Controls how much Ralph itself logs (e.g. error, warn, info, debug): iteration progress, timing, errors, and other operational messages. The user can set log level explicitly (e.g. via config or CLI); that value overrides any shortcut.

- **Show AI command output** — Controls whether the AI process’s stdout is streamed to the terminal in real time. Ralph always captures that output for signal scanning; this setting only determines if the user sees it. Default is true in normal runs so the user can watch the AI work. When false, the user sees only Ralph’s logs and final summary.

**Quiet** is a shortcut for minimal output: it sets log level to a minimal level (e.g. error only) and show AI command output to false, so scripts and CI get only essential messages and no streamed AI output. Explicit log level or show AI command output override the shortcut where set. There is no separate “verbose” flag; the user gets more output by raising log level and/or enabling show AI command output.

## Verification

- When the AI CLI command is missing or invalid (e.g. alias not found, binary not on PATH), Ralph reports a clear error before the loop or review starts, so the user understands why the run could not start and can fix config or install the tool.
- On success signal: Ralph prints a completion message, reports iteration count and timing, exits 0.
- On failure threshold: Ralph prints the threshold value and consecutive failure count, exits with a distinct code (e.g. 1).
- On max iterations exhausted: Ralph prints the iteration count and limit, exits with a distinct code (e.g. 2).
- On interruption (e.g. SIGINT/SIGTERM): Ralph exits with a distinct code (e.g. 130).
- User sets log level and show AI command output as desired: Ralph’s logs follow the log level; when show AI command output is true, the AI’s output is streamed to the terminal while Ralph still captures it for signal scanning.
- User uses quiet: only minimal Ralph logs (e.g. errors) and no streamed AI output, unless the user explicitly overrides log level or show AI command output.
- User runs dry-run and sees the fully assembled prompt (preamble + prompt) printed without spawning an AI process.
- After a multi-iteration run, Ralph reports iteration statistics (e.g. min/max/mean duration).
- For review: exit code and report make it clear whether the review completed, the prompt had errors, or the run failed; the user can see what to do next.
- When AI-interpreted signal precedence is used and both signals appear in an iteration, the user can see (e.g. via logs or summary) that an interpretation run occurred and what outcome was used (success, failure, or fallback applied), so they can diagnose why the loop treated that iteration as it did.

## Non-outcomes

- Operational messages and the AI command stream go to stdout (the run's log). stderr is reserved for fatal or startup errors only. Ralph does not provide persistent log files.
- Ralph does not provide per-iteration diffs or file change tracking between iterations.
- Ralph does not integrate with external monitoring, alerting, or observability systems.
- Ralph does not support replaying or debugging past executions.

## Risks

| Risk | Mitigating Requirement |
|------|------------------------|
| User does not know why the command stopped | [R002 — Success report and exit zero](R002-success-report-and-exit-zero.md), [R003 — Failure threshold report and exit code](R003-failure-threshold-report-and-exit-code.md), [R004 — Max iterations report and exit code](R004-max-iterations-report-and-exit-code.md), [R005 — Distinct exit code on interrupt](R005-distinct-exit-code-on-interrupt.md) |
| User cannot see how the run performed | [R002 — Success report and exit zero](R002-success-report-and-exit-zero.md), [R008 — Iteration statistics](R008-iteration-statistics.md) |
| User cannot control verbosity or AI output visibility | [R006 — Log level and show AI output](R006-log-level-and-show-ai-output.md) |
| Review outcome (completed vs errors vs run failed) is ambiguous | [R009 — Review exit code and report](R009-review-exit-code-and-report.md) |
| AI command missing or invalid reported obscurely | [R001 — Clear error missing AI command](R001-clear-error-missing-ai-command.md) |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R001](R001-clear-error-missing-ai-command.md) | The system reports a clear, user-facing error when the AI command or alias is missing or invalid before run or review. | ready |
| [R002](R002-success-report-and-exit-zero.md) | The system reports a completion message, iteration count, and timing on success and exits 0. | ready |
| [R003](R003-failure-threshold-report-and-exit-code.md) | The system reports the failure threshold and consecutive failure count when exiting due to failure threshold and uses a distinct exit code. | ready |
| [R004](R004-max-iterations-report-and-exit-code.md) | The system reports iteration count and limit when max iterations are exhausted and uses a distinct exit code. | ready |
| [R005](R005-distinct-exit-code-on-interrupt.md) | The system exits with a distinct code on user interrupt (e.g. SIGINT/SIGTERM). | ready |
| [R006](R006-log-level-and-show-ai-output.md) | The system respects configured log level and show-AI-output setting for what is emitted to the user. | ready |
| [R007](R007-dry-run-shows-assembled-prompt.md) | The system supports a dry-run mode that prints the assembled prompt without invoking the AI. | ready |
| [R008](R008-iteration-statistics.md) | The system reports iteration statistics after a multi-iteration run. | ready |
| [R009](R009-review-exit-code-and-report.md) | The system makes review outcome clear via report and presentation so the user understands whether the review completed, the prompt had errors, or the run failed; exit code semantics follow the review command contract (see prompt review outcome). | ready |
