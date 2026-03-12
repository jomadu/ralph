# Prompt review summary

## 1. Signal and state

**Assessment:** Strong. The prompt defines clear, machine-parseable signals: `<promise>SUCCESS</promise>`, `<promise>FAILURE</promise>`, and no signal when more work remains. State is externalized in work tracking (bd); the procedure instructs the agent to query ready work at the start (PHASE 1) and to close the completed task in PHASE 4, so a fresh process can re-query and continue correctly.

**Gap:** The prompt does not explicitly state that each run is a **fresh process** with no in-memory state from prior runs. Making that explicit (and tying it to “re-read AGENTS.md and re-query work tracking every run”) would make statefulness and re-entrancy unambiguous for Ralph.

---

## 2. Iteration awareness

**Assessment:** Partially present. The prompt says “one task per iteration” and “do NOT emit a signal so the loop runs again” when more work remains. It does not state that Ralph runs the procedure in a **new process each iteration**, so the agent may not assume it must re-read docs and re-query work tracking at the start of every run.

**Recommendation:** Add a short “Execution context” (or equivalent) block that states: each iteration is a fresh process; the agent has no memory of previous runs; at the start of each run the agent must re-read AGENTS.md and re-query work tracking to get current state.

---

## 3. Scope and convergence

**Assessment:** Good. Scope is explicit (“one task per iteration”; no batching). Completion is checkable: SUCCESS when there is no ready work or when the chosen task is fully done and work tracking is updated; continue when one task was completed but more ready work exists. The loop converges when the work tracking system reports no ready work and the agent emits SUCCESS.

---

## 4. Subjective completion criteria and escape

**Assessment:** Completion is mostly objective (driven by ready work and task closure). Remaining risk: the agent might believe a task is done without actually closing it in work tracking, or might repeatedly continue with similar outcomes (e.g., always picking the same type of task or misreading “ready work”).

**Recommendation:** Add a brief escape note: if the agent has continued multiple times with similar summaries (e.g., “completed one task, more work remains”), it should re-query ready work, verify the last task was closed, and consider whether the selected task or completion criteria need to be re-evaluated. This reduces the chance of getting stuck in a no-progress loop.

---

## Overall

The prompt is well suited for Ralph: clear signals, external state, and convergent scope. Two small improvements would make it more robust: (1) explicit iteration context (fresh process, re-read state each run), and (2) a short escape technique when multiple continues occur with similar outcomes. The suggested revision and diff implement these.
