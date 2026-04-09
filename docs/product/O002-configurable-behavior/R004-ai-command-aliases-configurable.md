# R004: AI Command Aliases Configurable

**Outcome:** O002 — Configurable Behavior

## Requirement

The system supports configurable AI command aliases: short names that expand to a full AI CLI command. Built-in aliases exist by default; the user can add or override them in config. Alias resolution uses the same layer order as other config (R001); the selected alias is chosen via config or CLI.

## Detail

Users run Ralph with an AI CLI (e.g. Claude, Cursor agent, Kiro). Instead of typing the full command every time, they can use a short alias (e.g. `claude`, `cursor-agent`) that expands to the full command. Built-in aliases for known AI CLIs exist so the tool works out of the box; the user can add custom aliases or override built-ins in config files so the same name works across global or workspace config.

**Layer order:** Alias definitions and the chosen alias (e.g. `loop.ai_cmd_alias`) follow the same layer order as R001: defaults, global file, workspace file, explicit file, environment, prompt-level overrides, command-line options. User-defined aliases in config override built-in aliases for the same name. The effective AI command for a run is resolved from the effective loop settings (direct command if set, otherwise alias expansion); see the [config component](../../engineering/components/config.md) for the canonical schema and built-in list.

**Selection:** The user sets the alias (or direct AI command) in config (root or per prompt), via environment variables, or via CLI options. When both a direct command and an alias are set, the direct command takes precedence.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No user aliases; only built-ins | Run and list use built-in aliases (e.g. `claude`, `cursor-agent`). |
| User defines an alias with the same name as a built-in | User alias overrides the built-in for that name. |
| User defines an alias in workspace; same name in global | Workspace alias wins (layer merge per R001). |
| Config sets ai_cmd_alias to a name that does not exist | Run fails with a clear error (unknown alias). |
| Config sets ai_cmd (direct command) and ai_cmd_alias | Direct command is used; alias is ignored for that run. |
| Alias selected via CLI for a run | That run uses the CLI-specified alias (or direct command); overrides config. |

### Examples

#### Built-in alias

**Input:** No config file. User runs `ralph run` with default alias (e.g. `cursor-agent` from defaults) or with an alias chosen by the tool.

**Expected output:** The run invokes the AI CLI command that the built-in alias expands to (e.g. the Cursor agent command). No config required.

**Verification:** Run without config; confirm the expected AI command is invoked (e.g. via dry-run or log).

#### User override of built-in

**Input:** Workspace config defines `aliases.claude: "claude-custom --my-flag"`. Root loop sets `ai_cmd_alias: claude`. User runs a prompt.

**Expected output:** The run uses `claude-custom --my-flag` (user alias overrides built-in).

**Verification:** Run and confirm the custom command is used (e.g. via dry-run or observability).

#### Direct command overrides alias

**Input:** Config sets `loop.ai_cmd: "custom-ai --batch"` and `loop.ai_cmd_alias: claude`. User runs a prompt.

**Expected output:** The run uses `custom-ai --batch`; the alias is ignored when a direct command is set.

**Verification:** Run and confirm the direct command is used.

#### List shows aliases

**Input:** Config defines alias "my-ai" with command "my-cli --non-interactive". User runs `ralph list` (or list aliases).

**Expected output:** List shows built-in aliases and "my-ai" (and optionally the expansion), so the user can see what is available.

**Verification:** Run the list command; confirm user and built-in aliases appear (R006).

## Acceptance criteria

- [ ] The system provides built-in AI command aliases (e.g. for known AI CLIs) so the tool works without a config file.
- [ ] The user can define additional aliases in config (global, workspace, or explicit file); user aliases override built-ins for the same name.
- [ ] Alias definitions and the selected alias (direct command or alias name) follow the same config layer order as R001.
- [ ] The user can select which alias or direct command to use via config (root or per prompt), environment variables, or command-line options; direct command takes precedence over alias when both are set.
- [ ] When the selected alias name is unknown (not built-in and not in resolved config), the system reports a clear error and does not proceed with the run.
- [ ] Run and review use the resolved AI command (alias expansion or direct command); list shows available aliases (R006).

## Dependencies

- R001 — Config layer resolution (alias definitions and selection use the same layers and override order).
