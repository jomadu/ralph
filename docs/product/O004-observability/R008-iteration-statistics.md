# R008: Iteration statistics

**Outcome:** O004 — Observability

## Requirement

The system reports iteration statistics after a multi-iteration run.

## Detail

After a run that executed more than one iteration, the user can see how the run performed—e.g. min/max/mean duration per iteration or total duration and count. This supports tuning and diagnosis. Single-iteration runs may report simplified or equivalent timing (e.g. single duration) consistent with "how it performed."

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Two or more iterations | Statistics reported (e.g. min/max/mean duration, or total time and count); format is implementation-defined but user-observable |
| Single iteration | At least timing for that run is reported (may be part of R002 success report or a single-iteration summary); no requirement for min/max/mean when N=1 |
| Run exits on failure threshold or max iterations | Statistics cover the iterations that ran before exit |
| Quiet mode | Statistics are part of the essential outcome; they are still reported (e.g. in summary) per O004 verification |

### Examples

#### Multi-iteration run

**Input:** Run completes successfully after five iterations with varying iteration durations.

**Expected output:** After the completion message, the user sees iteration statistics—e.g. min, max, and mean duration per iteration, or total duration and iteration count—so they can assess performance.

**Verification:** User can see how long iterations took and use that to tune or diagnose.

## Acceptance criteria

- [ ] After a run with two or more iterations, the system reports iteration statistics (e.g. min/max/mean duration or equivalent).
- [ ] The statistics are visible to the user (e.g. in logs or summary) and support understanding how the run performed.
- [ ] Single-iteration runs report at least the timing for that run (may be combined with R002 success report).
- [ ] Statistics are reported for runs that exit on failure threshold or max iterations, covering the iterations that executed.
