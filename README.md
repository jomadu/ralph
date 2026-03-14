# Ralph

A dumb loop that pipes a prompt to an AI CLI, lets it work, and repeats.

Ralph is a loop runner — an implementation of the [Ralph Wiggum technique](https://ghuntley.com/ralph/) by [Geoffrey Huntley](https://github.com/ghuntley). The idea: put an AI coding agent in a `while` loop, give it a prompt, and let it build your project one task at a time with a fresh context window each iteration. State continuity comes from the filesystem, not conversation history.

For the full technique and playbook, see [How to Ralph Wiggum](https://github.com/ghuntley/how-to-ralph-wiggum).

Ralph (this tool) adds structure on top of the raw bash loop: configurable signals for success/failure detection, per-prompt loop settings, multi-CLI support, and a config layer so you can tune the loop without editing scripts.

## Quick Start

Create a prompt file:

`prompts/build.md`
```markdown
Study AGENTS.md for build commands and project context.
Pick the most important task.
Write tests first. Run them — they should fail.
Implement fully. No placeholders, no stubs, no partial implementations.
Run tests again — they should pass.
If all work is done, output: <promise>SUCCESS</promise>
If blocked, output: <promise>FAILURE</promise>
```

Wire it up in `ralph-config.yml`:

```yaml
loop:
  ai_cmd_alias: claude

prompts:
  build:
    path: "./prompts/build.md"
```

Run it:

```bash
ralph run build              # run the loop
ralph run build -n 20        # cap at 20 iterations
ralph run build --unlimited  # run until signal or failure threshold
ralph run build --dry-run    # preview assembled prompt
ralph run -f ./prompts/fix.md  # one-off prompt from file
cat prompt.md | ralph run    # pipe via stdin
```

Each iteration: read prompt, pipe to AI CLI, scan output for signals, stop or loop.

## Install and Uninstall

**Prerequisites:** `curl`. Installs from release artifacts only. Supported: Linux, macOS, Windows (amd64, arm64).

Install latest:

```bash
curl -fsSL https://raw.githubusercontent.com/maxdunn/ralph/main/scripts/install.sh | sh
```

Install specific version:

```bash
curl -fsSL https://raw.githubusercontent.com/maxdunn/ralph/main/scripts/install.sh | sh -s -- 1.0.0
```

Verify:

```bash
ralph version
```

Uninstall (removes binary + install state, keeps config):

```bash
curl -fsSL https://raw.githubusercontent.com/maxdunn/ralph/main/scripts/uninstall.sh | sh
```

Installs to `~/bin` by default. Ensure the install directory is on your PATH.

## How It Works

Each iteration:

1. Read the prompt file mapped to the alias
2. Wrap it with a preamble (iteration count, optional context)
3. Pipe the assembled prompt to the AI CLI's stdin
4. Capture output, scan for success/failure signals
5. On success signal: exit 0
6. On failure signal: increment consecutive failure counter
7. On consecutive failure threshold: abort (exit 1)
8. On max iterations reached: exit 2
9. Otherwise → next iteration

Fresh process per iteration. No conversation history carried between runs. State continuity comes from the filesystem — the AI reads and writes files, and the next iteration's AI sees those changes.

## Configuration

Five layers, highest precedence first:

| Layer | Source |
|-------|--------|
| CLI flags | `--max-iterations`, `--ai-cmd`, etc. |
| Environment variables | `RALPH_LOOP_*` |
| Workspace config | `./ralph-config.yml` |
| Global config | `$RALPH_CONFIG_HOME/ralph-config.yml`, `$XDG_CONFIG_HOME/ralph/ralph-config.yml`, or `~/.config/ralph/ralph-config.yml` |
| Built-in defaults | Compiled into the binary |

### `ralph-config.yml`

| Key | Type | Description |
|-----|------|-------------|
| `loop.ai_cmd_alias` | string | AI command alias to use (e.g. `claude`) |
| `loop.default_max_iterations` | integer | Max iterations before exit |
| `loop.failure_threshold` | integer | Consecutive failures before abort |
| `loop.iteration_timeout` | integer | Per-iteration timeout in seconds (0 = no timeout) |
| `loop.show_ai_output` | boolean | Stream AI output to terminal |
| `loop.preamble` | boolean | Enable/disable preamble injection |
| `loop.log_level` | string | Log verbosity: `debug`, `info`, `warn`, `error` |
| `loop.signals.success` | string | Signal string indicating task success |
| `loop.signals.failure` | string | Signal string indicating task failure |
| `ai_cmd_aliases.<name>` | string | AI CLI command string for a named alias |
| `prompts.<name>.path` | string | Path to the prompt file |
| `prompts.<name>.name` | string | Display name |
| `prompts.<name>.description` | string | Description of the prompt |
| `prompts.<name>.loop` | object | Per-prompt loop overrides (same keys as `loop.*`) |

Example:

```yaml
loop:
  default_max_iterations: 5
  failure_threshold: 3
  iteration_timeout: 300
  show_ai_output: false
  ai_cmd_alias: claude
  signals:
    success: "<promise>SUCCESS</promise>"
    failure: "<promise>FAILURE</promise>"

ai_cmd_aliases:
  claude: "claude -p --dangerously-skip-permissions"
  kiro: "kiro-cli chat --no-interactive --trust-all-tools"
  copilot: "copilot --yolo"
  cursor-agent: "agent -p --force --output-format stream-json --stream-partial-output"

prompts:
  build:
    path: "./prompts/build.md"
    name: "Build"
    description: "Run the main build loop"
    loop:
      default_max_iterations: 10
      failure_threshold: 5

```

Each prompt alias maps a name to a file and can override any loop setting.

### Built-in Aliases

Ralph ships with built-in AI command aliases. Use them via `loop.ai_cmd_alias` in config or `--ai-cmd-alias` on the CLI. User-defined aliases in config override a built-in with the same name.

| Alias | Command |
|-------|---------|
| `claude` | `claude -p --dangerously-skip-permissions` |
| `kiro` | `kiro-cli chat --no-interactive --trust-all-tools` |
| `copilot` | `copilot --yolo` |
| `cursor-agent` | `agent -p --force --output-format stream-json --stream-partial-output` |

### Environment Variables

| Variable | Effect |
|----------|--------|
| `RALPH_CONFIG_HOME` | Directory for global config file lookup |
| `RALPH_LOOP_AI_CMD` | Direct AI command string |
| `RALPH_LOOP_AI_CMD_ALIAS` | AI command alias name |
| `RALPH_LOOP_DEFAULT_MAX_ITERATIONS` | Max iterations |
| `RALPH_LOOP_FAILURE_THRESHOLD` | Consecutive failures before exit |
| `RALPH_LOOP_ITERATION_TIMEOUT` | Per-iteration timeout in seconds (0 = no timeout) |
| `RALPH_LOOP_STREAMING` | Stream AI output to terminal |
| `RALPH_LOOP_PREAMBLE` | Enable/disable preamble injection |
| `RALPH_LOOP_LOG_LEVEL` | Log level (`debug`, `info`, `warn`, `error`) |

## Signals

Ralph scans AI CLI output for configurable signal strings to determine iteration outcome.

| Signal | Default | Meaning |
|--------|---------|---------|
| Success | `<promise>SUCCESS</promise>` | Task complete, stop looping |
| Failure | `<promise>FAILURE</promise>` | Blocked, increment failure counter |

Your prompt tells the AI what to emit. Ralph's signal config tells the scanner what to look for — use whatever strings your prompt expects.

If both signals appear in the same output, failure wins.

## CLI

### `ralph run`

Run the iteration loop. The prompt is read once at loop start and reused for every iteration.

```
ralph run <alias> [flags]          Prompt from config alias
ralph run -f <path> [flags]        Prompt from file
cat prompt.md | ralph run [flags]  Prompt from stdin
```

Exactly one prompt source (alias, `--file`, or stdin). When no alias or `-f` flag is provided and stdin is not a TTY, Ralph reads from stdin.

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--file` | `-f` | path | Read prompt from file |
| `--max-iterations` | `-n` | int | Override max iterations |
| `--unlimited` | `-u` | — | Run until signal or failure threshold |
| `--failure-threshold` | — | int | Consecutive failures before abort |
| `--iteration-timeout` | — | int | Per-iteration timeout in seconds (0 = no timeout) |
| `--max-output-buffer` | — | int | Max output buffer in bytes |
| `--no-preamble` | — | — | Disable preamble injection |
| `--dry-run` | `-d` | — | Assemble and print prompt, then exit |
| `--ai-cmd` | — | string | Direct AI command string |
| `--ai-cmd-alias` | — | string | AI command alias from config |
| `--signal-success` | — | string | Success signal string |
| `--signal-failure` | — | string | Failure signal string |
| `--context` | `-c` | string | Inject context into preamble (repeatable) |
| `--verbose` | `-v` | — | Stream AI output to terminal |
| `--quiet` | `-q` | — | Suppress non-error output |
| `--log-level` | — | string | `debug`, `info`, `warn`, `error` |
| `--config` | — | path | Explicit config file (skips global/workspace) |

### `ralph review`

Review a prompt and produce a report with a suggested revision.

```
ralph review <alias> [flags]
ralph review -f <path> [flags]
cat prompt.md | ralph review [flags]
```

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--file` | `-f` | path | Read prompt from file |
| `--report` | — | path | Write report to this path |
| `--prompt-output` | — | path | Write suggested revision to this path |
| `--apply` | — | — | Write the suggested revision to a file |
| `--yes` | `-y` | — | Skip confirmation when applying |
| `--quiet` | `-q` | — | Suppress non-error output |
| `--log-level` | — | string | `debug`, `info`, `warn`, `error` |
| `--config` | — | path | Explicit config file (skips global/workspace) |

### `ralph list`

List prompts and/or AI command aliases from resolved config.

```
ralph list                         List all prompts and aliases
ralph list prompts                 List prompt aliases
ralph list aliases                 List AI command aliases
```

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--config` | — | path | Explicit config file (skips global/workspace) |

### `ralph show`

Show effective config or detail for a prompt or alias.

```
ralph show config                  Show resolved configuration
ralph show prompt <name>           Show prompt detail
ralph show alias <name>            Show alias detail
```

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--provenance` | — | — | Include which layer supplied each value |
| `--config` | — | path | Explicit config file (skips global/workspace) |

### `ralph version`

Print version string and exit.

```
ralph version
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success signal received |
| 1 | Failure threshold reached or abort |
| 2 | Max iterations exhausted |
| 130 | Interrupted (SIGINT/SIGTERM) |

## Cursor Agent Wrapper

The built-in `cursor-agent` alias invokes the raw `agent` command with JSON streaming output. This works, but the JSON is noisy for signal scanning and terminal readability. Cursor's docs include a [wrapper script](https://cursor.com/docs/cli/headless#real-time-progress-tracking) that parses the stream-json output into clean text — use that as your alias so Ralph can scan signals cleanly.

Set up a custom alias in your config pointing to the wrapper:

```yaml
ai_cmd_aliases:
  cursor: "./scripts/cursor-wrapper.sh"
```

Then run with it:

```bash
ralph run build --ai-cmd-alias cursor
```

## Acknowledgments

Ralph is named after and inspired by the [Ralph Wiggum technique](https://ghuntley.com/ralph/) originated by [Geoffrey Huntley](https://github.com/ghuntley) — the idea of putting an AI coding agent in a bash `while` loop and letting it build software autonomously, one task per iteration, with a fresh context window each time. For the full technique, playbook, and prompt templates, see [How to Ralph Wiggum](https://github.com/ghuntley/how-to-ralph-wiggum).

## License

GPL-3.0. See [LICENSE](LICENSE).
