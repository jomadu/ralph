# R5: Install and Uninstall Documentation

**Outcome:** O6 — Install and Uninstall

## Requirement

A single authoritative reference exists for how to install and how to uninstall Ralph (using the install and uninstall scripts), so that users and automation can follow consistent steps.

## Specification

The project provides documentation that serves as the authoritative reference for install and uninstall. This can be a section of the main README, a dedicated INSTALL or docs page, or equivalent. The documentation is the source of truth for R1 (artifact and location), R2 (install/uninstall scripts), R3 (invokability), and R4 (uninstall behavior).

**Content:**

- **Install:** The docs describe (for the install script, R2):
  - Prerequisites: `curl` (required to download the release binary). The script is run from the repo or via curl; no Go or build step.
  - That the install script **only** installs from release artifacts (downloads a pre-built binary from GitHub releases; no build from source).
  - Optional **version** positional argument: e.g. `install.sh 1.0.0` or `install.sh v1.0.0` to install a specific version; omit for latest release.
  - Steps to run the install script (e.g. from repo or via curl), and what it does (downloads binary for OS/arch, places it, records state for uninstall).
  - Where the binary is placed (R1) and how to override (e.g. `--dir` or env).
  - How to ensure Ralph is invokable (R3): e.g. "place in a directory on PATH," "open a new terminal," or "add this path to PATH."
  - How to verify: e.g. run `ralph version` and expect exit 0 and version output.
- **Uninstall:** The docs describe:
  - Steps to remove Ralph (run the uninstall script).
  - What is removed (binary and install state per R1); what is not removed (e.g. user config unless stated otherwise).
  - That no broken references (e.g. PATH pointing at removed binary) should remain (R4).
- **Platform support:** Documentation states which platforms (OS/arch) the install script supports: Linux (amd64, arm64), macOS (amd64, arm64), Windows (amd64, arm64) where the script can run (e.g. Git Bash). Gaps can be stated if any.

**Location:**

- The authoritative reference is in the repository (e.g. README.md, docs/INSTALL.md) or in published docs linked from the repository, so that contributors and users have one place to look.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Install script behavior or options change (e.g. new default path) | Documentation is updated to reflect new steps, artifact set (R1), and uninstall behavior. |
| A method is deprecated | Documentation updated to mark deprecated or remove; existing users can still find uninstall steps (e.g. in changelog or archived doc). |
| Platform-specific quirks (e.g. Windows PATH) | Documented for the platforms the project supports so users can complete install and uninstall. |

### Examples

#### README install section

**Input:** README has an "Install" section.

**Expected:** Install section describes how to run the install script, where the binary is placed, how to ensure it's on PATH (if needed), and how to verify (`ralph version`). Uninstall section describes how to run the uninstall script and what is left (e.g. "User config is not removed").

**Verification:** A new user on a documented platform can install and uninstall using only the README.

#### Dedicated INSTALL.md

**Input:** docs/INSTALL.md with "Install" and "Uninstall" sections, and README links to it ("See [Install guide](docs/INSTALL.md).").

**Expected:** INSTALL.md contains the same content as above; README remains the entry point that points to the authoritative guide.

**Verification:** Single authoritative place is INSTALL.md; README directs users there.

## Acceptance criteria

- [ ] There is a single authoritative place (or linked set) that describes install and uninstall
- [ ] Install steps, artifact/location (R1), invokability (R3), and uninstall steps (R4) are documented
- [ ] Documentation states or links to which platforms (OS/arch) the install script supports
- [ ] Verification step (e.g. `ralph version`) is documented for install

## Dependencies

- R1, R2, R3, R4 — This requirement is the documentation layer that describes the behavior specified in those requirements; it does not define new behavior, it captures it in one place.
