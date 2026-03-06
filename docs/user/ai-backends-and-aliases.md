# AI backends and aliases

**Intent:** [O3 — Backend agnosticism](../intent/O3-backend-agnosticism/README.md); [R1 — Built-in command aliases](../intent/O3-backend-agnosticism/R1-builtin-aliases.md); [R3 — User-defined command aliases](../intent/O3-backend-agnosticism/R3-user-defined-aliases.md); [R6 — Command source precedence](../intent/O3-backend-agnosticism/R6-command-source-precedence.md)

Ralph runs your prompt through an AI CLI that accepts stdin and writes to stdout. You can use **built-in aliases** (known AI tools), **user-defined aliases** (your own names and commands in config), or a **direct command** (path or command string). This page explains how to choose or override the backend.

## Three ways to specify the AI command

| Method | What you set | Use when |
|--------|----------------|----------|
| **Built-in alias** | A name Ralph already knows: `claude`, `kiro`, `copilot`, `cursor-agent` | You use one of these tools and want the correct flags without config |
| **User-defined alias** | A name and command in `ai_cmd_aliases` in your config | You want a custom name, different flags, or a tool not in the built-in list |
| **Direct command** | A full command string via `--ai-cmd` or `loop.ai_cmd` | You want to point at a script or binary without defining an alias |

**Precedence:** If you set both a direct command and an alias (e.g. `--ai-cmd "..."` and `--ai-cmd-alias claude`), the **direct command wins**. So you can override an alias for a single run with `--ai-cmd`.

## Built-in aliases

Ralph ships with four built-in aliases. No config is required; they are always available.

| Alias | Resolved command | Notes |
|-------|------------------|--------|
| `claude` | `claude -p --dangerously-skip-permissions` | Pipe mode; skips interactive permission prompts |
| `kiro` | `kiro-cli chat --no-interactive --trust-all-tools` | Chat subcommand; non-interactive, trust tools |
| `copilot` | `copilot --yolo` | Unattended execution |
| `cursor-agent` | `agent -p --force --output-format stream-json --stream-partial-output` | Raw Cursor Agent; stdout is JSON. For plain-text stdout and signal scanning, see [Cursor Agent workaround](cursor-agent-workaround.md) |

**Example:** Run with Claude (no config):

```bash
ralph run build --ai-cmd-alias claude
```

## User-defined aliases

Define your own aliases in `ai_cmd_aliases` in your config file. User-defined aliases **merge** with built-in ones. If you use the same name as a built-in (e.g. `claude`), your value **overrides** the built-in for that name.

**Config (workspace or global):**

```yaml
# ./ralph-config.yml or ~/.config/ralph/ralph-config.yml
ai_cmd_aliases:
  my-ai: "my-ai-cli --headless --stdin"
  claude: "claude -p --model claude-3-5-sonnet"   # overrides built-in claude
```

Then:

```bash
ralph run build --ai-cmd-alias my-ai
ralph run build --ai-cmd-alias claude   # uses your command, not built-in
```

Config merge rules (same as the rest of Ralph): workspace overrides global for the same key. With `--config /path/to/file.yml`, only that file is used (no workspace/global merge).

## Direct command

Bypass aliases and pass a full command string. Useful for a one-off script or binary.

**CLI:**

```bash
ralph run build --ai-cmd "/path/to/scripts/cursor-wrapper.sh"
ralph run build --ai-cmd "custom-ai --headless --stdin"
```

**Config:**

```yaml
loop:
  ai_cmd: "/path/to/your/agent-or-wrapper"
```

**Environment:** `RALPH_LOOP_AI_CMD` (same precedence as other env vars).

Direct command always wins over an alias when both are set at the same layer (e.g. both in config).

## Where to set alias vs direct command

Same precedence as [Configuration](configuration.md): CLI → environment → prompt-level loop → workspace config → global config. There is **no default** AI command; you must set either `ai_cmd_alias` or `ai_cmd` (or both; then direct command wins). If neither is set, Ralph exits with an error before starting the loop.

| Source | Alias | Direct command |
|--------|--------|-----------------|
| CLI | `--ai-cmd-alias <name>` | `--ai-cmd "<string>"` |
| Environment | `RALPH_LOOP_AI_CMD_ALIAS` | `RALPH_LOOP_AI_CMD` |
| Config | `loop.ai_cmd_alias` (root or under `prompts.<alias>.loop`) | `loop.ai_cmd` |

## Cursor Agent and the optional wrapper

The built-in `cursor-agent` alias runs the raw `agent` command. Output is newline-delimited JSON, not plain text. You can either configure signal strings to match content inside that JSON, or use an **optional wrapper script** that parses the JSON and emits plain text so normal signal scanning works.

See **[Cursor Agent workaround](cursor-agent-workaround.md)** for how to point Ralph at the wrapper (config or `--ai-cmd`), dependencies (`jq`, `agent`), and usage.
