# Config

## Responsibility

The config component resolves configuration from defined layers into a single effective configuration used by the run-loop, review, list, and show paths. It supplies prompt and alias definitions, loop behavior settings (iterations, failure threshold, timeout, signals, precedence mode, preamble, AI command, streaming, log level), and read-only semantics unless the user opts in to writes (e.g. apply). It does not write files unless the user explicitly requests an action that writes (e.g. review `--apply`); config file resolution only reads.

Implements the requirements assigned to this component in the [engineering README](../README.md).

## Interfaces

**Consumes**

- Defaults (built-in values).
- Optional global config file (platform-specific user config directory; path may be overridden by a documented env var).
- Optional workspace config file (current working directory).
- Optional explicit config file path (when user supplies it via CLI option; then global and workspace are not loaded).
- Environment variables (override file-based values).
- Command-line overrides for the current run.

**Produces**

- **Effective config** — A single merged configuration for the current context: root loop settings, prompt definitions (with optional per-prompt overrides), alias definitions. Used by run-loop, review, list, and show.
- **Effective config for display** — The same resolved config, optionally with provenance (which layer supplied each value), for `ralph show config`. The show config command does not pass CLI overrides; display uses file/env/prompt layers only.

**Called by**

- CLI: to resolve config before dispatching to run, review, list, or show.

**Single entrypoint (T1.7, O002/R007)**

- **Resolve(getenv, cwd, configPath, promptName)** — Resolve effective config for the current context. Returns (*Effective, ok, error). The returned Effective includes built-in aliases; user aliases override built-ins for the same name. Use this from run-loop, review, list, and show. When promptName is empty, root loop is returned; when promptName is set and the prompt exists, Effective.Loop includes that prompt’s overrides. When promptName is set but not found, returns (nil, false, nil).

## Implementation spec

### Layer order (lowest to highest priority)

1. **Defaults** — Built-in values so the tool works without a config file.
2. **Global config file** — User-level; optional; if missing, skipped.
3. **Workspace config file** — Project-level in cwd; optional; if missing, skipped.
4. **Explicit config file** — When the user specifies a config file path via the documented CLI option, only that file is used for file-based config; global and workspace are not loaded. The file must exist or the system reports an error.
5. **Environment variables** — Override file-based config. A documented env var can also control the global config file path.
6. **Prompt-level overrides** — In config files, each prompt can specify its own loop settings; those apply when running or listing that prompt and override root loop settings; env and CLI still override for that run.
7. **Command-line options** — Override all other layers for that run.

### Built-in defaults

When no config file is present or a setting is omitted from all layers, the following built-in values apply so the tool works without a config file (O002/R001, R002). Implementation: `internal/config/defaults.go`.

| Setting | Default |
|--------|--------|
| max_iterations | 10 |
| failure_threshold | 3 |
| timeout_seconds | 0 (no per-iteration timeout) |
| success_signal | `<promise>SUCCESS</promise>` |
| failure_signal | `<promise>FAILURE</promise>` |
| signal_precedence | `static` |
| preamble | (empty; no preamble injection) |
| streaming | true |
| log_level | `info` |

Built-in AI command aliases (e.g. `claude`, `kiro`, `copilot`, `cursor-agent`) are defined in the config package and merged with user-defined aliases; user aliases override built-ins for the same name.

### Config file structure (canonical schema)

Config files are YAML. The following structure is the authoritative shape implementers must support. Unknown keys may be ignored or rejected per implementation policy; the listed keys are required for product behavior.

**Root-level keys**

