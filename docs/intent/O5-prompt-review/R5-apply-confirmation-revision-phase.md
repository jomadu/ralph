# R5: Apply with Confirmation and Revision Phase

**Outcome:** O5 — Prompt review

## Requirement

The system supports applying the suggested revision to the user's prompt. Apply is requested by the user (e.g. `--apply`). When apply is requested, the system may run a second AI invocation (revision phase) to produce the revised prompt. The revision-phase prompt instructs the AI where to read the report from (review output path) and where to write the revised prompt (prompt output path); both paths are interpolated. When apply is requested with confirmation (e.g. without `-y`), the system prompts the user to confirm before writing. With `-y` (or equivalent), the system applies without interactive confirmation. When the prompt source is stdin, apply requires the user to supply `--prompt-output` so the system has a path to write to; using `--apply` with stdin and without `--prompt-output` is invalid and results in exit 2.

## Specification

**Apply flag:** User requests apply via `--apply` (or equivalent flag). When set, the system will write the suggested revision to the prompt output path (R4): either the source file (alias or `-f`) or `--prompt-output` when required (stdin) or when user explicitly directs output elsewhere.

**Confirmation:** When apply is requested and the user has not passed a non-interactive confirmation flag (e.g. `-y` / `--yes`), Ralph must prompt the user for confirmation before writing (e.g. "Apply revision to <path>? [y/N]"). On confirmation (e.g. y/yes), proceed to write. On decline or EOF, do not write; exit 0 or 1 based on review result (revision not applied). When `-y` (or equivalent) is passed, skip the prompt and write the revision without interaction.

**Revision phase:** To produce the revised prompt content for apply, Ralph may run a second AI invocation (revision phase). The revision-phase prompt must include: (1) the review output path (from R3) — "Read the report from: <path>"; (2) the prompt output path (from R4) — "Write the revised prompt to: <path>". Both paths are interpolated so the AI reads the report from the file produced in the review phase and writes the new prompt to the apply destination. The revision phase is only run when apply is requested and after the review phase has completed and the report file exists (R9). If the report file is missing, R9 yields exit 2 and the revision phase is not run.

**Stdin + apply validation:** Before running the review phase (or at latest before the revision phase), if input source is stdin and `--apply` is set, Ralph must check that `--prompt-output <path>` is provided. If not, Ralph must not run the revision phase; report an error (e.g. "stdin input with --apply requires --prompt-output") and exit 2 (R8).

**Apply destination (recap):** Alias → alias's configured path; `-f <path>` → that path; stdin → must have `--prompt-output`. Optional `--prompt-output` with alias/file can override the apply destination to a different file (revision written there, source unchanged) — per R4.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| `--apply` without `-y` and TTY not available | Treat as declined or require `-y`; do not hang. If no TTY and no `-y`, exit without applying (e.g. exit 0 or 1, revision not written); or exit 2 with message that `-y` is required for non-interactive apply. |
| Stdin + `--apply` + no `--prompt-output` | Exit 2 before revision phase; do not run revision phase. |
| Revision phase AI does not write to prompt output path | R8/R9-style handling: if implementation verifies revision file (like R9 for report), missing revision file → exit 2; else document that apply "best effort" and exit code reflects review result. Recommend: verify revision file exists after revision phase when apply was requested; if missing, exit 2. |
| User declines confirmation | Do not write; exit 0 or 1 from review result only. |

### Examples

#### Apply with confirmation

**Input:** `ralph review build --apply` (no `-y`); review completes; user types `y` at prompt.

**Expected output:** Ralph prompts "Apply revision to ./prompts/build.md? [y/N]"; on `y`, runs revision phase (if used), then writes revised content to `./prompts/build.md`.

#### Apply without confirmation

**Input:** `ralph review -f prompt.md --apply -y`

**Expected output:** After review, revision is written to `prompt.md` without prompting.

#### Stdin + apply requires --prompt-output

**Input:** `cat p.md | ralph review --apply` (no `--prompt-output`)

**Expected output:** Ralph reports error and exits 2; revision phase is not run.

## Acceptance criteria

- [ ] User can request apply (e.g. `--apply`); the system then writes the suggested revision to the appropriate path (source file for alias/file, or `--prompt-output` when set or required).
- [ ] When confirmation is required, the system prompts the user before writing the revision; on confirmation, the revision is written.
- [ ] User can bypass confirmation (e.g. `--apply -y`); the revision is written without prompting.
- [ ] The revision phase (if used) receives a prompt that includes the review output path and the prompt output path so the AI knows where to read the report and where to write the revised prompt.
- [ ] When prompt source is stdin and user specifies `--apply` without `--prompt-output`, the system reports an error and exits 2; it does not attempt to apply.
- [ ] Apply is supported for all three input modes (alias, file, stdin) subject to the stdin + `--prompt-output` requirement.

## Dependencies

- R3 (review output path) and R4 (prompt output path) supply the paths interpolated into the revision-phase prompt. R8 (failure handling) covers invalid apply (e.g. stdin + apply without `--prompt-output`) with exit 2.
