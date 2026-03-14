# R004: Detect success signal and exit zero

**Outcome:** O001 — Iterative Completion

## Requirement

The system detects the configured success signal in the AI output and exits with the documented success code.

## Detail

The user configures a success signal (e.g. a string, regex, or pattern) that the AI is expected to emit when the task is done. The system captures the AI output (per configuration) and scans it for that signal. When the success signal is found in the output of an iteration, the system treats the iteration as successful, reports completion, and exits with the documented success code. Detection is based on the configured rule (e.g. substring match, line match, or regex); the implementation may live in engineering.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Success signal appears in output | The system detects it, reports completion, exits with the documented success code. |
| Success signal does not appear | The system does not exit with success for that iteration; continues or exits per R005/R007. |
| Success and failure signals both present | Precedence applies per R006 (or R008 if AI-interpreted); outcome is either success or failure, not undefined. |
| Output is empty | No success signal; iteration is not successful. |
| Signal is configured as case-sensitive vs case-insensitive | Behavior matches configuration. |

### Examples

#### Success signal found

**Input:** Success signal configured as substring `DONE`. The AI outputs several lines including "Status: DONE". The system runs one iteration.

**Expected output:** The system detects "DONE" in the output, reports completion (user-observable), and exits with the documented success code.

**Verification:** Exit is with the documented success code; user-observable message indicates success or completion.

#### Success signal not present

**Input:** Success signal configured as `DONE`. The AI outputs "Still working..." and exits. The system runs the iteration.

**Expected output:** The system does not exit with the success code; it increments failure count or starts the next iteration per R005/R007.

**Verification:** Exit is not with the success code (unless max iterations or another exit condition applies); no completion message for success.

## Acceptance criteria

- [ ] The system scans the captured AI output for the configured success signal.
- [ ] When the success signal is detected in an iteration's output, the system reports completion and exits with the documented success code.
- [ ] When the success signal is not detected, the system does not exit with the success code on that basis alone (other rules may still cause exit).
- [ ] Detection respects configuration (e.g. literal vs regex, case sensitivity).

## Dependencies

- R006 (and optionally R008) — When both success and failure signals appear, precedence is defined; R004 applies after precedence resolves to success.
