# R004: Detect success signal and exit zero

**Outcome:** O001 — Iterative Completion

## Requirement

The system detects the configured success signal in the AI output and exits with the documented success code.

## Detail

The user configures a success signal (e.g. a string, regex, or pattern) that the AI is expected to emit when the task is done. The system captures the AI output (per configuration) and scans **only the last non-empty line** of that output for the configured success signal. The exact definition of "last non-empty line" is given in the [run-loop component spec](../../engineering/components/run-loop.md) (split on newline, trim each line, take the last line that is non-empty after trim; if there is no non-empty line, no success signal is detected). When the success signal is found on that line, the system treats the iteration as successful, reports completion, and exits with the documented success code. Detection on that line is based on the configured rule (e.g. substring match, line match, or regex); the implementation may live in engineering.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Success signal appears on the last non-empty line | The system detects it, reports completion, exits with the documented success code. |
| Success signal appears only on an earlier line (not the last non-empty line) | The system does not treat the iteration as success; continues or exits per R005/R007. |
| Success signal does not appear | The system does not exit with success for that iteration; continues or exits per R005/R007. |
| Success and failure signals both present (on that line) | Precedence applies per R006 (or R008 if AI-interpreted); outcome is either success or failure, not undefined. |
| Output is empty (no non-empty line) | No success signal; iteration is not successful. |
| Signal is configured as case-sensitive vs case-insensitive | Behavior matches configuration. |

### Examples

#### Success signal found

**Input:** Success signal configured as substring `DONE`. The AI outputs several lines with "Status: DONE" on the **last non-empty line** (e.g. "Still working...\nStatus: DONE"). The system runs one iteration.

**Expected output:** The system detects "DONE" on the last non-empty line, reports completion (user-observable), and exits with the documented success code.

**Verification:** Exit is with the documented success code; user-observable message indicates success or completion.

#### Success signal not present

**Input:** Success signal configured as `DONE`. The AI outputs "Still working..." and exits. The system runs the iteration.

**Expected output:** The system does not exit with the success code; it increments failure count or starts the next iteration per R005/R007.

**Verification:** Exit is not with the success code (unless max iterations or another exit condition applies); no completion message for success.

#### Success signal only on an earlier line

**Input:** Success signal configured as `DONE`. The AI outputs "Status: DONE\nStill working..." (so "DONE" appears only on the first line; the last non-empty line is "Still working...").

**Expected output:** The system does not treat the iteration as success (the last non-empty line is scanned, and it does not contain the signal); it continues or exits per R005/R007.

**Verification:** No success exit; no completion message for success.

## Acceptance criteria

- [ ] The system scans only the last non-empty line of the captured AI output for the configured success signal; earlier lines are not used for success detection (see run-loop component spec for definition of last non-empty line).
- [ ] When the success signal is detected on that line in an iteration's output, the system reports completion and exits with the documented success code.
- [ ] When the success signal is not detected (including when it appears only on an earlier line), the system does not exit with the success code on that basis alone (other rules may still cause exit).
- [ ] Detection respects configuration (e.g. literal vs regex, case sensitivity).

## Dependencies

- R006 (and optionally R008) — When both success and failure signals appear, precedence is defined; R004 applies after precedence resolves to success.
