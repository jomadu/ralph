# R8: Environment Variable Reference

**Outcome:** O2 — Configurable Behavior

## Requirement

The system supports a defined set of environment variables with the `RALPH_*` prefix that override or supply configuration values. This requirement is the authoritative reference for which variables exist and what config key each maps to; it does not define precedence or merge semantics (those are in R1, R5, R6, and the per-feature requirements).

## Specification

Ralph reads the following environment variables from the process environment when resolving configuration. Each variable maps to a configuration key. Unset variables are ignored (the value comes from a lower-precedence layer per R1). The precedence order is: CLI flags > environment variables > workspace config > global config > built-in defaults (R1).

**General rules:**

- Variable names use the `RALPH_` prefix. Only the variables listed below are supported; other `RALPH_*` names are ignored (no warning).
- Values are parsed according to the type of the target config key (see R3 for validation). Invalid values produce validation errors at load time.
- Boolean env vars: when set, common truthy values (e.g. `true`, `1`, `yes`) resolve to true; `false`, `0`, `no`, or empty resolve to false. Implementation may accept additional forms.

### Supported variables

| Variable | Maps to | Type / notes |
|----------|---------|----------------|
| `RALPH_CONFIG_HOME` | Global config directory | string. Directory path used to locate `ralph-config.yml` (see R5). Resolved before config load. If set, Ralph looks for `$RALPH_CONFIG_HOME/ralph-config.yml`. Does not map to a key in the merged config struct; it changes where the global config file is read from. |
| `RALPH_LOOP_AI_CMD` | `loop.ai_cmd` | string. Direct AI command. When set, overrides or supplies the command string (see O3/R6). |
| `RALPH_LOOP_AI_CMD_ALIAS` | `loop.ai_cmd_alias` | string. Name of an AI command alias. |
| `RALPH_LOOP_ITERATION_MODE` | `loop.iteration_mode` | string. Enum: `max-iterations` \| `unlimited` (see O1/R4). |
| `RALPH_LOOP_DEFAULT_MAX_ITERATIONS` | `loop.default_max_iterations` | int. Positive integer; minimum 1. |
| `RALPH_LOOP_FAILURE_THRESHOLD` | `loop.failure_threshold` | int. Positive integer; minimum 1. |
| `RALPH_LOOP_ITERATION_TIMEOUT` | `loop.iteration_timeout` | int. Seconds; 0 means no timeout. |
| `RALPH_LOOP_LOG_LEVEL` | `loop.log_level` | string. Enum: `debug` \| `info` \| `warn` \| `error` (see O4/R5). |
| `RALPH_LOOP_SHOW_AI_OUTPUT` | `loop.show_ai_output` | bool. When true, stream AI CLI output to the terminal (see O4/R3). When **unset**, default is **true** (per O4/R3). |

No other `RALPH_*` variables are defined. Future variables may be added by extending this table in a later requirement.

### Edge cases

| Condition | Expected Behavior |
|----------|-------------------|
| Variable unset → default true | For `RALPH_LOOP_SHOW_AI_OUTPUT`: when unset, default is true; AI output is streamed (per O4/R3). |
| Variable set to empty string | Treated as set; value overlay applied (R1). Validation (R3) may reject empty for required/min-length fields. |
| Variable set to invalid value (e.g. non-numeric for an int key) | Config validation fails at load time (R3); Ralph exits with code 1 and a clear message. |
| Variable not in this list but with `RALPH_` prefix | Ignored; no warning. No overlay. |
| Same key set by both env var and CLI flag | CLI wins (R1). |

### Examples

#### Override timeout via env

**Input:**
`RALPH_LOOP_ITERATION_TIMEOUT=60`, no config file. User runs `ralph run build`.

**Expected output:**
Resolved `loop.iteration_timeout` is 60 (from env). Per-iteration timeout of 60 seconds is applied (O1/R3).

**Verification:**
- Debug log (if enabled) shows `loop.iteration_timeout = 60 (source: env)`.

#### Boolean env var

**Input:**
`RALPH_LOOP_SHOW_AI_OUTPUT=true`. User runs `ralph run build` without `-v`.

**Expected output:**
Resolved `loop.show_ai_output` is true. AI output is streamed to the terminal (O4/R3).

**Verification:**
- AI CLI stdout/stderr appear in the terminal during the run.

## Acceptance criteria

- [ ] Every variable in the table above is read from the environment and overlays the corresponding config key when set
- [ ] Unset or unsupported variables do not alter config
- [ ] Invalid values for a variable produce validation errors (R3) and exit 1
- [ ] This document is the single authoritative list of supported RALPH_* environment variables

## Dependencies

- R1 (configuration provenance) — env layer and precedence.
- R3 (config validation) — validation applies to env-supplied values.
- R5 (silent skip) — RALPH_CONFIG_HOME affects global config path resolution.
- Per-feature requirements (O1/R4, O3/R6, O4/R3, O4/R5) — semantics of each key.
- R6 (per-prompt overrides) — prompt-level loop overrides are overridden by env vars.
