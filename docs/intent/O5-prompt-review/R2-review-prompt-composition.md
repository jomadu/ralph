# R2: Review Prompt Composition

**Outcome:** O5 — Prompt review

## Requirement

The system composes the prompt given to the AI for the review phase by combining Ralph's review instructions (built-in or configured) with the user's prompt to be reviewed. The composed prompt instructs the AI where to write the report (review output path). The system does not interpret the AI's stdout to derive the report; the AI is instructed to write the report to a file at a known path. Review instructions are owned by Ralph (or configuration), not supplied by the user; the user supplies only the prompt content to be reviewed.

## Specification

(To be specified in Step 5.)

## Acceptance criteria

- [ ] The prompt sent to the AI for review includes both Ralph's review instructions and the user's prompt content.
- [ ] The review instructions direct the AI to write the report to a path that Ralph provides (the path is interpolated into the prompt).
- [ ] Ralph does not parse or interpret the AI's raw output to construct the report; the report is produced by the AI writing to the specified path.
- [ ] The source of review instructions is Ralph (built-in or configured), not the user's prompt file.
- [ ] The user's contribution to the review input is only the prompt to be reviewed (alias, file, or stdin content).

## Dependencies

- R3 (review output path) defines the path interpolated into the review prompt. R1 (input modes) supplies the user prompt content.
