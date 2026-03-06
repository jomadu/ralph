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
| `cursor-agent` | `agent -p --force --output-format stream-json --stream-partial-output` | Raw Cursor Agent CLI; works from any CWD. Output is newline-delimited JSON. Optional wrapper script (see user docs) parses JSON and emits plain text for conventional signal scanning. |

The cursor-agent built-in alias points at the raw `agent` command so it works from any working directory. Cursor Agent emits newline-delimited JSON with typed messages (`system`, `assistant`, `tool_call`, `result`). Users may configure signal strings to match content within JSON lines, or use an optional wrapper script (documented in user-facing docs) that invokes the same command, parses JSON with `jq`, and emits plain text to stdout for signal scanning.

## Verification

- User runs `ralph run build --ai-cmd-alias claude`, then `ralph run build --ai-cmd-alias kiro`. Both execute successfully using the same prompt file.
- User defines a custom alias in config for a proprietary AI CLI: `my-tool: "my-ai-cli --headless --stdin"`. Ralph resolves and executes it.
- User passes `--ai-cmd "custom-tool --flag1 --flag2"` directly on the command line, bypassing aliases entirely. Ralph parses and executes the command.
- The AI CLI process inherits the user's environment variables (API keys, PATH) and working directory.
- User runs with the `cursor-agent` alias. Ralph invokes the raw agent command; output is JSON. User may use signal strings that match JSON content or the optional wrapper (user docs) for plain-text stdout.

## Non-outcomes

- Ralph does not validate that the AI CLI is installed or functional before execution. A missing binary produces a process start error on the first iteration.
- Ralph does not manage AI CLI authentication, API keys, or session tokens.
- Ralph does not normalize output formats across different AI CLIs — that responsibility belongs to the alias command or wrapper script.
- Ralph does not execute commands through a shell. Commands are parsed and exec'd directly. Users needing pipes, redirects, or shell expansion must wrap their command in a script.
- Ralph does not provide vendor-specific adapters in core. The cursor-agent built-in is the raw CLI; an optional wrapper script is documented as a user-facing workaround for plain-text stdout.

## Risks

| Risk | Mitigating Requirement |
|----------|----------------------|
| Each AI CLI has different flags for non-interactive stdin-based usage | [R1 — Built-in command aliases](R1-builtin-aliases.md) |
| Command strings require shell-style quoting (e.g., `--model "claude-3-5-sonnet"`) but Ralph doesn't use a shell | [R2 — Shell-style command parsing](R2-command-parsing.md) |
| User wants to use an AI CLI not in the built-in alias list | [R3 — User-defined command aliases](R3-user-defined-aliases.md) |
| AI CLI requires environment variables (API keys, config paths) to function | [R4 — Process environment inheritance](R4-environment-inheritance.md) |
| An AI CLI emits structured output (JSON) instead of plain text | [R1 — Built-in command aliases (raw command + optional wrapper in user docs)](R1-builtin-aliases.md) |
| User doesn't know which aliases are available or what commands they expand to | [R5 — AI command alias resolution with clear errors](R5-alias-resolution-errors.md) |
| Direct command string and alias are both specified, behavior is ambiguous | [R6 — Command source precedence](R6-command-source-precedence.md) |
| User uses optional wrapper; wrapper has its own dependency (jq) | Wrapper documents and reports missing dependency; user-facing docs describe the workaround. |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R1](R1-builtin-aliases.md) | Built-in command aliases | draft |
| [R2](R2-command-parsing.md) | Shell-style command parsing | draft |
| [R3](R3-user-defined-aliases.md) | User-defined command aliases | draft |
| [R4](R4-environment-inheritance.md) | Process environment inheritance | draft |
| [R5](R5-alias-resolution-errors.md) | AI command alias resolution with clear errors | draft |
| [R6](R6-command-source-precedence.md) | Command source precedence | draft |
