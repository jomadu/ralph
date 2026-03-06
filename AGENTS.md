# AGENTS.md

## Work Tracking System

This project does not use a dedicated issue-tracking CLI. Use markdown TODOs, task lists, or your preferred method for tracking work. `PLAN.md` and `TASK.md` provide feature context.

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
