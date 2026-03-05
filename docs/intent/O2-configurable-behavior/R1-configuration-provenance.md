# R1: Configuration Provenance Tracking

**Outcome:** O2 — Configurable Behavior

## Requirement

The system tracks the source layer of each resolved configuration value so the user can determine where the active value came from. Every resolved value knows whether it was set by a built-in default, global config file, workspace config file, environment variable, or CLI flag.

## Specification

Every configuration value resolved by Ralph carries a provenance tag identifying which layer set it. Provenance is metadata attached to each value during the merge process — it is not a separate data structure maintained alongside the config.

**Provenance layers** (in precedence order, highest first):

| Layer | Tag | Source |
|-------|-----|--------|
| CLI flag | `cli` | Command-line flag parsed by cobra |
| Environment variable | `env` | `RALPH_*` variable read from the process environment |
| Explicit config file | `file` | Path from `--config` (replaces both global and workspace layers; see R5) |
| Workspace config | `workspace` | `./ralph-config.yml` |
| Global config | `global` | `~/.config/ralph/ralph-config.yml` (resolved via `RALPH_CONFIG_HOME` > `$XDG_CONFIG_HOME/ralph/` > `~/.config/ralph/`) |
| Built-in default | `default` | Compiled into the binary |

**Merge algorithm:**

1. Start with a config struct where every field is populated from built-in defaults, each tagged `default`.
2. If `--config <path>` was provided: parse that file as the sole file-based layer, tagging each value `file`. Skip steps 3 and 4.
3. If the global config file exists, parse it and overlay its values. Each value set by this file is tagged `global`.
4. If the workspace config file exists, parse it and overlay its values. Each value set by this file is tagged `workspace`.
5. For each supported `RALPH_*` environment variable, if the variable is set in the process environment, parse its value and overlay. Tag each as `env`.
6. For each CLI flag explicitly provided by the user, overlay its value. Tag each as `cli`.

At each overlay step, only keys explicitly present in that layer are applied. A key absent from a layer does not reset the value from a lower layer. The provenance tag always reflects the highest-precedence layer that set the value.

**Exposure:**

Provenance is surfaced through two channels:

- **Debug logging:** At `debug` log level, after config resolution completes, Ralph logs every resolved key with its final value and provenance tag. Format: `config: <key> = <value> (source: <tag>)`. For sensitive or long values, the value may be truncated but the provenance tag is always shown.
- **Dry-run output:** When `--dry-run` is active, the resolved configuration section includes provenance for each value displayed.

Provenance is internal metadata. It is not serialized to disk, not exposed via a public API, and not included in non-debug log output.

**Provenance for per-prompt overrides:**

When a prompt alias includes loop overrides (R6), the prompt-level values participate in the merge between file-based config and environment variables. A prompt-level override is tagged with the same layer as its containing config file (`global`, `workspace`, or `file`), with a qualifier: e.g., `workspace:prompt(build)` or `file:prompt(build)`. This distinguishes a root-level `loop.failure_threshold: 3` (tagged `workspace`) from a prompt-level `prompts.build.loop.failure_threshold: 5` (tagged `workspace:prompt(build)`).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| A key is set in both global and workspace config files | Workspace wins; provenance tag is `workspace` |
| A key is set only by built-in default (no config file, env var, or flag sets it) | Provenance tag is `default` |
| Environment variable is set to an empty string | The variable is considered set; its value (empty string) is overlaid and tagged `env`. Validation (R3) may then reject the value. |
| `--config` provides an explicit config file | Values from that file are tagged `file`. Global and workspace layers are not loaded (R5). The file path is visible in debug logging. |
| A prompt-level override and a root-level value both exist in the same workspace config | The prompt-level override applies when running that alias (tagged `workspace:prompt(<alias>)`); the root-level value applies for other aliases (tagged `workspace`) |
| CLI flag overrides a prompt-level override | CLI wins; provenance tag is `cli` |

### Examples

#### Debug log output after config resolution

**Input:**
Global config sets `loop.failure_threshold: 3`. Workspace config sets `loop.failure_threshold: 5`. User passes `--failure-threshold 10`.

**Expected output:**
Debug log includes:
```
config: loop.failure_threshold = 10 (source: cli)
```

**Verification:**
- Run with `--log-level debug` and confirm the line appears in stderr
- The resolved value is 10, not 3 or 5

#### Dry-run provenance display

**Input:**
No config files exist. User runs `ralph run build --dry-run --max-iterations 8`.

**Expected output:**
Dry-run output shows `default_max_iterations = 8 (source: cli)` and all other values show `(source: default)`.

**Verification:**
- Every value in the dry-run config section has a provenance annotation
- Only `default_max_iterations` shows `cli`; everything else shows `default`

## Acceptance criteria

- [ ] Each resolved config value carries its provenance: built-in default, global config, workspace config, environment variable, or CLI flag
- [ ] Provenance information is available via debug-level logging
- [ ] When multiple layers set the same key, the highest-precedence layer wins and the provenance reflects that winning layer
- [ ] Precedence order is: CLI flags > environment variables > workspace config > global config > built-in defaults

## Dependencies

_None identified._
