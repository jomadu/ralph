# R3: Per-Iteration Timeout

**Outcome:** O1 — Iterative Completion

## Requirement

The system enforces a configurable time limit on each AI CLI process invocation, terminating processes that exceed the limit. The timeout applies independently to each iteration, not to the total loop duration.

## Specification

When `iteration_timeout` is set to a positive integer, Ralph enforces a per-iteration time limit on the AI CLI process.

**Mechanism:**

1. Spawn the AI CLI process
2. Start a timer for `iteration_timeout` seconds
3. If the process exits before the timer fires: cancel the timer, proceed normally
4. If the timer fires before the process exits:
   a. Send SIGTERM to the child process
   b. Wait up to 5 seconds for the process to exit gracefully
   c. If the process has not exited after 5 seconds, send SIGKILL
   d. Wait for the process to exit (SIGKILL cannot be caught)

After a timed-out process is terminated, it exits with a non-zero code. Per R1, this is logged as a crash, but the post-exit flow is the same as any iteration: the output buffer is scanned for signals (R2), and the iteration outcome drives loop logic (R4, R5). The timeout itself adds no special handling beyond killing the process.

**Configuration:**

- Field: `iteration_timeout`
- Type: non-negative integer (seconds)
- Default: `0` (no timeout — process runs until it exits, is interrupted per R7, or the system kills it)
- A value of `0` means no timer is started

**Interaction with interruption handling (R7):**

If a SIGINT/SIGTERM arrives from the user during the timeout's 5-second grace period, interruption handling (R7) takes over. The child is terminated per R7's rules, output is discarded (not scanned), and Ralph exits 130.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Process exits 1 second before timeout | Timer is cancelled; normal iteration processing |
| Process exceeds timeout and exits on SIGTERM | Output captured up to that point is scanned for signals; normal loop logic applies |
| Process exceeds timeout and ignores SIGTERM | SIGKILL sent after 5-second grace period; process forcefully killed |
| Timeout set to 0 | No timer is started; process runs until natural exit, interruption (R7), or system-level termination |
| Process writes success signal then hangs past timeout | Timeout fires, process is killed, but the success signal is already in the buffer; iteration outcome is success |
| SIGINT arrives during the timeout's 5-second grace period | Interruption handling (R7) takes over; exit 130 |

### Examples

#### Process times out with partial output

**Input:**
`iteration_timeout: 30`. AI CLI runs for 45 seconds, producing output for the first 25 seconds, then hangs.

**Expected output:**
At 30 seconds, Ralph sends SIGTERM. The process exits. Signal scanning runs on the output captured during the first 25 seconds. If no signal is found, iteration outcome is no-signal; loop continues to the next iteration.

**Verification:**
- Iteration completes in approximately 30 seconds (plus up to 5 seconds if SIGTERM is slow)
- Next iteration starts
- Log output indicates the timeout was enforced

#### Process completes within timeout

**Input:**
`iteration_timeout: 60`. AI CLI runs for 15 seconds and exits normally.

**Expected output:**
Timer is cancelled. Iteration proceeds through normal signal scanning and loop logic. The timeout has no effect on this iteration.

**Verification:**
- No timeout-related log messages
- Iteration completes in approximately 15 seconds

## Acceptance criteria

- [ ] When iteration_timeout is set to a positive value, Ralph kills the AI CLI process if it runs longer than the specified duration in seconds
- [ ] Partial output from a timed-out process is captured and scanned for signals
- [ ] A timed-out iteration counts as one completed iteration toward the max iteration limit
- [ ] When iteration_timeout is 0 or unset, no time limit is enforced
- [ ] The timeout applies to each iteration independently — a 60-second timeout means each iteration gets 60 seconds, not 60 seconds total

## Dependencies

_None identified._
