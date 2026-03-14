# R004: Full non-interactive config

**Outcome:** O010 — Automation

## Requirement

Ralph allows full non-interactive configuration of behavior (e.g. timeouts, iteration limits, prompts, AI commands) via config and environment variables so headless and scripted use require no interactive setup.

## Detail

Users running Ralph from scripts or CI must be able to configure everything needed for the run or review without any interactive prompts or setup steps. Configuration is supplied via config file(s) and/or environment variables (per O002 config layer resolution). The set of behaviors that can be configured includes at least: prompt source (alias or path), AI command or alias, iteration limits (max iterations), failure threshold, timeouts (if applicable), and for review: report path, apply behavior, and any options that would otherwise require confirmation. Headless and scripted invocations can therefore run with no user input: config (and optionally env) fully specify behavior. Interactive flows (e.g. apply with confirmation) remain available when not using the non-interactive options; automation uses config and env so that no interactive setup is required.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| All behavior set via config file | Loop or reviewer runs with no interactive prompts; no config wizard or prompts. |
| Override via environment variables where supported | Env takes precedence per O002; script can override without editing config file. |
| Missing required config (e.g. no AI command) | Ralph fails with clear error (O004) and documented exit code; does not prompt for input. |
| Config file path via env or CLI | User can point to config in a headless way (e.g. documented env or CLI option for config path). |
| Prompt from alias or file path in config | No prompt to "choose" prompt; alias or path is configured. |
| Apply with non-interactive option in config or env | Apply can be enabled without confirmation when so configured. |

### Examples

#### Full config for loop

**Input:** Config file specifies: prompt alias, AI command alias, max iterations, failure threshold. User runs the run command with that config in CI; no TTY.

**Expected output:** Loop runs with prompt and AI command from config; stops at configured iteration or failure limits; exits with the documented success or failure/exhaustion code. No prompts for config or confirmation.

**Verification:** No interactive prompt; behavior matches config; exit code is documented.

#### Review and apply in CI

**Input:** Config specifies report path and apply-without-confirmation (or equivalent). CI runs the review command with that config.

**Expected output:** Review runs; report written to configured path; revision applied without confirmation prompt. Exit code is one of the documented set (R002).

**Verification:** No confirmation prompt; config fully drives behavior; suitable for headless run.

#### Missing AI command

**Input:** Config has no AI command (and none in env); script runs the run command.

**Expected output:** Ralph exits with clear error that AI command is missing (O004/R001); the documented failure code. No prompt to "enter" command.

**Verification:** Failure is explicit and non-interactive; script can detect and handle.

## Acceptance criteria

- [ ] Behavior required for loop and reviewer (prompt source, AI command, iteration limits, failure threshold, timeouts where applicable, report path, apply behavior) can be set via config file and/or environment variables.
- [ ] When config (and env) are fully specified, loop and reviewer can run to completion without any interactive prompts or setup steps.
- [ ] When required config is missing, Ralph fails with a clear error and documented exit code instead of prompting for input.
- [ ] The mechanism to supply config (config file path via env or CLI) is usable in a headless way (no interactive choice of config path required for the common case).

## Dependencies

- O002 — Config layer resolution (config file, env) defines how config is merged; this requirement ensures the set of options available there is sufficient for full non-interactive use.
