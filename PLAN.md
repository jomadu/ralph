# Ralph Implementation Plan

This plan turns the intent specification (`docs/intent/`) into well-scoped, prioritized, bite-sized tasks. Implementation is assumed to be **Go** (per spec: cobra, yaml, exec semantics). Tasks are ordered by dependency; later phases assume earlier phases are done.

---

## Principles

- **Bite-sized:** Each task is completable in one or two sessions and verifiable (manual or, when tests exist, automated).
- **Well-scoped:** Each task names the requirement(s) it implements and has clear done criteria from the spec.
- **Prioritized:** P0 = unblocks others / critical path; P1 = core user value; P2 = polish or secondary surface.

---

## Phases Overview

| Phase | Focus | Outcome |
|-------|--------|---------|
| **0** | Bootstrap | Go project, CLI skeleton, version subcommand |
| **1** | Config load & merge | Config struct, file loading, provenance, validation, env overlay |
| **2** | CLI & flags | Cobra commands and flags wired to config merge |
| **3** | Prompt & AI command | Prompt modes, fail-fast, command parsing, aliases, precedence, errors, env inheritance |
| **4** | Loop core | Spawn, buffer, signal scan, crash handling, limits, interruption, timeout |
| **5** | Observability | Exit codes, progress, statistics, verbose, log level, dry-run |
| **6** | Polish | Per-prompt overrides, list commands, provenance display, env/docs, full CLI |

---

## Task List

### Phase 0 — Bootstrap

| ID | Task | Spec(s) | Priority | Dependencies |
|----|------|--------|----------|--------------|
| 0.1 | Create Go module and directory layout (`cmd/ralph`, `internal/`) | — | P0 | — |
| 0.2 | Implement `ralph version` subcommand (version string to stdout, exit 0) | O4/R7 | P0 | 0.1 |
| 0.3 | Add root command and `run` / `list` / `version` subcommands (no behavior yet) | O2/R9 | P0 | 0.1 |

### Phase 1 — Config Load & Merge

| ID | Task | Spec(s) | Priority | Dependencies |
|----|------|--------|----------|--------------|
| 1.1 | Define config struct (loop, prompts, ai_cmd_aliases) with YAML tags and built-in defaults | O2/R1, R3, R6 | P0 | 0.1 |
| 1.2 | Implement config file discovery: global (RALPH_CONFIG_HOME → XDG → ~/.config/ralph), workspace (./ralph-config.yml) | O2/R5 | P0 | 1.1 |
| 1.3 | Implement silent skip for absent default config files; error if `--config <path>` missing | O2/R5 | P0 | 1.2 |
| 1.4 | Implement config merge with provenance (default → global → workspace or file → env → CLI); tag each value with layer | O2/R1 | P0 | 1.3 |
| 1.5 | Implement RALPH_* environment variable overlay for all supported keys | O2/R8 | P0 | 1.4 |
| 1.6 | Add unknown-key detection (strict decode vs permissive) and warn at load; do not block startup | O2/R2 | P1 | 1.4 |
| 1.7 | Implement config validation (schema + semantic: ai_cmd_alias resolution); fail-fast, collect all errors, exit 1 | O2/R3 | P0 | 1.4, 1.5 |

### Phase 2 — CLI & Flags

| ID | Task | Spec(s) | Priority | Dependencies |
|----|------|--------|----------|--------------|
| 2.1 | Wire `--config` and config load into run path; merge CLI flags after env in config resolution | O2/R1, R5, R9 | P0 | 1.7, 0.3 |
| 2.2 | Add `ralph run` flags: loop control (--max-iterations, --unlimited, --failure-threshold, --iteration-timeout, --max-output-buffer, --preamble/--no-preamble, --dry-run), signals (--signal-success, --signal-failure), context (--context repeatable) | O1/R3–R8, O2/R9, O4/R4 | P0 | 2.1 |
| 2.3 | Add `ralph run` flags: AI command (--ai-cmd, --ai-cmd-alias), output (--verbose, --quiet, --log-level) | O3/R6, O4/R3, R5, O2/R9 | P0 | 2.1 |
| 2.4 | Implement prompt input mode resolution: positional alias vs --file/-f vs stdin; error if ambiguous or none | O1/R9 | P0 | 2.1 |

### Phase 3 — Prompt & AI Command

