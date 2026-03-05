# Intent Model

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

### Obstacles

For each outcome, ask **"what could prevent this from being true?"** Each answer is an obstacle. Obstacles are the primary tool for finding requirements you'd otherwise discover during implementation.

Every obstacle must either:
- Map to an existing requirement that mitigates it, or
- Surface a new requirement that needs to be added

If an obstacle has no mitigating requirement, the decomposition is incomplete.

Examples of obstacles for a loop runner:
- *"The AI CLI crashes mid-execution"* — surfaces a requirement for process crash recovery
- *"Both success and failure signals appear in the output"* — surfaces a requirement for signal precedence rules
- *"The user doesn't know which config value is active"* — surfaces a requirement for provenance tracking

### Completeness

An outcome is fully decomposed when:
- Every "how" question has been answered with a requirement
- Every obstacle has a mitigating requirement
- Every requirement traces back to the outcome via "why"
- You cannot describe a realistic failure scenario that no requirement addresses

Decomposition is a judgment call, not a formula. The goal is sufficiency for building, not exhaustive enumeration.

## Consistency

Internal consistency is a first-class concern at every phase. Each layer of the intent tree introduces more documents and more surface area for contradiction. A single index page is trivially self-consistent. Multiple outcome files can silently diverge from each other and from the index. Requirements across outcomes can conflict. Specifications can prescribe incompatible behaviors.

At every phase boundary, before locking, perform a consistency review across two dimensions:

- **Vertical consistency** — Does each document agree with the layer above it? Outcome detail files must match the index. Requirements must match their parent outcome. Specifications must match their requirement.
- **Horizontal consistency** — Do documents at the same level agree with each other? Outcome files must not make contradictory claims or assume incompatible models of the system. Requirements across different outcomes must not prescribe conflicting behaviors. Specifications must not define mechanisms that cannot coexist.

The expansion from one file to many is where inconsistency enters. When the product is a single index page, contradictions are visible on sight. The moment each outcome gets its own file, contradictions hide — one outcome's framing can drift from another's, and no single document reveals the conflict. The same is true when requirements fan out across outcomes. Each phase review must deliberately reunify the separate documents and read them as a set, not just individually.

Inconsistency at any layer invalidates everything below it. Catching it early is cheap; catching it in specifications is expensive; catching it in implementation is worse.

## Phased Development

Build the intent tree in three phases. Complete and review each phase before starting the next. Ambiguity compounds across layers — a vague outcome produces wrong requirements, which produce wrong specifications. Each phase is an annealing step: apply heat (scrutiny), let it settle (review), then lock it in before building the next layer on top.

### Phase 1a: Outcome Index

Write the root `README.md` first. One table — each outcome as a one-line statement with its verification criteria. No directories, no detail. This is the product on a single page.

**Review criteria:**
- Each outcome is a present-tense assertion about the world, not a feature description
- Outcomes don't overlap — if two outcomes could share a requirement, they may be the same outcome
- Verification criteria are observable by a user, not by a test suite
- The set of outcomes is complete — together they describe the whole product
- The set is minimal — removing any outcome would leave a gap

Lock the index before expanding. If the one-liners aren't right, the detail won't save them.

### Phase 1b: Outcome Detail

Expand each outcome into its own directory and README. Why it matters, full verification, non-outcomes, obstacles. The index is the contract; the detail is the justification.

**Review criteria:**
- The detail is consistent with the one-liner in the index — if they diverge, fix the index first
- Obstacles are realistic, not hypothetical
- Non-outcomes are clear enough that someone could push back on scope and point to this list

**Consistency review:** Read all outcome files as a set. Check that no two outcomes make contradictory claims, imply overlapping scope, or assume incompatible models of the user, the system, or the domain. Each outcome was written in isolation — this review is the first time they are read together, and it's where silent divergence surfaces.

Lock outcome detail before proceeding. Changes after this point ripple through everything below.

### Phase 2: Requirements

Decompose each outcome into requirements using the why/how/how-else chain. Map every obstacle to a mitigating requirement. Write the requirement statement and acceptance criteria. No specifications yet — stay at the "what," not the "how."

