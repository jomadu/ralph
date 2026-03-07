# R3: Verbose Output Streaming

**Outcome:** O4 — Observability

## Requirement

The system streams the AI CLI's output to the terminal in real time when "show AI command output" is enabled, while still capturing it in the output buffer for signal scanning. This lets the user watch the AI work without sacrificing loop control. AI output streaming is controlled by the resolved value of **show AI command output** (default **true**). It can be turned off via `--no-ai-cmd-output` or `-q` (R5). It is independent of log level — `--log-level` affects Ralph's own operational messages but does not suppress or enable AI output streaming.

## Specification

Ralph invokes the AI CLI as a child process (O3); the child's stdout and stderr are the output streams. Ralph always captures the full output of each iteration into an output buffer so that signal scanning (O1/R2) can run after the process exits. When the resolved value of **show AI command output** is **true**, Ralph additionally mirrors the child's stdout and stderr to **stdout** (the run's log; see O4 README Non-outcomes) as bytes are produced. Mirroring is in addition to capture; it does not replace it.

**Control of streaming:**

Show AI command output (stream AI output to the terminal) is determined by the resolved value of `show_ai_output`, with the following precedence (highest wins):

1. **CLI:** `--no-ai-cmd-output` — forces streaming **off** for the run (output captured only, not mirrored). Overrides config, env, and other flags.
2. **CLI:** `-q` / `--quiet` (R5) — also sets show AI command output to **false** for the run. So quiet mode suppresses both Ralph's log verbosity and AI output.
3. **CLI:** `-v` / `--verbose` — sets show AI command output to **true** for the run (unless overridden by `--no-ai-cmd-output`). Overrides config and env when present.
4. **Environment:** `RALPH_LOOP_SHOW_AI_OUTPUT` — when set, parsed as a boolean (e.g. `true`, `1`, `yes` → true; `false`, `0`, `no`, empty → false). See O2/R8 for env var semantics.
5. **Config:** `loop.show_ai_output` — boolean. Can be set globally or per-prompt (O2/R6).
6. **Default:** `true` — when no CLI flag, env var, or config sets it, AI output **is** streamed to stdout.

When the resolved value is `true`, Ralph mirrors the child's stdout and stderr to **stdout**. When `false`, output is captured only (no mirroring).

- **Not controlled by:** `--log-level` (R5). Log level affects only Ralph's own operational messages, not whether AI output is streamed.

**Behavior:**

1. **When streaming is enabled (resolved `show_ai_output` is true):** For each iteration, while the AI CLI process is running, every byte (or line, if buffered for display) read from the child's stdout and stderr is written to **stdout** in real time, in addition to being appended to the iteration's output buffer. After the process exits, the buffer is scanned for signals (O1/R2). Order of stdout vs stderr when interleaved is implementation-defined (e.g., merge in read order or separate streams).
2. **When streaming is disabled:** No bytes from the child are written to the terminal. All output is still captured into the buffer for signal scanning. The user sees only Ralph's operational output (e.g., progress per R6, statistics per R2), subject to log level (R5).

**Invariants:**

- The same output that is scanned for signals is the output that was captured; if streaming is enabled, the streamed bytes and the buffered bytes are the same (streaming is a mirror, not a tee that diverges).
- Streaming does not affect exit code or loop logic (R1, O1).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| `ralph run build` (default) | AI output **streamed** to stdout (default show_ai_output = true); also captured for signal scan |
| `ralph run build -q` | AI output **not** streamed (-q sets show AI command output to false); captured only |
| `ralph run build --no-ai-cmd-output` | AI output not streamed; captured only (explicit off) |
| `ralph run build -v` | AI stdout/stderr mirrored to stdout; also captured; -v forces true |
| `ralph run build --verbose` | Same as `-v` |
| `ralph run build -v --no-ai-cmd-output` | Verbose logs (R5), but AI output **not** streamed; --no-ai-cmd-output wins |
| `ralph run build -q -v` or `-v -q` | -q wins for show AI output → false; log level from -q (error) unless --log-level overrides |
| `loop.show_ai_output: true` in config, no flags | AI output streamed (config or default) |
| `RALPH_LOOP_SHOW_AI_OUTPUT=true` in env, no flags | AI output streamed (env enables it) |
| `ralph run build --log-level debug` (no -v) | Ralph's debug messages shown; AI output streamed by default |
| `ralph run build --log-level warn -v` | AI output streamed; Ralph's info/debug suppressed; warn/error shown |
| AI CLI produces only stdout | Only stdout mirrored to stdout (and captured) |
| AI CLI produces only stderr | Only stderr mirrored to stdout (and captured) |
| AI CLI produces both, interleaved | Both streams mirrored to stdout; capture contains both; order for display is implementation-defined |
| Child process crashes (O1/R1) | Partial output captured and scanned; if streaming on, partial output also streamed up to crash |
| Iteration timeout (O1/R3) | Output up to timeout captured and scanned; if streaming on, output streamed until process is killed |

### Examples

#### Default — AI output streamed to stdout

**Input:**
`ralph run build` (no flags). Prompt runs two iterations; second iteration emits success signal.

**Expected output:**
User sees Ralph's progress messages (e.g., "Iteration 1/10", "Iteration 2/10") and completion statistics on stdout (R6, R2; see O4 README). User also sees the raw output of the AI CLI streamed to stdout (default show AI command output = true).

**Verification:**
- AI-generated content (e.g., success signal text) appears in terminal output
- Ralph exits 0 (signals were found in the buffer)

#### Quiet — no AI output visible

**Input:**
`ralph run build -q`. Prompt runs two iterations; second iteration emits success signal.

**Expected output:**
User does not see the raw output of the AI CLI (-q sets show AI command output to false). Ralph's non-error messages are also suppressed (R5). User sees only errors if any, and exit code 0.

**Verification:**
- AI-generated content does not appear in terminal output
- Ralph still exits 0 (signals were found in the buffer)

#### With verbose — AI output streamed and captured

**Input:**
`ralph run build -v`. One iteration; AI writes "Working... <promise>SUCCESS</promise>" to stdout.

**Expected output:**
User sees "Working... <promise>SUCCESS</promise>" (or equivalent) in real time on stdout. Ralph then reports completion and exits 0.

**Verification:**
- The same text appears in the terminal during the run
- Ralph exit code is 0
- Success was detected (buffer contained the signal)

#### Verbose logs but no AI stream

**Input:**
`ralph run build -v --no-ai-cmd-output`.

**Expected output:**
Ralph's log level is verbose (debug). AI CLI output is **not** streamed; `--no-ai-cmd-output` wins. User sees Ralph's operational messages (including debug) but not the AI's stdout/stderr in real time.

**Verification:**
- Ralph debug/info messages visible
- AI output not visible in terminal (still captured for signal scan)

## Acceptance criteria

- [ ] Default show AI command output is **true** — AI output is streamed to stdout when no flags/config override
- [ ] `--no-ai-cmd-output` sets show AI command output to false (no mirroring; capture only)
- [ ] `-q` / `--quiet` (R5) also sets show AI command output to false
- [ ] With `-v` / `--verbose`, show AI command output is true unless `--no-ai-cmd-output` is set (CLI overrides config/env)
- [ ] Precedence: `--no-ai-cmd-output` → false; `-q` → false; `-v` → true; config/env; default true
- [ ] AI output is mirrored to **stdout** (run log), not stderr
- [ ] Output is simultaneously captured in the buffer for signal scanning after the process exits
- [ ] loop.show_ai_output (config) and RALPH_LOOP_SHOW_AI_OUTPUT (env) can enable or disable streaming; default is true
- [ ] --log-level does not affect streaming (e.g., -v --log-level warn streams AI output but suppresses Ralph's info/debug messages)
- [ ] Edge case: `-v --no-ai-cmd-output` yields verbose Ralph logs but no AI stream

## Dependencies

- O1/R2 (signal precedence) — signal scanning uses the same buffer that is populated from the child's output; streaming must not alter or bypass that buffer.
- O1/R6 (output buffer management) — capture semantics and buffer contents are defined there; this requirement adds mirroring only.
- O3 — AI CLI as child process with stdout/stderr; Ralph reads those streams.
- R5 (log level control) — log level governs Ralph's messages only; it does not enable or disable AI output streaming.
