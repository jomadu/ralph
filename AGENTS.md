# AGENTS.md

## Work Tracking System

**System:** beads (bd CLI)

This project uses **bd (beads)** for all issue tracking. Run `bd onboard` to get started. Do NOT use markdown TODOs, task lists, or other tracking methods.

**Query ready work:**
```bash
bd ready --json
```

**Update status (claim work atomically):**
```bash
bd update <id> --claim --json
```

**Close issue:**
```bash
bd close <id> --reason "Completed" --json
```

**Create issue:**
```bash
bd create "Issue title" --description="Detailed context" -t bug|feature|task -p 0-4 --json
bd create "Issue title" --description="What this issue is about" -p 1 --deps discovered-from:bd-123 --json
```

Issue types: `bug`, `feature`, `task`, `epic`, `chore`. Priorities: `0` (critical) to `4` (backlog). Use `--json` for programmatic use. Link discovered work with `discovered-from` dependencies. bd auto-syncs to `.beads/issues.jsonl`; no manual export/import needed.

## Feature Input

`TASK.md` contains feature requirements and specifications for Ralph.

## Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --claim  # Claim work atomically
bd close <id>         # Complete work
bd sync               # Sync with git
```

**Non-interactive shell commands:** Use `-f` / `-rf` with `cp`, `mv`, `rm` to avoid hanging on confirmation prompts (e.g. `cp -f source dest`, `rm -rf directory`). Use `-y` for apt-get, `BatchMode=yes` for ssh/scp, `HOMEBREW_NO_AUTO_UPDATE=1` for brew.

## Planning System

`PLAN.md` documents the current plan (create when needed). Agent can run `bd create` commands to file issues from a plan.

## Build/Test/Lint Commands

Implementation not yet present (Ralph is in specification phase):

- **Test:** Manual verification (no automated tests yet)
- **Build:** Not required (no build artifact yet)
- **Lint:** Not configured

When Go or another implementation is added, document commands here and keep in sync with the repository.

## Specification Definition

Specifications live in `docs/intent/`. Index at `docs/intent/README.md`. Methodology in `building-intent.md` at repository root.

Format: Outcome/requirement hierarchy. Each outcome has a directory `O<n>-<slug>/` with `README.md` (outcome, risks, requirement one-liners) and `R<n>-<slug>.md` (requirement + specification). Every specification traces to a requirement; every requirement traces to an outcome.

Exclude: `docs/intent/README.md` is the index, not a single spec.

Current state: Intent tree defined (O1–O4 and requirements). Implementation pending.

## Implementation Definition

Location: `scripts/`, root-level scripts (e.g. `cursor-wrapper.sh`), and (when present) any future `cmd/`, `internal/`, or equivalent.

Patterns:
- `scripts/*.sh` — Scripts and wrappers
- `cursor-wrapper.sh` — Cursor integration script

Excludes: `.git/`, `.beads/`, `docs/` (specifications), `AGENTS.md`, `PLAN.md`, `TASK.md`, `building-intent.md`

Implementation status: Specs and intent complete; runtime implementation (e.g. Go loop runner) not yet built per TASK.md.

## Audit Output

Audit results written to `AUDIT.md` at repository root.

## Quality Criteria

**Specifications:**
- All outcomes have README with outcome definition and risks (PASS/FAIL)
- All requirements have traceability to an outcome (PASS/FAIL)
- Examples and implementation notes where applicable (PASS/FAIL)

**Implementation:**
- When tests exist: test command passes (PASS/FAIL)
- When build exists: build command succeeds (PASS/FAIL)
- Documentation and examples match actual behavior (PASS/FAIL)

**Refactoring triggers:**
- Spec/implementation divergence
- Test failures (when tests exist)
- Unclear or broken documentation

## Operational Learnings

Last verified: (update when verified)

**Working:**
- bd (beads) for issue tracking; `bd ready --json`, `bd update <id> --claim`, `bd close <id> --reason "..."` work
- Intent structure in `docs/intent/` with O1–O4 and requirements
- TASK.md defines Ralph scope and design; building-intent.md defines spec methodology

**Not working:**
- No automated test or build yet (implementation pending)

**Rationale:**
- Beads used for dependency-aware, git-friendly, agent-optimized tracking
- Specification model follows outcome → requirement → specification from building-intent.md
- Non-interactive shell commands and Landing the Plane reduce agent hangs and stranded work

**Session completion (Landing the Plane):** When ending a work session, complete: (1) File issues for remaining work, (2) Run quality gates if code changed, (3) Update issue status, (4) `git pull --rebase`, `bd sync`, `git push` until `git status` shows "up to date with origin", (5) Clean up and hand off. Work is NOT complete until `git push` succeeds.
