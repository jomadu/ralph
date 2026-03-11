# R003: Run via alias, file, or standard input

**Outcome:** O001 — Iterative Completion

## Requirement

The system supports running the loop by alias, by prompt file path, or by reading the prompt from standard input.

## Detail

The user can start the iterative loop in three ways: (1) by alias — a configured name that resolves to both an AI command and optionally a prompt path; (2) by explicit prompt file path (the documented way to supply a prompt from a file); (3) by supplying the prompt via standard input (e.g. pipe or redirect). Exactly one source is used per run. How the AI command is determined in each case (alias vs default vs documented option) is a matter of configuration; the requirement is that all three entry points are supported so the user can choose the appropriate one.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Run by alias only | Alias resolves to AI command (and optionally prompt path); loop uses that command and prompt source. |
| Run with prompt file path | Prompt is loaded from the file; AI command from alias or default. |
| Run with prompt via standard input (e.g. pipe) | Prompt is read from the chosen input; AI command from alias or default. |
| Both file and standard input specified | One takes precedence; behavior is documented (e.g. file overrides standard input or vice versa). |
| Alias specifies prompt path | That path is used as prompt source when running by alias. |

### Examples

#### Run by alias

**Input:** Config defines alias `build` with a command and prompt path `./prompts/build.md`. The user runs the run command with that alias.

**Expected output:** The system uses the alias's command as the AI command and loads the prompt from the alias's prompt path; the loop runs with that command and prompt.

**Verification:** The loop executes; the AI receives the content of the prompt file; the process invoked is the one specified by the alias.

#### Run by file path

**Input:** The user runs the run command with a prompt file path.

**Expected output:** The system loads the prompt from the file; the loop runs (AI command from config or default).

**Verification:** Prompt content matches the file; the loop runs without requiring an alias that points to this file.

#### Run with prompt via standard input

**Input:** The user supplies the prompt via standard input (e.g. pipe).

**Expected output:** The system reads the prompt from the input once and runs the loop with that content.

**Verification:** The AI receives the supplied content; no prompt file path is required.

## Acceptance criteria

- [ ] The user can start the loop by specifying an alias (the run command with an alias name).
- [ ] The user can start the loop by specifying a prompt file path (the documented way to supply a prompt from a file).
- [ ] The user can start the loop by supplying the prompt via standard input (e.g. pipe or redirect).
- [ ] Exactly one prompt source is used per run; conflict between file and standard input is resolved in a documented way.
- [ ] When running by alias, the alias supplies or implies the AI command and may supply the prompt path.

## Dependencies

None.