- **loop** (object, optional) — Root loop behavior. All keys below can appear here and can be overridden per prompt.
  - **max_iterations** (integer, optional) — Maximum iterations before exit. See [Built-in defaults](#built-in-defaults); 0 or absent = use default.
  - **failure_threshold** (integer, optional) — Consecutive failures before exit.
  - **timeout_seconds** (integer, optional) — Per-iteration timeout; 0 or absent = no timeout.
  - **success_signal** (string, optional) — Substring or pattern that indicates success in AI output.
  - **failure_signal** (string, optional) — Substring or pattern that indicates failure.
  - **signal_precedence** (string, optional) — e.g. `static` (first match wins) or `ai_interpreted` when both signals appear.
  - **preamble** (string or boolean, optional) — Optional preamble injection; or enable/disable.
  - **streaming** (boolean, optional) — Whether to show AI command output in the terminal (default: true). Used by both `ralph run` and `ralph review`. Env: `RALPH_LOOP_STREAMING`. CLI: `--no-stream` only (turns off for that run; no flag to turn on—streaming is the default).
  - **log_level** (string, optional) — Log verbosity (e.g. debug, info, warn, error).
  - **ai_cmd** (string, optional) — Direct AI command string (e.g. a full CLI invocation). When both **ai_cmd** and **ai_cmd_alias** are set, the direct command takes precedence over the alias.
  - **ai_cmd_alias** (string, optional) — AI command alias name; must be a name from **aliases** or a built-in alias (e.g. `claude`, `cursor-agent`).
- **prompts** (object, optional) — Map of prompt name to prompt definition.
  - Each entry: **path** or **content** (file path or inline); optional **loop** overrides (same keys as root loop, including **ai_cmd** and **ai_cmd_alias**).
  - **Prompt path resolution:** A relative **path** is resolved relative to the **directory containing the config file that defined that prompt** (the layer that supplied the prompt when layers are merged). It is not relative to the current working directory. Absolute paths remain absolute. This allows config files to reference prompt files next to them or in a stable location relative to the config (e.g. `./prompts/build.md` or `prompts/build.md` from the same directory as the config file).
- **aliases** (object, optional) — Map of alias name to AI command string (or alias definition with **command**).
  - Each entry: **command** (string) — The AI CLI command line (e.g. `claude --non-interactive`).

Exact key names and nesting may be refined in implementation; this document defines the minimal set for layer resolution and loop/review behavior. Validation: invalid or out-of-range values (e.g. negative max_iterations) produce a clear error or documented fallback.

### Environment variables

The following environment variables affect configuration. They are applied in the environment layer (after file-based config, before CLI flags). Only the variables listed here are supported; other `RALPH_*` names are ignored. Values are parsed when set; invalid values produce a clear error.

**Config file location (global config directory)**

| Variable | Effect |
|----------|--------|
| `RALPH_CONFIG_HOME` | Directory used to locate the global config file. The file path is `$RALPH_CONFIG_HOME/ralph-config.yml`. When unset, fallback is `$XDG_CONFIG_HOME/ralph/ralph-config.yml`, then `~/.config/ralph/ralph-config.yml`. Does not set the explicit config file for the current invocation (that is only via CLI `--config`). |

**Loop settings (overlay onto resolved config)**

| Variable | Config key / effect | Type / notes |
|----------|---------------------|--------------|
| `RALPH_LOOP_AI_CMD` | Direct AI command string | string |
| `RALPH_LOOP_AI_CMD_ALIAS` | AI command alias name | string |
| `RALPH_LOOP_ITERATION_MODE` | Iteration mode (e.g. `max-iterations`, `unlimited`) | string |
| `RALPH_LOOP_DEFAULT_MAX_ITERATIONS` | Max iterations | integer; invalid value → error |
| `RALPH_LOOP_FAILURE_THRESHOLD` | Consecutive failures before exit | integer; invalid value → error |
| `RALPH_LOOP_ITERATION_TIMEOUT` | Per-iteration timeout in seconds; 0 = no timeout | integer; invalid value → error |
| `RALPH_LOOP_LOG_LEVEL` | Log level (e.g. `debug`, `info`, `warn`, `error`) | string |
| `RALPH_LOOP_STREAMING` | Whether to stream AI output to terminal (config key: `loop.streaming`) | boolean: `true`/`1`/`yes`/`on` → true; `false`/`0`/`no`/`off`/empty → false; invalid → error |
| `RALPH_LOOP_PREAMBLE` | Enable/disable preamble injection | boolean (same parsing as above) |

When a variable is unset, it does not override; the value from a lower layer (file or defaults) is used. For example, when `RALPH_LOOP_STREAMING` is unset, the effective value comes from config or default (typically true for normal runs).

### Prompt path resolution

When a prompt is defined in a config file with a **path** (file-based prompt), relative paths are resolved against the **directory of the config file that defined that prompt**:

- **Single explicit config** (`--config path/to/ralph.yml`): the prompt path is relative to the directory containing that file (e.g. `path/to/`).
- **Global + workspace merge:** the prompt comes from the layer that won the merge (workspace overrides global for the same name). The path is relative to that layer’s config file directory (global file dir or workspace file dir).
- Absolute paths are left unchanged. Paths are resolved at merge time so the effective config exposes a resolved (absolute) path for file-based prompts.

This keeps prompt locations tied to the config that defines them and avoids dependence on the process current working directory.

### Resolution rules

- For each setting, the effective value is the one from the highest-priority layer that supplies that setting.
- **AI command and alias:** **ai_cmd** and **ai_cmd_alias** use the same layer order as other loop settings. The effective value is the highest-priority layer that supplies it (root or prompt loop, then env, then CLI flags). When both a direct command and an alias are supplied, the direct command overrides the alias.
- When explicit config path is supplied, only that file is read; global and workspace are not read. Missing explicit file → error.
- Missing global or workspace file → skip that layer; no error.
- Prompt-level loop overrides apply when running or listing that prompt; env and CLI overrides still apply for that run.

### Show config and provenance

`ralph show config` prints the effective (resolved) config. With `--provenance`, it also reports which layer supplied each loop value. For **ai_cmd** and **ai_cmd_alias**, the reported layer is one of: **default** (no layer set it), **global**, **workspace**, **explicit**, **env**, or **prompt** (when the value came from a prompt’s loop override).
