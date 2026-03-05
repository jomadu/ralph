# R4: Dry-Run Mode

**Outcome:** O4 — Observability

## Requirement

The system prints the fully assembled prompt — preamble plus prompt content — to stdout without spawning an AI process, enabling the user to validate prompt assembly and configuration before committing to a run.

## Specification

Dry-run mode is triggered by `--dry-run` or `-d`. When triggered, Ralph performs configuration resolution and validation (O2), resolves the prompt source (alias, `-f` file, or stdin), loads the prompt content once, and assembles the prompt exactly as for the first iteration of a normal run (O1/R8: preamble plus prompt content). Ralph then prints the complete assembled prompt to stdout and exits 0. No AI CLI process is spawned; no loop runs; no signal scanning or iteration logic runs.

**Order of operations:**

1. Resolve config layers (O2): CLI flags, environment, workspace config, global config, defaults. Resolve prompt alias or `-f` path or stdin.
2. Validate configuration (O2/R3). If validation fails (invalid values, unknown keys that are errors, etc.), emit errors to stderr and exit with a non-zero code; do not print any prompt to stdout.
3. Resolve prompt source: if alias, load the file at the path defined for that alias; if `-f`, load that file; if stdin, read from stdin. If the prompt source is missing, unreadable, or empty (per O2/R4), fail-fast: error to stderr, non-zero exit; do not print prompt.
4. Build assembled prompt for "iteration 1": Generate preamble per O1/R8 (iteration 1, limit from config, optional context from `--context` flags). If preamble is disabled for this prompt, use only prompt content. Concatenate: `<preamble>\n\n<prompt_content>` or just `<prompt_content>`.
5. Print the assembled prompt to stdout, unchanged (no extra headers or footers). This is exactly what would be piped to the AI CLI's stdin on the first iteration of a normal run.
6. Exit 0.

**Output:** Only the assembled prompt goes to stdout. Errors and warnings go to stderr (consistent with R5). So: `ralph run build -d > assembled.txt` captures only the prompt; config errors appear on stderr and are not in `assembled.txt`.

**Invariants:**

- Dry-run does not spawn the AI CLI. No subprocess.
- Dry-run does not run the loop. No iterations, no signals, no exit codes 1 or 2 from loop termination (R1). Exit is 0 on success.
- The printed content is identical to what would be sent to the AI on iteration 1: same preamble format (O1/R8), same prompt content, same ordering.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| `ralph run build -d` with valid config and prompt | Config and prompt resolved; assembled prompt printed to stdout; exit 0; no AI process |
| `ralph run build --dry-run` | Same as `-d` |
| Invalid config (e.g., negative max iterations per O2/R3) | Validation fails; error on stderr; no prompt on stdout; non-zero exit |
| Missing prompt file (alias points to missing path) | Fail-fast (O2/R4); error on stderr; no prompt on stdout; non-zero exit |
| `-d` with `-f ./missing.md` | Same as missing prompt: fail-fast; no prompt printed |
| Dry-run with `--context "foo"` | Preamble includes CONTEXT section per O1/R8; printed prompt includes it |
| Dry-run with preamble disabled (`loop.preamble: false` or per-prompt) | Printed output is prompt content only, no preamble |
| Dry-run with `-n 20` | Preamble shows iteration 1 of 20 (config resolved; first iteration) |
| `ralph run build -d 2>/dev/null` | Only assembled prompt on stdout; errors hidden |
| Dry-run does not consume stdin for AI: prompt from alias | Stdin is not read (prompt from file); assembled prompt printed |

### Examples

#### Valid dry-run, preamble enabled

**Input:**
`ralph run build -d`. Alias `build` points to `./prompts/build.md`. Config has `default_max_iterations: 10`, preamble enabled. File `./prompts/build.md` contains `Fix the tests.\n`.

**Expected output:**
stdout contains exactly what would be piped to the AI on iteration 1, e.g.:

```
[RALPH] Iteration 1 of 10

Fix the tests.
```

(With two newlines between preamble and content per O1/R8.) Exit 0. No AI process started.

**Verification:**
- Exit code 0
- No child process for AI CLI
- First line of stdout is `[RALPH] Iteration 1 of 10` (or equivalent per O1/R8)

#### Invalid config — no prompt printed

**Input:**
`ralph run build -d` with `ralph-config.yml` containing `default_max_iterations: -1` (invalid per O2/R3).

**Expected output:**
Configuration validation fails. Error message(s) on stderr. Nothing or only partial/corrupt content on stdout (implementation may print nothing). Exit code non-zero.

**Verification:**
- Exit code ≠ 0
- stderr contains an error related to config validation
- No valid assembled prompt (with preamble) on stdout

#### Dry-run with context

**Input:**
`ralph run build -d --context "Focus on unit tests"`. Preamble enabled. Prompt content: `Refactor the module.\n`.

**Expected output:**
stdout contains preamble with CONTEXT section and prompt content, e.g.:

```
[RALPH] Iteration 1 of 5

CONTEXT:
Focus on unit tests

Refactor the module.
```

Exit 0.

**Verification:**
- "Focus on unit tests" appears in stdout between iteration line and "Refactor the module."
- Exit code 0

## Acceptance criteria

- [ ] With --dry-run or -d, Ralph resolves configuration, loads the prompt, assembles the preamble, and prints the complete assembled prompt to stdout
- [ ] No AI CLI process is spawned in dry-run mode
- [ ] The output shows exactly what would be piped to the AI CLI's stdin on the first iteration
- [ ] Configuration validation still runs in dry-run mode — invalid config produces errors before the prompt is displayed
- [ ] Dry-run exits with code 0 on success

## Dependencies

- O1/R8 (preamble injection) — defines the format and content of the assembled prompt (preamble + prompt content) for iteration 1; dry-run prints that same assembled prompt.
- O2 (config resolution and validation) — config is resolved before prompt load; validation (O2/R3) must pass before printing; fail-fast on invalid prompt source (O2/R4).
