# R002: Report content and format

**Outcome:** O005 — Prompt Review

## Requirement

The system produces a report with narrative feedback and a machine-parseable summary for every completed review.

## Detail

Every successful review run produces a report. The report contains (1) **narrative feedback** — human-readable commentary on the prompt's quality and structure — and (2) **a machine-parseable summary** so that scripts or CI can gate on the result (e.g. pass/fail, severity, or structured result). The format of the machine-parseable portion is defined so that consumers can reliably parse it; the narrative may be free-form or structured. When the review does not complete successfully (e.g. I/O error, reviewer failure), no report is produced and the documented failure code is used per R008.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Review completes successfully | Report is produced and written to file (R005); it contains both narrative and machine-parseable summary. |
| Review fails before producing output | No report file is required; the documented failure code is used (R008). |
| Prompt has no issues | Report still produced; narrative and summary can indicate "no issues" or "pass." |
| Prompt has severe issues | Summary allows CI/scripts to detect failure or severity; narrative explains issues. |
| Very long prompt | Report remains bounded and usable; summary is compact; narrative may be sectioned. |

### Examples

#### Report for prompt with issues

**Input:** User runs review on a prompt that lacks explicit success/failure signals. Report is written to default or user-chosen path.

**Expected output:** Report file contains narrative feedback (e.g. "The prompt does not specify how the AI should signal success or failure…") and a machine-parseable summary (e.g. a line or block with a result code, or structured fields such as `result: fail` or `signals: missing`).

**Verification:** Human can read narrative and understand the issue; a script can parse the summary and gate (e.g. use the documented prompt-errors or failure code when result is fail).

#### Report for compliant prompt

**Input:** User runs review on a prompt that meets the evaluation dimensions. Report is written to file.

**Expected output:** Report contains narrative indicating compliance or minor suggestions, and machine-parseable summary indicating pass or equivalent.

**Verification:** CI can parse summary and pass; user can read narrative for confirmation.

## Acceptance criteria

- [ ] Every completed review produces a report that includes both narrative feedback and a machine-parseable summary.
- [ ] The machine-parseable summary is documented or stable enough for scripts and CI to parse and gate on the result.
- [ ] The report is written to a file per R005 (path user-chosen or default).
- [ ] When the review does not complete successfully, a report is not required; the documented failure code is used per R008.

## Dependencies

- R005 — Report is written to a file (report output path and persistence).
- R008 — Exit codes (so "completed" vs "did not complete" is well-defined).
