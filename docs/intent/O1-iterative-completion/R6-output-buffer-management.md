# R6: Output Buffer Management

**Outcome:** O1 — Iterative Completion

## Requirement

The system captures AI CLI output into a bounded buffer, discarding the oldest content when the buffer is full to preserve the most recent output for signal scanning. The buffer is per-iteration — each iteration starts with a fresh buffer.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] Output from the AI CLI's stdout and stderr is captured into a single buffer
- [ ] When the buffer exceeds the configured max size, the oldest content is discarded (truncation from the beginning)
- [ ] Signal scanning operates on the buffer contents after the AI CLI process exits
- [ ] The default buffer size is 10MB (10,485,760 bytes)
- [ ] The buffer size is configurable through the standard configuration hierarchy
- [ ] Each iteration starts with a fresh, empty buffer
