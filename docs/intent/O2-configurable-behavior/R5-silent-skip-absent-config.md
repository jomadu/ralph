# R5: Silent Skip for Absent Config Files

**Outcome:** O2 — Configurable Behavior

## Requirement

The system silently skips config files that do not exist at their default locations, applying only the layers that are present. This allows Ralph to run with no config files at all (using built-in defaults) and avoids requiring users to create config files before first use.

## Specification

_To be specified._

### Edge cases

_To be specified._

### Examples

_To be specified._

## Acceptance criteria

- [ ] A missing global config file (~/.config/ralph/ralph-config.yml) does not produce an error or warning
- [ ] A missing workspace config file (./ralph-config.yml) does not produce an error or warning
- [ ] When no config files are present, built-in defaults are used for all values
- [ ] When --config specifies an explicit path that does not exist, Ralph exits with an error (the user explicitly named a file, so its absence is an error)

## Dependencies

_None identified._
