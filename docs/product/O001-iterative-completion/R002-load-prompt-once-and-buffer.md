# R002: Load prompt once and buffer

**Outcome:** O001 — Iterative Completion

## Requirement

The system loads the prompt once from the chosen source (file or standard input) and buffers it for all iterations, or fails before the loop with a clear message if the source is unavailable.

## Detail

At loop start, the system obtains the prompt from exactly one source: alias-defined prompt path, explicit file path (the documented way to supply a prompt file), or standard input. The prompt is read once and held in memory for the duration of the loop. All iterations use this buffered content (with any per-iteration preamble applied in memory). If the source cannot be read (file missing, permission denied, standard input closed or empty when expected), the system fails before starting the first iteration and reports the reason clearly.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| File path given; file does not exist | Fail before loop; clear message (e.g. file not found). |
| File path given; file not readable (permissions) | Fail before loop; clear message (e.g. permission denied). |
| File path given; file is empty | Buffer is empty; loop may start with empty prompt (or treat as config error per product policy). |
| Standard input chosen but closed or empty | Fail before loop or treat as empty buffer; behavior must be documented. |
| Standard input chosen; user supplies content (e.g. pipe) | Read once, buffer; use for all iterations. |
| Prompt loaded successfully | Buffer used for all iterations; file is not re-read between iterations. |

### Examples

#### File path, file exists and is readable

**Input:** The user runs the run command with a prompt from a file; the file exists and is readable.

**Expected output:** The system reads the file once, buffers its content, starts the loop, and uses the buffer for every iteration.

**Verification:** The loop runs; modifying the file on disk during the loop does not change what the AI receives in subsequent iterations.

#### File path, file missing

**Input:** The user runs the run command with a prompt file path; the file does not exist.

**Expected output:** The system exits with a documented non-success code before the first iteration. The message indicates the prompt source is unavailable (e.g. file not found).

**Verification:** Exit is with a non-success code; no AI process spawned for the main task; the message is about the prompt or source, not the AI command.

#### Input supplied via standard input, use for all iterations

**Input:** The user supplies the prompt via standard input (e.g. pipe).

**Expected output:** The system reads the input once, buffers the content, and uses that buffer for all iterations.

**Verification:** Only one read of the input source; multiple iterations run with the same prompt content.

## Acceptance criteria

- [ ] When the prompt source is a file and the file is missing or unreadable, the system fails before starting the loop with a clear message identifying the prompt source problem.
- [ ] When the prompt source is valid (file or standard input), the system reads it once and buffers the content.
- [ ] The buffered prompt is used for every iteration; the source is not re-read between iterations.
- [ ] Per-iteration preamble (e.g. iteration count) may be applied in memory to the buffered content; the underlying prompt content is unchanged for the duration of the loop.

## Dependencies

- R003 — Defines how the prompt source is chosen (alias, file path, standard input).
