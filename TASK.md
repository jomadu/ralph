# TASK: Prompt Review Feature

## Overview

Add a `ralph review` subcommand that performs static analysis on Ralph prompts using an AI to evaluate prompt quality and provide improvement recommendations.

## Motivation

Ralph prompts have specific characteristics that make them effective for iterative loop execution:
- Signal discipline (when to emit SUCCESS/FAILURE)
- Statefulness awareness (references to filesystem, work tracking)
- Iteration awareness (acknowledges loop context)
- Scope control (prevents doing everything at once)
- Convergence criteria (defines "done")

A reviewer helps users write better prompts by evaluating these qualities and providing actionable feedback.

## Design

### Architecture

The reviewer is itself an AI prompt that evaluates other prompts. This approach:
- Leverages AI's natural language understanding (no brittle regex rules)
- Adapts flexibly to different prompt styles
- Provides narrative feedback, not just pass/fail

### Execution Flow

```
ralph review <alias>  [--apply] [-y]
  ↓
1. Load target prompt (alias/file/stdin)
2. Load embedded review prompt (internal/review package, e.g. prompt.md via go:embed)
3. Resolve signal strings if reviewing an alias
4. Assemble: review prompt + context + target prompt
5. Spawn AI CLI once (no loop, no preamble, no signal scanning)
6. Capture output → stdout or --output file
7. Parse summary line (if present) and set exit code: 0 (no errors), 1 (errors), or 2 (review failed to run)
8. If --apply: require a writable path (alias or -f; stdin is invalid). Parse output for suggested revision block. If interactive and not -y: prompt "Apply suggested revision to <path>? [y/N]"; if user confirms (or -y), write revised content to that path. If no suggested revision in output, --apply is a no-op or exit 2 (implementation choice).
```

## Command Interface

### Syntax

```bash
ralph review <alias>                    # Review a configured prompt alias
ralph review -f prompt.md               # Review a file directly
cat prompt.md | ralph review            # Review from stdin
ralph review build --output report.txt  # Write to file instead of stdout
ralph review build --apply              # After review, prompt to apply suggested revision to the prompt file
ralph review build -f prompt.md --apply -y  # Apply suggested revision without confirmation (non-interactive)
ralph review build --ai-cmd-alias claude # Override AI backend
ralph review build --ai-cmd "custom-cli --flags" # Direct command override
```

### Flags

- `-c, --config <path>` - Same as `ralph run`; explicit config file path (optional). When set, used as sole file-based config; otherwise global + workspace discovery.
- `-f, --file <path>` - Read prompt from file (no alias required)
- `--output <path>` - Write review report to file instead of stdout (silent stdout)
- `--apply` - After the review, apply the suggested revision to the prompt file. Only valid when the prompt was loaded from an alias or from `-f <path>` (writes to that path). Invalid when input is stdin (no destination path). If the AI output contains no suggested revision block, Ralph does not overwrite the file (no-op or exit 2; implementation choice). When run interactively (TTY), Ralph prompts for confirmation before writing unless `-y`/`--yes` is set.
- `-y, --yes` - Non-interactive mode for `--apply`: skip confirmation and write the suggested revision to the prompt file. Ignored if `--apply` is not set. Intended for scripts and CI.
- `--ai-cmd <string>` - Direct AI command string
- `--ai-cmd-alias <string>` - AI command alias

### Prompt Resolution

Same as `ralph run`:
1. Positional alias → resolve from config `prompts.<alias>.path`
2. `-f <path>` → read file directly
3. stdin → read from pipe

### AI Backend Resolution

Same as `ralph run`:
1. CLI flags (`--ai-cmd` or `--ai-cmd-alias`) take precedence
2. Fall back to config `loop.ai_cmd` or `loop.ai_cmd_alias`
3. Error if no command configured

### Exit Codes

Machine-readable exit codes for scripts and CI:

