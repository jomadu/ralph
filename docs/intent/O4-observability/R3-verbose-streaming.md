# R3: Verbose Output Streaming

**Outcome:** O4 — Observability

## Requirement

The system streams the AI CLI's output to the terminal in real time when verbose mode is enabled, while still capturing it in the output buffer for signal scanning. This lets the user watch the AI work without sacrificing loop control. AI output streaming is controlled by the --verbose flag and is independent of log level — --log-level affects Ralph's own operational messages but does not suppress or enable AI output streaming.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] With --verbose or -v, AI CLI stdout and stderr are mirrored to the terminal as they are produced
- [ ] Output is simultaneously captured in the buffer for signal scanning after the process exits
- [ ] Without --verbose, AI CLI output is captured silently and not displayed to the terminal
- [ ] AI output streaming is controlled solely by the --verbose flag — --log-level does not affect it (e.g., --verbose --log-level warn streams AI output but suppresses Ralph's debug messages)
- [ ] --log-level debug without --verbose does not enable AI output streaming
