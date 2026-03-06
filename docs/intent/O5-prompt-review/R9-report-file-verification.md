# R9: Report File Verification

**Outcome:** O5 — Prompt review

## Requirement

After the review-phase AI invocation completes, the system verifies that the report file exists at the expected path (the path that was interpolated into the prompt per R3). If the report is not present at that path—for example because the AI did not write it, wrote to a different location, or the process failed before writing—the system treats this as a review failure and exits 2. This ensures the user and downstream steps (e.g. apply, CI) can rely on the report file being present when exit code is 0 or 1.

## Specification

**When verification runs:** Immediately after the review-phase AI process exits. Before any apply flow, before deriving exit 0 vs 1 from report content (R6), and before printing report to stdout (R7), Ralph checks that a file exists at the review output path (the same path that was interpolated into the review prompt per R3).

**Check:** A file must exist at that path. The check is existence only: a regular file (or platform equivalent) is present at the path. The specification does not require validating the report content structure (e.g. that it contains the machine-parseable summary); content validation is optional and may be specified elsewhere. If the file exists, Ralph proceeds: read the report for exit code derivation (R6), optionally print to stdout (R7), and if apply was requested, proceed to revision phase (R5). If the file does not exist, Ralph must not exit 0 or 1; it must report a review failure and exit 2 (R8). No apply or revision phase is run when verification fails.

**Reasons the file might be missing:** The AI did not write it; the AI wrote to a different path; the AI process crashed before writing; the path was wrong or the filesystem rejected the write. Regardless of reason, missing file → exit 2.

**Order of operations:** (1) Review phase AI runs. (2) AI process exits. (3) Ralph checks: file exists at review output path? (4) If no → exit 2, no apply. (5) If yes → read report, derive exit 0/1 (R6), optionally print (R7), if apply then revision phase (R5).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| File exists but is empty (0 bytes) | Implementation may treat as pass (existence only) or fail (empty = invalid). Recommend: pass (existence only); R6 parsing may then yield exit 1 if content is unparseable. |
| File exists but is not a regular file (e.g. directory) | Treat as failure: exit 2 (expected a file, not a directory). |
| Path exists but is a symlink to a missing file | Treat as failure: exit 2. |
| AI wrote to a sibling path (typo) | File at expected path is missing → exit 2. |

### Examples

#### Report present

**Input:** Review phase completes; AI wrote report to `/tmp/ralph-report-xyz.md`.

**Expected output:** Ralph verifies file exists, proceeds to parse for exit 0/1, optionally apply; does not exit 2 from verification.

#### Report missing

**Input:** Review phase completes; no file at the expected path (AI crashed or wrote elsewhere).

**Expected output:** Ralph reports failure (e.g. "report file not found at <path>"); exit 2; revision phase is not run.

## Acceptance criteria

- [ ] After the review-phase AI process exits, Ralph checks that a file exists at the review output path (from `--review-output` or the chosen temp path).
- [ ] If the file does not exist at that path, Ralph does not exit 0 or 1; it reports a review failure and exits 2.
- [ ] If the file exists at that path, Ralph may proceed (e.g. to derive exit 0 vs 1 from report content, or to proceed to apply if requested); verification does not require validating report content structure; content validation is optional (e.g. R6 parsing).
- [ ] Verification happens before any apply or revision phase so that apply is not attempted when the report was not produced.
- [ ] Exit 2 from report verification is consistent with R6 and R8 (review failure).

## Dependencies

- R3 (review output path) defines where the report is expected. R8 (failure handling) defines exit 2 for review failures; missing report is one such failure. R6 (exit codes) specifies that exit 2 means review failed to run or complete.
