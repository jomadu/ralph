# R3: Verbose Output Streaming

**Outcome:** O4 — Observability

## Requirement

The system streams the AI CLI's output to the terminal in real time when verbose mode is enabled, while still capturing it in the output buffer for signal scanning. This lets the user watch the AI work without sacrificing loop control.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] With --verbose or -v, AI CLI stdout and stderr are mirrored to the terminal as they are produced
- [ ] Output is simultaneously captured in the buffer for signal scanning after the process exits
- [ ] Without --verbose, AI CLI output is captured silently and not displayed to the terminal
- [ ] Verbose mode also enables debug-level log messages from Ralph itself
