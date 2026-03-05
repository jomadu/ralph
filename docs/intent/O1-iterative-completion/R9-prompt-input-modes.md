# R9: Prompt Input Modes

**Outcome:** O1 — Iterative Completion

## Requirement

The system accepts prompts from multiple sources — a configured alias, a direct file path flag, or stdin — loading the prompt content once at loop start and reusing it immutably across all iterations. This ensures consistent behavior regardless of filesystem changes during execution and supports one-off usage without config file setup.

## Specification

Ralph supports three mutually exclusive prompt input modes. The mode is determined at startup based on the command-line arguments and stdin state.

**Mode resolution (evaluated in order):**

1. If a positional argument `<alias>` is provided → **alias mode**
2. Else if `--file <path>` or `-f <path>` is provided → **file mode**
3. Else if stdin is not a TTY (piped input detected) → **stdin mode**
4. Else → **error** — no prompt source identified

Providing both a positional alias and `--file` is an error. Ralph exits with a clear message identifying the conflict.

**Mode behaviors:**

| Mode | Invocation | Source resolution |
|------|------------|-------------------|
| Alias | `ralph run <alias>` | Look up `prompts.<alias>.path` in resolved config. Read the file at that path. |
| File | `ralph run -f <path>` | Read the file at the specified path directly. No alias lookup. |
| Stdin | `cat prompt.md \| ralph run` | Read all bytes from stdin until EOF. |

**Read-once semantics:**

In all modes, the prompt content is read exactly once at loop start:

1. Resolve the source (alias → config path → file, or `-f` → file, or stdin → stream)
2. Read the entire content into a memory buffer
3. Store the buffer as the immutable prompt content for the loop
4. All iterations use this same buffer — the source is never re-read

For stdin mode, stdin is fully consumed at startup. Subsequent iterations do not attempt to read stdin again.

**Prompt content is treated as opaque bytes.** Ralph does not parse, validate, or transform the content. It is piped as-is (after optional preamble wrapping per R8) to the AI CLI's stdin.

**Validation:**

Prompt source validation — missing file, unreadable file, empty content — is handled by O2/R4 (Fail-fast on invalid prompt source). This requirement covers only mode resolution and read-once buffering.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| `ralph run build -f ./prompt.md` | Error: alias and file flag are mutually exclusive |
| `ralph run` with stdin as a TTY (interactive terminal) | Error: no prompt source identified |
| `ralph run -f ./prompt.md` with stdin also piped | File mode takes precedence; piped stdin data is ignored |
| Alias resolves to a path that doesn't exist | Handled by O2/R4 (fail-fast validation), not by this requirement |
| Prompt file changes on disk after loop starts | Changes are not seen; the buffered content from the initial read is used for all iterations |
| Prompt content is binary (not UTF-8 text) | Ralph does not interpret content; bytes are piped as-is |
| Stdin contains zero bytes (immediate EOF) | Handled by O2/R4 (fail-fast on empty prompt) |

### Examples

#### Alias mode

**Input:**
Config contains `prompts.build.path: "./prompts/build.md"`. User runs `ralph run build`.

**Expected output:**
Ralph looks up the `build` alias, reads `./prompts/build.md` once, buffers the content, and uses it for all iterations.

**Verification:**
- The file is read once at startup
- Modifying `./prompts/build.md` after startup does not affect subsequent iterations

#### File mode

**Input:**
User runs `ralph run -f ./my-prompt.md`.

**Expected output:**
Ralph reads `./my-prompt.md` once, buffers the content, and uses it for all iterations. No alias lookup is performed.

**Verification:**
- No config file is needed for prompt resolution
- The file is read once at startup

#### Stdin mode

**Input:**
User runs `cat notes.md | ralph run`.

**Expected output:**
Ralph detects that stdin is not a TTY, reads all of stdin until EOF, buffers the content, and uses it for all iterations.

**Verification:**
- stdin is fully consumed before the first iteration begins
- All iterations receive the same content

## Acceptance criteria

- [ ] `ralph run <alias>` reads the prompt file mapped to the alias in the resolved config
- [ ] `ralph run -f <path>` reads the specified file directly, without requiring an alias in config
- [ ] `cat prompt.md | ralph run` reads the prompt from stdin when no alias or file flag is provided
- [ ] In all modes, the prompt content is read once at loop start and buffered in memory
- [ ] The same buffered content is used for every iteration — changes to the prompt file on disk after loop start do not affect subsequent iterations
- [ ] Prompt source validation (missing, unreadable, or empty) is handled by O2/R4 — Fail-fast on invalid prompt source

## Dependencies

_None identified._
