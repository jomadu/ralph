# R003: Doc Accuracy

**Outcome:** O007 — User Documentation

## Requirement

Documentation matches actual product behavior.

## Detail

Described commands, flags, config keys, exit codes, and report formats behave as documented. When implementation changes, user-facing docs and release notes are updated in the same change or follow-up so drift is not left unresolved. Engineering specs are the source of truth for protocols; user docs reflect them accurately.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Pre-release behavior differs from stable | Docs label pre-release or version scope |
| Bug fix changes behavior | Docs updated if the described contract changes |

### Examples

**Input:** User follows README install steps on a supported platform.

**Expected output:** Installed binary matches documented invocation and version behavior.

**Verification:** Spot-check critical paths (install, run, review) against docs after releases.
