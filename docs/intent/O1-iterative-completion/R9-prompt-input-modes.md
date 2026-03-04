# R9: Prompt Input Modes

**Outcome:** O1 — Iterative Completion

## Requirement

The system accepts prompts from multiple sources — a configured alias, a direct file path flag, or stdin — loading the prompt content once at loop start and reusing it immutably across all iterations. This ensures consistent behavior regardless of filesystem changes during execution and supports one-off usage without config file setup.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] `ralph run <alias>` reads the prompt file mapped to the alias in the resolved config
- [ ] `ralph run -f <path>` reads the specified file directly, without requiring an alias in config
- [ ] `cat prompt.md | ralph run` reads the prompt from stdin when no alias or file flag is provided
- [ ] In all modes, the prompt content is read once at loop start and buffered in memory
- [ ] The same buffered content is used for every iteration — changes to the prompt file on disk after loop start do not affect subsequent iterations
- [ ] If the prompt source is empty, unreadable, or missing, Ralph exits with an error before starting the loop