| ID | Task | Spec(s) | Priority | Dependencies |
|----|------|--------|----------|--------------|
| 3.1 | Implement prompt source resolution and read-once buffer (alias → path → read; -f → read; stdin → read) | O1/R9 | P0 | 2.4, 1.4 |
| 3.2 | Fail-fast on invalid prompt source: unknown alias, missing/unreadable/empty file, empty stdin; clear errors, exit 1 | O2/R4 | P0 | 3.1 |
| 3.3 | Implement shell-style command parsing (quoted args, no shell); produce argv for exec | O3/R2 | P0 | — |
| 3.4 | Implement built-in AI command aliases (claude, kiro, copilot, cursor-agent) and merge with user ai_cmd_aliases | O3/R1, R3 | P0 | 1.1, 3.3 |
| 3.5 | Implement command source precedence (ai_cmd vs ai_cmd_alias, layer order); no default for ai_cmd/ai_cmd_alias | O3/R6 | P0 | 3.4, 2.3 |
| 3.6 | Alias resolution errors: unknown alias (name + list available), no command configured; exit 1 before loop | O3/R5 | P0 | 3.5 |
| 3.7 | Spawn AI process with inherited environment and current working directory; no filtering of env | O3/R4 | P0 | 3.3 |

### Phase 4 — Loop Core

| ID | Task | Spec(s) | Priority | Dependencies |
|----|------|--------|----------|--------------|
| 4.1 | Assemble prompt: preamble (R8) + buffered content; pipe to process stdin; capture stdout+stderr to bounded buffer (discard oldest), per-iteration buffer | O1/R8, R6 | P0 | 3.1, 3.7 |
| 4.2 | After process exit: signal scan on buffer (substring, failure wins over success); drive iteration outcome (success / failure / no-signal) | O1/R2 | P0 | 4.1 |
| 4.3 | Process crash handling: non-zero exit = crash; still scan buffer; same loop logic; one iteration; no retry | O1/R1 | P0 | 4.2 |
| 4.4 | Max iteration limit: check before each iteration; exit 2 on exhaustion; support --unlimited | O1/R4 | P0 | 4.2 |
| 4.5 | Consecutive failure counter: increment on failure, reset on success/no-signal; exit 1 when threshold reached | O1/R5 | P0 | 4.2 |
| 4.6 | Graceful interruption: SIGINT/SIGTERM → forward SIGTERM to child, 5s grace, then SIGKILL; discard output, exit 130 | O1/R7 | P0 | 4.1 |
| 4.7 | Per-iteration timeout: optional timer, SIGTERM then 5s then SIGKILL; scan partial output | O1/R3 | P1 | 4.1 |

### Phase 5 — Observability

| ID | Task | Spec(s) | Priority | Dependencies |
|----|------|--------|----------|--------------|
| 5.1 | Exit codes: 0 success, 1 failure threshold, 2 max iterations, 130 interrupt | O4/R1 | P0 | 4.2, 4.4, 4.5, 4.6 |
| 5.2 | Iteration progress: emit "Iteration N/M" or "N (unlimited)" at start of each iteration at info level to stderr | O4/R6 | P1 | 4.1 |
| 5.3 | Iteration statistics at completion: count, min/max/mean, stddev (Welford); stderr, info level; not on exit 130 | O4/R2 | P1 | 4.1 |
| 5.4 | Verbose mode: when show_ai_output true or -v, mirror child stdout/stderr to terminal while capturing | O4/R3 | P1 | 4.1 |
| 5.5 | Log level control: debug/info/warn/error, default info; --quiet = error; all Ralph log to stderr | O4/R5 | P1 | 5.2, 5.3 |
| 5.6 | Dry-run: resolve config + prompt, assemble iteration 1 (preamble + content), print to stdout, exit 0; no AI process | O4/R4 | P1 | 3.1, 4.1 (assembly only) |

### Phase 6 — Polish

| ID | Task | Spec(s) | Priority | Dependencies |
|----|------|--------|----------|--------------|
| 6.1 | Per-prompt loop overrides: merge prompts.<name>.loop.* into effective config for that alias | O2/R6 | P1 | 1.4, 2.2 |
| 6.2 | `ralph list prompts`: YAML output, sorted by alias; path always, name/description if set | O2/R7 | P1 | 1.4 |
| 6.3 | `ralph list aliases`: YAML output, built-in + user, sorted; resolved command per alias | O2/R7, O3/R1 | P1 | 3.4 |
| 6.4 | Provenance in debug log and in dry-run config section | O2/R1 | P2 | 1.4, 5.6 |
| 6.5 | Document and enforce RALPH_* env var set (single authoritative list per O2/R8) | O2/R8 | P2 | 1.5 |
| 6.6 | Ensure full CLI surface matches O2/R9 (all commands and flags discoverable via help) | O2/R9 | P2 | 2.2, 2.3, 6.2, 6.3 |

---

## Dependency Graph (Critical Path)

