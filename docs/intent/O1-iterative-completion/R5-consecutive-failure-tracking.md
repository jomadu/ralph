# R5: Consecutive Failure Tracking

**Outcome:** O1 — Iterative Completion

## Requirement

The system tracks consecutive iterations that produce a failure signal and aborts the loop when the count reaches a configurable threshold. A success signal or no-signal iteration resets the counter.

## Specification

_To be specified._

### Edge cases

_To be specified._

### Examples

_To be specified._

## Acceptance criteria

- [ ] Each iteration producing a failure signal increments the consecutive failure counter
- [ ] A success signal resets the consecutive failure counter to zero (though the loop stops on success anyway)
- [ ] A no-signal iteration resets the consecutive failure counter to zero
- [ ] When the counter reaches the failure threshold, Ralph aborts and exits with code 1
- [ ] The default failure threshold is 3
- [ ] The threshold is configurable through the standard configuration hierarchy

## Dependencies

_None identified._
