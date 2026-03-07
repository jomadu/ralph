# Plan: Output and Verbosity (TASK.md) → Intent + Implementation

This plan aligns the intent tree and implementation with the desired output and verbosity behavior in **TASK.md**. Tasks are scoped so each can be completed with only the context in this plan and the referenced files.

**Document status:** TASK.md and this plan (PLAN.md) are **working documents** for this output/verbosity change. Once all tasks T1–T8 are completed and verified, both files should be **deleted**. The source of truth then lives in the intent tree (`docs/intent/`) and the implementation; beads issues may be closed but remain in history.

---

## Context (read once)

**Source of truth:** [TASK.md](TASK.md) — output destination, two axes (log level + show AI command output), defaults, flags, precedence.

**Desired state in short:**
- **stdout** = run log (Ralph operational messages + AI command stream). **stderr** = fatal/startup errors only.
- **Defaults:** log level = info, show AI command output = **true**.
- **Flags:** `-q` (quiet log + no AI output), `-v` (verbose log + AI output unless `--no-ai-cmd-output`), `--log-level <level>`, `--no-ai-cmd-output`.
- **Precedence:** Log level: `--log-level` else `-q` else `-v` else config/env else info. Show AI output: `--no-ai-cmd-output` → false; `-q` → false; `-v` → true; else config/env else **true**.

**Gap (current intent vs TASK.md):** Intent currently says Ralph operational output and AI mirror go to **stderr** (or “implementation chooses”); default show_ai_output is **false**; no `--no-ai-cmd-output`; `-q` does not suppress AI output in spec. All of that must be updated to match TASK.md.

---

## Task dependency overview

```
Intent (spec) — must complete before implementation
──────────────────────────────────────────────────
T1  O4 README (output destination)          [P1]  (none)
T2  O4 R3 (default, flag, precedence, stdout) [P1]  T1
T3  O4 R5 (stdout, -q dual effect)           [P1]  T1
T4  O4 R2 + R6 (stats/progress → stdout)     [P2]  T3
T5  O2 R9 (CLI table: --no-ai-cmd-output)    [P2]  T2, T3
T6  O2 R8 (env default true)                 [P2]  T2

Implementation — after intent locked
────────────────────────────────────
T7  Output routing (stdout/stderr)            [P3]  T1,T2,T3,T4
T8  Default + flags (--no-ai-cmd-output, -q/-v) [P3] T2,T3,T5
```

**Priority:** P1 = must do first (foundation). P2 = next (dependent on P1). P3 = implementation after spec is locked.

---

## Tasks (well-scoped, prioritized, dependency-linked)

### T1 — O4 README: Output destination and non-outcomes  
**Priority:** 1 · **Depends on:** —  

**File:** `docs/intent/O4-observability/README.md`

**What to do:**  
In the **Non-outcomes** section, replace the bullet that says “Output goes to stderr/stdout for the current invocation only” with:

- Operational messages and the AI command stream go to **stdout** (the run’s log). **stderr** is reserved for fatal/startup errors only. Ralph does not provide persistent log files.

**Acceptance:**  
- That bullet is the only place in the Non-outcomes section that defines where output goes.  
- It explicitly states: run log → stdout; stderr → fatal/startup only; no persistent log files.

---

### T2 — O4 R3: Default true, --no-ai-cmd-output, precedence, mirror to stdout  
**Priority:** 1 · **Depends on:** T1  

**File:** `docs/intent/O4-observability/R3-verbose-streaming.md`

**What to do:**  
1. **Default:** Change the default for “show AI command output” from **false** to **true** everywhere it is stated (control resolution, edge cases, acceptance criteria).  
2. **Flag:** Document **`--no-ai-cmd-output`**: when set, streaming is **false** (user explicitly disables).  
3. **Precedence (show AI output), highest wins:**  
   - `--no-ai-cmd-output` → false  
   - `-q` / `--quiet` → false  
   - `-v` / `--verbose` → true  
   - config / env  
   - default **true**  
