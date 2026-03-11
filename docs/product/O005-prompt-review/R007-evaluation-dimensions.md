# R007: Evaluation dimensions

**Outcome:** O005 — Prompt Review

## Requirement

The review evaluates the prompt along four dimensions: signal and state, iteration awareness, scope and convergence, and subjective completion criteria; feedback in the report reflects these dimensions.

## Detail

The reviewer assesses the prompt so that users get actionable feedback for Ralph's execution model. The four dimensions are:

1. **Signal and state** — Whether the prompt defines clear success and failure signals that Ralph can detect, and whether statefulness (e.g. filesystem, work-tracking) is compatible with the loop model (fresh process per iteration).
2. **Iteration awareness** — Whether the prompt acknowledges that execution is multi-iteration with a fresh process each time, so the AI can re-read state, emit signals, and avoid assuming a single run.
3. **Scope and convergence** — Whether the task has a defined scope and completion criteria that are checkable in practice, so the loop can converge rather than run indefinitely.
4. **Subjective completion criteria** — When "done" is subjective (e.g. "good enough," "reads well"), whether the prompt includes techniques to escape local optima: variation, creative exploration, or stepping back (e.g. consider alternatives, challenge assumptions), so the AI does not get stuck in small repetitive tweaks.

The report (R002) structures or reflects feedback along these dimensions so the user can see what is strong and what needs improvement. The suggested revision (R003) addresses gaps in these dimensions. The reviewer does not enforce a single prompt style; it evaluates qualities that support Ralph's execution model.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Prompt addresses all dimensions well | Report indicates compliance or minor suggestions; revision may be minimal. |
| Prompt addresses none explicitly | Report still structured by dimension; each dimension gets feedback (e.g. "missing" or "unclear"). |
| Very short prompt | All dimensions may be "needs improvement"; report is still dimensioned. |
| Prompt has strong signals but weak scope | Feedback reflects strength on signal/state and weakness on scope/convergence. |
| Subjective "done" with no escape techniques | Dimension 4 feedback suggests adding variation or stepping-back techniques. |

### Examples

#### Prompt missing success/failure signals

**Input:** User runs review on a prompt that never mentions exit codes, markers, or how to signal done/fail.

**Expected output:** Report includes feedback under "Signal and state" (or equivalent) that success and failure signals are missing or unclear; suggested revision includes example language or structure for signals. Machine-parseable summary can indicate failure or dimension scores.

**Verification:** User can read narrative and understand the gap; revision text addresses it; CI can parse summary if needed.

#### Prompt with good iteration awareness

**Input:** Prompt states "You run in a loop; each run is a new process; re-read the task file and state before continuing."

**Expected output:** Report indicates iteration awareness is present or adequate; other dimensions may still have feedback. Revision may not change this part.

**Verification:** Feedback is dimension-specific; no false negative on iteration awareness.

## Acceptance criteria

- [ ] The review evaluates the prompt along the four dimensions: signal and state, iteration awareness, scope and convergence, and subjective completion criteria.
- [ ] The report (narrative and, where applicable, machine-parseable summary) reflects these dimensions so the user can see what is strong and what needs improvement.
- [ ] The suggested revision addresses gaps in these dimensions where appropriate.
- [ ] The reviewer does not enforce a single template; it evaluates qualities that support Ralph's execution model (per O005 non-outcomes).
- [ ] When a dimension is not addressed by the prompt, the report still provides feedback for that dimension (e.g. "missing" or "unclear").

## Dependencies

- R002 — Report content and format; feedback is delivered via the report.
- R003 — Suggested revision; revision applies the evaluation to suggested edits.
