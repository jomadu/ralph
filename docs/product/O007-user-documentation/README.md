# O007: User Documentation

## Who

Users of Ralph (new and existing) who need an authoritative reference for how to install, configure, run, update, upgrade, uninstall, and troubleshoot the product.

## Statement

Users have documentation that explains how to use the product.

## Why it matters

Without clear documentation, users guess at subcommands, config format, and behavior. They rely on examples that may be outdated or incomplete. Support burden increases and adoption drops. Documentation that is complete, accurate, and discoverable lets users self-serve and use Ralph as intended.

## Verification

- User can find documentation (e.g. in-repo docs, or a documented URL) that describes how to install Ralph, configure loop running (e.g. iteration limits, timeouts, signals), prompts, and commands, run the loop and review, update and upgrade to a newer version, uninstall, and interpret exit codes and output.
- User can look up CLI commands, flags, config keys, and environment variables in one place (or a clearly linked set of docs).
- Documentation matches actual behavior: described commands and options work as documented; exit codes and report formats align with the docs.
- User can resolve common problems (e.g. "prompt not found", "wrong exit code") by following the documented behavior and troubleshooting guidance where provided.

## Non-outcomes

- Documentation does not replace implementation specs; engineering holds schemas and protocols. User docs explain how to use the product, not how to build it.
- Ralph does not require a specific doc format or hosting (e.g. a dedicated docs site). The outcome is that usable documentation exists and is referenced from the product (e.g. README, release notes).
- Ralph does not embed interactive tutorials or in-app help beyond what is needed to run and configure it; the outcome is written documentation.

## Risks

| Risk | Mitigating Requirement |
|------|------------------------|
| Key topics (install, config, run, etc.) missing from docs | [R001 — Doc coverage](R001-doc-coverage.md) |
| User cannot find CLI, flags, or config reference in one place | [R002 — Doc discoverability](R002-doc-discoverability.md) |
| Documentation does not match actual behavior | [R003 — Doc accuracy](R003-doc-accuracy.md) |
| User stuck on common problems with no guidance | [R004 — Troubleshooting](R004-troubleshooting.md) |

## Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| [R001](R001-doc-coverage.md) | Documentation covers install, configuration, run, review, update, upgrade, uninstall, and exit codes | ready |
| [R002](R002-doc-discoverability.md) | User can look up CLI, flags, config, and environment variables in one place or a clearly linked set | ready |
| [R003](R003-doc-accuracy.md) | Documentation matches actual product behavior | ready |
| [R004](R004-troubleshooting.md) | User can resolve common problems using documented troubleshooting guidance | ready |