4. **Mirror destination:** Where the spec says Ralph mirrors the child’s stdout/stderr “to the terminal” or “Ralph’s stderr or stdout, or both — implementation chooses”, replace with: Ralph mirrors to **stdout** (the run’s log).  
5. **Edge cases / acceptance:** Add or update: default → streamed; `-q` → not streamed; `--no-ai-cmd-output` → not streamed; `-v --no-ai-cmd-output` → verbose logs, AI output not streamed.

**Acceptance:**  
- Default is true; `--no-ai-cmd-output` and `-q` both force false; precedence matches TASK.md.  
- Mirror destination is explicitly stdout.

---

### T3 — O4 R5: Log output to stdout, -q also suppresses AI output  
**Priority:** 1 · **Depends on:** T1  

**File:** `docs/intent/O4-observability/R5-log-level-control.md`

**What to do:**  
1. **Output destination:** Change “All Ralph operational log output … goes to **stderr**” to **stdout**. Update every edge case, example, and acceptance criterion that says “stderr” for Ralph log output to “stdout”.  
2. **-q behavior:** State explicitly that **`-q` / `--quiet`** also sets show AI command output to **false** (cross-reference R3).  
3. Keep existing behavior: `--log-level` overrides `-q`/`-v` for log level only; `-v` still enables AI streaming unless `--no-ai-cmd-output` is set.

**Acceptance:**  
- No remaining “stderr” for Ralph operational log output in R5.  
- -q is documented as affecting both log level and show AI output (with ref to R3).

---

### T4 — O4 R2 and R6: Statistics and progress to stdout  
**Priority:** 2 · **Depends on:** T3  

**Files:**  
- `docs/intent/O4-observability/R2-iteration-statistics.md`  
- `docs/intent/O4-observability/R6-iteration-progress.md`

**What to do:**  
In both files, replace every reference to statistics (R2) or progress messages (R6) going to **stderr** with **stdout**. Align wording with R5 (operational output to stdout). Do not change behavior of *when* stats or progress are emitted, only *where* (stdout).

**Acceptance:**  
- R2: “Statistics … to stderr” → “Statistics … to stdout”; same for examples and acceptance criteria.  
- R6: “Progress messages go to stderr” → “Progress messages go to stdout”; same for examples and acceptance criteria.

---

### T5 — O2 R9: CLI table — add --no-ai-cmd-output, update -q and -v  
**Priority:** 2 · **Depends on:** T2, T3  

**File:** `docs/intent/O2-configurable-behavior/R9-cli-interface-reference.md`

**What to do:**  
In the **Output control** table for `ralph run`:  
1. **Add row:** `--no-ai-cmd-output` | (no short) | — | Set show AI command output to false. | O4/R3  
2. **Update -q row:** Description → “Quiet: set log level to error and do not stream AI command output.” (Specified in O4/R5, O4/R3.)  
3. **Update -v row:** Description → “Verbose: set log level to debug and stream AI command output (unless --no-ai-cmd-output).” (Specified in O4/R3, O4/R5.)

**Acceptance:**  
- Table has four output-control entries: `--verbose`, `--quiet`, `--log-level`, `--no-ai-cmd-output`.  
- Descriptions match TASK.md and O4/R3, R5.

---

### T6 — O2 R8: Env default true for show_ai_output  
**Priority:** 2 · **Depends on:** T2  

**File:** `docs/intent/O2-configurable-behavior/R8-environment-variable-reference.md`

**What to do:**  
1. In the table row for **RALPH_LOOP_SHOW_AI_OUTPUT**, add a note: when **unset**, default is **true** (per O4/R3).  
2. In Edge cases, add: “Variable unset → default true; AI output is streamed.”

**Acceptance:**  
- Default true is documented for the env var.  
- Edge case “Unset → default true” is present.

---

### T7 — Implementation: Output routing (stdout = run log, stderr = fatal only)  
**Priority:** 3 · **Depends on:** T1, T2, T3, T4  

**Files:** Implementation (e.g. `internal/logger/`, `internal/runner/loop.go`, `cmd/ralph/` — wherever run is wired and output is written).

