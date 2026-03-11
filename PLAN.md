# Ralph Implementation Plan

This plan is derived from the **engineering specification** (`docs/engineering/README.md` and `docs/engineering/components/`), which is in turn derived from the **product documentation** (`docs/product/README.md` and outcome/requirement docs). Tasks are phased, dependency-ordered, and scoped for implementation. Each task includes links to relevant docs and a **beads priority** (0–4 per AGENTS.md). Use **bd** for tracking; create issues from this plan with dependencies so the dependency tree is reflected in bd.

**Status of this document:** This plan is a **working document** for the implementation effort. Once all tasks are completed (and the work is reflected in bd and the codebase), **PLAN.md should be removed**. The long-term source of truth is the product and engineering docs plus bd history and release notes; the plan itself is temporary scaffolding.

**Beads priorities:** 0 = Critical, 1 = High, 2 = Medium, 3 = Low, 4 = Backlog.

---

## Phase 0: Project foundation

Get a buildable, testable Go tree and align tooling with AGENTS.md. No product behavior yet.

| Task ID | Title | Scope | Docs | Deps | Priority |
|---------|--------|--------|------|------|----------|
| T0.1 | Go module and repo layout | Initialize `go.mod`, create `cmd/ralph/main.go` (minimal binary that exits 0), and `internal/` package layout. Ensure `make build` produces `bin/ralph`. | [Engineering README](docs/engineering/README.md), [AGENTS.md](AGENTS.md) Implementation Definition | — | 1 |
| T0.2 | Makefile and quality gates | Verify `make build`, `make test`, `make lint`, `make fmt` from repo root. Add or fix targets so CI and local quality gates match AGENTS.md. | [AGENTS.md](AGENTS.md) Build/Test/Lint | T0.1 | 1 |
| T0.3 | Create bd issues from plan | Create bd issues for each task in this plan with correct `-t` (feature/task/chore), `-p` (0–4), and dependency links (`--deps` block/dep or discovered-from as appropriate) so the dependency tree is in bd. | [AGENTS.md](AGENTS.md) Beads | T0.1 | 2 |

---

## Phase 1: Config component

Implement configuration layer resolution, effective config, and schema per engineering spec. Required by run-loop, CLI, and review.

| Task ID | Title | Scope | Docs | Deps | Priority |
|---------|--------|--------|------|------|----------|
| T1.1 | Config defaults and types | Define internal types for effective config (loop settings, prompts map, aliases map). Implement built-in defaults (max iterations, failure threshold, signal strings, log level, streaming, etc.) so the tool works without a config file. | [config](docs/engineering/components/config.md), [O002/R001](docs/product/O002-configurable-behavior/R001-config-layer-resolution.md), [O002/R002](docs/product/O002-configurable-behavior/R002-loop-behavior-configurable.md) | T0.2 | 1 |
| T1.2 | Config file locations (global, workspace) | Resolve paths for global config (`$RALPH_CONFIG_HOME/ralph-config.yml`, XDG/fallback) and workspace config (cwd). Skip missing files without error. Read YAML from those paths. | [config](docs/engineering/components/config.md) Layer order, [O002/R001](docs/product/O002-configurable-behavior/R001-config-layer-resolution.md) | T1.1 | 1 |
| T1.3 | Explicit config file and error on missing | When user supplies explicit config path (CLI), use only that file; do not load global/workspace. If file is missing or unreadable, return error (R005). | [config](docs/engineering/components/config.md), [O002/R005](docs/product/O002-configurable-behavior/R005-explicit-config-file-only.md) | T1.2 | 1 |
| T1.4 | Config schema (loop, prompts, aliases) | Parse and validate canonical YAML structure: `loop` (max_iterations, failure_threshold, timeout_seconds, success_signal, failure_signal, signal_precedence, preamble, streaming, log_level), `prompts` (path/content, optional loop overrides), `aliases` (command). Reject or error on invalid/out-of-range values. | [config](docs/engineering/components/config.md) Config file structure | T1.1 | 1 |
| T1.5 | Environment variable overlay | Apply `RALPH_CONFIG_HOME` and `RALPH_LOOP_*` env vars per config component table. Environment layer overrides file-based config; invalid values produce clear error. | [config](docs/engineering/components/config.md) Environment variables, [O010/R004](docs/product/O010-automation/R004-full-non-interactive-config.md) | T1.4 | 1 |
| T1.6 | Prompt-level overrides and merge order | Merge prompt-level loop overrides when resolving config for a given prompt. Final order: defaults → global → workspace → explicit file → env → prompt overrides → CLI (CLI wired in Phase 4). | [config](docs/engineering/components/config.md) Layer order, [O002/R003](docs/product/O002-configurable-behavior/R003-named-prompts-with-overrides.md) | T1.5 | 1 |
| T1.7 | Resolve effective config API | Expose a single function/entrypoint: resolve config for context (cwd, explicit path, env). Return effective config struct used by run-loop, review, list, show. | [config](docs/engineering/components/config.md) Interfaces, [O002/R007](docs/product/O002-configurable-behavior/R007-view-effective-config.md) | T1.6 | 1 |

