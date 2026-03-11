# Product and Engineering Documentation

This document defines the model for product and engineering documentation. **Product** documentation describes who we're building for, what we're building, and why — outcomes and requirements at the level of *intent*. **Engineering** documentation describes where in the system each capability lives and how the system is built — component-centric architecture, implementation specifications (schemas, interfaces, protocols), and references to product requirements. Together they provide a complete picture: product is the source of truth for intent (who, what, why); engineering is the source of truth for implementation (where and how — structure, placement, and the hard specs implementers build from).

## Structure

```
docs/
  product/
    README.md                        # Index of all outcomes
    O001-<slug>/                     # Outcome dirs use three-digit zero-padded IDs
      README.md                      # Outcome definition (risks, requirement one-liners)
      R001-<slug>.md                 # Requirement — complete for intent (IDs zero-padded to three digits)
  engineering/
    README.md                        # Overview: diagram, design principles, component list with O/R assignments
    components/
      <component>.md                 # Component: responsibility, interfaces, implementation specs; O/R only in README
      <component>/                   # Or component as directory (README + schema docs, design notes as needed)
```

Product answers *who*, *what*, and *why* at the level of intent. Engineering answers *where* and *how*: structure, placement, and the implementation specifications (e.g. config schema, APIs, protocols) that implementers build from. Engineering does not re-specify user-facing intent; it owns the hard specs for how the system is implemented and references product requirements by ID (e.g. O001/R002, O004/R001).

### Product Hierarchy

**Outcome** — A measurable change in the world when this product exists. Outcomes are user-facing and verifiable. They answer *who* we're building for, *what* is true when we succeed, and *why* it matters.

**Requirement** — A capability needed to deliver an outcome. Each requirement belongs to exactly one outcome. The requirement document is **complete for intent**: it contains the capability statement, acceptance criteria, edge cases, and examples that define *what* we're committing to and how we'll verify it. Detailed implementation specifications (schemas, wire formats, algorithms, strict interfaces) may live in engineering; when they do, the requirement states the capability and may reference the engineering doc. The requirement is the single source of truth for *intent*; engineering is the source of truth for the implementation spec. Requirements elaborate the *what* from the outcome: *"what must be true?"* in enough detail that stakeholders and implementers understand the commitment; implementers use engineering docs for the buildable spec.

Every requirement traces to an outcome. If a requirement has no outcome above it, it's unjustified.

### Engineering Hierarchy

**Architecture** — The structure of the system: major components, their responsibilities, and how they interact. Documented in `docs/engineering/` with a component-centric layout.

**Component** — A named part of the system (e.g. run-loop, config, backend, review). A component is documented either as a single file (`components/<component>.md`) or as a directory (`components/<component>/`). A component directory standardly includes a README.md as the primary doc, plus optional supporting files (e.g. a strict schema doc for config YAML, API contracts, state diagrams). Requirement assignments (which O/R IDs a component satisfies) live only in the engineering README; component docs do not duplicate that list. Each component's primary doc states responsibility, key interfaces, and the **implementation specifications** that implementers build from — schemas, protocols, data structures — and references the engineering README for the list of requirements it implements. Product is the source of truth for *intent* (what we're committing to); the engineering README is the source of truth for *which component implements which requirements*; the component doc is the source of truth for *how* that component is implemented (where it lives, how it connects, and the strict specs it adheres to).

## Consistency

### Product

Internal consistency is a first-class concern at every step. Each layer of the product tree introduces more documents and more surface area for contradiction.

At every step boundary, before locking, perform a consistency review across two dimensions:

- **Vertical consistency** — Does each document agree with the layer above it? Outcome detail files must match the index. Requirements must match their parent outcome.
- **Horizontal consistency** — Do documents at the same level agree with each other? Outcome files must not make contradictory claims or assume incompatible models of the system. Requirements across different outcomes must not prescribe conflicting behaviors.

