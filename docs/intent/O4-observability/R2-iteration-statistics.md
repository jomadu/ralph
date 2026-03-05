# R2: Iteration Statistics

**Outcome:** O4 — Observability

## Requirement

The system reports iteration timing statistics at loop completion, giving the user data to understand loop performance and tune configuration. Statistics include iteration count, min/max/mean duration, and standard deviation.

## Specification

_To be specified._

### Edge cases

_To be specified._

### Examples

_To be specified._

## Acceptance criteria

- [ ] After the loop completes (success, failure threshold, or max iterations), Ralph reports the total number of iterations executed
- [ ] Min, max, and mean iteration durations are reported
- [ ] Standard deviation of iteration durations is reported, calculated using Welford's online algorithm
- [ ] Statistics are reported to stderr so they do not interfere with stdout piping or capture
- [ ] Statistics are not reported on interruption (SIGINT/SIGTERM — exit 130)
- [ ] For single-iteration runs, standard deviation is reported as 0 or omitted

## Dependencies

_None identified._
