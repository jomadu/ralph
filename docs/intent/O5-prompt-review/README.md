# O5: Prompt review

## Statement

Prompts can be reviewed for quality and structure before or without running the loop. The user supplies the prompt by alias, file path, or stdin (e.g. pipe).

Ralph composes a review prompt that instructs the AI to produce the report — narrative feedback, machine-parseable summary, and a full suggested revision — and to write it to the appropriate location. Ralph runs one AI invocation. The review prompt is constructed so the AI outputs the report there; Ralph does not interpret the AI output to derive the report. The report is always saved to a file: the path from `--review-output` or, when not set, the system temporary storage (e.g. the temp directory). By default the output of the AI command and the report are exposed to stdout; this is configurable. The user can direct the revised prompt to a path via `--prompt-output`.

To update the source file with the suggested revision, the user requests apply (confirmation or `--apply -y`). Ralph may run the AI again (revision phase) to produce the revision. The review output path (from `--review-output` or the temp directory) and the prompt output path (from `--prompt-output`, or the source file path when applying to an alias or file) are interpolated into the prompt given to the AI during the revision phase, so the AI is instructed where to read the report from and where to write the revised prompt. When the prompt came from stdin, there is no source file to overwrite; with `--apply`, `--prompt-output` is required so that path can be interpolated into the revision-phase prompt. Apply otherwise works as normal.

## Why it matters

Without a reviewer, problems show up only when the loop runs: the AI never emits a success signal, fails repeatedly, or does too much in one iteration. The user inspects output, guesses that the prompt lacks signal discipline or convergence criteria, and edits by trial and error. There is no structured way to check whether a prompt instructs the AI to emit success/failure correctly, references filesystem or work-tracking state, acknowledges the loop, or defines “done.” A reviewer gives feedback before execution and enables CI or pre-commit checks so prompt quality can be evaluated without running the task. Configurable outputs (report, revised prompt) let users save results and—when using stdin—obtain a file with the recommended changes.

## Verification

User runs `ralph review` with prompt from alias, `-f <path>`, or stdin; receives a report (including suggested revision); with `--apply`, can write revision to the prompt file (confirm or `-y`); exit 0 (ok), 1 (errors in prompt), 2 (review failed or apply invalid).

**Input (alias, file path, or stdin)**

- User runs `ralph review <alias>`. Ralph loads the prompt for that alias and runs the reviewer; the report is always written to a file (see Output paths); by default the AI command output and the report are exposed to stdout (configurable). Exit code reflects result.
- User runs `ralph review -f ./prompts/my-prompt.md`. Ralph reads the file, runs the reviewer; report is written to a file; by default the AI command output and the report are exposed to stdout (configurable).
- User runs `cat prompt.md | ralph review` (or pipes prompt via stdin). Ralph reads stdin, runs the reviewer; report is written to a file; by default the AI command output and the report are exposed to stdout (configurable). Exit code reflects result.

**Output paths**

- The report is always saved to a file. When `--review-output` is set (e.g. `ralph review build --review-output report.txt`), the report is written to that path. When `--review-output` is not set, the report is written to the system temporary storage (e.g. the temp directory). By default the output of the AI command and the report are exposed to stdout; this is configurable. Exit code reflects result.
- User runs `ralph review -f prompt.md --prompt-output prompt-revised.md`. The suggested revised prompt is written to the specified path; the source file is not overwritten.
- User runs `cat prompt.md | ralph review --prompt-output revised.md`. The suggested revised prompt is written to `revised.md`. When input is stdin and the user requests apply, `--prompt-output` is required.

**Apply and output paths**

The review output path (`--review-output` or the temp directory) and the prompt output path (`--prompt-output` or the source file when applying) are interpolated into the prompt given to the AI during the revision phase.

- User runs `ralph review build --apply`. After the report, Ralph prompts to apply the suggested revision to the prompt file; on confirmation, the revised content is written to the path for that alias.
- User runs `ralph review -f prompt.md --apply -y`. The suggested revision is applied to `prompt.md` without prompting; exit code reflects review result; the file is updated with the suggested revision.
- User runs `ralph review` with prompt from stdin and `--apply --prompt-output revised.md`. Apply works as normal; `--prompt-output` is required so Ralph knows where to write the revised prompt (there is no source file to overwrite). If `--apply` is used with stdin but without `--prompt-output`, the system reports an error (e.g. exit code 2).

**Report and exit codes**

