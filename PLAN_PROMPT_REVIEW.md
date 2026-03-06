# Plan: O5 Prompt Review Implementation

This plan implements **O5 — Prompt review** so that prompts can be reviewed for quality and structure before or without running the loop. Outcome and requirements are specified in `docs/intent/O5-prompt-review/`.

**Verification (from O5 README):** User runs `ralph review` with prompt from alias, `-f <path>`, or stdin; receives a report (including suggested revision); with `--apply`, can write revision to the prompt file (confirm or `-y`); exit 0 (ok), 1 (errors in prompt), 2 (review failed or apply invalid).

**Existing implementation:** O1–O4 are implemented. Reuse: `internal/config` (LoadConfigWithProvenanceAndExplicit, Validate, ResolveAICommand, prompt aliases), `internal/prompt` (LoadPrompt, ResolveMode — note R1 resolution order differs: file wins over alias over stdin), `internal/cmdparse` (Parse), `cmd/ralph/main.go` (cobra structure). Review does **not** use the runner loop; it runs one (or two, for apply) AI invocation(s).

---

## Task dependency overview

```
T1 (review subcommand + R1 input modes)
  → T2 (R3 review output path)
  → T3 (R2 review prompt composition)
  → T4 (R8 failure handling)  ← can be done in parallel with T3 once T1/T2 exist
  → T5 (R9 report file verification)
  → T6 (R6 report format and exit codes)
  → T7 (R4 prompt output path)
  → T8 (R5 apply and revision phase)
  → T9 (R7 configurable review stdout)
```

---

## T1: Review subcommand and R1 input modes

**Priority:** 1 (must be first)  
**Dependencies:** None  
**Spec:** [R1 — Review command input modes](docs/intent/O5-prompt-review/R1-review-command-input-modes.md)

**Objective:** Add `ralph review` subcommand that accepts the prompt from exactly one of: alias, file (`-f`/`--file`), or stdin. Resolution order is **mutually exclusive**: (1) if `-f` is present, use file (ignore alias and stdin); (2) else if positional alias is present, use alias; (3) else use stdin. Load prompt once at start; do not re-read during review.

**Context:**

- **CLI surface:** `ralph review [alias]`, `ralph review -f <path>`, `ralph review` (stdin). Add flags: `--review-output`, `--apply`, `--prompt-output` (can be stubbed or validated later).
- **Resolution order differs from current `prompt.ResolveMode`:** R1 says when both alias and `-f` are present, **file wins** (prompt from file; for apply, destination is the `-f` path). Current `ResolveMode` returns an error for "ambiguous" alias+file. Either add a review-specific resolver (e.g. `ResolveModeReview(alias, filePath string) (Mode, error)`) that applies R1 order, or change `ResolveMode` and ensure `run` still gets the desired behavior (R1 explicitly says "alias is ignored" when `-f` is present).
- **Validation:** Alias not in config, file missing/unreadable, or file/alias source empty → fail fast exit 2 (O2 R4 / R8). Stdin with no data (empty pipe) → exit 2. Use existing config load and prompt load where possible; exit code 2 for these cases (T4 will centralize messages; here just wire failure path).
- **Output:** For T1, after loading prompt, you may stub the rest (e.g. print "review not yet implemented" and exit 0) so the command and input resolution are testable.

**Acceptance:**

- [ ] `ralph review <alias>` loads prompt for that alias and proceeds (or stubs).
- [ ] `ralph review -f <path>` reads prompt from file; no alias lookup for content.
- [ ] `ralph review` with piped stdin reads prompt from stdin.
- [ ] When both alias and `-f` are given, prompt is read from file; alias ignored.
- [ ] When stdin is connected and `-f` is given, stdin is ignored; prompt from file.
- [ ] Prompt is loaded once and not re-read during the same run.
- [ ] Missing alias, missing/unreadable/empty file, or empty stdin → exit 2 (message can be refined in T4).

---

## T2: R3 — Review output path

**Priority:** 2  
**Dependencies:** T1  
**Spec:** [R3 — Review output path](docs/intent/O5-prompt-review/R3-review-output-path.md)

