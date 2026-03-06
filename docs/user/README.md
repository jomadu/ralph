# User-facing documentation

This directory holds user-facing documentation for Ralph: how-tos, workarounds, and examples. It is **produced from or aligned with** the objective intent tree under `docs/intent/`. When a topic corresponds to a requirement or outcome, that linkage is noted so the intent tree remains the source of truth and user docs stay consistent with specifications.

## Topics

| Document | Description | Intent |
|----------|-------------|--------|
| [Cursor Agent workaround](cursor-agent-workaround.md) | Optional wrapper script for plain-text stdout and signal scanning when using the Cursor Agent backend | O3 — Backend agnosticism; R1 built-in aliases |

## Policy

- User-facing docs here should not contradict the intent tree. When behavior is specified in a requirement, user docs describe the same behavior in user terms.
- New topics that stem from the intent tree (e.g. workarounds referenced in a requirement) should be added here and linked from the relevant requirement or outcome where appropriate.
