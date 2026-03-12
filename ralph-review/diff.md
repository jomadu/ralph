# Diff: original vs revision

## Prose summary of changes

1. **Execution context (new block after opening paragraph)**  
   Added an explicit “Execution context (Ralph)” paragraph stating that each iteration runs in a fresh process, the agent has no memory of prior runs, and at the start of every run the agent must re-read AGENTS.md and re-query work tracking. This addresses **signal/state** and **iteration awareness** by making re-entrancy and state re-read explicit.

2. **Escape technique in PHASE 3 (DECIDE)**  
   Under “3. Decide signal”, added an **Escape** bullet: if the agent has continued multiple times with similar outcomes, it should re-query ready work, verify the last task was closed, and consider re-evaluating task selection or completion criteria. This addresses **subjective completion criteria / getting stuck** with a light-weight escape without changing the main flow.

No other sections were modified. Scope, convergence, signaling, and phases are unchanged.

---

## Unified diff (key sections)

```diff
--- original.md
+++ revision.md
@@ -5,6 +5,10 @@
 You are an AI coding agent executing a build procedure. This is an EXECUTABLE PROCEDURE: complete all phases and produce concrete outputs.

+**Execution context (Ralph):** Each iteration runs in a **fresh process**. You have no memory of prior runs. At the start of every run, re-read AGENTS.md and re-query work tracking to get current state; do not assume any in-memory state from a previous iteration.
+
 **Scope: One task per iteration.** In this run you must work on only the single next most important task. Do not batch or multi-task. Pick one, complete it, update work tracking, then signal or continue.

 **Success signaling:**
@@ -95,6 +99,8 @@
 - **Continue** (no signal): you completed work but more ready work exists and there are no blockers.

+**Escape:** If you have continued multiple times with similar outcomes (e.g., repeatedly reporting "completed one task, more work remains"), re-query ready work, verify the last task was actually closed in work tracking, and consider re-evaluating your task selection or completion criteria before proceeding.
+
 ---
```
