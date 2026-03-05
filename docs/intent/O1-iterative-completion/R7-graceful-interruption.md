# R7: Graceful Interruption Handling

**Outcome:** O1 — Iterative Completion

## Requirement

The system handles SIGINT and SIGTERM by terminating the current AI CLI process, waiting for it to exit, and then exiting with code 130. The interruption is clean — no partial iteration results are processed.

## Specification

Ralph installs handlers for SIGINT and SIGTERM at loop startup. The handlers coordinate clean shutdown by terminating the child process and exiting with code 130.

**Signal handling during AI CLI execution (child process running):**

1. First signal (SIGINT or SIGTERM) received:
   a. Set a flag indicating interruption is in progress
   b. Forward SIGTERM to the child process
   c. Start a 5-second grace period
2. If the child exits within the grace period:
   a. Discard the iteration — do not scan output for signals
   b. Exit 130
3. If the grace period expires and the child has not exited:
   a. Send SIGKILL to the child process
   b. Wait for the child to exit (SIGKILL cannot be caught)
   c. Exit 130
4. Second signal received during the grace period:
   a. Send SIGKILL to the child immediately
   b. Exit 130

**Signal handling between iterations (no child process running):**

Signal received → exit 130 immediately. No cleanup is needed.

**Key behaviors:**

- No signal scanning occurs on an interrupted iteration. The output buffer is discarded. This is the critical difference from crash recovery (R1), where the process exited on its own and output is always scanned.
- The 5-second SIGTERM → SIGKILL escalation matches the same pattern used in iteration timeout (R3). Both use the same grace period duration.
- Exit code 130 follows Unix convention (128 + signal number for SIGINT = 2).
- SIGINT and SIGTERM are handled identically. The signal type does not affect behavior.

**Interaction with iteration timeout (R3):**

If an iteration timeout fires and starts its own SIGTERM → grace period → SIGKILL sequence, and then a user signal (SIGINT/SIGTERM) arrives during that grace period, interruption handling takes over. The child is terminated, output is discarded (not scanned), and Ralph exits 130. The timeout's intent was to continue the loop; the user's intent is to stop entirely. The user's intent wins.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| SIGINT during AI CLI execution | SIGTERM forwarded to child; 5-second grace period; exit 130 |
| SIGTERM during AI CLI execution | Same as SIGINT — SIGTERM forwarded to child; grace period; exit 130 |
| Two rapid SIGINTs | First starts grace period; second sends SIGKILL immediately; exit 130 |
| SIGINT between iterations | No child to terminate; exit 130 immediately |
| SIGINT during timeout grace period (R3) | Interruption takes over; output discarded; exit 130 |
| Child exits on its own during interruption grace period with success signal in output | Output is not scanned (interruption discards the iteration); exit 130 |

### Examples

#### User presses Ctrl-C during iteration

**Input:**
AI CLI is running on iteration 3 of 10. User presses Ctrl-C (SIGINT).

**Expected output:**
Ralph sends SIGTERM to the AI CLI process. AI CLI exits within 5 seconds. Ralph exits with code 130. No iteration result is processed for iteration 3.

**Verification:**
- Ralph exit code is 130
- No success/failure message is printed for iteration 3

#### User presses Ctrl-C between iterations

**Input:**
Iteration 2 just completed (no-signal). Before iteration 3 starts, user presses Ctrl-C.

**Expected output:**
No child process is running. Ralph exits immediately with code 130.

**Verification:**
- Ralph exit code is 130
- Only 2 iterations executed

## Acceptance criteria

- [ ] On SIGINT or SIGTERM, Ralph sends a termination signal to the running AI CLI process
- [ ] Ralph waits for the AI CLI process to exit with a bounded timeout before forcing termination
- [ ] Ralph exits with code 130 after handling the interruption
- [ ] If no AI CLI process is running at the time of the signal (e.g., between iterations), Ralph exits immediately with code 130
- [ ] A second SIGINT/SIGTERM during the wait forces immediate exit

## Dependencies

_None identified._