**Objective:** The report is **always** written to a file. Report path is either `--review-output <path>` (resolved relative to CWD) or a new unique file in the system temp directory (e.g. `os.TempDir()` + unique filename). When temp is used, communicate the path to the user (e.g. stderr or a summary line) so they can find the report. Validate before spawning AI: parent directory writable for `--review-output`; temp dir available when not set. Invalid/unwritable path → exit 2.

**Context:**

- **Explicit path:** `--review-output ./report.txt` → report path is that file; allow overwrite (spec recommends).
- **Directory vs file:** If `--review-output` points to an existing directory, treat as invalid (exit 2); do not auto-choose a filename inside it.
- **Temp path:** Use platform temp dir; create a unique file path (e.g. `ralph-review-*`). After choosing, print something like "Report written to <path>" (or equivalent) so user can discover it.
- **Validation:** Before running the review-phase AI: if `--review-output` is set, ensure parent exists and is writable; if not set, ensure temp dir is available. On failure, exit 2 with a clear message.

**Acceptance:**

- [ ] With `--review-output <path>`, report path is that path (used when AI is instructed and for R9 verification).
- [ ] Without `--review-output`, report path is a unique file in system temp directory.
- [ ] When temp is used, user can discover the path (e.g. printed to stderr or stdout).
- [ ] Invalid or unwritable report path (or temp unavailable) → exit 2 before spawning AI.
- [ ] Report path is the single canonical location for the report file (R9 will verify it).

---

## T3: R2 — Review prompt composition

**Priority:** 3  
**Dependencies:** T1, T2  
**Spec:** [R2 — Review prompt composition](docs/intent/O5-prompt-review/R2-review-prompt-composition.md)

**Objective:** Compose the prompt sent to the AI for the **review phase** by: (1) Ralph's review instructions (embedded in binary via `//go:embed`), (2) a path directive telling the AI to write the report to the review output path (from T2), (3) the user's prompt content (from T1). Ralph does **not** parse AI stdout to derive the report; the report is whatever the AI writes to the file at that path.

**Context:**

- **Embedded instructions:** Add an `internal/review` (or similar) package. Use Go `embed` to embed one or more instruction files (e.g. `review_instructions.md`, `revision_instructions.md` for later apply). These are part of the binary, not read from the user's repo.
- **Composition order:** [Ralph instructions] + [path directive: "Write your report to: <path>"] + [user prompt]. Path is the exact string from T2 (review output path). Escape or quote if needed for the prompt (e.g. path with spaces).
- **Revision phase (later):** R5 will add a second invocation with embedded revision instructions and interpolated review output path + prompt output path. For T3, only the review phase is needed.
- **AI invocation:** Run one AI process (same mechanism as runner spawn: resolve AI command from config, parse with `cmdparse.Parse`, spawn with stdin set to the composed prompt). Do not use the loop; no preamble injection; no signal scanning. Capture stdout/stderr for R7 (default: show to user); report content is read from the **file** after the AI exits (R9).

**Acceptance:**

- [ ] Review-phase prompt includes embedded review instructions + path directive + user prompt.
- [ ] Path in the directive is the review output path from T2.
- [ ] Ralph does not parse AI stdout to build the report; report is defined as the file at the review path.
- [ ] Review (and later revision) instructions are embedded in the binary via `embed`, not read from user repo.
- [ ] Single AI invocation for review phase; AI is instructed to write report to the given path.

---

## T4: R8 — Review failure handling

**Priority:** 4  
**Dependencies:** T1, T2 (and effectively T3 so there is a run path)  
**Spec:** [R8 — Review failure handling](docs/intent/O5-prompt-review/R8-review-failure-handling.md)

**Objective:** Ensure all conditions that must yield **exit 2** are implemented with clear error messages. Exit 2 means "review did not complete successfully" or "apply invalid"; never exit 0 or 1 when these conditions hold.

**Context — conditions that must yield exit 2:**

