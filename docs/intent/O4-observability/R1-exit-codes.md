# R1: Distinct Exit Codes

**Outcome:** O4 — Observability

## Requirement

The system exits with distinct codes for each termination reason, enabling scripts and CI systems to programmatically determine why the loop stopped without parsing output.

## Specification

Ralph exits with exactly one of four exit codes. The code is determined solely by the termination reason. Exit code selection happens after the loop has decided to stop (success, failure threshold, max iterations) or when the user interrupts (O1/R7). No other factors — verbosity, log level, dry-run, or output content — affect the exit code.

**Termination reason → exit code mapping:**

| Termination reason | Exit code | When it occurs |
|--------------------|-----------|----------------|
| Success signal received | 0 | O1/R2: signal scan finds success (and no failure); loop stops |
| Failure threshold reached or abort | 1 | O1/R5: consecutive failure count reaches threshold; loop stops |
| Max iterations exhausted | 2 | O1/R4: iteration count reaches limit without success signal; loop stops |
| Interrupted (SIGINT or SIGTERM) | 130 | O1/R7: user signal received; child terminated or none running; no signal scan on interrupted iteration |

**Precedence:** Only one termination reason applies at a time. Success is decided per iteration (O1/R2); failure threshold and max iterations are evaluated after the iteration outcome. Interruption (O1/R7) takes precedence over normal completion — if SIGINT/SIGTERM is received, Ralph exits 130 regardless of what would have happened otherwise.

**Invariants:**

- No two termination reasons share an exit code.
- Exit code is identical whether `--verbose` is set or not, and regardless of `--log-level` or `--quiet`.
- Dry-run (R4) is a separate path: no loop runs, no termination reason from the loop; dry-run success is exit 0 by specification of R4.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Success on iteration 1 | Exit 0 |
| Success on iteration N (N > 1) | Exit 0 |
| Failure threshold reached on iteration 5 with threshold 3 | Exit 1 |
| Max iterations (e.g., 10) reached with no success signal | Exit 2 |
| SIGINT during AI CLI execution | Exit 130 (O1/R7); no success/failure/max-iter evaluation |
| SIGINT between iterations | Exit 130 immediately |
| Dry-run with valid config | Exit 0 (R4); no loop termination reason |
| `--verbose` or `-v` set | Exit code unchanged vs same run without verbose |
| `--log-level error` or `--quiet` | Exit code unchanged vs default log level |
| Both failure and success signals in output (O1/R2) | Iteration is failure; if threshold then reached → exit 1 |

### Examples

#### Success signal on first iteration

**Input:**
`ralph run build` with default config. Prompt runs one iteration; AI output contains `<promise>SUCCESS</promise>`.

**Expected output:**
Ralph detects success, stops the loop, exits with code 0.

**Verification:**
- `echo $?` (or equivalent) after Ralph exits is 0

#### Failure threshold reached

**Input:**
`ralph run build` with `failure_threshold: 3`. Three consecutive iterations produce output containing `<promise>FAILURE</promise>`.

**Expected output:**
After the third failure, Ralph exits with code 1.

**Verification:**
- Exit code is 1
- No fourth iteration is started

#### Max iterations exhausted

**Input:**
`ralph run build -n 5`. Five iterations run; no iteration output contains the success signal.

**Expected output:**
After the fifth iteration, Ralph exits with code 2.

**Verification:**
- Exit code is 2
- Exactly 5 iterations executed

#### Interrupted by SIGINT

**Input:**
`ralph run build`. User presses Ctrl-C during the second iteration while the AI CLI is running.

**Expected output:**
Ralph sends SIGTERM to the child (O1/R7), then exits with code 130. Iteration 2 is not counted as success or failure for loop logic.

**Verification:**
- Exit code is 130

## Acceptance criteria

- [ ] Exit code 0: success signal received
- [ ] Exit code 1: failure threshold reached or explicit abort
- [ ] Exit code 2: max iterations exhausted without a success signal
- [ ] Exit code 130: interrupted by SIGINT or SIGTERM
- [ ] No two termination reasons share an exit code
- [ ] Exit codes are consistent regardless of verbosity, log level, or other output settings

## Dependencies

- O1/R2 (signal precedence) — defines when "success" and "failure" iteration outcomes occur, which drive exit 0 and exit 1.
- O1/R4 (max iteration limit) — defines when "max iterations exhausted" occurs, which drives exit 2.
- O1/R5 (consecutive failure tracking) — defines when "failure threshold reached" occurs, which drives exit 1.
- O1/R7 (graceful interruption) — defines when Ralph exits 130 and that no signal scanning occurs on interrupt.
