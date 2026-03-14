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
- **Effective config for display** — The same resolved config, optionally with provenance (which layer supplied each value), for `ralph show config`.

**Called by**

- CLI: to resolve config before dispatching to run, review, list, or show.

## Implementation spec

### Layer order (lowest to highest priority)

1. **Defaults** — Built-in values so the tool works without a config file.
2. **Global config file** — User-level; optional; if missing, skipped.
3. **Workspace config file** — Project-level in cwd; optional; if missing, skipped.
4. **Explicit config file** — When the user specifies a config file path via the documented CLI option, only that file is used for file-based config; global and workspace are not loaded. The file must exist or the system reports an error.
5. **Environment variables** — Override file-based config. A documented env var can also control the global config file path.
6. **Prompt-level overrides** — In config files, each prompt can specify its own loop settings; those apply when running or listing that prompt and override root loop settings; env and CLI still override for that run.
7. **Command-line options** — Override all other layers for that run.

### Config file structure (canonical schema)

Config files are YAML. The following structure is the authoritative shape implementers must support. Unknown keys may be ignored or rejected per implementation policy; the listed keys are required for product behavior.

**Root-level keys**

- **loop** (object, optional) — Root loop behavior. All keys below can appear here and can be overridden per prompt.
  - **max_iterations** (integer, optional) — Maximum iterations before exit. Default documented; 0 or absent may mean "use product default."
  - **failure_threshold** (integer, optional) — Consecutive failures before exit.
  - **timeout_seconds** (integer, optional) — Per-iteration timeout; 0 or absent = no timeout.
  - **success_signal** (string, optional) — Substring or pattern that indicates success in AI output.
  - **failure_signal** (string, optional) — Substring or pattern that indicates failure.
  - **signal_precedence** (string, optional) — e.g. `static` (first match wins) or `ai_interpreted` when both signals appear.
  - **preamble** (string or boolean, optional) — Optional preamble injection; or enable/disable.
  - **streaming** (boolean, optional) — Whether to stream AI output to the terminal. Env: `RALPH_LOOP_STREAMING`. CLI: `--stream` / `--no-stream`.
  - **log_level** (string, optional) — Log verbosity (e.g. debug, info, warn, error).
- **prompts** (object, optional) — Map of prompt name to prompt definition.
  - Each entry: **path** or **content** (file path or inline); optional **loop** overrides (same keys as root loop).
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

### Resolution rules

- For each setting, the effective value is the one from the highest-priority layer that supplies that setting.
- When explicit config path is supplied, only that file is read; global and workspace are not read. Missing explicit file → error.
- Missing global or workspace file → skip that layer; no error.
- Prompt-level loop overrides apply when running or listing that prompt; env and CLI overrides still apply for that run.
