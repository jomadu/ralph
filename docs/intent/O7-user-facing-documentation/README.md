# O7: User-facing documentation

## Statement

Users have access to user-facing documentation that enables them to use Ralph effectively — including configuration, backends, workarounds, and how to get help — so they can succeed without reverse-engineering the product.

## Why it matters

Ralph’s behavior is specified in the intent tree, but that tree is written for builders and maintainers. Users need a separate surface: discoverable, task-oriented docs that explain how to run the loop, choose a backend, apply workarounds (e.g. Cursor Agent wrapper), configure signals, and understand install/uninstall. Without a dedicated outcome for user-facing documentation, those docs may be missing, scattered, or drift out of alignment with actual behavior. This outcome makes “we have user docs that work” a first-class result of the product.

## Verification

- User can find user-facing documentation (e.g. under `docs/user/`, or linked from README / CLI help) and use it to complete common tasks (e.g. run with Cursor Agent, override an alias, apply the wrapper workaround).
- Documentation covers key topics (configuration, AI backends/aliases, workarounds where relevant, install/uninstall) and is maintained so it stays aligned with product behavior or explicitly states limitations.
- A new user or maintainer can add or update a user-facing doc and know where it lives and how it relates to the intent tree (traceability or policy is defined under requirements).

## Non-outcomes

- The intent tree itself is not “user-facing” in this sense; it remains the source of truth for specification. User-facing docs are derived from or aligned with it, not a replacement.
- We do not require a single doc format or toolchain (e.g. static site generator); the outcome is that user docs exist, are discoverable, and are accurate.
- In-repo docs are in scope; external sites, videos, or community wikis are out of scope unless we explicitly document them as the canonical user surface.

## Risks

| Risk | Mitigating Requirement |
|------|------------------------|
| User docs are scattered or contributors don't know where to add them | [R1 — User-facing documentation location and structure](R1-documentation-location-and-structure.md) |
| Users cannot discover where the docs are | [R2 — Discoverability of user-facing documentation](R2-documentation-discoverability.md) |
| User docs are missing for important workflows (e.g. backends, workarounds) | [R3 — Coverage of user-facing documentation topics](R3-documentation-coverage.md) |
| User docs drift from actual behavior or intent | [R4 — Alignment of user-facing documentation with behavior](R4-documentation-alignment.md) |
| No clear policy for what belongs in user docs vs intent vs README | [R5 — Policy and maintenance for user-facing documentation](R5-documentation-policy-and-maintenance.md) |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R1](R1-documentation-location-and-structure.md) | User-facing documentation location and structure (canonical directory, index, topic files) | ready |
| [R2](R2-documentation-discoverability.md) | Discoverability of user-facing documentation (entry points so users can find docs) | ready |
| [R3](R3-documentation-coverage.md) | Coverage of user-facing documentation topics (configuration, backends, workarounds, install/uninstall) | ready |
| [R4](R4-documentation-alignment.md) | Alignment of user-facing documentation with behavior (maintain accuracy, state limitations) | ready |
| [R5](R5-documentation-policy-and-maintenance.md) | Policy and maintenance for user-facing documentation (what goes where, how to add/update, traceability) | ready |
