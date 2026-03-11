# R002: Install and First Command

**Outcome:** O008 — Discoverability

## Requirement

The user can find install steps and a first command that completes successfully.

## Detail

A new user needs to go from "I don't have the product" to "I ran a command and it completed." Documented install steps (or a documented way to run the product, e.g. via a script or container) are findable. A first command is documented (e.g. the run command or the review command with a path) and is achievable so that when the user follows the steps, that command can complete successfully. Success means the command runs and exits with an outcome the user can recognize (e.g. exit 0, or a documented non-zero that indicates "run completed" rather than "install/config broken").

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User does not have the product installed | Documented install (or run) steps are findable and executable. |
| User follows install then runs first command | The first command can complete successfully (subject to user having minimal prerequisites, e.g. config or prompt as documented). |
| Multiple ways to run (e.g. built binary, script, brew) | At least one path is documented; that path leads to a first command that can succeed. |
| First command requires a config or prompt | Docs state the minimal prerequisite (e.g. minimal config or sample prompt) so the user can satisfy it and run. |

### Examples

#### Install and run first command

**Input:** New user follows documented install steps, then runs the documented first command (e.g. the run command with a minimal config that defines an alias and a prompt).

**Expected output:** The product runs and completes (e.g. loop runs to success/failure or review finishes); exit code is 0 or as documented for "run completed."

**Verification:** User can reproduce: install → run first command → see a successful completion without guessing.

#### Documented run-without-install

**Input:** Docs describe running the product without installing (e.g. a documented script or run-from-source path). User follows that path and runs the documented first command.

**Expected output:** Same as above — the command completes successfully when prerequisites are met.

**Verification:** At least one documented path from "have the product" to "first command completed" is achievable.

## Acceptance criteria

- [ ] Documented install steps (or a documented way to run the product) are findable in docs or README.
- [ ] A first command (e.g. the run command or the review command) is documented and identified as the recommended path to first run.
- [ ] When the user follows the documented install/run path and meets documented minimal prerequisites, that first command can complete successfully (exit 0 or documented "run completed" outcome).
- [ ] Minimal prerequisites for the first command (e.g. config, prompt) are stated so the user can satisfy them.

## Dependencies

None. (Implementation may depend on O002 config and O001 run loop; this requirement is about discoverability of install and first command only.)
