# R1: Distinct Exit Codes

**Outcome:** O4 — Observability

## Requirement

The system exits with distinct codes for each termination reason, enabling scripts and CI systems to programmatically determine why the loop stopped without parsing output.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] Exit code 0: success signal received
- [ ] Exit code 1: failure threshold reached or explicit abort
- [ ] Exit code 2: max iterations exhausted without a success signal
- [ ] Exit code 130: interrupted by SIGINT or SIGTERM
- [ ] No two termination reasons share an exit code
- [ ] Exit codes are consistent regardless of verbosity, log level, or other output settings
