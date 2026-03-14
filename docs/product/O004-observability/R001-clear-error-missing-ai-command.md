# R001: Clear error when AI command is missing or invalid

**Outcome:** O004 — Observability

## Requirement

The system reports a clear, user-facing error when the AI command or alias is missing or invalid before the run or review starts.

## Detail

The user must understand why the run could not start so they can fix config or install the tool. The error is emitted before any loop or review execution begins. "Missing or invalid" includes: alias not found, binary not on PATH, command string that fails to resolve or execute.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Alias not defined in config | Clear error stating the alias is missing or unresolved; no loop/review started |
| Binary not on PATH | Clear error indicating the command cannot be found or executed |
| Command resolves but fails immediately (e.g. bad args) | User sees a clear error; behavior may be reported as startup failure so user can fix invocation |
| Review command uses same AI command resolution | Same clear error behavior before review runs |

### Examples

#### Alias not found

**Input:** Config references an alias that is not defined; user invokes the run command.

**Expected output:** A message to the user that the AI command or alias is missing or invalid, with enough context to identify the alias or command name.

**Verification:** User can read the message and understand they need to define the alias or fix the command; the run does not start.

#### Binary not on PATH

**Input:** Config specifies a command that is not installed or not on PATH; user invokes the run command.

**Expected output:** A clear error that the command could not be found or executed.

**Verification:** User understands the run could not start due to missing/invalid command; no loop execution.

## Acceptance criteria

- [ ] When the AI command or alias is missing or invalid, the system emits a clear, user-facing error message before starting the loop or review.
- [ ] The error message allows the user to identify what is wrong (e.g. which alias or command).
- [ ] The system does not start the run or review when the AI command is missing or invalid.
- [ ] Behavior applies to both run and review commands when they resolve the AI command the same way.
