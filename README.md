# Ralph

A dumb loop that pipes a prompt to an AI CLI, lets it work, and repeats.

Ralph is a loop runner, not a methodology. You bring the prompt. Ralph runs it in a fresh AI process per iteration, scans for completion signals, and stops when the task is done — or when it isn't going to be.

## Quick Start

```bash
# Define a prompt alias in ralph-config.yml
prompts:
  build:
    path: "./prompts/build.md"

# Run it
ralph run build

# Override iteration limit
ralph run build -n 20

# Run until success or failure threshold
ralph run build --unlimited

# Preview assembled prompt without executing
ralph run build --dry-run

# Run a one-off prompt from a file (no alias needed)
ralph run -f ./prompts/fix-tests.md

# Pipe a prompt via stdin
cat prompts/build.md | ralph run
```

## How It Works

Each iteration:

1. Read the prompt file mapped to the alias
2. Wrap it with a preamble (iteration count, optional context)
3. Pipe the assembled prompt to the AI CLI's stdin
4. Capture output, scan for success/failure signals
5. On success signal → exit 0
6. On failure signal → increment consecutive failure counter
7. On consecutive failure threshold → abort (exit 1)
8. On max iterations reached → exit 2
9. Otherwise → next iteration

Fresh process per iteration. No conversation history carried between runs. State continuity comes from the filesystem — the AI reads and writes files, and the next iteration's AI sees those changes.

## Configuration

Five layers, highest precedence first:

| Layer | Source |
|-------|--------|
| CLI flags | `--max-iterations`, `--ai-cmd`, etc. |
| Environment variables | `RALPH_LOOP_*` |
| Workspace config | `./ralph-config.yml` |
| Global config | `~/.config/ralph/ralph-config.yml` |
| Built-in defaults | Compiled into the binary |

### Example `ralph-config.yml`

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
  cursor-agent: "cursor-agent -p -f --stream-partial-output --output-format stream-json"

prompts:
  build:
    path: "./prompts/build.md"
    name: "Build"
    description: "Run the main build loop"
    loop:
      default_max_iterations: 10
      failure_threshold: 5

  bootstrap:
    path: "./prompts/bootstrap.md"
    name: "Bootstrap"
    description: "One-shot project setup"
    loop:
      default_max_iterations: 1
      preamble: false
```

Each prompt alias maps a name to a file and can override any loop setting.

## Signals

Ralph scans AI CLI output for configurable signal strings to determine iteration outcome.

| Signal | Default | Meaning |
|--------|---------|---------|
| Success | `<promise>SUCCESS</promise>` | Task complete, stop looping |
| Failure | `<promise>FAILURE</promise>` | Blocked, increment failure counter |

Your prompt tells the AI what to emit. Ralph's signal config tells the scanner what to look for — use whatever strings your prompt expects.

If both signals appear in the same output, failure wins.

## CLI

```
ralph run <alias> [flags]          Run a prompt by config alias
ralph run -f <path> [flags]        Run a prompt from a file path
cat prompt.md | ralph run [flags]  Run a prompt piped via stdin
ralph list prompts                 List available prompt aliases
ralph list aliases                 List available AI command aliases
ralph version                      Show version info

Flags:
  -f, --file path                 Read prompt from file (no alias required)
  -n, --max-iterations int        Override max iterations
  -u, --unlimited                 Run until signal or failure threshold
      --failure-threshold int     Consecutive failures before abort
      --iteration-timeout int     Per-iteration timeout in seconds (0 = no timeout)
      --max-output-buffer int     Max output buffer in bytes
      --no-preamble               Disable preamble injection
  -d, --dry-run                   Validate and show assembled prompt
      --ai-cmd string             Direct AI command string
      --ai-cmd-alias string       AI command alias
      --signal-success string     Success signal string
      --signal-failure string     Failure signal string
  -c, --context string            Inject context into preamble (repeatable)
  -v, --verbose                   Stream AI output to terminal
  -q, --quiet                     Suppress non-error output
      --log-level string          debug, info, warn, error
      --config path               Alternate config file path
```

The prompt is read once at loop start and reused for every iteration. When no alias or `-f` flag is provided, Ralph reads from stdin.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success signal received |
| 1 | Failure threshold reached or abort |
| 2 | Max iterations exhausted |
| 130 | Interrupted (SIGINT/SIGTERM) |

## License

GPL-3.0. See [LICENSE](LICENSE).
