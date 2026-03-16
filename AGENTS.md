# AGENTS.md

## Work Tracking System

This project uses **bd (beads)** for all issue tracking. Run `bd --help` for commands; `bd prime` for full workflow context. Do not use markdown TODOs or other trackers. `PLAN.md` and `TASK.md` provide feature context.

## Feature Input

`TASK.md` contains feature requirements and specifications for Ralph.

**Non-interactive shell commands:** Use `-f` / `-rf` with `cp`, `mv`, `rm` to avoid hanging on confirmation prompts (e.g. `cp -f source dest`, `rm -rf directory`). Use `-y` for apt-get, `BatchMode=yes` for ssh/scp, `HOMEBREW_NO_AUTO_UPDATE=1` for brew.

## Planning System

`PLAN.md` documents the current plan (create when needed).

## Build/Test/Lint Commands

Use the Makefile from the repository root:

- **Build:** `make build` — builds `bin/ralph` (override with `make build BINARY=<path>`). Optional `VERSION=<ver>` (e.g. from `git describe --tags`) sets binary version; default is `dev`.
- **Test:** `make test` — runs `go test ./...`
- **Lint:** `make lint` — runs `go vet ./...` and checks `gofmt -s` (fails if any file needs formatting)
- **Format:** `make fmt` — runs `gofmt -s -w .` to fix formatting

All: `make all` (default target is build). Clean: `make clean`.

**Release:** Conventional commits + semantic-release. Commit messages are enforced by the local commit-msg hook (husky + commitlint); run `npm install` to install the hook. CI runs commitlint on PR and push (commit range: PR = base..head, push = before..after; new-branch push uses --last). Branches `main` (stable), `rc`, `alpha`, `beta` (pre-releases). Run `npm run semantic-release` for a dry-run; CI runs semantic-release on push to those branches. Release process: conventional commits and semantic-release on push to main/rc/alpha/beta.

Keep this section in sync with the Makefile.

## Documentation (engineering as entry point)

Use **engineering** as the primary context for build/implementation. Product holds intent; engineering holds placement and implementation specs and links to product.

- **Engineering (start here):** `docs/engineering/README.md` — overview, high-level flow, **component list with O/R requirement IDs linking to product**. `docs/engineering/components/` — one file or directory per component; each has responsibility, interfaces, and **implementation specifications** (schemas, APIs, protocols). Requirement assignments live only in the engineering README; component docs link back to it.
- **Product (via links):** Intent lives in `docs/product/` (outcomes, requirements). Engineering README and component docs reference product by O/R ID (e.g. links to `../product/O001-.../R002-....md`). For a given task, follow those links to read the relevant product requirement(s) when you need outcome/acceptance-criteria detail; do not read the full product tree up front.
- **Methodology:** Product vs engineering roles, structure, and consistency rules: `building-intent.md`.

**Current state:** Engineering overview and components defined; runtime implementation (e.g. Go loop runner) not yet built per TASK.md.

## Implementation Definition

Location: `scripts/`, and (when present) any future `cmd/`, `internal/`, or equivalent.

Patterns:
- `scripts/*.sh` — Scripts and wrappers
- **Agent wrapper pattern:** For any agent that outputs structured/noisy data, users can configure a wrapper script (progress → stderr, assistant text → stdout) so Ralph can scan for signals; see [docs/agent-wrapper-pattern.md](docs/agent-wrapper-pattern.md). Cursor is one example; that doc links to Cursor’s headless docs.
- `testdata/` — Test fixture files (e.g. config YAMLs used only by tests). See O002/R005 for the convention; integration tests use `--config testdata/<fixture>.yml` from repo root.
- **Review report summary:** Machine-parseable line format and exit code derivation are specified in `docs/engineering/components/review.md` (and product O005/R002, O005/R008, O010/R003); parser in `internal/review/summary.go`.

Excludes: `.git/`, `docs/product/`, `docs/engineering/` (documentation only), `AGENTS.md`, `PLAN.md`, `TASK.md`, `building-intent.md`

Implementation status: Product and engineering docs complete; runtime implementation (e.g. Go loop runner) not yet built per TASK.md.

## Audit Output

Audit results written to `AUDIT.md` at repository root.

## Quality Criteria

**Product and engineering docs:**
- All outcomes have README with outcome definition and risks (PASS/FAIL)
- All requirements have traceability to an outcome (PASS/FAIL)
- Engineering README assigns every product requirement to at least one component; every O/R in engineering exists in product (PASS/FAIL)
- Examples and implementation notes where applicable (PASS/FAIL)

**Implementation:**
- When tests exist: test command passes (PASS/FAIL)
- When build exists: build command succeeds (PASS/FAIL)
- Documentation and examples match actual behavior (PASS/FAIL)

**Refactoring triggers:**
- Product/engineering doc and implementation divergence
- Test failures (when tests exist)
- Unclear or broken documentation

## Operational Learnings

Last verified: (update when verified)

**Working:**
- Product structure in `docs/product/` (O001–O011, three-digit IDs) and engineering structure in `docs/engineering/` (overview + components)
- TASK.md defines Ralph scope and design; building-intent.md defines product/engineering methodology (P1–P4, E1–E2)

**Not working:**
- No automated test or build yet (implementation pending)

**Rationale:**
- Product = intent (who, what, why); engineering = placement and implementation specs (where, how); both follow building-intent.md
- Non-interactive shell commands reduce agent hangs and stranded work

**Session completion:** When ending a work session: (1) Note remaining work in PLAN.md or file issues in bd, (2) Run quality gates if code changed, (3) `git pull --rebase`, `git push` until `git status` shows "up to date with origin", (4) Hand off with context for next session.

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
