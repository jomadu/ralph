# R002: Doc Discoverability

**Outcome:** O007 — User Documentation

## Requirement

The user can look up CLI commands, flags, config keys, and environment variables in one place or a clearly linked set of documentation.

## Detail

Documentation is organized so a user does not have to hunt across unrelated pages. The README, engineering component docs, and cross-links form a coherent map: primary user surface (README) links to canonical specs (CLI, config, exit codes) where detail lives. List/show/help commands align with documented names so what the user reads matches what they type.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Option exists only in engineering spec | User-facing doc links to or summarizes it |
| Deprecated flag | Documented as deprecated with replacement |

### Examples

**Input:** User wants all `ralph run` flags.

**Expected output:** README or linked doc section lists flags or points to the CLI spec.

**Verification:** Every flag in `ralph run --help` appears in user-facing docs or the linked CLI spec.
