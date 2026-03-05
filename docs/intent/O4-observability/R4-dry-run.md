# R4: Dry-Run Mode

**Outcome:** O4 — Observability

## Requirement

The system prints the fully assembled prompt — preamble plus prompt content — to stdout without spawning an AI process, enabling the user to validate prompt assembly and configuration before committing to a run.

## Specification

_To be specified._

### Edge cases

_To be specified._

### Examples

_To be specified._

## Acceptance criteria

- [ ] With --dry-run or -d, Ralph resolves configuration, loads the prompt, assembles the preamble, and prints the complete assembled prompt to stdout
- [ ] No AI CLI process is spawned in dry-run mode
- [ ] The output shows exactly what would be piped to the AI CLI's stdin on the first iteration
- [ ] Configuration validation still runs in dry-run mode — invalid config produces errors before the prompt is displayed
- [ ] Dry-run exits with code 0 on success

## Dependencies

_None identified._
