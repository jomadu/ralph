# AGENTS.md

## Work Tracking System

This project uses **bd (beads)** for all issue tracking. Run `bd --help` for commands; `bd prime` for full workflow context. Do not use markdown TODOs or other trackers. `PLAN.md` and `TASK.md` provide feature context.

## Feature Input

`TASK.md` contains feature requirements and specifications for Ralph.

**Non-interactive shell commands:** Use `-f` / `-rf` with `cp`, `mv`, `rm` to avoid hanging on confirmation prompts (e.g. `cp -f source dest`, `rm -rf directory`). Use `-y` for apt-get, `BatchMode=yes` for ssh/scp, `HOMEBREW_NO_AUTO_UPDATE=1` for brew.

## Planning System

`PLAN.md` documents the current plan (create when needed).

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
- `testdata/` — Test fixture files (e.g. config YAMLs used only by tests). See O2-R5 for the convention; integration tests use `--config testdata/<fixture>.yml` from repo root.

Excludes: `.git/`, `docs/` (specifications), `AGENTS.md`, `PLAN.md`, `TASK.md`, `building-intent.md`

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
- Intent structure in `docs/intent/` with O1–O4 and requirements
- TASK.md defines Ralph scope and design; building-intent.md defines spec methodology

**Not working:**
- No automated test or build yet (implementation pending)

**Rationale:**
- Specification model follows outcome → requirement → specification from building-intent.md
- Non-interactive shell commands reduce agent hangs and stranded work

**Session completion:** When ending a work session: (1) Note remaining work in PLAN.md or TODOs, (2) Run quality gates if code changed, (3) `git pull --rebase`, `git push` until `git status` shows "up to date with origin", (4) Hand off with context for next session.

<!-- BEGIN BEADS INTEGRATION -->
## Issue Tracking with bd (beads)

**IMPORTANT**: This project uses **bd (beads)** for ALL issue tracking. Do NOT use markdown TODOs, task lists, or other tracking methods.

### Why bd?

- Dependency-aware: Track blockers and relationships between issues
- Git-friendly: Dolt-powered version control with native sync
- Agent-optimized: JSON output, ready work detection, discovered-from links
- Prevents duplicate tracking systems and confusion

### Quick Start

**Check for ready work:**

```bash
bd ready --json
```

**Create new issues:**

```bash
bd create "Issue title" --description="Detailed context" -t bug|feature|task -p 0-4 --json
bd create "Issue title" --description="What this issue is about" -p 1 --deps discovered-from:bd-123 --json
```

**Claim and update:**

```bash
bd update <id> --claim --json
bd update bd-42 --priority 1 --json
```

**Complete work:**

```bash
bd close bd-42 --reason "Completed" --json
```

### Issue Types

- `bug` - Something broken
- `feature` - New functionality
- `task` - Work item (tests, docs, refactoring)
- `epic` - Large feature with subtasks
- `chore` - Maintenance (dependencies, tooling)

### Priorities

- `0` - Critical (security, data loss, broken builds)
- `1` - High (major features, important bugs)
- `2` - Medium (default, nice-to-have)
- `3` - Low (polish, optimization)
- `4` - Backlog (future ideas)

### Workflow for AI Agents

1. **Check ready work**: `bd ready` shows unblocked issues
2. **Claim your task atomically**: `bd update <id> --claim`
3. **Work on it**: Implement, test, document
4. **Discover new work?** Create linked issue:
   - `bd create "Found bug" --description="Details about what was found" -p 1 --deps discovered-from:<parent-id>`
5. **Complete**: `bd close <id> --reason "Done"`

### Auto-Sync

bd automatically syncs via Dolt:

- Each write auto-commits to Dolt history
- Use `bd dolt push`/`bd dolt pull` for remote sync
- No manual export/import needed!

### Important Rules

- ✅ Use bd for ALL task tracking
- ✅ Always use `--json` flag for programmatic use
- ✅ Link discovered work with `discovered-from` dependencies
- ✅ Check `bd ready` before asking "what should I work on?"
- ❌ Do NOT create markdown TODO lists
- ❌ Do NOT use external issue trackers
- ❌ Do NOT duplicate tracking systems

For more details, see README.md and docs/QUICKSTART.md.

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   bd dolt push
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds

<!-- END BEADS INTEGRATION -->
