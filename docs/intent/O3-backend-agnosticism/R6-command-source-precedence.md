# R6: Command Source Precedence

**Outcome:** O3 — Backend Agnosticism

## Requirement

The system resolves the AI command from multiple possible sources using a deterministic precedence order. A direct command string always takes precedence over an alias, regardless of where each is specified. This prevents ambiguity when both are configured.

## Specification

The AI command for a run is determined by resolving two values from the config merge: **effective `ai_cmd`** (direct command string) and **effective `ai_cmd_alias`** (alias name). Resolution uses the same precedence order as all loop settings (see O2-R6): CLI flags → environment variables → prompt-level loop overrides (when running a specific prompt) → workspace root loop → global root loop → built-in defaults.

**Built-in defaults for AI command:** There are no built-in defaults for `ai_cmd` or `ai_cmd_alias`. If no layer provides a value for either, command resolution fails (see R5).

**Within-layer rule:** When both `ai_cmd` and `ai_cmd_alias` are set at the same effective layer (e.g., both in workspace config), the direct command wins: the resolved command source is `ai_cmd`, and the alias is ignored for that run.

**Resolution steps:**

1. Compute the effective loop config for this run (merge order above; when running `ralph run <prompt-alias>`, prompt-level overrides apply).
2. From the effective config, read effective `ai_cmd` and effective `ai_cmd_alias`.
3. If effective `ai_cmd` is non-empty → use it as the command string. Resolution is "direct command"; do not look up an alias.
4. Else if effective `ai_cmd_alias` is non-empty → resolve the alias to a command string via the merged alias map (R1 + R3). Resolution is "alias &lt;name&gt;".
5. Else → no AI command configured; fail with R5 error before starting the loop.

**Sources and visibility:**

| Source | CLI | Environment | Config (loop) |
|--------|-----|-------------|---------------|
| Direct command | `--ai-cmd <string>` | `RALPH_LOOP_AI_CMD` | `loop.ai_cmd` (root or prompt-level) |
| Alias | `--ai-cmd-alias <string>` | `RALPH_LOOP_AI_CMD_ALIAS` | `loop.ai_cmd_alias` (root or prompt-level) |

When debug-level logging is enabled, the resolved command source must be logged (e.g., "direct command", "alias claude" with expanded command, or that resolution failed). This allows the user to verify which source was used.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| `--ai-cmd "cli --flag"` and `--ai-cmd-alias kiro` both provided | Direct command wins: use `cli --flag`. Alias is ignored. |
| Only `--ai-cmd-alias kiro` provided | Resolve alias `kiro` from merged alias map; use resulting command string. |
| No CLI flags; `RALPH_LOOP_AI_CMD=my-cli` set | Use `my-cli` as command string. |
| No CLI or env; config has `loop.ai_cmd: "x"` and `loop.ai_cmd_alias: claude` | Use `x` (direct command wins at same layer). |
| No CLI, env, or config values for either field | No AI command configured; R5 error, exit before loop. |
| Prompt-level `prompts.build.loop.ai_cmd_alias: kiro`; root has `loop.ai_cmd_alias: claude` | When running `build`, effective alias is `kiro`. When running another prompt with no override, effective alias is `claude`. |
| `--config path/to/file.yml` used | Only that file's loop and ai_cmd_aliases participate; no global/workspace merge for config (per O2-R6). |

### Examples

#### Direct command overrides alias at CLI

**Input:**
```bash
ralph run build --ai-cmd "custom-ai --headless" --ai-cmd-alias claude
```

**Expected output:**
Ralph uses `custom-ai --headless` as the command string. The alias `claude` is not consulted. Loop runs (or fails on first iteration if the binary is missing).

**Verification:**
- Run with `--log-level debug`. Log shows resolved command is the direct string (or equivalent).
- No alias resolution error; prompt is piped to `custom-ai --headless`.

#### Alias from config when no CLI flags

**Input:**
Config has `loop.ai_cmd_alias: kiro`. User runs `ralph run build` with no `--ai-cmd` or `--ai-cmd-alias`.

**Expected output:**
Ralph resolves alias `kiro` to the built-in (or user-overridden) command string for `kiro`, parses it per R2, and executes. Loop runs.

**Verification:**
- Debug log shows alias `kiro` and the expanded command.
- Process invoked is the resolved command (e.g. `kiro-cli chat --no-interactive --trust-all-tools` or user override).

#### No command configured

**Input:**
No config files; no env vars; user runs `ralph run build` without `--ai-cmd` or `--ai-cmd-alias`.

**Expected output:**
R5 error: no AI command configured. Message explains that an AI command must be set via config, env, or CLI. Ralph exits with code 1 before the loop starts.

**Verification:**
- Exit code 1. No iteration runs. Error message is actionable (see R5).

## Acceptance criteria

- [ ] If --ai-cmd is specified on the CLI, it is used regardless of any other setting
- [ ] If --ai-cmd-alias is specified on the CLI (and no --ai-cmd), the alias is resolved
- [ ] If neither CLI flag is specified, environment variables are checked (RALPH_LOOP_AI_CMD, RALPH_LOOP_AI_CMD_ALIAS)
- [ ] If no CLI flags or environment variables are set, config file values are used (loop.ai_cmd, loop.ai_cmd_alias)
- [ ] There is no built-in default for ai_cmd or ai_cmd_alias — if no layer provides a value, resolution fails (see O3-R5)
- [ ] At each precedence layer, a direct command (ai_cmd) takes precedence over an alias (ai_cmd_alias)
- [ ] The resolved command source is visible in debug-level logging

## Dependencies

_None identified._
