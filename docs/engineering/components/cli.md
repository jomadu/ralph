# CLI

**This document is the source of truth for all CLI commands and flags.** Implementers must use it to construct the whole CLI. Every command name, subcommand, flag (long and short form), argument, and behavioral rule is specified here so that nothing is left to guesswork. User-facing help and documentation must align with this spec.

---

## Responsibility

The CLI component parses command-line arguments and environment, resolves configuration (via the config component), and dispatches to the appropriate handler: run-loop for `ralph run`, review for `ralph review`, list/show for `ralph list` and `ralph show`, and a version printer for `ralph version`. It exposes help and non-interactive flags so scripts and CI can run without prompts. It does **not** implement install, uninstall, or upgrade as subcommands; those are documented procedures only (see O006, O011).

Implements the requirements assigned to this component in the [engineering README](../README.md).

---

## Top-level commands

The first argument after `ralph` is always one of: **run**, **review**, **list**, **show**, **version**. Nothing else is a valid top-level command. Unknown command → error to stderr, non-zero exit, suggest `--help`.

| Command | Purpose |
|--------|--------|
| `ralph run` | Run the iteration loop. Prompt via alias, file, or stdin. |
| `ralph review` | Review prompt (alias, file, stdin); produce a report directory (result.json, summary.md, original.md, revision.md, diff.md) and a suggested revision; optional apply. |
| `ralph list` | List prompts and/or aliases from resolved config. |
| `ralph show` | Show effective config or detail for a prompt/alias. |
| `ralph version` | Print version string and exit 0. |

**Help:** `ralph --help`, `ralph -h` — Print top-level help and exit.  
**Per-command help:** `ralph run --help`, `ralph review --help`, `ralph list --help`, `ralph show --help` — Print help for that command and exit.

---

## Global options

These options affect config resolution or process behavior and apply to **run**, **review**, **list**, and **show** (where relevant). They must be accepted and processed before dispatching to the subcommand handler.

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--config` | *(none)* | path | Use this file as the sole file-based config. Global and workspace config are not loaded. The file must exist or the CLI reports an error and exits non-zero (R005). Relative to cwd unless absolute. |
| `--help` | `-h` | — | Print help and exit. At top level: list commands. On a command: list that command's flags and usage. |
| `--version` | — | — | Print version and exit 0. Same effect as `ralph version`. |

**Note:** There is no short form for `--config` because `-c` is reserved for `--context` on `ralph run` (repeatable). Implementers must not add a short form for config that would conflict with that.

**Environment variable:** `RALPH_CONFIG_HOME` — If set, it is the **directory** used to locate the global config file. The global config file path is `$RALPH_CONFIG_HOME/ralph-config.yml` (filename is fixed). Fallback when unset: `$XDG_CONFIG_HOME/ralph/ralph-config.yml`, then `~/.config/ralph/ralph-config.yml`. This only affects *where* the global (user-level) config is read from; it does **not** set the explicit config file for the current invocation. For "use this one file and skip global/workspace," the user must pass `--config <path>`. When `--config <path>` is passed, that path is used and env does not override it. (The implementation does not provide an env var that means "use this path as the explicit config file"; that is CLI-only.)

---

## ralph run

**Purpose:** Run the iteration loop. Prompt is supplied once (alias, file path, or stdin) and buffered; the run-loop invokes the backend each iteration until a success or failure condition.

### Syntax

- `ralph run [flags]` — Prompt from **stdin** (stdin must not be a TTY, or behavior is defined: e.g. error “no prompt source”).
- `ralph run <alias> [flags]` — Prompt from the named prompt in resolved config (alias must exist).
- `ralph run --file <path> [flags]` or `ralph run -f <path> [flags]` — Prompt from file at `<path>` (file must exist and be readable).

**Exactly one** of: positional alias, `-f`/`--file`, or stdin. Combining alias with `-f`/`--file` is an error (clear message, non-zero exit). When no alias and no `-f` and stdin is a TTY, the CLI must error (no prompt source) and not start the loop.

### Flags (all optional)

Precedence: CLI flags override environment and config (config component layer order).

**Prompt source**

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--file` | `-f` | path | Read prompt from this file. Mutually exclusive with positional alias and stdin. File must exist. |

