# R5: Log Level Control

**Outcome:** O4 — Observability

## Requirement

The system supports configurable log verbosity levels, allowing the user to control how much operational output Ralph produces. Log output goes to stderr to keep stdout clean for prompt output (dry-run) and AI output piping.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] Supported log levels: debug, info, warn, error (in order of decreasing verbosity)
- [ ] Default log level is info
- [ ] --quiet sets the effective log level to error, suppressing all non-error output
- [ ] --verbose sets the effective log level to debug
- [ ] --log-level explicitly sets the log level and takes precedence over --quiet and --verbose
- [ ] All log output goes to stderr
