# Configuration

**Intent:** [O2 — Configurable behavior](../intent/O2-configurable-behavior/README.md); [R1 — Configuration provenance](../intent/O2-configurable-behavior/R1-configuration-provenance.md); [R8 — Environment variables](../intent/O2-configurable-behavior/R8-environment-variable-reference.md); [R9 — CLI interface](../intent/O2-configurable-behavior/R9-cli-interface-reference.md)

Ralph’s loop behavior (iterations, timeouts, signals, AI backend) is controlled by configuration files, environment variables, and CLI flags. You can change how the loop runs without editing your prompt file.

## Where configuration lives

| Layer | Location |
|-------|----------|
| **Workspace** | `./ralph-config.yml` in the directory where you run `ralph` |
| **Global** | `~/.config/ralph/ralph-config.yml` (or `$RALPH_CONFIG_HOME/ralph-config.yml` if set; otherwise `$XDG_CONFIG_HOME/ralph/ralph-config.yml`) |
| **Explicit file** | `ralph run build --config /path/to/ralph-config.yml` — this file replaces both workspace and global; only this file is loaded |

If a config file is missing, Ralph skips it (no error). So you can have only a workspace config, only a global config, or both; workspace overrides global for keys you set in both.

## Precedence (who wins)

From **highest** to **lowest**:

1. **CLI flags** — e.g. `ralph run build -n 20 --failure-threshold 5`
2. **Environment variables** — e.g. `RALPH_LOOP_ITERATION_TIMEOUT=60`
3. **Prompt-level overrides** — `loop` section under a prompt alias in the config file
4. **Workspace config** — `./ralph-config.yml`
5. **Global config** — `~/.config/ralph/ralph-config.yml`
6. **Built-in defaults**

So: to change the iteration limit for one run, use `-n 20`. To change it for all runs in a project, set `default_max_iterations` in `./ralph-config.yml`. To change it for one prompt alias only, add a `loop` section under that alias.

## Key options (task-oriented)

### Change how many iterations run

- **Config (root):** `loop.default_max_iterations: 10`
- **Config (per prompt):** under `prompts.<alias>.loop`, set `default_max_iterations: 10`
- **CLI:** `ralph run build -n 20` or `--max-iterations 20`
- **Env:** `RALPH_LOOP_DEFAULT_MAX_ITERATIONS=20`

### Change when the loop gives up (consecutive failures)

- **Config:** `loop.failure_threshold: 5`
- **CLI:** `ralph run build --failure-threshold 5`
- **Env:** `RALPH_LOOP_FAILURE_THRESHOLD=5`

### Change success/failure signal strings

- **Config:** `loop.signals.success: "SUCCESS"`, `loop.signals.failure: "FAILURE"` (defaults)
- **CLI:** `ralph run build --signal-success "DONE" --signal-failure "ABORT"`

These are the strings Ralph looks for in the AI CLI output to decide whether the task succeeded or failed.

### Change per-iteration timeout

- **Config:** `loop.iteration_timeout: 120` (seconds; `0` = no timeout)
- **CLI:** `ralph run build --iteration-timeout 60`
- **Env:** `RALPH_LOOP_ITERATION_TIMEOUT=60`

### Choose or override the AI backend

- **Config:** `loop.ai_cmd_alias: cursor-agent` (use a built-in or user-defined alias), or `loop.ai_cmd: "/path/to/command"` (direct command)
- **CLI:** `ralph run build --ai-cmd-alias claude` or `ralph run build --ai-cmd "/path/to/agent"`
- **Env:** `RALPH_LOOP_AI_CMD_ALIAS` or `RALPH_LOOP_AI_CMD`

See [AI backends and aliases](ai-backends-and-aliases.md) for built-in and user-defined aliases and direct command; see the [Cursor Agent workaround](cursor-agent-workaround.md) for that backend specifically.

### Logging and output

- **Config:** `loop.log_level: debug|info|warn|error`; `loop.show_ai_output: true` to stream AI output to the terminal
- **CLI:** `ralph run build --log-level debug`, `-v` / `--verbose`, `-q` / `--quiet`
- **Env:** `RALPH_LOOP_LOG_LEVEL`, `RALPH_LOOP_SHOW_AI_OUTPUT`

## Example: workspace config

```yaml
# ./ralph-config.yml
loop:
  default_max_iterations: 15
  failure_threshold: 3
  iteration_timeout: 300
  signals:
    success: "SUCCESS"
    failure: "FAILURE"
  ai_cmd_alias: cursor-agent

prompts:
  build:
    path: prompts/build.md
    description: "Build and test"
  bootstrap:
    path: prompts/bootstrap.md
    description: "One-shot bootstrap"
    loop:
      default_max_iterations: 1
      failure_threshold: 1
```

Here, `build` uses the root loop settings (15 iterations, threshold 3, cursor-agent). `bootstrap` overrides to 1 iteration and threshold 1 for that alias only.

## Full reference

- **Environment variables:** All supported `RALPH_*` variables and types are listed in [O2 R8 — Environment variable reference](../intent/O2-configurable-behavior/R8-environment-variable-reference.md).
- **CLI commands and flags:** Full CLI surface is in [O2 R9 — CLI interface reference](../intent/O2-configurable-behavior/R9-cli-interface-reference.md).
- **Provenance:** With `--log-level debug` or `--dry-run`, you can see which layer (cli, env, workspace, global, default) supplied each resolved value.
