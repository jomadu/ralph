# Plan: O7 User-facing Documentation

This plan implements **O7 — User-facing documentation** so users have access to documentation that enables them to use Ralph effectively. Outcome and requirements are in `docs/intent/O7-user-facing-documentation/`.

**Verification (from O7 README):** User can find user-facing documentation (e.g. under `docs/user/`, or linked from README) and use it to complete common tasks (e.g. run with Cursor Agent, override an alias, apply the wrapper workaround). Documentation covers key topics (configuration, backends, workarounds, install/uninstall) and is maintained so it stays aligned with product behavior. A new maintainer can add or update a user-facing doc and know where it lives and how it relates to the intent tree.

**Current state:** `docs/user/` exists with an index (`README.md`) and one topic (`cursor-agent-workaround.md`) that includes intent links. The index has a short Policy section. The root README does **not** yet link to `docs/user/`, so discoverability (R2) is incomplete. Coverage (R3) is partial: workarounds (Cursor Agent) are covered; configuration, backends/aliases overview, and install/uninstall are not yet dedicated user-doc topics or linked from the index. Policy (R5) is summarized in the index but not fully spelled out (what goes where, how to add/update, traceability). Alignment (R4) has not been done as a formal pass.

---

## Task dependency overview

```
T1 (R1 — structure and index)
  ├→ T2 (R2 — README link)
  ├→ T3 (R5 — policy)
  ├→ T4 (R3 — configuration topic)
  ├→ T5 (R3 — backends/aliases topic)
  └→ T6 (R3 — install/uninstall coverage)
        ↓
T7 (R4 — alignment and traceability)  [depends on T2, T3, T4, T5, T6]
```

T2–T6 can be done in parallel after T1. T7 should run after all content and entry points are in place.

---

## T1: R1 — User-facing documentation location and structure

**Priority:** 1 (must be first)  
**Dependencies:** None  
**Spec:** [R1 — User-facing documentation location and structure](docs/intent/O7-user-facing-documentation/R1-documentation-location-and-structure.md)

**Objective:** Confirm and, if needed, complete the canonical location and structure for user docs so that there is one place for user docs and contributors know where to add or update them. No new content topics in this task — only structure and index format.

**Context:**

- **Location:** The project uses `docs/user/` as the canonical directory. Ensure it exists and is the single root for task-oriented user docs.
- **Index:** `docs/user/README.md` must list all current topics with short descriptions and an "Intent" column (or equivalent) linking to the related outcome/requirement. The existing table (Document, Description, Intent) already does this for the Cursor Agent workaround; verify every topic is listed and intent links are correct.
- **Adding topics:** The index should include a brief note on how to add a new topic (e.g. "Add a new file under `docs/user/`, add a row to the Topics table with intent link if applicable") or an explicit pointer to the policy (T3). R1 says the index is the entry point; R5 will document the full policy — here we only ensure the index makes it obvious where new topics go and how they are registered.
- **Traceability:** Each topic document should reference the intent tree (outcome and/or requirement) it aligns with. Verify `cursor-agent-workaround.md` has correct intent links; no other topic docs exist yet.

**Acceptance:**

- [ ] `docs/user/` exists and is the designated location for user-facing docs.
- [ ] `docs/user/README.md` lists all current topics (at least Cursor Agent workaround) with description and intent link.
- [ ] Index includes a one-line or short "how to add a topic" note or pointer to policy (so contributors know where to put new docs and how to register them).
- [ ] Existing topic doc(s) include intent references at the topic level.

---

## T2: R2 — Discoverability (README link to user docs)

**Priority:** 2  
**Dependencies:** T1  
**Spec:** [R2 — Discoverability of user-facing documentation](docs/intent/O7-user-facing-documentation/R2-documentation-discoverability.md)

**Objective:** Add at least one entry point from the repository so users can discover user-facing documentation without guessing. The minimum is a link from the root README to `docs/user/` (or its index).

**Context:**

- **README:** The root `README.md` currently has no link to `docs/user/`. Add a section (e.g. "Documentation" or "User guides") or a clear link within an existing section (e.g. after "Configuration" or near the top) that points to `docs/user/` or `docs/user/README.md`. Link text must make it clear the target is *user* documentation (e.g. "User guides and workarounds", "User documentation").
- **Placement:** The link must be visible in the normal flow of the README, not only in a footer or contributor-only section. R2 allows "Documentation" or "Getting started" as natural places.
- **Optional:** CLI help or a `ralph docs` subcommand is not required for this task; README link is sufficient for R2.

**Acceptance:**

- [ ] Root README contains a link to `docs/user/` (or `docs/user/README.md`).
- [ ] Link text or surrounding context indicates the target is user documentation (how to use Ralph).
- [ ] A user opening the repository can find and open the user docs index without guessing paths.

---

## T3: R5 — Policy in user docs index

**Priority:** 3  
**Dependencies:** T1  
**Spec:** [R5 — Policy and maintenance for user-facing documentation](docs/intent/O7-user-facing-documentation/R5-documentation-policy-and-maintenance.md)

