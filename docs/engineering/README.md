# Engineering

This directory describes the **structure, placement, and implementation specifications** of Ralph. Product documentation (`docs/product/`) is the source of truth for *who* we're building for, *what* we're building, and *why* (outcomes and requirements at the level of intent). Engineering is the source of truth for *where* each capability lives and *how* the system is built: component responsibilities, interfaces, and the hard specs (schemas, APIs, protocols) that implementers build from. Component docs reference product requirements by ID (e.g. O001/R002); the assignments below are the single place for the full map.

## High-level flow

- **Entry:** User invokes `ralph` with a subcommand: `run`, `review`, `list`, `show`, or `version`. The CLI parses arguments and resolves configuration (layers: defaults, global file, workspace file, explicit file, environment, CLI flags).
- **Run path:** `ralph run` — Config resolves prompt source (alias, file, or stdin) and loop settings. The run-loop loads the prompt once, then iterates: it invokes the backend (AI CLI) with the assembled prompt on stdin, captures stdout, detects success or failure signals (or applies AI-interpreted precedence when configured), and continues or exits. Observability (exit codes, reports, log level, dry-run, iteration statistics) is produced along the run path.
- **Review path:** `ralph review` — Config resolves prompt source. The review component evaluates the prompt, produces a report and a suggested revision, and optionally applies the revision to a user-chosen path with confirmation. Review outcome (exit code and report) is clear so the user knows whether the review completed, the prompt had errors, or the run failed.
- **List and show:** `ralph list`, `ralph show config`, `ralph show prompt`, `ralph show alias` — Config supplies resolved prompts and aliases; the CLI exposes them. Effective (resolved) config is produced by config and shown via `ralph show config`.
- **Backend:** The run-loop (and review, when using AI for interpretation) invokes the user-chosen AI command via the backend component: prompt on stdin, stdout captured, environment and working directory inherited.

Install, uninstall, and upgrade are **documented procedures** (scripts, package manager), not subcommands (see [CLI structure](supplemental/cli-structure.md)).

## Components

Each component is responsible for the requirements listed. Component detail and implementation specs live under [components/](components/) (one file or directory per component). This README is the single place for requirement assignments; component docs do not duplicate the O/R list.

| Component | One-line description | Assigned requirements |
|-----------|------------------------|------------------------|
| [cli](components/cli.md) | Command parsing and dispatch; subcommands run, review, list, show, version; help and non-interactive flags | O001/R003, O002/R006, O005/R001, O005/R004, O005/R006, O008/R003, O009/R003, O010/R001 |
| [config](components/config.md) | Configuration layer resolution, effective config, defaults, global/workspace/explicit file, env, CLI overrides; prompt and alias definitions; read-only unless opt-in | O002/R001–R007, O009/R002, O010/R004 |
| [run-loop](components/run-loop.md) | Iteration loop: validate AI command, load prompt once, invoke backend, detect signals, continue or exit; run reports and exit codes | O001/R001–R009, O004/R001–R008 |
| [backend](components/backend.md) | Invoke AI CLI with prompt on stdin, capture stdout, inherit environment and cwd; support structured output via signals or wrapper | O003/R001–R003 |
| [review](components/review.md) | Prompt review: report and suggested revision, apply with confirmation; review exit code and report; machine-parseable summary | O005/R001–R008, O004/R009, O009/R001, O009/R003, O010/R003 |
| [documentation](components/documentation.md) | User-facing docs, install/uninstall/upgrade procedures, release notes, discoverability content, documented exit codes and contract | O006/R001–R004, O007/R001–R004, O008/R001, O008/R002, O008/R004, O010/R002, O011/R001–R004 |

## Consistency

- Every requirement ID above exists in the product tree (`docs/product/`).
- Every product requirement is assigned to at least one component.
- No component is empty; boundaries are distinct (no two components claim the same requirement in conflicting ways).

Component docs and implementation specs are added in E2; this overview is locked before expanding.
