# R008: Exit codes

**Outcome:** O005 — Prompt Review

## Requirement

Exit codes distinguish: (1) review completed with no errors, (2) review completed but the prompt has errors, and (3) review or apply did not complete successfully.

## Detail

Scripts and CI need to gate on review outcome without parsing the report. The system uses distinct exit codes (or a documented scheme) so that:

- **The documented success code:** Review completed; report and suggested revision were produced; the prompt has no reported errors (or only minor/acceptable issues per product definition). Apply, if requested, completed if applicable.
- **The documented prompt-errors code:** Review completed; report and suggested revision were produced; the prompt has one or more errors or serious issues according to the evaluation dimensions. Apply, if requested, may still have been performed; the code indicates "review succeeded but prompt has problems."
- **The documented failure code:** Review or apply did not complete successfully — e.g. invalid input source (R001), report write failed (R005), prompt supplied via standard input + apply without revision output path (R006), confirmation required but not given in non-interactive session, or reviewer/internal error. No guarantee that report or revision is complete or written.

The exact numeric values and their mapping are documented so that scripts can rely on them. This requirement states the three-way distinction for the review command.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Review completes, prompt has no issues | The documented success code is used. |
| Review completes, prompt has issues | The documented prompt-errors code is used. |
| Invalid alias or missing file (R001) | The documented failure code is used; no report. |
| Report path unwritable (R005) | The documented failure code is used. |
| Prompt supplied via standard input + apply without path (R006) | The documented failure code is used; do not apply. |
| User declines confirmation | Success or a dedicated code may be used; apply was not performed; behavior is defined. |
| Reviewer crashes or times out | The documented failure code is used. |

### Examples

#### Success — no prompt errors

**Input:** User runs the review command with the prompt supplied by file path; prompt is compliant; report and revision are produced; report is written to file.

**Expected output:** The documented success code is used. Report file exists; revision is available. CI can treat as pass.

**Verification:** Script or user can observe the documented success code; report content confirms no errors.

#### Prompt has errors

**Input:** User runs the review command with the prompt supplied by file path; prompt lacks success/failure signals; review completes; report and revision are produced.

**Expected output:** The documented prompt-errors code is used. Report file exists and describes the issues; revision suggests fixes. CI can treat as "prompt needs work."

**Verification:** Exit code is distinct from the documented success code and from the documented failure code; report explains issues.

#### Failure — prompt supplied via standard input + apply without path

**Input:** User runs the review command with the prompt supplied via standard input and requests apply, with no revision output path.

**Expected output:** The documented failure code is used. Error message that the revision output path is required. No revision file written.

**Verification:** The documented failure code is used; script can distinguish from success and prompt-errors; no silent success.

## Acceptance criteria

- [ ] Three outcomes are distinguishable by exit code: (1) review completed with no errors, (2) review completed but prompt has errors, (3) review or apply did not complete successfully.
- [ ] The mapping from exit code to outcome is documented so scripts and CI can gate reliably.
- [ ] When the review does not complete (invalid input, write failure, missing revision output path when prompt was from standard input + apply, etc.), the documented failure code is used.
- [ ] When the review completes and produces report and revision, the exit code is either the documented success code or the documented prompt-errors code depending on whether the prompt has reported errors; the documented failure code is not used in that case.

## Dependencies

- R001, R005, R006 — Failure conditions (invalid input, report write fail, prompt from standard input + apply with no path) result in the documented failure code.
- R002, R003, R005 — When review "completes," report and revision are produced; exit code then reflects success vs prompt errors.
