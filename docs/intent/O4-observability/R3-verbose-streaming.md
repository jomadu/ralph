# R3: Verbose Output Streaming

**Outcome:** O4 — Observability

## Requirement

The system streams the AI CLI's output to the terminal in real time when verbose mode is enabled, while still capturing it in the output buffer for signal scanning. This lets the user watch the AI work without sacrificing loop control. AI output streaming is controlled by the --verbose flag and is independent of log level — --log-level affects Ralph's own operational messages but does not suppress or enable AI output streaming.

## Specification

Ralph invokes the AI CLI as a child process (O3); the child's stdout and stderr are the output streams. Ralph always captures the full output of each iteration into an output buffer so that signal scanning (O1/R2) can run after the process exits. When the user enables verbose mode via `--verbose` or `-v`, Ralph additionally mirrors the child's stdout and stderr to the terminal (Ralph's stderr or stdout, or both — implementation chooses the destination for mirroring as long as the user sees the AI output in real time) as bytes are produced. Mirroring is in addition to capture; it does not replace it.

**Control of streaming:**

Verbose mode (stream AI output to the terminal) is determined by the resolved value of `show_ai_output`, with the following precedence (highest wins):

1. **CLI:** `--verbose` or `-v` forces streaming on for the run, regardless of config or env.
2. **Environment:** `RALPH_LOOP_SHOW_AI_OUTPUT` — when set, parsed as a boolean (e.g. `true`, `1`, `yes` → true; `false`, `0`, `no`, empty → false). See O2/R8 for env var semantics.
3. **Config:** `loop.show_ai_output` — boolean. Can be set globally or per-prompt (O2/R6).
4. **Default:** `false` — when no CLI flag, env var, or config sets it, AI output is not streamed.

When the resolved value is `true`, Ralph mirrors the child's stdout and stderr to the terminal. When `false`, output is captured only (no mirroring). The `--verbose` / `-v` flag overrides config and env for that run.

- **Not controlled by:** `--log-level` (R5). Log level affects only Ralph's own operational messages, not whether AI output is streamed.

**Behavior:**

1. **When streaming is enabled (resolved `show_ai_output` is true or `-v`/`--verbose` used):** For each iteration, while the AI CLI process is running, every byte (or line, if buffered for display) read from the child's stdout and stderr is written to the terminal in real time, in addition to being appended to the iteration's output buffer. After the process exits, the buffer is scanned for signals (O1/R2). Order of stdout vs stderr when interleaved is implementation-defined (e.g., merge in read order or separate streams).
2. **When streaming is disabled:** No bytes from the child are written to the terminal. All output is still captured into the buffer for signal scanning. The user sees only Ralph's operational output (e.g., progress per R6, statistics per R2), subject to log level (R5).

**Invariants:**

- The same output that is scanned for signals is the output that was captured; if streaming is enabled, the streamed bytes and the buffered bytes are the same (streaming is a mirror, not a tee that diverges).
- Streaming does not affect exit code or loop logic (R1, O1).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| `ralph run build` (no -v, default show_ai_output) | AI output not shown; captured for signal scan |
| `ralph run build -v` | AI stdout/stderr mirrored to terminal; also captured; CLI overrides config/env |
| `ralph run build --verbose` | Same as `-v` |
| `loop.show_ai_output: true` in config, no -v | AI output streamed (config enables it) |
| `RALPH_LOOP_SHOW_AI_OUTPUT=true` in env, no -v | AI output streamed (env enables it) |
| `ralph run build --log-level debug` (no -v) | Ralph's debug messages shown; AI output not streamed unless show_ai_output resolved true |
| `ralph run build --log-level warn -v` | AI output streamed; Ralph's info/debug suppressed; warn/error shown |
| `ralph run build --quiet -v` | R5: quiet sets log level to error; -v still streams AI output; Ralph's non-error output suppressed |
| AI CLI produces only stdout | Only stdout mirrored (and captured) |
| AI CLI produces only stderr | Only stderr mirrored (and captured) |
| AI CLI produces both, interleaved | Both streams mirrored; capture contains both; order for display is implementation-defined |
| Child process crashes (O1/R1) | Partial output captured and scanned; if -v, partial output also streamed up to crash |
| Iteration timeout (O1/R3) | Output up to timeout captured and scanned; if -v, output streamed until process is killed |

### Examples

#### Without verbose — no AI output visible

**Input:**
`ralph run build` with no `-v`. Prompt runs two iterations; second iteration emits success signal.

**Expected output:**
User sees Ralph's progress messages (e.g., "Iteration 1/10", "Iteration 2/10") and completion statistics on stderr (R6, R2). User does not see the raw output of the AI CLI (e.g., model text, tool calls).

**Verification:**
- AI-generated content (e.g., success signal text) does not appear in terminal output
- Ralph still exits 0 (signals were found in the buffer)

#### With verbose — AI output streamed and captured

**Input:**
`ralph run build -v`. One iteration; AI writes "Working... <promise>SUCCESS</promise>" to stdout.

**Expected output:**
User sees "Working... <promise>SUCCESS</promise>" (or equivalent) in real time. Ralph then reports completion and exits 0.

**Verification:**
- The same text appears in the terminal during the run
- Ralph exit code is 0
- Success was detected (buffer contained the signal)

#### Log level does not enable streaming

**Input:**
`ralph run build --log-level debug`. No `-v` flag.

**Expected output:**
Ralph's debug-level messages (if any) are shown. AI CLI output is not streamed to the terminal.

**Verification:**
- AI output is not visible
- Ralph may still exit 0 if success signal was in captured output

## Acceptance criteria

- [ ] With --verbose or -v, AI CLI stdout and stderr are mirrored to the terminal as they are produced (CLI overrides config/env)
- [ ] When resolved show_ai_output is true (config, env, or default), AI output is streamed; when false, it is not
- [ ] Output is simultaneously captured in the buffer for signal scanning after the process exits
- [ ] loop.show_ai_output (config) and RALPH_LOOP_SHOW_AI_OUTPUT (env) can enable streaming; default is false
- [ ] --log-level does not affect streaming (e.g., --verbose --log-level warn streams AI output but suppresses Ralph's debug messages)

## Dependencies

- O1/R2 (signal precedence) — signal scanning uses the same buffer that is populated from the child's output; streaming must not alter or bypass that buffer.
- O1/R6 (output buffer management) — capture semantics and buffer contents are defined there; this requirement adds mirroring only.
- O3 — AI CLI as child process with stdout/stderr; Ralph reads those streams.
- R5 (log level control) — log level governs Ralph's messages only; it does not enable or disable AI output streaming.