- The report includes narrative feedback (e.g. signal discipline, statefulness, scope, convergence), a machine-parseable summary so scripts or CI can gate on the result, and the full suggested revision.
- Exit code 0: review completed, no errors (or only warnings if specified). Exit code 1: review completed, one or more errors in the prompt. Exit code 2: review failed to run (config invalid, prompt load failure, AI missing or spawn failed) or apply invalid (e.g. `--apply` with stdin but without `--prompt-output`).

**Workflow**

- Ralph composes the review prompt (Ralph’s instructions plus the user’s prompt to be reviewed) so that it tells the AI where to write the report: the path from `--review-output` or, if unset, the system temporary storage (e.g. the temp directory). Ralph runs one AI process; the AI produces the report and writes it to that location. Ralph does not interpret the AI output to derive the report. The report is always saved to a file; by default the AI command output and the report are exposed to stdout (configurable). For the revision phase of the review (when the revised prompt is produced and written), Ralph interpolates both the review output path (from `--review-output` or the temp directory) and the prompt output path into the prompt given to the AI: the prompt output path is either the source file (when applying to an alias or `-f` file) or the path from `--prompt-output` (required when applying with stdin). The AI is instructed where to read the report from and where to write the revised prompt.

## Non-outcomes

- The reviewer does not run or modify the execution loop. It does not execute the user’s task.
- The review instructions (the prompt that tells the AI how to evaluate) are Ralph’s (built-in or configured), not the user’s; the user supplies only the prompt to be reviewed.
- The reviewer does not modify the user’s prompt file unless the user requests apply and confirms (or uses `-y`). Without apply, the reviewer only reports; the user edits manually. Apply is valid for alias, file path, or stdin; when the prompt is from stdin, `--prompt-output` is required with `--apply` to specify where to write the revised prompt.
- The reviewer does not enforce a single prompt style or template. It evaluates qualities that support Ralph’s execution model (signals, state, iteration awareness, scope, convergence), not a fixed format.
- The reviewer does not replace human judgment on content or correctness. It checks structure and discipline relevant to loop behavior.
- The reviewer is not a general-purpose markdown or prose linter. Evaluation is tuned for Ralph prompts and Ralph’s execution model (fresh process per iteration, filesystem state, preamble, signal scanning).

## Risks

| Risk | Mitigating Requirement |
|------|------------------------|
| User runs review with stdin and `--apply` but omits `--prompt-output` | [R5 — Apply with confirmation and revision phase](R5-apply-confirmation-revision-phase.md) |
| Report path (`--review-output` or temp) is unwritable or invalid | [R3 — Review output path](R3-review-output-path.md), [R8 — Review failure handling](R8-review-failure-handling.md) |
| AI does not write the report to the specified path | [R9 — Report file verification](R9-report-file-verification.md) |
| Revision phase fails (AI does not write revised prompt to path) | [R5 — Apply with confirmation and revision phase](R5-apply-confirmation-revision-phase.md), [R8 — Review failure handling](R8-review-failure-handling.md) |
| User does not know where the report was written when `--review-output` is unset | [R3 — Review output path](R3-review-output-path.md) |
| Review instructions (built-in or configured) missing or wrong | [R2 — Review prompt composition](R2-review-prompt-composition.md) |
| Config invalid, prompt source missing, or AI spawn fails | [R8 — Review failure handling](R8-review-failure-handling.md) |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R1](R1-review-command-input-modes.md) | Review command with alias, file, and stdin input modes | ready |
| [R2](R2-review-prompt-composition.md) | Review prompt composition (instructions + user prompt, path interpolation) | ready |
| [R3](R3-review-output-path.md) | Review output path (`--review-output` or temp; report always to file) | ready |
| [R4](R4-prompt-output-path.md) | Prompt output path (`--prompt-output` for revised prompt; required when apply + stdin) | ready |
| [R5](R5-apply-confirmation-revision-phase.md) | Apply with confirmation and revision phase (interpolation, stdin+apply validation) | ready |
| [R6](R6-report-format-exit-codes.md) | Report format and exit code derivation (narrative, machine-parseable, full revision; exit 0/1/2) | ready |
| [R7](R7-configurable-review-stdout.md) | Configurable review output to stdout | ready |
| [R8](R8-review-failure-handling.md) | Review failure handling (invalid config, missing prompt, spawn failure, invalid apply → exit 2) | ready |
| [R9](R9-report-file-verification.md) | Report file verification (report exists at expected path after run; else exit 2) | ready |
