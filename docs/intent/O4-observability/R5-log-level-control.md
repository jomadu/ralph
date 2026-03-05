# R5: Log Level Control

**Outcome:** O4 — Observability

## Requirement

The system supports configurable log verbosity levels, allowing the user to control how much operational output Ralph produces. Log levels govern Ralph's own operational messages (iteration progress, provenance, warnings, errors) and do not control AI output streaming, which is managed independently by the --verbose flag (see O4/R3).

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] Supported log levels: debug, info, warn, error (in order of decreasing verbosity)
- [ ] Default log level is info
- [ ] --quiet sets the effective log level to error, suppressing all non-error output from Ralph
- [ ] --verbose sets the effective log level to debug (in addition to enabling AI output streaming per O4/R3)
- [ ] --log-level explicitly sets the log level and takes precedence over --quiet and --verbose for log verbosity, but does not affect AI output streaming
- [ ] All log output goes to stderr
