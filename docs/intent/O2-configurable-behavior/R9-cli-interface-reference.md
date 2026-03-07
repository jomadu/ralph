# R9: CLI Interface Reference

**Outcome:** O2 — Configurable Behavior

## Requirement

The system exposes a command-line interface that is the primary way users invoke Ralph. This requirement is the authoritative reference for top-level commands, subcommands, and flags. Semantics of each are specified in the requirement that owns the behavior (e.g. O1/R4 for `--max-iterations`); this requirement enumerates the full CLI surface so that documentation and help text have a single source of truth.

## Specification

### Top-level commands

| Command | Description | Specified in |
|---------|-------------|--------------|
| `ralph run` | Run the loop with a prompt (alias, file path, or stdin). | O1 (R9, R8, R4, etc.), O3, O4 |
| `ralph list` | List configured resources. Requires a subcommand. | O2/R7 |
| `ralph version` | Print version string and exit. | O4/R7 |
| `ralph --help`, `ralph -h` | Print help and exit. | (standard) |

### `ralph run` — prompt input

- `ralph run <alias>` — Run using the prompt alias `<alias>` (config must define `prompts.<alias>.path`). O1/R9.
- `ralph run --file <path>`, `ralph run -f <path>` — Run using the prompt file at `<path>` directly. O1/R9.
- `ralph run` with stdin not a TTY — Run using prompt content from stdin (read once, buffered). O1/R9.

Exactly one of: positional alias, `--file`/`-f`, or piped stdin. Combining alias with `--file` is an error.

### `ralph run` — flags

All flags are optional. Precedence: CLI flags override environment variables and config (R1, O2/R8).

**Loop control**

| Flag | Short | Type | Effect | Specified in |
|------|-------|------|--------|--------------|
| `--max-iterations` | `-n` | int | Override `loop.default_max_iterations` | O1/R4 |
| `--unlimited` | `-u` | — | Set `loop.iteration_mode` to `unlimited` | O1/R4 |
| `--failure-threshold` | — | int | Override `loop.failure_threshold` | O1/R5 |
| `--iteration-timeout` | — | int | Override `loop.iteration_timeout` (seconds; 0 = no timeout) | O1/R3 |
| `--max-output-buffer` | — | int | Override `loop.max_output_buffer` (bytes) | O1/R6 |
| `--preamble` | — | — | Enable preamble injection | O1/R8 |
| `--no-preamble` | — | — | Disable preamble injection | O1/R8 |
| `--dry-run` | `-d` | — | Validate config and print assembled prompt; do not execute loop | O4/R4 |

**AI command**

| Flag | Short | Type | Effect | Specified in |
|------|-------|------|--------|--------------|
| `--ai-cmd` | — | string | Override/set direct AI command string | O3/R6 |
| `--ai-cmd-alias` | — | string | Override/set AI command alias name | O3/R6 |

**Signals**

| Flag | Short | Type | Effect | Specified in |
|------|-------|------|--------|--------------|
| `--signal-success` | — | string | Override `loop.signals.success` | O1/R2, O2/R6 |
| `--signal-failure` | — | string | Override `loop.signals.failure` | O1/R2, O2/R6 |

**Context**

| Flag | Short | Type | Effect | Specified in |
|------|-------|------|--------|--------------|
| `--context` | `-c` | string | Inline context string; repeatable. Injected into preamble CONTEXT section. Not interpreted as a file path. | O1/R8 |

**Output control**

| Flag | Short | Type | Effect | Specified in |
|------|-------|------|--------|--------------|
| `--verbose` | `-v` | — | Verbose: log level debug, stream AI output unless `--no-ai-cmd-output` | O4/R3, O4/R5 |
| `--quiet` | `-q` | — | Quiet: log level error, do not stream AI command output | O4/R5, O4/R3 |
| `--log-level` | — | string | Set `loop.log_level` (debug, info, warn, error) | O4/R5 |
| `--no-ai-cmd-output` | — | — | Set show AI command output to false | O4/R3 |

**Configuration**

| Flag | Short | Type | Effect | Specified in |
|------|-------|------|--------|--------------|
| `--config` | — | path | Use this file as the sole file-based config (replaces global and workspace) | O2/R1, R5 |

### `ralph list` subcommands

| Subcommand | Description | Specified in |
|------------|-------------|--------------|
| `ralph list prompts` | List prompt aliases (YAML to stdout) | O2/R7 |
| `ralph list aliases` | List AI command aliases (YAML to stdout) | O2/R7 |

`ralph list` with no subcommand prints usage/help listing these subcommands.

### Global / help

- `ralph version` — Print version to stdout, exit 0. O4/R7.
- `ralph --help`, `ralph -h` — Print help and exit. Standard CLI behavior; help output must list `run`, `list`, and `version` so users can discover them.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Unknown top-level command | Error message and non-zero exit (e.g. 1); suggest --help. |
| `ralph run` with no alias, no -f, stdin is TTY | Error: no prompt source (O1/R9). |
| `ralph run build -f ./p.md` | Error: alias and --file mutually exclusive (O1/R9). |
| Unknown flag (e.g. `--unknown`) | Rejected by CLI parser; error or help. |
| `--context` value that looks like a path | Treated as literal inline string; no file is read (O1/R8). |

## Acceptance criteria

- [ ] All commands and flags in this document are implemented and discoverable via help
- [ ] No flag or command is part of the public CLI surface without being listed here (or in an amendment to this requirement)
- [ ] Help output for `ralph run` and `ralph list` is consistent with this reference

## Dependencies

- Individual requirements (O1/R3, R4, R5, R6, R8, R9; O2/R1, R5, R6, R7; O3/R6; O4/R3, R4, R5, R7) define the behavior of each flag and command; this requirement only enumerates the interface.
