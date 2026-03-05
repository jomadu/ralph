# Building Intent

This document defines the model for documenting product intent. Intent documentation starts with outcomes and cascades into buildable specification.

## Structure

```
docs/
  intent/
    README.md                        # Index of all outcomes
    O<n>-<slug>/
      README.md                      # Outcome definition
      R<n>-<slug>.md                 # Requirement + specification
```

## Hierarchy

**Outcome** — A measurable change in the world when this product exists. Outcomes are user-facing and verifiable. They answer *"what is true when we succeed?"*

**Requirement** — A capability needed to deliver an outcome. Each requirement belongs to exactly one outcome. Requirements answer *"what must the system do?"*

**Specification** — Buildable detail embedded within a requirement. Schemas, formats, algorithms, edge cases. Specifications answer *"how exactly does it work?"*

Every specification traces to a requirement. Every requirement traces to an outcome. If a spec has no requirement above it, it's unjustified. If a requirement has no outcome above it, it's unjustified.

## Decomposition

Outcomes don't arrive pre-decomposed. The technique for deriving requirements from outcomes is the **why/how/how-else chain**, drawn from goal-oriented requirements engineering (GORE).

### The Chain

Starting from an outcome, ask three questions repeatedly:

1. **"How is this outcome achieved?"** — Each distinct answer is a candidate requirement. If the answer is still abstract, it's a sub-outcome that needs further decomposition. If it's concrete enough to specify and build, it's a requirement.

2. **"How else?"** — After the first answer, ask again. This prevents locking into a single design prematurely, surfaces alternatives, and catches missing capabilities. Stop when additional answers no longer feel distinct.

3. **"Why does this requirement exist?"** — For each candidate requirement, trace back up. It must connect to the outcome. If you can't articulate why, the requirement is either misplaced (belongs under a different outcome) or unjustified (shouldn't exist).

### Risk Analysis

For each outcome, ask **"what could prevent this from being true?"** Each answer is a risk. Risks surface requirements you'd otherwise discover during implementation. Each risk is documented in the outcome README with a reference to the requirement that mitigates it, providing traceability from risk to solution.

Every risk must either map to an existing requirement or surface a new one. If a risk has no mitigating requirement, the decomposition is incomplete.

Examples of risks and the requirements they produce:
- *"The AI CLI crashes mid-execution"* → a requirement for process crash recovery
- *"Both success and failure signals appear in the output"* → a requirement for signal precedence rules
- *"The user doesn't know which config value is active"* → a requirement for provenance tracking

### Completeness

An outcome is fully decomposed when:
- Every "how" question has been answered with a requirement
- Every risk you can identify maps to a mitigating requirement
- Every requirement traces back to the outcome via "why"
- You cannot describe a realistic failure scenario that no requirement addresses

Decomposition is a judgment call, not a formula. The goal is sufficiency for building, not exhaustive enumeration.

## Consistency

Internal consistency is a first-class concern at every step. Each layer of the intent tree introduces more documents and more surface area for contradiction. A single index page is trivially self-consistent. Multiple outcome files can silently diverge from each other and from the index. Requirements across outcomes can conflict. Specifications can prescribe incompatible behaviors.

At every step boundary, before locking, perform a consistency review across two dimensions:

- **Vertical consistency** — Does each document agree with the layer above it? Outcome detail files must match the index. Requirements must match their parent outcome. Specifications must match their requirement.
- **Horizontal consistency** — Do documents at the same level agree with each other? Outcome files must not make contradictory claims or assume incompatible models of the system. Requirements across different outcomes must not prescribe conflicting behaviors. Specifications must not define mechanisms that cannot coexist.

The expansion from one file to many is where inconsistency enters. When the product is a single index page, contradictions are visible on sight. The moment each outcome gets its own file, contradictions hide — one outcome's framing can drift from another's, and no single document reveals the conflict. The same is true when requirements fan out across outcomes. Each step's review must deliberately reunify the separate documents and read them as a set, not just individually.

