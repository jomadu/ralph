# R2: Signal Precedence Rules

**Outcome:** O1 — Iterative Completion

## Requirement

The system resolves conflicting signals in a single iteration's output using a deterministic precedence rule: failure wins over success. Signal scanning uses substring matching against the full captured output after the AI CLI process exits.

## Specification

_To be specified._

### Edge cases

_To be specified._

### Examples

_To be specified._

## Acceptance criteria

- [ ] When only the success signal appears in output, the iteration is treated as success and the loop stops
- [ ] When only the failure signal appears in output, the iteration is treated as failure and the consecutive failure counter is incremented
- [ ] When both success and failure signals appear in the same output, the iteration is treated as failure
- [ ] When neither signal appears in output, the iteration is treated as no-signal — the consecutive failure counter is reset to zero and the loop proceeds to the next iteration
- [ ] Signal scanning uses substring matching — the signal string can appear anywhere in the output
- [ ] Signal scanning occurs after the AI CLI process exits, not during streaming

## Dependencies

_None identified._
