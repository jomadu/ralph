# Plan: Implement `ralph show prompt-guide`

**Status:** In progress  
**Goal:** Implement the `ralph show prompt-guide` subcommand so users get prompt-writing guidance from the CLI. The command outputs the **full guide verbatim** from a single source of truth. When all tasks are done, delete this plan.

**Beads (bd):** Tasks published. T1 → `ralph-0cx`, T2 → `ralph-qg1` (blocks: ralph-0cx), T3 → `ralph-dld` (blocks: ralph-qg1), T4 → `ralph-uc5` (blocks: ralph-qg1). Run `bd ready` for unblocked work; `bd show ralph-0cx` etc. for detail.

**Canonical specs (implement to these):**

- **CLI behavior:** [docs/engineering/components/cli.md](docs/engineering/components/cli.md) — section "ralph show", syntax for `show prompt-guide [flags]`, `--markdown`, and show error handling.
- **Product:** [O007/R005](docs/product/O007-user-documentation/R005-prompt-writing-guidance.md) (doc coverage), [O008/R005](docs/product/O008-discoverability/R005-prompt-writing-guide-discoverable.md) (discoverability).
- **Content source (single source of truth):** [docs/writing-ralph-prompts.md](docs/writing-ralph-prompts.md). The CLI must output this document **verbatim** so the guidance is defined in one place only. Users find the doc at `docs/writing-ralph-prompts.md`; they must not need to look inside Go code to read it.

**Mental model of the guide (writing-ralph-prompts.md):**

- **Signal and state** — Clear success/failure signals Ralph can detect; state that works with a fresh process each iteration (persist and re-read from disk).
- **Iteration awareness** — Assume multi-iteration and fresh process; re-read state each run; emit signals. Do *not* prescribe behavior by iteration or pass count (avoids iteration artifacts in the repo). Each run is conceptually fresh; the AI may optionally use on-disk history (e.g. git) to investigate.
- **Scope and convergence** — Defined scope and checkable completion criteria so the loop can converge.
- **Subjective completion** — When "done" is subjective, use per-run techniques (variation, consider alternatives and the *existing* structure; pick the best, which may be keeping the current one). Avoid pass-count or iteration-count rules and avoid churn from always choosing a "new" alternative.

**Behavior summary:**

- `ralph show prompt-guide [flags]` — Output the **full** [Writing Ralph prompts](docs/writing-ralph-prompts.md) guide **verbatim** (same content as the doc). No short summary; one source of truth. Supports optional `--markdown` (output is already markdown; flag can be used for consistency or future formatting). No config resolution; exit 0.

**Implementation: single source of truth, no doc in Go internals**

- **Source of truth:** `docs/writing-ralph-prompts.md` — authors and users use this path only.
- **How the binary gets the content:** Embed the guide in the binary so `ralph show prompt-guide` works from any working directory (e.g. after install). Go `//go:embed` can only embed files under the package directory. So at **build time**, copy `docs/writing-ralph-prompts.md` into a directory that the Go package can embed (e.g. `cmd/ralph/embed/writing-ralph-prompts.md` or a small package under `internal/`). The Makefile (or a pre-build step) runs this copy before `go build`. No symlinks; the copy is generated, not hand-maintained. The canonical file remains in `docs/`; users and README point to `docs/writing-ralph-prompts.md` only.

---

## Tasks

### T1. Register `show prompt-guide` subcommand and update show help — P1

**Dependencies:** None.

**Context:**

- In `cmd/ralph/main.go`, `showCmd()` currently adds only `showConfigCmd()`, `showPromptCmd()`, `showAliasCmd()`.
- Add a new subcommand so that `ralph show prompt-guide` is valid. The first argument after `show` must be one of: config, prompt, alias, **prompt-guide** (cli.md).
- Update the `show` root command’s `Long` description to include `show prompt-guide` (e.g. "Use 'show config', 'show prompt [name]', 'show alias [name]', or 'show prompt-guide'."). This satisfies O008/R005 discoverability via help.
- The subcommand should reject unexpected positional args (e.g. `ralph show prompt-guide foo` → error), consistent with `show config` which errors on args. No config resolution; no `--config` needed for this subcommand.

