# O009: Predictability

## Who

Users who expect Ralph to leave their files and config unchanged unless they explicitly request a change (e.g. applying a review revision).

## Statement

Ralph only changes user content when the user explicitly requests it (e.g. by requesting that a suggested revision be applied and confirming).

## Why it matters

Unexpected edits to prompt files, config, or other user content cause loss of trust and broken workflows. Users must be able to run the loop, run the reviewer, list resources, and use other features knowing that their prompt files and configuration are not modified unless they opt in — for example by requesting that a review revision be applied and confirming (or using a non-interactive option where supported). Predictability makes the product safe to use in scripts and alongside other tools.

## Verification

- User runs the loop or runs the reviewer without requesting that a revision be applied. The prompt file on disk is unchanged; any in-memory preamble or composed prompt is not written back.
- User runs the reviewer with a prompt from a file path without requesting apply. That file is unchanged; the report and suggested revision are produced but not written to the user's prompt file unless the user requests apply.
- User requests that the suggested revision be applied and confirms (or uses a non-interactive option where supported). Only then does Ralph write the revised prompt to the chosen path (e.g. the source file or a path the user specified when the prompt came from stdin).
- Config files are read but not rewritten by Ralph; no automatic migration or normalization of user config unless the user invokes a documented opt-in flow that is described as modifying config.

## Non-outcomes

- Ralph may write to paths the user explicitly specifies for report or revised prompt output when the user has chosen those paths; the outcome is that the user's source prompt and config are not changed without explicit request.
- Ralph does not guarantee that third-party tools (e.g. the AI CLI) do not modify files; the outcome is about Ralph's own behavior.
- Creating or updating Ralph's own state (e.g. caches) in documented locations is out of scope of this outcome; the focus is user-owned content (prompts, config).

## Risks

| Risk | Mitigating Requirement |
|------|------------------------|
| Ralph writes to user files when the user did not request apply | [R001 — Explicit apply for writes](R001-explicit-apply-for-writes.md) |
| Config is rewritten or migrated without user opt-in | [R002 — Config read-only unless opt-in](R002-config-read-only-unless-opt-in.md) |
| In-memory or composed prompt is written back to source file without explicit apply | [R001 — Explicit apply for writes](R001-explicit-apply-for-writes.md) |
| User believes they are only reviewing but a revision is applied (ambiguous UX) | [R003 — Review–apply separation and confirmation](R003-review-apply-separation-and-confirmation.md) |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R001](R001-explicit-apply-for-writes.md) | Ralph writes to user prompt content (files or user-specified output paths) only when the user has explicitly requested that a revision be applied and has confirmed or used a documented non-interactive apply option. | ready |
| [R002](R002-config-read-only-unless-opt-in.md) | Ralph reads user config but does not rewrite or migrate it unless the user invokes a documented opt-in flow that is described as modifying config. | ready |
| [R003](R003-review-apply-separation-and-confirmation.md) | Ralph separates review (report and suggestions only, no writes) from apply (write revision) and requires confirmation for apply in interactive use unless a non-interactive option is used. | ready |
