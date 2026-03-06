# R1: Built-in Command Aliases

**Outcome:** O3 — Backend Agnosticism

## Requirement

The system ships with built-in command aliases for known AI CLIs, encoding the correct flags and invocation protocols for non-interactive, stdin-based execution. Built-in aliases eliminate the need for users to reverse-engineer each tool's invocation requirements. For AI CLIs that emit non-standard output (e.g., structured JSON), the built-in alias resolves to the raw CLI command so it works from any working directory; an optional wrapper script that normalizes output is documented as a user-facing workaround (see user docs).

## Specification

Ralph ships with a fixed set of **built-in** AI command aliases. These are available without any user configuration. Built-in aliases are defined in code or bundled config; they are not read from the user's config files.

**Minimum set:** The built-in alias map must include at least these four keys and the following command strings (or equivalent that achieves non-interactive, stdin-based execution):

| Alias | Command string |
|-------|----------------|
| `claude` | `claude -p --dangerously-skip-permissions` |
| `kiro` | `kiro-cli chat --no-interactive --trust-all-tools` |
| `copilot` | `copilot --yolo` |
| `cursor-agent` | `agent -p --force --output-format stream-json --stream-partial-output` (raw Cursor Agent CLI; requires `agent` on PATH). For plain-text stdout and signal scanning, users may override with a wrapper script — see user-facing docs (e.g. Cursor Agent workaround). |

**Override semantics:** User-defined aliases (R3) are merged with built-in aliases. If the user defines an alias with the same name as a built-in (e.g. `claude`), the user-defined value overrides the built-in for that name. The merged map is what R5/R6 use for resolution.

**Command string format:** Each built-in alias value is a single string that will be parsed per R2. It must be a complete, invocable command line (program name or path plus flags). No shell is used; the string is parsed into argv and exec'd (R2, R4).

**cursor-agent and optional wrapper:**

- The `cursor-agent` built-in alias resolves to the raw command `agent -p --force --output-format stream-json --stream-partial-output`. This ensures the alias works from any working directory (no dependency on a script path in the repo). The Cursor Agent emits newline-delimited JSON; signal strings may be configured to match content within JSON lines (see O1 signal specs), or the user may use an optional wrapper script that parses JSON and emits plain text for conventional signal scanning.
- **Optional wrapper:** A wrapper script (e.g. `scripts/cursor-wrapper.sh` in the Ralph repo) that invokes the same agent command, parses JSON with `jq`, and emits plain text to stdout is documented as a user-facing workaround: [docs/user/cursor-agent-workaround.md](../../user/cursor-agent-workaround.md). Users who want plain-text stdout can set `ai_cmd` or a user-defined alias to the path of that script. The wrapper is not an embedded resource; it is documented under user-facing docs derived from the intent tree.

**Availability:** Built-in aliases are always present. They do not require the user to create a config file or set any options. Listing available aliases (e.g. for R5 error messages or `ralph list aliases`) must include built-in aliases plus any user-defined aliases from config.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User has no config file | Built-in aliases (claude, kiro, copilot, cursor-agent) are available. `--ai-cmd-alias claude` works. |
| User defines `ai_cmd_aliases.claude: "claude -p --model other"` | Merged map uses user value for `claude`. Built-in value is overridden. |
| `agent` not on PATH | Process start fails (e.g. "executable file not found"). Same as any missing binary. |
| User overrides `cursor-agent` with wrapper path (e.g. `ai_cmd_aliases.cursor-agent: "./scripts/cursor-wrapper.sh"`) | Ralph invokes the wrapper; wrapper emits plain text; signal scanning works on wrapper stdout. Wrapper must be present at that path (CWD-relative or absolute). |
| `cursor-agent` alias (built-in) used as-is | Ralph invokes raw `agent ...`; output is JSON. Signal strings may be set to match JSON content, or user may switch to wrapper for plain-text stdout. |

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

#### cursor-agent with built-in alias

**Input:**
User runs `ralph run build --ai-cmd-alias cursor-agent`. `agent` is on PATH.

**Expected output:**
Ralph resolves `cursor-agent` to `agent -p --force --output-format stream-json --stream-partial-output`, parses it (R2), and executes. Loop runs; stdout is JSON. User may configure signal strings to match JSON lines or use the optional wrapper (user docs) for plain-text stdout.

## Acceptance criteria

- [ ] Built-in aliases include at minimum: claude, kiro, copilot, cursor-agent
- [ ] Each alias resolves to a complete command string with all necessary flags for non-interactive, stdin-based operation
- [ ] Built-in aliases are available without any user configuration
- [ ] User-defined aliases with the same name as a built-in alias override the built-in
- [ ] The cursor-agent built-in alias resolves to the raw command `agent -p --force --output-format stream-json --stream-partial-output` and works from any working directory
- [ ] User-facing docs (from intent tree) document the optional Cursor Agent wrapper script as a workaround for plain-text stdout

## Dependencies

_None identified._