1. **Config invalid or required missing** — e.g. malformed YAML, missing required key. Message to stderr; exit 2 before loading prompt or spawning AI.
2. **Prompt source invalid (alias or file)** — Unknown alias ("unknown prompt alias \"<alias>\""), file not found, file unreadable, file empty (0 bytes). Fail fast; exit 2 before any AI invocation.
3. **Stdin empty** — Input mode is stdin but no data read. Message e.g. "no prompt content from stdin"; exit 2.
4. **AI spawn/crash** — AI binary not found, spawn fails, or process exits in a way that prevents review from completing. Exit 2 (distinguish from "AI ran but reported errors in prompt" → exit 1).
5. **Invalid apply** — Stdin + `--apply` without `--prompt-output`. Message e.g. "stdin input with --apply requires --prompt-output <path>"; exit 2 before revision phase (validation before review phase is acceptable).
6. **Review output path invalid/unwritable** — Per R3: unwritable parent, path is directory, temp unavailable. Exit 2 before spawning AI.
7. **Report file missing after review** — Per R9: file not at expected path after AI exits. Exit 2 (T5).
8. **Revision file missing after apply** — If we verify revision was written when apply was requested and file is missing, exit 2 (T8).

**Acceptance:**

- [ ] Each of the above conditions produces a clear message and exit 2.
- [ ] Exit 2 is not used when review completed and report exists; only 0 (no errors in prompt) or 1 (errors in prompt) in that case.
- [ ] Error messages allow user/script to correct the condition (e.g. missing `--prompt-output` when using stdin + apply).

---

## T5: R9 — Report file verification

**Priority:** 5  
**Dependencies:** T3, T4  
**Spec:** [R9 — Report file verification](docs/intent/O5-prompt-review/R9-report-file-verification.md)

