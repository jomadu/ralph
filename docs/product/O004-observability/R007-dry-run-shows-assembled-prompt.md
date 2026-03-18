# R007: Dry-run shows assembled prompt

**Outcome:** O004 — Observability

## Requirement

The system supports a dry-run mode that prints the assembled prompt without invoking the AI.

## Detail

The user can verify what prompt will be sent to the AI: a single CONTEXT section (Ralph loop description and iteration when preamble is enabled, plus any invoker-provided context with an explicit label) and the INSTRUCTIONS section, without running the loop or spawning an AI process. This supports understanding why the loop or review would behave as it does and debugging config or prompt assembly. No AI process is started in dry-run.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Dry-run with default config | CONTEXT section (Ralph loop + iteration) and INSTRUCTIONS section printed to stdout; no AI invocation |
| Dry-run with preamble disabled | CONTEXT section only if invoker passed -c (with invoker label); then INSTRUCTIONS section |
| Dry-run with -c "Hello world" | CONTEXT section includes Ralph loop/iteration (if enabled) and "Context provided by the invoker of this Ralph run:" followed by the given text; then INSTRUCTIONS section. No duplicate CONTEXT title in the body. Section boundaries use `# --- CONTEXT ---` / `# --- INSTRUCTIONS ---` (see run-loop engineering spec). |
| Dry-run and log level | The system may emit minimal log context (e.g. that dry-run is active); primary output is the assembled prompt |
| Dry-run with missing AI command | Behavior may follow R001 (error before run); dry-run may still show assembled prompt without needing a valid AI command |

### Examples

#### Dry-run output

**Input:** User invokes the run command in dry-run mode with a prompt file; preamble is enabled (default).

**Expected output:** The fully assembled prompt (CONTEXT section with Ralph loop/iteration and optionally invoker context, then INSTRUCTIONS section) is printed so the user can see exactly what would be sent to the AI. No AI process is spawned.

**Verification:** User sees the assembled prompt; no iteration or AI invocation occurs.

## Acceptance criteria

- [ ] A dry-run option is available for the run command (and the review command if applicable).
- [ ] When dry-run is used, the system prints the fully assembled prompt (CONTEXT section when non-empty, then INSTRUCTIONS section) that would be sent to the AI.
- [ ] No AI process is invoked in dry-run mode.
- [ ] The output is sufficient for the user to understand what would be run and to debug prompt assembly.
