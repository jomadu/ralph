# R3: Review Output Path

**Outcome:** O5 — Prompt review

## Requirement

The system always writes the review report to a file. The report path is either the value of `--review-output` when provided, or a path in system temporary storage (e.g. the platform temp directory) when `--review-output` is not set. The user can determine where the report was written (e.g. by being told the path when the default temp location is used). Invalid or unwritable report paths are handled as failures (see R8).

## Specification

**Report path determination:**

- If `--review-output <path>` is provided: the report path is that path (after resolving relative to current working directory). Ralph instructs the AI to write the report to this path (per R2) and uses it for verification (R9).
- If `--review-output` is not provided: the report path is a new file in the system temporary directory. The path must be unique per run (e.g. a new temp file or a uniquely named file in the platform temp dir, such as `os.TempDir()` + unique filename). The implementation must use the platform's standard temp directory (e.g. `$TMPDIR` on Unix, or the OS default).

**User discovery when default (temp) is used:** After choosing the temp path, Ralph must communicate it to the user so they can find the report. For example: print the path to stderr or to stdout (before or after the report content, per R7), or include it in a final summary line. The exact channel is implementation-defined but must be documented; the user must be able to discover the path without guessing.

**Persistence:** The report is always written to a file at the chosen path. There is no mode where the report exists only in memory or only on stdout; the file is the canonical output (the AI is instructed to write there, and R9 verifies the file exists).

**Validation:** Before running the review-phase AI, Ralph must ensure the report path is usable: parent directory exists and is writable when using `--review-output`; temp directory is available when not. If the path is invalid or unwritable, Ralph must not start the AI; fail with exit 2 (R8).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| `--review-output` points to a non-existent directory | Fail before spawning AI; exit 2; message indicates path invalid or unwritable. |
| `--review-output` points to an existing file | Overwrite allowed; report is written to that file (or implementation may refuse overwrite and exit 2 — if so, document). Recommended: allow overwrite for simplicity. |
| `--review-output` path is a directory, not a file | Treat as invalid: exit 2; path must refer to a file (or non-existent path whose parent is writable). Do not auto-choose a filename inside the directory. |
| Temp directory is missing or unwritable | Fail before spawning AI; exit 2. |
| No `--review-output`; user runs in a context where temp is not available | Same as above; exit 2. |

### Examples

#### Explicit path

**Input:** `ralph review build --review-output ./report.txt`

**Expected output:** Report is written to `./report.txt` (relative to CWD). User knows path from the flag. R9 verifies file at `./report.txt` after AI exits.

#### Default temp path

**Input:** `ralph review build` (no `--review-output`)

**Expected output:** Report is written to a file in the system temp directory (e.g. `/tmp/ralph-review-abc123.md` or platform equivalent). Ralph communicates this path to the user (e.g. "Report written to /tmp/ralph-review-abc123.md"). R9 verifies the file at that path.

## Acceptance criteria

- [ ] When the user supplies `--review-output <path>`, the report is written to that path (when the review phase succeeds and the AI writes the report there).
- [ ] When the user does not supply `--review-output`, the report is written to a path in system temporary storage.
- [ ] When the default (temp) location is used, the user can discover the report path (e.g. it is printed or otherwise communicated).
- [ ] The report is always persisted to a file; there is no mode where the report exists only in memory or only on stdout without a file.
- [ ] Invalid or unwritable report path is handled as a review failure (exit 2 per R8).

## Dependencies

- R2 (review prompt composition) interpolates this path into the prompt so the AI knows where to write. R9 (report file verification) checks that the report exists at the expected path after the run.
