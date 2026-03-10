# O006: Install and Uninstall

## Who

Users who want to add Ralph to their system (e.g. for use from a shell) and remove it cleanly when no longer needed.

## Statement

Users can install Ralph on their system and uninstall it cleanly.

## Why it matters

Without a defined install and uninstall story, users must discover how to obtain the binary, where to put it, and how to make Ralph invocable from the shell. They may copy the binary by hand, guess at PATH, or rely on ad-hoc scripts that differ by platform. Uninstall is unclear — removing one file may leave config, symlinks, or PATH entries pointing at nothing. A clear install/uninstall outcome makes Ralph a normal system tool: predictable to add and remove. Keeping the process simple (e.g. release-only install, no build-from-source) reduces failure modes and keeps behavior consistent. Upgrade and update are covered by the update upgrade outcome.

## Verification

- User follows documented install steps (including optional version choice). After install, they can invoke Ralph from a new shell and confirm it runs successfully (e.g. a version check or a simple command).
- User can run Ralph's commands from the shell without specifying a path to the binary; the installed artifact is on PATH or otherwise invocable as documented.
- User follows documented uninstall steps. After uninstall, Ralph is no longer invocable from the shell (or the install location no longer contains the binary). No broken PATH entries or orphaned files that were part of the documented install scope remain.

## Non-outcomes

- Ralph does not require a built-in installer subcommand (e.g. a "install" command inside Ralph). Installation is supported by documented install and uninstall steps (e.g. scripts) that operate on release artifacts.
- Installation does not build from source or accept a local binary path; it uses a pre-built binary from a release (e.g. latest or a specified version).
- Ralph does not manage or install the AI CLI backends (Claude, Kiro, etc.). Those remain the user's responsibility.
- Ralph does not define a single mandatory install location. The documented install method defines where the binary is placed; the outcome is that the supported method works and is documented.
- Upgrade and update are out of scope for this outcome; they are covered by the update upgrade outcome.

## Risks

| Risk | Mitigating Requirement |
|------|------------------------|
| User has no clear way to install Ralph | [R001 — Documented install](R001-documented-install.md) |
| Binary not on PATH or not invocable after install | [R002 — Invocable after install](R002-invocable-after-install.md) |
| Uninstall leaves broken PATH or orphaned artifacts | [R003 — Documented uninstall](R003-documented-uninstall.md) |
| User cannot confirm install succeeded | [R002 — Invocable after install](R002-invocable-after-install.md) |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R001](R001-documented-install.md) | Documented install steps result in Ralph being invocable from a shell | draft |
| [R002](R002-invocable-after-install.md) | After install, user can run Ralph commands without specifying the binary path | draft |
| [R003](R003-documented-uninstall.md) | Documented uninstall steps remove Ralph with no broken PATH or orphaned install artifacts | draft |
| [R004](R004-clean-uninstall.md) | After uninstall, Ralph is not invocable and no leftover binary remains in the documented install scope | draft |
