# R005: Prompt-writing guide discoverable

**Outcome:** O008 — Discoverability

## Requirement

A user who does not know how to write a Ralph prompt can discover prompt-writing guidance via documentation and/or a CLI command.

## Detail

The prompt-writing guide (O007/R005) is discoverable so that a new or uncertain user can find it without searching the repo. At least one of the following is true: (1) the path to first run or discoverability content (e.g. README, Quick Start) links to or mentions the guide, or (2) a CLI command (e.g. `ralph show prompt-guide`) outputs the full guide verbatim (single source of truth: docs/writing-ralph-prompts.md). The goal is that a user asking "how do I write a good prompt?" can get an answer from the docs or by running a command. The list command and help may reference the guide or the show prompt-guide command where appropriate.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User reads README or path to first run | They can find a link or mention of the prompt-writing guide (or the command to show it). |
| User runs the prompt-guide command | They see the full guide (same content as docs/writing-ralph-prompts.md). |
| User runs `ralph show --help` or `ralph --help` | Help text exposes the existence of `show prompt-guide` (or equivalent) so they can discover the command. |

### Examples

#### User discovers the guide from the README

**Input:** New user opens the README and looks for how to get started or what to run.

**Expected output:** README (or linked discoverability content) includes a reference to the prompt-writing guide—e.g. "New to writing prompts? See [Writing Ralph prompts](docs/writing-ralph-prompts.md)" or "Run `ralph show prompt-guide` for the full guide."

**Verification:** User can reach the guide or the command from the main entry point.

#### User runs ralph show prompt-guide

**Input:** User runs `ralph show prompt-guide` (or the implemented command name).

**Expected output:** CLI outputs the full "Writing Ralph prompts" guide verbatim (same content as docs/writing-ralph-prompts.md). Output may support an option (e.g. `--markdown`) for saving or piping to a pager.

**Verification:** User gets the full guidance without opening a file; the doc is the single source of truth and the command emits it in full.

## Acceptance criteria

- [ ] The prompt-writing guide is discoverable from at least one of: path to first run, README or discoverability content, or a CLI command.
- [ ] If a CLI command is provided (e.g. `ralph show prompt-guide`), it outputs the full guide verbatim; help for `show` (or equivalent) documents the command.
- [ ] A user asking "how do I write a good prompt?" can find the guide or the command from the documented discoverability entry points.

## Dependencies

- [O007/R005](../O007-user-documentation/R005-prompt-writing-guidance.md) — The guide that is to be discovered.
