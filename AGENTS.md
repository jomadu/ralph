# AGENTS.md

## Issue tracking (beads)

This project uses **bd (beads)** for all issue tracking. Do not use markdown TODOs or other trackers. Quality gates are the build, test, and lint commands below; run them when code changes. When ending a session: file or update issues in bd, run quality gates if needed, then `git pull --rebase`, `bd dolt push`, `git push` until up to date.

**Non-interactive shell commands:** Use `-f`/`-rf` with `cp`, `mv`, `rm`; `-y` for apt-get; `BatchMode=yes` for ssh/scp; `HOMEBREW_NO_AUTO_UPDATE=1` for brew.

## Build / test / lint

From repository root (see Makefile):

- **Build:** `make build` → `bin/ralph` (optional `BINARY=`, `VERSION=`)
- **Test:** `make test` → `go test ./...`
- **Lint:** `make lint` → `go vet` + `gofmt -s` check + `lint-docs` (remark-validate-links on `docs/`); `make fmt` to fix formatting; `make lint-docs` for docs only.

Release: conventional commits + semantic-release; commit-msg hook via `npm install`. Branches: main, rc, alpha, beta.

## Documentation

**Engineering (start here):** `docs/engineering/README.md` — overview, flow, component list with O/R IDs linking to product. `docs/engineering/components/` — one file or directory per component; responsibility, interfaces, implementation specs. Requirement assignments live only in the engineering README.

**Product (via links):** Intent in `docs/product/` (outcomes, requirements). Follow O/R links from engineering when you need outcome or acceptance-criteria detail.

**Methodology:** `building-intent.md` — product vs engineering roles and consistency rules.

## Implementation definition

**In scope:** `scripts/`, and (when present) `cmd/`, `internal/`, or equivalent. `testdata/` — test fixtures; integration tests use `--config testdata/<fixture>.yml` from repo root.

**Patterns:** `scripts/*.sh` for scripts and wrappers. Agent wrapper pattern: [docs/agent-wrapper-pattern.md](docs/agent-wrapper-pattern.md). Review report summary format and exit-code derivation: `docs/engineering/components/review.md` (O005/R002, O005/R008, O010/R003); parser in `internal/review/summary.go`.

**Excludes:** `.git/`, `docs/product/`, `docs/engineering/`, `AGENTS.md`, `PLAN.md`, `TASK.md`, `building-intent.md`.

Audit output: `AUDIT.md` at repo root.