**Acceptance:**

- `ralph show prompt-guide` runs without error (implementation of output can be a stub in T2).
- `ralph show --help` lists `prompt-guide` and the Long text mentions it.
- `ralph show prompt-guide extra` exits non-zero with a clear error.

---

### T2. Implement verbatim guide output for `show prompt-guide` — P1

**Dependencies:** T1.

**Context:**

- Output must be the **full** [docs/writing-ralph-prompts.md](docs/writing-ralph-prompts.md) document **verbatim** (single source of truth). No short summary; the guide is the only definition of the four dimensions and their wording.
- Implementation: embed the guide in the binary. Because `//go:embed` only embeds files under the package (or subdirs), add a **build-time copy** of the doc into an embeddable location. For example: Makefile target or pre-build step copies `docs/writing-ralph-prompts.md` to e.g. `cmd/ralph/embed/writing-ralph-prompts.md`; the `show prompt-guide` command reads and prints that embedded content. The canonical file remains `docs/writing-ralph-prompts.md`; the copy is generated so users never look inside Go for the doc.
- No config resolution; do not call `config.Resolve`. Exit 0 on success.

**Acceptance:**

- `ralph show prompt-guide` prints to stdout the full content of docs/writing-ralph-prompts.md (byte-for-byte identical to the doc, or logically equivalent if any normalization is applied).
- Exit code 0.
- No config file or env required to run.

---

### T3. Add `--markdown` flag — P2

**Dependencies:** T2.

**Context:**

- cli.md: "Supports optional `--markdown` to emit markdown (e.g. for saving or piping to a pager)." Only for `show prompt-guide`.
- Add a boolean flag `--markdown` to the `show prompt-guide` subcommand. The output is already the full guide (markdown). When `--markdown` is set, output the same content (e.g. for saving or piping to a pager); when not set, output the same content. Flag is retained for CLI consistency and for scripts that want to request markdown explicitly.
- Help text for the flag should mention saving or piping to a pager (per cli.md).

**Acceptance:**

- `ralph show prompt-guide --markdown` prints the full guide (same as without the flag).
- `ralph show prompt-guide` (no flag) prints the full guide.
- `ralph show prompt-guide --help` documents `--markdown`.

---

### T4. Tests for `show prompt-guide` — P2

**Dependencies:** T2 (and T3 if testing --markdown).

**Context:**

- Add tests in the same style as existing show/config tests if present, or in `cmd/ralph` / appropriate test file. See [AGENTS.md](AGENTS.md): test command is `make test`.
- Cover: (1) `ralph show prompt-guide` exits 0 and stdout contains the full guide content — all four dimension names (signal and state, iteration awareness, scope and convergence, subjective completion), the summary table, and the closing line that references `ralph show prompt-guide`. (2) With `--markdown`, stdout is the same (full guide; contains markdown such as `##`, `**`). (3) Invalid usage (e.g. unexpected positional arg) exits non-zero.
- Tests may invoke the binary or test the command RunE in isolation; follow existing patterns in the repo. If the implementation embeds a build-time copy of the doc, tests can assert that key phrases from the guide appear in output rather than requiring the exact bytes of docs/writing-ralph-prompts.md at test time.

**Acceptance:**

- `make test` passes.
- At least one test covers success path (stdout contains the four dimensions and full-guide content) and one covers invalid-usage or flag behavior.

---

## Definition of done (then delete this plan)

- [ ] T1–T4 complete.
- [ ] Build copies docs/writing-ralph-prompts.md into embed location before go build (single source of truth in docs/).
- [ ] `make build` and `make test` pass.
- [ ] Manual smoke: `ralph show prompt-guide` outputs the full guide verbatim; `ralph show prompt-guide --markdown` same; `ralph show --help` mentions prompt-guide.
- [ ] Delete PLAN.md.