---

## Phase 2: Backend component

Implement invocation of the AI CLI with prompt on stdin and stdout capture. Required by run-loop and review.

| Task ID | Title | Scope | Docs | Deps | Priority |
|---------|--------|--------|------|------|----------|
| T2.1 | Backend interface and exec contract | Define interface: invoke(command, promptBytes, cwd, env) → (stdout, exitCode, error). No shell; exec-style invocation. Stdin receives prompt; stream closed after write. Stdout captured in full. | [backend](docs/engineering/components/backend.md) Interfaces, [O003/R001](docs/product/O003-backend-agnosticism/R001-invoke-ai-cli-stdin-capture-stdout.md) | T0.2 | 1 |
| T2.2 | Inherit environment and working directory | Pass through process environment and cwd to child unless overridden by config (document behavior). No default overrides in this task. | [backend](docs/engineering/components/backend.md) Invocation contract, [O003/R002](docs/product/O003-backend-agnosticism/R002-inherit-env-and-cwd.md) | T2.1 | 1 |
| T2.3 | Built-in AI command aliases | Ship built-in alias definitions (claude, kiro, copilot, cursor-agent) per backend component table. Config layer resolves alias name → command string; backend receives resolved command. | [backend](docs/engineering/components/backend.md) Built-in AI command aliases, [O002/R004](docs/product/O002-configurable-behavior/R004-ai-command-aliases-configurable.md) | T1.7, T2.1 | 1 |
| T2.4 | Per-iteration timeout (backend or run-loop) | Support optional timeout for a single invocation. Either backend kills process after N seconds or run-loop passes timeout and backend respects it. Document placement. | [backend](docs/engineering/components/backend.md) Timeout, [O001](docs/product/O001-iterative-completion/README.md) | T2.1 | 2 |

---

## Phase 3: Run-loop component

Implement the iteration loop: validate command, load prompt once, invoke backend, detect signals, exit with correct codes and reports.

