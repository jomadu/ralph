# R7: Configurable Review Output to Stdout

**Outcome:** O5 — Prompt review

## Requirement

The system allows the user to control whether the AI command output and the review report are exposed to stdout. By default, both the AI command output and the report are written to stdout (in addition to the report being saved to a file per R3). The user can change this behavior (e.g. via a flag or config) so that stdout is suppressed or limited for scripting or quieter operation, while the report remains written to its file.

## Specification

**Default behavior:** When a review run executes, both (1) the raw output of the AI command (stdout/stderr from the AI process) and (2) the report content (the file content at the review output path, after the AI has written it) are exposed to the user's stdout (or stderr for AI stream, if implementation separates them). So by default the user sees the AI stream and then the report; the report is also always saved to the file per R3.

**Configurable suppression:** The user can change this so that:
- The AI command output is not printed to stdout (e.g. `--quiet` or `show_ai_output: false` in a review-specific or global config), and/or
- The report content is not printed to stdout (e.g. report-only-to-file: only the file is written, nothing echoed to stdout).

Mechanism: a CLI flag (e.g. `--quiet` for no AI output, `--no-print-report` or `--report-to-file-only` for not printing report) and/or config keys under the review or global config (e.g. `review.show_ai_output`, `review.print_report`). At least one of flag or config must be available; the requirement is that the user can achieve "report only to file, nothing to stdout" for scripting.

**Invariants:** (1) Configuring stdout never prevents the report from being written to the file at the review output path (R3). (2) Behavior is the same for alias, file, and stdin input modes — only the source of the prompt differs; stdout behavior is independent of input mode.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User sets "no stdout" for both AI output and report | Report is still written to file; stdout receives nothing (or only minimal progress/stderr per implementation). |
| Verbose / default: both AI output and report to stdout | Report content may appear twice (once from AI stream if AI echoes it, once from Ralph reading file and printing) — acceptable; or implementation may deduplicate (e.g. only print report from file). Document. |
| Non-interactive / no TTY | Same rules; flags/config control what goes to stdout. |

### Examples

#### Default

**Input:** `ralph review build`

**Expected output:** AI command output streams to stdout (or stderr); after AI exits, report file is read and its content is also written to stdout (or only path message if report is long — implementation may choose). Report file exists at R3 path.

#### Report only to file

**Input:** `ralph review build --report-to-file-only` (or equivalent)

**Expected output:** Report is written to the R3 path; report content is not printed to stdout. AI output may still be shown unless also disabled.

## Acceptance criteria

- [ ] By default, when review runs, the AI command output and the report content are exposed to stdout (and the report is also saved to the path per R3).
- [ ] The user can configure or flag so that the AI command output and/or the report are not printed to stdout (e.g. report-only-to-file mode).
- [ ] Configuring stdout does not prevent the report from being written to the file at the review output path (R3); it only affects what is sent to stdout.
- [ ] Behavior is consistent for alias, file, and stdin input modes.

## Dependencies

- R3 (review output path) ensures the report is always in a file; this requirement only governs what additionally goes to stdout.
