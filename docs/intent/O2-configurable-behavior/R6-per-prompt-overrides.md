# R6: Per-Prompt Loop Setting Overrides

**Outcome:** O2 — Configurable Behavior

## Requirement

The system allows each prompt alias to override any loop configuration value independently. Prompt-level overrides take precedence over root-level loop config but are still overridden by environment variables and CLI flags. Each alias's overrides are isolated — one alias's settings do not affect another.

## Specification

Each prompt alias may include a `loop` section that overrides root-level loop configuration. Prompt-level overrides are a subset of the merge chain — they sit between file-level config and environment variables in precedence.

**Full precedence order for a given prompt alias:**

```
CLI flags                           (highest)
Environment variables (RALPH_*)
Prompt-level loop overrides         ← this requirement
Workspace root loop config
Global root loop config
Built-in defaults                   (lowest)
```

**Overridable fields:**

The prompt `loop` section accepts exactly the same keys as the root `loop` section:

| Key | Type | Effect |
|-----|------|--------|
| `iteration_mode` | string | Override iteration mode for this alias |
| `default_max_iterations` | int | Override max iterations for this alias |
| `failure_threshold` | int | Override failure threshold for this alias |
| `iteration_timeout` | int | Override per-iteration timeout for this alias |
| `max_output_buffer` | int | Override output buffer size for this alias |
| `ai_cmd` | string | Override direct AI command for this alias |
| `ai_cmd_alias` | string | Override AI command alias for this alias |
| `preamble` | bool | Override preamble injection for this alias |
| `signals.success` | string | Override success signal for this alias |
| `signals.failure` | string | Override failure signal for this alias |

Fields not present in the prompt's `loop` section inherit from the root config (which itself was resolved from global + workspace layers).

**Merge semantics:**

1. Resolve the root `loop` config by merging built-in defaults → global config → workspace config (per R1).
2. When `ralph run <alias>` is invoked, look up the prompt's `loop` section.
3. For each field present in the prompt's `loop` section, overlay it onto the resolved root config. Fields not present in the prompt's `loop` section retain the root-resolved value.
4. Apply environment variable overrides on top of the result.
5. Apply CLI flag overrides on top of the result.

The merge is field-level, not section-level. If a prompt sets `loop.failure_threshold: 5` and `loop.signals.success: "DONE"`, only those two fields are overridden — all other fields come from the root config.

**Isolation:**

Each prompt alias's overrides are independent. Configuring `prompts.build.loop.failure_threshold: 5` has no effect on `prompts.bootstrap` or any other alias. When Ralph runs a different alias, that alias's own overrides (or lack thereof) apply.

**Config file scoping:**

In default mode (no `--config` flag), prompt-level overrides can appear in both global and workspace config files. The same precedence applies: workspace-level prompt overrides take precedence over global-level prompt overrides for the same alias. Prompt definitions are merged by alias name across files — if global config defines `prompts.build` and workspace config also defines `prompts.build`, the workspace definition's fields overlay the global definition's fields (field-level merge, same as root `loop`).

When `--config <path>` is provided, the explicit file is the sole file-based config source (R5). Only prompt definitions in that file participate. No global or default workspace config is consulted.

**Prompt-only fields:**

The three prompt-specific fields (`path`, `name`, `description`) are not loop overrides. They belong to the prompt alias definition and do not participate in loop config merging. `path` is required (R3 validates this). `name` and `description` are optional and used only for display (R7).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Prompt sets `loop.failure_threshold: 5`; root is `loop.failure_threshold: 3` | Resolved value is 5 when running that alias |
| Prompt sets `loop.failure_threshold: 5`; CLI passes `--failure-threshold 10` | Resolved value is 10 — CLI overrides prompt-level |
| Prompt sets `loop.signals.success: "DONE"` but does not set `loop.signals.failure` | Success signal is `"DONE"`; failure signal inherits from root config (default: `<promise>FAILURE</promise>`) |
| Prompt defines no `loop` section at all | All loop values inherited from root config — the alias runs with the same settings as any unconfigured alias |
| Two prompts set different `ai_cmd_alias` values | Each alias uses its own value when executed. The values are independent. |
| Global config defines `prompts.build.loop.preamble: false`; workspace config does not define `prompts.build` at all | The global prompt definition applies: preamble is false for the `build` alias |
| Global config defines `prompts.build.path: ./a.md`; workspace defines `prompts.build.path: ./b.md` | Workspace wins: path is `./b.md` |
| Prompt sets `loop.ai_cmd: "my-cli"` and also sets `loop.ai_cmd_alias: "claude"` | Both are set; `ai_cmd` takes precedence over `ai_cmd_alias` at runtime (consistent with root-level behavior) |
| Environment variable `RALPH_LOOP_FAILURE_THRESHOLD=7` is set while prompt has `loop.failure_threshold: 5` | Resolved value is 7 — env vars override prompt-level settings |

### Examples

#### Prompt with custom signals and low iteration limit

**Input:**
Root config:
```yaml
loop:
  default_max_iterations: 10
  signals:
    success: "<promise>SUCCESS</promise>"
    failure: "<promise>FAILURE</promise>"
```

Prompt config:
```yaml
prompts:
  bootstrap:
    path: "./prompts/bootstrap.md"
    loop:
      default_max_iterations: 1
      preamble: false
```

User runs `ralph run bootstrap`.

**Expected output:**
Loop runs with `default_max_iterations: 1`, `preamble: false`. Signal strings are inherited from root: `<promise>SUCCESS</promise>` and `<promise>FAILURE</promise>`.

**Verification:**
- Loop executes at most 1 iteration
- No preamble is injected
- Signal scanning uses the default signal strings

#### CLI override on top of prompt override

**Input:**
Prompt config:
```yaml
prompts:
  build:
    path: "./prompts/build.md"
    loop:
      failure_threshold: 5
```

User runs `ralph run build --failure-threshold 2`.

**Expected output:**
Resolved `failure_threshold` is 2 (CLI flag), not 5 (prompt-level).

**Verification:**
- Debug log shows `loop.failure_threshold = 2 (source: cli)`
- The loop aborts after 2 consecutive failures, not 5

## Acceptance criteria

- [ ] A prompt alias can override any loop setting: iteration_mode, default_max_iterations, failure_threshold, iteration_timeout, max_output_buffer, ai_cmd, ai_cmd_alias, preamble, signals.success, signals.failure
- [ ] Prompt-level overrides take precedence over root loop config
- [ ] Environment variables and CLI flags still override prompt-level settings
- [ ] Unspecified prompt-level values inherit from the root loop section
- [ ] Each alias's overrides are independent — configuring one alias does not affect any other alias

## Dependencies

_None identified._
