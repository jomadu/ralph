# R4: Prompt Output Path

**Outcome:** O5 — Prompt review

## Requirement

The system supports directing the suggested revised prompt to a user-specified path via `--prompt-output`. When the user does not request apply, this path is where the revised prompt (suggested revision) is written without modifying the source. When the user requests apply and the source is an alias or a file, the path for applying can be the source file (so `--prompt-output` is optional). When the user requests apply and the source is stdin, `--prompt-output` is required so the system knows where to write the revised prompt; absence of `--prompt-output` with `--apply` and stdin is invalid (handled by R5/R8).

## Specification

**Semantics of `--prompt-output <path>`:**

- **When apply is not requested:** `--prompt-output <path>` directs where the suggested revised prompt is written. The AI (in the review phase or a dedicated revision phase) is instructed to write the revised prompt to this path. The source (alias path, `-f` path, or stdin) is never modified. If `--prompt-output` is omitted and apply is not requested, the revised prompt may still be included in the report (per R6) but need not be written to a separate file by path; implementation may still write to a temp path for internal use or omit a separate file — the requirement is that when the user sets `--prompt-output`, the revision is written there.
- **When apply is requested and source is alias or file:** The path for applying is the source file (the alias's configured path or the `-f` path). `--prompt-output` is optional; if provided, it can override the apply destination (revision written to `--prompt-output` instead of source), or implementation may define that with apply, the revision always goes to the source unless a different convention is specified. To avoid ambiguity: when apply is requested and source is alias or file, the default apply destination is the source file; `--prompt-output` may be used to direct the revised prompt to a different path instead of overwriting the source (so user can do apply-to-different-file). When apply is requested and source is stdin, `--prompt-output` is required and is the only possible destination; absence is invalid (exit 2, R5/R8).

**Path interpolation:** The "prompt output path" (resolved as above) is the path Ralph interpolates into the revision-phase prompt so the AI knows where to write the revised prompt (R5). Resolution: (1) If apply and stdin: must have `--prompt-output` → use it. (2) If apply and alias/file: use source path unless `--prompt-output` is set, then use `--prompt-output`. (3) If no apply and user set `--prompt-output`: use it for where to write the suggested revision. (4) If no apply and user did not set `--prompt-output`: no path need be interpolated for revision file output (revision may only appear in report).

**Validation:** If apply is requested with stdin and `--prompt-output` is not provided, Ralph must not start the revision phase; report error and exit 2 (R8). Validation order: after input mode and apply flag are known, before running review or revision phase.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Stdin + `--apply` without `--prompt-output` | Invalid; exit 2; do not run revision phase. |
| Alias + `--apply` without `--prompt-output` | Apply destination is the alias's prompt path (from config). |
| `-f ./p.md` + `--apply` without `--prompt-output` | Apply destination is `./p.md`. |
| Stdin + `--apply` + `--prompt-output out.md` | Apply destination is `out.md`; revision written there. |
| Any input + no apply + `--prompt-output out.md` | Suggested revision is written to `out.md`; source unchanged. |
| `--prompt-output` path is unwritable or invalid | Fail when revision phase would write (or before); exit 2 (R8). |

### Examples

#### Apply to source (alias)

**Input:** `ralph review build --apply -y` (alias `build` → `./prompts/build.md`)

**Expected output:** Revision is written to `./prompts/build.md`. No `--prompt-output` required.

#### Apply with stdin (required flag)

**Input:** `cat prompt.md | ralph review --apply --prompt-output revised.md -y`

**Expected output:** Revision is written to `revised.md`. Without `--prompt-output`, same command would exit 2.

## Acceptance criteria

- [ ] User can set `--prompt-output <path>` so the suggested revised prompt is written to that path; the source prompt (alias or file) is not modified when apply is not requested.
- [ ] When input is from stdin and the user does not request apply, `--prompt-output` can be used to write the suggested revision to a file.
- [ ] When input is from alias or file and the user requests apply, the revision can be written to the source (alias path or `-f` path) so `--prompt-output` is optional for apply.
- [ ] When input is from stdin and the user requests apply, `--prompt-output` is required; the system treats apply-with-stdin and no `--prompt-output` as invalid (exit 2).
- [ ] The prompt output path, when determined (from `--prompt-output` or source), is interpolated into the revision-phase prompt so the AI knows where to write the revised prompt (per R5).

## Dependencies

- R5 (apply and revision phase) uses this path for apply behavior and revision-phase interpolation. R8 defines exit 2 for invalid apply (e.g. stdin + apply without `--prompt-output`).
