# O002: Configurable Behavior

## Who

Users who run Ralph in different contexts (one-shot bootstrap, long build loops, cautious exploration) and need different iteration limits, failure thresholds, timeouts, or signal strings without maintaining variant copies of the same prompt.

## Statement

Loop execution adapts to the user's constraints without changing the prompt file.

## Why it matters

Different tasks need different loop parameters. A one-shot bootstrap needs one iteration with no preamble. A long build loop needs many iterations with a high failure threshold. A cautious exploration needs a short timeout and low threshold. Without external configuration, these differences live inside the prompt file — the user maintains variant copies of the same prompt for different environments, projects, or risk tolerances. Configuration separates loop behavior from prompt content so the same prompt file works across contexts.

## Configuration layers

Configuration is merged from multiple layers. Later layers override earlier ones for the same setting. This lets the user keep shared defaults, override per project or per prompt, and still override once at run time without editing files.

1. **Defaults** — Built-in values so Ralph works out of the box (e.g. iteration limit, failure threshold, signal strings, and built-in AI commands). No config file required.

2. **Global config file** — User-level config that applies across all projects. Stored in the user’s config directory (platform-specific; may be overridden by an environment variable). Optional: if the file is missing, Ralph skips it and continues. Suited for personal preferences (default AI command, timeouts, log level) that should apply everywhere.

3. **Workspace config file** — Project-level config in the current working directory. Optional; if missing, skipped. When both global and workspace exist, workspace overrides global for the same settings. Suited for project-specific prompts, loop limits, or AI command so the same repo behaves consistently for everyone working in it.

4. **Explicit config file** — When the user points Ralph at a specific config file (e.g. via a CLI option), that file is the only file-based source: global and workspace are not loaded. The file must exist or Ralph reports an error. Suited for tests, CI, or running with an alternate config without changing the current directory or user config.

5. **Environment variables** — Variables that override file-based config without editing files. Useful for scripts and CI (e.g. set a timeout or AI command for a single job) or for temporarily changing behavior. An environment variable can also control where Ralph looks for the global config file, so the user can isolate or relocate their user-level config.

6. **Prompt-level overrides** — In config files, each prompt can specify its own loop settings (e.g. a different failure threshold or signal strings for that prompt only). Those overrides apply when running or listing that prompt and take precedence over the root loop settings, but are still overridden by environment variables and CLI flags for that run.

7. **CLI flags** — Command-line options that override all other layers for that run. Suited for one-off overrides (e.g. run with a different iteration limit or AI command without touching config or env).

## Configuration scope

**Loop behavior** — Maximum iterations, iteration mode (bounded vs unlimited), consecutive failure threshold, per-iteration timeout, output limits, success and failure signal strings, whether to inject a preamble, which AI command or AI command alias to use, whether to stream AI output to the terminal, log level, and max output buffer. Configurable at the root (default for all prompts), per prompt in config files, and overridable by environment variables and CLI for a run.

**Prompts** — The user defines named prompts in config files: a name used when running or listing, the path to the prompt file, optional display name and description for listing, and optional loop overrides so one prompt can have different limits or signals than another. Prompt file paths in config are **relative to the config file that defines them** (not the current working directory). Prompt definitions and their overrides live only in config files; they are not overridable by environment or CLI (the run still uses the chosen prompt’s resolved settings, which may then be overridden by env/CLI for loop-wide settings).

**AI commands** — Short names that expand to a full AI CLI command. Built-in commands for known AI CLIs exist by default; the user can add or override them in config files so the same name (e.g. for a proprietary tool) works across global or workspace config.

**Where config is loaded from** — A user-level (global) config directory, a project-level (workspace) config file in the current working directory, or a single explicit file path supplied by the user. Documentation defines the exact locations and how to override them (e.g. via environment).

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