**Objective:** Document the full policy for user-facing documentation in the repository so maintainers and contributors know what belongs in user docs vs intent vs README, how to add or update a topic, and how traceability to the intent tree works.

**Context:**

- **Where:** Policy can live in `docs/user/README.md` (e.g. a "Policy" or "For contributors" section), in a separate file under `docs/user/`, or in CONTRIBUTING. R5 says it must be discoverable by someone adding or changing user-facing content; the user docs index is a natural place.
- **What to include:**
  - **What goes where:** User-facing docs (`docs/user/`) = task-oriented content for users (how to configure, run, choose backend, workarounds, install/uninstall). Intent tree (`docs/intent/`) = source of truth for behavior; not end-user how-tos. README = high-level description, quick start, link(s) to user docs; may summarize or link to install/uninstall.
  - **Adding a topic:** Add a new file under `docs/user/`, add a row to the index Topics table, and include an intent link (outcome/requirement) when the topic implements or explains a specific requirement (e.g. workarounds).
  - **Updating:** When behavior or intent changes, update any user doc that describes that behavior so docs stay aligned (R4). The policy should state that linked user docs must be updated when the corresponding requirement changes.
  - **Traceability:** Topics that explain a specific requirement must have topic-level intent link; the index may list intent per topic (already in place).
- **Length:** Keep the policy concise but complete enough that a new contributor can follow it. You may reference O7/R5 for full specification.

**Acceptance:**

- [ ] Policy is written in the repo (e.g. in `docs/user/README.md` or a linked doc) and is discoverable from the user docs index.
- [ ] Policy states what belongs in user docs vs intent vs README.
- [ ] Policy describes how to add a new user-facing topic (where to put it, how to register in index, intent link when applicable).
- [ ] Policy states that when a requirement changes, linked user docs should be updated to stay aligned.

---

## T4: R3 — Configuration topic

**Priority:** 4  
**Dependencies:** T1  
**Spec:** [R3 — Coverage of user-facing documentation topics](docs/intent/O7-user-facing-documentation/R3-documentation-coverage.md)

**Objective:** Add a user-facing document that covers configuration: where config files live, key loop options, and precedence (config vs env vs CLI), so users can complete the task "configure the loop" without reading the intent tree.

**Context:**

- **Content:** Where config files are (workspace `./ralph-config.yml`, global `~/.config/ralph/ralph-config.yml`), key options (e.g. `default_max_iterations`, `failure_threshold`, `signals.success`/`signals.failure`, `ai_cmd_alias`), and that CLI flags override env override config. Task-oriented (e.g. "To change the iteration limit, set `default_max_iterations` in your workspace config or use `ralph run build -n 20`"). May point to README or CLI help for exhaustive flag list.
- **Intent:** Reference O2 (configurable behavior) and optionally O2/R8 (env vars), O2/R9 (CLI) or summarize. Add intent link at topic level per R5.
- **Index:** Add the new topic to the Topics table in `docs/user/README.md` with description and intent link.
- **Format:** One new file under `docs/user/` (e.g. `configuration.md` or `config-and-signals.md`). Naming should be clear and consistent with existing style (e.g. `cursor-agent-workaround.md`).

**Acceptance:**

- [ ] A new topic document under `docs/user/` covers configuration (config file location, key options, precedence).
- [ ] Content is task-oriented so a user can achieve "configure the loop" (e.g. set iterations, signals, backend alias).
- [ ] Topic includes intent link (e.g. O2).
- [ ] Index in `docs/user/README.md` lists the new topic with description and intent.

---

## T5: R3 — AI backends and aliases topic

**Priority:** 5  
**Dependencies:** T1  
**Spec:** [R3 — Coverage of user-facing documentation topics](docs/intent/O7-user-facing-documentation/R3-documentation-coverage.md)

**Objective:** Add a user-facing document that covers how to choose or override the AI CLI: built-in aliases, user-defined aliases, and direct command, so users can complete tasks like "run with a different backend" or "add a custom alias". The Cursor Agent workaround remains a separate topic; this one is the general overview with a pointer to workarounds where relevant.

**Context:**

- **Content:** What built-in aliases exist (claude, kiro, copilot, cursor-agent) and what they resolve to (one-line each); how to set `ai_cmd_alias` or `ai_cmd` in config or via `--ai-cmd` / `--ai-cmd-alias`; how to add a custom alias in `ai_cmd_aliases`. Mention that some backends have optional workarounds (e.g. Cursor Agent) and link to the relevant user doc(s).
- **Intent:** Reference O3 (backend agnosticism) and optionally O3/R1 (built-in aliases), O3/R3 (user-defined). Add intent link at topic level.
- **Index:** Add the new topic to the Topics table in `docs/user/README.md`.
- **Overlap:** Do not duplicate the full Cursor Agent workaround content here; link to `cursor-agent-workaround.md` for that. This doc is the "how to choose/override backend" overview.

**Acceptance:**