**Loop control**

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--max-iterations` | `-n` | int | Override max iterations for this run. Must be ≥ 1 or defined (e.g. 0 = use config default). |
| `--unlimited` | `-u` | — | Run until success signal or failure threshold; no iteration cap. Overrides max-iterations for this run. |
| `--failure-threshold` | — | int | Consecutive failures before exit. Override for this run. |
| `--iteration-timeout` | — | int | Per-iteration timeout in seconds. 0 = no timeout. Override for this run. |
| `--max-output-buffer` | — | int | Max output buffer in bytes for capturing AI stdout. Override for this run. |
| `--no-preamble` | — | — | Disable preamble injection for this run. |
| `--dry-run` | `-d` | — | Do not invoke the AI. Assemble prompt (with preamble if enabled) and print it; then exit 0. No backend invocation. |

**AI command**

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--ai-cmd` | — | string | Direct AI command string for this run (overrides config alias). |
| `--ai-cmd-alias` | — | string | AI command alias name from config for this run. Overrides config default. |

When `--ai-cmd` or `--ai-cmd-alias` are omitted, the value is taken from the environment overlay (e.g. `RALPH_LOOP_AI_CMD`, `RALPH_LOOP_AI_CMD_ALIAS`) if set; otherwise from resolved config (root loop or per-prompt loop). Thus flags override env and config; env overrides config.

**Signals**

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--signal-success` | — | string | Success signal string for this run (substring or pattern in AI output). |
| `--signal-failure` | — | string | Failure signal string for this run. |
| `--signal-precedence` | — | string | When both signals appear: `static` (e.g. first match wins) or `ai_interpreted` (one extra AI invocation to decide). Override for this run. |

**Context / preamble**

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--context` | `-c` | string | Inline context string injected into preamble (e.g. CONTEXT section). Repeatable. Value is literal text; not read as a file path. |

**Output and observability**

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--verbose` | `-v` | — | Increase verbosity: log level debug. |
| `--quiet` | `-q` | — | Minimal output: log level error-only; do not show AI command output. |
| `--log-level` | — | string | Log level: `debug`, `info`, `warn`, `error`. Overrides config and shortcuts when set. |
| `--no-stream` | — | — | Do not show AI command output in the terminal. Default is to show it; this flag turns it off for this run. |

**Configuration (global)**

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--config` | — | path | Explicit config file; only this file for file-based config (see Global options). |

### Run: error handling

- Missing or invalid AI command/alias (e.g. alias not in config, binary not on PATH) → clear error before loop starts; documented error exit code (see run-loop).
- Explicit `--config` path missing or unreadable → error, do not start run; non-zero exit.
- Invalid or out-of-range flag value (e.g. negative `--max-iterations`) → error and exit non-zero.
- Unknown flag → parser error and exit non-zero.

---

## ralph review

**Purpose:** Review a prompt (alias, file, or stdin). Produce a report directory (result.json, summary.md, original.md, revision.md, diff.md) and a suggested revision; optionally write the revision to a path with confirmation (or non-interactive flag).

### Syntax

- `ralph review [flags]` — Prompt from **stdin** (non-TTY or defined behavior).
- `ralph review <alias> [flags]` — Prompt from named prompt in resolved config.
- `ralph review --file <path> [flags]` or `ralph review -f <path> [flags]` — Prompt from file.

**Exactly one** of: positional alias, `-f`/`--file`, or stdin. Same mutual-exclusion and error rules as run.

### Flags (all optional except when apply + stdin)

**Prompt source**

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--file` | `-f` | path | Read prompt from this file. Mutually exclusive with alias and stdin. |

**Report and revision output**

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--report` | — | path | Report directory path. The path is interpolated into the review prompt so the AI knows where to create the five files (`result.json`, `summary.md`, `original.md`, `revision.md`, `diff.md`). If omitted, default is `./ralph-review/` (under current working directory). If the path exists and is a file (not a directory), error and exit 2. If the path does not exist, Ralph creates it as a directory before invoking the AI. Path must be writable. |
| `--prompt-output` | — | path | When using `--apply`, write the revision to this path. **Required** when prompt is from stdin and `--apply` is set; if missing in that case, error and exit 2. When prompt is from file/alias, may default to source file (with confirmation) or require path; behavior must be documented. |
| `--apply` | — | — | Request that the suggested revision be written to a file. In interactive mode, confirmation is required before overwriting (unless `--yes`). In non-interactive mode, use `--yes` to apply without confirmation or error with a clear message if confirmation would be needed. |
| `--yes` | `-y` | — | Non-interactive apply: do not prompt for confirmation; apply revision when `--apply` is set. If confirmation would be required and session is non-interactive and `--yes` is not set, exit 2 with clear message. |