| Task ID | Title | Scope | Docs | Deps | Priority |
|---------|--------|--------|------|------|----------|
| T3.1 | Validate AI command before loop | Resolve AI command (alias or direct) and verify executable is available (e.g. on PATH). If missing or invalid, return clear error and do not start loop. Documented error exit code. | [run-loop](docs/engineering/components/run-loop.md), [O001/R001](docs/product/O001-iterative-completion/R001-validate-ai-command-before-loop.md), [O004/R001](docs/product/O004-observability/R001-clear-error-missing-ai-command.md) | T1.7, T2.3 | 1 |
| T3.2 | Load prompt once and buffer | Load prompt from resolved source (alias → file path, file path, or stdin). Buffer in memory. Fail before loop if source unavailable. Do not re-read between iterations. | [run-loop](docs/engineering/components/run-loop.md), [O001/R002](docs/product/O001-iterative-completion/R002-load-prompt-once-and-buffer.md) | T1.7 | 1 |
| T3.3 | Assemble prompt (preamble injection) | Build assembled prompt: optional preamble (e.g. iteration count, context) + buffered prompt content. Preamble configurable via config/CLI; no preamble written to user file. | [run-loop](docs/engineering/components/run-loop.md), [O002/R002](docs/product/O002-configurable-behavior/R002-loop-behavior-configurable.md) | T3.2 | 1 |
| T3.4 | Success signal detection and exit 0 | Scan captured stdout for configured success signal (substring or pattern). On match: emit completion message, iteration count, timing; exit with documented success code (e.g. 0). | [run-loop](docs/engineering/components/run-loop.md), [O001/R004](docs/product/O001-iterative-completion/R004-detect-success-signal-exit-zero.md), [O004/R002](docs/product/O004-observability/R002-success-report-and-exit-zero.md) | T2.1, T3.3 | 1 |
| T3.5 | Failure signal and consecutive-failure threshold | Detect configured failure signal; increment consecutive-failure count; if count ≥ threshold, report and exit with distinct failure-threshold code. Otherwise start next iteration. | [run-loop](docs/engineering/components/run-loop.md), [O001/R005](docs/product/O001-iterative-completion/R005-detect-failure-signal-continue-or-exit.md), [O004/R003](docs/product/O004-observability/R003-failure-threshold-report-and-exit-code.md) | T3.4 | 1 |
| T3.6 | Signal precedence (static) | When both success and failure signals appear in same output, apply static precedence (e.g. first match wins or defined order). Document behavior. | [run-loop](docs/engineering/components/run-loop.md), [O001/R006](docs/product/O001-iterative-completion/R006-signal-precedence.md) | T3.4, T3.5 | 1 |
| T3.7 | Max iterations exit | When iteration count reaches max iterations without success, report and exit with distinct max-iterations code. | [run-loop](docs/engineering/components/run-loop.md), [O001/R007](docs/product/O001-iterative-completion/R007-exit-on-max-iterations.md), [O004/R004](docs/product/O004-observability/R004-max-iterations-report-and-exit-code.md) | T3.5 | 1 |
| T3.8 | Process exit without signal (crash/kill) | If AI process exits without success or failure signal, treat iteration as failure; increment consecutive-failure count; continue or exit per threshold. Distinguish where documented. | [run-loop](docs/engineering/components/run-loop.md), [O001/R009](docs/product/O001-iterative-completion/R009-process-exit-without-signal.md) | T3.5 | 1 |
| T3.9 | Interrupt (SIGINT/SIGTERM) distinct exit code | On user interrupt, exit with distinct documented code (e.g. 130). No ambiguous exit. | [run-loop](docs/engineering/components/run-loop.md), [O004/R005](docs/product/O004-observability/R005-distinct-exit-code-on-interrupt.md) | T3.7 | 1 |
| T3.10 | Dry-run: print assembled prompt, no backend | When dry-run is set, assemble prompt (with preamble if enabled), print to stdout or logs, exit 0. Do not invoke backend. | [run-loop](docs/engineering/components/run-loop.md), [O004/R007](docs/product/O004-observability/R007-dry-run-shows-assembled-prompt.md) | T3.3 | 2 |
| T3.11 | AI-interpreted signal precedence (optional) | When configured, both signals present: invoke AI once with built-in Ralph prompt to interpret output and decide success/failure; if unclear, apply defined fallback. | [run-loop](docs/engineering/components/run-loop.md), [O001/R008](docs/product/O001-iterative-completion/R008-ai-interpreted-signal-precedence.md) | T3.6, T2.1 | 2 |
| T3.12 | Iteration statistics | After multi-iteration run, report statistics (e.g. min/max/mean duration per iteration) when configured. | [run-loop](docs/engineering/components/run-loop.md), [O004/R008](docs/product/O004-observability/R008-iteration-statistics.md) | T3.7 | 2 |
| T3.13 | Log level and show AI output (streaming) | Respect log level for Ralph’s logs. When “show AI output” is true, stream AI stdout to terminal while still capturing for signal scan. Quiet = minimal log + no stream. | [run-loop](docs/engineering/components/run-loop.md), [O004/R006](docs/product/O004-observability/R006-log-level-and-show-ai-output.md) | T2.1, T3.4 | 2 |

---

## Phase 4: CLI — run path

Command parsing, global options, and `ralph run` with full flag set and dispatch to run-loop.

