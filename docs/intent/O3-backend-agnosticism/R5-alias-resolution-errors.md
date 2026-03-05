# R5: AI Command Alias Resolution with Clear Errors

**Outcome:** O3 — Backend Agnosticism

## Requirement

The system resolves AI command aliases to executable command strings and produces clear, actionable error messages when resolution fails. The user should never see a cryptic "command not found" — the error should explain what alias was requested, that it wasn't found, and what aliases are available.

## Specification

Command resolution runs at startup, after config merge and before the first iteration. The effective command source is chosen per R6. Two kinds of resolution failure are defined here: **unknown alias** (user specified an alias name that is not in the merged alias map) and **no command configured** (no effective `ai_cmd` or `ai_cmd_alias` after merge).

**Merged alias map:** The set of aliases available for resolution is the merge of built-in aliases (R1) and user-defined aliases (R3). User-defined aliases override built-in aliases by name. Resolution looks up the effective `ai_cmd_alias` value in this merged map.

**When resolution succeeds:**
- Effective `ai_cmd` non-empty → the command string is that value (parsed per R2). No alias lookup.
- Effective `ai_cmd_alias` non-empty and present in merged map → the command string is the map value (parsed per R2).

**When resolution fails:**

| Failure | Condition | Required error behavior |
|---------|-----------|-------------------------|
| Unknown alias | User specified an alias name (via CLI, env, or config) that is not a key in the merged alias map | Exit with code 1. Error message must: (1) name the unknown alias, (2) state that it was not found or is not defined, (3) list the available aliases (keys of the merged map). Message must be written to stderr. |
| No command configured | After applying R6, neither effective `ai_cmd` nor effective `ai_cmd_alias` is set (or alias is set but resolution is not applicable until after we know we need an alias) | Exit with code 1. Error message must explain that an AI command must be configured (via `loop.ai_cmd` or `loop.ai_cmd_alias` in config, `RALPH_LOOP_AI_CMD` or `RALPH_LOOP_AI_CMD_ALIAS` in environment, or `--ai-cmd` / `--ai-cmd-alias` on the CLI). Message must be written to stderr. |

Resolution errors occur before the loop begins. No iteration runs. Ralph does not start the AI process.

**List of available aliases:** When reporting "unknown alias", the list of available aliases must reflect the merged map (built-in + user-defined, with user overrides applied). Order is implementation-defined; stability is not required.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| `--ai-cmd-alias typo` and no alias named `typo` exists | Error: unknown alias "typo"; list available aliases. Exit 1. |
| `loop.ai_cmd_alias: claude` in config; user renames built-in in code so `claude` is removed | Same as unknown alias: error names "claude" and lists current available aliases. |
| No `ai_cmd`, no `ai_cmd_alias` in any layer | Error: no AI command configured; explain how to set one. Exit 1. |
| `--ai-cmd ""` (empty string) provided | Empty string is "present but empty". Treated as no direct command; fall through to alias. If no alias either, no command configured. |
| `--ai-cmd-alias ""` (empty string) | Treated as no alias specified. If no direct command, no command configured. |
| Config has `ai_cmd_aliases: {}` and no built-ins (hypothetical) | Merged map has only built-ins (R1). If user sets alias to a name not in built-ins, unknown alias. |
| User defines alias `claude` in config (override) | Merged map has `claude` → user value. Resolution of `ai_cmd_alias: claude` uses user value. No error. |

### Examples

#### Unknown alias

**Input:**
```bash
ralph run build --ai-cmd-alias not-an-alias
```

**Expected output:**
(stderr) An error message such as: `unknown ai_cmd_alias "not-an-alias": no matching alias defined` (or equivalent). Followed by a list of available aliases, e.g. `available aliases: claude, kiro, copilot, cursor-agent` (and any user-defined aliases).

**Verification:**
- Exit code 1. No iteration runs.
- Message contains the requested alias name and the list of available aliases.

#### No command configured

**Input:**
Fresh install; no config files; no env; `ralph run -f ./prompt.md` (no `--ai-cmd` or `--ai-cmd-alias`).

**Expected output:**
(stderr) An error message explaining that an AI command must be configured, and how (config keys, env vars, CLI flags).

**Verification:**
- Exit code 1. Loop does not start.
- User can fix by setting e.g. `--ai-cmd-alias claude` or adding config.

#### Valid alias resolves

**Input:**
```bash
ralph run build --ai-cmd-alias kiro
```

**Expected output:**
Ralph resolves `kiro` to the command string from the merged alias map, parses it (R2), and starts the loop. No resolution error.

**Verification:**
- No stderr message about unknown alias or missing command. Loop proceeds (or fails later for other reasons, e.g. binary not found).

## Acceptance criteria

- [ ] When a valid alias is specified, Ralph resolves it to the corresponding command string
- [ ] When an unknown alias is specified, Ralph exits with an error that names the unknown alias and lists available aliases
- [ ] When no AI command is configured (no alias and no direct command), Ralph exits with an error explaining that an AI command must be configured
- [ ] Resolution errors occur at startup, before the loop begins

## Dependencies

_None identified._
