# Writing Ralph prompts

This guide explains how to write prompts that work well with Ralph's execution model. **Ralph evaluates prompts along the same four dimensions described here when you run `ralph review`.** If you address these dimensions in your prompt, the loop is more likely to converge, and review will give you fewer gaps to fix.

**What the model receives:** On `ralph run`, your prompt file is wrapped in an **INSTRUCTIONS** block. Ralph may prepend a **CONTEXT** block (loop description, iteration line, and any `-c` context). Each block starts with a title line such as `# --- CONTEXT ---` and `# --- INSTRUCTIONS ---` (a leading `#` avoids CLIs that treat a bare `---` line as YAML frontmatter). Use `ralph run <alias> --dry-run` to see the exact assembled text.

## 1. Signal and state

Ralph injects a **preamble** (when enabled) that explains the loop and current iteration; your prompt can focus on the task and signals. Ralph needs **clear success and failure signals** it can detect (e.g. exit codes, or markers in stdout). Your prompt should tell the AI how to signal "task done" and "task failed" in a way Ralph can read.

**Emit the success or failure signal on the last line of your response.** Ralph only scans the **last non-empty line** of the AI output for these signals. If the outcome appears earlier (e.g. the AI explains the protocol or echoes a marker in the middle of the output), Ralph will not treat it as the final outcome. Putting the real outcome on the last line avoids false positives and ensures the loop stops or continues correctly.

State (files, work-tracking) should be **compatible with a fresh process each iteration**. Ralph starts a new process for each loop iteration; the AI does not keep in-memory state between runs. So the prompt should assume the AI will re-read files or state from disk each time.

**Strong:** "When the task is complete, exit 0. If you cannot complete it (e.g. blocked or invalid input), exit 1. Persist progress in `./state.json` and re-read it at the start of each run. Put your final outcome (success or failure) on the last line of your output so Ralph can detect it."

**Weak:** "Finish the task." (No exit code or marker; Ralph cannot tell when to stop.)

---

## 2. Iteration awareness

Ralph’s **preamble** (injected before your prompt when enabled) already tells the AI that execution is multi-iteration with a fresh process each time. Your prompt **does not need to repeat that**. Focus on:

- **Avoid iteration-index-dependent behavior** — Do not prescribe behavior that depends on which iteration or "pass" the run is (e.g. "if you have made more than two passes, do X"). That can lead to **iteration artifacts**: pass counters, iteration logs, or "I'm stuck on attempt N" in the repo.
- **State that works with re-reads** — Assume the AI will re-read task and state from disk each run (the preamble explains the fresh process). So the prompt should not assume in-memory state across runs.

**Strong:** "Re-read the task file and state (e.g. state.json) at the start of each run. Emit success (exit 0) or failure (exit 1) on the last line so the loop can stop."

**Weak:** "Do the task." (No signals; the AI may not re-check state or emit what Ralph needs.)

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
| **Signal and state** | Clear success/failure signals Ralph can detect; emit the outcome signal on the last line (Ralph scans only the last non-empty line); state that works with a fresh process each iteration. |
| **Iteration awareness** | Preamble explains the loop; prompt need not repeat it. Avoid pass-count or iteration-dependent behavior (avoids artifacts). State compatible with re-read each run. |
| **Scope and convergence** | Defined scope and checkable completion criteria so the loop can converge. |
| **Subjective completion** | When "done" is subjective, add variation or stepping-back techniques that are per-run (not "after N passes") to avoid getting stuck and avoid artifacts. |

These are the same four dimensions **`ralph review`** uses. You can run `ralph review` on your prompt to get AI-generated feedback and a suggested revision along these dimensions. For a short reminder in the terminal, run `ralph show prompt-guide`.
