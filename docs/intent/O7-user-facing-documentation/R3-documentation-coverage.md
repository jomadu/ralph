# R3: Coverage of User-facing Documentation Topics

**Outcome:** O7 — User-facing Documentation

## Requirement

User-facing documentation covers the topics that enable users to use Ralph effectively: configuration, AI backends and aliases, workarounds where relevant, and install/uninstall, so that important workflows are documented and the outcome "users succeed without reverse-engineering" is met.

## Specification

**Topics to cover:** The set of user-facing docs (under the location and structure of R1) must cover at least the following areas. Coverage may be in one or more documents; each area must be addressable by a user trying to accomplish the corresponding task.

| Area | Purpose | Example content |
|------|---------|-----------------|
| **Configuration** | How to configure the loop (iterations, timeouts, signals, config file location). | Where config files live; key options; precedence (config vs env vs CLI). May reference O2/R8, O2/R9 or summarize. |
| **AI backends / aliases** | How to choose or override the AI CLI (built-in aliases, user-defined aliases, direct command). | What built-in aliases exist; how to set `ai_cmd_alias` or `ai_cmd`; how to add a custom alias. May reference O3. |
| **Workarounds** | Where a backend or workflow needs extra steps (e.g. wrapper script for Cursor Agent). | Optional wrapper for Cursor Agent (plain-text stdout); how to point Ralph at the wrapper; dependencies (jq, agent). Other workarounds as they arise from other outcomes. |
| **Install / uninstall** | How to install and uninstall Ralph so it is invocable and removable. | May be in README or a dedicated doc; must align with O6 and its documentation requirement (O6/R5). Can link to the authoritative install/uninstall reference. |

**Task-oriented:** Content is written so a user can complete a task (e.g. "Run with Cursor Agent," "Override the cursor-agent alias with the wrapper," "Install Ralph on macOS"). Purely reference material (e.g. exhaustive flag list) may live in the intent tree or CLI help; user docs focus on how to achieve goals.

**Gaps:** If a topic is not yet documented (e.g. a new backend workaround), that is a gap to be filled; the requirement is that the *set* of topics above is covered. New topics that emerge from other outcomes (e.g. a new O3 workaround) should be added to user docs and to the index (R1, R5).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Install/uninstall is documented only in README | Satisfies coverage if README (or a linked doc) is the authoritative reference per O6/R5 and is linked from or listed in user docs index. |
| A backend has no workaround | No workaround doc is required for that backend. Coverage requires workarounds "where relevant" (e.g. Cursor Agent wrapper). |
| New outcome adds a user-visible workflow | When the outcome is implemented, user docs should be updated to cover the new workflow (R4, R5). |

### Examples

#### Cursor Agent workaround documented

**Input:** User wants to use Cursor Agent with plain-text stdout for signal scanning.

**Expected:** A user-facing doc (e.g. `docs/user/cursor-agent-workaround.md`) explains the built-in alias (raw command), the optional wrapper, how to point Ralph at the wrapper (config or `--ai-cmd`), and dependencies. The topic is listed in the user docs index (R1).

**Verification:** User can complete the task "use Cursor Agent with the wrapper" using only user docs.

#### Configuration overview

**Input:** User wants to set iteration limit and signal strings.

**Expected:** User docs cover (in one or more pages) where to put config, key options (e.g. `default_max_iterations`, `signals.success`), and that CLI flags override config. May point to full reference (intent or CLI) for exhaustive detail.

**Verification:** User can achieve "configure the loop" without reading the intent tree.

## Acceptance criteria

- [ ] Configuration (loop options, config file, precedence) is covered in user-facing docs
- [ ] AI backends and aliases (built-in, user-defined, direct command) are covered
- [ ] Workarounds that are referenced in the intent tree (e.g. Cursor Agent wrapper) have a corresponding user-facing doc
- [ ] Install and uninstall are covered (in user docs or via clear link to authoritative reference per O6)
- [ ] Content is task-oriented so users can complete common workflows

## Dependencies

- R1 — Topics live under the location and structure of R1; index lists them.
- O6/R5 — Install/uninstall content must align with the authoritative install/uninstall documentation.
- Other outcomes (e.g. O3) — Workarounds or behaviors specified there may require a user-facing topic; coverage includes those that are user-visible.
