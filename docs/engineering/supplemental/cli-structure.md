# CLI structure (engineering supplemental)

This document defines the top-level command structure for Ralph. It is engineering supplemental: it will feed into the formal engineering docs when the full `docs/engineering/` tree is in place. Product outcomes define *what* the CLI must support; this doc specifies *how* commands are named and grouped to avoid breaking changes.

## Design principles

1. **Verb-first** — Top-level subcommands are actions: `ralph run`, `ralph review`, `ralph list`, `ralph show`, `ralph version`. No noun-first groups (e.g. `ralph config show`) unless we later introduce a cohesive set of actions on one object and document migration.
2. **Migration-safe** — Names and shape are chosen so we can extend (e.g. flags, output formats) without renaming commands or breaking scripts.
3. **List and show use subcommands** — `ralph list` with no subcommand shows all (prompts and aliases). `ralph list prompts` and `ralph list aliases` list only that type. Similarly, `ralph show` takes a required object: `ralph show config`, `ralph show prompt [name]`, `ralph show alias [name]`. Same verb–noun shape for both verbs.
4. **Show is a first-class verb** — `ralph show` takes an object (what to show). It can show detailed information about config, prompts, and aliases. New “what to show” objects can be added under `show` without new top-level verbs.
5. **No lifecycle in binary** — Install, uninstall, and upgrade are documented processes (scripts, package manager), not Ralph subcommands (per O006, O011).

## Command layout

| Command | Purpose |
|--------|--------|
| **ralph run** | Run the loop. Prompt via alias (named prompt), `-f` file, or stdin. Flags: `--dry-run`, config file, log level, quiet, etc. |
| **ralph review** | Review prompt (alias, `-f` file, stdin). Report and suggested revision; `--apply` with confirmation or non-interactive flag. |
| **ralph list** | List prompts and/or aliases from resolved config. With no subcommand: show all. Subcommands: **prompts**, **aliases**. See below. |
| **ralph show** | Show detailed information. Takes a required object: **config**, **prompt** [name], **alias** [name]. See below. |
| **ralph version** | Print version (install verification, scripting). |

**ralph list**

- **ralph list** — List all (prompts and aliases) from resolved config.
- **ralph list prompts** — List only prompts.
- **ralph list aliases** — List only AI command aliases.

**ralph show**

- **ralph show config** — Effective (resolved) config for current context; optional provenance. Same config resolution as run/list (cwd, `--config`, env, etc.). Satisfies O002/R007.
- **ralph show prompt [name]** — Detailed information about a prompt (from resolved config). Name identifies the prompt; when omitted, behavior is defined (e.g. show all or error).
- **ralph show alias [name]** — Detailed information about an AI command alias (from resolved config). Name identifies the alias; when omitted, behavior is defined (e.g. show all or error).

**Help:** `ralph --help`, `ralph run --help`, `ralph show --help`, etc. (standard).

**Not in CLI:** `ralph install`, `ralph uninstall`, `ralph upgrade` — documented procedures only.

## Global vs per-command

Options that affect config resolution (e.g. config file path, env such as `RALPH_CONFIG`) are global and apply to run, review, list, and show so behavior is consistent. Exact flag names and env vars are implementation detail; this doc only fixes the command set and the principle that resolution is shared.

## Future-proofing

- **Show objects** — New detail views live under `ralph show`: e.g. `ralph show config`, `ralph show prompt`, `ralph show alias`. Additional objects (e.g. `ralph show prompt-template`) can be added without new top-level verbs.
- **List** — Subcommands `prompts` and `aliases` scope the list; `ralph list` with no subcommand shows all. Flags (e.g. `--output=json`) can extend format without new subcommands.
- **Dry-run and apply** — Remain flags on `run` and `review` respectively; no separate top-level commands.
- **More config actions** — If we add e.g. validate config, we can add `ralph validate-config` (verb) or a `validate` object under a future group; the current design does not depend on it.

## Comparison with kubectl

kubectl uses a **verb–noun** structure: the first word after `kubectl` is a verb (get, create, run, delete, describe, logs, exec, apply, …), and the **noun** (resource type) is the next argument: `kubectl get pods`, `kubectl create deployment`, `kubectl describe pod my-pod`. So the pattern is `kubectl <verb> [resource-type] [name] [flags]`.

Ralph’s design **aligns with that philosophy** at the top level: we are verb-first. The first word after `ralph` is always an action (run, review, list, show, version). Both **list** and **show** use a verb–noun pattern: `ralph list` (all), `ralph list prompts`, `ralph list aliases`; `ralph show config`, `ralph show prompt [name]`, `ralph show alias [name]`. So the verb is fixed and the next token is the object (or omitted for list = all). Same idea as `kubectl get pods` / `kubectl get deployments`, scoped to a small set of objects. New list or detail views extend these verbs with new subcommands.

## Traceability

- **O001** — run (iterative completion).
- **O002** — list (R006), show config (R007); config resolution shared.
- **O004** — observability via flags on run/review (dry-run, log level, quiet).
- **O005** — review (and --apply).
- **O006, O008** — version and help for install/discoverability; no install subcommand.
- **O011** — stable command names and options within non-breaking versions.
