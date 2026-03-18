# R003: Named Prompts with Overrides

**Outcome:** O002 — Configurable Behavior

## Requirement

The system supports named prompts in config with path (or content), optional display name and description, and optional loop overrides. Prompt definitions live in config; path resolution is relative to the config file that defines the prompt (not the current working directory).

## Detail

Users define prompts in config files under a **prompts** map. Each entry has a name used when running or listing (e.g. `ralph run build`), the path to the prompt file or inline content, optional display name and description for listing, and optional loop overrides so one prompt can have different limits or signals than another.

**Prompt path resolution:** A relative **path** is resolved relative to the **directory containing the config file that defined that prompt** (the layer that supplied the prompt when layers are merged). It is not relative to the current working directory. Absolute paths remain absolute. This allows config files to reference prompt files next to them or in a stable location relative to the config (e.g. `./prompts/build.md` or `prompts/build.md` from the same directory as the config file). See the [config component](../../engineering/components/config.md) for the canonical schema.

**Loop overrides:** Each prompt can specify its own loop settings (e.g. failure threshold, signal strings, AI command alias). When running or listing that prompt, those overrides apply over the root loop settings; environment and command-line options still override for that run (R001, R002).

**Display name and description:** Optional fields for listing (e.g. `ralph list`); when present, the list command shows them so users can discover prompts without reading config.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Prompt defined with relative path in workspace config | Path resolved relative to workspace config file directory. |
| Prompt defined with relative path when explicit config file is used | Path resolved relative to the explicit config file's directory. |
| Same prompt name in global and workspace | Workspace definition wins (layer merge per R001). |
| Prompt has loop overrides | When running or listing that prompt, prompt overrides apply over root; env/CLI still override for that run. |
| Prompt has no loop overrides | Root loop settings apply when running that prompt. |
| Prompt path does not exist at run time | Run fails with a clear error (e.g. file not found). List may still show the prompt name; behavior is defined. |
| Prompt defined with inline content instead of path | System supports content-based prompt; resolution and overrides apply the same. |

### Examples

#### Relative path from config directory

**Input:** Workspace config at `./ralph-config.yml` defines prompt "build" with path `./prompts/build.md`. User runs `ralph run build` from the repo root.

**Expected output:** Prompt file is loaded from `./prompts/build.md` relative to the directory containing the config file that defined "build" (e.g. same directory as `ralph-config.yml` if workspace supplied the prompt). Loop runs with that prompt.

**Verification:** Run from repo root; confirm the correct file is used. Move the config file and ensure path is still resolved relative to the config file directory, not cwd.

#### Per-prompt failure threshold

**Input:** Root loop sets failure threshold to 5. Prompt "cautious" has a loop override setting failure threshold to 1. User runs `ralph run cautious`.

**Expected output:** Loop exits after 1 consecutive failure when running "cautious".

**Verification:** Run "cautious" and trigger one failure; loop exits. Run another prompt without override; it uses the root value of 5.

#### List shows display name and description

**Input:** Config defines prompt "ci" with path `ci.md`, display name "CI pipeline", description "Runs the CI prompt." User runs `ralph list`.

**Expected output:** List output includes "ci" (and optional "CI pipeline" and description) so the user can discover what the prompt is for.

**Verification:** Run the list command; confirm prompt name and optional display name/description appear.

## Acceptance criteria

- [ ] The system supports named prompts in config: each prompt has a name, path (or content), and optional display name, description, and loop overrides.
- [ ] Relative prompt paths are resolved relative to the directory of the config file that defined that prompt (not the current working directory).
- [ ] When running or listing a prompt, that prompt's loop overrides (if any) apply over root loop settings; environment and command-line options still override for that run.
- [ ] When display name or description are defined for a prompt, the list command can show them (R006).
- [ ] Prompt definitions and their overrides live only in config files; they are not overridable by environment or CLI (the run still uses the chosen prompt's resolved settings, which may then be overridden by env/CLI for loop-wide settings).

## Dependencies

- R001 — Config layer resolution (prompts are merged from layers; path resolution uses the defining layer's config file directory).
- R002 — Loop behavior configurable (per-prompt loop overrides use the same settings and override order).
