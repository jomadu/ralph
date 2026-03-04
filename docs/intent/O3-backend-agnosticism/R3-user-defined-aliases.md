# R3: User-Defined Command Aliases

**Outcome:** O3 — Backend Agnosticism

## Requirement

The system allows users to define custom AI command aliases in config files. User-defined aliases merge with built-in aliases, and user-defined aliases with the same name as a built-in alias override the built-in. This supports proprietary, internal, or newly released AI CLIs that Ralph doesn't ship aliases for.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] Users can define aliases under the `ai_cmd_aliases` section in config files
- [ ] User-defined aliases merge with built-in aliases — both are available simultaneously
- [ ] A user-defined alias with the same name as a built-in alias overrides the built-in
- [ ] Aliases defined in workspace config override aliases with the same name defined in global config
- [ ] User-defined alias values are parsed using the same shell-style command parsing as direct commands
