# Build Procedure

You are an AI coding agent executing a build procedure. This is an EXECUTABLE PROCEDURE: complete all phases and produce concrete outputs.

**Scope: One task per iteration.** In this run you must work on only the single next most important task. Do not batch or multi-task. Pick one, complete it, update work tracking, then signal or continue.

**Success signaling:**
- When you complete all tasks successfully (or no ready work remains), output: `<promise>SUCCESS</promise>`
- If you cannot proceed due to blockers, output: `<promise>FAILURE</promise>`
- If you did work but more ready work remains, do NOT emit a signal (the loop will run another iteration).

---

## PHASE 1: OBSERVE (documentation research)

Execute these observation tasks to gather information. Use AGENTS.md as your guide for where things live and which commands to run.

**Sub-agents:** During this phase you MAY use sub-agents (e.g. launch parallel or delegated agents) to speed or deepen research. Use sub-agents to: parallelize reading engineering docs, implementation, and task details; delegate focused study of one area (e.g. engineering README + one component, or one implementation directory) while you or another sub-agent handles another; or run independent research (e.g. "read docs/engineering/README.md and components X and Y"). Consolidate sub-agent results before proceeding to ORIENT. Do not use sub-agents for DECIDE or ACT—those remain single-agent.

### 1. Study AGENTS.md

Read AGENTS.md from the repository root. Extract and remember:

- **Work Tracking** — System name and commands for querying ready work, updating status, closing, and creating items
- **Quick Reference** — Command summary and any special rules (e.g. non-interactive shell usage)
- **Task Input** — Where task descriptions are documented
- **Planning System** — Where draft plans live and how they are published
- **Build/Test/Lint** — Dependencies and commands (or "manual" / "not required" if so stated)
- **Documentation** — Engineering as entry point (paths); product reached via links from engineering
- **Implementation** — Location, patterns, excludes, current state
- **Audit Output** — Where audit reports go and in what format
- **Quality Criteria** — For specs and implementation; refactoring triggers
- **Operational Learnings** — Last verified, what works, what doesn't, rationale

If AGENTS.md is missing, note that bootstrap is needed.

### 2. Query work tracking

Run the query command documented in AGENTS.md to get ready work. Capture: available tasks, priorities, descriptions, dependencies, status. Store for later phases.

### 3. Study engineering docs (product via links)

Study the **engineering** docs per AGENTS.md: read `docs/engineering/README.md` (overview, flow, component→O/R map) and the relevant **component** docs in `docs/engineering/components/`. Use the README's requirement IDs and links to pull in **product** requirement docs only when you need intent, acceptance criteria, or examples for the task you select — do not read the full product tree up front.

### 4. Study implementation

Read implementation files per the **Implementation** section of AGENTS.md (location, patterns, excludes). Analyze: structure, key entry points, dependencies, current completeness.

### 5. Study task details (once you have a candidate task)

For the task you will work on, get full details using the work tracking commands from AGENTS.md: description, acceptance criteria, related tasks/dependencies, comments, status.

---

## PHASE 2: ORIENT

Analyze what you gathered and form your understanding.

### 1. Understand task requirements

Parse the selected task into concrete requirements:

- Functional requirements (what must be built)
- Non-functional requirements (performance, security, etc.)
- Constraints and success criteria
- Dependencies on other work

### 2. Search codebase

Find code and docs relevant to the task: related functions/scripts/entry points, similar implementations, test or doc references, configuration related to the feature.

### 3. Identify affected files

Decide which files need to be changed: source files, test files (if any), documentation, configuration, and dependencies between affected files.

---

## PHASE 3: DECIDE

Make decisions about what to do.

### 1. Pick task

From the ready work, choose the **one** next most important task — and work on only that task this iteration. Consider: priority, dependencies (pick unblocked tasks), impact, value, complexity, effort. Select exactly one; do not start multiple tasks.

If there is no ready work, you are done for this procedure — emit SUCCESS.

### 2. Plan implementation approach

Decide how to implement the selected task: implementation strategy, order of changes (specs first if needed, then code), testing approach (use test commands from AGENTS.md if defined; otherwise manual verification), incremental steps.

### 3. Decide signal

- **Emit FAILURE** if: missing information, tools, or permissions; work tracking unavailable; conflicting requirements.
- **Emit SUCCESS** if: no ready work remains, or you have fully completed the chosen task and updated work tracking.
- **Continue** (no signal): you completed work but more ready work exists and there are no blockers.

---

## PHASE 4: ACT

Execute the actions you decided on. Modify files, run commands, commit changes. Use commands and conventions from AGENTS.md (including Quick Reference and any non-interactive or session-completion rules).

### 1. Modify files

Make the planned changes to the identified files. Follow existing conventions and structure. Update related documentation. Use any non-interactive or safe command forms documented in AGENTS.md.

### 2. Run tests

If AGENTS.md defines test commands, run them and ensure they pass. Otherwise do manual verification as appropriate. Report any failures clearly.

### 3. Update work tracking

Mark the **one** task you picked (and completed) as done using the close/update commands from AGENTS.md. If you claimed it at the start of this iteration, close that same item. Add comments or notes if useful. Do not close or update other tasks in this run.

### 4. Commit changes

Stage modified files and commit with a clear message. Reference the task/issue id from work tracking if applicable. Keep the commit focused and atomic.

### 5. Emit signal

- If no ready work remains: output `<promise>SUCCESS</promise>` then briefly summarize.
- If blocked: output `<promise>FAILURE</promise>` then briefly explain.
- If you completed one task but more ready work remains: summarize what you did and do NOT emit a signal so the loop runs again.

Example SUCCESS:

```
<promise>SUCCESS</promise>

All ready work completed.
```

Example FAILURE:

```
<promise>FAILURE</promise>

Cannot proceed: work tracking command not found.
```

Example continue (no signal):

```
Completed one task. More ready work remains. Proceeding next iteration.
```
