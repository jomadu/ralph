# O5: Prompt review

## Statement

Ralph prompts can be reviewed for quality and structure before or without running the loop, with actionable feedback, an optional suggested revision of the entire prompt, and a machine-readable result.

## Why it matters

Without a reviewer, users discover prompt problems only when they run the loop: the AI never emits a success signal, or it fails repeatedly, or it does too much in one iteration. The user inspects output, guesses that the prompt lacks signal discipline or convergence criteria, and edits by trial and error. There is no structured way to check whether a prompt instructs the AI to emit success/failure correctly, references filesystem or work-tracking state, acknowledges the loop, or defines "done." A reviewer gives feedback before execution and enables CI or pre-commit checks so prompt quality can be evaluated without running the task.

## Verification

- User runs `ralph review <alias>`. Ralph loads the prompt for that alias, runs the reviewer (one-shot AI evaluation), and prints a report to stdout; exit code 0 (no errors), 1 (errors), or 2 (review failed to run).
- User runs `ralph review -f ./prompts/my-prompt.md`. Ralph reads the file, runs the reviewer, and prints the report; exit code reflects result.
- User pipes a prompt: `cat prompt.md | ralph review`. Ralph reads stdin, runs the reviewer, and prints the report; exit code reflects result.
- User runs `ralph review build --output report.txt`. The report is written to the file; stdout is silent; exit code still reflects result.
- User runs `ralph review build --apply`. After the report, Ralph prompts to apply the suggested revision to the prompt file; on confirmation, the revised content is written to the path for that alias.
- User runs `ralph review -f prompt.md --apply -y`. The suggested revision is applied to `prompt.md` without prompting (non-interactive). Exit code reflects review result; the file is updated when a suggested revision was present.
- The report includes narrative feedback (e.g. signal discipline, statefulness, scope, convergence), a machine-parseable summary so scripts or CI can gate on exit code, and may include a suggested revision of the entire prompt (full revised text for the user to copy or apply).

## Non-outcomes

- The reviewer does not run or modify the execution loop. It does not execute the user's task.
- The reviewer does not modify the user's prompt file unless the user requests it: with `--apply`, the user confirms (or uses `-y` in non-interactive mode) to write the suggested revision. Without `--apply`, the reviewer only reports; the user edits the prompt manually. `--apply` is only valid when the prompt was loaded from an alias or from a file path (`-f`), not from stdin (no destination).
- The reviewer does not enforce a single prompt style or template. It evaluates qualities that support Ralph's execution model (signals, state, iteration awareness, scope, convergence), not a fixed format.
- The reviewer does not replace human judgment on content or correctness. It checks structure and discipline relevant to loop behavior.
- The reviewer is not a general-purpose markdown or prose linter. Evaluation is tuned for Ralph prompts and Ralph's execution model (fresh process per iteration, filesystem state, preamble, signal scanning).

## Risks

| Risk | Mitigating Requirement |
|------|------------------------|
| User cannot supply the prompt from alias, file, or stdin in a consistent way | [R1 — Prompt input modes for review](R1-prompt-input-modes.md) |
| Review does not run (config invalid, prompt load failure, AI command missing or spawn fails) | [R9 — Review failure handling](R9-review-failure-handling.md) |
| AI output has no parseable summary line so exit code cannot be set correctly | [R4 — Review report content and parseable format](R4-report-content-format.md), [R5 — Machine-readable exit code](R5-machine-readable-exit-code.md) |
| Suggested revision cannot be extracted from output for apply or display | [R4 — Review report content and parseable format](R4-report-content-format.md) |
| User applies suggested revision by accident (no chance to decline) | [R7 — Apply confirmation](R7-apply-confirmation.md) |
| User expects --apply to work when prompt was piped from stdin | [R8 — Apply invalid for stdin](R8-apply-invalid-stdin.md) |
| Wrong file overwritten when applying (e.g. alias resolves elsewhere) | [R6 — Apply suggested revision to file](R6-apply-revision.md) |
| Report lacks narrative feedback or suggested revision in a usable form | [R4 — Review report content and parseable format](R4-report-content-format.md) |
| No AI command or config available for review subcommand | [R3 — Config and AI backend for review](R3-config-ai-backend.md) |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R1](R1-prompt-input-modes.md) | Prompt input modes for review | draft |
| [R2](R2-one-shot-execution.md) | One-shot review execution | draft |
| [R3](R3-config-ai-backend.md) | Config and AI backend for review | draft |
| [R4](R4-report-content-format.md) | Review report content and parseable format | draft |
| [R5](R5-machine-readable-exit-code.md) | Machine-readable exit code | draft |
| [R6](R6-apply-revision.md) | Apply suggested revision to file | draft |
| [R7](R7-apply-confirmation.md) | Apply confirmation | draft |
| [R8](R8-apply-invalid-stdin.md) | Apply invalid for stdin | draft |
| [R9](R9-review-failure-handling.md) | Review failure handling | draft |
