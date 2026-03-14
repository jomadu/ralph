# R002: Success report and exit zero

**Outcome:** O004 — Observability

## Requirement

The system reports a completion message, iteration count, and timing on success and exits with the documented success code.

## Detail

When the loop stops because the configured success signal was detected in the AI output, the user must see that the run completed successfully and understand how long it took and how many iterations ran. A completion message (e.g. to stdout or logs) and iteration count give the user observable evidence of success. The process exits with the documented success code so scripts and CI can treat the run as successful without parsing output.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Success signal detected on first iteration | Completion message; iteration count 1; timing for that iteration; process exits with documented success code |
| Success signal detected after multiple iterations | Completion message; total iteration count; timing (e.g. total or per-iteration summary); process exits with documented success code |
| User has quiet or minimal log level | Completion message and essential success info still visible (e.g. completion + count); process exits with documented success code |
| Loop and review both use success code for success | Run command: success exits with documented success code per this requirement; review command success is defined in O005 |

### Examples

#### Single-iteration success

**Input:** User invokes the run command; success signal is detected in the first iteration's output.

**Expected output:** A completion message (e.g. "Completed successfully" or equivalent), the iteration count (1), and timing (e.g. duration of the iteration). The process exits with the documented success code.

**Verification:** User sees that the run succeeded and how long it took; the process exits with the documented success code so automation can detect success.

#### Multi-iteration success

**Input:** User invokes the run command; success signal is detected on the third iteration.

**Expected output:** Completion message, iteration count (3), and timing (e.g. total elapsed or per-iteration). The process exits with the documented success code.

**Verification:** User understands the loop ran 3 iterations before success; scripts can rely on the documented success code for success.

## Acceptance criteria

- [ ] When the loop exits because the success signal was detected, the system prints or logs a completion message.
- [ ] The system reports the iteration count (number of iterations run before success).
- [ ] The system reports timing (e.g. elapsed time for the run or iteration summary) so the user can see how the run performed.
- [ ] The process exits with the documented success code on success so scripts and CI can branch on outcome.
- [ ] The success exit code is documented or consistent so automation can rely on it.

## Dependencies

- O001/R004 — Success is defined as detection of the configured success signal and exit with the documented success code; this requirement covers the observability of that outcome (message, count, timing).
