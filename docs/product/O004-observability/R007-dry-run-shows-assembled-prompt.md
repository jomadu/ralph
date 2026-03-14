# R007: Dry-run shows assembled prompt

**Outcome:** O004 — Observability

## Requirement

The system supports a dry-run mode that prints the assembled prompt without invoking the AI.

## Detail

The user can verify what prompt will be sent to the AI (preamble + prompt) without running the loop or spawning an AI process. This supports understanding why the loop or review would behave as it does and debugging config or prompt assembly. No AI process is started in dry-run.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Dry-run with default config | Fully assembled prompt (preamble + prompt) printed to stdout (or as designated); no AI invocation |
| Dry-run with custom preamble | Assembled output includes the configured preamble and the prompt file content |
| Dry-run and log level | The system may emit minimal log context (e.g. that dry-run is active); primary output is the assembled prompt |
| Dry-run with missing AI command | Behavior may follow R001 (error before run); dry-run may still show assembled prompt without needing a valid AI command |

### Examples

#### Dry-run output

**Input:** User invokes the run command in dry-run mode with a prompt file and config that adds a preamble.

**Expected output:** The fully assembled prompt (preamble + prompt content) is printed so the user can see exactly what would be sent to the AI. No AI process is spawned.

**Verification:** User sees the assembled prompt; no iteration or AI invocation occurs.

## Acceptance criteria

- [ ] A dry-run option is available for the run command (and the review command if applicable).
- [ ] When dry-run is used, the system prints the fully assembled prompt (preamble + prompt) that would be sent to the AI.
- [ ] No AI process is invoked in dry-run mode.
- [ ] The output is sufficient for the user to understand what would be run and to debug prompt assembly.