**What to do:**  
1. Route **all** Ralph operational output (iteration progress per R6, iteration statistics per R2, log-leveled messages per R5) to **stdout**.  
2. Use **stderr** only for output that is not part of the normal run log (e.g. fatal startup errors: config load failure, invalid args, missing binary, usage errors).  
3. When “show AI command output” is true, mirror the AI CLI’s stdout/stderr to **stdout** (same stream as Ralph’s operational output).  
4. Ensure `ralph run <alias> > run.log` captures the full run log (Ralph messages + AI stream in order).

**Acceptance:**  
- Default `ralph run` produces a single logical run log on stdout; stderr is unused on the happy path.  
- Fatal/startup errors go to stderr only.  
- Redirect test: `ralph run ... > run.log` yields a complete, ordered run log in `run.log`.

---

### T8 — Implementation: Default show_ai_output true and flags (--no-ai-cmd-output, -q, -v)  
**Priority:** 3 · **Depends on:** T2, T3, T5  

**Files:** Implementation (e.g. `internal/config/` defaults and overlay, `cmd/ralph/` flag parsing and wiring).

**What to do:**  
1. **Default:** Resolved `show_ai_output` when nothing is set = **true**.  
2. **CLI flag:** Add **`--no-ai-cmd-output`**; when set, resolved show_ai_output = false.  
3. **Precedence (show AI output):** `--no-ai-cmd-output` → false; `-q` → false; `-v` → true; else config/env; else **true**.  
4. **-q:** Set log level to error **and** show_ai_output to false.  
5. **-v:** Set log level to debug **and** show_ai_output to true **unless** `--no-ai-cmd-output` is set.  
6. **Log level precedence:** `--log-level` (if present) else `-q` else `-v` else config/env else default info.

**Acceptance:**  
- Default run streams AI output.  
- `-q` suppresses both log verbosity and AI output.  
- `-v --no-ai-cmd-output` gives verbose logs but no AI stream.  
- `--log-level` only affects log level, not show_ai_output.  
- Behavior matches TASK.md and updated O4/R3, R5 and O2/R9.

---

## Beads (work tracking)

**Publish tasks:** Run from repo root:

```bash
./scripts/publish-plan-beads.sh
```

This creates beads for T1–T8 with titles and descriptions from this plan. After running, **add dependencies** so `bd ready` reflects the task graph. Run `bd list` to get issue IDs (they look like `ralph-xxx`), then for each dependent task run:

```bash
bd update <issue-id> --deps blocks:<dep-id>
```

**Dependency map (issue-id depends on dep-id):**

| Task | Issue (this id) | Depends on (blocks:) |
|------|-----------------|----------------------|
| T2   | &lt;T2-id&gt;   | &lt;T1-id&gt;        |
| T3   | &lt;T3-id&gt;   | &lt;T1-id&gt;        |
| T4   | &lt;T4-id&gt;   | &lt;T3-id&gt;        |
| T5   | &lt;T5-id&gt;   | &lt;T2-id&gt;, &lt;T3-id&gt; |
| T6   | &lt;T6-id&gt;   | &lt;T2-id&gt;        |
| T7   | &lt;T7-id&gt;   | &lt;T1-id&gt;, &lt;T2-id&gt;, &lt;T3-id&gt;, &lt;T4-id&gt; |
| T8   | &lt;T8-id&gt;   | &lt;T2-id&gt;, &lt;T3-id&gt;, &lt;T5-id&gt; |

If `bd update` accepts only one `--deps` at a time, run it once per dependency (e.g. for T7 run four times with blocks:T1-id, blocks:T2-id, blocks:T3-id, blocks:T4-id).

---

## Verification (after all tasks)

- **Intent:** Re-read O4 README, R3, R5, R2, R6 and O2 R8, R9; confirm no remaining “stderr” for run log, default true for show AI output, and that --no-ai-cmd-output and -q/-v behavior and precedence match TASK.md.  
- **Implementation:** Run default (AI output streamed, log level info); run -q (no AI output, only errors); run -v --no-ai-cmd-output (verbose logs, no AI stream); run `ralph run ... > run.log` and confirm full log on stdout.
