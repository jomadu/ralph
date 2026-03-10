# O010: Automation

## Who

Users who run Ralph from scripts, CI pipelines, or other automation and need reliable exit codes, non-interactive behavior, and stable interfaces.

## Statement

Users can run Ralph from scripts and CI.

## Why it matters

Ralph is not only for interactive use. Scripts and CI need to invoke the loop or the reviewer, interpret exit codes, and optionally consume machine-parseable output (e.g. review summary). If Ralph requires interactive confirmation, uses unclear exit codes, or changes behavior in ways that break scripts, automation fails. This outcome ensures Ralph is a good citizen in automated environments.

## Verification

- User runs the loop or the reviewer from a script without user input. The command completes (success or failure) and exits with a documented code; no interactive prompts block completion unless the user has explicitly chosen a flow that requires confirmation (e.g. applying a revision without a non-interactive option).
- Exit codes are documented and stable so scripts can branch on outcome: e.g. success, failure threshold or prompt errors, exhaustion or review/apply failure, interruption. The meaning of each code is consistent across releases within the compatibility contract.
- User runs the reviewer and obtains a machine-parseable summary (or report) so CI can gate on prompt quality without scraping free text.
- Environment variables and config allow full non-interactive configuration (e.g. timeouts, iteration limits, prompts, AI commands) so no interactive setup is required in headless environments.

## Non-outcomes

- Ralph does not provide a dedicated "CI mode" or separate binary; the same product is scriptable when used with the documented options and config.
- Ralph does not guarantee compatibility with every CI platform's quirks; the outcome is that the contract (exit codes, non-interactive behavior, parseable output where documented) supports scripting and CI.
- Optional interactive flows (e.g. apply with confirmation) remain available; automation uses the non-interactive variants where documented.
