# R5: Apply with Confirmation and Revision Phase

**Outcome:** O5 — Prompt review

## Requirement

The system supports applying the suggested revision to the user's prompt. Apply is requested by the user (e.g. `--apply`). When apply is requested, the system may run a second AI invocation (revision phase) to produce the revised prompt. The revision-phase prompt instructs the AI where to read the report from (review output path) and where to write the revised prompt (prompt output path); both paths are interpolated. When apply is requested with confirmation (e.g. without `-y`), the system prompts the user to confirm before writing. With `-y` (or equivalent), the system applies without interactive confirmation. When the prompt source is stdin, apply requires the user to supply `--prompt-output` so the system has a path to write to; using `--apply` with stdin and without `--prompt-output` is invalid and results in exit 2.

## Specification

(To be specified in Step 5.)

## Acceptance criteria

- [ ] User can request apply (e.g. `--apply`); the system then writes the suggested revision to the appropriate path (source file for alias/file, or `--prompt-output` when set or required).
- [ ] When confirmation is required, the system prompts the user before writing the revision; on confirmation, the revision is written.
- [ ] User can bypass confirmation (e.g. `--apply -y`); the revision is written without prompting.
- [ ] The revision phase (if used) receives a prompt that includes the review output path and the prompt output path so the AI knows where to read the report and where to write the revised prompt.
- [ ] When prompt source is stdin and user specifies `--apply` without `--prompt-output`, the system reports an error and exits 2; it does not attempt to apply.
- [ ] Apply is supported for all three input modes (alias, file, stdin) subject to the stdin + `--prompt-output` requirement.

## Dependencies

- R3 (review output path) and R4 (prompt output path) supply the paths interpolated into the revision-phase prompt. R8 (failure handling) covers invalid apply (e.g. stdin + apply without `--prompt-output`) with exit 2.
