# R009: Process exit without signal

**Outcome:** O001 — Iterative Completion

## Requirement

When the AI process exits without emitting the configured success or failure signal (e.g. crash, kill, abnormal exit), the system treats the iteration as a failure, increments the consecutive-failure count, and continues or exits according to the failure threshold; the user can distinguish this condition from signal-based failure where documented.

## Detail

The AI process may exit without producing the success or failure signal — for example it crashes, is killed, times out, or exits without the configured markers. The system treats such an iteration as a failure: it increments the consecutive-failure count and either starts the next iteration or exits based on the failure threshold (same as R005). Thus the loop remains bounded and predictable. Where documented, the user can distinguish "no signal" from "failure signal present" in reporting (e.g. "exited due to process crash" vs "exited due to failure signal threshold") so they can debug or adjust behavior.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Process exits 0 but no success signal in output | Treat as failure; increment count; continue or exit per threshold. |
| Process crashes (e.g. segfault, kill -9) | Treat as failure; increment count; continue or exit per threshold. |
| Process times out (if applicable) | Treat as failure; increment count; continue or exit per threshold. |
| Process exits non-zero, no failure signal in output | Treat as failure; increment count. |
| Distinguishability | Where documented, user can tell "no signal" from "failure signal" (e.g. in logs or exit reason). |

### Examples

#### Process crash mid-run

**Input:** AI process is killed (e.g. SIGKILL) before emitting any signal. Failure threshold = 2.

**Expected output:** The system treats the iteration as failure, increments count to 1, starts the next iteration (or exits if threshold is 1).

**Verification:** Consecutive-failure count increases; loop continues or exits per threshold; no hang.

#### Exit without signal; user can distinguish

**Input:** The AI process exits without emitting the success or failure signal (e.g. non-success exit). The system is configured to report exit reason.

**Expected output:** The system treats the iteration as failure, increments count. Where documented, the report indicates that the process exited without emitting the configured signal (e.g. "process exited without success/failure signal").

**Verification:** Behavior same as failure for loop purposes; where documented, the reason is distinguishable from signal-based failure.

## Acceptance criteria

- [ ] When the AI process exits without the configured success or failure signal appearing in its output, the system treats the iteration as a failure.
- [ ] The system increments the consecutive-failure count for such iterations and applies the same continue/exit logic as for failure-signal (R005).
- [ ] Where documented, the user can distinguish "process exited without signal" from "failure signal detected" (e.g. for debugging or reporting).
- [ ] No iteration is left undefined; "no signal" always maps to failure for the purpose of the loop.

## Dependencies

- R005 — Same failure count and threshold logic; R009 covers the "no signal" case.