Inconsistency at any layer invalidates everything below it. Catching it early is cheap; catching it in implementation is worse.

### Engineering

- Every requirement ID in the engineering README's component list must exist in the product tree (assignments live only in the README).
- Component docs must not contradict each other (e.g. two components claiming to implement the same requirement in incompatible ways).
- Component docs do not duplicate the O/R assignment list; the engineering README is the single place for that.
- When product requirements change, engineering docs (including implementation specs) must be updated so component responsibilities and references remain accurate.
- When an engineering implementation spec changes in a way that affects user-facing contract, the product requirement must be updated.

## Phased Development

Product and engineering each follow a **compress, then expand** pattern: write one-liners (or index-level content) first, lock, then expand into full documents. Compressed forms are cheap to review for consistency; expanded forms are reviewed against the locked compression above them.

### Product (P1–P4)

Build the product tree in four steps. Each step produces artifacts, reviews them, and locks them before the next step begins.

#### P1: Outcome Index

Write the root `docs/product/README.md`. One table — each outcome as a one-line statement. No directories, no detail. This is the product on a single page. Verification criteria are *not* in the index; they are defined in each outcome's README (P2). Keeping verification in the outcome document avoids one-line compression and drift between index and detail.

**Review:**
- Each outcome is a present-tense assertion about the world, not a feature description
- Each outcome is clear about *who* it is for (users, roles, or personas) and *why* it matters
- Outcomes don't overlap — if two outcomes could share a requirement, they may be the same outcome
- The set of outcomes is complete — together they describe the whole product (who, what, why)
- The set is minimal — removing any outcome would leave a gap

Lock the index before expanding.

#### P2: Outcome Detail

Expand each outcome into its own directory and README. Who it's for, statement, why it matters, full verification, non-outcomes. These files define the problem space (who, what, why). They do not reference requirements or risks — those come later.

**Review:**
- Each outcome detail is consistent with its one-liner in the index — if they diverge, fix the index first
- Each outcome explicitly states *who* it is for and *why* it matters
- Verification criteria are defined here (in the outcome README) and are observable by a user, not by a test suite
- Non-outcomes are clear enough that someone could push back on scope and point to this list
- Read all outcome files as a set: no two outcomes make contradictory claims, imply overlapping scope, or assume incompatible models of who we're building for, the system, or the domain

Lock outcome detail before proceeding. Changes after this point ripple through everything below.

#### P3: Requirement Index

Outcomes don't arrive pre-decomposed. Use the following decomposition to derive requirements and risks; then for each outcome, append two tables to its README: a risks table mapping each risk to its mitigating requirement, and a requirements table with one-line summaries. No requirement documents yet.

**Decomposition**

