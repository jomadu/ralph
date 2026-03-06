# TASK: Prompt Linter Feature

## Overview

Add a `ralph lint` subcommand that performs static analysis on Ralph prompts using an AI to evaluate prompt quality and provide improvement recommendations.

## Motivation

Ralph prompts have specific characteristics that make them effective for iterative loop execution:
- Signal discipline (when to emit SUCCESS/FAILURE)
- Statefulness awareness (references to filesystem, work tracking)
- Iteration awareness (acknowledges loop context)
- Scope control (prevents doing everything at once)
- Convergence criteria (defines "done")

A linter helps users write better prompts by evaluating these qualities and providing actionable feedback.

## Design

### Architecture

The linter is itself an AI prompt that evaluates other prompts. This approach:
- Leverages AI's natural language understanding (no brittle regex rules)
- Adapts flexibly to different prompt styles
- Provides narrative feedback, not just pass/fail

### Execution Flow

```
ralph lint <alias>
  ↓
1. Load target prompt (alias/file/stdin)
2. Load embedded linter prompt (prompts/lint.md via go:embed)
3. Resolve signal strings if linting an alias
4. Assemble: linter prompt + context + target prompt
5. Spawn AI CLI once (no loop, no preamble, no signal scanning)
6. Capture output → stdout or --output file
7. Exit 0
```

## Command Interface

### Syntax

```bash
ralph lint <alias>                    # Lint a configured prompt alias
ralph lint -f prompt.md               # Lint a file directly
cat prompt.md | ralph lint            # Lint from stdin
ralph lint build --output report.txt # Write to file instead of stdout
ralph lint build --ai-cmd-alias claude # Override AI backend
ralph lint build --ai-cmd "custom-cli --flags" # Direct command override
```

### Flags

- `-f, --file <path>` - Read prompt from file (no alias required)
- `--output <path>` - Write lint report to file instead of stdout (silent stdout)
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

Always exit 0 (linter ran successfully). The lint findings are in the output, not the exit code.

## Linter Prompt

### Location

`prompts/lint.md` in the repository, embedded at build time via `//go:embed`.

### Content Structure

The linter prompt should instruct the AI to evaluate:

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
```

### Context Injection

When linting an alias, include signal string context:

```
The user has configured these signal strings:
- Success: <promise>SUCCESS</promise>
- Failure: <promise>FAILURE</promise>

Evaluate whether the prompt correctly instructs the AI to use these signals.
```

When linting a file/stdin (no alias), omit signal context:

```
The user may configure custom signal strings. Evaluate whether the prompt
provides clear signal discipline regardless of the specific strings used.
```

### Ralph Architecture Context

The linter prompt should explain Ralph's execution model:
- Fresh AI process per iteration (no conversation history)
- State continuity via filesystem only
- Preamble injection (iteration count, optional context)
- Signal scanning determines loop outcome
- Max iterations and failure threshold limits

## Implementation Notes

### No Loop Execution

Unlike `ralph run`, `ralph lint` spawns the AI CLI exactly once:
- No iteration loop
- No preamble injection
- No signal scanning
- No max iterations or failure threshold

It's a simple one-shot execution: assemble prompt → spawn AI → capture output → display.

### No Verbose/Quiet Flags

`ralph lint` does not respect `-v`, `-q`, or `--log-level` flags. The AI's output is the lint report and should not be mixed with Ralph's own logging.

### Signal String Resolution

When linting an alias:
1. Check for per-prompt signal overrides (`prompts.<alias>.loop.signals`)
2. Fall back to global config (`loop.signals`)
3. Include resolved strings in linter prompt context

When linting a file/stdin:
- No signal strings available
- Linter prompt evaluates generically

### Output Handling

- Default: Write AI output to stdout
- `--output <path>`: Write AI output to file, silent stdout
- No parsing or interpretation of the AI's response
- Pass through verbatim

## Open Questions

**Q1: Where does the linter prompt file live during development?**
- Option A: `prompts/lint.md` (alongside user prompts)
- Option B: `internal/lint/prompt.md` (with implementation)

**Q2: Package structure?**
- Option A: New `internal/lint` package with embedded prompt and assembly logic
- Option B: Add `LintPrompt()` to `internal/runner`
- Option C: Implement directly in `cmd/ralph/main.go`

**Q3: Should this be a new outcome (O5) in the intent tree?**
- Requires full decomposition: outcome → requirements → specifications
- Or is this a feature addition to existing outcomes?

## Next Steps

1. Decide on open questions (Q1-Q3)
2. Write the linter prompt (`prompts/lint.md`)
3. Add outcome O5 to intent tree (if needed)
4. Implement `ralph lint` subcommand
5. Test with existing prompts (build.md, etc.)
6. Document in README.md

## References

- Existing prompt: `prompts/build.md` (complex, well-structured example)
- Config system: `internal/config/loader.go` (for signal resolution)
- Prompt loading: `internal/prompt/loader.go` (reuse for target prompt)
- Command parsing: `internal/cmdparse/parser.go` (for AI backend)
- Process spawning: `internal/runner/spawn.go` (one-shot execution)