| Code | Meaning |
|------|---------|
| 0 | No errors (warnings and suggestions are allowed; review ran and found no errors) |
| 1 | One or more **errors** (e.g. missing signal discipline) |
| 2 | Review failed to run (config error, prompt load failure, AI spawn failure, etc.) |

Ralph parses the reviewer's output for a summary line (see **Review output contract** below) to set 0 vs 1. Tool/config failures exit 2 before or without relying on AI output.

## Review Prompt

### Location

The review prompt lives with the implementation: inside the `internal/review` package (e.g. `internal/review/prompt.md`), embedded at build time via `//go:embed`.

### Content Structure

The review prompt should instruct the AI to evaluate:

**Core Criteria:**
- **Signal discipline**: Does the prompt instruct when to emit SUCCESS/FAILURE/nothing?
- **Statefulness awareness**: Does it tell the AI where to find state (files, work tracking)?
- **Iteration awareness**: Does it acknowledge it's in a loop?
- **Scope control**: Does it prevent the AI from trying to do everything at once?
- **Convergence**: Does it define what "done" means?

**Anti-patterns:**
- Conversational language ("let's discuss...")
- Missing success criteria
- Ambiguous completion conditions
- Instructions that assume conversation history

**Severity Levels:**
- **Error**: Missing signal discipline entirely
- **Warning**: Conversational language, missing convergence criteria, no scope control
- **Suggestion**: Could be more explicit about iteration awareness, could reference specific file paths

**Output Format:**
- Summary section (✓/⚠/✗ for each criterion)
- Detailed narrative with recommendations
- Error/warning/suggestion counts
- **Suggested revision**: A full revised version of the user's prompt that incorporates the reviewer's recommendations. The user can copy this into their prompt file; the reviewer does not write to disk. When present, this section should be clearly delimited (e.g. a "Suggested revision" or "Revised prompt" heading and the full prompt text) so it can be extracted or skipped by parsers. The machine-readable summary line (see below) must still appear in a parseable position (e.g. after the revision block or in a fixed format) so Ralph can set the exit code.

**Review output contract (for machine-readable exit code):**  
The review prompt must require the AI to emit a single parseable summary line so Ralph can set exit 0 vs 1. For example, a line of the form `REVIEW: errors=N warnings=M suggestions=P` (with N, M, P non-negative integers) on a single line, or an equivalent convention (e.g. a small JSON block). Ralph parses this line/block after capture; if the line is missing or unparseable, treat as 0 errors (exit 0) so that older or non-conforming review prompts still succeed. The exact format is specified when implementing the review prompt and parser.

Example output:
```
Analyzing prompt: build (./prompts/build.md)

✓ Signal discipline: Found explicit success/failure conditions
✓ Statefulness: References work tracking and filesystem
⚠ Iteration awareness: Implicit only (no explicit mention of iteration count)
✓ Scope control: Enforces "one task per iteration"
✓ Convergence: Clear completion criteria

1 warning, 0 errors

Overall: GOOD

Recommendations:
- Consider explicitly acknowledging iteration count in the prompt

Suggested revision:
--- BEGIN REVISED PROMPT ---
<full revised prompt text here>
--- END REVISED PROMPT ---

REVIEW: errors=0 warnings=1 suggestions=0
```

### Context Injection

When reviewing an alias, include signal string context:

```
The user has configured these signal strings:
- Success: <promise>SUCCESS</promise>
- Failure: <promise>FAILURE</promise>

Evaluate whether the prompt correctly instructs the AI to use these signals.
```

When reviewing a file/stdin (no alias), omit signal context:

```
The user may configure custom signal strings. Evaluate whether the prompt
provides clear signal discipline regardless of the specific strings used.
```

### Ralph Architecture Context

The review prompt should explain Ralph's execution model:
- Fresh AI process per iteration (no conversation history)
- State continuity via filesystem only
- Preamble injection (iteration count, optional context)
- Signal scanning determines loop outcome
- Max iterations and failure threshold limits

## Implementation Notes

### Config and --config

