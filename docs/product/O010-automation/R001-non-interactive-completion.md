# R001: Non-interactive completion

**Outcome:** O010 — Automation

## Requirement

Ralph completes loop and reviewer runs without requiring interactive input when invoked with documented non-interactive options or config, and exits with a documented exit code.

## Detail

When run from scripts or CI, the loop and the reviewer must not block on prompts for user input (e.g. confirmation to apply a revision, or interactive choices). The user invokes Ralph with the documented non-interactive options or config (e.g. a flag to skip confirmation, or config that disables interactive behavior). Under those conditions, the command runs to completion—success or failure—and exits with an exit code whose meaning is documented (see R002). Flows that by design require confirmation (e.g. apply revision in interactive mode) remain available when the user does not use the non-interactive variant; automation uses the non-interactive variant so that no input is required.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Loop run with non-interactive config/options | Completes without prompting; exits with documented code (success, exhaustion, failure threshold, etc.). |
| Reviewer run with non-interactive options, no apply | Completes without prompting; report and revision produced; exit code per R002. |
| Reviewer run with apply and non-interactive apply option | Apply is performed without confirmation prompt; exits with documented code. |
| User invokes without non-interactive option where apply would require confirmation | Ralph may prompt for confirmation or exit with a code indicating "confirmation required"; behavior is documented so scripters know to use the non-interactive option. |
| Pipe or redirect of stdin for prompt content | Not treated as "interactive input" for confirmation; non-interactive completion still required when non-interactive options are used. |
| Timeout or iteration limit reached | Command completes (does not hang); exit code reflects exhaustion or timeout per R002. |

### Examples

#### Loop from script

**Input:** Script runs the run command with config that specifies prompt alias, AI command, and iteration limit; no TTY. No non-interactive option needed if loop has no confirmation points, or user passes a documented non-interactive option.

**Expected output:** Loop runs until success, failure threshold, or max iterations; then process exits with the documented success code or the documented failure/exhaustion code as appropriate. No prompt for user input.

**Verification:** Process exits; script check yields a documented code; no hang waiting for input.

#### Reviewer from CI without apply

**Input:** CI runs the review command with report output path (no apply). No TTY.

**Expected output:** Review runs; report is written to the specified path; process exits with the documented success code or the documented failure (or prompt-errors) code as appropriate.

**Verification:** Report file exists; exit code is one of the documented set; no prompt for confirmation.

#### Reviewer with apply in non-interactive mode

**Input:** Script runs the review command with apply and the documented non-interactive apply option. Config or option indicates apply without confirmation.

**Expected output:** Review runs; suggested revision is applied to the configured or specified path; process exits with the documented success or failure code. No confirmation prompt.

**Verification:** Revision was written; exit code is documented; no interactive prompt was shown.

## Acceptance criteria

- [ ] When the loop is invoked with documented non-interactive options or config, it completes (success or failure) without requiring interactive input and exits with a documented exit code.
- [ ] When the reviewer is invoked with documented non-interactive options or config (and, if apply is requested, the documented non-interactive apply option), it completes without requiring interactive input and exits with a documented exit code.
- [ ] The set of options or config that constitute "non-interactive" for loop and reviewer is documented so script and CI authors know how to avoid blocking.
- [ ] Exit codes are documented (per R002) so scripts can interpret success vs failure vs exhaustion etc.

## Dependencies

- R002 — Documented stable exit codes (so "documented exit code" is well-defined).
