# R005: Report to file

**Outcome:** O005 — Prompt Review

## Requirement

The report is always written to a file; the user can choose the path or accept a default location.

## Detail

Every completed review produces a report (R002). That report is always saved to a file. The user may specify the report output path; if they do not, the system uses a default (e.g. a temporary file, or a file in the current directory with a documented name). The default location is documented so the user knows where to find the report. If the user-specified path cannot be written (e.g. permission denied, directory does not exist), the system reports an error and does not complete the review successfully; the documented failure code is used per R008.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User specifies a valid path | Report is written to that path. |
| User does not specify path | Report is written to the documented default (e.g. temp file or cwd). |
| Specified path is not writable (permission, read-only FS) | Error; treat as failure and exit with the documented failure code (R008). Report is not considered delivered unless written to a file. |
| Specified path is a directory | Error or defined behavior (e.g. write report as a file inside that directory with a default name); no silent misuse. |
| Default path would overwrite existing file | Behavior is defined (overwrite, or unique name, or error); user can avoid by specifying path. |

### Examples

#### User-chosen path

**Input:** User runs the review command with the prompt supplied by file path and specifies the report output path (e.g. to a user-chosen path). Review completes.

**Expected output:** Report is written to the specified path; suggested revision is also produced (R003). Exit code per R008.

**Verification:** File exists and contains the report (narrative + machine-parseable summary per R002).

#### Default location

**Input:** User runs the review command with the prompt supplied by file path and does not specify the report output path. Review completes.

**Expected output:** Report is written to the default location (documented in CLI or docs). User can find and open the file.

**Verification:** A report file exists at the documented default; content matches R002.

#### Unwritable path

**Input:** User runs the review command and specifies a report output path that is not writable (e.g. permission denied).

**Expected output:** System reports an error (e.g. cannot write report); does not silently succeed. The documented failure code is used per R008.

**Verification:** No report file is created at that path; user sees clear error; the documented failure code is used.

## Acceptance criteria

- [ ] Every completed review writes the report to a file; there is no "report only to stdout" mode that bypasses file output for the report.
- [ ] User can specify the report output path; when specified and writable, the report is written there.
- [ ] When the user does not specify a path, the report is written to a documented default location.
- [ ] When the chosen or default path cannot be written, the system reports an error and exits with the documented failure code per R008; behavior is defined for invalid path (e.g. directory) so there is no silent misuse.

## Dependencies

- R002 — Report content and format; this requirement covers persistence and path.
- R008 — Exit code on write failure.
