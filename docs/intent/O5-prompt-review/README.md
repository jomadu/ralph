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
