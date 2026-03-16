# R005: Detect failure signal and continue or exit

**Outcome:** O001 — Iterative Completion

## Requirement

The system detects the configured failure signal, increments the consecutive-failure count, and either starts a new iteration or exits based on the failure threshold.

## Detail

The user configures a failure signal (e.g. a string or pattern) that indicates the AI did not yet achieve the task. The system captures the AI output (per configuration) and scans **only the last non-empty line** of that output for the configured failure signal. The exact definition of "last non-empty line" is given in the [run-loop component spec](../../engineering/components/run-loop.md) (split on newline, trim each line, take the last line that is non-empty after trim; if there is no non-empty line, no failure signal is detected). When the system finds the failure signal on that line (and precedence, if applicable, resolves to failure), it increments a consecutive-failure counter. If the counter is below the configured failure threshold, the system starts a new iteration. If the counter reaches the threshold, the system exits with a distinct exit code that indicates the failure threshold was reached (and may report that the threshold was reached). The counter is reset when a success signal is detected (iteration succeeds). Configuration defines the threshold (e.g. 1 = exit on first failure; N = allow N consecutive failures before exit).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Failure signal appears on the last non-empty line | The system detects it, increments consecutive-failure count, and continues or exits per threshold. |
| Failure signal appears only on an earlier line (not the last non-empty line) | The system does not treat the iteration as failure on that basis; no increment from that output; continues per other rules (e.g. R009 if no signal on last line). |
| Failure signal detected; count below threshold | Increment count; start next iteration. |
| Failure signal detected; count at threshold | Increment count (or after increment equals threshold); exit with the documented failure-threshold exit code. |
| Success signal detected on a later iteration | Consecutive-failure count resets; next failure starts count from 1. |
| Both success and failure in same output (on that line) | Precedence (R006 or R008) decides; if failure wins, R005 applies. |
| Max iterations reached before threshold | R007 applies; exit due to max iterations. |
| Threshold is 0 or 1 | First failure causes exit (or documented special behavior). |

### Examples

#### Failure signal on the last non-empty line

**Input:** Failure signal = `FAIL`. The AI outputs "Still working...\nFAIL" (so "FAIL" is on the last non-empty line). The system runs.

**Expected output:** The system detects the failure signal on the last non-empty line, increments consecutive-failure count, and continues or exits per threshold.

**Verification:** Consecutive-failure count increments; behavior matches threshold configuration.

#### Failure signal only on an earlier line

**Input:** Failure signal = `FAIL`. The AI outputs "FAIL\nStill working..." (so "FAIL" appears only on the first line; the last non-empty line is "Still working...").

**Expected output:** The system does not treat the iteration as failure on that basis (the last non-empty line is scanned, and it does not contain the failure signal); no increment from that output. If no success signal on that line, R009 (process exit without signal) or other rules may apply.

**Verification:** No failure-signal increment from that iteration; behavior consistent with last-line-only scanning.

#### Failure then retry

**Input:** Failure threshold = 3. Failure signal = `FAIL`. Iteration 1 output contains `FAIL` on the last non-empty line. The system runs.

**Expected output:** The system increments consecutive-failure count to 1, starts iteration 2 (without exiting).

**Verification:** Second iteration runs; exit code not yet set; count is 1.

#### Failure threshold reached

**Input:** Failure threshold = 2. Iterations 1 and 2 both emit failure signal. The system runs.

**Expected output:** After iteration 2, consecutive-failure count reaches 2; the system exits with the documented failure-threshold exit code (and may report threshold reached).

**Verification:** Exit is with the documented failure-threshold code; the user can see that exit was due to failure threshold (where documented).

#### Success resets count

**Input:** Iteration 1 fails (failure signal); iteration 2 succeeds (success signal). Iteration 3 fails again.

**Expected output:** After iteration 2, count resets. After iteration 3, count is 1; if threshold > 1, the system starts iteration 4.

**Verification:** Consecutive-failure count does not carry across a successful iteration.

## Acceptance criteria

- [ ] The system scans only the last non-empty line of the captured AI output for the configured failure signal; earlier lines are not used for failure detection (see run-loop component spec for definition of last non-empty line).
- [ ] When the failure signal is detected on that line (and precedence resolves to failure), the system increments the consecutive-failure count.
- [ ] If the count is below the configured failure threshold, the system starts a new iteration.
- [ ] If the count reaches the failure threshold, the system exits with the documented failure-threshold exit code.
- [ ] When a success signal is detected, the consecutive-failure count is reset to zero.
- [ ] Threshold and failure signal are configurable.

## Dependencies

- R004 — Success detection resets the failure count.
- R006 (and optionally R008) — Precedence when both signals present; R005 applies when outcome is failure.
