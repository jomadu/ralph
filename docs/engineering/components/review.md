# Review

## Responsibility

The review component implements `ralph review`: it accepts the prompt from alias, file path, or stdin; invokes the reviewer (e.g. via the backend with review-specific instructions or a dedicated review flow); produces a **report directory** with five files — `result.json` (structured status), `summary.md` (narrative), `original.md` (prompt as submitted), `revision.md` (suggested revision), `diff.md` (diff); writes the report to a user-chosen directory (or default); and optionally applies the revision to a path with confirmation (or non-interactive flag). It sets the process exit code based on review outcome (success, prompt errors, or failure to complete). It does not list or show config; the CLI dispatches those to other behavior.

Implements the requirements assigned to this component in the [engineering README](../README.md).

## Interfaces

**Consumes**

- **Prompt source** — Resolved by CLI/config: alias (resolved to prompt content), file path, or stdin. Invalid or missing source must be detected before review runs; the component reports clear error and exit code 2.
- **Resolved config** — For AI command/alias when review uses the backend; report output directory (or path that denotes the directory), revision output path (when apply is requested), apply/confirmation options, and **streaming** (show AI output). AI command/alias may be set in config (root or per-prompt), environment, or CLI; the CLI passes the already-resolved value. When prompt is from stdin and apply is requested, revision output path is still required.
- **Flags** — e.g. `--apply`, revision output path (`--prompt-output` or equivalent), non-interactive, `--no-stream` (do not show AI command output; default is to show it; same semantics as run). When prompt is from stdin and apply is requested, revision output path is required; if missing, error and exit 2.

**Produces**

- **Report directory** — User-chosen or default (e.g. `./ralph-review/`). Directory contains: `result.json`, `summary.md`, `original.md`, `revision.md`, `diff.md`. All files must be written for exit codes 0 or 1; if directory creation or any required write fails, exit 2.
- **Revision file** — When apply is requested and path is valid: revision written to user-chosen path; content is the same as `revision.md` in the report dir. Subject to confirmation when overwriting and when interactive. When prompt was from stdin and apply is requested without revision path: error, no write, exit 2.
- **Exit code** — 0 (review completed, no prompt errors), 1 (review completed, prompt has errors), 2 (review or apply did not complete successfully). Derivation: after the AI returns, read `result.json` from the report directory to choose 0 vs 1; missing or invalid `result.json` → exit 2.

**Calls**

- Config: prompt source and paths already resolved by CLI; review receives them.
- Backend: when the review flow uses the AI to evaluate the prompt and generate report and revision (e.g. via embedded review instructions and prompt content). When streaming is enabled, the component passes a stream writer (e.g. stdout) so the backend tees AI stdout to the user in real time (O004/R006); outcome is still derived from the report directory.
- (Optional) After invoke, verify report directory contains `result.json`; read it for exit code derivation. For apply, read `revision.md` from the report directory.

## Implementation spec

### Report directory and file formats

- **result.json** — JSON object with `status` (string: "ok" | "errors" | "warnings"), optional `errors` (number), optional `warnings` (number). No narrative. Canonical for exit code and CI. **Written by the AI** in the report directory.
- **summary.md**, **original.md**, **revision.md**, **diff.md** — As named: narrative feedback, exact prompt as submitted, full suggested revision, diff between original and revision. All **written by the AI** in the report directory.

### Exit code derivation

Read `result.json` from the report directory after invoke. `status=ok` and `errors=0` (or absent) → 0; `status=errors` or `errors>=1` → 1; `status=warnings` with `errors=0` → 0 (if policy allows). Missing or invalid `result.json` → exit 2.

### Apply and confirmation

- Revision content = **content read from `revision.md`** in the report directory. When the user requests apply, that content is written to `--prompt-output`. When the prompt was supplied via **stdin**, the user must supply the revision output path; if not, error and exit 2.
- When overwriting an existing file (or applying to the same file as the source), the system must confirm with the user in an interactive session unless a non-interactive flag is set. In non-interactive mode, apply either is skipped with a clear error or proceeds without confirmation per product definition; the behavior must be documented.
- Invalid path (e.g. directory, unwritable) → error and exit 2.

### Invocation inputs

- Exactly one prompt source per invocation: alias (resolved to content), file path, or stdin. Conflicting sources (e.g. file and stdin) → defined behavior (precedence or usage error); no silent ambiguity.
- Invalid alias or missing file → fail before running review; clear error; exit 2.

### Review invocation (agent-based)

Review **requires** an AI command (no fallback). The component invokes the backend with a single assembled “review prompt”:

- **Prompt assembly:** The review prompt is assembled with the **report directory path interpolated from run options** (e.g. placeholder `{{REPORT_DIR}}` replaced with the absolute path). The prompt instructs the AI to create the five files **in that directory** and to **not** put file contents in its response; the AI responds only with a short confirmation (e.g. "Created the review report at <path>.").
- **Embedded review instructions** (Ralph-owned): instruct the AI to evaluate the user’s prompt along the four dimensions (O005/R007) — signal and state, iteration awareness, scope and convergence, subjective completion — and to create the five report files on disk.
- **User prompt content** (the prompt under review), clearly delimited in the review prompt (e.g. in a fenced block or section).

If no AI command is configured (no value from config, env, or CLI flags such as `loop.ai_cmd`/`loop.ai_cmd_alias`, `RALPH_LOOP_AI_CMD`/`RALPH_LOOP_AI_CMD_ALIAS`, or `--ai-cmd`/`--ai-cmd-alias`) or backend invocation fails, review fails with a clear error and exit 2.

### Streaming (show AI output)

Review respects the same **streaming** (show-AI-output) setting as run (config `loop.streaming`, env, and CLI `--stream` / `--no-stream`; O004/R006). When streaming is true, the component invokes the backend with a non-nil stream writer (e.g. process stdout) so the AI command's stdout is streamed to the user in real time. When streaming is false (e.g. `--quiet` or `--no-stream`), the stream writer is nil and the user sees only the tool's logs and final report path. Exit code and apply content are always derived from the report directory (result.json, revision.md); streaming only affects visibility of the AI's raw output.

### Expected AI behavior

- **No structured stdout format for report body:** Ralph does not parse AI stdout for report content. It **reads result.json** (and when applying, **revision.md**) from the report directory after the AI returns. Failure to find or parse `result.json` → exit 2.
- The AI creates the five files in the report directory; it does not emit report body content in its response. Ralph derives exit code and apply content from files on disk.
