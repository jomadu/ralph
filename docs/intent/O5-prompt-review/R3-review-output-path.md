# R3: Review Output Path

**Outcome:** O5 — Prompt review

## Requirement

The system always writes the review report to a file. The report path is either the value of `--review-output` when provided, or a path in system temporary storage (e.g. the platform temp directory) when `--review-output` is not set. The user can determine where the report was written (e.g. by being told the path when the default temp location is used). Invalid or unwritable report paths are handled as failures (see R8).

## Specification

(To be specified in Step 5.)

## Acceptance criteria

- [ ] When the user supplies `--review-output <path>`, the report is written to that path (when the review phase succeeds and the AI writes the report there).
- [ ] When the user does not supply `--review-output`, the report is written to a path in system temporary storage.
- [ ] When the default (temp) location is used, the user can discover the report path (e.g. it is printed or otherwise communicated).
- [ ] The report is always persisted to a file; there is no mode where the report exists only in memory or only on stdout without a file.
- [ ] Invalid or unwritable report path is handled as a review failure (exit 2 per R8).

## Dependencies

- R2 (review prompt composition) interpolates this path into the prompt so the AI knows where to write. R9 (report file verification) checks that the report exists at the expected path after the run.
