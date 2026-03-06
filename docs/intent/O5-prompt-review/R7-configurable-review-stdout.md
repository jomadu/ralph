# R7: Configurable Review Output to Stdout

**Outcome:** O5 — Prompt review

## Requirement

The system allows the user to control whether the AI command output and the review report are exposed to stdout. By default, both the AI command output and the report are written to stdout (in addition to the report being saved to a file per R3). The user can change this behavior (e.g. via a flag or config) so that stdout is suppressed or limited for scripting or quieter operation, while the report remains written to its file.

## Specification

(To be specified in Step 5.)

## Acceptance criteria

- [ ] By default, when review runs, the AI command output and the report content are exposed to stdout (and the report is also saved to the path per R3).
- [ ] The user can configure or flag so that the AI command output and/or the report are not printed to stdout (e.g. report-only-to-file mode).
- [ ] Configuring stdout does not prevent the report from being written to the file at the review output path (R3); it only affects what is sent to stdout.
- [ ] Behavior is consistent for alias, file, and stdin input modes.

## Dependencies

- R3 (review output path) ensures the report is always in a file; this requirement only governs what additionally goes to stdout.
