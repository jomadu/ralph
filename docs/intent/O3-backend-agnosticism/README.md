# O3: Backend Agnosticism

## Statement

Any stdin-accepting AI CLI serves as the execution backend.

## Why it matters

The AI CLI landscape is fragmented. Teams use Claude, Kiro, Copilot, Cursor, or proprietary internal tools — and they switch between them. A loop runner locked to one CLI is a loop runner locked to one vendor. Ralph's value is the loop, not the AI. By treating the AI CLI as a pluggable process that accepts a prompt on stdin and produces output on stdout, Ralph works with whatever the user already has installed and avoids coupling to any vendor's invocation protocol.

The difficulty is that "accepts stdin, produces stdout" is the contract, but every AI CLI has its own way of getting there. Claude requires `-p` for pipe mode and `--dangerously-skip-permissions` to avoid interactive confirmation. Kiro needs a `chat` subcommand, `--no-interactive`, and `--trust-all-tools`. Cursor agent outputs structured JSON that requires parsing to extract text — it can't be consumed raw. Copilot has `--yolo` for unattended execution. Each tool has different flags for non-interactive mode, different approaches to permission handling, and different output formats. The built-in aliases encode this knowledge so the user doesn't have to reverse-engineer each tool's invocation protocol.

### Known AI CLI Commands

| Alias | Resolved command | Key flags |
|-------|-----------------|-----------|
| `claude` | `claude -p --dangerously-skip-permissions` | `-p`: pipe/stdin mode. `--dangerously-skip-permissions`: skip interactive permission prompts |
| `kiro` | `kiro-cli chat --no-interactive --trust-all-tools` | `chat`: subcommand for conversation. `--no-interactive`: suppress prompts. `--trust-all-tools`: skip tool confirmation |
| `copilot` | `copilot --yolo` | `--yolo`: unattended execution without confirmation |
| `cursor-agent` | [`cursor-wrapper.sh`](../../../scripts/cursor-wrapper.sh) | Wraps `agent -p --force --output-format stream-json --stream-partial-output`. Outputs structured JSON — requires a wrapper script using `jq` to parse JSON lines, extract text from `assistant` messages, route tool call progress to stderr, and emit plain text to stdout for signal scanning |

The cursor-agent case is the hardest. Unlike the other CLIs which emit plain text on stdout, cursor agent emits newline-delimited JSON with typed messages (`system`, `assistant`, `tool_call`, `result`). The text content is nested inside `message.content[0].text` fields on `assistant` type messages. Tool calls appear as separate JSON objects with nested structures (`writeToolCall`, `readToolCall`). A wrapper script must: read the prompt from stdin, pipe it to `agent`, parse each JSON line with `jq`, accumulate text deltas to stdout, route tool progress to stderr, and report final statistics. This wrapper introduces a `jq` dependency that only cursor-agent users pay.

## Verification

- User runs `ralph run build --ai-cmd-alias claude`, then `ralph run build --ai-cmd-alias kiro`. Both execute successfully using the same prompt file.
- User defines a custom alias in config for a proprietary AI CLI: `my-tool: "my-ai-cli --headless --stdin"`. Ralph resolves and executes it.
- User passes `--ai-cmd "custom-tool --flag1 --flag2"` directly on the command line, bypassing aliases entirely. Ralph parses and executes the command.
- The AI CLI process inherits the user's environment variables (API keys, PATH) and working directory.
- User runs with the `cursor-agent` alias. Ralph invokes the wrapper script, which parses JSON output and emits plain text. Signal scanning works on the emitted text.

## Non-outcomes

- Ralph does not validate that the AI CLI is installed or functional before execution. A missing binary produces a process start error on the first iteration.
- Ralph does not manage AI CLI authentication, API keys, or session tokens.
- Ralph does not normalize output formats across different AI CLIs — that responsibility belongs to the alias command or wrapper script.
- Ralph does not execute commands through a shell. Commands are parsed and exec'd directly. Users needing pipes, redirects, or shell expansion must wrap their command in a script.
- Ralph does not provide vendor-specific adapters in core. The cursor-agent wrapper script is the model: a shim that normalizes a non-conforming CLI to stdin-text-in / stdout-text-out.

## Obstacles

| Obstacle | Mitigating Requirement |
|----------|----------------------|
| Each AI CLI has different flags for non-interactive stdin-based usage | R1 — Built-in command aliases |
| Command strings require shell-style quoting (e.g., `--model "claude-3-5-sonnet"`) but Ralph doesn't use a shell | R2 — Shell-style command parsing |
| User wants to use an AI CLI not in the built-in alias list | R3 — User-defined command aliases |
| AI CLI requires environment variables (API keys, config paths) to function | R4 — Process environment inheritance |
| An AI CLI emits structured output (JSON) instead of plain text | R1 — Built-in command aliases (wrapper script) |
| User doesn't know which aliases are available or what commands they expand to | R5 — AI command alias resolution with clear errors |
| Direct command string and alias are both specified, behavior is ambiguous | R6 — Command source precedence |
| Wrapper script has its own dependency (jq) that may not be installed | R1 — Built-in command aliases (wrapper reports missing dependency) |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| R1 | Built-in command aliases | draft |
| R2 | Shell-style command parsing | draft |
| R3 | User-defined command aliases | draft |
| R4 | Process environment inheritance | draft |
| R5 | AI command alias resolution with clear errors | draft |
| R6 | Command source precedence | draft |
