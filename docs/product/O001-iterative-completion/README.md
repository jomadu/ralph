# O001: Iterative Completion

## Who

Users who run AI-driven tasks that have a clear "done" condition and want that condition reached without manually re-invoking the AI after each run. Examples: build until tests pass; refactor until lint is clean; burn down a task list until it is fully empty; improve an artifact until some quality criteria is met (e.g. develop a threat model for a repository).

## Statement

An AI-driven task reaches verified completion through iterative execution.

## Why it matters

Without a loop runner, the user manually invokes the AI CLI, reads its output, decides whether the task is done, and re-runs it. Every iteration is a manual cycle: read, judge, re-invoke. That applies whether the task is "build until tests pass," "burn down this task list until it's empty," or "improve this artifact until it meets quality criteria (e.g. a complete threat model)." For any multi-step task that requires several passes, the manual loop is the bottleneck. The AI can work autonomously within a single invocation, but it can't re-invoke itself. Ralph closes that gap — the user defines what "done" looks like (a signal), and Ralph handles the re-invocation until that signal appears or a limit is reached.

## Verification

- When the chosen AI command is not available (e.g. missing binary, invalid alias), Ralph fails before starting the loop with a clear message, so the user does not see a confusing process start failure on the first iteration.
- User runs `ralph run <alias>`. Ralph spawns a fresh AI CLI process, pipes the assembled prompt to its stdin, and captures output.
- User runs `ralph run -f ./my-prompt.md`. Ralph reads the file once, buffers it, and uses it for all iterations.
- User runs `cat prompt.md | ralph run`. Ralph reads stdin once, buffers it, and uses it for all iterations.
- The AI works, emits the configured success signal in its output, and exits. Ralph detects the success signal, reports completion, and exits 0.
- If the AI emits a failure signal instead, Ralph increments the consecutive failure counter and starts a new iteration.
- If the user runs a task that requires multiple iterations to converge, Ralph executes them without manual intervention.
- When both success and failure signals appear in the same output, Ralph applies a defined precedence (default). Optionally, the user can enable AI-interpreted precedence: Ralph invokes the AI once with a built-in prompt asking it to interpret the iteration output and decide success or failure; if that run does not yield a clear answer, Ralph applies a defined fallback (e.g. treat as failure or use static precedence).

## Non-outcomes

- Ralph does not inject or dictate the signal the AI should expose. It adapts to whatever signal the prompt (and thus the user) decides to expose, by allowing the user to configure which success and failure signals to look for in the output.
- Ralph does not carry conversation history between iterations. State continuity comes from the filesystem — the AI reads and writes files, and the next iteration's AI sees those changes.
- Ralph does not validate that the AI's work is correct. It trusts the signal. If the AI says success, Ralph stops.
- Ralph does not provide built-in prompts or impose a methodology. The user owns the prompt entirely.
- Ralph does not modify the user's prompt file on disk. Preamble wrapping happens in memory before piping.
- Ralph does not re-read the prompt between iterations. The prompt is loaded once at loop start and reused. The preamble may change per iteration (e.g. iteration count), but the underlying prompt content is immutable for the duration of the loop.
- Ralph does not retry a failed AI CLI process within the same iteration. A crash counts as one iteration.
- The built-in prompt used for optional AI-interpreted signal precedence is owned by Ralph (not user-editable). When that option is used, exactly one extra AI invocation is made per ambiguous iteration; there are no retries of the interpretation step.

## Risks

| Risk | Mitigating Requirement |
|------|------------------------|
| AI command missing or invalid; user sees process failure on first iteration | [R001 — Validate AI command before loop](R001-validate-ai-command-before-loop.md) |
| Prompt source missing or unreadable; loop starts then fails | [R002 — Load prompt once and buffer](R002-load-prompt-once-and-buffer.md) |
| Success and failure signals both appear in output; outcome ambiguous | [R006 — Signal precedence](R006-signal-precedence.md), [R008 — AI-interpreted signal precedence](R008-ai-interpreted-signal-precedence.md) |
| Loop runs without bounded exit | [R005 — Detect failure signal and continue or exit](R005-detect-failure-signal-continue-or-exit.md), [R007 — Exit on max iterations](R007-exit-on-max-iterations.md) |
| AI process crashes or exits without success/failure signal; behavior undefined | [R009 — Process exit without signal](R009-process-exit-without-signal.md) |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R001](R001-validate-ai-command-before-loop.md) | The system validates that the chosen AI command or alias is available before starting the loop. | draft |
| [R002](R002-load-prompt-once-and-buffer.md) | The system loads the prompt once from the chosen source (file or stdin) and buffers it for all iterations, or fails before the loop with a clear message if the source is unavailable. | draft |
| [R003](R003-run-via-alias-file-or-stdin.md) | The system supports running the loop by alias, by prompt file path, or by reading the prompt from stdin. | draft |
| [R004](R004-detect-success-signal-exit-zero.md) | The system detects the configured success signal in AI output and exits with code 0. | draft |
| [R005](R005-detect-failure-signal-continue-or-exit.md) | The system detects the configured failure signal, increments the consecutive-failure count, and either starts a new iteration or exits based on the failure threshold. | draft |
| [R006](R006-signal-precedence.md) | The system applies a defined precedence when both success and failure signals are present in the same output (default behavior). | draft |
| [R007](R007-exit-on-max-iterations.md) | The system exits when the maximum iteration count is reached. | draft |
| [R008](R008-ai-interpreted-signal-precedence.md) | The system may optionally resolve signal precedence by invoking the AI once with a built-in Ralph prompt that asks the AI to interpret the iteration output and decide success or failure; if the interpretation run does not yield a clear answer, the system applies a defined fallback (e.g. treat as failure or use static precedence). | draft |
| [R009](R009-process-exit-without-signal.md) | When the AI process exits without emitting the configured success or failure signal (e.g. crash, kill, abnormal exit), the system treats the iteration as a failure, increments the consecutive-failure count, and continues or exits according to the failure threshold; the user can distinguish this condition from signal-based failure where documented. | draft |
