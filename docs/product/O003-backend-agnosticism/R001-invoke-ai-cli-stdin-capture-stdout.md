# R001: Supply Prompt as Input and Capture AI Command Output

**Outcome:** O003 — Backend Agnosticism

## Requirement

The system supplies the assembled prompt as input to the user-chosen AI command (alias or direct) and captures that command's output.

## Detail

The user selects an AI backend either by name (alias defined in config or built-in) or by passing a direct command string. Ralph resolves the alias to the concrete invocation if needed, validates that the command is present or resolvable before starting the loop (or review), then invokes it with the assembled prompt supplied as input and captures all output. The AI command is invoked in the documented way (no shell); users who need pipes, redirects, or shell expansion must use a wrapper script.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Alias not defined in config and not a built-in | Clear error before loop/review starts; alias cannot be resolved |
| Resolved command not found or not available | Clear error before loop/review starts; command not available |
| User passes direct command string on the command line | Ralph parses and uses it; no alias lookup |
| Empty or invalid direct command | Clear error before execution |
| Prompt is empty | System still invokes the command with empty input (behavior is defined; edge case is explicit) |
| Process exits before consuming full input | Captured output is whatever was produced; exit is reported per normal failure handling |

### Examples

#### Alias resolution and execution

**Input:** Config defines alias `my-ai` → `my-cli --non-interactive`; user runs Ralph with `--ai my-ai` and a prompt file.

**Expected output:** Ralph resolves `my-ai` to the corresponding command, checks that the command is available, then runs it with the assembled prompt as input and captures its output.

**Verification:** Run with `--ai my-ai`; the process receives the prompt as input; output is captured for signal scanning.

#### Direct command bypasses alias

**Input:** User runs `ralph run -- "custom-tool --mode batch"` with a prompt file.

**Expected output:** Ralph uses the direct command string; no alias lookup; prompt supplied as input, output captured.

**Verification:** Invoke with a direct command; confirm the same process receives the prompt as input and its output is captured.

#### Missing command fails early

**Input:** User selects an alias that resolves to a command that is not installed, or passes a direct command that is not available.

**Expected output:** Ralph reports a clear error that the command is missing or cannot be resolved, before starting the loop or review.

**Verification:** Configure or pass a missing command; error message is shown and no iteration starts.

## Acceptance criteria

- [ ] The system resolves a user-specified alias (config or built-in) to a concrete command before execution.
- [ ] The system accepts a direct command string (e.g. from the command line) and uses it without alias resolution.
- [ ] The system validates that the resolved or direct command is present/available before starting the loop or review; if not, the user receives a clear error and no iteration starts.
- [ ] The system supplies the assembled prompt as input to the chosen AI command.
- [ ] The system captures the full output of the AI command for use by the loop (e.g. signal scanning) or review.
- [ ] The AI command is invoked in the documented way (no shell).

## Dependencies

None.
