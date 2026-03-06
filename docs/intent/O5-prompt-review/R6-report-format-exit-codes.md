# R6: Report Format and Exit Code Derivation

**Outcome:** O5 — Prompt review

## Requirement

The review report includes narrative feedback on the prompt (e.g. signal discipline, statefulness, scope, convergence), a machine-parseable summary so scripts or CI can gate on the result, and the full suggested revision. Ralph's exit code reflects the outcome: 0 when review completed with no errors (or only warnings if specified), 1 when review completed but the prompt has one or more errors, and 2 when the review failed to run or apply was invalid (e.g. config invalid, prompt load failure, AI spawn failure, or `--apply` with stdin and no `--prompt-output`).

## Specification

**Report content (three parts):**

1. **Narrative feedback** — Free-form or structured text evaluating the prompt for qualities relevant to Ralph's execution model: signal discipline (success/failure emission), statefulness and filesystem/state assumptions, scope and iteration awareness, convergence (how "done" is defined). Format is not strictly mandated; the AI produces human-readable feedback.
2. **Machine-parseable summary** — A structured block that scripts or CI can parse to determine pass/fail or error count. Format must be specified so implementers and parsers agree. Options: (a) a single line or block with a well-defined format (e.g. `ralph-review: status=ok|errors|warnings`, `errors=N`, `warnings=N`); (b) a small YAML/JSON block in the report file under a known heading (e.g. `## Summary` with key-value pairs). Ralph (or the report schema) must define one canonical format; the specification here is that such a block exists and is documented so exit code 0 vs 1 can be derived from it (see below).
3. **Full suggested revision** — The complete text of the prompt as the AI suggests it should be revised. It may appear inline in the report (e.g. under a "Suggested revision" section) or as a reference; the requirement is that the full revision is present in the report (or at the path the AI was instructed to write it to, when revision is written to file per R4). So either the report file contains the full revision, or the report references it and the AI has written it to `--prompt-output` when used — either way the "full suggested revision" is available to the user/CI.

**Exit code derivation:**

- **0** — Review completed successfully; the prompt has no errors (or only warnings, if policy treats warnings as non-fatal). Derived from the machine-parseable summary: e.g. `status=ok` or `errors=0` (and optionally `warnings` allowed).
- **1** — Review completed; one or more errors were identified in the prompt. Derived from the summary: e.g. `status=errors` or `errors>=1`. Ralph reads the report file (after R9 verification) and parses the machine-parseable section to set exit code 0 vs 1.
- **2** — Review did not complete successfully (config invalid, prompt source missing, AI spawn failed, report file missing per R9, invalid apply per R5/R8, etc.). Never 0 or 1 when any of those conditions hold.

Ralph must not exit 0 or 1 unless the report file exists at the expected path (R9). If R9 fails, exit 2. After R9 passes, Ralph parses the report file to extract the machine-parseable summary and sets exit code 0 or 1 accordingly. If the report file exists but the machine-parseable block is missing or unparseable, implementation must define behavior: e.g. treat as exit 1 (assume errors) or exit 2 (review output invalid); recommend exit 1 so CI fails safe.

**Schema for machine-parseable summary:** Define a single format. Example (minimal): a line matching `ralph-review:\s*status=(ok|errors|warnings)(\s+errors=(\d+))?(\s+warnings=(\d+))?` so that `status=ok` and `errors=0` imply exit 0; `status=errors` or `errors>0` imply exit 1; `status=warnings` with `errors=0` may imply exit 0 if policy allows. The exact regex or schema is implementation-defined but must be documented so the AI instructions (R2) can ask the AI to emit it and Ralph can parse it.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Report file exists but machine-parseable block missing or malformed | Treat as exit 1 (conservative) or exit 2 (invalid output); document. Recommend exit 1. |
| Report contains both narrative and revision; revision also written to `--prompt-output` | Acceptable; full revision is in report and/or at path. |
| Policy: warnings only, no errors | Exit 0 if implementation supports "warnings only" policy; else exit 1. Document policy. |

### Examples

#### Exit 0

**Input:** Review runs; report written; summary says `status=ok` or `errors=0`.

**Expected output:** Ralph exits 0.

#### Exit 1

**Input:** Review runs; report written; summary says `errors=2`.

**Expected output:** Ralph exits 1.

#### Exit 2 (report missing)

**Input:** Review phase completes but R9 finds no file at review output path.

**Expected output:** Ralph exits 2; no 0/1.

## Acceptance criteria

- [ ] The report contains narrative feedback relevant to Ralph's execution model (e.g. signals, state, iteration awareness, scope, convergence).
- [ ] The report contains a machine-parseable summary so automation (scripts, CI) can determine pass/fail or error count.
- [ ] The report contains the full suggested revision of the prompt.
- [ ] Exit code 0: review completed successfully; no errors in the prompt (or only warnings if policy allows).
- [ ] Exit code 1: review completed; one or more errors were identified in the prompt.
- [ ] Exit code 2: review did not complete successfully (e.g. config invalid, prompt source missing or unreadable, AI not available or spawn failed) or apply was invalid (e.g. stdin + apply without `--prompt-output`).
- [ ] Exit code semantics are consistent with R8 (failure handling) and R9 (report file verification); missing report at expected path yields exit 2.

## Dependencies

- R8 defines which conditions produce exit 2. R9 defines report file verification; failure there also yields exit 2. The report format is produced by the AI per R2/R3; Ralph parses the machine-parseable part for exit code derivation (see Specification above).