Inconsistency at any layer invalidates everything below it. Catching it early is cheap; catching it in specifications is expensive; catching it in implementation is worse.

## Phased Development

Build the intent tree in five steps. Each step produces artifacts, reviews them, and locks them before the next step begins. The pattern is **compress, then expand** — write one-liners first, then expand into full documents. This happens twice: once for outcomes, once for requirements. Compressed forms are cheap to review for consistency because everything fits on one page. Expanded forms are reviewed against the locked compression above them.

### Step 1: Outcome Index

Write the root `README.md`. One table — each outcome as a one-line statement with its verification criteria. No directories, no detail. This is the product on a single page.

**Review:**
- Each outcome is a present-tense assertion about the world, not a feature description
- Outcomes don't overlap — if two outcomes could share a requirement, they may be the same outcome
- Verification criteria are observable by a user, not by a test suite
- The set of outcomes is complete — together they describe the whole product
- The set is minimal — removing any outcome would leave a gap

Lock the index before expanding. If the one-liners aren't right, the detail won't save them.

### Step 2: Outcome Detail

Expand each outcome into its own directory and README. Statement, why it matters, full verification, non-outcomes. These files define the problem space. They do not reference requirements or risks — those come later.

**Review:**
- Each outcome detail is consistent with its one-liner in the index — if they diverge, fix the index first
- Non-outcomes are clear enough that someone could push back on scope and point to this list
- Read all outcome files as a set: no two outcomes make contradictory claims, imply overlapping scope, or assume incompatible models of the user, the system, or the domain

Lock outcome detail before proceeding. Changes after this point ripple through everything below.

### Step 3: Requirement Index

Decompose each outcome into requirements using the why/how/how-else chain and risk analysis (see Decomposition). For each outcome, append two tables to its README: a risks table mapping each risk to its mitigating requirement, and a requirements table with one-line summaries. No requirement documents yet.

This is the compressed form of the requirements layer. All requirement one-liners and risk mappings across all outcomes should be reviewable as a set, the same way the outcome index was.

**Review:**
- Every requirement traces to exactly one outcome
- Requirements are capabilities ("the system detects X"), not implementations ("use a regex to scan for X")
- No requirement is redundant with another under the same outcome
- Every risk maps to a mitigating requirement
- The set of requirements under each outcome is sufficient — you can't describe a realistic failure that nothing addresses
- Read all requirement one-liners across all outcomes as a single set: no two requirements prescribe contradictory capabilities or overlap in scope

Lock requirement one-liners before expanding.

### Step 4: Requirement Detail

Expand each requirement into its own file within its outcome directory. Requirement statement, acceptance criteria. No specifications yet — stay at the "what," not the "how."

**Review:**
- Each requirement document is consistent with its one-liner in the outcome README — if they diverge, fix the one-liner first
- Acceptance criteria are concrete and testable
- Read all requirement documents across all outcomes as a single set: no two requirements prescribe contradictory behaviors, make incompatible assumptions, or define the same concept differently
- Requirements remain vertically consistent with their parent outcome detail — they do not introduce scope, assumptions, or framing that the outcome doesn't support

Lock requirements before proceeding. Specification changes are cheap; requirement changes are not.

### Step 5: Specifications

Fill in the specification section of each requirement. Schemas, formats, algorithms, edge cases, error handling. This is where implementation detail lives.

**Review:**
- An engineer or AI agent can build from this specification without asking clarifying questions
- Edge cases are enumerated, not hand-waved
- The specification doesn't exceed what the requirement asks for (gold-plating)
- Read all specifications as a single set: no specification prescribes behavior, formats, schemas, or mechanisms that conflict with any other specification
- Each specification remains vertically consistent with its parent requirement — it implements what the requirement asks for, nothing more and nothing less

### Working Sessions

This document is a methodology for generating and reviewing the intent tree. It is not a dependency for working within the tree once it exists. The intent documents are self-describing — an agent or engineer reading the outcome index, an outcome README, and its requirement files has everything needed to do the work. The templates, phasing, and review criteria are construction scaffolding, not load-bearing walls.

