# Release notes and stable contract

This document describes **where to find release notes** for Ralph and **what constitutes the stable contract** that scripts and CI can rely on. When the contract or behavior changes, release notes describe the change and how to adapt.

## Where to find release notes

Ralph publishes release notes for each release so users can see what changed and how to adapt config or scripts.

- **Location:** [GitHub Releases](https://github.com/maxdunn/ralph/releases) for this repository. Each tagged release (created by semantic-release on push to `main`, `rc`, `alpha`, or `beta`) has a release page with notes.
- **Content:** For each release, intentional behavior changes, deprecations, and changes that affect config, CLI, or scripts are described. Patch releases that only fix bugs or add optional behavior may have minimal notes. Pre-releases (alpha/beta) may have lighter notes; the intent is that users can learn about changes.
- **Finding them:** From the repo, go to the Releases section or follow the link in the README ([Release notes and stable contract](../README.md#release-notes-and-stable-contract)). After upgrading, check the release notes for the version you moved to.

When a release changes the **stable contract** (exit codes, summary format, or config schema), the release notes for that release explain the change and any migration or workaround.

## Stable contract (for scripts and CI)

Within the compatibility contract (e.g. same major version), the following are stable so that scripts and CI can rely on them. Exact values and semantics are documented; they do not change meaning within a non-breaking range. New codes or keys may be added in a backward-compatible way; existing ones are not repurposed.

| Contract | Canonical doc | Summary |
|----------|----------------|---------|
| **Exit codes** | [docs/exit-codes.md](exit-codes.md) | Run: 0 (success), 2 (error before loop), 3 (max iterations), 4 (failure threshold), 130 (interrupt). Review: 0 (OK), 1 (prompt errors), 2 (did not complete). |
| **Review result format** | [docs/engineering/components/review.md](engineering/components/review.md) | Review result: result.json in report directory (status, errors, warnings). Report directory contains result.json, summary.md, original.md, revision.md, diff.md. Canonical spec: docs/engineering/components/review.md. |
| **Config** | [docs/engineering/components/config.md](engineering/components/config.md) | Layer order, config file locations, schema (loop, prompts, aliases), env vars (`RALPH_CONFIG_HOME`, `RALPH_LOOP_*`). |
| **Non-interactive use** | README and [exit-codes.md](exit-codes.md) | Flags and env for CI; `--yes` for review apply; exit semantics for gating. |

When any of these change in a way that could break scripts or docs, the change and migration guidance are documented in the release notes for that release. See also [docs/exit-codes.md](exit-codes.md) for full exit-code semantics and automation guidance.
