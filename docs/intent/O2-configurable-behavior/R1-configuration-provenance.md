# R1: Configuration Provenance Tracking

**Outcome:** O2 — Configurable Behavior

## Requirement

The system tracks the source layer of each resolved configuration value so the user can determine where the active value came from. Every resolved value knows whether it was set by a built-in default, global config file, workspace config file, environment variable, or CLI flag.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] Each resolved config value carries its provenance: built-in default, global config, workspace config, environment variable, or CLI flag
- [ ] Provenance information is available via debug-level logging
- [ ] When multiple layers set the same key, the highest-precedence layer wins and the provenance reflects that winning layer
- [ ] Precedence order is: CLI flags > environment variables > workspace config > global config > built-in defaults
