# R4: Prompt Output Path

**Outcome:** O5 — Prompt review

## Requirement

The system supports directing the suggested revised prompt to a user-specified path via `--prompt-output`. When the user does not request apply, this path is where the revised prompt (suggested revision) is written without modifying the source. When the user requests apply and the source is an alias or a file, the path for applying can be the source file (so `--prompt-output` is optional). When the user requests apply and the source is stdin, `--prompt-output` is required so the system knows where to write the revised prompt; absence of `--prompt-output` with `--apply` and stdin is invalid (handled by R5/R8).

## Specification

(To be specified in Step 5.)

## Acceptance criteria

- [ ] User can set `--prompt-output <path>` so the suggested revised prompt is written to that path; the source prompt (alias or file) is not modified when apply is not requested.
- [ ] When input is from stdin and the user does not request apply, `--prompt-output` can be used to write the suggested revision to a file.
- [ ] When input is from alias or file and the user requests apply, the revision can be written to the source (alias path or `-f` path) so `--prompt-output` is optional for apply.
- [ ] When input is from stdin and the user requests apply, `--prompt-output` is required; the system treats apply-with-stdin and no `--prompt-output` as invalid (exit 2).
- [ ] The prompt output path, when determined (from `--prompt-output` or source), is interpolated into the revision-phase prompt so the AI knows where to write the revised prompt (per R5).

## Dependencies

- R5 (apply and revision phase) uses this path for apply behavior and revision-phase interpolation. R8 defines exit 2 for invalid apply (e.g. stdin + apply without `--prompt-output`).
