# R2: Iteration Statistics

**Outcome:** O4 — Observability

## Requirement

The system reports iteration timing statistics at loop completion, giving the user data to understand loop performance and tune configuration. Statistics include iteration count, min/max/mean duration, and standard deviation.

## Specification

When the loop completes for a reason that produces exit 0, 1, or 2 (success, failure threshold, or max iterations), Ralph reports iteration timing statistics to stdout before exiting. Statistics are not reported when Ralph exits 130 (interruption per O1/R7); the iteration is discarded and no completion summary is produced.

**When statistics are reported:**

- After the iteration that produces a success signal → report, then exit 0.
- After the iteration that causes the consecutive failure count to reach the threshold → report, then exit 1.
- After the final iteration when max iterations is reached without success → report, then exit 2.
- On SIGINT/SIGTERM → do not report statistics; exit 130.

**What is reported:**

- **Iteration count:** Total number of iterations executed in this run (integer ≥ 1).
- **Min duration:** Shortest iteration wall-clock duration in seconds (or appropriate unit; see format).
- **Max duration:** Longest iteration wall-clock duration.
- **Mean duration:** Sum of iteration durations divided by iteration count.
- **Standard deviation:** Population standard deviation of iteration durations, computed using Welford's online algorithm so that a single pass over the durations suffices and numeric stability is maintained.

**Duration definition:** Per-iteration duration is the wall-clock time from the start of piping input to the AI CLI process until the process has exited and output has been captured (i.e., the time the iteration took, not including post-scan or loop overhead). Duration is measured in seconds; implementation may use fractional seconds (e.g., milliseconds) internally and display in a human-readable form (e.g., "1.234s" or "1.23s").

**Output destination:** All statistics output goes to stdout so that the run log (Ralph operational messages and AI command stream) is captured in one place. This aligns with R5 (Ralph's log output to stdout) and R6 (progress to stdout).

**Log level:** The iteration statistics block is emitted at **info** log level (R5). It is therefore visible at default log level and suppressed when the effective log level is **warn** or **error** (e.g. `--quiet` or `--log-level warn`).

**Single-iteration runs:** When iteration count is 1, standard deviation is undefined (single sample). The spec requires either reporting standard deviation as 0 or omitting it. If omitted, the statistics block must still include iteration count and min/max/mean (min, max, and mean are identical for one iteration).

**Format:** Exact format (e.g., line prefixes, labels, units) is implementation-defined as long as the required quantities are present and identifiable. Example shape: iteration count, then min/max/mean (and optionally stddev) with clear labels.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Success on iteration 1 | Report 1 iteration; min = max = mean; stddev 0 or omitted |
| Success on iteration 7 | Report 7 iterations; min, max, mean, stddev all reported |
| Failure threshold reached on iteration 4 | Report 4 iterations; statistics for those 4 iterations |
| Max iterations (10) exhausted | Report 10 iterations; statistics for all 10 |
| SIGINT during iteration 3 | No statistics reported; exit 130 |
| SIGINT between iterations | No statistics reported; exit 130 |
| Single iteration, success | Mean = min = max; stddev 0 or omitted |
| Single iteration, failure threshold 1 | One iteration counted; stddev 0 or omitted |
| Iteration timeout (O1/R3) then next iteration | Timed-out iteration's duration is up to timeout value; included in stats |
| Zero-duration iteration (process exits immediately) | Min can be 0 or very small; mean/max/stddev still computed |

### Examples

#### Multi-iteration success

**Input:**
`ralph run build` with default config. Iterations 1–3 run; iteration 3 output contains success signal. Iteration durations: 12.1s, 15.3s, 10.0s.

**Expected output:**
Before exiting 0, Ralph prints to stdout something equivalent to: 3 iterations executed; min duration (e.g., 10.0s); max (e.g., 15.3s); mean (e.g., 12.47s); stddev (e.g., 2.65s or similar).

**Verification:**
- Statistics appear on stdout
- Iteration count is 3
- Min ≤ mean ≤ max; stddev ≥ 0

#### Single-iteration run

**Input:**
`ralph run bootstrap` with `default_max_iterations: 1`. Single iteration runs and emits success signal. Duration 8.5s.

**Expected output:**
Ralph reports 1 iteration; min = max = mean = 8.5s (or equivalent); standard deviation is 0 or omitted.

**Verification:**
- Iteration count is 1
- No contradiction (e.g., stddev 0 or field absent)

#### Interruption — no statistics

**Input:**
`ralph run build`. User presses Ctrl-C during iteration 2.

**Expected output:**
Ralph exits 130. No iteration statistics block is printed to stdout.

**Verification:**
- Exit code 130
- stdout does not contain a completion statistics summary for this run

## Acceptance criteria

- [ ] After the loop completes (success, failure threshold, or max iterations), Ralph reports the total number of iterations executed
- [ ] Min, max, and mean iteration durations are reported
- [ ] Standard deviation of iteration durations is reported, calculated using Welford's online algorithm
- [ ] Statistics are reported to stdout (run log) so they are part of the single run log stream
- [ ] Statistics are not reported on interruption (SIGINT/SIGTERM — exit 130)
- [ ] For single-iteration runs, standard deviation is reported as 0 or omitted
- [ ] Statistics block is emitted at info log level and is suppressed when effective log level is warn or error (e.g. --quiet)

## Dependencies

- O1/R7 (graceful interruption) — defines exit 130 and that no completion processing (including statistics) is done on interrupt.
- O1 iteration loop — iteration count and per-iteration start/end are defined by the loop; duration is measured over the same iteration boundaries (e.g., from start of piping to AI CLI until process exit).
- R5 (log level control) — statistics are emitted at info level and suppressed when effective log level is warn or error.
