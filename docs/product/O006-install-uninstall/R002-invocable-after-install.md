# R002: Invocable after install

**Outcome:** O006 — Install and Uninstall

## Requirement

After install, the user can run Ralph commands from the shell without specifying the path to the binary.

## Detail

The product is on the user's PATH or documented location (or otherwise invocable as documented). The user can open a new shell and invoke the product by the documented command name. This allows the user to confirm that install succeeded and to use Ralph as a normal system tool; the user does not need to type the full path for normal use.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User opens a new shell after install | The product (or documented command) resolves and runs. |
| User has multiple shells (e.g. login vs non-login) | Invocability holds in shells that load the same environment (e.g. PATH or documented env) as described in install docs. |
| Install location is not on default PATH | Documentation explains how to make the product invocable (e.g. add to PATH, profile snippet) so that the user can run it. |
| Command name differs from binary name (e.g. symlink) | Documentation states the invocable command name; that name works after install. |

### Examples

#### Version check in new shell

**Input:** User completed documented install. User opens a new terminal and invokes the product (e.g. a version check).

**Expected output:** The product runs and prints version information (or equivalent success output); exit code 0.

**Verification:** The product is found and runs; output indicates Ralph is running. User has confirmed install succeeded.

#### Run a Ralph subcommand without path

**Input:** User runs the product with a subcommand (e.g. run) from a directory where they use Ralph.

**Expected output:** The product executes the subcommand; user did not need to specify the full path to the binary.

**Verification:** Command succeeds (or fails for reasons unrelated to install, e.g. missing config); invocability is established.

## Acceptance criteria

- [ ] After completing the documented install steps, the user can invoke the product from a new shell using the documented command name without specifying the binary path.
- [ ] A simple invocation (e.g. a version check or equivalent) allows the user to confirm that install succeeded.
- [ ] Invocability is achieved via the user's PATH or the mechanism described in the install documentation.
- [ ] Behavior holds in any shell that has been configured per the install docs (e.g. profile sourced, environment set as documented).

## Dependencies

- O006/R001 — Documented install steps must exist; invocability is the result of following them.
