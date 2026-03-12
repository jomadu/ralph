# R002: Report content and format

**Outcome:** O005 — Prompt Review

## Requirement

The system produces a report directory whose contents include narrative feedback (summary.md) and a machine-parseable result (result.json), plus the original prompt (original.md), suggested revision (revision.md), and a diff (diff.md) for every completed review.

## Detail

Every successful review run produces a report. The report directory contains: (1) **result.json** — machine-parseable status for CI; (2) **summary.md** — narrative feedback; (3) **original.md** — prompt as submitted; (4) **revision.md** — suggested revision; (5) **diff.md** — diff between original and revision. The format of result.json is defined so that consumers can reliably parse it; the narrative in summary.md may be free-form or structured. When the review does not complete successfully (e.g. I/O error, reviewer failure), no report directory is produced and the documented failure code is used per R008.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Review completes successfully | Report directory is produced and written (R005); it contains result.json, summary.md, original.md, revision.md, diff.md. |
| Review fails before producing output | No report directory is required; the documented failure code is used (R008). |
| Prompt has no issues | Report still produced; narrative and summary can indicate "no issues" or "pass." |
| Prompt has severe issues | Summary allows CI/scripts to detect failure or severity; narrative explains issues. |
| Very long prompt | Report remains bounded and usable; summary is compact; narrative may be sectioned. |

### Examples

#### Report for prompt with issues

**Input:** User runs review on a prompt that lacks explicit success/failure signals. Report is written to default or user-chosen path.

**Expected output:** Report directory contains summary.md (narrative feedback, e.g. "The prompt does not specify how the AI should signal success or failure…") and result.json (e.g. status/errors/warnings); scripts can read result.json to gate.

**Verification:** Human can read narrative and understand the issue; a script can parse the summary and gate (e.g. use the documented prompt-errors or failure code when result is fail).

#### Report for compliant prompt

**Input:** User runs review on a prompt that meets the evaluation dimensions. Report is written to file.

**Expected output:** result.json indicates pass; summary.md contains narrative indicating compliance or minor suggestions.

**Verification:** CI can parse summary and pass; user can read narrative for confirmation.

## Acceptance criteria

- [ ] Every completed review produces a report directory that includes summary.md (narrative) and result.json (machine-parseable status).
- [ ] The machine-parseable summary is documented or stable enough for scripts and CI to parse and gate on the result.
- [ ] The report directory is written per R005 (path user-chosen or default).
- [ ] When the review does not complete successfully, a report is not required; the documented failure code is used per R008.

## Dependencies

- R005 — Report directory is written (report output path and persistence).
- R008 — Exit codes (so "completed" vs "did not complete" is well-defined).