| Task ID | Title | Scope | Docs | Deps | Priority |
|---------|--------|--------|------|------|----------|
| T4.1 | CLI skeleton and top-level commands | Implement top-level commands: run, review, list, show, version. Unknown command → error stderr, non-zero exit, suggest `--help`. Help for each command. | [cli](docs/engineering/components/cli.md) Top-level commands, [O008/R003](docs/product/O008-discoverability/R003-list-and-help.md) | T0.2 | 1 |
| T4.2 | Global options: --config, --help, --version | Parse `--config <path>`, `--help`/`-h`, `--version`. Pass explicit config path to config resolution. No short form for `--config` (reserve `-c` for run context). | [cli](docs/engineering/components/cli.md) Global options | T4.1 | 1 |
| T4.3 | ralph run: prompt source (alias, file, stdin) | Exactly one of: positional alias, `-f`/`--file`, or stdin. Mutual exclusion; error if combined or TTY stdin with no source. | [cli](docs/engineering/components/cli.md) ralph run, [O001/R003](docs/product/O001-iterative-completion/R003-run-via-alias-file-or-stdin.md) | T4.2, T1.7 | 1 |
| T4.4 | ralph run: loop and AI flags | Parse `--max-iterations`/`-n`, `--unlimited`/`-u`, `--failure-threshold`, `--iteration-timeout`, `--max-output-buffer`, `--no-preamble`, `--dry-run`/`-d`, `--ai-cmd`, `--ai-cmd-alias`, `--signal-success`, `--signal-failure`, `--signal-precedence`, `-c`/`--context` (repeatable). Override effective config for the run. | [cli](docs/engineering/components/cli.md) ralph run Flags | T4.3 | 1 |
| T4.5 | ralph run: output and observability flags | Parse `--verbose`/`-v`, `--quiet`/`-q`, `--log-level`, `--stream`, `--no-stream`. Wire to run-loop observability. | [cli](docs/engineering/components/cli.md) ralph run | T4.4, T3.13 | 1 |
| T4.6 | Wire CLI run to config + run-loop | Resolve config (cwd, `--config`, env), resolve prompt source, call run-loop with effective config and overrides. Propagate exit code from run-loop. | [cli](docs/engineering/components/cli.md) Interfaces, [O010/R001](docs/product/O010-automation/R001-non-interactive-completion.md) | T4.5, T1.7, T3.9 | 1 |
| T4.7 | RALPH_CONFIG_HOME and run error handling | Honor `RALPH_CONFIG_HOME` for global config path. Run path: missing/invalid AI command, missing explicit config file, invalid flag values → clear error, non-zero exit. | [cli](docs/engineering/components/cli.md) Environment variables, Error handling | T4.6 | 1 |

---

## Phase 5: Review component

Implement `ralph review`: report with narrative and machine-parseable summary, suggested revision, apply with confirmation, exit codes.

| Task ID | Title | Scope | Docs | Deps | Priority |
|---------|--------|--------|------|------|----------|
| T5.1 | Review invocation inputs (alias, file, stdin) | Accept same prompt sources as run: alias, file path, stdin. Exactly one per invocation. Invalid or missing → error before review, exit 2. | [review](docs/engineering/components/review.md), [O005/R001](docs/product/O005-prompt-review/R001-review-invocation-inputs.md) | T1.7, T2.1 | 1 |
| T5.2 | Report content: narrative + summary line + revision | Produce report file with: (1) narrative feedback, (2) single machine-parseable summary line per spec (`ralph-review: status=ok|errors|warnings` with optional errors=/warnings=), (3) full suggested revision. | [review](docs/engineering/components/review.md) Report format, [O005/R002](docs/product/O005-prompt-review/R002-report-content-and-format.md), [O005/R003](docs/product/O005-prompt-review/R003-suggested-revision.md) | T5.1 | 1 |
| T5.3 | Report to file and revision output path | Write report to user-chosen path (or default). Support `--report` and `--prompt-output`. When prompt from stdin and `--apply` set, require `--prompt-output`; else error exit 2. | [review](docs/engineering/components/review.md), [O005/R005](docs/product/O005-prompt-review/R005-report-to-file.md), [O005/R006](docs/product/O005-prompt-review/R006-revision-output-path.md) | T5.2 | 1 |
| T5.4 | Apply with confirmation and non-interactive | When user requests apply: write revision to chosen path. Interactive: confirm before overwrite unless `--yes`. Non-interactive: if confirmation would be needed and `--yes` not set, exit 2 with clear message. | [review](docs/engineering/components/review.md) Apply and confirmation, [O005/R004](docs/product/O005-prompt-review/R004-apply-with-confirmation.md), [O009/R001](docs/product/O009-predictability/R001-explicit-apply-for-writes.md), [O009/R003](docs/product/O009-predictability/R003-review-apply-separation-and-confirmation.md) | T5.3 | 1 |
| T5.5 | Review exit code derivation (0, 1, 2) | After report file is written, parse machine-parseable summary: status=ok and errors=0 → 0; status=errors or errors≥1 → 1; missing/malformed → 1 (fail-safe). Report write failure or apply/precondition failure → 2. | [review](docs/engineering/components/review.md) Exit code derivation, [O005/R008](docs/product/O005-prompt-review/R008-exit-codes.md), [O004/R009](docs/product/O004-observability/R009-review-exit-code-and-report.md), [O010/R003](docs/product/O010-automation/R003-machine-parseable-review-summary.md) | T5.2, T5.3 | 1 |
| T5.6 | Evaluation dimensions (signal, state, scope, convergence) | Review evaluates prompt on: signal/state discipline, iteration awareness, scope and convergence, subjective completion criteria. Embedded reviewer instructions (Ralph-owned); no user-editable review prompt. | [review](docs/engineering/components/review.md), [O005/R007](docs/product/O005-prompt-review/R007-evaluation-dimensions.md) | T5.1 | 2 |
| T5.7 | Summary parser (internal/review/summary.go) | Implement parser for `ralph-review: status=...` line to derive exit code. Used by review component and documented for CI. | [review](docs/engineering/components/review.md) Report format, [AGENTS.md](AGENTS.md) Implementation Definition | T5.5 | 2 |

