# Releases and versioning

Ralph uses [semantic versioning](https://semver.org/) and [conventional commits](https://www.conventionalcommits.org/) to drive releases.

## For users

- **Stable releases** are published from the `main` branch (e.g. `v1.0.0`).
- **Pre-releases** are published from:
  - `rc` → release candidates (e.g. `v1.0.0-rc.1`)
  - `alpha` → alpha pre-releases (e.g. `v1.0.0-alpha.1`)
  - `beta` → beta pre-releases (e.g. `v1.0.0-beta.1`)

Releases and changelogs are created automatically on GitHub when commits are pushed to these branches. Each release includes **pre-built binaries** for common platforms (Linux, macOS, Windows; amd64 and arm64). The install script **only** installs from these release artifacts (no build from source). You can install the latest release (`./scripts/install.sh`) or a specific version (`./scripts/install.sh 1.0.0`). The `ralph version` command prints the version of the binary you are running.

## For maintainers

- **Commit messages** must follow the [Conventional Commits](https://www.conventionalcommits.org/) format (e.g. `feat: add X`, `fix: Y`, `docs: Z`). CI and an optional local git hook (husky + commitlint) enforce this.
- **Release flow:** Push to `main` (or `rc`/`alpha`/`beta`) with conventional commits. CI runs [semantic-release](https://github.com/semantic-release/semantic-release), which:
  - Determines the next version from commit types
  - Creates a git tag and [GitHub release](https://github.com/jomadu/ralph/releases) (release notes are generated there; no in-repo changelog file is updated)
  - Builds binaries for Linux (amd64, arm64), macOS (amd64, arm64), and Windows (amd64, arm64) and attaches them to the release (asset names: `ralph-<version>-<os>-<arch>` or `ralph-<version>-windows-<arch>.exe`)
- **Local checks:** Run `npm ci` then `npx commitlint --from main --to HEAD` to validate commits before pushing. For a dry-run of the next release: `npx semantic-release --no-ci`.
