# R6: Output Buffer Management

**Outcome:** O1 — Iterative Completion

## Requirement

The system captures AI CLI output into a bounded buffer, discarding the oldest content when the buffer is full to preserve the most recent output for signal scanning. The buffer is per-iteration — each iteration starts with a fresh buffer.

## Specification

Each iteration allocates a fresh output buffer. The buffer captures all bytes written to stdout and stderr by the AI CLI process into a single unified stream. When the buffer exceeds `max_output_buffer` bytes, the oldest content is discarded to maintain the size constraint.

**Capture behavior:**

- Ralph redirects both stdout and stderr of the AI CLI process to its capture mechanism
- Bytes from stdout and stderr are interleaved in arrival order into a single buffer
- Output is captured as raw bytes with no encoding transformation or line-based processing

**Truncation behavior:**

- The buffer retains the most recent `max_output_buffer` bytes of output
- When new data would cause the total to exceed the limit, content from the beginning of the buffer is discarded
- After the AI CLI process exits, signal scanning (R2) operates on the buffer's current contents — which is the trailing portion of the full output if truncation occurred

**Implication for signal placement:**

Because truncation discards from the beginning, a signal emitted early in the output may be truncated away if the process subsequently produces enough output to overflow the buffer. Signals emitted near the end of output are always preserved. This is by design — the buffer is bounded to prevent memory exhaustion, and the trade-off is that early signals in very large outputs may be lost.

**Configuration:**

- Field: `max_output_buffer`
- Type: positive integer (bytes, ≥ 1)
- Default: `10485760` (10 MiB)

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Output exceeds buffer size | Oldest bytes discarded; signal scanning sees only the last `max_output_buffer` bytes |
| Success signal emitted early, then truncated out by subsequent output | Signal not detected; iteration outcome is no-signal |
| Process produces exactly `max_output_buffer` bytes | No truncation; entire output is in the buffer |
| Process produces zero bytes | Buffer is empty; signal scanning finds nothing; outcome is no-signal |
| stderr and stdout arrive interleaved | Both captured in arrival order; signal scanning operates on the merged stream |
| Signal string spans the truncation boundary (first half discarded, second half retained) | Signal not detected — R2 requires a complete substring match |

### Examples

#### Signal truncated out of buffer

**Input:**
`max_output_buffer: 1024`. AI CLI writes `<promise>SUCCESS</promise>` in the first 100 bytes, then writes 2000 more bytes of output.

**Expected output:**
The buffer contains the last 1024 bytes. The success signal was in the first 100 bytes, which have been discarded. Signal scanning finds no signal. Iteration outcome is no-signal.

**Verification:**
- Ralph does not exit 0
- Loop continues to next iteration

#### Normal output within buffer bounds

**Input:**
`max_output_buffer: 10485760` (10 MiB). AI CLI writes 50 KiB of output ending with `<promise>SUCCESS</promise>`.

**Expected output:**
Entire output fits in the buffer. Signal scanning detects success. Ralph exits 0.

**Verification:**
- Ralph exit code is 0

## Acceptance criteria

- [ ] Output from the AI CLI's stdout and stderr is captured into a single buffer
- [ ] When the buffer exceeds the configured max size, the oldest content is discarded (truncation from the beginning)
- [ ] Signal scanning operates on the buffer contents after the AI CLI process exits
- [ ] The default buffer size is 10MB (10,485,760 bytes)
- [ ] The buffer size is configurable through the standard configuration hierarchy
- [ ] Each iteration starts with a fresh, empty buffer

## Dependencies

_None identified._
