# R2: Supported Install Methods

**Outcome:** O6 — Install and Uninstall

## Requirement

At least one supported, documented method exists for getting Ralph onto the user's system so that they can install without reverse-engineering the project.

## Specification

The project supports one or more install methods. Each method is documented (R5) and has a defined artifact set and location (R1). Support means: the method is tested or maintained so that users following the docs can complete install and achieve post-install invokability (R3).

**Minimum:**

- At least one method must be supported. For this project, the supported method is **install and uninstall scripts**: the user runs an install script (e.g. `scripts/install.sh`) that obtains or uses the Ralph binary and places it in a defined location (and ensures invokability, e.g. on PATH). The user runs an uninstall script (e.g. `scripts/uninstall.sh`) to remove the installed artifact(s). No package manager (e.g. Homebrew) is required or documented.

The install script may assume a pre-built binary is available (e.g. from a release, or built locally) or may build from source; the exact behavior is implementation-defined. The requirement is that install and uninstall scripts exist, are documented, and are consistent with R1 and R3.

**Platforms:**

- Documentation should state which platforms (OS/arch) the install script supports (e.g. macOS arm64/amd64, Linux amd64). Gaps (e.g. "Windows not yet supported") are acceptable if stated.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User's platform has no supported method | Documented; user may build from source if supported for that platform, or use an unsupported workaround at their own risk. |
| Install or uninstall script behavior changes (e.g. new default path) | Documentation updated; existing users can still run uninstall script for the location it previously used or recorded. |
| User runs install script twice (e.g. different target dir) | Install location (R1) is defined by the script; uninstall removes what that script installed (e.g. one target per run). |

### Examples

#### Install and uninstall scripts

**Input:** Project provides `scripts/install.sh` and `scripts/uninstall.sh`.

**Expected:** Docs describe: how to run the install script (e.g. from repo root or via curl), what it installs (binary at a defined or user-specified location), and how to verify (`ralph version`). R1 artifact set: single binary (and any references the script adds). Uninstall: run `scripts/uninstall.sh`, which removes the binary and cleans up what the install script added.

**Verification:** User on a documented platform can run the install script, then run `ralph version`; can uninstall by running the uninstall script.

## Acceptance criteria

- [ ] At least one install method is supported and documented
- [ ] The install and uninstall scripts are consistent with R1 (artifact and location defined) and lead to R3 (invokability)
- [ ] Supported platforms (OS/arch) for the install script are documented or clearly implied

## Dependencies

- R1 — Install artifact and location definition (install script has a defined scope)
- R5 — Install and uninstall documentation (methods are documented there)
