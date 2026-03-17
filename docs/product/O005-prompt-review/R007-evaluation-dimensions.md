# R007: Evaluation dimensions

**Outcome:** O005 — Prompt Review

## Requirement

The review evaluates the prompt along four dimensions: signal and state, iteration awareness, scope and convergence, and subjective completion criteria; feedback in the report reflects these dimensions.

## Detail

The reviewer assesses the prompt so that users get actionable feedback for Ralph's execution model. The four dimensions are:

1. **Signal and state** — Whether the prompt defines clear success and failure signals that Ralph can detect, and whether statefulness (e.g. filesystem, work-tracking) is compatible with the loop model (fresh process per iteration). Whether the prompt instructs the AI to emit the success or failure marker **on the last line** of its response (so Ralph’s last-line-only detection correctly treats it as the outcome). Review feedback should note if this is missing or unclear.
2. **Iteration awareness** — Ralph injects a preamble that explains multi-iteration execution and a fresh process each time, so the prompt does not need to repeat that. The review checks: does the prompt avoid prescribing behavior by iteration or pass count (to avoid iteration artifacts)? Is statefulness compatible with re-reading from disk each run?
3. **Scope and convergence** — Whether the task has a defined scope and completion criteria that are checkable in practice, so the loop can converge rather than run indefinitely.
4. **Subjective completion criteria** — When "done" is subjective (e.g. "good enough," "reads well"), whether the prompt includes techniques to escape local optima: variation, creative exploration, or stepping back (e.g. consider alternatives and the existing version, then pick the best). These should be phrased as per-run behavior (e.g. "consider two alternatives and the existing structure; pick the best, which may be keeping the current one") rather than conditional on pass or iteration count, to avoid iteration artifacts and unnecessary churn.

The report (R002) structures or reflects feedback along these dimensions so the user can see what is strong and what needs improvement. The suggested revision (R003) addresses gaps in these dimensions. The reviewer does not enforce a single prompt style; it evaluates qualities that support Ralph's execution model.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Prompt addresses all dimensions well | Report indicates compliance or minor suggestions; revision may be minimal. |
| Prompt addresses none explicitly | Report still structured by dimension; each dimension gets feedback (e.g. "missing" or "unclear"). |
| Very short prompt | All dimensions may be "needs improvement"; report is still dimensioned. |
| Prompt has strong signals but weak scope | Feedback reflects strength on signal/state and weakness on scope/convergence. |
| Prompt has clear signals but does not say to put them on the last line | Feedback under "Signal and state" suggests adding that the AI should emit the outcome on the last line (so Ralph’s last-line-only detection treats it as the final outcome). |
| Subjective "done" with no escape techniques | Dimension 4 feedback suggests adding variation or stepping-back techniques. |

### Examples

#### Prompt missing success/failure signals

**Input:** User runs review on a prompt that never mentions exit codes, markers, or how to signal done/fail.

**Expected output:** Report includes feedback under "Signal and state" (or equivalent) that success and failure signals are missing or unclear; suggested revision includes example language or structure for signals. Machine-parseable summary can indicate failure or dimension scores.

**Verification:** User can read narrative and understand the gap; revision text addresses it; CI can parse summary if needed.

#### Prompt with good iteration awareness

**Input:** Prompt avoids pass-count or iteration-dependent behavior and assumes state is re-read each run (e.g. "Re-read state at start of each run; emit success or failure on the last line"). It need not explain the loop — the preamble does that.

**Expected output:** Report indicates iteration awareness is present or adequate; other dimensions may still have feedback. Revision may not change this part.

**Verification:** Feedback is dimension-specific; no false negative on iteration awareness.

#### Prompt has clear signals but does not say to put them on the last line

**Input:** User runs review on a prompt that defines success/failure markers (e.g. `<promise>SUCCESS</promise>`) but never instructs the AI to emit the outcome on the last line of its response.

**Expected output:** Report includes feedback under "Signal and state" that the prompt should instruct the AI to emit the outcome signal on the last line (so Ralph’s last-line-only detection correctly treats it as the outcome). Suggested revision adds or implies "emit the outcome on the last line" (or equivalent).

**Verification:** User can read narrative and understand the gap; revision text addresses it.

## Acceptance criteria

- [ ] The review evaluates the prompt along the four dimensions: signal and state, iteration awareness, scope and convergence, and subjective completion criteria.
- [ ] The report (narrative and, where applicable, machine-parseable summary) reflects these dimensions so the user can see what is strong and what needs improvement.
- [ ] The suggested revision addresses gaps in these dimensions where appropriate.
- [ ] The reviewer does not enforce a single template; it evaluates qualities that support Ralph's execution model (per O005 non-outcomes).
- [ ] When a dimension is not addressed by the prompt, the report still provides feedback for that dimension (e.g. "missing" or "unclear").

## Dependencies

- R002 — Report content and format; feedback is delivered via the report.
- R003 — Suggested revision; revision applies the evaluation to suggested edits.
