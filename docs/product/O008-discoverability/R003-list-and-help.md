# R003: List and Help

**Outcome:** O008 — Discoverability

## Requirement

The user can discover available prompts, aliases, and subcommands via the list command or help so they can see what to run.

## Detail

A new user needs to see what is available: which prompts and AI command aliases are defined in the resolved config, and which top-level commands (run, review, list, show, version) exist. The list command (e.g. `ralph list`, `ralph list prompts`, `ralph list aliases`) exposes prompts and aliases from the effective config. Help (e.g. `ralph --help`, `ralph run --help`) exposes subcommands and flags. Together these give the user enough information to choose a prompt or alias and run a command without reading the codebase.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No config or empty config | List shows empty or minimal output; help still shows subcommands and usage. |
| User runs list before any config | List reflects resolved config (defaults + any found files); user sees what would be used. |
| User runs help at top level | Top-level help lists commands (run, review, list, show, version) and points to per-command help. |
| User runs help for a subcommand | Per-command help shows flags and usage for that command. |

### Examples

#### List prompts and aliases

**Input:** User runs `ralph list` (or `ralph list prompts`, `ralph list aliases`) with a config that defines prompts and aliases.

**Expected output:** Output includes names of prompts and aliases; format is documented (e.g. table, YAML, JSON). User can see what to pass to `ralph run <alias>` or which alias to use.

**Verification:** User can identify at least one prompt or alias name and use it in a subsequent run or review command.

#### Top-level help

**Input:** User runs `ralph --help`.

**Expected output:** Help lists run, review, list, show, version and indicates how to get help for each (e.g. `ralph run --help`).

**Verification:** User can choose a subcommand and know how to get more help.

## Acceptance criteria

- [ ] The list command (and list prompts / list aliases) outputs available prompts and aliases from resolved config in a documented format.
- [ ] Top-level help lists all subcommands and how to get per-command help.
- [ ] Per-command help documents flags and usage for run, review, list, show, version.
- [ ] A user can go from "I have Ralph" to "I know which prompt or alias to run" using only list and help.

## Dependencies

- O002 (config) — List and show depend on resolved config; config layer resolution and prompt/alias definitions must exist.
