# R007: Exit on max iterations

**Outcome:** O001 — Iterative Completion

## Requirement

The system exits when the maximum iteration count is reached.

## Detail

The user (or configuration) sets a maximum number of iterations (e.g. 10). The system counts iterations (each AI process run is one iteration). When the count reaches the maximum, the system stops the loop and exits. The exit code is distinct so that "stopped due to max iterations" is distinguishable from "completed successfully" (the documented success code). This bounds the loop so that runaway execution does not occur even if neither success nor failure signal appears, or if the task never converges.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Success signal on iteration N < max | The system exits with the documented success code before reaching max (R004). |
| Max iterations reached without success | The system exits with a distinct code; the user can tell exit was due to max iterations (where documented). |
| Max = 1 | After one iteration, if no success, the system exits (failure or max-iterations exit). |
| Failure threshold reached before max | The system may exit due to R005 before max iterations. |
| Both max and failure threshold reached on same iteration | One exit reason is chosen and reported (e.g. max iterations takes precedence or failure does); behavior is documented. |

### Examples

#### Max iterations reached

**Input:** Max iterations = 3. No success signal in iterations 1, 2, or 3. Failure signal appears each time; failure threshold is 5.

**Expected output:** After iteration 3, the system exits with the documented max-iterations exit code because max iterations was reached.

**Verification:** Exactly 3 AI invocations; exit is with the documented max-iterations code; message or docs indicate max iterations.

#### Success before max

**Input:** Max iterations = 10. Success signal appears in iteration 2.

**Expected output:** The system exits with the documented success code after iteration 2; max iterations is not reached.

**Verification:** Only 2 iterations run; exit is with the documented success code.

## Acceptance criteria

- [ ] The system counts iterations and stops when the count equals the configured maximum.
- [ ] When the loop stops due to max iterations, the system exits with the documented max-iterations exit code.
- [ ] The maximum iteration count is configurable (e.g. default and override).
- [ ] When success is detected before max is reached, the system exits with the documented success code and does not run extra iterations.

## Dependencies

- R004 — Success can cause exit before max.
- R005 — Failure threshold can cause exit before max.