- **Why/how/how-else chain** (from goal-oriented requirements engineering): Starting from an outcome, ask (1) *"How is this outcome achieved?"* — each distinct answer is a candidate requirement (if still abstract, decompose further; if concrete enough to build, it's a requirement). (2) *"How else?"* — ask again to avoid locking into one design and to surface gaps. (3) *"Why does this requirement exist?"* — every candidate must trace back to the outcome; if not, it's misplaced or unjustified.
- **Risk analysis:** For each outcome, ask *"What could prevent this from being true?"* Each answer is a risk. Risks surface requirements you'd otherwise discover late. Each risk is documented in the outcome README with the requirement that mitigates it. Every risk must map to an existing requirement or surface a new one. Examples: *"The AI CLI crashes mid-execution"* → process crash recovery; *"Both success and failure signals appear"* → signal precedence rules; *"The user doesn't know which config value is active"* → provenance tracking.
- **Completeness:** An outcome is fully decomposed when every "how" has a requirement, every risk has a mitigating requirement, every requirement traces back via "why," and you can't name a realistic failure that nothing addresses. Sufficiency for building is the goal, not exhaustive enumeration.

This is the compressed form of the requirements layer. All requirement one-liners and risk mappings across all outcomes should be reviewable as a set.

**Review:**
- Every requirement traces to exactly one outcome
- Requirements are capabilities ("the system detects X"), not implementations ("use a regex to scan for X")
- No requirement is redundant with another under the same outcome
- Every risk maps to a mitigating requirement
- The set of requirements under each outcome is sufficient — you can't describe a realistic failure that nothing addresses
- Read all requirement one-liners across all outcomes as a single set: no two requirements prescribe contradictory capabilities or overlap in scope

Lock requirement one-liners before expanding.

#### P4: Requirement Detail (Complete for Intent)

Expand each requirement into its own file within its outcome directory. The requirement document is **complete for intent**: it contains the capability statement, acceptance criteria, edge cases, and examples that define what we're committing to and how we'll verify it. Implementation specifications (e.g. a strict config schema, wire format, algorithm) may live in engineering; when they do, the requirement states the capability and may reference the engineering doc. The requirement is the single source of truth for *intent* for that capability; engineering holds the single source of truth for the implementation spec.

**Review:**
- Each requirement document is consistent with its one-liner in the outcome README — if they diverge, fix the one-liner first
- The requirement is complete for intent — clear what we're committing to; acceptance criteria are testable; implementers use product for intent and engineering for implementation specs to build
- Edge cases are enumerated where relevant; examples are concrete (input, expected output, verification)
- The requirement doesn't exceed what the outcome justifies (no gold-plating)
- Read all requirement documents across all outcomes as a single set: no two requirements prescribe contradictory behaviors, make incompatible assumptions, or define the same concept differently
- Requirements remain vertically consistent with their parent outcome detail

Lock requirements before proceeding. Requirement changes ripple; keep them stable once locked.

### Engineering (E1–E2)

Engineering is developed around the product requirements: components are named and requirements are assigned to them. Engineering documentation can begin once the product has at least a requirement index (P3) — ideally P4 is locked so requirement IDs are stable. Two steps: overview (compress), then component detail (expand).

#### E1: Overview

Write `docs/engineering/README.md`. Architecture on one page: purpose of engineering docs, high-level diagram or flow (e.g. CLI → config → run path / review path), a list of components with one-line descriptions, and for each component its **assigned requirement IDs** (O/R). The README is the single place for the full map: component names, one-liners, and requirement assignments. No per-component detail files yet.

**Decomposition**

Use the product requirement set (and outcome index) to derive the component set. Cluster requirements that hang together — by flow (e.g. everything in the run path), by concern (e.g. config resolution, backend invocation), or by user-facing boundary (e.g. review vs. run). For each cluster, name a component and write a one-liner (what this part of the system is responsible for). Ask *"What part of the system is responsible for this requirement?"* to assign each requirement to a component. Refine until the set is distinct, covers the product, and has no overlapping responsibilities. Record the result in the engineering README: component list with one-liners and O/R assignments (e.g. `run-loop — Runs the iteration loop; decides continue/exit — O001/R002, O001/R004, O001/R005, O004/R001`).

**Review:**
- Components are distinct; no two components have the same responsibility
- The set gives a plausible coverage of the system (every product area has a home)
- The overview is consistent with the product outcome index (no components that imply outcomes or requirements not in product)
- Every product requirement is assigned to at least one component
- No component is empty (every component has at least one requirement)
- Boundaries are clear — no two components claim the same requirement in conflicting ways
- Every O/R in the README exists in the product tree

Lock the engineering README before expanding.

#### E2: Component Detail and Implementation Specs

Flesh out each component as a file (`components/<component>.md`) or directory (`components/<component>/` with README.md standard). Include responsibility, interfaces, and the **implementation specifications** that implementers build from — e.g. the canonical config YAML schema, API contracts, state machines, protocols. Do not duplicate the O/R assignment list (that stays in the engineering README). Reference product requirements by following the assignments in the README for intent; the component doc is the authoritative spec for *how* this component is implemented.

**Review:**
- Each component doc matches its one-liner in the engineering README — if they diverge, fix the README first
- Component docs do not list O/R IDs; the engineering README is the single place for assignments
- Implementation specs (schemas, interfaces) are complete enough that an implementer can build from them
- Interfaces are consistent across components (what one component produces, another consumes as documented)

Lock component docs before treating the architecture as stable. Engineering docs are updated as the system evolves (new requirements, code structure changes); when updating, keep the E1 → E2 hierarchy consistent. When an implementation spec change affects user-facing contract, update the product requirement.

### Working Sessions

This document is a methodology for generating and reviewing the product and engineering trees. It is not a dependency for working within them once they exist. The product documents are self-describing for intent; the engineering documents hold implementation specs and point at product.

**Product:** Chunk by outcome. A session writing or reviewing requirements should load: (1) the outcome index, (2) the single outcome README being worked on, (3) the requirement files under that outcome only. One session per outcome, then a final session for cross-outcome consistency review.

**Engineering:** Chunk by component or by flow. Load the engineering README plus the component doc(s) being written or updated. Cross-component consistency: ensure every referenced O/R exists in product and no two components claim incompatible responsibility for the same requirement. Implementers read product first for intent and acceptance, then engineering for the buildable spec (schemas, interfaces, protocols).

#### Context by Step (Product)

**P1 — Outcome Index:** This methodology document; the repository's main README (product context: who we're building for, high-level what and why).

