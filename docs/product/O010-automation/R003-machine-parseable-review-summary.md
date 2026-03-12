# R003: Machine-parseable review summary

**Outcome:** O010 — Automation

## Requirement

Ralph provides a machine-parseable review result (e.g. result.json in the report directory) so CI and scripts can gate on review outcome without scraping free text.

## Detail

The review report directory includes result.json with status (and optional errors/warnings counts) (see O005/R002). result.json is written inside the report directory at a documented path (e.g. <report-dir>/result.json). The format is documented so that scripts and CI can parse it reliably and gate on pass/fail or severity without parsing narrative prose. result.json remains compact and parseable (JSON). CI can then run the reviewer, read result.json, and pass or fail the job based on the result. Ralph does not require scraping of free-form text to determine outcome.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Review completes successfully, no prompt errors | Summary indicates pass or equivalent; parseable. |
| Review completes, prompt has errors | Summary indicates fail or severity; parseable so CI can gate. |
| Review does not complete (e.g. I/O failure) | No report/summary required; the documented failure code (R002) indicates failure. |
| Report written to user-chosen path | result.json is written inside the report directory at a documented path (e.g. <report-dir>/result.json). |
| Very long narrative | result.json remains compact and parseable (JSON). |

### Examples

#### CI gates on summary

**Input:** CI runs the review command with report output path; review completes; prompt has errors.

**Expected output:** Report directory contains result.json; CI script reads result.json and does not need to scrape narrative.

**Verification:** Scripts can parse the summary (e.g. a documented line or field) to yield the outcome; CI can set job exit code based on it; no free-text scraping required.

#### Pass outcome

**Input:** CI runs review on a compliant prompt; report written to file.

**Expected output:** result.json indicates status ok / no errors. CI can read result.json and pass the job.

**Verification:** Same parsing approach works for pass and fail; format is consistent.

## Acceptance criteria

- [ ] Every completed review produces a report directory that includes result.json (machine-parseable).
- [ ] Scripts and CI can determine review outcome (e.g. pass/fail or severity) by reading result.json (or the process exit code) without relying on free-form narrative text.
- [ ] The result.json schema and report directory layout are documented so that consumers can parse reliably.
- [ ] When the review does not complete successfully, result.json is not required; the documented failure code (R002) indicates failure.

## Dependencies

- O005/R002 — Report content and format (narrative + machine-parseable summary) are defined there; this requirement ensures the summary is sufficient for CI/script gating and that the format is documented for automation.
