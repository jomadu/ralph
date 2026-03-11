# R006: List Prompts and Commands

**Outcome:** O002 — Configurable Behavior

## Requirement

The system provides listing commands that show available prompts and AI commands. Users can see what prompts and aliases are defined in the resolved config and their names and descriptions.

## Detail

Users need to discover what prompts and AI commands are available without reading config files. The listing command uses the same config resolution as the run command (defaults, global, workspace, or explicit file plus environment and command-line options). The output shows:

- **Prompts** — Names (and optional display name), optional description, and path or other identifying info as documented. Prompts are those defined in the resolved config (R003).
- **AI commands** — Alias names (and optionally a short description or expansion if useful). Includes built-in and user-defined aliases from the resolved config (R004).

Listing is read-only and does not modify config or run a loop. It helps users verify that their config is loaded and that prompt/alias names are correct before running.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No prompts defined in config | Listing shows empty prompts or a clear "none" state; no error. |
| No custom AI aliases; only built-ins | Listing shows built-in commands. |
| Explicit config file specified | Listing uses that config only (same as run); shows prompts and commands from that file. |
| Prompt path invalid for a defined prompt | Listing may still show the prompt name; running it fails with path error (R003). Or listing reports the issue; behavior is defined. |
| Multiple config layers merge to define prompts | Listed prompts are those in the resolved config (after merge). |

### Examples

#### List prompts and commands

**Input:** Config defines prompts "build" and "test" and alias "my-ai". User runs the listing command (no explicit config).

**Expected output:** Output includes "build" and "test" (with optional display name, description, and path) and "my-ai" (and other built-in aliases as designed).

**Verification:** Run the listing command; confirm both prompts and the alias appear.

#### List with explicit config

**Input:** User runs the listing command with the documented config file option pointing to a specific file (e.g. one that defines a single prompt "ci-only").

**Expected output:** Listing shows only prompts and commands from that file (e.g. "ci-only" and any aliases in the file or built-ins). No prompts from global or workspace.

**Verification:** Run the listing command with the same explicit config file; output matches that file's definitions.

## Acceptance criteria

- [ ] The system provides at least one listing command (or subcommand) that shows available prompts from the resolved config.
- [ ] The same or related listing shows available AI commands (aliases) from the resolved config, including built-ins.
- [ ] The listing command uses the same config resolution as the run command (R001); when an explicit config file is used, only that file's definitions (plus built-ins where applicable) are shown.
- [ ] Output includes names for prompts and commands; optional display name and description for prompts when defined (R003).
- [ ] When no prompts or no custom commands are defined, listing behaves without error (empty or minimal output as designed).

## Dependencies

- R001 — Config layer resolution (listing uses resolved config).
- R003 — Named prompts (listing shows prompts defined in config).
- R004 — AI command aliases (listing shows aliases defined in config and built-ins).