```
0.1 → 0.2, 0.3
1.1 → 1.2 → 1.3 → 1.4 → 1.5, 1.6, 1.7
2.1 ← 1.7, 0.3
2.2, 2.3, 2.4 ← 2.1
3.1 ← 2.4, 1.4
3.2 ← 3.1
3.3 → 3.4 → 3.5 → 3.6
3.7 ← 3.3
4.1 ← 3.1, 3.7
4.2 ← 4.1
4.3, 4.4, 4.5 ← 4.2
4.6 ← 4.1
4.7 ← 4.1
5.1 ← 4.2, 4.4, 4.5, 4.6
5.2, 5.3, 5.4, 5.5 ← 4.1
5.6 ← 3.1, assembly
6.1 ← 1.4, 2.2
6.2, 6.3 ← 1.4, 3.4
6.4, 6.5, 6.6 ← various
```

---

## Suggested bd Issue Titles

Use these with `bd create "Title" --description="..." -t task -p <0|1|2> --json`. Link dependencies with `--deps blockee:bd-<id>` when the blocking issue exists.

**Phase 0**
- `Bootstrap: Go module and directory layout (cmd/ralph, internal/)`
- `Implement ralph version subcommand (O4/R7)`
- `Add root and run/list/version subcommands (O2/R9 skeleton)`

**Phase 1**
- `Config struct with YAML tags and built-in defaults (O2)`
- `Config file discovery: global and workspace paths (O2/R5)`
- `Silent skip absent config; error for missing --config file (O2/R5)`
- `Config merge with provenance tags (O2/R1)`
- `RALPH_* environment variable overlay (O2/R8)`
- `Unknown config key warnings at load (O2/R2)`
- `Config validation: schema and semantic, fail-fast (O2/R3)`

**Phase 2**
- `Wire --config and CLI flags into config merge (O2/R1, R5, R9)`
- `ralph run loop and output flags (O1, O4)`
- `ralph run AI command and log flags (O3, O4)`
- `Prompt input mode resolution: alias vs -f vs stdin (O1/R9)`

**Phase 3**
- `Prompt source read-once buffer for alias, file, stdin (O1/R9)`
- `Fail-fast on invalid prompt source (O2/R4)`
- `Shell-style command parsing for exec (O3/R2)`
- `Built-in and user AI command aliases (O3/R1, R3)`
- `Command source precedence ai_cmd vs ai_cmd_alias (O3/R6)`
- `Alias resolution errors: unknown alias, no command (O3/R5)`
- `Process environment and cwd inheritance (O3/R4)`

**Phase 4**
- `Loop: assemble prompt, spawn process, capture output buffer (O1/R8, R6)`
- `Signal scan and iteration outcome (O1/R2)`
- `Process crash recovery and loop continuity (O1/R1)`
- `Max iteration limit and --unlimited (O1/R4)`
- `Consecutive failure tracking and exit 1 (O1/R5)`
- `Graceful SIGINT/SIGTERM handling, exit 130 (O1/R7)`
- `Per-iteration timeout (O1/R3)`

**Phase 5**
- `Distinct exit codes 0/1/2/130 (O4/R1)`
- `Iteration progress messages at info (O4/R6)`
- `Iteration statistics at completion (O4/R2)`
- `Verbose AI output streaming (O4/R3)`
- `Log level control and stderr (O4/R5)`
- `Dry-run: print assembled prompt, no AI process (O4/R4)`

**Phase 6**
- `Per-prompt loop overrides (O2/R6)`
- `ralph list prompts (O2/R7)`
- `ralph list aliases (O2/R7)`
- `Provenance in debug and dry-run (O2/R1)`
- `RALPH_* env var reference complete (O2/R8)`
- `CLI surface and help per O2/R9 (O2/R9)`

---

## How to Use This Plan

1. **Create issues:** Run `bd create "Title" --description="Spec refs and scope from PLAN.md" -t task -p 0|1|2 --json` for each task; add `--deps blockee:bd-<id>` when a blocking issue exists.
2. **Pick work:** Use `bd ready --json` and claim tasks in dependency order (Phase 0 → 1 → …).
3. **Update plan:** When the implementation evolves, adjust this plan and sync bd issues (close, create, or update).
4. **Session end:** Per AGENTS.md, file any new work as bd issues, update status, sync, push.

---

## Next Steps

1. Run `bd onboard` if beads is not yet initialized.
2. Create bd issues for **Phase 0** (0.1, 0.2, 0.3); no dependencies.
3. Create issues for **Phase 1** with deps on 0.1 where needed.
4. Start implementation with task **0.1** (Go module and layout), then **0.2** (version), then **0.3** (subcommands).
5. After Phase 0, add build/test/lint commands to AGENTS.md when they exist.
