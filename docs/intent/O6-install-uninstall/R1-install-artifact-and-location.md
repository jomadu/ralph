# R1: Install Artifact and Location Definition

**Outcome:** O6 — Install and Uninstall

## Requirement

The system (or distribution) defines what is installed when Ralph is installed and where those artifacts live, so that install and uninstall have a clear scope and users know what they are adding or removing.

## Specification

Ralph's install surface is defined so that the supported install method (install script, R2) has a well-defined artifact set and location. This requirement does not mandate a single global location; it mandates that the scope is documented and consistent.

**Artifact set:**

- **Primary artifact:** A single executable binary named `ralph` (or platform-equivalent, e.g. `ralph.exe` on Windows). This binary is self-contained for normal operation (run, review, list, version, etc.). Optional assets (e.g. shell completions, man pages) may be included in the artifact set for a given method if documented.
- **User-owned config:** Ralph does not install user config. User config lives in workspace or global config paths (e.g. `./ralph-config.yml`, `~/.config/ralph/ralph-config.yml`) and is created or edited by the user. These are not part of the "installed" set for uninstall (R4).

**Location:**

- For the supported install method (install script), documentation (R5) must state where the binary (and any optional artifacts) are placed. Example: a directory on PATH chosen by the user or by the script (e.g. `~/bin`, `/usr/local/bin`). The location must be defined so that uninstall knows what to remove.

**Scope for uninstall:**

- Uninstall (R4) removes only what the install script placed. The definition of "what was installed" is the artifact set and location described in the documentation. User-created config files are not removed by uninstall unless the docs explicitly say so (e.g. "removes binary only").

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User runs install script with different target dir on a second run | Script uses a defined or user-specified location; uninstall removes what the script recorded or uses the same logic to find and remove the binary. |
| User moves the binary after install | Invokability (R3) may break if PATH or symlinks pointed at the old path. Uninstall as documented applies to the original install location; moved binary is out of scope. |
| Optional assets (completions, man pages) added by the script | Documented as part of the artifact set; uninstall script removes them. |

### Examples

#### Documented artifact set for "manual binary" method

**Input:** Documentation for "Install by downloading the release binary."

**Expected content (conceptual):** "Place the `ralph` binary in a directory on your PATH (e.g. `~/bin` or `/usr/local/bin`). Only the binary is installed. Uninstall by deleting that binary."

**Verification:** User knows exactly one file was added and where; uninstall is delete-one-file.

#### Documented artifact set for install script

**Input:** Documentation for "Install via the install script."

**Expected content (conceptual):** "The install script places the `ralph` binary in a defined location (e.g. `~/bin` or a directory the user chooses). Only the binary is installed. Uninstall by running the uninstall script, which removes that binary (and any references the install script added)."

**Verification:** User knows install location and that uninstall is via the uninstall script.

## Acceptance criteria

- [ ] What is installed (binary and any optional artifacts) is documented
- [ ] Where artifacts are placed (or how the user chooses the location) is documented
- [ ] User config (e.g. `ralph-config.yml`) is not part of the installed artifact set for uninstall purposes unless the docs explicitly say so

## Dependencies

_None. R2 (supported install methods) and R5 (documentation) reference this definition; this requirement defines the scope that those refer to._
