# O003: Backend Agnosticism

## Who

Users and teams who use different AI CLIs (Claude, Kiro, Copilot, Cursor, or proprietary tools) and want one loop runner that works with whatever they have installed, without coupling to a single vendor.

## Statement

Any stdin-accepting AI CLI serves as the execution backend.

## Why it matters

The AI CLI landscape is fragmented. Teams use different tools and switch between them. A loop runner locked to one CLI is locked to one vendor. Ralph's value is the loop, not the AI. By treating the AI CLI as a pluggable process that accepts a prompt on stdin and produces output on stdout, Ralph works with whatever the user already has installed. The difficulty is that every AI CLI has its own flags for non-interactive mode, permission handling, and output format. Built-in aliases encode this knowledge so the user doesn't have to reverse-engineer each tool's invocation protocol.

## Verification

- User runs the same prompt with different AI CLIs (e.g. by switching alias or command). Both execute successfully using the same prompt file.
- User defines a custom alias in config for a proprietary AI CLI. Ralph resolves and executes it.
- User passes a direct command string on the command line, bypassing aliases. Ralph parses and executes the command.
- Ralph attempts to validate that the chosen AI command (or alias resolution) is installed or available before execution. If the command is missing or the alias cannot be resolved, the user gets a clear error before the loop (or review) starts, rather than a process start failure on the first iteration.
- The AI CLI process inherits the user's environment variables and working directory.
- Some AI CLIs write stdout in a structured format (e.g. newline-delimited JSON) rather than plain text. Ralph still works with them: the user configures success and failure signal strings to match content within that output, or uses a wrapper they write and an AI command alias that invokes it, so Ralph sees plain text for signal scanning.

## Non-outcomes

- Ralph does not guarantee that the AI CLI will succeed at runtime (e.g. API keys, network, or tool-specific errors). It validates that the command is present or resolvable before starting; runtime failures are still possible and are reported when they occur.
- Ralph does not manage AI CLI authentication, API keys, or session tokens.
- Ralph does not normalize output formats across different AI CLIs — that responsibility belongs to the command or wrapper the user invokes.
- Ralph does not write or provide wrappers that adapt structured output (e.g. JSON) to plain text. The user is responsible for writing such wrappers and for creating AI command aliases that call them if they want plain-text signal scanning. Otherwise Ralph would have to create and maintain adapters for how each built-in agent structures its output.
- Ralph does not execute commands through a shell. Commands are parsed and exec'd directly. Users needing pipes, redirects, or shell expansion must wrap their command in a script.
- Ralph does not provide vendor-specific adapters in core beyond built-in aliases (which invoke the raw CLI); adapting structured output is the user's responsibility.

## Risks

| Risk | Mitigating Requirement |
|------|------------------------|
| AI CLI invoked with wrong environment or working directory | [R002 — Inherit env and cwd](R002-inherit-env-and-cwd.md) |
| Structured or non–plain-text CLI output not detectable by the loop | [R003 — Structured output via signals or wrapper](R003-structured-output-via-signals-or-wrapper.md) |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R001](R001-invoke-ai-cli-stdin-capture-stdout.md) | The system invokes the user-chosen AI command (alias or direct) with the assembled prompt on stdin and captures stdout. | draft |
| [R002](R002-inherit-env-and-cwd.md) | The system inherits the user's environment and working directory when invoking the AI CLI. | draft |
| [R003](R003-structured-output-via-signals-or-wrapper.md) | The system works with AI CLIs that produce structured or non–plain-text output when the user configures signals or uses a wrapper. | draft |
