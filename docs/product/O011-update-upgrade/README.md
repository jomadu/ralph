# O011: Update Upgrade

## Who

Users who already have Ralph installed and want to upgrade to a specific version or update within a non-breaking version, without losing config, breaking existing prompts, or invalidating documented workflows.

## Statement

Users can upgrade to a chosen version or update within a non-breaking version, and when they do, existing config and workflows continue to work.

## Why it matters

Without a supported way to upgrade or update, users are stuck on an old version or must uninstall and re-install by hand. Upgrades that silently change behavior or invalidate config force users to fix things after every release. That increases support cost and discourages upgrades. A clear upgrade path — the ability to move to a newer version and backward compatibility for config and documented behavior — keeps existing users safe when they adopt new versions.

## Verification

- User can upgrade to a specific version via the documented process (e.g. re-run install or upgrade steps with that version). After upgrade, the new version is invocable and behaves as documented.
- User can update within a non-breaking version (e.g. latest patch or minor on the same major line) so they get fixes and improvements without changing the compatibility contract. The documented process makes it clear how to do this.
- After upgrade or update, existing config and prompts continue to work as-is; no required migration step blocks normal use. When format or behavior changes require migration, the user is given documentation or tools they run themselves — Ralph does not rewrite user config automatically.
- Documented commands, options, and exit codes that scripts or docs rely on remain valid across non-breaking upgrades, or changes are documented in release notes with migration guidance.
- User can read release notes (or equivalent) to learn about intentional behavior changes or deprecations and adjust config or scripts if needed.

## Non-outcomes

- Ralph does not automatically migrate or rewrite user config to new formats when it finds it (e.g. on upgrade or at runtime). The outcome is backward compatibility: existing config keeps working as-is. When migration is needed, the user is given documentation or optional tools they run themselves; Ralph does not modify config files without explicit user action.
- Ralph does not promise infinite backward compatibility; major versions may introduce breaking changes with documented migration. The outcome is that upgrades within the same major line (or as documented) do not break existing config or workflows without notice.
- Ralph does not perform automatic background or unattended updates; the user initiates upgrade or update when they choose, via the documented process.
- Third-party tools (AI CLIs, wrappers) are outside this outcome; the focus is Ralph's own config, CLI contract, and documented workflows.
