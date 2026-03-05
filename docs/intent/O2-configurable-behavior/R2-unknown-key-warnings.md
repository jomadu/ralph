# R2: Unknown Key Warnings

**Outcome:** O2 — Configurable Behavior

## Requirement

The system warns the user when a config file contains keys that are not part of the known schema, without preventing startup. This catches typos and forward-compatibility issues while remaining non-blocking.

## Specification

After parsing a YAML config file, Ralph compares every key in the parsed document against the known schema. Any key that does not correspond to a known field produces a warning. This applies recursively — unknown keys are detected at every nesting level within known sections.

**Known schema structure:**

```
loop:
  iteration_mode
  default_max_iterations
  failure_threshold
  iteration_timeout
  max_output_buffer
  log_level
  show_ai_output
  ai_cmd
  ai_cmd_alias
  preamble
  signals:
    success
    failure

ai_cmd_aliases:
  <any key>          # ai_cmd_aliases is a free-form map; keys are not validated

prompts:
  <any key>:         # prompt alias names are free-form
    path
    name
    description
    loop:            # same schema as root loop
      (same keys as loop above)
```

**Detection mechanism:**

The config schema is defined as typed Go structs with `yaml` struct tags. The `yaml.v3` decoder's `KnownFields(true)` setting rejects any YAML key that does not map to a struct field, producing an error that identifies the unknown key and its line number. Ralph uses this to detect unknown keys without a custom tree-walking algorithm.

Because `KnownFields` causes the decode to fail rather than warn, Ralph decodes each config file twice:

1. **Strict decode** — decode with `KnownFields(true)`. If this succeeds, there are no unknown keys. If it fails, collect the error(s) as warnings.
2. **Permissive decode** — decode without `KnownFields` to obtain the actual config values, silently discarding unknown keys.

Free-form maps are handled by their Go types: `ai_cmd_aliases` is typed as `map[string]string` and `prompts` is typed as `map[string]PromptConfig`, so any alias name is accepted at that map level. Within each `PromptConfig` struct, only `path`, `name`, `description`, and `loop` are defined as struct fields, so unknown keys within a prompt alias are caught. Within `prompts.<name>.loop`, the struct mirrors the root `loop` struct, so the same fields are recognized.

**Warning format:**

```
warning: unknown config key "<dotted.key.path>" in <file-path>
```

Examples:
```
warning: unknown config key "loop.retries" in ./ralph-config.yml
warning: unknown config key "prompts.build.priority" in ~/.config/ralph/ralph-config.yml
```

**Behavior rules:**

- Warnings are emitted at `warn` log level, so they appear at the default log level (`info`).
- Warnings are emitted once per unknown key per file. If the same unknown key appears in both global and workspace config, two warnings are emitted (one per file).
- Unknown keys do not prevent Ralph from starting. The unknown key and its value are silently discarded after the warning.
- Warnings are emitted during config loading, before validation (R3). A file can produce both unknown key warnings and validation errors.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Unknown key nested inside a known section (e.g., `loop.retries`) | Warning emitted: `unknown config key "loop.retries"` |
| Unknown top-level section (e.g., `plugins:`) | Warning emitted: `unknown config key "plugins"` |
| Unknown key inside `ai_cmd_aliases` (e.g., `ai_cmd_aliases.my-tool: "..."`) | No warning — `ai_cmd_aliases` is a free-form map |
| Unknown key inside a prompt alias's own fields (e.g., `prompts.build.priority: 1`) | Warning emitted: `unknown config key "prompts.build.priority"` |
| Unknown key inside a prompt alias's loop section (e.g., `prompts.build.loop.retries`) | Warning emitted: `unknown config key "prompts.build.loop.retries"` |
| Same unknown key in both global and workspace config | Two separate warnings, one per file, each identifying its file path |
| A YAML file with only unknown keys and no valid keys | All keys produce warnings; Ralph starts with built-in defaults |
| Key is valid but at the wrong nesting level (e.g., `failure_threshold` at root instead of under `loop`) | Warning emitted — keys are only recognized at their correct position in the schema |

### Examples

#### Typo in config key

**Input:**
Workspace config `./ralph-config.yml`:
```yaml
loop:
  max_iteratoins: 10
  failure_threshold: 3
```

**Expected output:**
```
warning: unknown config key "loop.max_iteratoins" in ./ralph-config.yml
```
Ralph starts normally. `default_max_iterations` uses its built-in default (5) since the typo'd key is not recognized.

**Verification:**
- Warning appears in stderr
- Ralph does not exit with an error
- `default_max_iterations` resolves to 5 (built-in default), not 10

#### Future config key from a newer version

**Input:**
User upgrades their config file but downgrades Ralph. Config contains:
```yaml
loop:
  failure_threshold: 3
  retry_delay: 5
```

**Expected output:**
```
warning: unknown config key "loop.retry_delay" in ./ralph-config.yml
```
Ralph starts normally, ignoring the unknown key.

**Verification:**
- Warning appears in stderr
- Ralph does not crash or refuse to start

## Acceptance criteria

- [ ] Unrecognized top-level keys in a config file produce a warning at load time
- [ ] Unrecognized nested keys within known sections also produce warnings
- [ ] Each warning identifies the key name and the config file path where it was found
- [ ] The presence of unknown keys does not prevent Ralph from starting or running

## Dependencies

_None identified._
