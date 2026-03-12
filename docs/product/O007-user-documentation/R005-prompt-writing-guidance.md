# R005: Prompt-writing guidance

**Outcome:** O007 — User Documentation

## Requirement

Users can find guidance on how to write a well-formed Ralph prompt, using the same criteria that `ralph review` uses when evaluating prompts.

## Detail

Documentation includes a dedicated guide (e.g. "Writing Ralph prompts") that explains what makes a prompt work well with Ralph's execution model. The criteria in the guide are the same four dimensions used by the review component (O005/R007): signal and state, iteration awareness, scope and convergence, and subjective completion criteria. The guide is written for users who are new to writing prompts or who want to improve their prompts before running or reviewing. It does not replace or duplicate the evaluation-dimensions requirement; it presents those dimensions in user-friendly language with examples (do/don't or strong vs weak) so users can self-serve. The doc may optionally include a minimal "good enough" prompt template. The guide must explicitly state that these are the same criteria `ralph review` uses.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User has never written a Ralph prompt | Guide is readable without prior knowledge of the codebase; dimensions are explained in plain language. |
| User wants to check their prompt before running | User can use the guide to self-assess; they may also run `ralph review` for AI-evaluated feedback. |
| Evaluation dimensions (O005/R007) change | The writing guide is updated so it continues to match review criteria. |

### Examples

#### User looks up how to write a prompt

**Input:** User is about to run Ralph for the first time and does not know how to structure a prompt.

**Expected output:** User finds the "Writing Ralph prompts" guide (or equivalent title). The guide describes the four dimensions (signal and state, iteration awareness, scope and convergence, subjective completion criteria) in user terms and states that these are the same criteria `ralph review` uses. User can read it and apply the guidance to draft or revise a prompt.

**Verification:** User can locate the guide, understand the criteria, and see how they map to writing a prompt; the criteria align with O005/R007.

#### User wants to know what review will check

**Input:** User plans to run `ralph review` and wants to know what the reviewer evaluates.

**Expected output:** The writing guide (and/or review help) states that review evaluates prompts along the same four dimensions described in the guide. User can prepare their prompt accordingly.

**Verification:** No mismatch between what the guide recommends and what review evaluates.

## Acceptance criteria

- [ ] A user-facing guide exists (e.g. "Writing Ralph prompts") that explains how to write a well-formed Ralph prompt.
- [ ] The guide uses the same four dimensions as O005/R007: signal and state, iteration awareness, scope and convergence, subjective completion criteria.
- [ ] The guide explicitly states that these are the same criteria `ralph review` uses.
- [ ] The guide is written in user-friendly language and includes examples (e.g. do/don't or strong vs weak) where appropriate.
- [ ] When O005/R007 is updated, the guide is updated so it stays aligned with review criteria.

## Dependencies

- [O005/R007](../O005-prompt-review/R007-evaluation-dimensions.md) — The four evaluation dimensions are the single source of truth; the writing guide presents them for users.
