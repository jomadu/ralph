# R2: Unknown Key Warnings

**Outcome:** O2 — Configurable Behavior

## Requirement

The system warns the user when a config file contains keys that are not part of the known schema, without preventing startup. This catches typos and forward-compatibility issues while remaining non-blocking.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] Unrecognized top-level keys in a config file produce a warning at load time
- [ ] Unrecognized nested keys within known sections also produce warnings
- [ ] Each warning identifies the key name and the config file path where it was found
- [ ] The presence of unknown keys does not prevent Ralph from starting or running
