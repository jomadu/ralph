# R1: Built-in Command Aliases

**Outcome:** O3 — Backend Agnosticism

## Requirement

The system ships with built-in command aliases for known AI CLIs, encoding the correct flags and invocation protocols for non-interactive, stdin-based execution. Built-in aliases eliminate the need for users to reverse-engineer each tool's invocation requirements. For AI CLIs that emit non-standard output (e.g., structured JSON instead of plain text), the alias resolves to a wrapper script that normalizes the output.

## Specification

Ralph ships with a fixed set of **built-in** AI command aliases. These are available without any user configuration. Built-in aliases are defined in code or bundled config; they are not read from the user's config files.

**Minimum set:** The built-in alias map must include at least these four keys and the following command strings (or equivalent that achieves non-interactive, stdin-based execution):

| Alias | Command string |
|-------|----------------|
| `claude` | `claude -p --dangerously-skip-permissions` |
| `kiro` | `kiro-cli chat --no-interactive --trust-all-tools` |
| `copilot` | `copilot --yolo` |
| `cursor-agent` | Path to the bundled wrapper script: `scripts/cursor-wrapper.sh` (relative to Ralph's installation root or resolved such that the script is invocable). The wrapper reads stdin, invokes the underlying agent with the correct flags, parses JSON output, and emits plain text to stdout. |

**Override semantics:** User-defined aliases (R3) are merged with built-in aliases. If the user defines an alias with the same name as a built-in (e.g. `claude`), the user-defined value overrides the built-in for that name. The merged map is what R5/R6 use for resolution.

**Command string format:** Each built-in alias value is a single string that will be parsed per R2. It must be a complete, invocable command line (program name or path plus flags). No shell is used; the string is parsed into argv and exec'd (R2, R4).

**cursor-agent wrapper:**

- The `cursor-agent` alias must resolve to a command that invokes the wrapper script (e.g. absolute or relative path to `cursor-wrapper.sh`). The wrapper is responsible for: reading the prompt from stdin, piping to the underlying agent, parsing newline-delimited JSON, extracting text from `assistant` message content, routing tool-call progress to stderr, and emitting plain text to stdout for signal scanning.
- **Missing dependency:** The wrapper script may depend on external tools (e.g. `jq`, `agent`). If the wrapper detects a missing dependency at startup, it must exit with a non-zero code and emit a clear error to stderr (e.g. "jq is required but not installed" and install hint). Ralph does not parse wrapper stderr for this; the wrapper's exit code and message are the contract. Ralph treats wrapper failure like any other process failure (e.g. exit code propagated; no special "missing dependency" handling in Ralph core).

**Availability:** Built-in aliases are always present. They do not require the user to create a config file or set any options. Listing available aliases (e.g. for R5 error messages or `ralph list aliases`) must include built-in aliases plus any user-defined aliases from config.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User has no config file | Built-in aliases (claude, kiro, copilot, cursor-agent) are available. `--ai-cmd-alias claude` works. |
| User defines `ai_cmd_aliases.claude: "claude -p --model other"` | Merged map uses user value for `claude`. Built-in value is overridden. |
| Wrapper script not found at resolved path | Process start fails (e.g. "no such file or directory"). No special message from Ralph; same as any missing binary. |
| Wrapper runs but `jq` is missing | Wrapper exits non-zero and prints error to stderr. Ralph sees non-zero exit; no special handling. |
| `cursor-agent` alias points to wrapper; wrapper emits plain text | Ralph pipes prompt to wrapper stdin; signal scanning reads wrapper stdout. Behavior is the same as for any other alias that emits plain text. |

### Examples

#### Built-in alias without config

**Input:**
No config files. User runs `ralph run build --ai-cmd-alias kiro`.

**Expected output:**
Ralph resolves `kiro` to `kiro-cli chat --no-interactive --trust-all-tools` (or the built-in value), parses it (R2), and executes. Loop runs.

**Verification:**
- No "unknown alias" error. Process invoked is the built-in kiro command.

#### User overrides built-in

**Input:**
Config: `ai_cmd_aliases: { claude: "claude -p --model claude-3-5-sonnet" }`. User runs `ralph run build --ai-cmd-alias claude`.

**Expected output:**
Resolved command string is the user value `claude -p --model claude-3-5-sonnet`, not the built-in. Loop runs with that command.

**Verification:**
- Debug log or process listing shows the user's command string. Built-in is not used.

#### cursor-agent and missing jq

**Input:**
User runs `ralph run build --ai-cmd-alias cursor-agent`. `jq` is not installed; wrapper script is present.

**Expected output:**
Wrapper script runs, detects missing `jq`, prints a clear error to stderr (e.g. "jq is required but not installed" and install hint), and exits non-zero. Ralph sees the exit code; iteration fails. No need for Ralph to parse the message — the wrapper owns the message.

**Verification:**
- Wrapper stderr contains an actionable message about the missing dependency. Ralph does not need to interpret it; user can fix by installing jq.

## Acceptance criteria

- [ ] Built-in aliases include at minimum: claude, kiro, copilot, cursor-agent
- [ ] Each alias resolves to a complete command string with all necessary flags for non-interactive, stdin-based operation
- [ ] Built-in aliases are available without any user configuration
- [ ] User-defined aliases with the same name as a built-in alias override the built-in
- [ ] The cursor-agent alias resolves to a wrapper script ([`cursor-wrapper.sh`](../../../scripts/cursor-wrapper.sh)) that parses structured JSON output and emits plain text suitable for signal scanning
- [ ] If a wrapper script has a missing dependency (e.g., jq), it reports the missing dependency clearly

## Dependencies

_None identified._