**AI command**

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--ai-cmd` | — | string | Direct AI command for this review (overrides config and env). |
| `--ai-cmd-alias` | — | string | AI command alias name for this review (overrides config and env). |

When omitted, the AI command/alias is taken from the environment overlay if set; otherwise from resolved config (root or per-prompt loop). Flags override env and config; env overrides config.

**Output and observability**

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--verbose` | `-v` | — | Increase verbosity: log level debug. Same semantics as run. |
| `--quiet` | `-q` | — | Minimal output: log level error-only; do not show AI command output. Same semantics as run. |
| `--log-level` | — | string | Log level: `debug`, `info`, `warn`, `error`. Overrides shortcuts when set. |
| `--no-stream` | — | — | Do not show AI command output in the terminal. Default is to show it; same semantics as run. |

**Configuration (global)**

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--config` | — | path | Explicit config file (global). |

### Review: error handling

- Invalid or missing prompt source (alias unknown, file missing, stdin TTY with no input) → error before review; exit 2.
- Stdin + `--apply` and no `--prompt-output` → error, exit 2.
- Report directory unwritable, or path exists and is not a directory (e.g. existing file), or directory cannot be created → error, exit 2.
- Apply requested, confirmation required, non-interactive and no `--yes` → error, exit 2.
- Unknown flag or invalid value → parser error, non-zero exit.

---

## ralph list

**Purpose:** List prompts and/or AI command aliases from the resolved config. Uses the same config resolution as run (defaults, global, workspace, explicit file, env, CLI). Read-only; does not run the loop or modify config.

### Syntax

- `ralph list [flags]` — List **all**: both prompts and aliases from resolved config.
- `ralph list prompts [flags]` — List only prompts (names and optional display name, description, path).
- `ralph list aliases [flags]` — List only AI command aliases (names and optional expansion/description).

No positional arguments after the subcommand. Invalid subcommand (e.g. `ralph list foo`) → error and non-zero exit.

### Flags

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--config` | — | path | Explicit config file (global). |
| `--help` | `-h` | — | Print list command help and exit. |

### List: output

Output format (e.g. table, YAML, JSON) is implementation-defined but must include names for prompts and aliases; prompts may include display name, description, and path as defined in config. Empty list is valid (no error).

---

## ralph show

**Purpose:** Show detailed information: effective (resolved) config, or a single prompt or alias by name. Uses the same config resolution as run and list.

### Syntax

- `ralph show config [flags]` — Output the effective config for the current context (same resolution as run). Accepts only `--provenance`, `--prompt`, and global `--config`; no run-style flags. Optionally include provenance (which layer supplied each value).
- `ralph show prompt [name] [flags]` — Show detailed information for the prompt named `name`. **When `name` is omitted:** the implementation errors with "name required" and suggests `ralph show prompt <name>` or `ralph list prompts`; name is required.
- `ralph show alias [name] [flags]` — Show detailed information for the alias named `name`. **When `name` is omitted:** the implementation errors with "name required" and suggests `ralph show alias <name>` or `ralph list aliases`; name is required.
- `ralph show prompt-guide [flags]` — Output the full [Writing Ralph prompts](../../writing-ralph-prompts.md) guide **verbatim** (same content as docs/writing-ralph-prompts.md). Single source of truth: the doc is defined only in docs/; the CLI emits that content so users get the full guide from the command. Supports optional `--markdown` (output is already markdown; flag for consistency or saving/piping to a pager). No config resolution required; exits 0. Implements O008/R005.

**Required:** The first argument after `show` must be one of: **config**, **prompt**, **alias**, **prompt-guide**. Unknown object (e.g. `ralph show foo`) → error and non-zero exit.

