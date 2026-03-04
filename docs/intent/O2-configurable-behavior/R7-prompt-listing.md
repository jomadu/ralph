# R7: Resource Listing Command

**Outcome:** O2 — Configurable Behavior

## Requirement

The system provides a command to list configured resources — prompt aliases and AI command aliases — so users can discover what is available without reading config files. Both built-in and user-defined resources are shown.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] `ralph list prompts` outputs all prompt aliases defined in the resolved configuration (merged global + workspace)
- [ ] Each prompt entry shows the alias key, display name (if set), description (if set), and prompt file path
- [ ] If no prompts are configured, the output clearly indicates that no prompts are available
- [ ] `ralph list aliases` outputs all AI command aliases (built-in and user-defined, merged)
- [ ] Each AI command alias entry shows the alias name and the resolved command string
- [ ] User-defined aliases that override a built-in are indicated as such
- [ ] The output is human-readable and formatted for terminal display