**Review criteria:**
- Every requirement traces to exactly one outcome
- Every obstacle has a mitigating requirement
- Requirements are capabilities ("the system detects X"), not implementations ("use a regex to scan for X")
- Acceptance criteria are concrete and testable
- No requirement is redundant with another under the same outcome
- The set of requirements under each outcome is sufficient — you can't describe a realistic failure that nothing addresses

**Consistency review:** Read all requirements across all outcomes as a single set. Check that no two requirements — even under different outcomes — prescribe contradictory behaviors, make incompatible assumptions, or define the same concept differently. A requirement written under O1 may silently conflict with one under O3; neither file reveals this on its own. This review must also verify that requirements remain vertically consistent with their parent outcome detail — they should not introduce scope, assumptions, or framing that the outcome doesn't support.

Lock requirements before proceeding. Specification changes are cheap; requirement changes are not.

### Phase 3: Specifications

Fill in the specification section of each requirement. Schemas, formats, algorithms, edge cases, error handling. This is where implementation detail lives.

**Review criteria:**
- An engineer or AI agent can build from this specification without asking clarifying questions
- Edge cases are enumerated, not hand-waved
- The specification doesn't exceed what the requirement asks for (gold-plating)

**Consistency review:** Read all specifications as a single set. Check that no specification prescribes behavior, formats, schemas, or mechanisms that conflict with any other specification — including those under different outcomes. Specifications are the most detailed layer and the most likely to introduce subtle incompatibilities (e.g., two specs defining the same data structure differently, or assuming contradictory ordering guarantees). Each specification must also remain vertically consistent with its parent requirement — it should implement what the requirement asks for, nothing more and nothing less.

### Why This Order Matters

Each layer is a lossy compression of the intent above it. Requirements are a compression of outcomes; specifications are a compression of requirements. If you write all three at once, errors in the upper layers silently propagate downward and get baked into detail that feels authoritative but is wrong. Phasing forces you to get the intent right before you commit to the mechanism.

## Index — `docs/intent/README.md`

The root index is a single page that maps the entire product. One entry per outcome: a one-line statement and its verification criteria.

```markdown
# Intent

## Outcomes

| ID | Outcome | Verification |
|----|---------|--------------|
| O1 | Statement of what is true when this outcome is achieved | How a user proves it |
| O2 | ... | ... |
```

## Outcome — `O<n>-<slug>/README.md`

Each outcome directory has a README that fully defines the outcome.

### Fields

**Statement** — One sentence. What is true when this outcome is delivered. Written as a present-tense assertion, not a feature description.

> *"An AI-driven task reaches a verified completion state through iterative execution."*
>
> Not: *"The system provides an iteration loop with signal scanning."*

**Why it matters** — The pain without this outcome. What goes wrong today.

**Verification** — How a user (not a test suite) knows this outcome was delivered. Observable, demonstrable evidence.

**Non-outcomes** — What this outcome explicitly does not cover. Prevents scope creep and clarifies boundaries.

**Obstacles** — What could prevent this outcome from being true. Each obstacle must map to a mitigating requirement. Obstacles are discovered by asking *"what could go wrong?"* and are the primary mechanism for surfacing missing requirements.

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

## Obstacles

| Obstacle | Mitigating Requirement |
|----------|----------------------|
| What could prevent this outcome | R<n> — <Title> |
| ... | ... |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| R1 | One-line summary | draft / ready / built / verified |
| R2 | ... | ... |
```

## Requirement — `R<n>-<slug>.md`

Each requirement file contains both the requirement statement and its buildable specification. One file, two sections — the "what" and the "how" stay together.

### Fields

**Outcome** — Which outcome this requirement serves (traceability upward).

**Requirement** — What the system must do. Written as a capability statement.

**Specification** — The buildable detail. Schemas, formats, algorithms, edge cases, error handling. An engineer or AI agent reads this section and implements from it.

**Acceptance criteria** — Concrete conditions that must be true for this requirement to be considered met.

### Template

```markdown
# R<n>: <Title>

**Outcome:** O<n> — <Outcome title>

## Requirement

<What the system must do. Capability statement.>

## Specification

<Buildable detail. Schemas, formats, algorithms, edge cases.>

<This section can be as long as needed. Use subsections, code blocks, tables, diagrams.>

## Acceptance criteria

- [ ] <Concrete, testable condition>
- [ ] <...>
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
