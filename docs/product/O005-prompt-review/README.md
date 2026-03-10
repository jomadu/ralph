# O005: Prompt Review

## Who

Users who want to check prompt quality and structure before or without running the loop — including CI or pre-commit — and who may want a suggested revision applied to the prompt file.

## Statement

Prompts can be reviewed for quality and structure before or without running the loop; the user gets a report and an optional suggested revision, and can request that the revision be written to a file (with confirmation when appropriate).

## Why it matters

Without a reviewer, problems show up only when the loop runs: the AI never emits a success signal, fails repeatedly, or does too much in one iteration. The user inspects output and edits by trial and error. There is no structured way to check whether a prompt instructs the AI to emit success/failure correctly, references filesystem or work-tracking state, or defines "done." A reviewer gives feedback before execution and enables CI or pre-commit checks. Configurable outputs (report, revised prompt) let users save results and, when the prompt is supplied via stdin, obtain a file with the recommended changes.

## Verification

- User runs the review command with the prompt supplied by alias, file path, or stdin; receives a report that includes a suggested revision. The user can request that the revision be written to a file; when doing so, the user confirms (or uses a non-interactive option where supported). Exit codes distinguish: review completed with no errors, review completed but the prompt has errors, or review or apply did not complete successfully.
- The report is always saved to a file. The user can choose where the report is written or accept a default (e.g. a temporary location). The user can direct the revised prompt to a chosen path.
- When the prompt was supplied via stdin and the user requests that the revision be written, the user must specify where to write it (there is no source file to overwrite); if they do not, the system reports an error and does not apply.
- The report includes narrative feedback and a machine-parseable summary so scripts or CI can gate on the result.
- The review evaluates prompts along these dimensions:
  - **Signal and state** — Clear success and failure signals Ralph can detect, and statefulness (e.g. filesystem, work-tracking) that works with the loop model.
  - **Iteration awareness** — The prompt acknowledges that execution is multi-iteration with a fresh process each time, so the AI can behave accordingly (e.g. re-read state, emit signals).
  - **Scope and convergence** — The task has a defined scope and completion criteria that are checkable in practice, so the loop can converge rather than run indefinitely.
  - **Subjective completion criteria** — When "done" is subjective (e.g. "good enough" or "reads well"), the prompt includes techniques that help escape local optima: variation, creative exploration, or brainstorming and stepping back (e.g. consider alternatives, challenge assumptions) so the AI does not get stuck in small repetitive tweaks or a mediocre solution.

## Non-outcomes

- The reviewer does not run or modify the execution loop. It does not execute the user's task.
- The review (and revision) instructions are Ralph's — embedded in the binary; not read from the user's repository. The user supplies only the prompt to be reviewed.
- The reviewer does not modify any file unless the user explicitly requests that the revision be applied and confirms (or uses a non-interactive option). Without that request, the reviewer only reports.
- The reviewer does not enforce a single prompt style or template. It evaluates qualities that support Ralph's execution model (signals, state, iteration awareness, scope, convergence), not a fixed format.
- The reviewer is not a general-purpose markdown or prose linter. Evaluation is tuned for Ralph prompts and Ralph's execution model.
