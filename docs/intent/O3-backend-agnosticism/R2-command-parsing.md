# R2: Shell-Style Command Parsing

**Outcome:** O3 — Backend Agnosticism

## Requirement

The system parses command strings using shell-style quoting rules — supporting quoted arguments and escaped characters — without invoking a shell. Commands are exec'd directly as processes, not passed through sh/bash.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] Command strings with double-quoted arguments are parsed correctly (e.g., `--model "claude-3-5-sonnet"` produces two tokens: `--model` and `claude-3-5-sonnet`)
- [ ] Command strings with single-quoted arguments are parsed correctly
- [ ] Escaped characters within quoted strings are handled
- [ ] The parsed command is executed directly via exec, not through a shell
- [ ] Shell features (pipes, redirects, glob expansion, variable substitution) are not interpreted and do not cause silent misbehavior