---

## Phase 6: CLI — review, list, show, version

Complete CLI surface for review, list, show, and version.

| Task ID | Title | Scope | Docs | Deps | Priority |
|---------|--------|--------|------|------|----------|
| T6.1 | ralph review subcommand and flags | Syntax: `ralph review [alias]`, `ralph review -f <path>`, `ralph review` (stdin). Flags: `--file`/`-f`, `--report`, `--prompt-output`, `--apply`, `--yes`/`-y`, `--quiet`, `--log-level`, `--config`. Wire to review component. | [cli](docs/engineering/components/cli.md) ralph review | T4.1, T5.5 | 1 |
| T6.2 | ralph list prompts and aliases | `ralph list`, `ralph list prompts`, `ralph list aliases`. Output names (and optional display name, description, path for prompts). Use resolved config. | [cli](docs/engineering/components/cli.md) ralph list, [O002/R006](docs/product/O002-configurable-behavior/R006-list-prompts-and-commands.md) | T4.1, T1.7 | 1 |
| T6.3 | ralph show config | `ralph show config` outputs effective config. Optional `--provenance` to include which layer supplied each value. | [cli](docs/engineering/components/cli.md) ralph show, [O002/R007](docs/product/O002-configurable-behavior/R007-view-effective-config.md) | T1.7, T4.1 | 2 |
| T6.4 | ralph show prompt and show alias | `ralph show prompt [name]`, `ralph show alias [name]` with defined behavior when name omitted. Error on unknown name. | [cli](docs/engineering/components/cli.md) ralph show | T6.2 | 2 |
| T6.5 | ralph version | Print version string to stdout, exit 0. Version from build (e.g. ldflags). | [cli](docs/engineering/components/cli.md) ralph version | T4.1 | 1 |

---

## Phase 7: Observability and polish

Log level, streaming, iteration stats, and any run/review output formatting refinements.

| Task ID | Title | Scope | Docs | Deps | Priority |
|---------|--------|--------|------|------|----------|
| T7.1 | Document run and review exit codes | Publish exact numeric exit codes for run (success, failure threshold, max iterations, interrupt, error) and review (0, 1, 2) in user docs and automation docs. | [run-loop](docs/engineering/components/run-loop.md), [review](docs/engineering/components/review.md), [O010/R002](docs/product/O010-automation/R002-documented-stable-exit-codes.md) | T4.7, T6.1 | 2 |
| T7.2 | Quiet and verbose behavior consistency | Ensure `--quiet` and `--verbose` behave per spec across run and review; no interactive prompts when non-interactive + appropriate flags. | [cli](docs/engineering/components/cli.md), [O004/R006](docs/product/O004-observability/R006-log-level-and-show-ai-output.md) | T4.5, T6.1 | 2 |
| T7.3 | Provenance in show config (optional) | If not done in T6.3: implement `--provenance` for `ralph show config` so each value can show layer (default, global, workspace, explicit, env, cli, prompt). | [config](docs/engineering/components/config.md), [O002/R007](docs/product/O002-configurable-behavior/R007-view-effective-config.md) | T6.3 | 3 |

---

## Phase 8: Documentation and release

User-facing docs, install/uninstall/upgrade procedures, discoverability, and release contract.

