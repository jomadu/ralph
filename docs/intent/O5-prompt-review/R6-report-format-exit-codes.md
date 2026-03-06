# R6: Report Format and Exit Code Derivation

**Outcome:** O5 — Prompt review

## Requirement

The review report includes narrative feedback on the prompt (e.g. signal discipline, statefulness, scope, convergence), a machine-parseable summary so scripts or CI can gate on the result, and the full suggested revision. Ralph's exit code reflects the outcome: 0 when review completed with no errors (or only warnings if specified), 1 when review completed but the prompt has one or more errors, and 2 when the review failed to run or apply was invalid (e.g. config invalid, prompt load failure, AI spawn failure, or `--apply` with stdin and no `--prompt-output`).

## Specification

(To be specified in Step 5.)

## Acceptance criteria

- [ ] The report contains narrative feedback relevant to Ralph's execution model (e.g. signals, state, iteration awareness, scope, convergence).
- [ ] The report contains a machine-parseable summary so automation (scripts, CI) can determine pass/fail or error count.
- [ ] The report contains the full suggested revision of the prompt.
- [ ] Exit code 0: review completed successfully; no errors in the prompt (or only warnings if policy allows).
- [ ] Exit code 1: review completed; one or more errors were identified in the prompt.
- [ ] Exit code 2: review did not complete successfully (e.g. config invalid, prompt source missing or unreadable, AI not available or spawn failed) or apply was invalid (e.g. stdin + apply without `--prompt-output`).
- [ ] Exit code semantics are consistent with R8 (failure handling) and R9 (report file verification); missing report at expected path yields exit 2.

## Dependencies

- R8 defines which conditions produce exit 2. R9 defines report file verification; failure there also yields exit 2. The report format is produced by the AI per R2/R3; Ralph may parse the machine-parseable part for exit code derivation (to be specified in Step 5).
