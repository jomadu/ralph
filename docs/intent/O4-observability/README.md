# O4: Observability

## Statement

The user knows why the loop stopped and how it performed.

## Why it matters

A loop runner that silently exits is a black box. The user re-runs the command, checks file diffs, reads AI output, and tries to reconstruct what happened. Did it succeed? Did it hit the failure threshold? Did it time out? Did the signal never appear? Without clear termination reporting and execution statistics, diagnosing loop behavior is guesswork and tuning config is trial and error.

## Verification

- On success signal: Ralph prints a completion message, reports iteration count and timing, exits 0.
- On failure threshold: Ralph prints the threshold value and consecutive failure count, exits 1.
- On max iterations exhausted: Ralph prints the iteration count and limit, exits 2.
- On SIGINT/SIGTERM: Ralph exits 130.
- User runs `ralph run build -v` and sees the AI's output streamed to the terminal in real time while Ralph still captures it for signal scanning.
- User runs `ralph run build -d` and sees the fully assembled prompt (preamble + prompt file) printed to stdout without any AI process being spawned.
- After a multi-iteration run, Ralph reports min/max/mean/stddev of iteration durations.

## Non-outcomes

- Ralph does not provide persistent logging, log files, or structured log output (JSON lines, etc.). Output goes to stderr/stdout for the current invocation only.
- Ralph does not provide per-iteration diffs or file change tracking between iterations.
- Ralph does not integrate with external monitoring, alerting, or observability systems.
- Ralph does not support replaying or debugging past executions.

## Risks

| Risk | Mitigating Requirement |
|----------|----------------------|
| User scripts around Ralph can't distinguish termination reasons programmatically | [R1 — Distinct exit codes](R1-exit-codes.md) |
| User doesn't know how long iterations took or whether performance is degrading | [R2 — Iteration statistics](R2-iteration-statistics.md) |
| User wants to see what the AI is doing during execution | [R3 — Verbose output streaming](R3-verbose-streaming.md) |
| User wants to validate prompt assembly before committing to a long run | [R4 — Dry-run mode](R4-dry-run.md) |
| Normal operation produces too much output and obscures the summary | [R5 — Log level control](R5-log-level-control.md) |
| User wants to suppress all non-error output for scripting | [R5 — Log level control](R5-log-level-control.md) |
| User doesn't know what iteration the loop is on during execution | [R6 — Iteration progress reporting](R6-iteration-progress.md) |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R1](R1-exit-codes.md) | Distinct exit codes | ready |
| [R2](R2-iteration-statistics.md) | Iteration statistics | ready |
| [R3](R3-verbose-streaming.md) | Verbose output streaming | ready |
| [R4](R4-dry-run.md) | Dry-run mode | ready |
| [R5](R5-log-level-control.md) | Log level control | ready |
| [R6](R6-iteration-progress.md) | Iteration progress reporting | ready |
