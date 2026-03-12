# Writing Ralph prompts

This guide explains how to write prompts that work well with Ralph's execution model. **Ralph evaluates prompts along the same four dimensions described here when you run `ralph review`.** If you address these dimensions in your prompt, the loop is more likely to converge, and review will give you fewer gaps to fix.

## 1. Signal and state

Ralph runs your AI command in a loop. To know when to stop, it needs **clear success and failure signals** it can detect (e.g. exit codes, or markers in stdout). Your prompt should tell the AI how to signal "task done" and "task failed" in a way Ralph can read.

State (files, work-tracking) should be **compatible with a fresh process each iteration**. Ralph starts a new process for each loop iteration; the AI does not keep in-memory state between runs. So the prompt should assume the AI will re-read files or state from disk each time.

**Strong:** "When the task is complete, exit 0. If you cannot complete it (e.g. blocked or invalid input), exit 1. Persist progress in `./state.json` and re-read it at the start of each run."

**Weak:** "Finish the task." (No exit code or marker; Ralph cannot tell when to stop.)

---

## 2. Iteration awareness

Execution is **multi-iteration**: Ralph may invoke the AI many times, and each run is a **new process** with no in-memory state from prior runs. The prompt should assume this execution model: re-read the task and any state from disk at the start of each run, emit the right signals, and avoid logic that only works in a single invocation.

**Do not** prescribe behavior that depends on which iteration or "pass" the run is (e.g. "if you have made more than two passes, do X"). That would require the AI to track or infer pass count and can lead to **iteration artifacts** — the AI writing pass numbers, iteration logs, or "I'm stuck on attempt N" into the repository. Each run should be **conceptually fresh**: the AI may optionally use on-disk history (e.g. git) to investigate, but the prompt should not require or assume iteration-index awareness.

**Strong:** "You run in a loop; each run is a new process. Re-read the task file and state (e.g. state.json) before continuing. Emit success (exit 0) or failure (exit 1) so the loop can stop."

**Weak:** "Do the task." (No mention of loop or re-reading; the AI may assume one shot and not re-check state.)

---

## 3. Scope and convergence

The task should have a **defined scope** and **completion criteria that are checkable in practice**. That way the loop can converge instead of running indefinitely. Vague or open-ended "done" conditions make it hard for Ralph to know when to stop.

**Strong:** "Implement the API described in SPEC.md. Done when: (1) all tests in `make test` pass, and (2) the checklist in SPEC.md is satisfied. Exit 0 when both are true; exit 1 if blocked or if you decide the spec cannot be met."

**Weak:** "Improve the codebase until it's good." (No checkable criteria; loop may not converge.)

---

## 4. Subjective completion criteria

When "done" is subjective (e.g. "good enough," "reads well"), the AI can get stuck in small repetitive tweaks. The prompt should include **techniques to escape local optima**: variation, creative exploration, or stepping back (e.g. consider alternatives, challenge assumptions). That helps the AI avoid infinite micro-edits.

Keep these techniques **iteration-agnostic**: describe what to do in a given run, not what to do "after N passes." Prescribing pass- or iteration-count rules encourages the AI to emit **artifacts** — pass counters, iteration logs, or meta-commentary — into the repo. Prefer per-run behavior that includes the **current version** in the comparison: e.g. "When revising a section, consider two alternative structures and the existing structure; pick the best (which may be keeping the current one)." That avoids churn from always choosing a "new" alternative when the existing structure is already best.

**Strong:** "Revise the doc until it reads well. When revising a section, consider two alternative structures and the existing structure; pick the best (which may be the current one). When you are satisfied, exit 0; if stuck, exit 1 and explain why."

**Weak:** "Make the doc read well. Keep editing until it's done." (No escape technique; risk of endless small edits.)

---

## Summary

| Dimension | What to address |
|-----------|------------------|
| **Signal and state** | Clear success/failure signals Ralph can detect; state that works with a fresh process each iteration. |
| **Iteration awareness** | Prompt assumes multi-iteration and fresh process; re-read state each run, emit signals. Do not prescribe behavior by iteration or pass count (avoids iteration artifacts in the repo). |
| **Scope and convergence** | Defined scope and checkable completion criteria so the loop can converge. |
| **Subjective completion** | When "done" is subjective, add variation or stepping-back techniques that are per-run (not "after N passes") to avoid getting stuck and avoid artifacts. |

These are the same four dimensions **`ralph review`** uses. You can run `ralph review` on your prompt to get AI-generated feedback and a suggested revision along these dimensions. For a short reminder in the terminal, run `ralph show prompt-guide`.

---

## Tailoring the build prompt to this repo

The build prompt (`prompts/build.md`) is kept **tight and repo-specific** so the agent spends tokens on the task, not on rediscovering layout. Tailoring choices:

1. **Bake in the layout** — Instead of "read AGENTS.md and extract work tracking, docs, implementation," the prompt states: work tracking is bd with specific commands; spec is `docs/engineering/` + components, product via O/R links; implementation is `cmd/`, `internal/`, `scripts/`. The agent still uses AGENTS.md for build/test/lint details and session rules, but no longer re-derives the map every run.

2. **One source of truth for signals** — Success/failure/continue and the `<promise>` format are defined once at the top; later sections only say when to emit which signal. Keeps the prompt consistent with the four dimensions (signal and state, iteration awareness).

3. **Phases collapsed into three steps** — (1) Get the next task from bd, load only the context needed for that task. (2) Plan and implement; run quality gates and close the task in bd. (3) Emit the correct signal. No separate OBSERVE/ORIENT/DECIDE/ACT phases; the flow is "pull task → do task → signal."

4. **Task selection is explicit** — "Next highest-priority ready item," "exactly one task per run," and "re-query bd if unsure" keep scope and convergence clear and avoid iteration-count logic.

5. **No sub-agent instructions** — Optional parallelization is omitted from the streamlined prompt; the agent can still use tools as needed, but the procedure does not prescribe when to delegate.

When you change the repo (e.g. new tooling, new doc structure), update the repo layout block in `prompts/build.md` so the prompt stays accurate.
