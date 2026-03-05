# R4: Maximum Iteration Limit

**Outcome:** O1 — Iterative Completion

## Requirement

The system enforces a configurable upper bound on the number of iterations, stopping the loop when the limit is reached without a success signal. An unlimited mode disables the limit, running until a signal or failure threshold is reached.

## Specification

Ralph maintains an iteration counter, starting at 0 and incrementing before each iteration begins. When `iteration_mode` is `"max-iterations"`, Ralph checks the counter before starting each iteration. If the counter exceeds `default_max_iterations`, the loop stops and Ralph exits with code 2.

**Algorithm:**

```
iteration = 0
loop:
    iteration += 1
    if iteration_mode == "max-iterations" and iteration > max_iterations:
        exit(2)
    // ... execute iteration ...
    // ... check success (exit 0), failure threshold (exit 1) ...
    goto loop
```

The iteration limit check runs *before* the iteration executes. Success detection (R2) and failure threshold (R5) run *after*. This means:

- If the last allowed iteration produces a success signal, the loop exits 0 — not 2
- If the last allowed iteration produces a failure that hits the threshold, the loop exits 1 — not 2
- Exit code 2 only occurs when all allowed iterations completed without a success signal and the next iteration would exceed the limit

**Configuration:**

- Field: `default_max_iterations`
- Type: positive integer (≥ 1)
- Default: `5`
- CLI override: `--max-iterations N` or `-n N`

- Field: `iteration_mode`
- Type: enum — `"max-iterations"` | `"unlimited"`
- Default: `"max-iterations"`
- CLI override: `--unlimited` or `-u` (sets mode to `"unlimited"`)

When `iteration_mode` is `"unlimited"`, the iteration counter still increments (it is used in preamble injection per R8 and in logging) but the limit check is skipped. The loop runs until a success signal, failure threshold (R5), or interruption (R7).

**`--unlimited` and `--max-iterations` interaction:**

These flags set different fields. `--unlimited` sets `iteration_mode` to `"unlimited"`. `--max-iterations N` sets `default_max_iterations` to N. When both are provided, `iteration_mode` is `"unlimited"` and the max iterations value is stored but not enforced.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| `default_max_iterations: 1` | Exactly one iteration runs. If no success signal, exit 2. |
| `--max-iterations 10 --unlimited` | `iteration_mode` is `"unlimited"`; the limit of 10 is stored but not enforced. Loop runs until success, failure threshold, or interruption. |
| Unlimited mode with no success signal and no failure signals | Loop runs indefinitely until interrupted (R7). |
| Last allowed iteration produces a success signal | Exit 0 (success), not exit 2 (exhaustion). |
| Last allowed iteration produces a failure that reaches the threshold | Exit 1 (failure threshold per R5), not exit 2 (exhaustion). |

### Examples

#### Exhaustion after 5 iterations

**Input:**
`default_max_iterations: 5`, `iteration_mode: "max-iterations"`. All 5 iterations produce no-signal output.

**Expected output:**
Iterations 1–5 execute. Before iteration 6 would start, the counter (6) exceeds the limit (5). Ralph exits with code 2.

**Verification:**
- Ralph exit code is 2
- Exactly 5 iterations executed

#### Success on the last iteration

**Input:**
`default_max_iterations: 3`. Iterations 1 and 2 produce no-signal output. Iteration 3 produces the success signal.

**Expected output:**
Iteration 3 executes, success signal detected. Ralph exits 0.

**Verification:**
- Ralph exit code is 0 (not 2)
- 3 iterations executed

## Acceptance criteria

- [ ] The loop executes at most the configured number of iterations
- [ ] When the limit is reached without a success signal, Ralph exits with code 2
- [ ] The default limit is 5 iterations
- [ ] The limit is configurable through the standard configuration hierarchy
- [ ] Unlimited mode (--unlimited / -u) disables the iteration limit, running until a success signal, failure threshold, or interruption
- [ ] The iteration limit must be at least 1

## Dependencies

_None identified._
