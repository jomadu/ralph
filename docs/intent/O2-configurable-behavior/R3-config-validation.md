# R3: Configuration Validation at Load Time

**Outcome:** O2 — Configurable Behavior

## Requirement

The system validates all resolved configuration values at load time and fails fast with clear, actionable error messages for invalid values. Validation runs after all layers are merged, so the user sees errors for the final resolved values, not intermediate ones.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] default_max_iterations less than 1 produces a validation error
- [ ] failure_threshold less than 1 produces a validation error
- [ ] Empty success or failure signal strings produce a validation error
- [ ] Invalid log_level values (not one of debug, info, warn, error) produce a validation error
- [ ] Negative iteration_timeout produces a validation error
- [ ] max_output_buffer less than 1 produces a validation error
- [ ] iteration_mode not one of "max-iterations" or "unlimited" produces a validation error
- [ ] A prompt alias with an empty path produces a validation error
- [ ] Validation errors prevent the loop from starting
- [ ] Error messages identify the invalid value, the field name, and the config source (provenance)
