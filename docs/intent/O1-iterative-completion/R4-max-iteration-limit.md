# R4: Maximum Iteration Limit

**Outcome:** O1 — Iterative Completion

## Requirement

The system enforces a configurable upper bound on the number of iterations, stopping the loop when the limit is reached without a success signal. An unlimited mode disables the limit, running until a signal or failure threshold is reached.

## Specification

_To be specified._

### Edge cases

_To be specified._

### Examples

_To be specified._

## Acceptance criteria

- [ ] The loop executes at most the configured number of iterations
- [ ] When the limit is reached without a success signal, Ralph exits with code 2
- [ ] The default limit is 5 iterations
- [ ] The limit is configurable through the standard configuration hierarchy
- [ ] Unlimited mode (--unlimited / -u) disables the iteration limit, running until a success signal, failure threshold, or interruption
- [ ] The iteration limit must be at least 1

## Dependencies

_None identified._
