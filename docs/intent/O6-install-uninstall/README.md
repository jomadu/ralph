# O6: Install and Uninstall

## Statement

Users can install Ralph on their system so it is invocable (e.g. from a shell) and uninstall it cleanly, with no broken references or leftover artifacts. Install and uninstall are simple: the install script only installs from release artifacts; the uninstall script removes what was installed.

## Why it matters

Without a defined install and uninstall story, users must discover how to obtain the binary, where to put it, and how to make `ralph` available. They may copy the binary by hand, guess at PATH, or rely on ad-hoc scripts that differ by platform. Uninstall is unclear — removing one file may leave config, symlinks, or PATH entries pointing at nothing. A clear install/uninstall outcome makes Ralph a normal system tool: predictable to add and remove. Keeping the scripts simple (release-only install, no build-from-source) reduces failure modes and keeps behavior consistent.

## Verification

- User follows documented install steps (run the install script, optionally with a version). After install, they open a new shell and run `ralph version`; it prints the version and exits 0.
- User runs `ralph run`, `ralph list prompts`, and other subcommands; all work without specifying a path to the binary.
- User follows documented uninstall steps. After uninstall, `ralph` is no longer found (or the previous install path no longer contains the binary). No broken PATH entries or orphaned files that were part of the documented install scope remain.

## Non-outcomes

- Ralph does not require a system-wide "installer" binary (e.g. `ralph install`) to be built. Installation is supported by install/uninstall scripts that operate only on release artifacts.
- The install script does not build from source or accept a local binary path; it only downloads a pre-built binary from a GitHub release (latest or a specified version).
- Ralph does not manage or install the AI CLI backends (Claude, Kiro, etc.). Those remain the user's responsibility.
- Ralph does not define a single mandatory install location. The install script (and docs) define where the binary is placed; the outcome is that the supported method works and is documented.
- Auto-update or self-update is out of scope. Install/uninstall covers getting Ralph onto the system and removing it, not upgrading in place.

## Risks

| Risk | Mitigating Requirement |
|------|------------------------|
| User does not know what gets installed or where | [R1 — Install artifact and location definition](R1-install-artifact-and-location.md) |
| User cannot obtain Ralph (no clear method) | [R2 — Supported install methods](R2-supported-install-methods.md) |
| After install, `ralph` is not invocable (e.g. not on PATH) | [R3 — Post-install invokability](R3-post-install-invokability.md) |
| Uninstall leaves files or broken PATH/config references | [R4 — Uninstall behavior](R4-uninstall-behavior.md) |
| Install or uninstall steps are missing, wrong, or platform-inconsistent | [R5 — Install and uninstall documentation](R5-install-uninstall-documentation.md) |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R1](R1-install-artifact-and-location.md) | Install artifact and location definition (what is installed and where) | ready |
| [R2](R2-supported-install-methods.md) | Supported install methods (at least one documented way to get Ralph onto the system) | ready |
| [R3](R3-post-install-invokability.md) | Post-install invokability (`ralph` invocable from a shell after install) | ready |
| [R4](R4-uninstall-behavior.md) | Uninstall behavior (removal of installed artifacts, no broken state) | ready |
| [R5](R5-install-uninstall-documentation.md) | Install and uninstall documentation (authoritative reference) | ready |