As the tree grows from outcomes to requirements to specifications, the total context exceeds what a single working session can hold. The intent structure is designed to be chunked by outcome. A session writing or reviewing specifications should load:

1. The outcome index — the full product on one page
2. The single outcome README being worked on — its risks and requirement one-liners
3. The requirement files under that outcome only

Other outcomes are not needed and dilute focus. One session per outcome, then a final session for cross-outcome consistency review.

The same pattern applies at any step where the volume of documents exceeds a single session's capacity. At Step 4 with many requirements, chunk by outcome. At Step 3 with many outcomes, chunk by outcome group. The compressed forms (outcome index, requirement one-liners) always fit in a single session and serve as the anchor for each focused session.

#### Context by Step

Each step requires different context. Earlier steps are lightweight — the tree barely exists. Later steps require tracing a path from the root down to the specific artifact being written. The principle: **load the locked ancestors of whatever you're creating, plus siblings for horizontal consistency.**

**Step 1 — Outcome Index:**
- This methodology document
- The repository's main README (product context)

Nothing else exists yet. The session's job is to compress the entire product into one table.

**Step 2 — Outcome Detail:**
- This methodology document
- The locked outcome index (`docs/intent/README.md`)
- All outcome READMEs written so far in this step (for horizontal consistency)

If the number of outcomes is small enough to hold in one session, write them all and review as a set. If not, write one at a time, then run a dedicated consistency session that loads all outcome READMEs together.

**Step 3 — Requirement Index:**
- This methodology document
- The outcome index
- The outcome README being decomposed

Work one outcome at a time. After all outcomes have requirement one-liners and risk tables, run a consistency session loading all outcome READMEs to review requirement one-liners across the full product.

**Step 4 — Requirement Detail:**
- The outcome index (product-level anchor)
- The parent outcome README (requirement one-liners and risks — your anchor)
- All requirement files under that outcome (for horizontal consistency within the outcome)

The methodology document is optional here — the locked artifacts above provide the structural guidance. Work one outcome's requirements per session. After all requirement files exist, run a cross-outcome consistency session loading requirement files that are most likely to interact.

**Step 5 — Specifications:**
- The outcome index
- The parent outcome README
- The requirement file being specified
- Related requirement files under the same outcome (for format and schema consistency)

By this step, context is narrowest: you're filling in buildable detail for a single requirement, anchored by the locked layers above it. If specifications across different outcomes must agree on shared formats or schemas, run a focused consistency session loading just those specific requirement files.

**Consistency review sessions** at any step follow a different pattern — they load documents *across* outcomes rather than *down* into one. Load the outcome index plus all documents at the layer being reviewed. These sessions read broadly rather than deeply.

### Why This Order Matters

The five steps alternate between compression and expansion. Compressed forms (outcome one-liners, requirement one-liners) are easy to hold in your head and cheap to review for consistency. Expanded forms (outcome detail, requirement documents, specifications) add depth but also surface area for contradiction. By locking the compressed form before expanding, you ensure that the detail is anchored to a reviewed, consistent summary. If you write all layers at once, errors in the upper layers silently propagate downward and get baked into detail that feels authoritative but is wrong.

## Index — `docs/intent/README.md`

Link each outcome ID to that outcome's README. From the index, use `./O<n>-<slug>/README.md`.

```markdown
# Intent

## Outcomes

| ID | Outcome | Verification |
|----|---------|--------------|
| [O1](./O1-<slug>/README.md) | Statement of what is true when this outcome is achieved | How a user proves it |
| [O2](./O2-<slug>/README.md) | ... | ... |
```

## Outcome — `O<n>-<slug>/README.md`

Each outcome directory has a README that fully defines the outcome. Link to requirement docs with `R<n>-<slug>.md`.

### Fields

**Statement** — One sentence. What is true when this outcome is delivered. Written as a present-tense assertion, not a feature description.

