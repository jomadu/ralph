# R8: Review Failure Handling

**Outcome:** O5 — Prompt review

## Requirement

The system treats certain conditions as review failures and exits with code 2. These include: invalid or missing configuration, missing or unreadable prompt source (for alias or file input), failure to spawn or run the AI process, invalid apply request (e.g. `--apply` with stdin and no `--prompt-output`), and unwritable or invalid report path. The user and scripts can rely on exit 2 to mean "review did not complete successfully" or "apply was invalid," distinct from exit 0 (success) and exit 1 (review completed but prompt has errors).

## Specification

(To be specified in Step 5.)

## Acceptance criteria

- [ ] When configuration is invalid or required config is missing, Ralph reports the failure and exits 2.
- [ ] When prompt source is invalid (alias not found, file missing or unreadable), Ralph fails fast and exits 2.
- [ ] When the AI command cannot be spawned or fails in a way that prevents the review from completing, Ralph exits 2.
- [ ] When the user requests `--apply` with stdin but does not provide `--prompt-output`, Ralph reports an error and exits 2.
- [ ] When the review output path (from `--review-output` or temp) is invalid or unwritable, Ralph exits 2 (or fails before running the AI).
- [ ] Exit 2 is used consistently for these failure cases so scripts can distinguish them from exit 0 (ok) and exit 1 (errors in prompt).
- [ ] Error messages or logs give the user enough information to correct the condition (e.g. missing `--prompt-output` when using stdin + apply).

## Dependencies

- R6 (exit codes) defines exit 2 semantics; this requirement enumerates the conditions that must yield exit 2. R3 (review output path) and R9 (report file verification) define path-related failures that also result in exit 2. R5 defines the stdin + apply + `--prompt-output` validation.
