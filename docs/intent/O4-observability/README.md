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

## Obstacles

| Obstacle | Mitigating Requirement |
|----------|----------------------|
| User scripts around Ralph can't distinguish termination reasons programmatically | R1 — Distinct exit codes |
| User doesn't know how long iterations took or whether performance is degrading | R2 — Iteration statistics |
| User wants to see what the AI is doing during execution | R3 — Verbose output streaming |
| User wants to validate prompt assembly before committing to a long run | R4 — Dry-run mode |
| Normal operation produces too much output and obscures the summary | R5 — Log level control |
| User wants to suppress all non-error output for scripting | R5 — Log level control |
| User doesn't know what iteration the loop is on during execution | R6 — Iteration progress reporting |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| R1 | Distinct exit codes | draft |
| R2 | Iteration statistics | draft |
| R3 | Verbose output streaming | draft |
| R4 | Dry-run mode | draft |
| R5 | Log level control | draft |
| R6 | Iteration progress reporting | draft |
