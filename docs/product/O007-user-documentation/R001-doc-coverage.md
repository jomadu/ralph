# R001: Doc coverage

**Outcome:** O007 — User Documentation

## Requirement

Documentation covers install, configuration, run, review, update, upgrade, uninstall, and exit codes.

## Detail

User-facing documentation exists for each of the major usage areas so that a user can find how to get Ralph onto their system, set it up, run the loop and review, change versions, remove it, and interpret how the program exited. Coverage is by topic; the doc set need not be a single file as long as each topic is addressed and reachable (see R002 for discoverability).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| A topic (e.g. upgrade) is not yet implemented in the product | Docs may state that the feature is planned or omit it; they must not describe behavior that does not exist. |
| Multiple ways to install (e.g. brew, binary) | Documentation covers at least one supported install path; additional paths may be documented. |
| Config and run behavior vary by version | Docs may reference "current" or versioned behavior; coverage of configuration options and run/review behavior is required for the documented version(s). |
| Exit codes or report format change | Docs describe the actual exit codes and report format; when they change, docs are updated (see R003). |

### Examples

#### User looks up how to configure the loop

**Input:** User has Ralph installed and wants to set iteration limit or timeout.

**Expected output:** User finds user documentation that describes configuration (e.g. where to configure and which options control iteration limits, timeouts, signals). The description is sufficient to set the value and understand its effect.

**Verification:** User can locate the configuration topic and apply a documented option; behavior matches the description.

#### User looks up how to interpret exit status

**Input:** User ran the run command and got exit code 1; they want to know what it means.

**Expected output:** User documentation describes exit codes (e.g. 0 = success, non-zero = failure or threshold) and how they relate to loop outcome, so the user can interpret the result.

**Verification:** User can find the documented exit code meanings and understand why the run ended.

## Acceptance criteria

- [ ] Documentation includes install: at least one way to get Ralph onto the user's system.
- [ ] Documentation includes configuration: how to set up loop-related options (e.g. iteration limits, timeouts, signals) via the documented configuration mechanism.
- [ ] Documentation includes run and review: how to run the loop and how to run review using the documented commands and options.
- [ ] Documentation includes update and upgrade: how to update or upgrade to a chosen version (or states limitation if not supported).
- [ ] Documentation includes uninstall: how to remove Ralph cleanly (or states limitation if not documented).
- [ ] Documentation includes exit codes: what exit codes mean and how they relate to loop/review outcome.
- [ ] Each of the above topics is findable (via structure or links; see R002) and described in enough detail to use the product.

## Dependencies

None.
