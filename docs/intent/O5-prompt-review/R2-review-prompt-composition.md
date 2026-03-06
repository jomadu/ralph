# R2: Review Prompt Composition

**Outcome:** O5 — Prompt review

## Requirement

The system composes the prompt given to the AI for the review phase by combining Ralph's review instructions (embedded in the binary) with the user's prompt to be reviewed. The composed prompt instructs the AI where to write the report (review output path). The system does not interpret the AI's stdout to derive the report; the AI is instructed to write the report to a file at a known path. Review instructions are owned by Ralph (embedded in the binary), not supplied by the user; the user supplies only the prompt content to be reviewed.

## Specification

**Composition structure:** The prompt sent to the AI for the review phase is a single document (or stream) formed by concatenating or otherwise combining, in order:

1. **Ralph's review instructions** — Text that tells the AI how to evaluate the prompt (e.g. criteria: signal discipline, statefulness, scope, convergence, iteration awareness) and what to produce (narrative feedback, machine-parseable summary, full suggested revision). Source: files embedded in the Ralph binary at build time via Go's `embed` (e.g. `//go:embed` in the Ralph codebase); not read from the user's repository. The user does not supply the instructions — only the prompt content to be reviewed.
2. **Path directive** — A clear instruction to the AI to write the report to a specific path. The path is the review output path from R3 (either `--review-output <path>` or the chosen temp path). It must be interpolated into the prompt as a concrete path string (e.g. "Write your report to the following path: <path>") so the AI knows exactly where to write. The AI is expected to write the report to that path (file); Ralph does not parse stdout to construct the report.
3. **User's prompt to be reviewed** — The exact content loaded per R1 (alias, file, or stdin). No modification or wrapping beyond inclusion in the composed prompt.

Ralph does not interpret the AI's raw stdout to derive the report. The report is defined as the content the AI writes to the file at the review output path. If the AI also prints to stdout, that is separate (and may be shown to the user per R7); the canonical report is the file.

**Source of review instructions:** One or more files embedded in the Ralph binary at build time using Go's `embed` (e.g. under an `internal/review` package, with `//go:embed`). These files are part of the Ralph distribution; they are not read from the repository in which the user runs `ralph review`. The user does not supply the review instructions — only the prompt content to be reviewed.

**Revision-phase prompt:** When apply is requested, the prompt for the revision phase (instructing the AI to read the report and write the revised prompt to the prompt output path) also uses embedded content: a revision-instructions file embedded in the binary via Go `embed`, with the review output path and prompt output path interpolated. Revision instructions are not read from the user's repository.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User prompt content is very large | No truncation required by this requirement; pass through. Implementation may impose limits (document). |
| Review output path contains special characters or spaces | Interpolate the path as given; the AI is instructed to write to that path; escaping is implementation-defined (e.g. quote in the prompt). |

### Examples

#### Review phase prompt shape

**Input:** User runs `ralph review build`; alias `build` resolves to `./prompts/build.md`; review output path is `/tmp/ralph-report-xyz.md`.

**Expected output:** The AI receives a prompt that (1) contains Ralph's review instructions, (2) says "Write your report to: /tmp/ralph-report-xyz.md" (or equivalent), (3) contains the full content of `./prompts/build.md`. The AI writes the report to `/tmp/ralph-report-xyz.md`; Ralph does not parse AI stdout to build the report.

## Acceptance criteria

- [ ] The prompt sent to the AI for review includes both Ralph's review instructions and the user's prompt content.
- [ ] The review instructions direct the AI to write the report to a path that Ralph provides (the path is interpolated into the prompt).
- [ ] Ralph does not parse or interpret the AI's raw output to construct the report; the report is produced by the AI writing to the specified path.
- [ ] The source of review (and revision) instructions is Ralph — embedded in the binary via Go embed — not the user's prompt file or the run-time repository.
- [ ] The user's contribution to the review input is only the prompt to be reviewed (alias, file, or stdin content).

## Dependencies

- R3 (review output path) defines the path interpolated into the review prompt. R1 (input modes) supplies the user prompt content.
