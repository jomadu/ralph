# R4: Fail-Fast on Invalid Prompt Source

**Outcome:** O2 — Configurable Behavior

## Requirement

The system verifies that the prompt source is valid and produces usable content before starting the loop. Missing files, unreadable files, and empty input cause an immediate, clear error rather than a failure on the first iteration. This applies to all prompt input modes: alias, file flag, and stdin.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] When `ralph run <alias>` is invoked, Ralph checks that the mapped prompt file exists and is readable before starting the loop
- [ ] If the alias's file does not exist, Ralph exits with an error message naming the missing file and the alias that referenced it
- [ ] When `ralph run -f <path>` is invoked, Ralph checks that the specified file exists and is readable before starting the loop
- [ ] If the file exists but is not readable (permission denied), Ralph exits with a clear error message
- [ ] When reading from stdin, Ralph exits with an error if stdin is empty (zero bytes)
- [ ] All checks happen at startup, not deferred to the first iteration
