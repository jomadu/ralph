# Review

## Responsibility

The review component implements `ralph review`: it accepts the prompt from alias, file path, or stdin; invokes the reviewer (e.g. via the backend with review-specific instructions or a dedicated review flow); produces a report containing narrative feedback and a machine-parseable summary plus the full suggested revision; writes the report to a user-chosen path; and optionally applies the revision to a path with confirmation (or non-interactive flag). It sets the process exit code based on review outcome (success, prompt errors, or failure to complete). It does not list or show config; the CLI dispatches those to other behavior.

Implements the requirements assigned to this component in the [engineering README](../README.md).

## Interfaces

**Consumes**

- **Prompt source** — Resolved by CLI/config: alias (resolved to prompt content), file path, or stdin. Invalid or missing source must be detected before review runs; the component reports clear error and exit code 2.
- **Resolved config** — For AI command/alias when review uses the backend; report output path, revision output path (when apply is requested), and apply/confirmation options.
- **Flags** — e.g. `--apply`, revision output path (`--prompt-output` or equivalent), non-interactive. When prompt is from stdin and apply is requested, revision output path is required; if missing, error and exit 2.

**Produces**

- **Report file** — Written to user-chosen path (or default). Contains narrative feedback, machine-parseable summary line, and full suggested revision (inline or by reference). Report file must exist at the expected path for exit codes 0 or 1; if report write fails, exit 2.
- **Revision file** — When apply is requested and path is valid: revision written to path, subject to confirmation when overwriting and when interactive. When prompt was from stdin and apply is requested without revision path: error, no write, exit 2.
- **Exit code** — 0 (review completed, no prompt errors), 1 (review completed, prompt has errors), 2 (review or apply did not complete successfully). Derivation: after report file is verified to exist, parse the machine-parseable summary to choose 0 vs 1; otherwise 2.

**Calls**

- Config: prompt source and paths already resolved by CLI; review receives them.
- Backend: when the review flow uses the AI to evaluate the prompt and generate report and revision (e.g. via embedded review instructions and prompt content).
- (Optional) Internal report verifier and summary parser: verify report file exists; parse summary line for exit code derivation.

## Implementation spec

### Report format

The report file contains three parts:

1. **Narrative feedback** — Human-readable evaluation of the prompt (e.g. signal discipline, statefulness, scope, convergence). Format is not strictly mandated; the reviewer produces readable feedback.
2. **Machine-parseable summary** — A single line that scripts and CI can parse to derive exit code 0 vs 1. Canonical format (see also `internal/review/summary.go`):
   - Line matching: `ralph-review:\s*status=(ok|errors|warnings)(\s+errors=(\d+))?(\s+warnings=(\d+))?`
   - Examples: `ralph-review: status=ok`, `ralph-review: status=errors errors=2`, `ralph-review: status=warnings warnings=1 errors=0`
   - Exit code derivation: `status=ok` and `errors=0` (or absent) → 0; `status=errors` or `errors>=1` → 1; `status=warnings` with `errors=0` → 0 (if policy allows). Missing or malformed summary → 1 (fail-safe for CI).
3. **Full suggested revision** — The complete suggested prompt text. May be inline in the report or written to the revision output path when apply is used; the full revision must be available to the user/CI.

### Exit code derivation

- **0** — Review completed; report file exists; machine-parseable summary indicates no errors (e.g. `status=ok`, `errors=0`).
- **1** — Review completed; report file exists; summary indicates one or more errors (e.g. `status=errors`, `errors>=1`). Ralph parses the report file (after verifying it exists) to set 0 vs 1.
- **2** — Review or apply did not complete: invalid prompt source, report write failure, stdin + apply without revision output path, confirmation required in non-interactive session, or reviewer/internal error. Never 0 or 1 in these cases.

If the report file does not exist at the expected path after the review phase, exit 2. Do not exit 0 or 1 without a present, parseable report.

### Apply and confirmation

- When the user requests apply, the revision is written to the user-chosen path. When the prompt was supplied via **stdin**, the user must supply the revision output path; if not, error and exit 2.
- When overwriting an existing file (or applying to the same file as the source), the system must confirm with the user in an interactive session unless a non-interactive flag is set. In non-interactive mode, apply either is skipped with a clear error or proceeds without confirmation per product definition; the behavior must be documented.
- Invalid path (e.g. directory, unwritable) → error and exit 2.

### Invocation inputs

- Exactly one prompt source per invocation: alias (resolved to content), file path, or stdin. Conflicting sources (e.g. file and stdin) → defined behavior (precedence or usage error); no silent ambiguity.
- Invalid alias or missing file → fail before running review; clear error; exit 2.
