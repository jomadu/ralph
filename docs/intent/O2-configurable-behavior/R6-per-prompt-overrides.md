# R6: Per-Prompt Loop Setting Overrides

**Outcome:** O2 — Configurable Behavior

## Requirement

The system allows each prompt alias to override any loop configuration value independently. Prompt-level overrides take precedence over root-level loop config but are still overridden by environment variables and CLI flags. Each alias's overrides are isolated — one alias's settings do not affect another.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] A prompt alias can override any loop setting: iteration_mode, default_max_iterations, failure_threshold, iteration_timeout, max_output_buffer, ai_cmd, ai_cmd_alias, preamble, signals.success, signals.failure
- [ ] Prompt-level overrides take precedence over root loop config
- [ ] Environment variables and CLI flags still override prompt-level settings
- [ ] Unspecified prompt-level values inherit from the root loop section
- [ ] Each alias's overrides are independent — configuring one alias does not affect any other alias
