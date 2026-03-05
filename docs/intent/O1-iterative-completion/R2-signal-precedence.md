# R2: Signal Precedence Rules

**Outcome:** O1 — Iterative Completion

## Requirement

The system resolves conflicting signals in a single iteration's output using a deterministic precedence rule: failure wins over success. Signal scanning uses substring matching against the full captured output after the AI CLI process exits.

## Specification

After the AI CLI process exits (normally or via crash per R1), Ralph scans the output buffer for signal strings. The scan determines the iteration outcome.

**Algorithm:**

1. Let `output` = the full contents of the output buffer (after any truncation per R6)
2. Let `failure_signal` = the configured failure signal string
3. Let `success_signal` = the configured success signal string
4. `has_failure` = `failure_signal` is a substring of `output`
5. `has_success` = `success_signal` is a substring of `output`
6. Determine iteration outcome:
   - If `has_failure` → **failure** (regardless of `has_success`)
   - Else if `has_success` → **success**
   - Else → **no-signal**

**Matching rules:**

- Substring matching: the signal string may appear anywhere in the output, at any position, on any line
- Case-sensitive: `SUCCESS` does not match `success`
- No regex, no pattern matching, no wildcards — strict byte-level substring comparison
- The signal string must appear in full. A partial match (e.g., signal split across the buffer truncation boundary per R6) does not count
- Multiple occurrences of the same signal have no additional effect

**Iteration outcome effects:**

| Outcome | Loop behavior |
|---------|---------------|
| success | Loop exits. Ralph exits 0. |
| failure | Consecutive failure counter increments (R5). If threshold reached, exit 1. Otherwise, next iteration. |
| no-signal | Consecutive failure counter resets to 0 (R5). Next iteration (unless max iterations reached per R4, then exit 2). |

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Output contains both the success and failure signal strings | Iteration outcome is failure (failure takes precedence) |
| Output is empty (0 bytes) | Neither signal found; outcome is no-signal |
| Signal string appears inside a larger token (e.g., `foo<promise>SUCCESS</promise>bar`) | Signal is detected — substring match has no boundary requirements |
| Signal string appears in stderr, not stdout | Signal is detected — both streams feed the same buffer (R6) |
| Signal string is split across the buffer truncation boundary (beginning discarded per R6) | Signal is not detected — only complete substring matches count |
| Output contains the signal string because the AI echoed the prompt | Signal is detected — Ralph does not distinguish the source of the text. The prompt and preamble should not contain signal strings. |

### Examples

#### Only success signal

**Input:**
Buffer contains `"Task complete. <promise>SUCCESS</promise>\n"`.

**Expected output:**
`has_failure` = false, `has_success` = true → iteration outcome is success. Ralph exits 0.

**Verification:**
- Ralph exit code is 0

#### Both signals present

**Input:**
Buffer contains `"Found error: <promise>FAILURE</promise>\nFixed it: <promise>SUCCESS</promise>\n"`.

**Expected output:**
`has_failure` = true → iteration outcome is failure, regardless of success also being present. Consecutive failure counter increments.

**Verification:**
- Ralph does not exit 0
- Consecutive failure counter increments by 1

#### No signal

**Input:**
Buffer contains `"Still working on the implementation...\n"`.

**Expected output:**
Neither signal found → no-signal. Consecutive failure counter resets to 0. Loop continues to next iteration.

**Verification:**
- Ralph does not exit
- Consecutive failure counter is 0

## Acceptance criteria

- [ ] When only the success signal appears in output, the iteration is treated as success and the loop stops
- [ ] When only the failure signal appears in output, the iteration is treated as failure and the consecutive failure counter is incremented
- [ ] When both success and failure signals appear in the same output, the iteration is treated as failure
- [ ] When neither signal appears in output, the iteration is treated as no-signal — the consecutive failure counter is reset to zero and the loop proceeds to the next iteration
- [ ] Signal scanning uses substring matching — the signal string can appear anywhere in the output
- [ ] Signal scanning occurs after the AI CLI process exits, not during streaming

## Dependencies

_None identified._
