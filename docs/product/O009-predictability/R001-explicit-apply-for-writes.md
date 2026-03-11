# R001: Explicit apply for writes

**Outcome:** O009 — Predictability

## Requirement

Ralph writes to user prompt content (prompt files or user-specified output paths for revisions) only when the user has explicitly requested that a revision be applied and has confirmed (in interactive use) or used a documented non-interactive apply option.

## Detail

User prompt content includes: (1) prompt files on disk that Ralph reads for the loop or for review, and (2) user-specified output paths for revised prompt content (e.g. when applying a review revision to a chosen path). Ralph must not write to any of these unless the user has explicitly requested that a revision be applied. In interactive use, that request must be accompanied by confirmation (e.g. prompt or explicit confirm step). In non-interactive use, a documented non-interactive apply option (e.g. documented flag) constitutes the explicit request. In-memory or composed prompts (e.g. preamble + prompt buffer) are never written back to the source prompt file unless the user requests apply and confirms or uses the non-interactive option.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User runs the loop (no review, no apply) | Prompt file on disk is never written to; any composed prompt is not persisted to the source file. |
| User runs the reviewer without requesting apply | No write to the source prompt file or to any path; report and suggested revision may be produced and shown or written to user-chosen report/revision output paths per O005. |
| User runs review and then requests apply but does not confirm (interactive) | Ralph does not write the revision to the prompt file or user path until the user confirms. |
| User runs review and uses a documented non-interactive apply option (e.g. documented flag) | Ralph writes the revision to the appropriate path (source file or user-specified revision output path per O005/R006) without prompting. |
| Prompt supplied via stdin; user requests apply | Ralph writes only to the user-specified revision output path (no default overwrite of a file); per O005/R006. |
| User specifies a revision output path and then requests apply | Ralph writes to that path when the user requests that the revision be applied and has confirmed (interactive) or used the non-interactive option. |

### Examples

#### Loop run leaves prompt file unchanged

**Input:** User invokes the run command (or equivalent) with a prompt file configured. The loop runs one or more iterations.

**Expected output:** The prompt file on disk is unchanged. Any in-memory composed prompt (e.g. preamble + buffer) is not written back to that file.

**Verification:** Compare file mtime or content before and after the run; no change to the prompt file.

#### Review without apply leaves source file unchanged

**Input:** User runs the review command on a prompt file (e.g. ./my-prompt.md) and views the report and suggested revision. User does not request that the revision be applied.

**Expected output:** `my-prompt.md` is unchanged. Report and suggested revision may be written to paths the user specified for those outputs (e.g. report file, revision to stdout or a separate file), but the source prompt file is not modified.

**Verification:** Content of the source prompt file is identical before and after the review command.

#### Apply with confirmation writes revision

**Input:** User runs review, then explicitly requests that the revision be applied and confirms (e.g. prompt or explicit confirm step, such as "Apply? [y/N]" → "y").

**Expected output:** Ralph writes the revised prompt to the chosen path (e.g. the source file or a user-specified output path). The write occurs only after the explicit request and confirmation.

**Verification:** After apply, the target file contains the revised content; the apply step was explicitly requested and confirmed.

#### Non-interactive apply with documented option

**Input:** User runs the review command with the documented apply option and a non-interactive option (e.g. documented flag) in a non-interactive session. No interactive confirmation is possible.

**Expected output:** Ralph writes the revision to the appropriate path (source file or user-specified revision output path) without prompting.

**Verification:** Target file contains the revision; no interactive prompt was shown; the option is documented as modifying user content.

## Acceptance criteria

- [ ] When the user runs the loop or the reviewer without requesting that a revision be applied, Ralph does not write to the user's prompt file or to any path used as the source of the prompt.
- [ ] In-memory or composed prompt (e.g. preamble + buffer) is never written back to the source prompt file unless the user has requested apply and confirmed or used the non-interactive apply option.
- [ ] When the user requests that a revision be applied in interactive use, Ralph writes to the chosen path only after the user has confirmed (e.g. explicit confirm step or prompt).
- [ ] When the user uses a documented non-interactive apply option (e.g. documented flag), Ralph writes the revision to the appropriate path without requiring interactive confirmation; the option is documented as modifying user content.
- [ ] When the prompt was supplied via stdin and the user requests apply, Ralph writes only to a user-specified revision output path (no silent overwrite of a file).

## Dependencies

- O005 (Prompt Review) defines the review flow, report, and suggested revision; O005/R006 defines revision output path behavior when the prompt comes from stdin or when the user specifies an output path. R001 constrains when Ralph may write to those paths.
