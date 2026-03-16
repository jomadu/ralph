# Documentation

## Responsibility

The documentation component is the body of user-facing documentation and procedures that support install, uninstall, upgrade, discoverability, user guidance, and automation contracts. It is not a runtime component; it is the set of docs and documented procedures that the product team maintains so that users can install Ralph, get to a first successful run, use the product correctly, run it from scripts and CI with stable exit codes and non-interactive behavior, and upgrade without breaking config or workflows. Implementation-wise, this "component" is realized as docs in the repository (and any published sites), install/uninstall/upgrade scripts or package-manager instructions, and release notes. The engineering README does not list O/R in component docs; requirement assignments for this component are in the [engineering README](../README.md).

Implements the requirements assigned to this component in the [engineering README](../README.md).

## Interfaces

**Consumes**

- Product intent (outcomes and requirements) and engineering implementation specs (config schema, exit codes, report format) so that documentation stays accurate.

**Produces**

- **User documentation** — Covers how to use Ralph: run, review, list, show, version; config file and layers; prompt sources; loop behavior; review report and apply; exit codes and automation. Must cover the capabilities required by the product so users can operate the product from the docs alone.
- **Install procedure** — Documented steps (and/or scripts, package manager) to install Ralph so the binary is invocable (e.g. `ralph version` succeeds). No `ralph install` subcommand; install is a procedure.
- **Uninstall procedure** — Documented steps to remove Ralph and clean up (e.g. config, caches) so uninstall is complete and documented.
- **Upgrade procedure** — How to upgrade to a chosen version or update within a non-breaking version; backward compatibility and migration are documented (per O011). No `ralph upgrade` subcommand; upgrade is a procedure.
- **Release notes** — For changes that affect users or automation (new features, breaking changes, migrations). Users can understand what changed and how to adapt.
- **Discoverability content** — What Ralph is, why to use it, and how to get to a first successful run (install, first command). May live in README, docs site, or both.
- **Stable contract documentation** — Exit codes (run and review), non-interactive flags, review result format (result.json in report directory) and report directory layout, and config contract so scripts and CI can rely on them. When the contract changes (e.g. new exit code, new result field), docs and release notes are updated.

## Implementation spec

### Doc coverage

The primary user-facing doc for subcommands, config, and exit codes is the repository [README](../../../README.md); the canonical exit-code contract is [docs/exit-codes.md](../../../exit-codes.md). Together they satisfy O007/R001.

Documentation must cover at least:

- All subcommands and their options (run, review, list, show, version) and the fact that install/uninstall/upgrade are procedures, not subcommands.
- Config: layer order, config file location (global, workspace), explicit config option, and the canonical config structure (loop, prompts, aliases) so users can author valid config.
- Prompt sources: alias, file path, stdin; how to specify them for run and review.
- Loop behavior: iterations, failure threshold, timeout, signals, precedence, preamble, AI command/alias, streaming, log level.
- Review: invocation (alias, file, stdin); report **directory** and the five files (result.json, summary.md, original.md, revision.md, diff.md); **result.json** as the machine-parseable outcome and how to use it for CI (or exit code); apply and confirmation; revision output path requirement when prompt is from stdin.
- **Writing Ralph prompts:** A user-facing guide ([docs/writing-ralph-prompts.md](../../writing-ralph-prompts.md)) that explains how to write a well-formed Ralph prompt using the same four dimensions as `ralph review` (O005/R007): signal and state, iteration awareness, scope and convergence, subjective completion criteria. The guide is the single source of truth; `ralph show prompt-guide` outputs it verbatim. Discoverable from the path to first run or README and via that command (O007/R005, O008/R005).
- Exit codes: run (success, failure threshold, max iterations, interrupt, error) and review (0, 1, 2) with exact values and semantics so automation can gate reliably. The canonical user and automation doc is [docs/exit-codes.md](../../exit-codes.md); README summarizes and links to it.
- Non-interactive use: flags and environment so CI/scripts can run without prompts; behavior when confirmation would be required in non-interactive mode.

### Release notes and stable contract

