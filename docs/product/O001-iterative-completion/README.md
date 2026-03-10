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

## Non-outcomes

- Ralph does not inject or dictate the signal the AI should expose. It adapts to whatever signal the prompt (and thus the user) decides to expose, by allowing the user to configure which success and failure signals to look for in the output.
- Ralph does not carry conversation history between iterations. State continuity comes from the filesystem — the AI reads and writes files, and the next iteration's AI sees those changes.
- Ralph does not validate that the AI's work is correct. It trusts the signal. If the AI says success, Ralph stops.
- Ralph does not provide built-in prompts or impose a methodology. The user owns the prompt entirely.
- Ralph does not modify the user's prompt file on disk. Preamble wrapping happens in memory before piping.
- Ralph does not re-read the prompt between iterations. The prompt is loaded once at loop start and reused. The preamble may change per iteration (e.g. iteration count), but the underlying prompt content is immutable for the duration of the loop.
- Ralph does not retry a failed AI CLI process within the same iteration. A crash counts as one iteration.
