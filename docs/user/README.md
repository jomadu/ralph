# User-facing documentation

This directory is the **canonical location** for user-facing documentation for Ralph: how-tos, workarounds, and examples. It is produced from or aligned with the intent tree under [docs/intent/](../intent/). When a topic corresponds to a requirement or outcome, that linkage is noted so the intent tree remains the source of truth and user docs stay consistent with specifications.

## Topics

| Document | Description | Intent |
|----------|-------------|--------|
| [Cursor Agent workaround](cursor-agent-workaround.md) | Optional wrapper script for plain-text stdout and signal scanning when using the Cursor Agent backend | [O3](../intent/O3-backend-agnosticism/README.md) — Backend agnosticism; [R1](../intent/O3-backend-agnosticism/R1-builtin-aliases.md) built-in aliases |
| [Review report summary format](review-report-format.md) | Machine-parseable summary line and exit codes for `ralph review` report output | [O5](../intent/O5-prompt-review/README.md) — Prompt review; [R6](../intent/O5-prompt-review/R6-report-format-exit-codes.md) report format |
| [Ralph binary and build prompt](ralph-binary-and-build-prompt.md) | How to build the ralph binary, put it on PATH, and use the build procedure prompt for O5 tasks | [O5](../intent/O5-prompt-review/README.md) — Prompt review (developer workflow) |

## How to add a topic

1. Create a new file under `docs/user/` (e.g. `docs/user/my-topic.md`). Use kebab-case for filenames.
2. Add a row to the **Topics** table above: document link, short description, and intent link (outcome and/or requirement from `docs/intent/`) when applicable.
3. In the topic file, add an **Intent** line at the top linking to the relevant outcome/requirement so traceability is explicit.

## Policy

- User-facing docs here should not contradict the intent tree. When behavior is specified in a requirement, user docs describe the same behavior in user terms.
- New topics that stem from the intent tree (e.g. workarounds referenced in a requirement) should be added here and linked from the relevant requirement or outcome where appropriate.
