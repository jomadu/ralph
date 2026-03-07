# R2: Supported Install Methods

**Outcome:** O6 — Install and Uninstall

## Requirement

At least one supported, documented method exists for getting Ralph onto the user's system so that they can install without reverse-engineering the project.

## Specification

The project supports one or more install methods. Each method is documented (R5) and has a defined artifact set and location (R1). Support means: the method is tested or maintained so that users following the docs can complete install and achieve post-install invokability (R3).

**Minimum:**

- At least one method must be supported. For this project, the supported method is **install and uninstall scripts**: the user runs an install script (e.g. `scripts/install.sh`) that **obtains the Ralph binary only from release artifacts** (GitHub releases) and places it in a defined location (and ensures invokability, e.g. on PATH). The user runs an uninstall script (e.g. `scripts/uninstall.sh`) to remove the installed artifact(s). No package manager (e.g. Homebrew) is required or documented.

**Install script behavior:**

- The install script **only** installs from release artifacts. It does not build from source and does not accept a path to a local binary. It downloads a pre-built binary for the user's OS and architecture from a GitHub release.
- The script accepts an optional **version** positional argument (e.g. `1.0.0` or `v1.0.0`). If omitted, the script installs the latest release. This allows users to pin a specific version.
- The script requires `curl` (or equivalent) to download the release asset. Prerequisites are documented (R5).

**Platforms:**

- Documentation must state which platforms (OS/arch) the install script supports. Release artifacts are built for Linux (amd64, arm64), macOS (amd64, arm64), and Windows (amd64, arm64); the install script supports those same platforms where the script can run (e.g. Linux, macOS, and Windows environments such as Git Bash). Gaps are acceptable if stated.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User's platform has no release artifact (e.g. unsupported OS/arch) | Script exits with a clear error; documentation states supported platforms. User may obtain a binary by other means (e.g. build from source manually) but that is outside the supported install method. |
| User requests a version that does not exist or has no asset for their platform | Script exits with a clear error (e.g. download failed, release not found). |
| Install or uninstall script behavior changes (e.g. new default path) | Documentation updated; existing users can still run uninstall script for the location it previously used or recorded. |
| User runs install script twice (e.g. different target dir or different version) | Install location (R1) is defined by the script; uninstall removes what the script recorded (one install location per run). Installing a different version overwrites the binary at that location. |

### Examples

#### Install and uninstall scripts

**Input:** Project provides `scripts/install.sh` and `scripts/uninstall.sh`.

**Expected:** Docs describe: how to run the install script (e.g. from repo or via curl), optional version argument (install latest vs. a specific version), that the script downloads the binary from GitHub releases only (no build from source), what it installs (binary at a defined or user-specified location), and how to verify (`ralph version`). R1 artifact set: single binary (and any state the script writes for uninstall). Uninstall: run `scripts/uninstall.sh`, which removes the binary and the install state.

**Verification:** User on a documented platform with network access can run the install script (optionally with a version), then run `ralph version`; can uninstall by running the uninstall script.

## Acceptance criteria

- [ ] At least one install method is supported and documented
- [ ] The install and uninstall scripts are consistent with R1 (artifact and location defined) and lead to R3 (invokability)
- [ ] Supported platforms (OS/arch) for the install script are documented or clearly implied

## Dependencies

- R1 — Install artifact and location definition (install script has a defined scope)
- R5 — Install and uninstall documentation (methods are documented there)
