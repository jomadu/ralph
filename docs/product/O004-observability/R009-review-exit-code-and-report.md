# R009: Review exit code and report

**Outcome:** O004 — Observability

## Requirement

The system makes the review outcome clear via report and presentation so the user understands whether the review completed, the prompt had errors, or the run failed; exit code semantics follow the review command contract defined in the prompt review outcome (O005).

## Detail

The review command has distinct outcomes: review completed (with or without findings), the prompt had errors (e.g. structural or validation failures), or the run failed (e.g. AI command or execution failure). The user must be able to tell which case occurred and what to do next. The report content and structure are defined by the prompt review outcome; this requirement covers observability: clarity of outcome and exit code.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Review completes successfully, no errors in prompt | Report indicates completion; exit code per O005 contract (documented success code) |
| Review completes but prompt has errors/findings | Report clearly indicates prompt had errors or findings; user can see what to fix; exit code per O005 (distinct for errors) |
| Review run fails (e.g. AI command missing, crash) | Report or presentation indicates the run failed; user understands it was not a "review completed" or "prompt errors" outcome; exit code distinct per O005 |
| Ambiguity between "prompt errors" and "run failed" | Presentation and exit code distinguish the two so the user knows whether to fix the prompt or fix the environment/command |

### Examples

#### Review completed, prompt has errors

**Input:** User invokes the review command; the prompt has structural or validation errors reported by the review logic.

**Expected output:** Report states or clearly implies that the prompt had errors; user can see what was wrong. Exit code reflects "prompt had errors" per the review contract (O005).

**Verification:** User understands the outcome is "prompt had errors," not "review run failed" or "review completed with no issues."

#### Review run failed

**Input:** User invokes the review command; the AI command is missing (or execution fails).

**Expected output:** Report or message indicates the run failed (e.g. could not complete review). Exit code is distinct from "review completed" and "prompt had errors" per O005.

**Verification:** User understands they need to fix the environment or command, not the prompt content.

## Acceptance criteria

- [ ] The review command produces a report (or equivalent presentation) that makes the outcome clear: review completed, prompt had errors, or run failed.
- [ ] The user can distinguish "review completed" (with or without findings) from "prompt had errors" from "run failed."
- [ ] Exit codes for the review command follow the contract defined in the prompt review outcome (O005) and are distinct for these cases.
- [ ] The user can determine what to do next (e.g. fix prompt vs. fix config/command) from the report and exit code.

## Dependencies

- O005 (Prompt review) defines the review command contract, report content, and exit code semantics; this requirement ensures observability of that contract.
