# R1: Review Command Input Modes

**Outcome:** O5 — Prompt review

## Requirement

The system accepts the prompt to be reviewed from three sources: a configured prompt alias, a file path (e.g. `-f <path>`), or standard input (e.g. piped content). The `ralph review` command resolves the prompt from exactly one of these sources per invocation and uses it as the input to the review workflow. Ralph loads the prompt once and does not re-read it during the review.

## Specification

**Resolution order (mutually exclusive, one source per invocation):**

1. If `-f` / `--file <path>` is present, the prompt is read from that file path. Alias positional argument and stdin are ignored for prompt content.
2. Else if a positional argument is present (e.g. `ralph review build`), it is treated as a prompt alias; the prompt is loaded from the path configured for that alias (same resolution as `ralph run <alias>` per O2).
3. Else the prompt is read from stdin (e.g. `cat prompt.md | ralph review` or `ralph review` with stdin connected).

Ralph loads the prompt content once at the start of the review workflow (before composing the review prompt or spawning the AI). The same in-memory content is used for the entire review run; the source is not re-read mid-run. For alias and file, validation and fail-fast behavior follow O2 R4 (missing/unreadable/empty → exit 2 before any AI invocation).

**CLI surface:** `ralph review [alias]` and `ralph review -f <path>` and `ralph review` (stdin). Flags such as `--review-output`, `--apply`, `--prompt-output` are independent of input mode.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User passes both alias and `-f <path>` | File wins; prompt is read from `-f` path; for apply, destination is the `-f` path (alias is ignored). |
| User runs `ralph review -f ./p.md` and stdin is also connected | Stdin is ignored; prompt is read from `./p.md`. |
| User runs `ralph review` with no args and no stdin | Stdin is empty; treat as invalid prompt source; fail fast with exit 2 (per R8). |
| Alias not in config | Fail fast, exit 2 (O2 R4 / R8). |
| File path missing or unreadable | Fail fast, exit 2 (O2 R4 / R8). |
| File or alias source is empty (0 bytes) | Fail fast, exit 2 (O2 R4 / R8). |

### Examples

#### Alias input

**Input:** `ralph review build` with alias `build` configured to `./prompts/build.md`.

**Expected output:** Ralph loads content from `./prompts/build.md`, runs reviewer with that content; report and exit code per R3–R9.

**Verification:** Report content reflects the text of `./prompts/build.md`; no read from stdin.

#### File input

**Input:** `ralph review -f ./my/prompt.md`.

**Expected output:** Ralph loads content from `./my/prompt.md`, runs reviewer; no alias lookup.

**Verification:** Report reflects `./my/prompt.md` content; path need not be under config.

#### Stdin input

**Input:** `cat prompt.md | ralph review`.

**Expected output:** Ralph reads prompt from stdin, runs reviewer. No alias or file path used for content.

**Verification:** Report reflects piped content; with `--apply`, `--prompt-output` is required (R4/R5).

## Acceptance criteria

- [ ] User can run `ralph review <alias>`; Ralph loads the prompt associated with that alias and runs the reviewer with that content.
- [ ] User can run `ralph review -f <path>`; Ralph reads the prompt from the file at that path and runs the reviewer.
- [ ] User can pipe prompt content into `ralph review` (e.g. `cat prompt.md | ralph review`); Ralph reads from stdin and runs the reviewer with that content.
- [ ] Exactly one input source is used per invocation; alias, file path, and stdin are mutually exclusive in resolution (precedence or flag semantics are specified elsewhere).
- [ ] The prompt is loaded once at the start of the review and not re-read from the source during the same run.

## Dependencies

- Configuration and alias resolution (O2). Prompt source validation and fail-fast behavior (O2 R4) apply when alias or file is used.
