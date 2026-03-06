# R4: Uninstall Behavior

**Outcome:** O6 — Install and Uninstall

## Requirement

The user can remove Ralph from the system in a way that removes the installed artifact(s) for the install method they used and does not leave broken references (e.g. PATH entries or config that point at a missing binary).

## Specification

The project provides an uninstall procedure (the uninstall script, per R2). Executing it removes what the install script installed (as defined in R1) and does not leave the system in a broken state.

**What uninstall removes:**

- The binary (and any optional artifacts) that the install method placed, at the location(s) defined for that method (R1). User-created config files (e.g. `~/.config/ralph/ralph-config.yml`, `./ralph-config.yml`) are not removed unless the method explicitly documents that (e.g. "uninstall also removes global config"). Default expectation: uninstall removes only the artifacts installed by that method; user config is left in place.

**Broken references:**

- Uninstall must not leave references that point at the removed binary. Examples of broken references: a PATH entry that was the only content of a directory that is now empty; a symlink in a PATH directory that points to the removed binary. Resolution: either the uninstall procedure removes such references (e.g. remove the symlink, or remove the empty directory from PATH if the install added it), or the install method does not create them (e.g. user manually added an existing directory to PATH and placed the binary there — uninstall is "delete the binary," and the PATH entry still points at a valid directory). The requirement is that after uninstall, the user is not left with a PATH entry or similar that points at nothing or at a removed binary.

**How uninstall is performed:**

- The project provides an uninstall script (e.g. `scripts/uninstall.sh`) that removes the binary and any references the install script added (e.g. symlinks). The procedure is documented (R5). No requirement for a separate `ralph uninstall` subcommand; the script is sufficient.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User installed to a directory shared with other binaries (e.g. `/usr/local/bin`) | Uninstall removes only the `ralph` binary from that directory; other files untouched. PATH remains valid. |
| User added a new directory to PATH only for Ralph and put only `ralph` there | Uninstall docs: delete the binary. Optionally document "you may remove this directory from PATH." Removing the directory from PATH is user's choice; deleting the binary is required for uninstall. |
| Uninstall script run after manual delete of binary | Script detects missing binary or idempotently removes what it can; no broken state if script is idempotent or reports "not installed." |
| User deleted the binary manually before "uninstall" | No procedure to run; no broken state if no PATH entry pointed only at that file. Docs may say "if you already deleted the binary, uninstall is complete." |
| User installed to a custom path and runs uninstall script | Uninstall script removes the binary from the location the install script used (or from a recorded path); other copies elsewhere are out of scope. |

### Examples

#### Uninstall manual binary install

**Input:** User installed by copying `ralph` to `~/bin`. Uninstall docs: "Delete the `ralph` binary from the directory where you installed it (e.g. `rm ~/bin/ralph`)."

**Expected:** User deletes the file. `ralph` is no longer in that directory. If `~/bin` is on PATH and contained only `ralph`, the directory may be empty; PATH still points at a valid directory. No broken reference.

**Verification:** `which ralph` (or equivalent) returns nothing; `ralph` command not found.

#### Uninstall via uninstall script

**Input:** User ran the install script, which placed `ralph` in `~/bin`. Uninstall docs: "Run `scripts/uninstall.sh`."

**Expected:** Uninstall script removes the binary from `~/bin` (and any symlinks or PATH-related changes the install script made). No broken references.

**Verification:** `ralph` not found; `~/bin/ralph` does not exist.

## Acceptance criteria

- [ ] The uninstall procedure (script) is documented
- [ ] Executing the uninstall script removes the artifacts the install script installed (per R1)
- [ ] After uninstall, no broken references remain (e.g. PATH or symlinks pointing at the removed binary)
- [ ] User config files are not removed by uninstall unless the method explicitly documents that

## Dependencies

- R1 — Install artifact and location definition (defines what "installed" means)
- R2 — Supported install methods (install and uninstall scripts)
- R5 — Install and uninstall documentation (uninstall procedure is documented there)