`ralph review` uses the same config loader and `--config` behavior as `ralph run`: `LoadConfigWithProvenanceAndExplicit(configFlag)`. When `--config <path>` is provided, that path is the sole file-based config; when omitted, global and workspace config files are discovered. This keeps behavior consistent across subcommands.

### No Loop Execution

Unlike `ralph run`, `ralph review` spawns the AI CLI exactly once:
- No iteration loop
- No preamble injection
- No signal scanning
- No max iterations or failure threshold

It's a simple one-shot execution: assemble prompt → spawn AI → capture output → display.

### No Verbose/Quiet Flags

`ralph review` does not respect `-v`, `-q`, or `--log-level` flags. The AI's output is the review report and should not be mixed with Ralph's own logging.

### Signal String Resolution

When reviewing an alias:
1. Check for per-prompt signal overrides (`prompts.<alias>.loop.signals`)
2. Fall back to global config (`loop.signals`)
3. Include resolved strings in review prompt context

When reviewing a file/stdin:
- No signal strings available
- Review prompt evaluates generically

### Output Handling

- Default: Write AI output to stdout
- `--output <path>`: Write AI output to file, silent stdout
- Ralph parses the output for the review summary line/block (see **Review output contract**) to set exit code 0 vs 1; the rest of the AI response is passed through verbatim to stdout or the output file

### Apply and Confirmation

- When `--apply` is set and the prompt was loaded from an alias or `-f <path>`, Ralph parses the captured output for the suggested revision block (e.g. content between `--- BEGIN REVISED PROMPT ---` and `--- END REVISED PROMPT ---`). If found:
  - **Interactive (TTY):** Ralph prompts the user, e.g. `Apply suggested revision to <path>? [y/N]`. Accept `y`/`yes` (case-insensitive) as confirm; anything else or Enter declines. On confirm, write the revised content to the path; on decline, do nothing and exit with the same code as the review (0 or 1).
  - **Non-interactive (`-y`/`--yes`):** Skip the prompt and write the revised content to the path.
- If `--apply` is set but input was stdin, Ralph exits 2 (no destination path).
- If `--apply` is set but the AI output contains no parseable suggested revision, Ralph does not write; exit 2 or no-op with review exit code (implementation choice; recommend exit 2 so scripts can detect missing revision).

## Decisions

**Q1: Where does the review prompt file live during development?**  
**Decision:** With the implementation. The review prompt lives in the `internal/review` package (e.g. `internal/review/prompt.md`), embedded via `//go:embed`, so prompt and assembly logic stay together.

**Q2: Package structure?**  
**Decision:** Option A — a new `internal/review` package with the embedded prompt and assembly logic (load target prompt, resolve signal context, assemble review prompt + context + target, one-shot spawn, parse summary for exit code).

**Q3: Should this be a new outcome (O5) in the intent tree?**  
**Decision:** Yes. Add a new outcome O5 with full decomposition: outcome statement, risks, requirements, and specifications, following the methodology in `building-intent.md`.

## Next Steps

1. **Build the intent (O5)** — Go through the building-intent process (`building-intent.md`): start with outcome O5 (one-line in the intent index, then outcome detail). Derive requirements from this task using the why/how/how-else chain and risk analysis. Expand each requirement with specifications (schemas, formats, edge cases, examples). Lock each step with consistency review before expanding.
2. **Construct an implementation plan** — Create a plan (e.g. extend `PLAN.md` or a dedicated plan) that integrates the O5 requirements into the implementation: phased tasks, dependency order, and spec references, so that building the reviewer follows the same pattern as the existing phases.
3. **Implement** — Write the review prompt in `internal/review`, implement the package and `ralph review` subcommand per the plan, test with existing prompts, and document in README.md.

## References

- Existing prompt: `prompts/build.md` (complex, well-structured example)
- Config system: `internal/config/loader.go` (for signal resolution)
- Prompt loading: `internal/prompt/loader.go` (reuse for target prompt)
- Command parsing: `internal/cmdparse/parser.go` (for AI backend)
- Process spawning: `internal/runner/spawn.go` (one-shot execution)
