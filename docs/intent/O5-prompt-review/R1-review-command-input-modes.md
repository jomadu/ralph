# R1: Review Command Input Modes

**Outcome:** O5 — Prompt review

## Requirement

The system accepts the prompt to be reviewed from three sources: a configured prompt alias, a file path (e.g. `-f <path>`), or standard input (e.g. piped content). The `ralph review` command resolves the prompt from exactly one of these sources per invocation and uses it as the input to the review workflow. Ralph loads the prompt once and does not re-read it during the review.

## Specification

(To be specified in Step 5.)

## Acceptance criteria

- [ ] User can run `ralph review <alias>`; Ralph loads the prompt associated with that alias and runs the reviewer with that content.
- [ ] User can run `ralph review -f <path>`; Ralph reads the prompt from the file at that path and runs the reviewer.
- [ ] User can pipe prompt content into `ralph review` (e.g. `cat prompt.md | ralph review`); Ralph reads from stdin and runs the reviewer with that content.
- [ ] Exactly one input source is used per invocation; alias, file path, and stdin are mutually exclusive in resolution (precedence or flag semantics are specified elsewhere).
- [ ] The prompt is loaded once at the start of the review and not re-read from the source during the same run.

## Dependencies

- Configuration and alias resolution (O2). Prompt source validation and fail-fast behavior (O2 R4) apply when alias or file is used.
