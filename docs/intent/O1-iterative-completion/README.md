# O1: Iterative Completion

## Statement

An AI-driven task reaches verified completion through iterative execution.

## Why it matters

Without a loop runner, the user manually invokes the AI CLI, reads its output, decides whether the task is done, and re-runs it. Every iteration is a manual cycle: read, judge, re-invoke. For multi-step tasks that require several passes (build until tests pass, refactor until lint is clean, implement until the feature works), this manual loop is the bottleneck. The AI can work autonomously within a single invocation, but it can't re-invoke itself. Ralph closes that gap — the user defines what "done" looks like (a signal), and Ralph handles the re-invocation until that signal appears or a limit is reached.

## Verification

- User runs `ralph run <alias>`. Ralph spawns a fresh AI CLI process, pipes the assembled prompt to its stdin, and captures output.
- User runs `ralph run -f ./my-prompt.md`. Ralph reads the file once, buffers it, and uses it for all iterations.
- User runs `cat prompt.md | ralph run`. Ralph reads stdin once, buffers it, and uses it for all iterations.
- The AI works, emits `<promise>SUCCESS</promise>` in its output, and exits.
- Ralph detects the success signal, reports completion, and exits 0.
- If the AI emits a failure signal instead, Ralph increments the consecutive failure counter and starts a new iteration.
- If the user runs a task that requires 7 iterations to converge, Ralph executes all 7 without manual intervention.

## Non-outcomes

- Ralph does not carry conversation history between iterations. State continuity comes from the filesystem — the AI reads and writes files, and the next iteration's AI sees those changes.
- Ralph does not validate that the AI's work is correct. It trusts the signal. If the AI says SUCCESS, Ralph stops.
- Ralph does not provide built-in prompts or impose a methodology. The user owns the prompt entirely.
- Ralph does not modify the user's prompt file on disk. Preamble wrapping happens in memory before piping.
- Ralph does not re-read the prompt between iterations. The prompt is loaded once at loop start and reused. The preamble changes per iteration (iteration count, etc.), but the underlying prompt content is immutable for the duration of the loop.
- Ralph does not retry a failed AI CLI process within the same iteration. A crash counts as one iteration.

## Risks

| Risk | Mitigating Requirement |
|----------|----------------------|
| The AI CLI crashes or exits non-zero mid-execution | [R1 — Process crash recovery with partial output capture](R1-process-crash-recovery.md) |
| Both success and failure signals appear in the same output | [R2 — Signal precedence rules](R2-signal-precedence.md) |
| The AI process hangs indefinitely | [R3 — Per-iteration timeout](R3-iteration-timeout.md) |
| The AI never emits a signal and iterations run forever | [R4 — Maximum iteration limit](R4-max-iteration-limit.md) |
| The AI repeatedly fails without making progress | [R5 — Consecutive failure tracking](R5-consecutive-failure-tracking.md) |
| Output grows without bound and exhausts memory | [R6 — Output buffer management](R6-output-buffer-management.md) |
| User sends SIGINT during AI execution | [R7 — Graceful interruption handling](R7-graceful-interruption.md) |
| The AI doesn't know it's in a loop or what signals to emit | [R8 — Preamble injection](R8-preamble-injection.md) |
| User wants to run a one-off prompt without defining an alias in config | [R9 — Prompt input modes](R9-prompt-input-modes.md) |
| Prompt file changes on disk during loop execution, causing inconsistent behavior across iterations | [R9 — Prompt input modes (read once, buffer, reuse)](R9-prompt-input-modes.md) |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R1](R1-process-crash-recovery.md) | Process crash recovery with partial output capture | ready |
| [R2](R2-signal-precedence.md) | Signal precedence rules | ready |
| [R3](R3-iteration-timeout.md) | Per-iteration timeout | ready |
| [R4](R4-max-iteration-limit.md) | Maximum iteration limit | ready |
| [R5](R5-consecutive-failure-tracking.md) | Consecutive failure tracking | ready |
| [R6](R6-output-buffer-management.md) | Output buffer management | ready |
| [R7](R7-graceful-interruption.md) | Graceful interruption handling | ready |
| [R8](R8-preamble-injection.md) | Preamble injection | ready |
| [R9](R9-prompt-input-modes.md) | Prompt input modes | ready |
