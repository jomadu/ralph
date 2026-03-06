# R5: Policy and Maintenance for User-facing Documentation

**Outcome:** O7 — User-facing Documentation

## Requirement

A clear policy defines what belongs in user-facing documentation versus the intent tree versus the root README, and how to add or update a user-facing doc and how it relates to the intent tree, so that maintainers and contributors can act consistently and traceability is preserved.

## Specification

**Policy — what goes where:**

- **User-facing docs (R1 location, e.g. `docs/user/`):** Task-oriented content for *users* of Ralph: how to configure, run, choose a backend, apply workarounds, install/uninstall. Written in user terms; may summarize or link to intent or CLI reference. Must align with behavior (R4) and be discoverable (R2).
- **Intent tree (`docs/intent/`):** Outcomes, requirements, and specifications for *builders* and maintainers. Source of truth for behavior. Not written as end-user how-tos. User docs may reference outcomes/requirements for traceability.
- **Root README:** High-level project description, quick start, install/uninstall summary or link, and link(s) to user docs (R2) and optionally to intent or contributing. README may duplicate a minimal subset of user doc content (e.g. one paragraph on config file location) for convenience; the canonical detail lives in user docs or intent as appropriate.

**Adding or updating a user-facing doc:**

- **New topic:** When a new user-visible workflow or workaround is introduced (e.g. from a new or changed requirement in another outcome), a corresponding user-facing doc (or section) should be added under the R1 location and registered in the index. The doc should reference the outcome/requirement it aligns with (traceability). R1 and R3 define structure and coverage.
- **Update:** When behavior or intent changes, user docs that describe that behavior are updated per R4. The policy does not require a specific review workflow; it requires that the relationship between user docs and intent is clear so that whoever changes a requirement knows to update the corresponding user doc (if any).
- **Removal:** If a topic becomes obsolete (e.g. a workaround is no longer needed), the user doc can be removed or archived; the index is updated.

**Traceability:** Each user-facing topic document may state which outcome and/or requirement it aligns with (e.g. "Intent: O3 — Backend agnosticism; R1 built-in aliases"). This is not mandatory for every sentence but is required at least at the topic level for topics that implement or explain a specific requirement (e.g. workarounds). The index (R1) may list the intent link per topic. Traceability supports alignment (R4): when the requirement changes, the linked user doc is the one to update.

**Documented in repo:** The policy above is written down in the repository so that contributors can follow it. It may live in the user docs index (`docs/user/README.md`), in a CONTRIBUTING or docs policy file, or in the intent tree (e.g. O7 outcome or this requirement). The policy must be discoverable by someone adding or changing user-facing content.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Contributor adds a new how-to | They place it under `docs/user/`, update the index, add intent link if the topic maps to a requirement, and follow R4 so content matches behavior. |
| Requirement in O3 (or another outcome) references "user docs" | That reference points to the R1 location (e.g. `docs/user/`) or a specific topic (e.g. cursor-agent-workaround.md). Policy here ensures such references stay valid. |
| Dispute over whether something is "user" vs "intent" | User = task-oriented, enables use of Ralph. Intent = specification, source of truth for behavior. If it's both (e.g. workaround), intent has the spec; user doc has the how-to and links to intent. |

### Examples

#### Policy in docs/user/README.md

**Input:** `docs/user/README.md` contains a "Policy" section stating that user docs should not contradict the intent tree, that new topics should be added under this directory and listed in the index, and that topics should link to the relevant outcome/requirement when applicable.

**Expected:** A contributor opening the index sees the policy and knows where to add a new topic and how it relates to the intent tree.

**Verification:** Policy is readable in the repo; R1 index and this requirement are satisfied.

#### New workaround from O3

**Input:** O3/R1 (or another requirement) is updated to mention an optional wrapper for a new backend; the wrapper script is added to the repo.

**Expected:** A user-facing doc is added (e.g. `docs/user/new-backend-workaround.md`), listed in the index with intent link to O3/R1 (or the relevant requirement), and written so it matches the implementation. R3 coverage and R4 alignment are satisfied.

**Verification:** User can find the new workaround in the index; doc traces to intent; content matches behavior.

## Acceptance criteria

- [ ] Policy is documented in the repository (e.g. in user docs index or CONTRIBUTING) stating what belongs in user docs vs intent vs README
- [ ] Policy describes how to add a new user-facing topic (where to put it, how to register in index, traceability to intent)
- [ ] Policy describes that user docs must align with behavior (R4) and that when a requirement changes, linked user docs should be updated
- [ ] Topic-level traceability to outcome/requirement is required for topics that explain or implement a specific requirement (e.g. workarounds)

## Dependencies

- R1 — Policy applies to the location and structure defined in R1; index is where topics are registered.
- R2 — Policy does not duplicate discoverability rules but may reference them (README link to user docs).
- R3 — Coverage defines which topics exist; policy defines how they are added and maintained.
- R4 — Alignment is the ongoing result of following the policy (update docs when behavior or intent changes).
