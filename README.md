# Ralph

A dumb loop that pipes a prompt to an AI CLI, lets it work, and repeats.

Ralph is a loop runner, not a methodology. You bring the prompt. Ralph runs it in a fresh AI process per iteration, scans for completion signals, and stops when the task is done — or when it isn't going to be.

## What is Ralph and why use it?

- **What:** Ralph is a loop runner for AI-driven tasks. You give it a prompt (file, alias, or stdin); it invokes your AI CLI with that prompt on each iteration, captures output, and looks for success or failure signals you configure.
- **Why:** It replaces manual "run the AI, read the output, decide, re-run" with an automated loop that stops when the AI emits a success signal or hits a failure threshold. No conversation history between iterations — state lives on disk; each run is a fresh process.

If you're new here: install (below), then follow [Path to first run](#path-to-first-run) to get from "I have Ralph" to a first command that completes successfully.

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

## Install and Uninstall

Ralph can be installed with the provided script so the `ralph` binary is on your PATH. Uninstall removes only the binary and install state; your config (e.g. `~/.config/ralph/ralph-config.yml`) is not removed.

**Prerequisites:** `curl`. The install script **only** installs from release artifacts (no build from source). Supported: Linux, macOS, Windows (amd64, arm64); script runs on macOS/Linux or Windows (e.g. Git Bash).

**Install:**

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
3. Ensure the install directory is on your PATH (e.g. add it to your shell profile, or use `~/bin` if it is already on PATH).
4. Open a new terminal and verify:
   ```bash
   ralph version
   ```
   You should see version output and exit 0.

## Path to first run

From "I have Ralph" to a first command that completes successfully:

