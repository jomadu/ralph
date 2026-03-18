# O002: Configurable Behavior

## Who

Users who run Ralph in different contexts (one-shot bootstrap, long build loops, cautious exploration) and need different iteration limits, failure thresholds, timeouts, or signal strings without maintaining variant copies of the same prompt.

## Statement

Loop execution adapts to the user's constraints without changing the prompt file.

## Why it matters

Different tasks need different loop parameters. A one-shot bootstrap needs one iteration with no preamble. A long build loop needs many iterations with a high failure threshold. A cautious exploration needs a short timeout and low threshold. Without external configuration, these differences live inside the prompt file — the user maintains variant copies of the same prompt for different environments, projects, or risk tolerances. Configuration separates loop behavior from prompt content so the same prompt file works across contexts.

## Configuration layers

Configuration is merged from multiple layers with a defined override order; the full layer list and order are specified in [R001](R001-config-layer-resolution.md). Layers (summary): defaults, global file, workspace file, explicit file, environment, prompt-level overrides, CLI flags.

## Configuration scope

- **Loop behavior** — [R002](R002-loop-behavior-configurable.md) (and [R001](R001-config-layer-resolution.md) for layers).
- **Prompts** — [R003](R003-named-prompts-with-overrides.md).
- **AI commands** — [R004](R004-ai-command-aliases-configurable.md).
- **Where config is loaded from** — [R001](R001-config-layer-resolution.md), [R005](R005-explicit-config-file-only.md), and user documentation.

## Verification

- User sets iteration limits or overrides in global or workspace config and on the command line. The loop runs according to those constraints; command-line overrides win over config.
- User defines a prompt with a custom failure threshold and signal strings for that prompt in config. Those values take effect when running that prompt without affecting others.
- User sets environment variables. Ralph applies them without any config file change; environment overrides file-based config.
- User points Ralph at a specific config file. Only that file is used; global and workspace config are not loaded. If the file is missing, Ralph reports an error.
- User runs the list command and sees available prompts and AI commands with names and descriptions; they can list all, only prompts, or only aliases as the product allows.
- User can view the effective (resolved) config for the current context, including which layer supplied each value when supported.

## Non-outcomes

- Ralph does not provide a GUI, interactive config editor, or a config set/get subcommand. Configuration is files, environment variables, and flags. Read-only viewing of effective (resolved) config is in scope (R007).
- Ralph does not support runtime config changes during loop execution. Config is resolved once at startup.
- Ralph does not validate prompt file content — only that the file exists and is readable.
- Ralph does not support config inheritance between prompts. Each prompt independently overrides the root loop section where applicable.

## Risks

| Risk | Mitigating Requirement |
|------|------------------------|
| Unclear which config value applies for a setting | [R001 — Config layer resolution](R001-config-layer-resolution.md), [R007 — View effective config](R007-view-effective-config.md) |
| User cannot see which config value is active for a run | [R007 — View effective config](R007-view-effective-config.md) |
| Explicit config file missing when specified | [R005 — Explicit config file only](R005-explicit-config-file-only.md) |
| Per-prompt overrides not applied when running that prompt | [R003 — Named prompts with overrides](R003-named-prompts-with-overrides.md) |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R001](R001-config-layer-resolution.md) | The system resolves configuration from defined layers (defaults, global file, workspace file, explicit file, environment, prompt-level overrides, CLI flags) with a defined override order. | ready |
| [R002](R002-loop-behavior-configurable.md) | The system allows loop behavior (iterations, failure threshold, timeout, signals, preamble, AI command, streaming, log level, max output buffer) to be configured at root and per prompt. | ready |
| [R003](R003-named-prompts-with-overrides.md) | The system supports named prompts in config with path, optional display name and description, and optional loop overrides. | ready |
| [R004](R004-ai-command-aliases-configurable.md) | The system supports configurable AI command aliases. | ready |
| [R005](R005-explicit-config-file-only.md) | When the user specifies an explicit config file, the system uses only that file and reports an error if it is missing. | ready |
| [R006](R006-list-prompts-and-commands.md) | The system provides a list command that shows available prompts and aliases; the user can list all, only prompts, or only aliases. | ready |
| [R007](R007-view-effective-config.md) | The user can view the effective (resolved) configuration for the current context, including optional provenance (which layer supplied each value). | ready |
