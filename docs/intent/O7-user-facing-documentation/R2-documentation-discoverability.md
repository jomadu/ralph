# R2: Discoverability of User-facing Documentation

**Outcome:** O7 — User-facing Documentation

## Requirement

Users can discover where user-facing documentation lives and how to reach it without reverse-engineering the repository, so that the outcome "users have access to documentation" is achievable in practice.

## Specification

**Entry points:** At least one of the following (or equivalent) must direct users to the user-facing documentation:

- **Repository README:** The main README (e.g. root `README.md`) includes a link or section that points to the user-facing documentation (e.g. "User documentation" or "Guides" linking to `docs/user/` or its index). The link is visible in the normal flow of the README (e.g. in a "Documentation" or "Getting started" section), not only in a footer or contributor-only section.
- **Index in place:** The user docs directory has an index (per R1) so that anyone who lands in `docs/user/` (e.g. via README link or repository browser) can see the list of topics and navigate to them.

**Optional but not required:** CLI help (e.g. `ralph --help` or a `ralph docs` subcommand) may reference or link to the docs; it is not mandatory for this requirement. The minimum is that a user who opens the repository can find the user docs from the README or by navigating to the designated directory.

**Clarity:** The entry point (e.g. README sentence or link text) must make it clear that the target is *user* documentation (how to use Ralph), not only contributor or specification content. For example, "See [user guides](docs/user/)" or "Documentation: [User docs](docs/user/)" is sufficient; a generic "Documentation" link that mixes intent tree and user docs should still allow users to identify the user-facing surface (e.g. by naming or structure).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User arrives at the repo for the first time | They can find a link to user docs from the README (or equivalent) and open the index. |
| User is in a packaged distribution (e.g. release tarball) with no web README | If the repo is bundled, the same README and `docs/user/` are present; discoverability is unchanged. |
| Multiple entry points (README + CLI help) | Both may point to user docs; at least one (README or index) is required. |

### Examples

#### README links to docs/user

**Input:** Root `README.md` has a section "Documentation" with a link to `docs/user/README.md` labeled "User guides and workarounds."

**Expected:** A user reading the README can click or follow the link and reach the user docs index. No need to guess that user docs live under `docs/user/`.

**Verification:** README contains a clear link to the user docs index; link text or context indicates it is for end users.

#### No README link

**Input:** README does not link to `docs/user/`; only `docs/intent/` is linked.

**Expected:** Requirement is not satisfied. At least one discoverability entry point (README or equivalent) must point to user docs.

**Verification:** Add a README link (or another agreed entry point) to `docs/user/`.

## Acceptance criteria

- [ ] At least one entry point (e.g. README) directs users to the user-facing documentation (R1 location/index)
- [ ] The entry point makes it clear that the target is user documentation (how to use Ralph)
- [ ] A user opening the repository can find and open the user docs index without guessing paths

## Dependencies

- R1 — Discoverability points at the location and structure defined in R1.
