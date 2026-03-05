# R7: Graceful Interruption Handling

**Outcome:** O1 — Iterative Completion

## Requirement

The system handles SIGINT and SIGTERM by terminating the current AI CLI process, waiting for it to exit, and then exiting with code 130. The interruption is clean — no partial iteration results are processed.

## Specification

_To be specified._

### Edge cases

_To be specified._

### Examples

_To be specified._

## Acceptance criteria

- [ ] On SIGINT or SIGTERM, Ralph sends a termination signal to the running AI CLI process
- [ ] Ralph waits for the AI CLI process to exit with a bounded timeout before forcing termination
- [ ] Ralph exits with code 130 after handling the interruption
- [ ] If no AI CLI process is running at the time of the signal (e.g., between iterations), Ralph exits immediately with code 130
- [ ] A second SIGINT/SIGTERM during the wait forces immediate exit

## Dependencies

_None identified._
