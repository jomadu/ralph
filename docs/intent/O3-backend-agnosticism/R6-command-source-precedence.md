# R6: Command Source Precedence

**Outcome:** O3 — Backend Agnosticism

## Requirement

The system resolves the AI command from multiple possible sources using a deterministic precedence order. A direct command string always takes precedence over an alias, regardless of where each is specified. This prevents ambiguity when both are configured.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] If --ai-cmd is specified on the CLI, it is used regardless of any other setting
- [ ] If --ai-cmd-alias is specified on the CLI (and no --ai-cmd), the alias is resolved
- [ ] If neither CLI flag is specified, environment variables are checked (RALPH_LOOP_AI_CMD, RALPH_LOOP_AI_CMD_ALIAS)
- [ ] If no CLI flags or environment variables are set, config file values are used (loop.ai_cmd, loop.ai_cmd_alias)
- [ ] There is no built-in default for ai_cmd or ai_cmd_alias — if no layer provides a value, resolution fails (see O3-R5)
- [ ] At each precedence layer, a direct command (ai_cmd) takes precedence over an alias (ai_cmd_alias)
- [ ] The resolved command source is visible in debug-level logging