| Task ID | Title | Scope | Docs | Deps | Priority |
|---------|--------|--------|------|------|----------|
| T8.1 | Doc coverage: subcommands, config, exit codes | Documentation covers: run, review, list, show, version; config layers and schema; prompt sources; loop behavior; review report and apply; exit codes (run + review); non-interactive use. | [documentation](docs/engineering/components/documentation.md), [O007/R001](docs/product/O007-user-documentation/R001-doc-coverage.md) | T6.5, T7.1 | 2 |
| T8.2 | Doc discoverability and accuracy | Single place or linked set for CLI, flags, config, env vars. Docs match actual behavior; update when implementation contract changes. | [documentation](docs/engineering/components/documentation.md), [O007/R002](docs/product/O007-user-documentation/R002-doc-discoverability.md), [O007/R003](docs/product/O007-user-documentation/R003-doc-accuracy.md) | T8.1 | 2 |
| T8.3 | Install and uninstall procedures | Documented install steps (and/or scripts) so Ralph is invocable after install. Documented uninstall so binary and install state removed; no broken PATH. | [documentation](docs/engineering/components/documentation.md), [O006](docs/product/O006-install-uninstall/README.md) | T8.1 | 2 |
| T8.4 | Upgrade and backward compatibility | Document upgrade to chosen version and update within non-breaking version. Backward compatibility and migration (when needed) documented; release notes for changes. | [documentation](docs/engineering/components/documentation.md), [O011](docs/product/O011-update-upgrade/README.md) | T8.2 | 2 |
| T8.5 | Discoverability: what and why, first run | Content so a new user can find what Ralph is, why it exists, install steps, and a first command that completes successfully. Path to first successful run documented. | [documentation](docs/engineering/components/documentation.md), [O008](docs/product/O008-discoverability/README.md) | T8.2 | 2 |
| T8.6 | Troubleshooting | Document common problems (e.g. prompt not found, wrong exit code, config not found) and how to resolve them. | [documentation](docs/engineering/components/documentation.md), [O007/R004](docs/product/O007-user-documentation/R004-troubleshooting.md) | T8.2 | 3 |
| T8.7 | Release notes and stable contract | Release notes for each release; documented stable contract (exit codes, summary format, config) so scripts and CI can rely on them. | [documentation](docs/engineering/components/documentation.md), [O010/R002](docs/product/O010-automation/R002-documented-stable-exit-codes.md), [O011/R004](docs/product/O011-update-upgrade/R004-release-notes-for-changes.md) | T7.1, T8.4 | 2 |

---

## Dependency graph (summary)

- **Phase 0:** T0.1 → T0.2; T0.1 → T0.3.
- **Phase 1:** T0.2 → T1.1 → T1.2 → T1.3; T1.1 → T1.4 → T1.5 → T1.6 → T1.7.
- **Phase 2:** T0.2 → T2.1 → T2.2; T1.7, T2.1 → T2.3; T2.1 → T2.4.
- **Phase 3:** T1.7, T2.3 → T3.1; T1.7 → T3.2 → T3.3; T2.1, T3.3 → T3.4 → T3.5 → T3.6, T3.7, T3.8; T3.7 → T3.9; T3.3 → T3.10; T3.6, T2.1 → T3.11; T3.7 → T3.12; T2.1, T3.4 → T3.13.
- **Phase 4:** T0.2 → T4.1 → T4.2 → T4.3 → T4.4 → T4.5 → T4.6 → T4.7; T1.7, T3.9, T3.13 feed into T4.6/T4.5.
- **Phase 5:** T1.7, T2.1 → T5.1 → T5.2 → T5.3 → T5.4, T5.5; T5.1 → T5.6; T5.5 → T5.7.
- **Phase 6:** T4.1, T5.5 → T6.1; T4.1, T1.7 → T6.2; T1.7, T4.1 → T6.3; T6.2 → T6.4; T4.1 → T6.5.
- **Phase 7:** T4.7, T6.1 → T7.1; T4.5, T6.1 → T7.2; T6.3 → T7.3.
- **Phase 8:** T6.5, T7.1 → T8.1 → T8.2; T8.1 → T8.3; T8.2 → T8.4, T8.5, T8.6; T7.1, T8.4 → T8.7.

---

## How to use this plan with bd

1. Create an epic or parent issue for "Ralph implementation per engineering spec" if desired.
2. For each task T0.1–T8.7, run:
   - `bd create "T0.1: Go module and repo layout" --description="... (paste Scope + Docs)" -t task -p 1 --json`
   - For tasks with Deps, add `--deps block:T0.1` (or the bd-id of the dependency) when creating the dependent issue.
3. Use `bd ready` to pick unblocked work; implement, then close with `bd close <id> --reason "Done"`.
4. When adding discovered work, link with `--deps discovered-from:<parent-id>`.

This plan is the single reference for phased implementation; bd holds the live dependency tree and status.
