# R3: Configuration Validation at Load Time

**Outcome:** O2 — Configurable Behavior

## Requirement

The system validates all resolved configuration values at load time and fails fast with clear, actionable error messages for invalid values. Validation runs after all layers are merged, so the user sees errors for the final resolved values, not intermediate ones.

## Specification

Validation runs once, after all config layers have been merged into a single resolved config struct. Ralph validates the final resolved values, not individual layers. This means a workspace config with `failure_threshold: 0` does not error if a CLI flag overrides it to `3` before validation runs.

**Validation sequence:**

1. Merge all layers (built-in defaults → global config → workspace config → per-prompt overrides → env vars → CLI flags)
2. Run unknown key warnings (R2)
3. Validate the resolved config
4. If any validation errors exist, print all errors and exit with code 1
5. If validation passes, proceed to prompt source validation (R4) and then loop execution

**Validation is split into two layers:**

### Schema validation

The configuration schema is defined as a JSON Schema. Schema validation covers structural correctness: types, allowed values, ranges, and required fields. Error messages are produced by the schema validator and are not hand-crafted by Ralph.

| Field | Type | Constraint |
|-------|------|------------|
| `loop.iteration_mode` | string | Enum: `"max-iterations"`, `"unlimited"` |
| `loop.default_max_iterations` | int | Minimum: 1 |
| `loop.failure_threshold` | int | Minimum: 1 |
| `loop.iteration_timeout` | int | Minimum: 0 (0 means no timeout), or absent |
| `loop.max_output_buffer` | int | Minimum: 1 |
| `loop.log_level` | string | Enum: `"debug"`, `"info"`, `"warn"`, `"error"` |
| `loop.signals.success` | string | MinLength: 1 |
| `loop.signals.failure` | string | MinLength: 1 |
| `loop.ai_cmd` | string | MinLength: 1 (when present) |
| `loop.ai_cmd_alias` | string | MinLength: 1 (when present) |
| `prompts.<name>.path` | string | Required, MinLength: 1 |

Schema validation errors are reported as the validator produces them. Ralph does not rewrite or wrap these messages.

### Semantic validation

Semantic validation covers constraints that require cross-referencing resolved state — things a JSON Schema cannot express. Ralph owns these error messages.

| Check | Error |
|-------|-------|
| `loop.ai_cmd_alias` references an alias name not present in the resolved `ai_cmd_aliases` map | `invalid ai_cmd_alias "<value>": no matching alias defined` |

**Per-prompt validation:**

Each prompt alias's `loop` overrides are validated with the same rules as the root `loop` section. Validation runs on the resolved values after prompt-level overrides have been merged with root-level config. Only fields explicitly present in the prompt's `loop` section are validated at the prompt level — inherited values were already validated at the root level.

When `ralph run <alias>` is invoked, validation includes the prompt-specific resolved values. When `ralph list prompts` is invoked, validation covers root config and all prompt definitions (paths non-empty), but does not resolve per-prompt loop overrides since no specific alias is being executed.

**Error reporting:**

- All validation errors — schema and semantic — are collected and reported together, not fail-on-first. The user sees every problem in one pass.
- Schema validation errors use whatever format the JSON Schema validator produces. Ralph does not reformat them.
- Semantic validation errors include the field name, the invalid value, and the provenance tag (R1) so the user knows which layer to fix. Format: `validation error: <message> (source: <provenance-tag>)`
- Errors are written to stderr.
- Ralph exits with code 1 after printing all errors. The loop does not start.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Multiple validation errors across different fields | All errors are reported in a single output; Ralph does not exit after the first error |
| `iteration_timeout` is 0 | Valid — means no timeout. Distinct from a negative value, which is invalid. |
| `iteration_timeout` is nil/absent | Valid — means no timeout. The built-in default is nil (no timeout). |
| `default_max_iterations` is 0 | Validation error — must be >= 1 |
| `signals.success` and `signals.failure` are set to the same string | Valid from a schema perspective — signal precedence (O1-R2) handles the overlap at runtime. No validation error. |
| `ai_cmd` and `ai_cmd_alias` are both set | Valid — `ai_cmd` takes precedence over `ai_cmd_alias` at runtime. No validation error. |
| `ai_cmd_alias` references a user-defined alias not yet in `ai_cmd_aliases` | Validation error — the alias must resolve at load time |
| A prompt's `loop.failure_threshold` is invalid but the root `loop.failure_threshold` is valid | Validation error for the prompt-level value, reported with the prompt's provenance |
| Environment variable `RALPH_LOOP_DEFAULT_MAX_ITERATIONS` is set to non-numeric text | String-to-int conversion (`strconv`) fails; treated as a validation error with a clear message |
| Config file has valid YAML syntax but a value of the wrong type (e.g., `failure_threshold: "high"`) | Struct-based YAML decoding fails on type mismatch; treated as a load-time error with the field name and file path |

### Examples

#### Multiple schema validation errors

**Input:**
Workspace config:
```yaml
loop:
  default_max_iterations: 0
  failure_threshold: -1
  log_level: verbose
  signals:
    success: ""
```

**Expected output:**
Four schema validation errors are reported (one per invalid field). The exact messages are produced by the JSON Schema validator. Ralph exits with code 1.

**Verification:**
- All four errors appear in stderr
- Each error identifies the offending field and value
- Ralph does not start the loop
- Exit code is 1

#### Invalid value overridden by valid CLI flag

**Input:**
Workspace config has `default_max_iterations: 0`. User runs `ralph run build -n 5`.

**Expected output:**
No validation error for `default_max_iterations` — the CLI flag (5) overrides the workspace value (0) before validation runs. The resolved value is 5, which is valid.

**Verification:**
- Ralph starts normally
- Debug log shows `default_max_iterations = 5 (source: cli)`

## Acceptance criteria

- [ ] Schema validation rejects values that violate the JSON Schema constraints (type, range, enum, required, minLength)
- [ ] Semantic validation rejects `ai_cmd_alias` values that do not resolve to a known alias
- [ ] All errors (schema and semantic) are collected and reported together before exit
- [ ] Validation errors prevent the loop from starting (exit code 1)
- [ ] Semantic validation errors identify the invalid value, the field name, and the config source (provenance)

## Dependencies

_None identified._
