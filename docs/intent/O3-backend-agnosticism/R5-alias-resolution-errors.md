# R5: AI Command Alias Resolution with Clear Errors

**Outcome:** O3 — Backend Agnosticism

## Requirement

The system resolves AI command aliases to executable command strings and produces clear, actionable error messages when resolution fails. The user should never see a cryptic "command not found" — the error should explain what alias was requested, that it wasn't found, and what aliases are available.

## Specification

_To be specified._

### Edge cases

_To be specified._

### Examples

_To be specified._

## Acceptance criteria

- [ ] When a valid alias is specified, Ralph resolves it to the corresponding command string
- [ ] When an unknown alias is specified, Ralph exits with an error that names the unknown alias and lists available aliases
- [ ] When no AI command is configured (no alias and no direct command), Ralph exits with an error explaining that an AI command must be configured
- [ ] Resolution errors occur at startup, before the loop begins

## Dependencies

_None identified._
