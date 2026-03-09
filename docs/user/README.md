# User-facing documentation

This directory is the **canonical location** for user-facing documentation for Ralph: how-tos, workarounds, and examples. It is produced from or aligned with the intent tree under [docs/intent/](../intent/). When a topic corresponds to a requirement or outcome, that linkage is noted so the intent tree remains the source of truth and user docs stay consistent with specifications.

**Limitations:** Some described features (e.g. `ralph review`) are specified in the intent tree but not yet implemented; those docs state that or link to the spec. Install/uninstall is documented in the root README and tested on macOS and Linux; Windows is not yet supported.

## Topics

| Document | Description | Intent |
|----------|-------------|--------|
| [Configuration](configuration.md) | Config file location, key loop options, and precedence (config vs env vs CLI) | [O2](../intent/O2-configurable-behavior/README.md) — Configurable behavior |
| [AI backends and aliases](ai-backends-and-aliases.md) | Built-in aliases, user-defined aliases, and direct command; how to choose or override the AI CLI | [O3](../intent/O3-backend-agnosticism/README.md) — Backend agnosticism |
| [Cursor Agent workaround](cursor-agent-workaround.md) | Optional wrapper script for plain-text stdout and signal scanning when using the Cursor Agent backend | [O3](../intent/O3-backend-agnosticism/README.md) — Backend agnosticism; [R1](../intent/O3-backend-agnosticism/R1-builtin-aliases.md) built-in aliases |
| [Review report summary format](review-report-format.md) | Machine-parseable summary line and exit codes for `ralph review` report output | [O5](../intent/O5-prompt-review/README.md) — Prompt review; [R6](../intent/O5-prompt-review/R6-report-format-exit-codes.md) report format |
| [Ralph binary and build prompt](ralph-binary-and-build-prompt.md) | How to build the ralph binary, put it on PATH, and use the build procedure prompt for O5 tasks | [O5](../intent/O5-prompt-review/README.md) — Prompt review (developer workflow) |
| [Install and uninstall](../README.md#install-and-uninstall) | How to install Ralph on your system and uninstall it cleanly | [O6](../intent/O6-install-uninstall/README.md) — Install and uninstall; [R5](../intent/O6-install-uninstall/R5-install-uninstall-documentation.md) documentation |
| [Releases and versioning](release.md) | How versions and releases work; conventional commits and pre-release branches | — |
| [Publishing the first pre-release](first-prerelease.md) | Step-by-step: create rc/alpha/beta from impl, push, and publish first pre-release | — |

## Policy

This section defines what belongs where, how to add or update user-facing docs, and how they relate to the intent tree. Full specification: [O7 R5 — Policy and maintenance](../intent/O7-user-facing-documentation/R5-documentation-policy-and-maintenance.md).

### What goes where

- **User-facing docs (`docs/user/`):** Task-oriented content for *users* of Ralph: how to configure, run, choose a backend, apply workarounds, install and uninstall. Written in user terms; may summarize or link to intent or CLI reference. Must align with product behavior and be discoverable (e.g. linked from the root README).
- **Intent tree (`docs/intent/`):** Outcomes, requirements, and specifications for *builders* and maintainers. Source of truth for behavior. Not written as end-user how-tos. User docs may reference outcomes/requirements for traceability.
- **Root README:** High-level project description, quick start, install/uninstall summary or link, and link(s) to user docs. README may duplicate a minimal subset of user doc content for convenience; the canonical detail lives in user docs or intent as appropriate.

### Adding a topic

1. Create a new file under `docs/user/` (e.g. `docs/user/my-topic.md`). Use kebab-case for filenames.
2. Add a row to the **Topics** table in this index: document link, short description, and intent link (outcome and/or requirement from `docs/intent/`) when the topic implements or explains a specific requirement.
3. In the topic file, add an **Intent** line at the top linking to the relevant outcome/requirement so traceability is explicit.

### Updating and alignment

- When behavior or intent changes, update any user doc that describes that behavior so docs stay aligned with the product. When a requirement in the intent tree changes, update the linked user doc(s) that trace to it.
- User-facing docs must not contradict the intent tree. Where a topic traces to a requirement, that requirement is authoritative for behavior; the user doc translates it into user-oriented language.
- If a topic becomes obsolete (e.g. a workaround is no longer needed), remove or archive the doc and update the index.
