# R005: Distinct exit code on interrupt

**Outcome:** O004 — Observability

## Requirement

The system exits with a distinct code on user interrupt (e.g. SIGINT/SIGTERM).

## Detail

When the user or system interrupts the process (e.g. Ctrl+C sends SIGINT, or a process manager sends SIGTERM), the process must exit with a distinct exit code so that scripts and CI can tell "interrupted" from success, failure threshold, or max iterations exhausted. The code follows platform convention for interruption (e.g. for SIGINT) where applicable; the important point is that the code is distinct and documented so automation can recognize interruption.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User presses Ctrl+C (SIGINT) during loop | Process exits with a distinct exit code that indicates interruption; no success/failure/exhaustion code |
| Process receives SIGTERM | Process exits with a distinct code for interruption (or documented convention for SIGTERM) |
| Interrupt during review command | Same principle: distinct exit code for interrupt so it is not confused with review success, prompt errors, or review/apply failure |
| Interrupt before first iteration starts | Process exits with interrupt code; user understands the run was stopped by signal |
| Same code as success, failure threshold, or exhaustion | No; the interrupt code is distinct so scripts can branch on "was interrupted" |

### Examples

#### User interrupts with Ctrl+C

**Input:** User invokes the run command; after two iterations they interrupt (e.g. Ctrl+C / SIGINT).

**Expected output:** Process exits promptly with a distinct exit code that indicates interruption. No completion or failure-threshold message; user understands the run was interrupted.

**Verification:** The process exits with the documented interrupt code so scripts can distinguish interruption from success, failure threshold, or exhaustion.

#### SIGTERM from process manager

**Input:** CI or process manager sends SIGTERM during a run.

**Expected output:** The process exits with a distinct code that indicates interruption. Documentation states what code is used for SIGTERM (or platform convention).

**Verification:** Automation can detect that the process was terminated by signal, not by normal success/failure/exhaustion.

## Acceptance criteria

- [ ] When the process receives SIGINT (e.g. Ctrl+C), it exits with a distinct exit code that indicates interruption (e.g. platform convention for SIGINT).
- [ ] When the process receives SIGTERM or other configured interrupt signals, it exits with a distinct code (or documented convention) so that "interrupted" is distinguishable from other outcomes.
- [ ] The interrupt exit code(s) are documented so scripts and CI can recognize interruption.
- [ ] The code is not reused for success, failure threshold, or max iterations exhausted.

## Dependencies

- O010/R002 — Documented stable exit codes include interruption; this requirement defines the observability (distinct exit code) for interrupt so that the full set (success, failure, exhaustion, interruption) is consistent and documented.
