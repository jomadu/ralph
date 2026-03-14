# R006: List Prompts and Commands

**Outcome:** O002 — Configurable Behavior

## Requirement

The system provides a list command that shows available prompts and AI command aliases from the resolved config. The user can list all, only prompts, or only aliases. Users see names and descriptions for what is listed.

## Detail

Users need to discover what prompts and aliases are available without reading config files. The list command uses the same config resolution as the run command (defaults, global, workspace, or explicit file plus environment and command-line options). The output shows:

- **Prompts** — When listing prompts (or all), names (and optional display name), optional description, and path or other identifying info as documented. Prompts are those defined in the resolved config (R003).
- **AI commands (aliases)** — When listing aliases (or all), alias names (and optionally a short description or expansion if useful). Includes built-in and user-defined aliases from the resolved config (R004).

List is read-only and does not modify config or run a loop. It helps users verify that their config is loaded and that prompt/alias names are correct before running.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No prompts defined in config | List shows empty prompts or a clear "none" state; no error. |
| No custom AI aliases; only built-ins | List (or list aliases) shows built-in aliases. |
| Explicit config file specified | List uses that config only (same as run); shows prompts and aliases from that file. |
| Prompt path invalid for a defined prompt | List may still show the prompt name; running it fails with path error (R003). Or list reports the issue; behavior is defined. |
| Multiple config layers merge to define prompts | Listed prompts are those in the resolved config (after merge). |

### Examples

#### List all (prompts and aliases)

**Input:** Config defines prompts "build" and "test" and alias "my-ai". User runs the list command with no explicit config file.

**Expected output:** Output includes "build" and "test" (with optional display name, description, and path) and "my-ai" (and other built-in aliases as designed).

**Verification:** Run the list command; confirm both prompts and the alias appear.

#### List with explicit config

**Input:** User runs the list command with the documented config file option pointing to a specific file (e.g. one that defines a single prompt "ci-only").

**Expected output:** List shows only prompts and aliases from that file (e.g. "ci-only" and any aliases in the file or built-ins). No prompts from global or workspace.

**Verification:** Run the list command with the same explicit config file; output matches that file's definitions.

## Acceptance criteria

- [ ] The system provides a list command that shows prompts and aliases from the resolved config.
- [ ] The user can list all, only prompts, or only aliases.
- [ ] The list command uses the same config resolution as the run command (R001); when an explicit config file is used, only that file's definitions (plus built-ins where applicable) are shown.
- [ ] Output includes names for prompts and aliases; optional display name and description for prompts when defined (R003).
- [ ] When no prompts or no custom aliases are defined, list behaves without error (empty or minimal output as designed).

## Dependencies

- R001 — Config layer resolution (list uses resolved config).
- R003 — Named prompts (list shows prompts defined in config).
- R004 — AI command aliases (list shows aliases defined in config and built-ins).
