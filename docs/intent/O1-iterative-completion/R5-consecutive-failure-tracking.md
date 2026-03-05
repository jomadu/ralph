# R5: Consecutive Failure Tracking

**Outcome:** O1 — Iterative Completion

## Requirement

The system tracks consecutive iterations that produce a failure signal and aborts the loop when the count reaches a configurable threshold. A success signal or no-signal iteration resets the counter.

## Specification

Ralph maintains a counter `consecutive_failures`, initialized to 0 before the loop starts.

**After each iteration's signal scanning (R2):**

| Iteration outcome | Counter action | Loop action |
|-------------------|---------------|-------------|
| success | N/A — loop exits 0 | Exit 0 |
| failure | `consecutive_failures += 1` | If `consecutive_failures >= failure_threshold`, exit 1. Otherwise, next iteration. |
| no-signal | `consecutive_failures = 0` | Next iteration (subject to max iteration check per R4). |

The threshold check runs immediately after incrementing the counter on a failure outcome. If the threshold is reached, Ralph exits 1 without starting another iteration.

**Why no-signal resets the counter:**

A no-signal iteration means the AI ran and produced output without asserting success or failure. This typically means the AI is making progress (modifying files, exploring approaches) but isn't ready to commit to a verdict. Resetting the counter on no-signal provides a pressure release — the loop continues as long as the AI is working, and only aborts when the AI repeatedly asserts failure.

A crash with no signal in its output (R1) is also a no-signal outcome, so it resets the counter.

**Configuration:**

- Field: `failure_threshold`
- Type: positive integer (≥ 1)
- Default: `3`

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| `failure_threshold: 1` | First failure signal aborts the loop immediately with exit 1 |
| Alternating failure and no-signal iterations | Counter never exceeds 1; loop continues until max iterations or success |
| Crash (non-zero exit) with no signal in output | Iteration outcome is no-signal (R2); counter resets to 0 |
| Crash (non-zero exit) with failure signal in output | Iteration outcome is failure (R2); counter increments |
| Failure threshold reached on the last allowed iteration | Exit 1 (failure threshold), not exit 2 (exhaustion). Failure threshold is checked after the iteration, before the next iteration's max-iteration check (R4). |

### Examples

#### Three consecutive failures

**Input:**
`failure_threshold: 3`. Iterations 1, 2, and 3 all produce the failure signal.

**Expected output:**
After iteration 1: counter = 1. After iteration 2: counter = 2. After iteration 3: counter = 3. `3 >= 3` → Ralph exits with code 1.

**Verification:**
- Ralph exit code is 1
- 3 iterations executed

#### Failure streak broken by no-signal

**Input:**
`failure_threshold: 3`, `default_max_iterations: 5`. Iteration sequence: failure, failure, no-signal, failure, failure.

**Expected output:**
Counter progression: 1, 2, 0, 1, 2. Counter never reaches 3. All 5 iterations execute. Ralph exits with code 2 (exhaustion per R4).

**Verification:**
- Ralph exit code is 2 (not 1)
- 5 iterations executed

## Acceptance criteria

- [ ] Each iteration producing a failure signal increments the consecutive failure counter
- [ ] A success signal resets the consecutive failure counter to zero (though the loop stops on success anyway)
- [ ] A no-signal iteration resets the consecutive failure counter to zero
- [ ] When the counter reaches the failure threshold, Ralph aborts and exits with code 1
- [ ] The default failure threshold is 3
- [ ] The threshold is configurable through the standard configuration hierarchy

## Dependencies

_None identified._
