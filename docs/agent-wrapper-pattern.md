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
   - **Config:** In `ralph-config.yml`, add an alias whose command is your script, e.g. `"/path/to/my-wrapper.sh"` or `"bash /path/to/my-wrapper.sh"`. Use `--ai-cmd-alias <name>` or `RALPH_LOOP_AI_CMD_ALIAS`.
   - **CLI:** `ralph run build --ai-cmd "/path/to/my-wrapper.sh"`.

Ralph pipes the assembled prompt to the wrapper’s stdin and reads the wrapper’s stdout for signal scanning; stderr is your progress output.

## Example: Cursor agent

The built-in `cursor-agent` alias runs the Cursor CLI directly. That works for signal scanning but doesn’t show real-time progress. To get progress **and** signal scanning, use a wrapper that invokes the Cursor agent with stream-json and splits output as above.

**Cursor’s canonical docs (output format and example script):**  
[https://cursor.com/docs/cli/headless#real-time-progress-tracking](https://cursor.com/docs/cli/headless#real-time-progress-tracking)

That page describes `--output-format stream-json` and `--stream-partial-output`, and includes a full example that parses JSON lines with `jq`, sends progress to the terminal, and could be adapted to send assistant text to stdout and everything else to stderr for use with Ralph.

Other agents that emit structured or multi-channel output can follow the same pattern: implement a wrapper that separates “text to scan” (stdout) from “progress” (stderr), then point Ralph at the wrapper.
