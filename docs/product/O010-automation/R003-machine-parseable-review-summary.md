# R003: Machine-parseable review summary

**Outcome:** O010 — Automation

## Requirement

Ralph provides a machine-parseable review summary (or equivalent) so CI and scripts can gate on review outcome without scraping free text.

## Detail

The review report includes a machine-parseable summary (see O005/R002). That summary is emitted in a defined format—e.g. a dedicated line or block in the report file, or a separate output—so that scripts and CI can parse it reliably and gate on pass/fail or severity without parsing narrative prose. The format is documented (e.g. in user docs). CI can then run the reviewer, parse the summary, and pass or fail the job based on the result. Ralph does not require scraping of free-form text to determine outcome.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Review completes successfully, no prompt errors | Summary indicates pass or equivalent; parseable. |
| Review completes, prompt has errors | Summary indicates fail or severity; parseable so CI can gate. |
| Review does not complete (e.g. I/O failure) | No report/summary required; the documented failure code (R002) indicates failure. |
| Report written to user-chosen path | Summary is part of the report file or otherwise obtainable at a documented location. |
| Very long narrative | Summary portion remains compact and parseable; format is stable. |

### Examples

#### CI gates on summary

**Input:** CI runs the review command with report output path; review completes; prompt has errors.

**Expected output:** Report file contains a machine-parseable summary (e.g. a documented line or field indicating fail, or a structured block). CI script parses that line/block and does not need to scrape the narrative.

**Verification:** Scripts can parse the summary (e.g. a documented line or field) to yield the outcome; CI can set job exit code based on it; no free-text scraping required.

#### Pass outcome

**Input:** CI runs review on a compliant prompt; report written to file.

**Expected output:** Report contains machine-parseable summary indicating pass or "no errors." CI can parse and pass the job.

**Verification:** Same parsing approach works for pass and fail; format is consistent.

## Acceptance criteria

- [ ] Every completed review produces output that includes a machine-parseable summary (or equivalent) whose format is documented.
- [ ] Scripts and CI can determine review outcome (e.g. pass/fail or severity) by parsing the summary without relying on free-form narrative text.
- [ ] The summary format is documented so that consumers can parse it reliably (e.g. line format, field names, or schema).
- [ ] When the review does not complete successfully, a summary is not required; the documented failure code (R002) indicates failure.

## Dependencies

- O005/R002 — Report content and format (narrative + machine-parseable summary) are defined there; this requirement ensures the summary is sufficient for CI/script gating and that the format is documented for automation.
