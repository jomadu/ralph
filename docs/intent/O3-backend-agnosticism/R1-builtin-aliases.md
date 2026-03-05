# R1: Built-in Command Aliases

**Outcome:** O3 — Backend Agnosticism

## Requirement

The system ships with built-in command aliases for known AI CLIs, encoding the correct flags and invocation protocols for non-interactive, stdin-based execution. Built-in aliases eliminate the need for users to reverse-engineer each tool's invocation requirements. For AI CLIs that emit non-standard output (e.g., structured JSON instead of plain text), the alias resolves to a wrapper script that normalizes the output.

## Specification

_To be specified._

### Edge cases

_To be specified._

### Examples

_To be specified._

## Acceptance criteria

- [ ] Built-in aliases include at minimum: claude, kiro, copilot, cursor-agent
- [ ] Each alias resolves to a complete command string with all necessary flags for non-interactive, stdin-based operation
- [ ] Built-in aliases are available without any user configuration
- [ ] User-defined aliases with the same name as a built-in alias override the built-in
- [ ] The cursor-agent alias resolves to a wrapper script ([`cursor-wrapper.sh`](../../../scripts/cursor-wrapper.sh)) that parses structured JSON output and emits plain text suitable for signal scanning
- [ ] If a wrapper script has a missing dependency (e.g., jq), it reports the missing dependency clearly

## Dependencies

_None identified._