### Flags

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--config` | — | path | Explicit config file (global). |
| `--provenance` | — | — | For `show config` only: show which layer supplied each loop value (default, global, workspace, explicit file, env, prompt override). The "cli" layer does not appear when viewing config because `show config` does not accept run-style flags; CLI overrides apply only when running. |
| `--prompt` | — | string | For `show config` only: show effective config for the named prompt (includes prompt-level loop overrides). |
| `--markdown` | — | — | For `show prompt-guide` only: output the full guide (already markdown; flag for saving or piping to a pager). |
| `--help` | `-h` | — | Print show command help and exit. |

### Show: error handling

- Missing object (e.g. `ralph show`) → error, suggest usage; non-zero exit.
- Unknown object type → error, non-zero exit.
- `show prompt <name>` or `show alias <name>` with unknown name → clear error, non-zero exit.
- Explicit config file missing when `--config` used → error, non-zero exit.

---

## ralph version

**Syntax:** `ralph version` or `ralph version [flags]`.

**Behavior:** Print the version string to stdout and exit 0. No arguments required. Flags (if any) are implementation-defined (e.g. `--short`); the only required behavior is that with no args or with `--help` it prints version and exits 0.

---

## Non-interactive behavior

When the process is not interactive (e.g. no TTY, or CI), the CLI must **not** block on user prompts. For **review**:

- If the user requests `--apply` and confirmation would be required (e.g. overwrite), the CLI must either: (a) proceed without prompting when `--yes` is set, or (b) exit 2 with a clear message that confirmation is required and suggest `--yes` for non-interactive use.
- Detection of non-interactive context (e.g. `!isatty(stdin)` or explicit env such as `CI=true`) is implementation-defined but must be documented; the presence of `--yes` must always suppress confirmation when `--apply` is set.

For **run**, no confirmation prompts are required; non-interactive behavior is satisfied by config and flags (O010).

---

## Environment variables (summary)

- **RALPH_CONFIG_HOME** — Directory for the global config file; actual file is `$RALPH_CONFIG_HOME/ralph-config.yml`. Does not set the explicit config file for the current run; use `--config <path>` for that.
- **RALPH_LOOP_*** — Env vars that override loop settings (e.g. `RALPH_LOOP_AI_CMD`, `RALPH_LOOP_STREAMING`, `RALPH_LOOP_DEFAULT_MAX_ITERATIONS`, `RALPH_LOOP_LOG_LEVEL`). Exact set and mapping are in the config component; must be documented so that full non-interactive config is possible (O010/R004).

---

## Error handling (CLI-level)

- **Unknown top-level command** — Error to stderr, non-zero exit (e.g. 1), suggest `ralph --help`.
- **Unknown subcommand** (e.g. `ralph list foo` when `foo` is not `prompts` or `aliases`) — Error, non-zero exit.
- **Unknown flag** — Rejected by parser; error and exit non-zero.
- **Missing required argument** — E.g. `ralph show` with no object; error and exit non-zero.
- **Mutually exclusive options** — E.g. alias + `-f` on run or review; error and exit non-zero.
- **Explicit config file missing** — Error (R005), do not proceed with run/list/show; non-zero exit.

Exit codes for **run** and **review** are defined in the run-loop and review components; the CLI propagates them. For parse and dispatch errors (unknown command, bad args), the CLI uses a consistent non-zero code (e.g. 1) that is documented.

---

## Interfaces

**Consumes:** Process arguments (argv), environment variables. Resolved configuration from the config component (after CLI invokes config resolution with cwd, `--config`, env).

**Produces:** Dispatched invocation to run-loop, review, list, or show with resolved context; version string to stdout for `ralph version`; help text to stdout for `--help`; exit code from dispatched command or from CLI on parse/usage error.

**Calls:** Config component (resolve configuration); run-loop (subcommand `run`); review (subcommand `review`); list/show handlers (subcommands `list`, `show`). CLI does not mutate config; it only reads.

---

## Implementation note

Help text (`ralph --help`, `ralph run --help`, etc.) and user-facing CLI documentation must align with this document. Implementers should treat this file as the single source from which flag names, short forms, types, and behavioral rules are derived so that the built CLI matches the spec without drift.
