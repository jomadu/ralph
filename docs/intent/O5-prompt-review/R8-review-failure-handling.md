# R8: Review Failure Handling

**Outcome:** O5 — Prompt review

## Requirement

The system treats certain conditions as review failures and exits with code 2. These include: invalid or missing configuration, missing or unreadable prompt source (for alias or file input), failure to spawn or run the AI process, invalid apply request (e.g. `--apply` with stdin and no `--prompt-output`), and unwritable or invalid report path. The user and scripts can rely on exit 2 to mean "review did not complete successfully" or "apply was invalid," distinct from exit 0 (success) and exit 1 (review completed but prompt has errors).

## Specification

**Conditions that must yield exit 2 (review failure or invalid apply):**

1. **Configuration invalid or required config missing** — e.g. config file is malformed (YAML parse error), or a required key for review is missing. Ralph must report the failure (message to stderr or stdout) and exit 2. Do not proceed to load prompt or spawn AI.
2. **Prompt source invalid (alias or file input)** — Alias not found in config; file path missing; file unreadable (permission); file empty (0 bytes). Per O2 R4 semantics: fail fast with a clear message (e.g. "unknown prompt alias \"x\"", "prompt file not found: path", "prompt file not readable: path", "prompt file is empty: path"). Exit 2 before any AI invocation.
3. **Stdin empty** — When input mode is stdin and no data is read (empty pipe or no stdin). Exit 2 with message (e.g. "no prompt content from stdin").
4. **AI command cannot be spawned or fails before review completes** — AI binary not found, spawn fails, or process exits in a way that prevents the review from completing (e.g. crash before writing report). Exit 2. Distinguish from "AI ran but reported errors in the prompt" (exit 1).
5. **Invalid apply request** — Stdin + `--apply` without `--prompt-output`. Message must tell the user what is wrong (e.g. "stdin input with --apply requires --prompt-output <path>"). Exit 2 before revision phase.
6. **Review output path invalid or unwritable** — Per R3: `--review-output` path has unwritable parent, or path is a directory; temp dir unavailable. Exit 2 before spawning AI.
7. **Report file missing after review phase** — Per R9: file not found at the expected path after AI exits. Exit 2.
8. **Revision phase apply: revision file missing** — If implementation verifies that the revision was written when apply was requested, and the file is missing after revision phase, exit 2 (consistent with R5 edge case).

**Error messages:** Each failure must produce a message that allows the user or script to correct the condition. Examples: missing alias → "unknown prompt alias \"<alias>\""; stdin + apply no path → "stdin input with --apply requires --prompt-output <path>"; unwritable path → "cannot write report to <path>: <reason>".

**Consistency:** Exit 2 is used only for "review did not complete" or "apply invalid." Exit 0 and 1 are used only when the review completed and the report file exists (R9 passed); then 0 = no errors in prompt, 1 = errors in prompt. Scripts can rely on exit 2 to mean "do not trust the report or apply."

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Config file missing (optional config) | If config is optional, proceed with defaults; if required for review, exit 2. Document. |
| Multiple failures (e.g. invalid config and invalid path) | Report first encountered failure; exit 2; do not cascade. |
| AI exits non-zero but wrote the report | If R9 passes (file exists), derive exit 0 or 1 from report content; do not use exit 2 for AI exit code alone. |

### Examples

#### Alias not found

**Input:** `ralph review unknown-alias`

**Expected output:** Message like "unknown prompt alias \"unknown-alias\""; exit 2.

#### Stdin + apply without --prompt-output

**Input:** `cat p.md | ralph review --apply`

**Expected output:** Message like "stdin input with --apply requires --prompt-output <path>"; exit 2.

## Acceptance criteria

- [ ] When configuration is invalid or required config is missing, Ralph reports the failure and exits 2.
- [ ] When prompt source is invalid (alias not found, file missing or unreadable), Ralph fails fast and exits 2.
- [ ] When the AI command cannot be spawned or fails in a way that prevents the review from completing, Ralph exits 2.
- [ ] When the user requests `--apply` with stdin but does not provide `--prompt-output`, Ralph reports an error and exits 2.
- [ ] When the review output path (from `--review-output` or temp) is invalid or unwritable, Ralph exits 2 (or fails before running the AI).
- [ ] Exit 2 is used consistently for these failure cases so scripts can distinguish them from exit 0 (ok) and exit 1 (errors in prompt).
- [ ] Error messages or logs give the user enough information to correct the condition (e.g. missing `--prompt-output` when using stdin + apply).

## Dependencies

- R6 (exit codes) defines exit 2 semantics; this requirement enumerates the conditions that must yield exit 2. R3 (review output path) and R9 (report file verification) define path-related failures that also result in exit 2. R5 defines the stdin + apply + `--prompt-output` validation.
