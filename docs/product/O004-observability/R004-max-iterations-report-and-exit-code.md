# R004: Max iterations report and exit code

**Outcome:** O004 — Observability

## Requirement

The system reports iteration count and limit when max iterations are exhausted and uses a distinct exit code.

## Detail

When the loop stops because the maximum iteration count was reached without detecting the success signal, the user must see how many iterations ran and what the limit was. A distinct exit code for exhaustion allows scripts and CI to distinguish "exhausted" from success, failure threshold, or interrupt.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Max iterations is 5; fifth iteration completes without success signal | Report shows iteration count (5) and limit (5); process exits with distinct exit code for exhaustion |
| Max iterations is 1; first iteration does not produce success signal | Report shows count 1 and limit 1; process exits with distinct exit code for exhaustion |
| User has quiet or minimal log level | Essential report (count and limit) still visible; exit code unchanged |
| Same exit code as failure threshold or interrupt | No; the exit code for max-iterations exhausted is distinct from those so scripts can tell exhaustion from failure or interrupt |

### Examples

#### Exiting on max iterations

**Input:** Config sets max iterations to 3; three iterations run and none produces the success signal.

**Expected output:** Message or log line indicating exit due to max iterations, the iteration count (3), and the limit (3). The process exits with a distinct exit code for exhaustion (not the success code or the failure-threshold code).

**Verification:** User understands the loop stopped because the iteration limit was reached; scripts can branch on exit code (e.g. exhaustion vs failure).

#### Single iteration limit

**Input:** Config sets max iterations to 1; the first iteration does not emit the success signal.

**Expected output:** Report shows count 1 and limit 1; process exits with distinct exit code for exhaustion.

**Verification:** The user understands the run exhausted the limit without success; scripts can branch on the distinct exhaustion exit code.

## Acceptance criteria

- [ ] When the loop exits because the maximum iteration count was reached without success, the system prints or logs the iteration count and the configured limit.
- [ ] The process exits with a distinct exit code for exhaustion that is not used for success, failure threshold, or interrupt.
- [ ] The exit code is documented or consistent so scripts can rely on it for exhaustion.
- [ ] The user can tell that the loop stopped due to exhaustion, not failure threshold or interrupt.

## Dependencies

- O001/R007 — Exit on max iterations defines when the loop stops; this requirement covers the observability (report and exit code) of that outcome.
