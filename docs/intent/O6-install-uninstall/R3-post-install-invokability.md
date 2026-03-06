# R3: Post-Install Invokability

**Outcome:** O6 — Install and Uninstall

## Requirement

After a successful install, the user can invoke Ralph from a shell (e.g. run `ralph` or `ralph version`) without specifying the full path to the binary, so that Ralph behaves as a normal CLI tool.

## Specification

The documented install steps (running the install script, R2) must result in the user being able to invoke the Ralph binary from a new shell without typing its full path. "Invoke" means the shell finds the binary by name (e.g. `ralph`); typically this is achieved by placing the binary in a directory that is on the user's PATH, or by documenting a single path the user can add to PATH.

**Typical mechanisms:**

- **Binary on PATH:** The install script places the binary in a directory that is already on PATH (e.g. `/usr/local/bin`, `~/bin`) or prompts the user for a directory and ensures it is on PATH (or documents adding it). After install, opening a new shell and running `ralph` runs the installed binary.
- **Explicit path:** If the script installs to a path not on PATH, the documentation must tell the user how to add that path to PATH or how to invoke by full path (e.g. alias or symlink from a PATH directory). The outcome is still "user can run `ralph` (or equivalent) from a shell."

**Verification by user:**

- After install, the user opens a new shell (or sources their profile if the install modified it) and runs `ralph version`. Exit 0 and version output confirm the correct binary is on PATH (or otherwise invokable).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User installs but does not open a new shell | Current shell may not see updated PATH; docs should mention "open a new terminal" or "reload your shell config." |
| Multiple Ralph binaries on PATH (e.g. old install + new) | First one in PATH wins. Uninstall (R4) of the intended copy removes that copy; user may need to fix PATH order or remove the other copy. |
| PATH not set as expected (e.g. minimal env) | Documented install steps assume a typical interactive shell PATH; edge cases (cron, CI) may require explicit path or PATH setup in that environment. |

### Examples

#### Install script places binary in ~/bin

**Input:** Install script places `ralph` in `~/bin` (or docs say "Ensure `~/bin` is on your PATH" if the script uses that path).

**Expected:** User runs install script, opens new shell, runs `ralph version`. Command succeeds if the install script put the binary on PATH (e.g. in `~/bin` and `~/bin` is on PATH).

**Verification:** Exit 0; version string on stdout.

#### Install script places binary in /usr/local/bin

**Input:** Install script places `ralph` in `/usr/local/bin` (or another directory that is typically on PATH).

**Expected:** User runs install script (possibly with sudo if required for that directory), opens new shell, runs `ralph version`; succeeds.

**Verification:** Exit 0; version string on stdout.

## Acceptance criteria

- [ ] The documented install steps lead to the user being able to run `ralph` (or the documented command name) from a new shell without the full path
- [ ] Documentation states or implies that a new shell (or reload) may be needed after PATH changes

## Dependencies

- R2 — Supported install methods (invokability is required after any supported install)
