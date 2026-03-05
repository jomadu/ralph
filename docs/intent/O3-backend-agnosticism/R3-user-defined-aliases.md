# R3: User-Defined Command Aliases

**Outcome:** O3 — Backend Agnosticism

## Requirement

The system allows users to define custom AI command aliases in config files. User-defined aliases merge with built-in aliases, and user-defined aliases with the same name as a built-in alias override the built-in. This supports proprietary, internal, or newly released AI CLIs that Ralph doesn't ship aliases for.

## Specification

Users define custom AI command aliases in config files under the **`ai_cmd_aliases`** section. The section is a map: keys are alias names (strings), values are command strings (strings). Key and value are both required per entry; schema validation (O2-R3) applies to structure.

**Config file scope:** User-defined aliases are loaded from the same config files as the rest of Ralph config: global config and workspace config when not using `--config`, or the single file specified by `--config` (O2-R6). Workspace config and global config are merged; workspace entries override global entries for the same key. The exact file paths for "global" and "workspace" are defined in O2 (e.g. `ralph-config.yml` in workspace root, `~/.config/ralph/ralph-config.yml` or XDG equivalent for global).

**Merge with built-in:** The **merged alias map** used for resolution (R5, R6) is: built-in aliases (R1) as the base, then user-defined aliases overlaid. For any key present in both, the user-defined value wins. Keys only in built-in or only in user-defined all appear in the merged map. So: user-defined aliases merge with built-in (both available), and user-defined override built-in when the name is the same.

**Command string parsing:** Each user-defined alias value is a command string. It is parsed using the same shell-style rules as direct commands (R2) when the alias is resolved and executed. Ralph does not parse or validate the command string at config load time beyond schema (e.g. min length if required); parse errors occur at resolution/execution time when that alias is used.

**Naming:** Alias names are the map keys. No reserved names; the user may override any built-in by redefining it in config. Alias names are case-sensitive (implementation may use string keys as-is).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Workspace config has `ai_cmd_aliases.my-tool: "my-ai --stdin"` | Merged map includes `my-tool` → `my-ai --stdin`. Both built-ins and `my-tool` are available. |
| Global config has `ai_cmd_aliases.claude: "claude -p"`; workspace has no ai_cmd_aliases | Merged map: built-ins plus `claude` → `claude -p`. User override for `claude` only. |
| Workspace and global both define `ai_cmd_aliases.other: "..."` with different values | Workspace value wins (workspace overrides global for same key). |
| User defines alias with value `cli --model "model name"` | Value stored as one string; when resolved, parsed per R2 into argv. Quoting is preserved in the string and interpreted by R2. |
| Config has `ai_cmd_aliases: {}` | Empty user map. Merged map = built-in aliases only. |
| `--config other.yml` used; other.yml has ai_cmd_aliases | Only other.yml's ai_cmd_aliases participate. No merge with default workspace/global files. |

### Examples

#### Custom alias

**Input:**
Workspace `ralph-config.yml`:
```yaml
ai_cmd_aliases:
  my-ai: "my-ai-cli --headless --stdin"
```
User runs `ralph run build --ai-cmd-alias my-ai`.

**Expected output:**
Ralph resolves `my-ai` to `my-ai-cli --headless --stdin`, parses it (R2), and executes. Loop runs.

**Verification:**
- No unknown-alias error. Process invoked is `my-ai-cli` with the given arguments.

#### User override and merge

**Input:**
Workspace config:
```yaml
ai_cmd_aliases:
  claude: "claude -p --model claude-3-5-sonnet"
  internal: "internal-ai --api-key-from-env"
```
User runs `ralph run build --ai-cmd-alias internal`.

**Expected output:**
Merged map has `claude` (user), `internal` (user), and all other built-ins (kiro, copilot, cursor-agent). Resolution of `internal` yields `internal-ai --api-key-from-env`. Loop runs.

**Verification:**
- `--ai-cmd-alias claude` uses user value; `--ai-cmd-alias kiro` uses built-in. Both available.

#### Workspace overrides global alias

**Input:**
Global config: `ai_cmd_aliases.proprietary: "proprietary-cli"`. Workspace config: `ai_cmd_aliases.proprietary: "proprietary-cli --config ./local.yml"`. User runs `ralph run build --ai-cmd-alias proprietary` from workspace.

**Expected output:**
Resolved command is the workspace value: `proprietary-cli --config ./local.yml`.

**Verification:**
- Workspace overlay wins for the same key.

## Acceptance criteria

- [ ] Users can define aliases under the `ai_cmd_aliases` section in config files
- [ ] User-defined aliases merge with built-in aliases — both are available simultaneously
- [ ] A user-defined alias with the same name as a built-in alias overrides the built-in
- [ ] Aliases defined in workspace config override aliases with the same name defined in global config
- [ ] User-defined alias values are parsed using the same shell-style command parsing as direct commands

## Dependencies

_None identified._
