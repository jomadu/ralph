# R003: Failure threshold report and exit code

**Outcome:** O004 — Observability

## Requirement

The system reports the failure threshold and consecutive failure count when exiting due to the failure threshold and uses a distinct exit code.

## Detail

When the loop stops because the configured number of consecutive failures was reached, the user must see the threshold value and how many consecutive failures occurred. A distinct exit code that indicates failure threshold allows scripts and CI to distinguish "failed due to threshold" from success or other outcomes.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Failure threshold is 3; third consecutive failure | Report shows threshold (3) and consecutive failure count (3); process exits with distinct exit code for failure threshold |
| Failure threshold 1; first failure | Report shows threshold 1 and count 1; process exits with distinct exit code for failure threshold |
| Interleaved success resets counter | When threshold is hit, reported count matches consecutive failures; process exits with distinct exit code for failure threshold |

### Examples

#### Exiting on failure threshold

**Input:** Config sets failure threshold to 2; two consecutive iterations produce a failure signal.

**Expected output:** Message or log line indicating exit due to failure threshold, the threshold value (2), and the consecutive failure count (2). The process exits with a distinct exit code that indicates failure threshold (not the documented success code).

**Verification:** User understands the loop stopped because the failure threshold was reached; scripts can branch on exit code.

## Acceptance criteria

- [ ] When the loop exits due to the failure threshold, the system prints or logs the configured failure threshold value.
- [ ] The system reports the consecutive failure count that triggered the exit.
- [ ] The process exits with a distinct exit code that indicates failure threshold, not used for success or other defined outcomes (e.g. max iterations, interrupt).
- [ ] The exit code is documented or consistent so scripts can rely on it.