**Objective:** After the review-phase AI process exits, verify that a **file** exists at the review output path. If not (AI didn't write it, wrote elsewhere, or crashed), do not exit 0 or 1; report failure and exit 2. Do not run apply/revision when verification fails.

**Context:**

- **When:** Immediately after the review-phase AI exits. Before deriving exit 0/1 from report content (T6), before printing report to stdout (T7), and before any apply flow (T8).
- **Check:** File must exist at the path (regular file). Spec: existence only; empty file can be treated as pass (recommended); R6 parsing may then yield exit 1 if content is unparseable. If path is a directory or symlink to missing file → exit 2.
- **Order:** (1) Review AI runs and exits. (2) Verify file exists. (3) If no → exit 2, no apply. (4) If yes → read report, derive 0/1 (T6), optionally print (T7), if apply then revision phase (T8).

**Acceptance:**

- [ ] After review-phase AI exits, Ralph checks for a file at the review output path.
- [ ] If file does not exist → exit 2; no apply; clear message (e.g. "report file not found at <path>").
- [ ] If file exists → proceed to parse for exit code and optionally apply.
- [ ] Verification runs before apply/revision phase.

---

## T6: R6 — Report format and exit code derivation

**Priority:** 6  
**Dependencies:** T5  
**Spec:** [R6 — Report format and exit codes](docs/intent/O5-prompt-review/R6-report-format-exit-codes.md)

**Objective:** The report contains (1) narrative feedback, (2) a **machine-parseable summary** so scripts/CI can gate on result, and (3) the full suggested revision. Ralph reads the report file (after R9 passes) and parses the machine-parseable section to set **exit 0** (no errors / only warnings if policy allows) or **exit 1** (one or more errors in prompt). Exit 2 is already handled by R8/R9.

**Context:**

- **Report content:** Produced by the AI per embedded instructions (T3). Ralph must define and document one canonical format for the machine-parseable block (e.g. a line like `ralph-review: status=ok|errors|warnings`, `errors=N`, `warnings=N`, or a small YAML/JSON block under a known heading). Embed this format in the review instructions so the AI emits it; implement parser in Ralph.
- **Exit derivation:** After R9: read report file; extract machine-parseable block; if `status=ok` and `errors=0` (and optional warnings policy) → exit 0; if `status=errors` or `errors>=1` → exit 1. If block missing or malformed, spec recommends exit 1 (fail-safe for CI).
- **Document:** Document the format (e.g. in AGENTS.md or `docs/`) so CI/scripts can parse it and so the embedded prompt can require it.

**Acceptance:**

- [ ] Report is expected to contain narrative, machine-parseable summary, and full suggested revision (AI produces per R2 instructions).
- [ ] Ralph parses the report file for the machine-parseable summary and sets exit 0 or 1 accordingly.
- [ ] Exit 0: review completed, no errors (or only warnings if implemented).
- [ ] Exit 1: review completed, one or more errors in prompt.
- [ ] Missing or unparseable summary → exit 1 (recommended) or document alternative.
- [ ] Format is documented for automation and for the embedded review instructions.

---

## T7: R4 — Prompt output path

**Priority:** 7  
**Dependencies:** T1 (and T8 will consume this)  
**Spec:** [R4 — Prompt output path](docs/intent/O5-prompt-review/R4-prompt-output-path.md)

**Objective:** Support `--prompt-output <path>` to direct where the **suggested revised prompt** is written. When apply is not requested, `--prompt-output` is where the revision is written without modifying the source. When apply **is** requested: for alias or file, apply destination defaults to the source path (alias path or `-f` path), and `--prompt-output` can override to a different file; for **stdin**, `--prompt-output` is **required** with `--apply` (otherwise invalid → exit 2). The resolved "prompt output path" is interpolated into the revision-phase prompt (T8).

**Context:**

- **Without apply:** `--prompt-output out.md` → write suggested revision to `out.md`; source unchanged.
- **With apply + alias/file:** Default apply destination = source file; optional `--prompt-output` can direct revision to a different path (apply-to-different-file).
- **With apply + stdin:** `--prompt-output` is required; absence is invalid (exit 2, already in T4). Resolved prompt output path = `--prompt-output`.
- **Validation:** Stdin + `--apply` without `--prompt-output` → exit 2 (T4). When revision phase would write, validate path is writable; if not, exit 2 (R8).
- **Interpolation:** The resolved path (for apply or for non-apply `--prompt-output`) is what gets interpolated into the revision-phase prompt in T8.

**Acceptance:**

- [ ] `--prompt-output <path>` writes suggested revision to that path when not applying; source unchanged.
- [ ] With apply + alias or file, revision can go to source (default) or to `--prompt-output` if set.
- [ ] With apply + stdin, `--prompt-output` is required; otherwise exit 2.
- [ ] Resolved prompt output path is available for revision-phase interpolation (T8).

---

## T8: R5 — Apply with confirmation and revision phase

**Priority:** 8  
**Dependencies:** T3, T5, T6, T7  
**Spec:** [R5 — Apply with confirmation and revision phase](docs/intent/O5-prompt-review/R5-apply-confirmation-revision-phase.md)

**Objective:** When user passes `--apply`, optionally run a **revision phase** (second AI invocation) that reads the report from the review output path and writes the revised prompt to the prompt output path (R4). Both paths are interpolated into the revision-phase prompt. With `--apply` and no `-y`, prompt the user for confirmation (e.g. "Apply revision to <path>? [y/N]"); with `-y`, apply without prompting. If no TTY and no `-y`, do not hang (e.g. treat as decline or exit 2 requiring `-y` for non-interactive).

**Context:**

- **Revision-phase prompt:** Use a second embedded file (revision instructions) via `embed`. Interpolate: (1) review output path — "Read the report from: <path>", (2) prompt output path — "Write the revised prompt to: <path>". AI reads report from file, writes revised prompt to file.
- **When to run:** Only when apply is requested and after review phase completed and R9 passed (report file exists). Stdin+apply without `--prompt-output` already exited 2 in T4/T7.
- **Confirmation:** If `--apply` and not `-y`: prompt "Apply revision to <path>? [y/N]". On y/yes → run revision phase (if used) and write. On decline or EOF → do not write; exit 0 or 1 from review result only. If no TTY and no `-y`, either treat as decline or exit 2 with "use -y for non-interactive apply" (spec allows either).
- **Writing:** After revision phase, read the file the AI was instructed to write to (prompt output path); if implementing verification, ensure file exists and optionally overwrite the source path with that content (or the revision phase instructs AI to write to the final path; then verification is "file exists at prompt output path").
- **Revision file missing:** If AI does not write to prompt output path, exit 2 per R8 (recommended: verify revision file exists after revision phase when apply was requested).

**Acceptance:**

- [ ] `--apply` triggers (after successful review) confirmation or direct apply when `-y`.
- [ ] Revision-phase prompt includes review output path and prompt output path; AI is instructed to read report and write revised prompt to that path.
- [ ] Without `-y`, user is prompted; on confirm, revision is written; on decline, exit 0/1 without writing.
- [ ] With `-y`, revision is written without prompting.
- [ ] Stdin + apply without `--prompt-output` already exits 2 (T7/T4).
- [ ] Apply works for alias, file, and stdin (stdin requires `--prompt-output`).
- [ ] If revision phase does not produce a file at prompt output path, exit 2 (recommended).

---

## T9: R7 — Configurable review output to stdout

**Priority:** 9  
**Dependencies:** T3, T5, T6, T7  
**Spec:** [R7 — Configurable review output to stdout](docs/intent/O5-prompt-review/R7-configurable-review-stdout.md)

**Objective:** By default, both (1) the AI command output (stdout/stderr from the AI process) and (2) the report content (read from the report file after AI exits) are exposed to the user (e.g. stdout). User can configure or use flags so that AI output and/or report are **not** printed to stdout (e.g. `--quiet` for no AI output, `--report-to-file-only` so report is only in the file). Report must always be written to the file (R3); this requirement only controls what goes to stdout.

**Context:**

- **Default:** Show AI output and report content (e.g. after AI exits, read report file and print to stdout). Report also saved to file.
- **Flags/config:** At least one of: CLI flag (e.g. `--quiet`, `--report-to-file-only`) or config (e.g. `review.show_ai_output`, `review.print_report`). User must be able to achieve "report only to file, nothing to stdout" for scripting.
- **Invariant:** Suppressing stdout never prevents writing the report to the file at the review output path.
- **Same behavior for all input modes:** Alias, file, stdin — only the prompt source differs; stdout behavior is the same.

**Acceptance:**

- [ ] By default, AI command output and report content are exposed to stdout; report also in file.
- [ ] User can suppress AI output and/or report printing (e.g. report-only-to-file).
- [ ] Report is always written to the file (R3); stdout config only affects what is printed.
- [ ] Behavior consistent for alias, file, and stdin input.

---

## Summary table

| Task | Requirement | Brief description |
|------|-------------|--------------------|
| T1 | R1 | Review subcommand; input modes (alias / file / stdin) with R1 resolution order |
| T2 | R3 | Review output path: `--review-output` or temp; communicate path; validate |
| T3 | R2 | Review prompt composition: embed instructions, path directive, user prompt; one AI run |
| T4 | R8 | All exit 2 conditions and error messages |
| T5 | R9 | Report file verification after AI exits |
| T6 | R6 | Report format; machine-parseable summary; exit 0/1 derivation |
| T7 | R4 | `--prompt-output`; apply destination resolution; stdin+apply requires it |
| T8 | R5 | `--apply`, `-y`, confirmation, revision phase with path interpolation |
| T9 | R7 | Configurable stdout (default: show AI + report; flags to suppress) |

---

## Beads (bd) issue mapping

| Plan task | Bead ID    | Dependencies (blocked by) |
|-----------|------------|----------------------------|
| T1        | ralph-e8a  | —                          |
| T2        | ralph-wrd  | ralph-e8a                  |
| T3        | ralph-orr  | ralph-e8a                  |
| T4        | ralph-5vr  | ralph-e8a                  |
| T5        | ralph-5b8  | ralph-orr, ralph-5vr       |
| T6        | ralph-xkg  | ralph-5vr *(see note)*     |
| T7        | ralph-bgo  | ralph-e8a                  |
| T8        | ralph-1a3  | ralph-orr, ralph-5b8, ralph-xkg, ralph-bgo |
| T9        | ralph-b1e  | ralph-orr, ralph-5b8, ralph-xkg, ralph-bgo |

**Note:** T6 (ralph-xkg) is stored in beads with a single dependency on T4 (ralph-5vr). Per the plan, T6 should be done after T5 (report file verification); complete T5 (ralph-5b8) before T6 when implementing.

**Check ready work:** `bd ready` or `bd ready --json`. **Claim:** `bd update <id> --claim`. **Close:** `bd close <id> --reason "Completed"`.
