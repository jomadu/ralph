# O008: Discoverability

## Who

A new user (or someone evaluating Ralph) who wants to understand what Ralph does and achieve a first successful run with minimal friction.

## Statement

A new user can discover what Ralph does and get to a first successful run.

## Why it matters

If the value proposition and first-run path are unclear, potential users leave before they see the product work. Discovery answers "what is this?" and "how do I try it?" — so a new user can go from zero to a working loop (or review) without reverse-engineering the repo or guessing at commands.

## Verification

- A new user can find a short description of what Ralph is (e.g. loop runner for AI-driven tasks) and why it exists (manual read–judge–re-run replaced by automated iteration until a signal).
- A new user can find steps to install Ralph (or run it in a documented way) and run a first command (e.g. `ralph run <alias>` or `ralph review -f <path>`) that completes successfully.
- Listing or help commands (e.g. `ralph list prompts`, `ralph --help`) expose available prompts and subcommands so the user can see what to run.
- The path from "I have Ralph" to "I just ran a loop and it exited 0" is documented and achievable without prior knowledge of the codebase.

## Non-outcomes

- Discoverability does not require interactive onboarding, videos, or in-app wizards. The outcome is that the product and docs make the first run achievable.
- Ralph does not ship with pre-loaded prompts or mandatory samples; the user may use a minimal prompt or one they already have, as long as the path to a first successful run is clear.
- Marketing or positioning copy is out of scope; the outcome is technical discoverability (what it is, how to run it once).

## Risks

| Risk | Mitigating Requirement |
|------|------------------------|
| User does not understand what Ralph is or why it exists | [R001 — What and why](R001-what-and-why.md) |
| User cannot find how to install or run a first command | [R002 — Install and first command](R002-install-and-first-command.md) |
| User does not know which prompts or subcommands to run | [R003 — List and help](R003-list-and-help.md) |
| Path from "have Ralph" to first successful run is unclear | [R004 — Path to first run](R004-path-to-first-run.md) |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R001](R001-what-and-why.md) | User can find a short description of what Ralph is and why it exists | draft |
| [R002](R002-install-and-first-command.md) | User can find install steps and a first command that completes successfully | draft |
| [R003](R003-list-and-help.md) | List or help commands expose available prompts and subcommands | draft |
| [R004](R004-path-to-first-run.md) | Path from having Ralph to a first successful run is documented and achievable | draft |
