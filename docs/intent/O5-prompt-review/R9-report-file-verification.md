# R9: Report File Verification

**Outcome:** O5 — Prompt review

## Requirement

After the review-phase AI invocation completes, the system verifies that the report file exists at the expected path (the path that was interpolated into the prompt per R3). If the report is not present at that path—for example because the AI did not write it, wrote to a different location, or the process failed before writing—the system treats this as a review failure and exits 2. This ensures the user and downstream steps (e.g. apply, CI) can rely on the report file being present when exit code is 0 or 1.

## Specification

(To be specified in Step 5.)

## Acceptance criteria

- [ ] After the review-phase AI process exits, Ralph checks that a file exists at the review output path (from `--review-output` or the chosen temp path).
- [ ] If the file does not exist at that path, Ralph does not exit 0 or 1; it reports a review failure and exits 2.
- [ ] If the file exists at that path, Ralph may proceed (e.g. to derive exit 0 vs 1 from report content, or to proceed to apply if requested); verification does not require validating report content structure (that can be specified in Step 5).
- [ ] Verification happens before any apply or revision phase so that apply is not attempted when the report was not produced.
- [ ] Exit 2 from report verification is consistent with R6 and R8 (review failure).

## Dependencies

- R3 (review output path) defines where the report is expected. R8 (failure handling) defines exit 2 for review failures; missing report is one such failure. R6 (exit codes) specifies that exit 2 means review failed to run or complete.
