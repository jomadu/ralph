# Agent wrapper pattern (progress vs. signal scanning)

Many AI CLIs emit **structured or noisy output** (JSON lines, tool calls, progress messages) rather than plain text. Ralph scans **stdout** for your configured success/failure signals (e.g. `<promise>SUCCESS</promise>`). If the agent writes everything to stdout, signal detection still works, but you may not see clear progress in the terminal.

A **wrapper script** can split the agent’s output so that:

- **Progress and tool output** go to **stderr** — you see them live (model name, tool calls, character count, etc.).
- **Only the assistant’s text** goes to **stdout** — Ralph scans it for your signals.

This pattern applies to **any** agent whose CLI outputs structured or multi-channel data. Ralph does not ship agent-specific wrappers; you implement a small script per agent and point Ralph at it via `--ai-cmd` or a config alias.

## General pattern

1. **Wrapper behavior**
   - Read the prompt from **stdin** and pass it to the AI CLI (with whatever flags that agent needs).
   - Parse the CLI output (e.g. JSON lines). For each message:
     - If it’s **assistant text** → print it to **stdout** (so Ralph can detect signals).
     - Otherwise (system init, tool_call, result, progress) → print human-readable progress to **stderr**.
   - The wrapper must be invocable as a single command (script path or `bash /path/to/script.sh`).

2. **Wire the wrapper in Ralph**
   - **Config:** In `ralph-config.yml`, set `loop.ai_cmd_alias` to an alias whose command is your script, or `loop.ai_cmd` to the script path (e.g. `"/path/to/my-wrapper.sh"`). You can also use `--ai-cmd-alias <name>` or `RALPH_LOOP_AI_CMD_ALIAS`, or `--ai-cmd` / `RALPH_LOOP_AI_CMD`.
   - **CLI:** `ralph run build --ai-cmd "/path/to/my-wrapper.sh"`.

Ralph pipes the assembled prompt to the wrapper’s stdin and reads the wrapper’s stdout for signal scanning; stderr is your progress output.

**Assembled prompt shape:** Ralph wraps your prompt with section title lines such as `# --- CONTEXT ---` (preamble and optional `-c` text) and `# --- INSTRUCTIONS ---` (your prompt file). The first line of the message is never a bare `---` line—some CLIs (e.g. Cursor `agent`) treat that as YAML frontmatter and fail immediately. See [run-loop component](engineering/components/run-loop.md) (section title lines) for the exact format and alternatives that were considered.

## Paths: global config, workspace, and relative scripts

**Prompt file paths** in config are resolved **relative to the YAML file that defines them** (then merged across layers). **AI command aliases are different:** the `command` string for an alias is used **as-is**. Ralph does not rewrite relative paths in alias commands to be relative to the global config directory, the workspace config file, or anything else.

The backend runs the resolved command **without a shell**. The first token is the executable. If that token is a **relative path** (e.g. `./scripts/cursor-wrapper.sh`), the OS resolves it from the **current working directory of the Ralph process** — typically **where you ran `ralph`** (`cd` matters), not where your config file lives.

Implications:

- A **global** alias like `command: "./cursor-wrapper.sh"` only works when your shell cwd happens to contain that path. Running the same command from another project or directory will look for a different file (or fail).
- **Workspace** config with `./scripts/wrapper.sh` is still **cwd-dependent**: it is not “relative to the workspace `ralph-config.yml`” unless you always run Ralph from the repo root and keep the script at that path there.

**Practical guidance:**

1. **Global or shared default** — Use an **absolute path** to the wrapper (e.g. `/Users/you/bin/cursor-ralph-wrapper.sh`), or install the wrapper on **`PATH`** and use a bare command name (e.g. `cursor-ralph-wrapper`).
2. **Per-repository** — Keep the wrapper in the repo and either use an absolute path in workspace config, rely on always running from the repo root with a stable relative path (convention), or put a small launcher on PATH that delegates to the repo.
3. **One-off** — `ralph run … --ai-cmd "/absolute/path/to/wrapper.sh"` avoids ambiguity.

See also: [config component](engineering/components/config.md) (aliases and prompt path resolution).

## Example: Cursor agent

The built-in `cursor-agent` alias runs the Cursor CLI directly. That works for signal scanning but doesn’t show real-time progress. To get progress **and** signal scanning, use a wrapper that invokes the Cursor agent with stream-json and splits output as above.

**Cursor’s canonical docs (output format and example script):**  
[https://cursor.com/docs/cli/headless#real-time-progress-tracking](https://cursor.com/docs/cli/headless#real-time-progress-tracking)

That page describes `--output-format stream-json` and `--stream-partial-output`, and includes a full example that parses JSON lines with `jq`, sends progress to the terminal, and could be adapted to send assistant text to stdout and everything else to stderr for use with Ralph.

**Duplicate stdout with Cursor:** If you append every `assistant` event’s `message.content[0].text` and also print that blob at the end, you can see the **same reply twice**. Partial streaming sometimes re-emits the full (or overlapping) text across multiple assistant lines. For Ralph, print **one** final payload on stdout: prefer the terminal stream-json **`result`** event’s **`result`** field (authoritative full text on success), and only fall back to concatenated assistant text if there is no success result (e.g. error exit). The repo includes **`scripts/cursor-wrapper.sh`** wired that way; copy or adapt it for your config.

Other agents that emit structured or multi-channel output can follow the same pattern: implement a wrapper that separates “text to scan” (stdout) from “progress” (stderr), then point Ralph at the wrapper.
