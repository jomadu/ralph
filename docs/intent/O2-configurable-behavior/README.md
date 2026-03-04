# O2: Configurable Behavior

## Statement

Loop execution adapts to the user's constraints without prompt modification.

## Why it matters

Different tasks need different loop parameters. A one-shot bootstrap needs 1 iteration with no preamble. A long build loop needs 20 iterations with a high failure threshold. A cautious exploration needs a short timeout and low threshold. Without external configuration, these differences live inside the prompt file — the user maintains variant copies of the same prompt for different environments, projects, or risk tolerances. Configuration separates loop behavior from prompt content so the same prompt file works across contexts.

## Verification

- User sets `default_max_iterations: 10` in `ralph-config.yml`, then overrides with `ralph run build -n 20` on the command line. The loop runs up to 20 iterations.
- User defines a prompt alias with `failure_threshold: 5` and custom signal strings. Those values take effect for that alias without affecting other aliases.
- User sets `RALPH_LOOP_ITERATION_TIMEOUT=60` in the environment. Ralph applies a 60-second timeout without any config file change.
- User runs `ralph list prompts` and sees all available prompt aliases with names and descriptions.
- User runs `ralph list aliases` and sees all available AI command aliases — both built-in and user-defined — with their resolved commands.

## Non-outcomes

- Ralph does not provide a GUI, interactive config editor, or `ralph config` subcommand. Configuration is files, env vars, and flags.
- Ralph does not support runtime config changes during loop execution. Config is resolved once at startup.
- Ralph does not validate prompt file content — only that the file exists and is readable.
- Ralph does not support config inheritance between prompt aliases. Each alias independently overrides the root `loop` section.

## Obstacles

| Obstacle | Mitigating Requirement |
|----------|----------------------|
| Multiple config layers set the same key and user doesn't know which value is active | R1 — Configuration provenance tracking |
| Config file has a typo in a key name and silently does nothing | R2 — Unknown key warnings |
| Config file has invalid values (negative iterations, empty signal) | R3 — Configuration validation at load time |
| Prompt source is missing, unreadable, or empty | R4 — Fail-fast on invalid prompt source |
| Missing config files cause startup errors | R5 — Silent skip for absent config files |
| User needs different signal strings, timeouts, or thresholds per prompt | R6 — Per-prompt loop setting overrides |
| User can't discover which prompts or AI command aliases are available | R7 — Resource listing command |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| R1 | Configuration provenance tracking | draft |
| R2 | Unknown key warnings | draft |
| R3 | Configuration validation at load time | draft |
| R4 | Fail-fast on invalid prompt source | draft |
| R5 | Silent skip for absent config files | draft |
| R6 | Per-prompt loop setting overrides | draft |
| R7 | Resource listing command | draft |
