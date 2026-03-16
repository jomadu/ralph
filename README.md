# Ralph

A dumb loop that pipes a prompt to an AI CLI, lets it work, and repeats. You bring the prompt; Ralph runs it until a success or failure signal appears—or a limit is reached—so multi-step work can reach verified completion without manual read–judge–re-run cycles.

## What is Ralph and why use it?

Ralph is a **loop runner** for AI-driven tasks, not a methodology. You supply a prompt and choose an AI CLI (e.g. Claude, Cursor). Ralph invokes it in a fresh process per iteration, scans the output for configurable success or failure signals, and stops when the task is done or when it isn’t going to be.

Without a loop runner, every iteration is manual: read the output, judge, re-invoke. Ralph closes that gap: you define what “done” looks like (a signal in the prompt), and Ralph handles re-invocation until that signal appears or a limit is reached. State lives on the filesystem—the AI reads and writes files, and the next iteration sees those changes. No conversation history is carried between runs.

## Path to first run

After [install](#install-and-uninstall), get to a first successful run as follows.

**1. Create a prompt file**

Example: create `prompts/build.md` from the repo root:

```bash
mkdir -p prompts
cat > prompts/build.md << 'EOF'
1. Study the project rules and workflow (including build/test commands if the project defines them).
2. Study the project's specifications; use product intent or acceptance criteria when needed for the task.
3. Get ready work from the issue tracker, pick one task, and claim it.
4. Complete the task. For code: use test-driven development and implement fully—no partial implementations.
5. When done: if no work remains in the queue, output `<promise>SUCCESS</promise>`; if blocked, output `<promise>FAILURE</promise>`; otherwise output nothing (loop continues).
EOF
```

**2. Define config so Ralph uses this prompt**

Create `ralph-config.yml` in the repo root (or in `~/.config/ralph/`) with a prompt alias:

```bash
cat > ralph-config.yml << 'EOF'
prompts:
  build:
    path: "./prompts/build.md"
EOF
```

**3. Choose your AI backend**

Ralph requires a resolved AI command; if none is set, it exits with code 2 before starting the loop. Either:

- **Use a built-in alias:** `claude`, `kiro`, `copilot`, or `cursor-agent`. Example: `ralph run build --ai-cmd-alias claude`, or set the alias in config or env. If your agent outputs structured or noisy data (e.g. JSON, tool calls), use a wrapper that sends progress to stderr and assistant text to stdout so Ralph can detect signals; see [Agent wrapper pattern](docs/agent-wrapper-pattern.md) (Cursor is one example).
- **Use your own command:** `--ai-cmd "your-command ..."` or `RALPH_LOOP_AI_CMD`, or define a custom alias under `aliases:` in config and use `--ai-cmd-alias` or `RALPH_LOOP_AI_CMD_ALIAS`.

**4. Run**

```bash
export RALPH_LOOP_AI_CMD_ALIAS=cursor-agent   # or claude, kiro, copilot

ralph run build --dry-run   # preview assembled prompt, no AI invocation
ralph run build -n 5        # run with max 5 iterations
ralph run build --unlimited # run until success signal or failure threshold
```

Prompt can also come from a file or stdin: `ralph run -f ./prompt.md` or `cat prompt.md | ralph run`. Use `ralph list` to see prompts and aliases from your config; use `ralph run --help` and `ralph --help` for full options.

## Install and Uninstall

Ralph can be installed via the provided script so the `ralph` binary is on your PATH. Install and uninstall are **documented procedures**, not subcommands (there is no `ralph install` or `ralph uninstall`).

**Prerequisites:** `curl`. The install script installs only from release artifacts (no build from source). Supported: Linux, macOS, Windows (amd64, arm64); script runs on macOS/Linux or Windows (e.g. Git Bash).

**Install**

1. **Latest release** (from repo or one-line):
   ```bash
   ./scripts/install.sh
   # or
   curl -fsSL https://raw.githubusercontent.com/maxdunn/ralph/main/scripts/install.sh | sh
   ```
   **Specific version** (e.g. `1.0.0` or `v1.0.0`):
   ```bash
   ./scripts/install.sh 1.0.0
   ./scripts/install.sh v1.0.0 --dir /usr/local/bin
   ```
   The script installs to `~/bin` by default and records the install location for uninstall.

2. Optional: use a different directory with `RALPH_INSTALL_DIR` or `--dir`:
   ```bash
   ./scripts/install.sh --dir /usr/local/bin
   ```
   If the directory is not writable (e.g. `/usr/local/bin`), run with `sudo` or choose a user directory like `~/bin`.

3. Ensure the install directory is on your PATH (e.g. add to your shell profile).

4. Verify in a new terminal:
   ```bash
   ralph version
   ```
   You should see version output and exit 0.

**Uninstall**

From anywhere (the script reads the install location from `~/.config/ralph/install-state`):

```bash
./scripts/uninstall.sh
```

This removes the `ralph` binary and the install state file. User config (e.g. `~/.config/ralph/ralph-config.yml`) is **not** removed. Install does not modify PATH or symlinks, so uninstall leaves no broken references.

**Upgrade:** Reinstall over the existing binary (e.g. run `install.sh` with the desired version) or use your package manager. Backward compatibility and any migration for breaking changes are described in [release notes](docs/release-notes.md).

## How it works

Before the loop starts, Ralph resolves the AI command from config or `--ai-cmd` / `--ai-cmd-alias`. If it cannot be resolved (e.g. not on PATH), Ralph exits with code 2 and does not start the loop.

Each iteration:

1. Use the prompt from the chosen source (alias → file, file path, or stdin)—loaded **once** at loop start and buffered for all iterations.
2. Optionally wrap it with a preamble (e.g. iteration count, context).
3. Pipe the assembled prompt to the AI CLI’s stdin.
4. Capture stdout and scan for configured success and failure signals.
5. **Success signal** → report completion, exit 0.
6. **Failure signal** → increment consecutive-failure count; if count ≥ failure threshold, exit 4; otherwise next iteration.
7. **Max iterations reached** → exit 3.
8. **No signal** (e.g. crash, timeout) → treated as failure; same threshold rules.
9. **Interrupt (e.g. Ctrl+C)** → exit 130.

Fresh process per iteration. No conversation history between runs. State continuity is via the filesystem—the AI reads and writes files, and the next iteration sees those changes.

## Configuration

Configuration is resolved from layers (lowest to highest priority):

1. **Defaults** — built-in values
2. **Global config file** — user-level
3. **Workspace config file** — project-level in cwd
4. **Explicit config file** — when you pass `--config` (only this file is used; global and workspace are not loaded)
5. **Environment variables** — e.g. `RALPH_LOOP_*`
6. **Prompt-level overrides** — per-prompt `loop` settings in config
7. **CLI flags** — override all of the above for that run

- **Global config:** `$RALPH_CONFIG_HOME/ralph-config.yml` if set; else `$XDG_CONFIG_HOME/ralph/ralph-config.yml` or `~/.config/ralph/ralph-config.yml`.
- **Workspace config:** `./ralph-config.yml` in the current working directory.
- **Explicit file:** `--config <path>` uses only that file (global and workspace are not loaded). The file must exist or Ralph errors.

Built-in defaults include `max_iterations: 10`, `failure_threshold: 3`, `success_signal: "<promise>SUCCESS</promise>"`, `failure_signal: "<promise>FAILURE</promise>"`, `signal_precedence: static`, `streaming: true`, `log_level: info`. Loop settings can be overridden with `RALPH_LOOP_*` environment variables (see [config spec](docs/engineering/components/config.md)).

### Example `ralph-config.yml`

```yaml
loop:
  max_iterations: 5
  failure_threshold: 3
  timeout_seconds: 300
  streaming: true
  success_signal: "<promise>SUCCESS</promise>"
  failure_signal: "<promise>FAILURE</promise>"
  signal_precedence: static

aliases:
  # Custom alias: Claude with a specific model (built-ins like claude, kiro, copilot, cursor-agent already exist)
  claude-sonnet: "claude -p --model claude-sonnet-4 --dangerously-skip-permissions"

prompts:
  build:
    path: "./prompts/build.md"
    display_name: "Build"
    description: "Run the main build loop"
    loop:
      max_iterations: 10
      failure_threshold: 5
```

Each prompt can override loop settings under its `loop` key. The AI command for a run is chosen by config or `--ai-cmd` / `--ai-cmd-alias`; if none is set, Ralph exits 2 with a clear error.

## Signals

Ralph scans AI CLI output for configurable success and failure signal strings.

| Signal   | Default                    | Meaning                                      |
|----------|----------------------------|----------------------------------------------|
| Success  | `<promise>SUCCESS</promise>` | Task complete; stop looping, exit 0.        |
| Failure  | `<promise>FAILURE</promise>` | Blocked; increment failure count or exit 4. |

Your prompt tells the AI what to emit. Ralph’s config (or flags) tell the scanner what to look for.

With `signal_precedence: static` (default), if **both** signals appear in the same output, **success wins**. Set `signal_precedence: ai_interpreted` (or `--signal-precedence ai_interpreted`) to have Ralph ask the AI once to interpret the outcome when both appear.

## CLI

**Global option:** `--config <path>` — use this file as the sole file-based config (global and workspace are not loaded). Applies to run, review, list, and show.

Full spec (all commands and flags): [docs/engineering/components/cli.md](docs/engineering/components/cli.md).

### ralph run

**Purpose:** Run the iteration loop. Prompt is supplied once (alias, file path, or stdin) and buffered; the run-loop invokes the AI each iteration until a success or failure condition (or a limit).

**Usage:** Exactly one prompt source per run.

- `ralph run <alias> [flags]` — Prompt from the named prompt in resolved config (alias must exist).
- `ralph run --file <path> [flags]` or `ralph run -f <path> [flags]` — Prompt from file at `<path>` (file must exist).
- `ralph run [flags]` — Prompt from **stdin** (e.g. `cat prompt.md | ralph run`). Stdin must not be a TTY or Ralph errors with no prompt source.

**Flags**

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--file` | `-f` | path | Read prompt from this file. Mutually exclusive with alias and stdin. |
| `--max-iterations` | `-n` | int | Override max iterations for this run. |
| `--unlimited` | `-u` | — | Run until success signal or failure threshold; no iteration cap. |
| `--failure-threshold` | — | int | Consecutive failures before exit. Override for this run. |
| `--iteration-timeout` | — | int | Per-iteration timeout in seconds. 0 = no timeout. |
| `--no-preamble` | — | — | Disable preamble injection for this run. |
| `--dry-run` | `-d` | — | Assemble prompt and print it; do not invoke the AI. Exit 0. |
| `--ai-cmd` | — | string | Direct AI command string for this run. |
| `--ai-cmd-alias` | — | string | AI command alias name from config for this run. |
| `--signal-success` | — | string | Success signal string for this run. |
| `--signal-failure` | — | string | Failure signal string for this run. |
| `--signal-precedence` | — | string | `static` or `ai_interpreted` when both signals appear. |
| `--context` | `-c` | string | Inline context injected into preamble. Repeatable. |
| `--verbose` | `-v` | — | Log level debug. |
| `--quiet` | `-q` | — | Minimal output; do not show AI command output. |
| `--log-level` | — | string | `debug`, `info`, `warn`, `error`. |
| `--no-stream` | — | — | Do not show AI command output in the terminal (default is to show it). |

**Exit codes**

| Code | Meaning |
|------|--------|
| 0 | Success signal detected; loop completed. |
| 2 | Error before loop (invalid or missing AI command, invalid config, or prompt source error). |
| 3 | Max iterations reached without success. |
| 4 | Failure threshold reached (consecutive failures). |
| 130 | Interrupted (e.g. SIGINT/Ctrl+C). |

### ralph review

**Purpose:** Review a prompt (alias, file, or stdin). Produce a report directory with five files (result.json, summary.md, original.md, revision.md, diff.md) and a suggested revision; optionally write the revision to a path with confirmation (or non-interactive flag).

**Usage:** Exactly one prompt source per run.

- `ralph review <alias> [flags]` — Prompt from named prompt in resolved config.
- `ralph review --file <path> [flags]` or `ralph review -f <path> [flags]` — Prompt from file.
- `ralph review [flags]` — Prompt from **stdin**. When using `--apply`, you must pass `--prompt-output <path>` or Ralph exits 2.

**Flags**

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--file` | `-f` | path | Read prompt from this file. Mutually exclusive with alias and stdin. |
| `--report` | — | path | Report directory path (default `./ralph-review/`). AI creates the five files here. Path must be writable; if it does not exist, Ralph creates it. |
| `--prompt-output` | — | path | When using `--apply`, write the revision to this path. **Required** when prompt is from stdin and `--apply` is set. |
| `--apply` | — | — | Write the suggested revision to a file. Confirmation required before overwriting unless `--yes`. |
| `--yes` | `-y` | — | Non-interactive apply: do not prompt for confirmation. If confirmation would be required and `--yes` is not set, exit 2. |
| `--verbose` | `-v` | — | Log level debug. |
| `--quiet` | `-q` | — | Minimal output; do not show AI command output. |
| `--log-level` | — | string | `debug`, `info`, `warn`, `error`. |
| `--no-stream` | — | — | Do not show AI command output in the terminal. |

**Exit codes**

| Code | Meaning |
|------|--------|
| 0 | Review completed; report directory written; no prompt errors (result.json indicates OK). |
| 1 | Review completed; report directory written; prompt has one or more errors (result.json indicates errors). |
| 2 | Review or apply did not complete (invalid prompt source, report write failure, stdin + apply without `--prompt-output`, confirmation required without `--yes` in non-interactive mode, or internal error). |

### ralph list

**Purpose:** List prompts and/or AI command aliases from the resolved config. Read-only; does not run the loop or modify config.

**Usage**

- `ralph list [flags]` — List both prompts and aliases from resolved config.
- `ralph list prompts [flags]` — List only prompts (names and optional display name, description, path).
- `ralph list aliases [flags]` — List only AI command aliases (names and optional expansion/description).

**Flags**

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--config` | — | path | Explicit config file (global). |
| `--help` | `-h` | — | Print list command help and exit. |

**Exit codes:** 0 on success. Parse or usage errors (e.g. invalid subcommand) use a consistent non-zero code (e.g. 1).

### ralph show

**Purpose:** Show effective (resolved) config, or detailed information for a single prompt or alias by name, or the full prompt-writing guide.

**Usage**

- `ralph show config [flags]` — Output the effective config for the current context. Optionally `--provenance` (which layer supplied each value) and `--prompt <name>` (effective config for that prompt).
- `ralph show prompt <name> [flags]` — Show detail for the prompt named `name`. Name is required; omit → error.
- `ralph show alias <name> [flags]` — Show detail for the alias named `name`. Name is required; omit → error.
- `ralph show prompt-guide [flags]` — Output the full [Writing Ralph prompts](docs/writing-ralph-prompts.md) guide verbatim.

**Flags**

| Flag | Short | Type | Effect |
|------|-------|------|--------|
| `--config` | — | path | Explicit config file (global). |
| `--provenance` | — | — | For `show config` only: show which layer supplied each value. |
| `--prompt` | — | string | For `show config` only: show effective config for the named prompt. |
| `--markdown` | — | — | For `show prompt-guide` only: output the full guide (for saving or piping to a pager). |
| `--help` | `-h` | — | Print show command help and exit. |

**Exit codes:** 0 on success. Missing object (e.g. `ralph show` with no config/prompt/alias/prompt-guide), unknown object type, or unknown prompt/alias name → error and non-zero exit.

### ralph version

**Purpose:** Print the version string and exit.

**Usage:** `ralph version` or `ralph version [flags]`. No arguments required.

**Flags:** Implementation-defined (e.g. `--short`, `--help`). With no args or `--help`, prints version and exits 0.

**Exit codes:** 0 on success.

---

**Stable contract for scripts and CI:** [docs/exit-codes.md](docs/exit-codes.md).

## Where to look

| Topic | Canonical spec |
|-------|----------------|
| **CLI (commands and flags)** | [docs/engineering/components/cli.md](docs/engineering/components/cli.md) |
| **Config and environment** | [docs/engineering/components/config.md](docs/engineering/components/config.md) |
| **Agent wrapper (progress vs. signals)** | [docs/agent-wrapper-pattern.md](docs/agent-wrapper-pattern.md) |
| **Exit codes** | [docs/exit-codes.md](docs/exit-codes.md) |

## Troubleshooting

- **Prompt not found / unknown alias** — Prompt source is exactly one of: alias (from config), `-f`/`--file`, or stdin. Use `ralph list prompts` to see defined prompts; check `--config` and config file locations if your alias isn’t found.
- **Config file not found** — With `--config <path>`, only that file is used and it must exist. Without it, global and workspace configs are optional (missing files are skipped).
- **Wrong or unexpected exit code** — See [docs/exit-codes.md](docs/exit-codes.md). Common causes for exit 2: missing AI command, stdin + `--apply` without `--prompt-output`, report path not writable or path is an existing file, or confirmation required in non-interactive mode without `--yes`.
- **AI command not found** — The AI CLI must be on PATH or set via `--ai-cmd`. Ralph validates before the loop and exits 2 with a clear error.
- **ralph: command not found** — Ensure the install directory is on your PATH. Check `~/.config/ralph/install-state` for the install path and run `ralph version` to verify.

## Writing Ralph prompts

Ralph evaluates prompts along four dimensions: signal and state, iteration awareness, scope and convergence, and subjective completion. For guidance on writing well-formed prompts, run:

```bash
ralph show prompt-guide
```

That outputs the full [Writing Ralph prompts](docs/writing-ralph-prompts.md) guide. You can also open the doc directly in the repo.

## License

GPL-3.0. See [LICENSE](LICENSE).
