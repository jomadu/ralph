# Build Procedure (Ralph repo)

You are an AI coding agent. Each run is a **fresh process**; re-read this prompt and re-query work tracking at the start of every run. Work on **exactly one** task per run: the next highest-priority ready item from the queue.

**Repo layout (use this; do not re-derive from AGENTS.md):**
- **Work tracking:** bd (beads). Ready work: `bd ready --json`. Claim: `bd update <id> --claim`. Close: `bd close <id> --reason "Done"`. Create/link: `bd create "..." --description="..." -p 0-4 --deps discovered-from:<id>`. Use `--json` for programmatic use.
- **Spec:** `docs/engineering/README.md` (overview, component→O/R map) and `docs/engineering/components/` (implementation specs). Intent/acceptance criteria: follow O/R links into `docs/product/` only when needed for the chosen task.
- **Implementation:** `cmd/`, `internal/`, `scripts/`; tests and fixtures in `testdata/` as noted in component docs. Excludes: `.git/`, `docs/`, `AGENTS.md`, `PLAN.md`, `TASK.md`, `building-intent.md`.
- **Build/test/lint:** From repo root, `make build`, `make test`, `make lint`, `make fmt`. See AGENTS.md for release/session rules.
- **Development approach:** Use test-driven development (TDD) for code: write or update tests to define the behavior, then implement to satisfy them. This applies primarily to code in `cmd/`, `internal/`, and `scripts/`; documentation work does not require TDD.

**Signals (Ralph detects these):**
- `<promise>SUCCESS</promise>` — Only when all ready work has been completed or closed (i.e. `bd ready --json` is empty).
- `<promise>FAILURE</promise>` — Blocked (missing info, tools, or work tracking unavailable).
- No signal — You completed one task and more ready work remains; loop continues.

**Task close-out:** When you complete a task (implementation done, or you determined it is already done / no-op / not applicable), you **must** close it in bd: `bd close <id> --reason "Done"`. Do not leave completed work unclosed.

Do not encode iteration or pass counts in state or output. When in doubt, re-run `bd ready --json` and confirm the task you closed is no longer in the list.

---

## 1. Get the next task

- Run `bd ready --json`. If empty, output `<promise>SUCCESS</promise>` and stop.
- Pick the **one** next most important task (priority, dependencies, impact). Claim it: `bd update <id> --claim`.
- Load task context: read the task description and any linked items; read `docs/engineering/README.md` and the component doc(s) that apply; follow O/R links into `docs/product/` only for that task’s requirements. Inspect implementation in `cmd/`, `internal/`, `scripts/` as needed.

## 2. Plan and do

- Turn the task into concrete steps (what to build or change, which files, how to verify).
- For **code** (cmd/, internal/, scripts/): use TDD — add or update tests first, then implement to make them pass. For docs-only work, implement directly.
- Implement: edit files, run `make test` / `make lint` if applicable.
- **Close the task:** When the task is complete (work done, or already done / no-op), run `bd close <id> --reason "Done"`. This is required so the task is marked done and does not reappear in `bd ready`.
- Commit with a clear message; reference the bd id if useful.

## 3. Signal

- All ready work completed or closed (bd ready empty): output `<promise>SUCCESS</promise>` and a one-line summary.
- Blocked: output `<promise>FAILURE</promise>` and a brief reason.
- Task done (and closed in bd) but more ready work exists: summarize what you did; do **not** emit a signal (loop will run again).
