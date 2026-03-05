# R1: Process Crash Recovery

**Outcome:** O1 — Iterative Completion

## Requirement

The system continues loop execution when an AI CLI process crashes or exits non-zero, preserving any output produced before the crash for signal scanning.

A crash does not receive special treatment beyond what the loop already provides. The partial output is scanned for signals using the same rules as a normal exit. If no signal is found, the iteration counts as a no-signal iteration — the consecutive failure counter is reset to zero, and the loop proceeds to the next iteration. A crash counts as one completed iteration toward the max iteration limit.

## Specification

When the AI CLI process exits, Ralph reads the exit code. Two categories exist:

- **Exit code 0:** normal exit
- **Exit code non-zero** (including killed by signal, e.g., SIGSEGV, OOM): crash

The loop logic is identical for both categories:

1. The output buffer contains all bytes written by the process before exit
2. Signal scanning (R2) runs on the buffer contents
3. The iteration outcome (success, failure, or no-signal) drives loop decisions per R2, R4, and R5

Ralph does not distinguish crashes from normal exits in loop control flow. The only difference is observability: a crash is logged at warn level with the exit code. If the process is killed by the OS (e.g., OOM killer sends SIGKILL), Ralph treats this the same as any other non-zero exit — it captures whatever output was buffered before the kill.

This is distinct from interruption handling (R7), where the user explicitly requests termination. In R7, output is discarded and no signal scanning occurs. In a crash, the process exited on its own, so Ralph proceeds with normal scanning.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Process exits with code 1 and output contains the success signal | Iteration outcome is success; loop exits 0. The non-zero exit code does not override signal detection. |
| Process exits with code 137 (SIGKILL) and no output was produced | Iteration outcome is no-signal; consecutive failure counter resets to 0 (R5); loop continues |
| Process exits with code 1 and output contains the failure signal | Iteration outcome is failure; consecutive failure counter increments (R5) |
| Process writes partial output then crashes mid-line | Signal scanning runs on whatever bytes were captured; a signal string that was only partially written before the crash is not detected |
| Process exits non-zero on the last allowed iteration | Crash counts as a completed iteration; if no success signal found, loop exits 2 (exhaustion per R4) |

### Examples

#### Crash with success signal in partial output

**Input:**
AI CLI exits with code 1 after writing `<promise>SUCCESS</promise>` to stdout.

**Expected output:**
Ralph detects the success signal in the buffer, reports success, and exits 0. The non-zero exit code of the child process does not override the signal.

**Verification:**
- Ralph exit code is 0
- Log output shows the iteration was scanned and success was found despite the crash

#### Crash with no output

**Input:**
AI CLI exits with code 139 (SIGSEGV) with zero bytes written to stdout/stderr.

**Expected output:**
Signal scanning finds neither success nor failure. Iteration outcome is no-signal. Consecutive failure counter resets to 0 (R5). Loop continues to next iteration.

**Verification:**
- Ralph does not exit; next iteration starts
- Consecutive failure counter is 0

## Acceptance criteria

- [ ] When the AI CLI process exits with a non-zero exit code, Ralph captures all output written to stdout/stderr before the exit
- [ ] Captured partial output is scanned for success and failure signals using the same logic as a normal exit
- [ ] A crash with no signal in the partial output resets the consecutive failure counter to zero
- [ ] A crash counts as one completed iteration toward the max iteration limit
- [ ] The loop proceeds to the next iteration after a crash (unless max iterations or failure threshold is reached)
- [ ] Ralph does not retry the AI CLI process within the same iteration

## Dependencies

_None identified._
