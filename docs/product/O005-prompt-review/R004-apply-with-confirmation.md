# R004: Apply with confirmation

**Outcome:** O005 — Prompt Review

## Requirement

The user can request that the revision be written to a file, with confirmation when appropriate or a non-interactive option where supported.

## Detail

Applying the suggested revision to a file is optional and user-initiated. When the user requests that the revision be written (via the documented apply option), the system either (1) prompts for confirmation before writing when the operation would overwrite an existing file or is otherwise destructive, or (2) respects a non-interactive option (e.g. a documented flag to apply without confirmation) so that CI or scripts can apply without a prompt. When confirmation would be required and the session is non-interactive, the system does not write without the explicit non-interactive option and instead reports that confirmation is required or fails with a clear message. This prevents accidental overwrites.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User requests apply, target path exists (overwrite) | System prompts for confirmation unless non-interactive option is set. |
| User requests apply, target path does not exist | System may write without confirmation (new file), or still prompt per product policy; behavior is defined. |
| Non-interactive session, user did not set non-interactive option | System does not write; reports that confirmation is required or exits with clear error. |
| Non-interactive option set | System writes without prompting when path is valid (subject to R006 when prompt was supplied via standard input). |
| User does not request apply | No file is written; revision is only in output. |
| Prompt supplied via standard input + apply without revision output path | Per R006: error, do not apply. |

### Examples

#### Apply with confirmation (overwrite)

**Input:** User runs the review command with the prompt supplied by file path and requests that the revision be written (e.g. to same path or default). Target file exists. Session is interactive.

**Expected output:** System prompts to confirm overwrite (or similar). If user confirms, revision is written; if user declines, no write occurs. User sees clear outcome.

**Verification:** No overwrite without user confirmation in interactive sessions; revision file matches suggested revision after confirm.

#### Non-interactive apply (CI)

**Input:** User runs the review command with prompt by file path, requests apply, specifies the revision output path, and uses the non-interactive option, in a non-interactive session. Review completes.

**Expected output:** Revision is written to the specified path without prompting. The documented success code or documented prompt-errors code is used per R008.

**Verification:** File is written; no hang waiting for input; CI can gate on exit code.

#### Apply not requested

**Input:** User runs the review command with the prompt supplied by file path and does not request apply.

**Expected output:** Report and suggested revision are produced; no file is written for the revision. User may copy revision from output or run again with the apply option if desired.

**Verification:** Only report file is written (per R005); revision is not applied to any file.

## Acceptance criteria

- [ ] User can request that the suggested revision be written to a file (via the documented apply option).
- [ ] When the write would overwrite an existing file (or is otherwise destructive), the system prompts for confirmation in interactive sessions unless a non-interactive option is set.
- [ ] A non-interactive option (e.g. documented flag to apply without confirmation) allows apply without confirmation so that CI and scripts can use it.
- [ ] When confirmation would be required and the session is non-interactive and the non-interactive option is not set, the system does not write and reports clearly (e.g. "confirmation required" or exit with error).
- [ ] Without an explicit user request to apply, the system does not write the revision to any file.

## Dependencies

- R003 — Suggested revision is what gets applied.
- R006 — Output path and stdin+apply path requirement; apply behavior is coordinated with path rules.
