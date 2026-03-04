# R8: Preamble Injection

**Outcome:** O1 — Iterative Completion

## Requirement

The system wraps the user's prompt with loop state — iteration number, iteration limit, and optional user-provided context — before piping the assembled content to the AI CLI. The preamble provides dynamic metadata that the static prompt file cannot know on its own. It is generated per iteration (iteration number changes), but the underlying prompt content is immutable across iterations.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] The preamble includes the current iteration number and the max iteration limit (or "unlimited" when in unlimited mode)
- [ ] When context is provided via --context flag(s), the preamble includes a CONTEXT section with the provided content
- [ ] When no context is provided, the CONTEXT section is omitted entirely
- [ ] Preamble injection is enabled by default
- [ ] Preamble can be disabled globally (loop.preamble: false) or per-prompt (prompts.<name>.loop.preamble: false)
- [ ] When preamble is disabled, the prompt file content is piped directly to the AI CLI without any wrapping
