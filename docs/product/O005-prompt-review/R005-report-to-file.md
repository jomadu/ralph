# R005: Report to file

**Outcome:** O005 — Prompt Review

## Requirement

The report is always written to a directory of files; the user can choose the directory path or accept a default location.

## Detail

Every completed review produces a report (R002). That report is always saved as a directory containing result.json, summary.md, original.md, revision.md, and diff.md. The user may specify the report output directory; if they do not, the system uses a default (e.g. `./ralph-review/` in the current working directory). The default location is documented so the user knows where to find the report. If the directory cannot be created or written (e.g. permission denied, path exists and is a file), the system reports an error and does not complete the review successfully; the documented failure code is used per R008.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User specifies a valid path | User specifies a valid directory path (or path to be created as directory). Ralph creates the directory if needed and passes the path to the AI in the prompt; the AI creates the five files inside. |
| User does not specify path | Report directory defaults to the documented path (e.g. ./ralph-review/); that path is interpolated into the prompt for the AI. |
| Specified path is not writable (permission, read-only FS) | Specified directory path cannot be created or is not writable → error; exit with the documented failure code (R008). Report is not considered delivered unless written. |
| Specified path is a directory | Desired: report directory. The AI is instructed to create the five files in it (path is in the prompt); overwrite if present per policy. |
| Default path would overwrite existing file | Default is a directory; if it exists as a directory, the AI writes the five files inside; if it exists as a file, error (path must be a directory). |

### Examples

#### User-chosen path

**Input:** User runs the review command with the prompt supplied by file path and specifies the report output path (e.g. to a user-chosen path). Review completes.

**Expected output:** Report directory is created/written at the specified path; five files are inside. Suggested revision is also produced (R003). Exit code per R008.

**Verification:** Directory exists and contains result.json, summary.md, original.md, revision.md, diff.md per R002.

#### Default location

**Input:** User runs the review command with the prompt supplied by file path and does not specify the report output path. Review completes.

**Expected output:** Report directory is written to the default directory (e.g. ./ralph-review/); user can find the files there.

**Verification:** Report directory exists at the documented default; content matches R002.

#### Unwritable path

**Input:** User runs the review command and specifies a report output path that is not writable (e.g. permission denied).

**Expected output:** System reports an error (e.g. cannot write report); does not silently succeed. The documented failure code is used per R008.

**Verification:** No report directory is created; five files are not written; user sees clear error; the documented failure code is used.

## Acceptance criteria

- [ ] Every completed review writes the report to a directory (five files); there is no "report only to stdout" mode that bypasses file output for the report.
- [ ] User can specify the report output directory; when specified and writable, the report is written there.
- [ ] When the user does not specify a path, the report directory is written to the documented default (e.g. ./ralph-review/).
- [ ] When the directory cannot be created or written, the system reports an error and exits with the documented failure code per R008; behavior is defined for invalid path (e.g. path exists as file) so there is no silent misuse.

## Dependencies

- R002 — Report content and format; this requirement covers persistence and path.
- R008 — Exit code on write failure.
