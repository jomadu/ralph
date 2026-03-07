# R6: Iteration Progress Reporting

**Outcome:** O4 — Observability

## Requirement

The system reports the current iteration number at the start of each iteration so the user knows where the loop is in its run. Progress messages are visible during normal operation and suppressed only when the user explicitly reduces verbosity.

## Specification

At the start of each iteration — before the prompt is assembled and piped to the AI CLI (O1/R8) — Ralph emits a progress message indicating the current iteration number and the iteration limit. The user can see where the loop is in its run without parsing AI output.

**Message content:**

- **Bounded mode:** When a max iteration limit is in effect (e.g., from config or `-n`), the message includes both the current iteration and the limit. Format is implementation-defined; example: `Iteration 3/10` or `Iteration 3 of 10`. The iteration number is 1-based and matches the iteration count used in the preamble (O1/R8) and in completion statistics (R2).
- **Unlimited mode:** When the loop is run with no limit (e.g., `--unlimited`), the message indicates the current iteration and that there is no limit. Format example: `Iteration 3 (unlimited)` or `Iteration 3 of unlimited`.

**Log level:** Progress messages are emitted at **info** log level (R5). They are visible when the effective log level is **info** or **debug**. They are **suppressed** when the effective log level is **warn** or **error** (e.g., when the user passes `--quiet` or `--log-level warn`).

**Output destination:** Progress messages go to **stdout** so they are part of the run log (Ralph operational messages and AI command stream in one stream). This is consistent with R2 and R5 (Ralph's operational output to stdout).

**Timing:** The message is emitted once per iteration, at the start of that iteration — before the AI process is spawned. So for a run that executes 3 iterations, the user sees three progress lines (e.g., "Iteration 1/10", "Iteration 2/10", "Iteration 3/10") in order.

**Dry-run:** Dry-run (R4) does not run the loop and does not execute any iteration; therefore no iteration progress messages are emitted. Dry-run only prints the assembled prompt to stdout.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Bounded run, iteration 1 of 10 | Message e.g. "Iteration 1/10" or "Iteration 1 of 10" at start of iteration 1 |
| Bounded run, iteration 5 of 10 | Message shows 5 and 10 |
| Unlimited run, iteration 3 | Message e.g. "Iteration 3 (unlimited)" or "Iteration 3 of unlimited" |
| Single-iteration run (max 1) | One progress message: "Iteration 1/1" (or equivalent) |
| Effective log level info or debug | Progress messages visible |
| Effective log level warn or error (e.g., --quiet) | Progress messages suppressed |
| Dry-run (`-d`) | No iterations run; no progress messages |
| Interruption (SIGINT) during iteration 2 | User saw "Iteration 1/10" and "Iteration 2/10"; then Ralph exits 130; no message for iteration 3 |
| Per-prompt override: max 1, preamble false | Progress message still emitted at start of that single iteration (if log level allows) |

### Examples

#### Progress visible at info level

**Input:**
`ralph run build` with default log level (info). Run executes 3 iterations; success on iteration 3.

**Expected output:**
At the start of each iteration, stdout contains a line like "Iteration 1/10", "Iteration 2/10", "Iteration 3/10" (or equivalent). Order and numbering match the iterations that run.

**Verification:**
- Three progress lines on stdout
- Each line indicates iteration number and limit
- Lines appear before the AI output (or before completion) for that iteration

#### Progress suppressed with --quiet

**Input:**
`ralph run build -q`. Run executes 2 iterations; success on iteration 2.

**Expected output:**
No "Iteration N/M" lines on stdout (R5: quiet sets log level to error; progress is info). User sees only errors (if any) and exit. Completion statistics (R2) may also be suppressed if at info level.

**Verification:**
- Exit code 0
- stdout does not contain iteration progress lines

#### Unlimited mode

**Input:**
`ralph run build --unlimited`. Run executes several iterations.

**Expected output:**
Progress messages show current iteration and indicate unlimited, e.g. "Iteration 1 (unlimited)", "Iteration 2 (unlimited)", etc. No denominator (no "of N").

**Verification:**
- Message format includes "unlimited" or equivalent
- Iteration number increments each time

## Acceptance criteria

- [ ] At the start of each iteration, Ralph prints the iteration number and the limit (e.g., "Iteration 3/10" or "Iteration 3 (unlimited)")
- [ ] Progress messages are emitted at info log level
- [ ] Progress messages go to stdout so they are part of the run log (with AI output when streamed)
- [ ] Progress messages are suppressed when log level is set above info (--quiet or --log-level warn/error)

## Dependencies

- O1/R8 (preamble injection) — iteration number in the progress message must match the iteration number in the preamble for the same iteration.
- R2 (iteration statistics) — iteration count reported at completion is the same notion of "iteration" (each start-of-iteration corresponds to one completed iteration if the loop runs to completion).
- R5 (log level control) — progress messages are at info level; they are suppressed when effective log level is warn or error.
- R4 (dry-run) — dry-run does not run iterations, so no progress messages are emitted.
