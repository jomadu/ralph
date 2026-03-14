# R004: Path to First Run

**Outcome:** O008 — Discoverability

## Requirement

The path from having Ralph to a first successful run is documented and achievable.

## Detail

A new user should be able to go from "I have Ralph installed" (or can run it via a documented method) to "I just ran the loop (or review) and it exited with a clear outcome" without reverse-engineering the repo or guessing. Documentation describes the steps: minimal config (if needed), choice of prompt source (alias, file, or stdin), and a first command that can complete successfully. The path is linear enough that a user without prior knowledge of the codebase can follow it and see a successful run (e.g. exit 0 or a documented "run completed" outcome).

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User has only the binary | Docs describe minimal config and a first command (e.g. run with a file prompt or stdin) that can succeed. |
| User has config but no prompt file | Docs state how to supply prompt (e.g. stdin, or create a minimal file) so the first run can complete. |
| User follows path step-by-step | Each step is documented; the final step is a run or review that completes with a documented outcome. |
| Multiple valid paths exist | At least one path is clearly recommended or enumerated so the user is not lost in options. |

### Examples

#### Documented path to first run

**Input:** New user has installed Ralph (per R002). User opens the docs (README or user docs) and looks for "first run" or "getting started."

**Expected output:** Documentation describes a sequence such as: (1) ensure minimal config or use defaults, (2) choose a prompt source (e.g. `ralph run --file <path>` or stdin), (3) run the command. The command can complete successfully (e.g. loop exits 0 on success signal or review exits 0/1 with report).

**Verification:** A user following the documented path can achieve a first successful run without reading source code or guessing.

#### From list to run

**Input:** User runs `ralph list` and sees a prompt name. User wants to run that prompt.

**Expected output:** Docs or help explain how to run by alias (e.g. `ralph run <alias>`). User runs the command and it completes with a documented outcome.

**Verification:** The path from "I see a prompt in list" to "I ran it and it finished" is documented and achievable.

## Acceptance criteria

- [ ] Documentation describes at least one path from "have Ralph" to "first successful run" (run or review).
- [ ] The path includes any minimal config, prompt source, or prerequisites so the user can satisfy them.
- [ ] A user following the path can complete a run or review and see a documented outcome (e.g. exit 0 or documented non-zero).
- [ ] The path does not require reading the codebase or guessing; it is achievable from docs and CLI help alone.

## Dependencies

- [R001](R001-what-and-why.md) — User must know what Ralph is before following the path.
- [R002](R002-install-and-first-command.md) — User must have Ralph installed or runnable.
- [R003](R003-list-and-help.md) — List and help support discovering what to run; the path may include using list/help.
