# R2: Shell-Style Command Parsing

**Outcome:** O3 — Backend Agnosticism

## Requirement

The system parses command strings using shell-style quoting rules — supporting quoted arguments and escaped characters — without invoking a shell. Commands are exec'd directly as processes, not passed through sh/bash.

## Specification

Command strings appear in two places: (1) the effective `ai_cmd` value (direct command from CLI, env, or config), and (2) alias values in the merged alias map (built-in and user-defined). Both are parsed with the same rules before execution.

**Parsing rules (shell-style, no shell):**

- **Whitespace:** Consecutive spaces/tabs separate arguments. Leading/trailing whitespace is trimmed before parsing.
- **Double-quoted strings:** Text between `"` and `"` forms a single argument. Within double quotes, `\"` produces a literal `"`, `\\` produces a literal `\`. Other backslash sequences are implementation-defined (e.g. `\n` may be newline) or treated as the next character. The closing `"` must be present; otherwise treat as parse error.
- **Single-quoted strings:** Text between `'` and `'` forms a single argument. Within single quotes, no escape sequences are recognized — every character is literal except the closing `'`. The closing `'` must be present; otherwise parse error.
- **Unquoted tokens:** Characters that are not whitespace, `"`, or `'` form a token. Consecutive unquoted tokens are separate arguments. No shell interpretation: no variable expansion, no glob expansion, no pipe/redirect handling.
- **Empty string:** An empty command string (after trim) is invalid; no arguments. Callers (R5/R6) avoid passing empty when a command is required.

**Output:** Parsing produces an ordered list of strings (argv). The first element is the program name (or path); the rest are arguments. This list is passed to exec without invoking a shell: Ralph spawns the process with argv and the inherited environment (R4).

**No shell:** The process is started via exec-family semantics (e.g. `execve`): argv[0] is the binary to run, argv[1..] are arguments. Pipes, redirects, `$VAR`, `*`, `&&`, `;` are not interpreted. Users who need shell features must wrap the command in a script and invoke the script as the command.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| `--model "claude-3-5-sonnet"` | Two arguments: `--model`, `claude-3-5-sonnet`. |
| `cli --flag 'single quoted'` | Three arguments: `cli`, `--flag`, `single quoted`. |
| `cli --path="value"` | Two arguments: `cli`, `--path=value` (no special handling for `=`; it's part of the token). |
| Command string contains `|` or `>` | They are literal characters, not pipe/redirect. Resulting argv contains `|` or `>` as part of an argument. |
| Command string contains `$HOME` | Literal string `$HOME`, no expansion. |
| Unclosed `"` or `'` | Parse error: fail with clear error (e.g. "unclosed quote"), do not start the process. |
| Empty command string | Zero arguments; invalid for execution. Treated as "no command" by callers (R5). |
| Only whitespace | Same as empty. |
| Backslash inside double quotes | Per rules: `\"` and `\\` have defined meaning; other sequences implementation-defined. |

### Examples

#### Double-quoted argument

**Input:**
Command string: `claude -p --model "claude-3-5-sonnet"`.

**Expected output:**
Parsed argv: `["claude", "-p", "--model", "claude-3-5-sonnet"]`. Four elements. No shell invoked; exec with this argv.

**Verification:**
- Process spawned with exactly four arguments. The model value is one argument.

#### Single-quoted with space

**Input:**
Command string: `my-cli --prompt 'Hello world'`.

**Expected output:**
Parsed argv: `["my-cli", "--prompt", "Hello world"]`. Three elements.

**Verification:**
- The third argument is the single string `Hello world` (one argument, not two).

#### No shell interpretation

**Input:**
Command string: `echo $PATH | cat`.

**Expected output:**
Parsed argv: `["echo", "$PATH", "|", "cat"]` (or equivalent — four arguments). The process executed is `echo` with arguments `$PATH`, `|`, `cat`. No pipe is created; no variable expansion.

**Verification:**
- If Ralph exec's `echo` directly, the child prints literal `$PATH`, `|`, `cat` (or the binary fails). Ralph does not invoke sh.

## Acceptance criteria

- [ ] Command strings with double-quoted arguments are parsed correctly (e.g., `--model "claude-3-5-sonnet"` produces two tokens: `--model` and `claude-3-5-sonnet`)
- [ ] Command strings with single-quoted arguments are parsed correctly
- [ ] Escaped characters within quoted strings are handled
- [ ] The parsed command is executed directly via exec, not through a shell
- [ ] Shell features (pipes, redirects, glob expansion, variable substitution) are not interpreted and do not cause silent misbehavior

## Dependencies

_None identified._
