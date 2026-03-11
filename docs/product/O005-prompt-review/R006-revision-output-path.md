# R006: Revision output path

**Outcome:** O005 — Prompt Review

## Requirement

The revised prompt can be written to a user-chosen path; when the prompt was supplied via standard input and the user requests apply, the user must specify the revision output path — if they do not, the system reports an error and does not apply.

## Detail

When the user requests that the revision be written to a file (R004), the destination is a user-chosen path. When the prompt was **supplied via standard input**, there is no source file to overwrite, so the system cannot infer a path. The user must supply the revision output path. If they request apply but do not supply a path when the prompt was supplied via standard input, the system reports an error and does not write the revision anywhere; the documented failure code is used per R008. When the prompt was supplied from a **file or alias**, the user may specify a path (which may be the same file, implying overwrite with confirmation per R004) or a different file; behavior for "apply to same file" vs "apply to other file" is defined and consistent with R004.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Prompt supplied via standard input + apply, no path given | Error; do not write revision; the documented failure code is used (R008). |
| Prompt supplied via standard input + apply, path given | Write revision to the given path (subject to confirmation if overwrite, per R004). |
| File/alias + apply, path given | Write revision to the given path; confirm if overwriting existing file per R004. |
| File/alias + apply, path omitted | Behavior defined: e.g. default to source file (with confirmation) or require path; product chooses and documents. |
| Given path is directory or unwritable | Error; do not write; clear message; the documented failure code is used. |

### Examples

#### Prompt supplied via standard input + apply without path (error)

**Input:** User runs the review command with the prompt supplied via standard input and requests apply. No revision output path is supplied.

**Expected output:** System reports an error that the revision output path is required when applying when the prompt was supplied via standard input. No file is written. The documented failure code is used per R008.

**Verification:** No revision file is created; the user sees a clear message that path is required for apply when prompt was from standard input; the documented failure code is used.

#### Prompt supplied via standard input + apply with path

**Input:** User runs the review command with the prompt supplied via standard input, requests apply, and specifies the revision output path. Review completes.

**Expected output:** Revision is written to the specified path. Report is written per R005. Exit code per R008.

**Verification:** The revision file contains the suggested revision; the original source was not a file so no overwrite of a "source" file.

#### File + apply to chosen path

**Input:** User runs the review command with the prompt supplied by file path, requests apply, and specifies the revision output path. Target path may or may not exist. User confirms if prompted (R004). Review completes.

**Expected output:** Revision is written to the user-chosen path; the source file is unchanged unless the user had chosen it as the revision output (then confirmation per R004).

**Verification:** Revision content is in the specified file; user had explicit control over destination.

## Acceptance criteria

- [ ] When the user requests apply and the prompt was supplied via standard input, the user must supply the revision output path; if they do not, the system reports an error and does not write the revision; the documented failure code is used.
- [ ] When the user requests apply and supplies a valid path (whether prompt was from standard input or file/alias), the revision is written to that path subject to confirmation rules in R004.
- [ ] When the revision output path is invalid (e.g. directory, unwritable), the system reports an error and does not write; the documented failure code is used.
- [ ] Behavior when prompt was from file/alias and user requests apply without specifying path is defined and documented (e.g. default to source with confirmation, or require path).

## Dependencies

- R003 — Suggested revision content.
- R004 — Apply and confirmation; this requirement defines path rules for apply.
- R008 — Exit codes when path is missing or write fails.
