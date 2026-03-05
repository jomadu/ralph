# R4: Fail-Fast on Invalid Prompt Source

**Outcome:** O2 — Configurable Behavior

## Requirement

The system verifies that the prompt source is valid and produces usable content before starting the loop. Missing files, unreadable files, and empty input cause an immediate, clear error rather than a failure on the first iteration. This applies to all prompt input modes: alias, file flag, and stdin.

## Specification

Before the loop starts, Ralph validates that the prompt source is valid and will produce usable content. This check runs after config validation (R3) and before the first iteration. Prompt source validation applies to all three input modes.

**Input modes:**

| Mode | Trigger | Source |
|------|---------|--------|
| Alias | `ralph run <alias>` | File path from `prompts.<alias>.path` in resolved config |
| File flag | `ralph run -f <path>` | File path provided directly via `--file` / `-f` flag |
| Stdin | `echo "..." \| ralph run -` | Standard input stream |

**Alias mode validation:**

1. Look up `<alias>` in the resolved `prompts` map. If the alias is not defined, exit with error: `unknown prompt alias "<alias>"`.
2. Resolve the `path` field. If relative, resolve relative to the working directory.
3. Check that the file exists. If not, exit with error: `prompt file not found: <resolved-path> (alias: <alias>)`.
4. Check that the file is readable (permission check). If not, exit with error: `prompt file not readable: <resolved-path> (alias: <alias>): permission denied`.
5. Read the file. If the file is empty (0 bytes), exit with error: `prompt file is empty: <resolved-path> (alias: <alias>)`.

**File flag mode validation:**

1. Resolve the path. If relative, resolve relative to the working directory.
2. Check that the file exists. If not, exit with error: `prompt file not found: <resolved-path>`.
3. Check that the file is readable. If not, exit with error: `prompt file not readable: <resolved-path>: permission denied`.
4. Read the file. If empty (0 bytes), exit with error: `prompt file is empty: <resolved-path>`.

**Stdin mode validation:**

1. Read all bytes from stdin until EOF.
2. If 0 bytes were read, exit with error: `stdin is empty: no prompt content provided`.

**Error behavior:**

- All prompt source errors exit with code 1.
- Errors are written to stderr.
- The loop does not start.
- Prompt source validation is a separate phase from config validation (R3). Config validation runs first; if it passes, prompt source validation runs next. They do not interleave — the user fixes config errors first, then prompt source errors.

**File reading:**

The prompt file is read once during validation. The contents are held in memory for the duration of the loop. Ralph does not re-read the file between iterations — the prompt content is fixed at startup. This means file changes during loop execution are not picked up.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Alias exists in config but its `path` field points to a nonexistent file | Error: `prompt file not found: <path> (alias: <alias>)`. Exit 1. |
| File path contains `~` (home directory shorthand) | Ralph does not expand `~` — the path is used as-is. If `~` is not resolved by the shell before reaching Ralph, the file will not be found. |
| File path is an absolute path | Used as-is, no resolution relative to working directory |
| File path is a relative path | Resolved relative to the current working directory at startup |
| Prompt file exists but contains only whitespace (spaces, newlines, tabs) | Valid — the file is non-empty (byte count > 0). Content quality is not Ralph's concern. |
| Prompt file is a symlink to a valid file | Ralph follows the symlink. Validation checks the target file's existence and readability. |
| Prompt file is a symlink to a nonexistent file | Error: `prompt file not found: <resolved-path>` (after symlink resolution) |
| Stdin is a pipe that has not yet received data and stays open | Ralph blocks waiting for EOF. This is standard behavior for stdin reading — the user must close the pipe. |
| Both `-f` and a positional alias are provided | CLI parsing determines which mode is active. `-f` takes precedence; the alias argument is ignored. |
| Alias is valid but file becomes unreadable between validation and first iteration | Extremely unlikely race condition. Ralph does not re-check — the file was read into memory during validation. |

### Examples

#### Missing prompt file via alias

**Input:**
Config defines `prompts.build.path: "./prompts/build.md"`. The file `./prompts/build.md` does not exist. User runs `ralph run build`.

**Expected output:**
```
error: prompt file not found: ./prompts/build.md (alias: build)
```
Ralph exits with code 1. No loop executes.

**Verification:**
- Exit code is 1
- Error message names both the file path and the alias
- No iteration output appears

#### Empty stdin

**Input:**
User runs `echo -n "" | ralph run -`.

**Expected output:**
```
error: stdin is empty: no prompt content provided
```
Ralph exits with code 1.

**Verification:**
- Exit code is 1
- Error message is specific to stdin being empty

#### Valid file via flag

**Input:**
User runs `ralph run -f ./my-prompt.md`. File exists, is readable, and contains 200 bytes of text.

**Expected output:**
Validation passes silently. Loop starts using the file contents as the prompt.

**Verification:**
- No error output
- Loop begins executing with the prompt content from `./my-prompt.md`

## Acceptance criteria

- [ ] When `ralph run <alias>` is invoked, Ralph checks that the mapped prompt file exists and is readable before starting the loop
- [ ] If the alias's file does not exist, Ralph exits with an error message naming the missing file and the alias that referenced it
- [ ] When `ralph run -f <path>` is invoked, Ralph checks that the specified file exists and is readable before starting the loop
- [ ] If the file exists but is not readable (permission denied), Ralph exits with a clear error message
- [ ] When reading from stdin, Ralph exits with an error if stdin is empty (zero bytes)
- [ ] All checks happen at startup, not deferred to the first iteration

## Dependencies

_None identified._
