# R009: Process exit without signal

**Outcome:** O001 — Iterative Completion

## Requirement

When the AI process exits successfully (exit code 0) but the last non-empty line of stdout contains neither the configured success nor failure signal, the system treats the iteration as **neutral** (“more work remains”): it does **not** increment the consecutive-failure count, resets the failure streak, and continues to the next iteration until max iterations or a success signal. When the process exits abnormally, times out, returns an invocation error, exits non-zero without success, or emits the failure signal on the last line, the system applies the same failure-threshold logic as R005 where applicable.

## Detail

The loop preamble instructs the model: when completion criteria are met, emit the success signal; when it cannot proceed, emit the failure signal; **when more work remains, emit no signal** so the loop continues. That only works if exit 0 with no signal on the last line is neutral—not counted as consecutive failure.

**Neutral iteration (exit 0, last non-empty line has neither success nor failure signal):** Continue; consecutive-failure count is reset to zero for this streak.

**Counts toward consecutive failures:** failure signal on the last non-empty line; process exit code non-zero without a success signal on that line; invocation errors (timeout, crash before completion, exec failure).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Process exits 0; last line has no success/failure signal | Neutral; failure streak reset; next iteration. |
| Process exits 0; failure signal on last line | Failure; increment count per R005. |
| Process exits non-zero; no success on last line | Failure; increment count. |
| Process crashes, killed, timeout, invoke error | Failure; increment count. |
| Empty stdout (no non-empty line) with exit 0 | Neutral (no signals detected). |

### Examples

#### Several neutral iterations then success

**Input:** Failure threshold = 3. Iterations 1–2: exit 0, output “Still working…”. Iteration 3: success signal on last line.

**Expected output:** Run completes successfully after 3 iterations; neutral iterations did not consume the failure budget.

#### Non-zero exit without success

**Input:** AI CLI exits 1 with last line “error”. Failure threshold = 2. Two such iterations in a row.

**Expected output:** Exit with failure-threshold code after 2 iterations.

## Acceptance criteria

- [ ] Exit 0 with neither success nor failure on the last non-empty line does not increment consecutive failures and resets the failure streak.
- [ ] Failure signal on the last line, non-zero exit without success, and invocation errors increment consecutive failures per R005 threshold rules.
- [ ] Where documented, the user can distinguish failure-threshold exits (failure signal vs non-zero exit vs invocation error) in reporting.

## Dependencies

- R005 — Failure signal and threshold; R009 refines “no signal on last line” when process exits cleanly.
