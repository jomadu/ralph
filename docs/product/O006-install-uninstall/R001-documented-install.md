# R001: Documented install

**Outcome:** O006 — Install and Uninstall

## Requirement

Documented install steps exist such that following them results in Ralph being invocable from a shell.

## Detail

The product provides clear, step-by-step instructions (or scripts) for installing Ralph. The documented method uses a pre-built binary from a release (e.g. latest or a specified version); it does not build from source or accept an arbitrary local binary path. The steps define where the product is placed and how the user makes it invocable (e.g. on the user's PATH or documented location). After a user completes the documented steps, they can invoke Ralph from a new shell without having to specify a path to the binary (see R002 for the invocability criterion).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User follows steps on a supported platform | Ralph becomes invocable as documented; user can verify (e.g. version check). |
| Documented method specifies a version (e.g. latest vs pinned) | Steps describe how to obtain that version; result is invocable. |
| User has no write access to default install location | Documentation either covers an alternative (e.g. user dir) or states prerequisites (e.g. write access). |
| Multiple supported platforms (e.g. macOS, Linux) | Each has documented steps or a single method that works across them. |

### Examples

#### Install using documented script

**Input:** User runs the documented install script (e.g. curl-pipe or download-then-run) with default options on a supported OS.

**Expected output:** The product is placed in the documented install location; the user's environment (e.g. PATH or documented mechanism) is updated so the product is invocable (e.g. from the shell) by the documented command name.

**Verification:** User opens a new shell and invokes the product (e.g. a version check); the product runs and prints version info.

#### Install a specific release version

**Input:** Documentation describes how to install a specific version (e.g. by URL or version flag). User follows those steps.

**Expected output:** The requested version is installed and invocable.

**Verification:** A version check (or equivalent) reports the installed version; user can invoke the product without specifying the binary path.

## Acceptance criteria

- [ ] Install documentation (or install script with docs) exists and describes steps that use a pre-built release binary (no build-from-source in scope).
- [ ] The steps define or imply where the product is placed and how the user invokes it (e.g. on the user's PATH or documented location).
- [ ] A user who follows the documented steps can subsequently invoke Ralph from a new shell (invocability is verified per R002).
- [ ] Documentation is sufficient for a user on a supported platform to complete install without guessing where to place the product or how to make it invocable.

## Dependencies

None.
