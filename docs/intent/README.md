# Intent

This directory holds the product intent tree for Ralph: outcomes, requirements, and specifications. Intent starts with outcomes and cascades into buildable specification. The methodology (phased development, decomposition, review criteria) is defined in [building-intent.md](../../building-intent.md).

## Structure

```
docs/intent/
  README.md                    ← you are here (outcome index)
  O<n>-<slug>/
    README.md                  Outcome definition, risks, requirement one-liners
    R<n>-<slug>.md             Requirement + specification
```

- **Outcome** — A measurable change in the world when this product exists. User-facing and verifiable.
- **Requirement** — A capability needed to deliver an outcome. Each requirement belongs to one outcome.
- **Specification** — Buildable detail inside a requirement: schemas, formats, algorithms, edge cases.

Every specification traces to a requirement; every requirement traces to an outcome.

## Outcomes

| ID | Outcome | Verification |
|----|---------|--------------|
| [O1](./O1-iterative-completion/README.md) | An AI-driven task reaches verified completion through iterative execution | User runs `ralph run <alias>`, Ralph executes fresh AI processes across iterations, detects a success signal in the output, and exits 0 |
| [O2](./O2-configurable-behavior/README.md) | Loop execution adapts to the user's constraints without prompt modification | User changes iteration limits, failure thresholds, timeouts, and signal strings via config or CLI flags — the same prompt file produces different loop behavior |
| [O3](./O3-backend-agnosticism/README.md) | Any stdin-accepting AI CLI serves as the execution backend | User runs the same prompt with different AI CLIs by changing a config value or flag, and Ralph works with each |
| [O4](./O4-observability/README.md) | The user knows why the loop stopped and how it performed | Exit code distinguishes success (0), failure threshold (1), exhaustion (2), and interruption (130); iteration statistics are reported at completion |