> *"An AI-driven task reaches a verified completion state through iterative execution."*
>
> Not: *"The system provides an iteration loop with signal scanning."*

**Why it matters** — The pain without this outcome. What goes wrong today.

**Verification** — How a user (not a test suite) knows this outcome was delivered. Observable, demonstrable evidence.

**Non-outcomes** — What this outcome explicitly does not cover. Prevents scope creep and clarifies boundaries.

**Risks** — What could prevent this outcome from being true. Each risk maps to a mitigating requirement, providing traceability from risk to solution. Risks are discovered through risk analysis during decomposition (see Decomposition).

### Template

```markdown
# O<n>: <Title>

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
| What could prevent this outcome | [R<n> — <Title>](R<n>-<slug>.md) |
| ... | ... |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R1](R1-<slug>.md) | One-line summary | draft / ready / built / verified |
| [R2](R2-<slug>.md) | ... | ... |
```

The risks and requirements tables are appended in Step 3. During Steps 1 and 2, the outcome README ends after non-outcomes.

## Requirement — `R<n>-<slug>.md`

Each requirement file contains both the requirement statement and its buildable specification. One file, two sections — the "what" and the "how" stay together.

### Fields

**Outcome** — Which outcome this requirement serves (traceability upward).

**Requirement** — What the system must do. Written as a capability statement.

**Specification** — The buildable detail. Schemas, formats, algorithms, edge cases, error handling. An engineer or AI agent reads this section and implements from it. This section can use subsections, code blocks, tables, and diagrams as needed. It includes two standard subsections:

- **Edge cases** — A structured table of boundary conditions and the system's expected behavior for each. Forces enumeration rather than prose — every edge case is a row, not a paragraph.
- **Examples** — Concrete scenarios with input, expected output, and how to verify the result. Examples bridge the gap between abstract acceptance criteria and buildable understanding. An implementer can run the example and check their work.

**Acceptance criteria** — Concrete conditions that must be true for this requirement to be considered met.

**Dependencies** *(optional)* — Other requirements or system capabilities that must exist before this requirement can be implemented. Omit if the requirement is self-contained.

### Template

```markdown
# R<n>: <Title>

**Outcome:** O<n> — <Outcome title>

## Requirement

<What the system must do. Capability statement.>

## Specification

<Buildable detail. Schemas, formats, algorithms, edge cases.>

<This section can be as long as needed. Use subsections, code blocks, tables, diagrams.>

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| <Boundary condition> | <How the system responds> |
| <Error scenario> | <How the system responds> |

### Examples

#### <Scenario name>

**Input:**
<Input data, command, or setup>

**Expected output:**
<Output data, behavior, or result>

**Verification:**
- <How to verify the outcome>

## Acceptance criteria

- [ ] <Concrete, testable condition>
- [ ] <...>

## Dependencies

- <R<n> — Requirement that must exist first>
- <Other prerequisite>
```

## Rules

### Naming

- Outcomes: `O<n>-<slug>/` — numbered for stable reference, slug for readability.
- Requirements: `R<n>-<slug>.md` — numbered within their outcome. `R1` in `O1` and `R1` in `O2` are different requirements.
- Slugs are lowercase, hyphenated, descriptive. They may change; IDs (`O1`, `R2`) are stable.

### Traceability

- Every requirement file declares its parent outcome.
- Every outcome README lists its requirements.
- The root index lists all outcomes.
- If you can't trace a specification back to an outcome, question whether it belongs.

### Lifecycle

- New outcomes are added when a new user-facing problem is identified.
- New requirements are added under existing outcomes when a new capability is needed.
- Requirements without outcomes are removed or reassigned.
- Outcomes without requirements are aspirational — they need decomposition before they can be built.
- Status tracking lives in the outcome README's requirement table:
  - `draft` — requirement identified, specification incomplete
  - `ready` — fully specified, can be built from
  - `built` — implemented, not yet verified
  - `verified` — acceptance criteria confirmed to pass