**P2 — Outcome Detail:** This methodology document; the locked outcome index (`docs/product/README.md`); all outcome READMEs written so far (for horizontal consistency).

**P3 — Requirement Index:** This methodology document; the outcome index; the outcome README being decomposed. After all outcomes have requirement one-liners and risk tables, run a consistency session loading all outcome READMEs.

**P4 — Requirement Detail:** The outcome index; the parent outcome README; all requirement files under that outcome (for horizontal consistency). Work one outcome's requirements per session. After all requirement files exist, run a cross-outcome consistency session.

**Consistency review sessions** load documents *across* outcomes: outcome index plus all documents at the layer being reviewed.

#### Context by Step (Engineering)

**E1 — Overview:** This methodology document; the product outcome index and requirement index (outcome READMEs with requirement one-liners and IDs). Derive components, one-liners, and O/R assignments; write the engineering README. Run a consistency session: every O/R assigned, no component empty, boundaries clear, every O/R exists in product.

**E2 — Component Detail:** The locked engineering README (E1); the product requirement docs for the requirements assigned to the component being written. Write the implementation spec (e.g. config schema, APIs) that implements those requirements. Work one component per session. For cross-component consistency, ensure interfaces align.

### Why This Order Matters

**Product (P1–P4):** The steps alternate between compression and expansion. Compressed forms (outcome one-liners, requirement one-liners) are easy to hold in your head and cheap to review for consistency. Expanded forms (outcome detail, complete requirement documents) add depth but also surface area for contradiction. By locking the compressed form before expanding, you ensure that the detail is anchored to a reviewed, consistent summary. If you write all layers at once, errors in the upper layers silently propagate downward.

**Engineering (E1–E2):** The same pattern. The overview (E1) produces the full map in one place: component names, one-liners, and O/R assignments in the engineering README. Lock it, then expand into component detail (E2) — responsibility and interfaces per component, anchored to that README.

## Product Templates

### Index — `docs/product/README.md`

Open with a short **product summary**: a statement that describes the product itself — who it is for, what it is, and what it does — and that the outcomes below are the measurable results of that product. Avoid meta-description of the documentation (e.g. "this index is…"); the README should read as a statement about the product. Then the outcomes table.

Link each outcome ID to that outcome's README. From the index, use `./O001-<slug>/README.md` (three-digit zero-padded IDs). Each outcome one-liner should be clear on *who* it is for, *what* is true when achieved, and *why* it matters. Verification is defined in each outcome's README, not in the index.

