# R1: User-facing Documentation Location and Structure

**Outcome:** O7 — User-facing Documentation

## Requirement

User-facing documentation lives in a defined location and structure within the repository so that there is one canonical place for user docs and contributors know where to add or update them.

## Specification

**Location:** User-facing documentation is stored under a designated directory in the repository. The project uses `docs/user/` as the canonical location. This directory is the single root for task-oriented user docs (how-tos, workarounds, examples). Other locations (e.g. root README, intent tree under `docs/intent/`) are not considered the "user-facing documentation" surface for this requirement; they may link to or summarize it.

**Structure:**

- **Index:** The directory contains an index (e.g. `docs/user/README.md`) that lists available topics with short descriptions and, where applicable, links to the related outcome or requirement in the intent tree. The index is the entry point for browsing user docs.
- **Topic documents:** Individual topics are stored as separate documents (e.g. `docs/user/cursor-agent-workaround.md`). Naming and format are consistent enough that a maintainer can add a new topic and know where to place it and how to register it in the index.
- **Traceability:** Each topic document may reference the intent tree (outcome and/or requirement) it aligns with, so that the relationship between user-facing content and specification is explicit.

**Out of scope:** This requirement does not mandate a specific file format (e.g. Markdown only), build step, or static site. It mandates a defined location and a clear structure so that user docs are findable and maintainable.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| New workaround or how-to is needed (e.g. from a new requirement in another outcome) | Contributor adds a new file under `docs/user/`, updates the index, and links to the relevant requirement if applicable. |
| Project has no user docs yet | Directory and index exist (possibly with a single placeholder or one topic); structure is documented so future topics have a home. |
| User doc is better placed in README (e.g. quick start) | README can contain or link to content; the canonical *body* of user-facing docs remains under the designated directory. Cross-links between README and `docs/user/` are acceptable. |

### Examples

#### Index with topic list

**Input:** `docs/user/README.md` exists.

**Expected:** README lists topics (e.g. "Cursor Agent workaround") with a short description and an optional "Intent" column or inline link to the outcome/requirement (e.g. O3, R1). Policy or guidance for adding topics is stated (or delegated to R5).

**Verification:** A new maintainer can open `docs/user/README.md` and see what user docs exist and where to add more.

#### New topic added

**Input:** A new backend workaround is specified in the intent tree (e.g. O3) and needs a user-facing page.

**Expected:** A new file is created under `docs/user/` (e.g. `docs/user/some-backend-workaround.md`). The index is updated to include the new topic and its intent link. The topic content is written in user terms and does not contradict the requirement.

**Verification:** User can find the new topic via the index; requirement traceability is preserved.

## Acceptance criteria

- [ ] A single canonical directory for user-facing documentation is defined and used (e.g. `docs/user/`)
- [ ] An index under that directory lists available topics and (where applicable) their link to the intent tree
- [ ] Topic documents live under the same directory; structure is documented or obvious so contributors know where to add new topics
- [ ] Topic documents may reference outcome/requirement for traceability

## Dependencies

_None. R5 (policy and maintenance) may reference this structure when defining how to add or update user docs._
