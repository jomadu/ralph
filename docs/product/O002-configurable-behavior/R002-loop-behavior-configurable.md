# R002: Loop Behavior Configurable

**Outcome:** O002 — Configurable Behavior

## Requirement

The system allows loop behavior (iterations, failure threshold, timeout, signals, signal precedence mode, preamble, AI command, streaming, log level) to be configured at root and per prompt. Root settings apply to all prompts unless overridden per prompt; resolution follows the layer order in R001.

## Detail

Loop execution adapts to the user's constraints without changing the prompt file. The following are configurable:

- **Iterations** — Maximum iterations; iteration mode (bounded vs unlimited).
- **Failure threshold** — Consecutive failures before the loop exits.
- **Timeout** — Per-iteration timeout.
- **Output limits** — Where applicable (e.g. output truncation).
- **Success and failure signal strings** — Strings the system looks for to classify iteration success or failure.
- **Signal precedence mode** — Static default vs optional AI-interpreted when both success and failure signals appear.
- **Preamble** — Whether to inject a preamble (e.g. for bootstrap vs long loops).
- **AI command / alias** — Which AI CLI command or alias to use.
- **Streaming** — Whether to stream AI output to the terminal.
- **Log level** — Log verbosity.

Configuration can be set at the root (default for all prompts), per prompt in config files (overrides root for that prompt), and overridable by environment variables and command-line options for a run. Per-prompt overrides apply when running or listing that prompt; environment and command-line options still override for that run (R001).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Root sets the iteration limit; prompt does not override | All prompts use the root iteration limit. |
| Root sets the iteration limit; prompt overrides with different value | When running that prompt, prompt value used (unless environment or command-line override). |
| Neither root nor prompt sets a loop setting | Default value for that setting applies. |
| Invalid value for a setting (e.g. negative iteration limit) | System reports an error or uses a defined fallback; behavior is documented. |
| Timeout zero or disabled | Defined behavior (e.g. no per-iteration timeout). |
| Empty or missing signal strings | Defined behavior (e.g. no signal-based success/failure detection for that run). |

### Examples

#### Root iteration limit

**Input:** Config root sets the configured iteration limit to 20. Run any prompt without override.

**Expected output:** Loop runs until success, 20 iterations, or failure threshold.

**Verification:** Run and observe loop stops at 20 iterations if no success or failure before then.

#### Per-prompt failure threshold

**Input:** Config root sets the failure threshold to 5. Prompt "cautious" has an override setting that sets that threshold to 1. Run prompt "cautious".

**Expected output:** Loop exits after 1 consecutive failure when running "cautious".

**Verification:** Run "cautious" and trigger one failure; loop exits. Run another prompt without override; it uses the root value of 5.

#### Command-line overrides resolved config

**Input:** Config sets the per-iteration timeout to 300. User runs with the documented timeout option set to 60 for that run.

**Expected output:** That run uses 60-second per-iteration timeout.

**Verification:** Run and observe iteration timeout at 60 seconds.

## Acceptance criteria

- [ ] All listed loop behavior settings (iterations, failure threshold, timeout, output limits, success/failure signals, signal precedence mode, preamble, AI command, streaming, log level) are configurable at the root level.
- [ ] The same settings can be overridden per prompt in config; when running that prompt, prompt overrides apply over root.
- [ ] Environment variables and command-line options can override resolved loop settings for that run.
- [ ] Default values exist for all loop settings so the tool works without a config file.
- [ ] Invalid or out-of-range values are handled with a defined behavior (error or documented fallback).

## Dependencies

- R001 — Config layer resolution (loop settings are resolved via the same layers and override order).
