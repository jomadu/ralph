# R006: Signal precedence

**Outcome:** O001 — Iterative Completion

## Requirement

The system applies a defined precedence when both success and failure signals are present on the same line (default behavior).

## Detail

Ralph scans **only the last non-empty line** of each iteration's output for success and failure signals (see the run-loop component spec for the definition of "last non-empty line"). When both the configured success signal and the configured failure signal appear **on that same line** (the last non-empty line), the outcome is ambiguous without a rule. The system applies a defined precedence (e.g. "success wins" or "failure wins" or "first occurrence wins") so that the iteration is classified as either success or failure, not both. This is the default, static behavior. Signals that appear only on earlier lines are not used for detection; precedence applies only when both signals are present on the last non-empty line. Optional AI-interpreted precedence (R008) can override this for a run when the user enables it; when that option is off or when the interpretation step does not yield a clear answer, this requirement's precedence is used.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Only success signal present (on last non-empty line) | Iteration is success (R004). |
| Only failure signal present (on last non-empty line) | Iteration is failure (R005). |
| Both present on the last non-empty line | Apply defined precedence; iteration is either success or failure. |
| Both present only on earlier lines (not on last non-empty line) | No signal detected for the iteration; R009 applies (process exit without signal). |
| Neither present (on last non-empty line) | R009 applies (process exit without signal). |
| Precedence "success wins" | Both on that line → treat as success. |
| Precedence "failure wins" | Both on that line → treat as failure. |
| Order-dependent rule (e.g. first wins) | Document and apply consistently for that line. |

### Examples

#### Success wins (example policy)

**Input:** Precedence = success wins. The last non-empty line of the output contains both "DONE" and "FAIL".

**Expected output:** The system treats the iteration as success; exits with the documented success code per R004.

**Verification:** Exit is with the documented success code; no new iteration started.

#### Failure wins (example policy)

**Input:** Precedence = failure wins. The last non-empty line of the output contains both "DONE" and "FAIL".

**Expected output:** The system treats the iteration as failure; increments count and continues or exits per R005.

**Verification:** Exit is with the documented failure-threshold code if threshold reached; or the next iteration runs.

## Acceptance criteria

- [ ] When both success and failure signals appear on the last non-empty line of the iteration output, the system classifies the iteration as either success or failure according to the defined precedence rule.
- [ ] R006 states that "both present" means both on the last non-empty line (same line); signals only on earlier lines are not used for precedence.
- [ ] The precedence rule is documented (e.g. success wins, failure wins, or order-based).
- [ ] No iteration is left ambiguous (both success and failure); exactly one outcome is used for R004/R005/R009.
- [ ] This behavior is the default when AI-interpreted precedence (R008) is not used or does not yield a clear result.

## Dependencies

- R004, R005 — Precedence decides which of these applies when both signals are present.