- **Release notes** — Published for each release (e.g. [GitHub Releases](https://github.com/jomadu/ralph/releases)); semantic-release creates releases on push to main/rc/alpha/beta. Each release describes intentional behavior changes, deprecations, and changes that affect config or scripts so users can adapt. The location and expectations are documented in [docs/release-notes.md](../../release-notes.md); README links to it.
- **Stable contract** — Documented in one place for scripts and CI: [docs/release-notes.md](../../release-notes.md) summarizes the contract (exit codes, review summary format, config, non-interactive use) and links to canonical specs ([docs/exit-codes.md](../../exit-codes.md), [review.md](review.md), [config.md](config.md)). When the contract changes, release notes for that release explain the change and migration.

### Procedures (install, uninstall, upgrade)

- **Install** — Document where the binary is installed, how to add it to PATH if needed, and how to verify (`ralph version`). Any install script or package (e.g. Homebrew, npm, direct binary) must be documented so users can install and invoke Ralph. The repository provides `scripts/install.sh`: installs from GitHub release artifacts to a configurable directory. Default directory is platform- and privilege-dependent: Linux — `/usr/local/bin` if writable else `~/.local/bin` (FHS/XDG); macOS — `/usr/local/bin` if writable else `~/bin`; Windows (Git Bash) — `~/bin`. Override with `RALPH_INSTALL_DIR` or `--dir`. The script does not modify PATH or create symlinks; the user adds the install directory to PATH. README documents install steps, version selection, and verification.
- **Uninstall** — Document removal of the binary and optional cleanup (config directory, caches). Uninstall is complete and documented so users can remove Ralph cleanly. The repository provides `scripts/uninstall.sh`: looks for the binary in standard locations (`/usr/local/bin`, `~/.local/bin`, `~/bin`) in that order and removes it from the first found; user config (e.g. `ralph-config.yml`) is not removed. If the user installed with `--dir` to a custom path, they must remove that binary manually. Because install does not modify PATH or symlinks, uninstall leaves no broken references.
- **Upgrade** — Document how to upgrade (e.g. reinstall over existing, package manager upgrade). Backward compatibility within a non-breaking version and any documented migration for breaking changes (per O011) are explained. Release notes link to upgrade and migration when relevant.

### Discoverability (single place or linked set)

User-facing docs must provide a single place or clearly linked set for:

- **CLI and flags** — Authoritative spec: [cli.md](cli.md). README summarizes subcommands and flags and links to this file.
- **Config and env** — Schema, layer order, and all env vars: [config.md](config.md). README summarizes and links.
- **Exit codes** — Canonical contract: [docs/exit-codes.md](../../exit-codes.md). README summarizes and links.

The README section "Where to look (CLI, config, env, exit codes)" is the discoverability entry point and links to each canonical source.

### Discoverability content (what, why, first run)

Per O008, a new user can find what Ralph is, why it exists, install steps, and a path to a first successful run:

- **What and why** — README section "What is Ralph and why use it?" gives a short description (loop runner for AI-driven tasks) and rationale (automated iteration until success/failure signal).
- **Install** — README section "Install and Uninstall" documents install (and uninstall); verification is `ralph version`.
- **First run** — README section "Path to first run" describes one path from "I have Ralph" to a first command that completes successfully: verify PATH, choose prompt source (file, stdin, or alias), run (e.g. `ralph run -f <path> -n 1` or stdin), with minimal prerequisites (AI CLI on PATH) stated.
- **List and help** — `ralph list` and `ralph run --help` / `ralph --help` expose prompts, aliases, and subcommands; README Quick Start and Subcommands summarize and link to [cli.md](cli.md).
- **Prompt-writing guide** — New users who do not know how to write a prompt can find [Writing Ralph prompts](../../writing-ralph-prompts.md) via README or path-to-first-run content, or run `ralph show prompt-guide` to get the full guide verbatim (O008/R005).

### Troubleshooting (O007/R004)

User-facing docs must help users resolve common problems. The README includes a **Troubleshooting** section that covers at least:

- **Prompt not found / unknown alias** — How prompt source is resolved (alias, file, stdin); using `ralph list` to verify; effect of `--config` and config file locations.
- **Config file not found** — Behavior when `--config` is used (no fallback); where global and workspace config are read; skipping missing files without error.
- **Wrong or unexpected exit code** — Reference to exit code semantics (run: 0, 2, 3, 4, 130; review: 0, 1, 2) and how to interpret them; common causes for exit 2 (e.g. missing AI command, stdin + apply without `--prompt-output`, report directory unwritable or path is an existing file). Optionally: report is now a directory; look in `result.json` for status and in `summary.md` for narrative.
- **AI command not found** — AI CLI must be on PATH or specified via `--ai-cmd`; Ralph validates before the loop and exits 2 with a clear error.
- **ralph: command not found** — PATH must include the install directory; how to verify (e.g. `ralph version`). Standard locations are `/usr/local/bin`, `~/.local/bin`, `~/bin`.

Troubleshooting content lives in the README and links to the canonical CLI, config, and exit-code docs so users can resolve issues without leaving the doc set.

### Consistency with implementation (doc accuracy)

When implementation specs change in a way that affects user or automation contract (e.g. new flag, exit code, config key, or report directory and file formats (result.json, summary.md, original.md, revision.md, diff.md)), the documentation component is updated so that user docs, procedures, and contract docs remain accurate. Update the engineering spec first (cli.md, config.md, run-loop.md, review.md, etc.), then README and [docs/exit-codes.md](../../exit-codes.md) as needed. Product requirements that reference "documented" behavior are satisfied by this component.
