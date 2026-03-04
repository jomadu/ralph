# R6: Iteration Progress Reporting

**Outcome:** O4 — Observability

## Requirement

The system reports the current iteration number at the start of each iteration so the user knows where the loop is in its run. Progress messages are visible during normal operation and suppressed only when the user explicitly reduces verbosity.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] At the start of each iteration, Ralph prints the iteration number and the limit (e.g., "Iteration 3/10" or "Iteration 3 (unlimited)")
- [ ] Progress messages are emitted at info log level
- [ ] Progress messages go to stderr so they do not interfere with AI output capture or stdout piping
- [ ] Progress messages are suppressed when log level is set above info (--quiet or --log-level warn/error)
