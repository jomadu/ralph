# R001: Validate AI command before loop

**Outcome:** O001 — Iterative Completion

## Requirement

The system validates that the chosen AI command or alias is available before starting the loop.

## Detail

The system resolves the user's chosen command (e.g. an alias from config or a direct command name) and checks that it is available and executable. If the chosen AI command is missing, unavailable, or not executable, the system fails before starting any AI process and before loading the prompt. The user receives a clear, actionable message (e.g. command not found, alias unresolved, or command not executable).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Alias not defined in config | Fail before loop; message indicates alias is missing or invalid. |
| Alias points to non-existent path | Fail before loop; message indicates command/path not found. |
| Binary exists but is not executable | Fail before loop; message indicates not executable. |
| Command available and executable | Validation passes; loop may start (subject to R002). |
| Empty or whitespace-only command | Fail before loop; message indicates invalid command. |

### Examples

#### Alias resolves to available binary

**Input:** Config defines alias `agent` → a command that is available and executable. The user runs the run command with that alias.

**Expected output:** Validation succeeds; the system proceeds to load the prompt and start the loop (or fails at R002 if the prompt source is invalid).

**Verification:** No process spawn error; either the loop starts or a later failure is due to prompt or signal, not the command.

#### Alias points to missing or invalid command

**Input:** Config defines alias `agent` → a path or command that does not exist or is not executable. The user runs the run command with that alias.

**Expected output:** The system exits with a documented non-success code before starting the loop. The user sees a clear message that the command is missing or invalid.

**Verification:** Exit is with a non-success code; the message mentions the command, alias, or path; no AI process was spawned.

## Acceptance criteria

- [ ] When the chosen AI command (alias or direct) is not available (missing command, invalid alias, not executable), the system exits with a documented non-success code before starting the loop.
- [ ] The failure message clearly indicates that the problem is the AI command (e.g. alias not found, command unavailable, or not executable).
- [ ] No AI process is spawned when validation fails.
- [ ] When the command is available and executable, validation passes and the system proceeds to prompt loading and loop (or fails at a later step with a distinct message).

## Dependencies

None.