- [ ] A new topic document under `docs/user/` covers AI backends and aliases (built-in, user-defined, direct command).
- [ ] Content is task-oriented; includes link to Cursor Agent workaround (and any other workaround docs) where relevant.
- [ ] Topic includes intent link (e.g. O3).
- [ ] Index lists the new topic with description and intent.

---

## T6: R3 — Install/uninstall coverage

**Priority:** 6  
**Dependencies:** T1  
**Spec:** [R3 — Coverage of user-facing documentation topics](docs/intent/O7-user-facing-documentation/R3-documentation-coverage.md)

**Objective:** Ensure install and uninstall are covered for users. R3 allows coverage either as a dedicated user doc or as a clear link from the user docs index to the authoritative install/uninstall reference (e.g. README or O6/R5 doc). Choose the approach that fits the current repo: if install/uninstall are documented in README or a dedicated doc, add a topic row in the index that links to that content; if not, add a short user-facing doc under `docs/user/` that explains how to install and uninstall (or links to the canonical reference once it exists per O6).

**Context:**

- **O6:** Install/uninstall behavior and documentation are specified in O6 (e.g. R5 install/uninstall documentation). If PLAN_INSTALL_UNINSTALL or O6 implementation has already produced an authoritative reference (e.g. README section or `docs/INSTALL.md`), link from user docs to it. If the authoritative reference does not yet exist, this task can add a minimal user doc (e.g. "Install and uninstall") that states the current situation (e.g. "Build from source and place the binary on PATH; see README" or "See [Install guide](...) when available") so the index has an install/uninstall entry and R3 coverage is satisfied.
- **Index:** Add an "Install and uninstall" topic to the Topics table. If the content lives in README or elsewhere, the Description can be "How to install and uninstall Ralph" and the link can point to the relevant section or file. Intent link: O6.
- **No duplicate authority:** Do not create a second authoritative install/uninstall spec; user docs either link to the O6/R5 reference or provide a short how-to that aligns with it.

**Acceptance:**

- [ ] Install and uninstall are covered: either a dedicated user doc under `docs/user/` or a clear link from the user docs index to the authoritative reference (README or O6 doc).
- [ ] User docs index has an entry for install/uninstall (description and intent link to O6).
- [ ] Content or link aligns with O6 (and O6/R5 when that exists); no contradictory instructions.

---

## T7: R4 — Alignment pass and traceability

**Priority:** 7  
**Dependencies:** T2, T3, T4, T5, T6  
**Spec:** [R4 — Alignment of user-facing documentation with behavior](docs/intent/O7-user-facing-documentation/R4-documentation-alignment.md)

**Objective:** Review all user-facing docs for alignment with current product behavior, ensure each topic has an intent link where required, and state any known limitations so that O7 verification (docs enable effective use and stay aligned) is met.

**Context:**

- **Scope:** Every document under `docs/user/` (index plus all topic docs added in T1–T6): Cursor Agent workaround, configuration topic, backends/aliases topic, install/uninstall coverage, and the index/README content for user docs.
- **Alignment:** For each topic, check that the described behavior matches the product (e.g. config precedence, alias names, paths, CLI flags). If something is not yet implemented or is platform-limited, add a short caveat (e.g. "Install script is currently tested on macOS and Linux only."). Fix any outdated or incorrect statements.
- **Traceability:** Confirm every topic that implements or explains a specific requirement has an intent link at the topic level (outcome and/or requirement). Confirm the index lists intent for each topic. Add any missing links.
- **Policy:** Confirm the policy (T3) is reflected in practice: new topics are under `docs/user/`, in the index, with intent links where applicable. No need to repeat the policy text; this task is a one-time alignment and traceability pass.
- **Ongoing:** R4 is about maintaining alignment over time; this task establishes the baseline. The policy (R5) already states that when requirements change, linked user docs should be updated.

**Acceptance:**

- [ ] Every user-facing topic doc has been reviewed for accuracy against current behavior; corrections or limitation statements added where needed.
- [ ] Every topic that explains a specific requirement has an intent link (topic level and/or in index).
- [ ] No user doc contradicts the intent tree or current implementation.
- [ ] Known limitations (e.g. platform support) are stated where relevant.

---

## Beads (bd) issue mapping (badges)

| Plan task | Bead ID    | Dependencies (blocked by) |
|-----------|-------------|---------------------------|
| T1        | ralph-d6i   | —                         |
| T2        | ralph-bzh   | ralph-d6i                 |
| T3        | ralph-d4i   | ralph-d6i                 |
| T4        | ralph-rig   | ralph-d6i                 |
| T5        | ralph-t25   | ralph-d6i                 |
| T6        | ralph-3o2   | ralph-d6i                 |
| T7        | ralph-2ms   | ralph-bzh, ralph-d4i, ralph-rig, ralph-t25, ralph-3o2 |

**Check ready work:** `bd ready` or `bd ready --json`. **Claim:** `bd update <id> --claim`. **Close:** `bd close <id> --reason "Completed"`.
