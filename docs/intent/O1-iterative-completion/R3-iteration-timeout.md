# R3: Per-Iteration Timeout

**Outcome:** O1 — Iterative Completion

## Requirement

The system enforces a configurable time limit on each AI CLI process invocation, terminating processes that exceed the limit. The timeout applies independently to each iteration, not to the total loop duration.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] When iteration_timeout is set to a positive value, Ralph kills the AI CLI process if it runs longer than the specified duration in seconds
- [ ] Partial output from a timed-out process is captured and scanned for signals
- [ ] A timed-out iteration counts as one completed iteration toward the max iteration limit
- [ ] When iteration_timeout is 0 or unset, no time limit is enforced
- [ ] The timeout applies to each iteration independently — a 60-second timeout means each iteration gets 60 seconds, not 60 seconds total
