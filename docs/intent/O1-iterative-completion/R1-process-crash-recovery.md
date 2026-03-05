# R1: Process Crash Recovery

**Outcome:** O1 — Iterative Completion

## Requirement

The system continues loop execution when an AI CLI process crashes or exits non-zero, preserving any output produced before the crash for signal scanning.

A crash does not receive special treatment beyond what the loop already provides. The partial output is scanned for signals using the same rules as a normal exit. If no signal is found, the iteration counts as a no-signal iteration — the consecutive failure counter is reset to zero, and the loop proceeds to the next iteration. A crash counts as one completed iteration toward the max iteration limit.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] When the AI CLI process exits with a non-zero exit code, Ralph captures all output written to stdout/stderr before the exit
- [ ] Captured partial output is scanned for success and failure signals using the same logic as a normal exit
- [ ] A crash with no signal in the partial output resets the consecutive failure counter to zero
- [ ] A crash counts as one completed iteration toward the max iteration limit
- [ ] The loop proceeds to the next iteration after a crash (unless max iterations or failure threshold is reached)
- [ ] Ralph does not retry the AI CLI process within the same iteration