1. **Verify Ralph is on your PATH** — Run `ralph version`; you should see a version string and exit 0.
2. **Choose a prompt source** — You need exactly one: a prompt alias (from config), a file path, or stdin. With no config, the simplest path is a one-off file or stdin.
3. **Run a first command** — Examples that can complete successfully (your AI CLI must be installed and on PATH; Ralph defaults to the `cursor-agent` alias if you don't set `--ai-cmd` or config):
   - **From a file:** Create a small prompt file (e.g. `echo 'Say hello and output <promise>SUCCESS</promise>' > /tmp/hello.md`), then run:
     ```bash
     ralph run -f /tmp/hello.md -n 1
     ```
     On success you'll see the success message and exit 0.
   - **From stdin:** `echo 'Say OK and output <promise>SUCCESS</promise>' | ralph run -n 1`
   - **Using an alias:** Add a prompt to your config (see [Example ralph-config.yml](#example-ralph-configyml)), then `ralph run <alias>` (e.g. `ralph run build`).
4. **Prerequisites:** Your chosen AI CLI (e.g. Cursor agent, Claude CLI) must be installed and on PATH. Ralph does not install it. Use `--ai-cmd-alias` or `--ai-cmd` to select a different command; see Configuration and [config](docs/engineering/components/config.md).

**Uninstall:**

Run from anywhere (the script reads the install location from `~/.config/ralph/install-state`):

```bash
./scripts/uninstall.sh
```

This removes the `ralph` binary from the directory where it was installed and removes the install state file. User config in `~/.config/ralph/` (e.g. `ralph-config.yml`) is **not** removed. No PATH or symlink changes are made by the install script, so uninstall does not leave broken references.

**Upgrade:**

To upgrade or update an existing install:

- **Upgrade to a chosen version** — Re-run the install script with the desired version (same as a fresh install to that version). The script overwrites the binary in the same install directory (recorded in `~/.config/ralph/install-state`). Example:
  ```bash
  ./scripts/install.sh 1.2.0
  # or, to match your current install directory:
  ./scripts/install.sh 1.2.0 --dir "$(cat ~/.config/ralph/install-state)"
  ```
  After upgrade, run `ralph version` to confirm the new version; your config and prompts are unchanged.

- **Update within a non-breaking version** — To get the latest patch or minor release on the same major line (e.g. stay on 1.x), run the install script with no version to fetch the latest release, or with a specific newer tag in the same major:
  ```bash
  ./scripts/install.sh
  ```
  Within a non-breaking range (same major version), existing config, prompts, and documented commands/exit codes continue to work; no migration is required. Ralph does not rewrite or migrate your config automatically. If a release introduces behavior changes or deprecations, they are described in release notes with migration guidance (see [Release notes and contract](#release-notes-and-stable-contract)).

- **Backward compatibility** — Config and prompt files that worked in an older patch/minor (same major) keep working after update. Documented commands, options, and exit codes remain valid across non-breaking upgrades. Breaking contract or config changes are documented in release notes with migration steps; major version bumps may introduce breaking changes as documented.

**Release notes and stable contract:** For each release, behavior changes and deprecations that affect config or scripts are described in the release notes (e.g. [GitHub Releases](https://github.com/maxdunn/ralph/releases)). See [docs/exit-codes.md](docs/exit-codes.md) for the stable exit-code contract and [docs/release-notes.md](docs/release-notes.md) for where to find release notes and what constitutes the stable contract (exit codes, review summary format, config). When that contract or config schema changes, release notes explain the change and how to adapt.

## How It Works

Each iteration:

1. Read the prompt file mapped to the alias
2. Wrap it with a preamble (iteration count, optional context)
3. Pipe the assembled prompt to the AI CLI's stdin
4. Capture output, scan for success/failure signals
5. On success signal → exit 0
6. On failure signal → increment consecutive failure counter
7. On consecutive failure threshold → exit 4
8. On max iterations reached → exit 3
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

### Config file schema

Config files are YAML. The canonical structure (see `docs/engineering/components/config.md`) is:

- **loop** (optional) — Root loop behavior: `max_iterations`, `failure_threshold`, `timeout_seconds`, `success_signal`, `failure_signal`, `signal_precedence`, `preamble`, `streaming`, `log_level`. Per-prompt overrides go under each prompt’s `loop`.
- **prompts** (optional) — Map of prompt name to definition: `path` or `content`; optional `display_name`, `description`, and `loop` overrides.
- **aliases** (optional) — Map of alias name to AI command string (or `{ command: "..." }`). Built-in aliases (e.g. `claude`, `cursor-agent`) are merged; user aliases override built-ins for the same name.

The default AI command for `ralph run` is chosen by: `--ai-cmd` or `--ai-cmd-alias` (or env `RALPH_LOOP_AI_CMD` / `RALPH_LOOP_AI_CMD_ALIAS`); if none is set, Ralph uses the `cursor-agent` alias. There is no `loop.ai_cmd_alias` in the file schema; use CLI or env to select an alias.

### Example `ralph-config.yml`

```yaml
loop:
  max_iterations: 5
  failure_threshold: 3
  timeout_seconds: 300
  streaming: false
  success_signal: "<promise>SUCCESS</promise>"
  failure_signal: "<promise>FAILURE</promise>"

aliases:
  claude: "claude -p --dangerously-skip-permissions"
  kiro: "kiro-cli chat --no-interactive --trust-all-tools"
  copilot: "copilot --yolo"
  cursor-agent: "agent -p --force --output-format stream-json --stream-partial-output"

prompts:
  build:
    path: "./prompts/build.md"
    display_name: "Build"
    description: "Run the main build loop"
    loop:
      max_iterations: 10
      failure_threshold: 5

  bootstrap:
    path: "./prompts/bootstrap.md"
    display_name: "Bootstrap"
    description: "One-shot project setup"
    loop:
      max_iterations: 1
      preamble: false
```

Each prompt entry maps a name to a file (or `content`) and can override any loop setting under its `loop` key.

## Signals

Ralph scans AI CLI output for configurable signal strings to determine iteration outcome.

| Signal | Default | Meaning |
|--------|---------|---------|
| Success | `<promise>SUCCESS</promise>` | Task complete, stop looping |
| Failure | `<promise>FAILURE</promise>` | Blocked, increment failure counter |

Your prompt tells the AI what to emit. Ralph's signal config tells the scanner what to look for — use whatever strings your prompt expects.

With default `signal_precedence: static`, if both signals appear in the same output, success is checked first — success wins. See `docs/engineering/components/run-loop.md` for `ai_interpreted` precedence.

## Subcommands

Ralph has five top-level commands: **run**, **review**, **list**, **show**, **version**. Install, uninstall, and upgrade are documented procedures (scripts or package manager), not subcommands.

### ralph run

Run the iteration loop. Prompt source: exactly one of alias, `--file`/`-f` &lt;path&gt;, or stdin. The prompt is read once at loop start and reused for every iteration.

```
ralph run [alias] [flags]          Prompt from config alias
ralph run -f <path> [flags]        Prompt from file
cat prompt.md | ralph run [flags]  Prompt from stdin
```

Flags (all optional): `-f, --file`, `-n, --max-iterations`, `-u, --unlimited`, `--failure-threshold`, `--iteration-timeout`, `--max-output-buffer`, `--no-preamble`, `-d, --dry-run`, `--ai-cmd`, `--ai-cmd-alias`, `--signal-success`, `--signal-failure`, `--signal-precedence`, `-c, --context` (repeatable), `-v, --verbose`, `-q, --quiet`, `--log-level`, `--stream`, `--no-stream`, `--config`. Use `ralph run --help` for full list.

### ralph review

Review a prompt (alias, file, or stdin). Produces a report (narrative + machine-parseable summary) and a suggested revision; optionally apply the revision with confirmation (or `--yes` in non-interactive mode).

```
ralph review [alias] [flags]
ralph review -f <path> [flags]
cat prompt.md | ralph review [flags]
```

Flags: `-f, --file`, `--report` (report output path; default `./ralph-review-report.txt`), `--prompt-output` (required when using `--apply` with stdin), `--apply`, `--yes`/`-y` (non-interactive apply), `-v, -q, --quiet`, `--log-level`, `--config`. For CI: use exit code 0/1/2 to gate; or parse the report for a line `ralph-review: status=ok|errors|warnings` (see [docs/exit-codes.md](docs/exit-codes.md) and `docs/engineering/components/review.md`).

### ralph list

List prompts and/or AI command aliases from resolved config. Same config resolution as run (global, workspace, or `--config`).

```
ralph list [prompts|aliases]       List all, or only prompts or only aliases
```

### ralph show

Show effective config or detail for a prompt or alias. Same config resolution as run.

```
ralph show config [flags]          Effective (resolved) config
ralph show prompt <name>           Detail for prompt <name>
ralph show alias <name>            Detail for alias <name>
```

Name is required for `show prompt` and `show alias`. Use `--provenance` with `show config` to see which layer supplied each value.

### ralph version

Print version string to stdout and exit 0. No arguments required.

## Environment variables

- **RALPH_CONFIG_HOME** — Directory for the global config file; actual file is `$RALPH_CONFIG_HOME/ralph-config.yml`. Does not set the explicit config file for the current run; use `--config <path>` for that.
- **RALPH_LOOP_*** — Override loop settings: `RALPH_LOOP_AI_CMD`, `RALPH_LOOP_AI_CMD_ALIAS`, `RALPH_LOOP_DEFAULT_MAX_ITERATIONS`, `RALPH_LOOP_FAILURE_THRESHOLD`, `RALPH_LOOP_ITERATION_TIMEOUT`, `RALPH_LOOP_LOG_LEVEL`, `RALPH_LOOP_STREAMING`, `RALPH_LOOP_PREAMBLE`, etc. See `docs/engineering/components/config.md` for the full set so that full non-interactive config is possible (e.g. in CI).

## Non-interactive use (scripts and CI)

- **run:** No confirmation prompts; use `--config`, env vars, and flags to drive behavior. Exit codes are stable (see Exit Codes and [docs/exit-codes.md](docs/exit-codes.md)).
- **review:** Use `--yes` when applying in non-interactive mode so Ralph does not block on confirmation; without `--yes`, Ralph exits 2 with a clear message. For stdin + apply, you must pass `--prompt-output` or Ralph exits 2.
- Detection of non-interactive context (e.g. no TTY or `CI=true`) is implementation-defined; `--yes` always suppresses confirmation for `--apply`.

## Exit Codes

Exit codes for `ralph run` and `ralph review` are stable for scripts and CI. Full semantics and automation guidance: [docs/exit-codes.md](docs/exit-codes.md).

**ralph run**

| Code | Meaning |
|------|---------|
| 0 | Success signal received |
| 2 | Error before loop (e.g. missing/invalid AI command) |
| 3 | Max iterations exhausted |
| 4 | Failure threshold reached |
| 130 | Interrupted (SIGINT/SIGTERM) |

**ralph review**

| Code | Meaning |
|------|---------|
| 0 | Review completed, no prompt errors |
| 1 | Review completed, prompt has errors |
| 2 | Review or apply did not complete |

## Where to look (CLI, config, env, exit codes)

- **CLI and flags** — Full command and flag spec: [docs/engineering/components/cli.md](docs/engineering/components/cli.md). User-facing summary is in this README (Subcommands, Configuration, Environment variables, Exit Codes).
- **Config file and layers** — Schema, layer order, and env overlay: [docs/engineering/components/config.md](docs/engineering/components/config.md). README summarizes in Configuration and Example.
- **Environment variables** — `RALPH_CONFIG_HOME` and `RALPH_LOOP_*` are listed in [docs/engineering/components/config.md](docs/engineering/components/config.md) (Environment variables). README lists them in Environment variables.
- **Exit codes** — Stable contract for run and review: [docs/exit-codes.md](docs/exit-codes.md). README summarizes in Exit Codes.
- **Release notes and stable contract** — Where to find release notes and what scripts/CI can rely on: [docs/release-notes.md](docs/release-notes.md). Upgrade and backward-compatibility procedure is in this README (Install and Uninstall → Upgrade); see [Release notes and stable contract](#release-notes-and-stable-contract) above.

When implementation or contract changes (e.g. new flag, config key, or exit code), update the engineering spec and then this README (and exit-codes.md if needed) so docs stay accurate.

## License

GPL-3.0. See [LICENSE](LICENSE).
