# R003: Documented uninstall

**Outcome:** O006 — Install and Uninstall

## Requirement

Documented uninstall steps remove Ralph such that no broken references or leftover artifacts remain within the documented install scope.

## Detail

The product provides clear, step-by-step instructions (or a script) for uninstalling Ralph. The user follows documented uninstall steps. The steps cover removing the product and reversing any changes the documented install made (e.g. environment or path modifications, symlinks, or config files that are part of the install scope). "Documented install scope" means only what the install documentation says it installs or modifies; user-created config or data outside that scope is not required to be removed by uninstall. After following the documented uninstall steps, there are no broken references to the former install location and no leftover artifacts that were part of the documented install (see R004 for the resulting state).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User follows uninstall steps after a normal install | Product removed; environment or invocability changes reverted; no leftover artifacts in documented scope. |
| Install added an environment entry (e.g. in profile) | Uninstall steps remove or revert that entry so no broken reference to the former install location remains. |
| Install created a symlink | Uninstall steps remove the symlink (and optionally the binary if that was the only copy). |
| User has modified config that install did not create | Uninstall does not require removal of user-owned config outside install scope; docs may note what is in scope. |
| Multiple install methods documented | Each has corresponding uninstall steps or a single uninstall procedure that covers the documented install scope. |

### Examples

#### Uninstall after default install

**Input:** User installed Ralph using the documented install method (e.g. script placed the product in the documented install location, profile updated). User follows documented uninstall steps.

**Expected output:** Product removed from the documented install location; profile (or equivalent) no longer adds the install directory to the user's environment, or the directory is removed so no broken reference remains. No leftover artifacts that were part of the documented install remain.

**Verification:** User opens a new shell; the product is not invocable (e.g. command not found). No broken references; no leftover artifacts in the documented install location (per R004).

#### Uninstall removes symlink

**Input:** Documented install created a symlink (e.g. in a standard bin directory). User follows documented uninstall steps.

**Expected output:** Symlink removed; if the binary was only referenced by that symlink, it is removed or uninstall steps describe removal so no leftover artifact remains in documented scope.

**Verification:** The product is no longer invocable (e.g. shell does not find it); no leftover symlink or binary in documented install scope.

## Acceptance criteria

- [ ] Uninstall documentation (or uninstall script with docs) exists and describes steps that remove the product and revert install-time changes within the documented install scope.
- [ ] Steps address removal of the product and any environment/profile or symlink changes introduced by the documented install.
- [ ] After following the steps, no broken references point at the former install location.
- [ ] After following the steps, no leftover artifacts that were part of the documented install scope remain.
- [ ] Scope of "what uninstall removes" is clear (aligned with what install docs say was installed or modified).

## Dependencies

- O006/R001 — Uninstall steps align with what the documented install does; both refer to the same install scope.
