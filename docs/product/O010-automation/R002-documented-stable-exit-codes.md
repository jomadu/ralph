# R002: Documented stable exit codes

**Outcome:** O010 — Automation

## Requirement

Ralph documents and maintains stable exit codes for success, failure (e.g. threshold, prompt errors), exhaustion or review/apply failure, and interruption, consistent within the compatibility contract.

## Detail

Scripts and CI need to branch on outcome. Ralph uses a documented set of exit codes (or code ranges) so that automation can distinguish: **success** (e.g. loop succeeded or review completed with no prompt errors); **failure** (e.g. loop exited due to failure threshold, or prompt had errors, or review/apply did not complete); **exhaustion** (e.g. max iterations reached without success); **interruption** (e.g. user or system interrupted the process). The exact numeric values and their semantics are documented (e.g. in user docs). Within the compatibility contract (e.g. major version), these codes and their meanings are stable so that scripts do not break when upgrading. New outcomes may be added in a backward-compatible way (e.g. new codes); existing codes are not repurposed.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Loop succeeds (success signal) | Exit code = the documented success code. |
| Loop exits on failure threshold | Exit code = the documented failure (or failure-threshold) code. |
| Loop exits on max iterations without success | Exit code = the documented exhaustion code; distinct from success and failure threshold. |
| Review completes, no prompt errors | Exit code = the documented success code; per O005/R008. |
| Review completes, prompt has errors | Exit code = the documented prompt-errors (or failure) code. |
| Review or apply did not complete (e.g. I/O, missing command) | Exit code = the documented failure code. |
| Process interrupted (SIGINT, etc.) | Exit code = the documented interruption code. |
| Compatibility contract (e.g. major version) | Existing documented codes do not change meaning; new codes may be added. |

### Examples

#### Script branches on loop outcome

**Input:** Script runs the run command; loop exits because the failure threshold was reached.

**Expected output:** Process exits with the documented failure (or failure-threshold) code. Documentation states that this code means "exited due to failure threshold" or equivalent.

**Verification:** Script can branch on the documented code and rely on the meaning; compatibility contract (e.g. major version) does not change the code for this outcome.

#### CI gates on review result

**Input:** CI runs the review command; review completes but prompt has errors.

**Expected output:** Process exits with the documented prompt-errors (or failure) code. Documentation states that this code means "review completed but prompt had errors."

**Verification:** CI can fail the job when exit code is the documented prompt-errors code; meaning is stable within compatibility contract.

#### Interruption

**Input:** User or script sends SIGINT to Ralph during the loop.

**Expected output:** Process exits with the documented interruption code. Documentation states that this code means "interrupted."

**Verification:** Script can distinguish interruption from success/failure/exhaustion; code is stable within compatibility contract.

## Acceptance criteria

- [ ] Exit codes for at least the following outcomes are documented: success, failure (including failure threshold and prompt errors), exhaustion (e.g. max iterations), and interruption.
- [ ] The documentation is sufficient for script and CI authors to branch on outcome (e.g. success vs failure vs exhaustion vs interruption).
- [ ] Within the compatibility contract (e.g. same major version), documented exit codes do not change meaning; existing codes are not repurposed.
- [ ] Review command exit semantics align with the prompt review outcome (O005) and are included in the documented set.

## Dependencies

- O001, O004 — Loop exit semantics (success, failure threshold, exhaustion) are defined there; this requirement ensures they are documented and stable for automation.
- O005 — Review exit semantics (success, prompt errors, failure) are defined there; this requirement ensures they are documented and stable for automation.