```markdown
# Product

<Short product summary: who the product is for, what it is, and what it does. This product produces the outcomes below.>

## Outcomes

| ID | Outcome |
|----|---------|
| [O001](./O001-<slug>/README.md) | Statement of what is true when this outcome is achieved |
| [O002](./O002-<slug>/README.md) | ... |
```

### Outcome — `O001-<slug>/README.md`

Each outcome directory has a README that fully defines the outcome. Link to requirement docs with `R001-<slug>.md` (three-digit zero-padded IDs).

**Fields:** Who (users/roles/personas), Statement, Why it matters, Verification, Non-outcomes, Risks (table), Requirements (table). Risks and requirements tables are appended in P3. During P1 and P2, the outcome README ends after non-outcomes. Verification here is outcome-level (user-observable evidence); requirement-level testability is captured in each requirement's Acceptance criteria.

```markdown
# O001: <Title>

## Who

<Who this outcome is for: users, roles, or personas.>

## Statement

<One sentence. Present tense. What is true when this outcome is delivered.>

## Why it matters

<The pain without this. What the user suffers today.>

## Verification

<How a user knows this outcome was delivered. Observable evidence, not test cases.>

## Non-outcomes

- <What this does not cover>
- <What this is not responsible for>

## Risks

| Risk | Mitigating Requirement |
|------|----------------------|
| What could prevent this outcome | [R001 — <Title>](R001-<slug>.md) |
| ... | ... |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R001](R001-<slug>.md) | One-line summary | draft / ready / built / verified |
| [R002](R002-<slug>.md) | ... | ... |
```

### Requirement — `R001-<slug>.md`

Each requirement file is **complete for intent**: it contains the requirement statement and all detail needed to define what we're committing to and how we'll verify it. No separate "Specification" section; use whatever structure makes the requirement unambiguous (e.g. Detail, Edge cases, Examples, plus Acceptance criteria and Dependencies). Where implementation specifications (schemas, formats, algorithms) live in engineering, the requirement states the capability and may reference the engineering doc.

**Fields:**

- **Outcome** — Which outcome this requirement serves (traceability upward).
- **Requirement** — What the system must do. Written as a capability statement.
- **Detail** — Whatever is needed to make the requirement complete for intent: edge cases table, examples, behavioral rules. Use subsections (e.g. Edge cases, Examples) as needed. Where the authoritative implementation spec (e.g. config schema, API) lives in engineering, reference it here. Implementers use this doc for intent and the engineering doc for the buildable spec.
- **Acceptance criteria** — Concrete conditions that must be true for this requirement to be considered met.
- **Dependencies** *(optional)* — Other requirements or system capabilities that must exist first. Omit if self-contained.

**Standard subsections when useful:**

- **Edge cases** — Table of boundary conditions and expected behavior.
- **Examples** — Concrete scenarios: Input, Expected output, Verification.

```markdown
# R001: <Title>

**Outcome:** O001 — <Outcome title>

## Requirement

<What the system must do. Capability statement.>

## Detail

<Detail needed for intent: edge cases, examples, behavioral rules. Use subsections as needed. If the implementation spec (e.g. schema, API) lives in engineering, add: "See [engineering doc](link) for the authoritative schema/interface.">

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| <Boundary condition> | <How the system responds> |

### Examples

#### <Scenario name>

**Input:** <Input data, command, or setup>

**Expected output:** <Output data, behavior, or result>

**Verification:** <How to verify the outcome>

## Acceptance criteria

- [ ] <Concrete, testable condition>
- [ ] <...>

## Dependencies

- <O001/R001 — Requirement that must exist first>
```

## Engineering Templates

### Overview — `docs/engineering/README.md`

The single engineering index. It holds everything that would have been a separate "component index": overview plus component list with requirement assignments.

- Purpose of engineering docs (structure, placement, and implementation specs; product holds who, what, and why at the level of intent).
- High-level diagram or description of the system (e.g. CLI → config → run path / review path; run loop, backend, etc.).
- Component list: for each component, one-line description and assigned requirement IDs (O/R). Links to component docs (files or directories under `components/`). Component docs hold the implementation specifications (e.g. config schema, APIs) that implementers build from.

