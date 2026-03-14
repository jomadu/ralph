# R004: Clean uninstall

**Outcome:** O006 — Install and Uninstall

## Requirement

After uninstall, the product is not invocable from the shell and no leftover artifacts remain in the documented install scope.

## Detail

The observable result of a successful uninstall is: (1) the user cannot invoke the product from the shell (e.g. the product is not found or no longer resolves), and (2) no leftover artifacts remain in the documented install location(s). This is the verifiable end state that the documented uninstall steps (R003) are intended to achieve. It does not require removal of user-created config or data outside the documented install scope.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User ran documented uninstall | The product is not invocable; no leftover artifacts in the documented install location(s). |
| Multiple copies of binary (e.g. user copied elsewhere) | Only the documented install scope is guaranteed clean; copies outside that scope are out of scope. |
| Environment contained multiple entries (e.g. install dir + manual path) | After uninstall, at least the install-related entry is gone so the product is not invocable via the documented mechanism; no broken references. |
| Shell still has the product in memory (e.g. hash table) | Invocability is defined for a new shell; after uninstall, a new shell must not find the product. |

### Examples

#### Not invocable after uninstall

**Input:** User completed documented uninstall. User opens a new shell and attempts to invoke the product (e.g. a version check).

**Expected output:** Shell reports command not found (or equivalent); the product does not run.

**Verification:** Exit code or shell message indicates the product is not available. User has confirmed Ralph is no longer on the system for normal use.

#### No leftover binary in install scope

**Input:** Documented install placed the product at the documented install location. User completed documented uninstall.

**Expected output:** The documented install location no longer contains the product (no leftover binary or artifacts).

**Verification:** Checking the documented install location shows no leftover binary or artifacts in documented scope.

## Acceptance criteria

- [ ] After the user completes the documented uninstall steps, the product cannot be invoked from a new shell (e.g. command not found).
- [ ] After uninstall, the documented install location(s) do not contain the product (no leftover artifacts in documented install scope).
- [ ] Verification can be done by opening a new shell and attempting to invoke the product, and by checking that the documented install location no longer has the product.
- [ ] Scope is limited to what the documented install and uninstall define; user-owned files outside that scope are not required to be removed.

## Dependencies

- O006/R003 — Documented uninstall steps must exist and be followed to achieve this state; R004 describes the resulting state.
