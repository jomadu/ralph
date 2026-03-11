# R001: What and Why

**Outcome:** O008 — Discoverability

## Requirement

The user can find a short description of what the product is and why it exists.

## Detail

A new user or evaluator needs to quickly understand the product's value: what it does (e.g. loop runner for AI-driven tasks) and why it exists (e.g. manual read–judge–re-run replaced by automated iteration until a signal). This content is discoverable without reading the codebase — e.g. in the repository README, docs, or the product's top-level help. The description is short and sufficient to answer "what is this?" so the user can decide whether to try it.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User lands on repo root | A short "what" and "why" are visible (e.g. in README or linked doc). |
| User invokes the product's top-level help (e.g. help or list) | That entry point includes or points to a brief description of what the product does. |
| No prior knowledge of the codebase | User can find the description via documented entry points (README, docs index, or CLI). |
| Documentation is split across files | At least one primary entry (e.g. main README or docs landing) states what the product is and why it exists. |

### Examples

#### Read repository README

**Input:** New user opens the repository and reads the main README.

**Expected output:** README (or a clearly linked doc) contains a short description of the product (e.g. loop runner for AI-driven tasks) and why it exists (e.g. automating iteration until success/failure signal).

**Verification:** Reviewer can point to specific sentences that answer "what is this?" and "why does it exist?"

#### Help or list entry point

**Input:** User invokes the product's top-level help or a way to list prompts/subcommands (e.g. help or list).

**Expected output:** The output includes a brief description of what the product does, or a pointer (e.g. "See README" or "See docs/...") to where that description lives.

**Verification:** User can learn what the product is without opening the repo in a browser.

## Acceptance criteria

- [ ] A short description of what the product is (e.g. loop runner for AI-driven tasks) is discoverable from a documented entry point (e.g. repo README, docs, or the product's help).
- [ ] A short description of why the product exists (e.g. manual read–judge–re-run replaced by automated iteration) is discoverable from the same or linked content.
- [ ] A new user can find this content without reverse-engineering the codebase or guessing at URLs.
- [ ] The description is concise (no requirement for length; "short" means sufficient to answer "what is this?" and "why try it?").

## Dependencies

None.