### Component — `docs/engineering/components/<component>.md` or `docs/engineering/components/<component>/`

A component may be a single file (e.g. `run-loop.md`) or a directory (e.g. `run-loop/`). Use a directory when one file is insufficient.

**When using a single file:** the file contains responsibility and interfaces (below). Do not list O/R IDs; those live only in the engineering README.

**When using a directory:** the component has its own folder. Standard practice is to include a README.md as the primary entry, with the same content as below. Additional files in the directory may cover **implementation specifications** (e.g. a strict config YAML schema doc, API contract), sub-components, data flow, or design notes.

**Fields (in the component's primary doc):**

- **Responsibility** — What this component does in one or a few sentences.
- **Requirements** — Do not duplicate the O/R list. The engineering README is the single source of truth for which requirements are assigned to this component. The component doc may include a one-line reference (e.g. "Implements the requirements assigned to this component in the [engineering README](../README.md).").
- **Interfaces** — Key boundaries: what this component consumes (e.g. config, prompt buffer) and produces (e.g. exit code, iteration outcome), and which other components it calls or is called by.
- **Implementation spec** *(or linked doc)* — The authoritative spec implementers build from: e.g. the canonical config YAML schema (all keys, types, nesting, validation), API contracts, state machines, protocols. This is where the "hard" specification lives; product requirements state intent and may reference this doc.

Optional: data flow, invariants, or notes that help implementers place code correctly. Intent (who, what, why) stays in product; implementation specs stay in engineering.

## Rules

### Naming

- **Product:** Outcomes `O<n>-<slug>/`; requirements `R<n>-<slug>.md` within their outcome. The numeric part `<n>` is **zero-padded to three digits** (e.g. `O001`, `O002`, `R001`, `R002`) so that directories and files sort correctly in the filesystem. Numbered for stable reference, slug for readability. `R001` in O001 and `R001` in O002 are different requirements.
- **Engineering:** Component names are lowercase, hyphenated if multi-word (e.g. `run-loop.md`, `config.md`). Slugs may change; product IDs (O001, R002) are stable.

### Traceability

- **Product:** Every requirement declares its parent outcome. Every outcome README lists its requirements. The product index lists all outcomes. Product is the single source of truth for intent (who, what, why).
- **Engineering:** Requirement assignments (component → O/R IDs) live only in the engineering README. Component docs reference that README and do not duplicate the list. Engineering is the single source of truth for implementation (which component implements which requirements, and the implementation specs — schemas, interfaces, protocols — that each component adheres to). When an implementation spec changes in a way that affects user-facing contract, the product requirement is updated.

### Lifecycle (Product)

- New outcomes are added when a new problem is identified for a user segment (new who/what/why).
- New requirements are added under existing outcomes when a new capability is needed.
- Requirements without outcomes are removed or reassigned.
- Outcomes without requirements are aspirational — they need decomposition before they can be built.
- Status in the outcome README's requirement table:
  - **draft** — requirement identified, detail incomplete
  - **ready** — requirement complete for intent; implementation spec (if any) lives in engineering
  - **built** — implemented, not yet verified
  - **verified** — acceptance criteria confirmed to pass

### Lifecycle (Engineering)

- Add a component when the product has new requirements or an area that doesn't fit existing components; update the engineering README (names, one-liners, O/R assignments) and add the component doc or directory, including implementation specs (e.g. schema doc) as needed.
- Change or merge components when product or code structure warrants it; update the README and any affected component docs so assignments, interfaces, and implementation specs stay consistent.
- Requirement assignments live only in the engineering README; when requirements or component boundaries change, update the README first, then component docs and any schema/interface specs.
- No separate status for components; a component is "done" when its doc (and any implementation spec docs) exists and matches the README.
