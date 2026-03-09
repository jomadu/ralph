# Publishing the first pre-release

These steps assume you have branches `main` and `impl` and want to publish the first **pre-release** (e.g. `v1.0.0-rc.1`) without cutting a stable release from `main`.

## 1. Ensure release tooling is ready

- **package-lock.json:** If you haven’t already, run `npm install` in the repo root and commit `package-lock.json` so CI can run `npm ci`.
- **Conventional commits:** All commits that will be included in the release must follow [Conventional Commits](https://www.conventionalcommits.org/) (e.g. `feat: ...`, `fix: ...`). CI runs commitlint on push; fix any failures.

## 2. Create a pre-release branch

Only these branches trigger releases: `main` (stable), `rc`, `alpha`, `beta` (pre-releases). Your branch `impl` is not in that list, so create one of the pre-release branches from the code you want to release.

**Option A — Create `rc` from `impl` (recommended for first RC):**

```bash
git fetch origin
git checkout impl
git pull origin impl
git checkout -b rc
git push -u origin rc
```

**Option B — Create `alpha` or `beta`:**

Same as above but use `alpha` or `beta` instead of `rc`. The first release from that branch will be e.g. `v1.0.0-alpha.1` or `v1.0.0-beta.1`.

## 3. Ensure there is a releasable commit

Semantic-release only creates a release when there is at least one commit that triggers a version bump (e.g. `feat:`, `fix:`, `perf:`). If the tip of `rc` (or `alpha`/`beta`) already has such a commit, you’re set. If not, add a small conventional commit and push:

```bash
# Example: add a commit that triggers a release
git commit --allow-empty -m "chore(release): trigger first pre-release"
# Or use a real change, e.g.:
# git add ...
# git commit -m "feat: add X"
git push origin rc
```

Note: `chore:` alone usually does **not** bump the version. For the very first release you typically need at least one `feat:` or `fix:`. If in doubt, use something like:

```bash
git commit --allow-empty -m "feat: initial pre-release"
git push origin rc
```

## 4. Let CI run

- Pushing to `rc` (or `alpha`/`beta`) runs the [Release workflow](.github/workflows/release.yml): commitlint, then semantic-release.
- Semantic-release will:
  - Determine the next version (e.g. `1.0.0-rc.1` for branch `rc`).
  - Update `CHANGELOG.md`, create a git tag, and open a GitHub release.
  - The **build** job will build binaries for Linux, macOS, and Windows (amd64/arm64) and attach them to that release.

## 5. Verify

- On GitHub: **Releases** should show a new pre-release (e.g. `v1.0.0-rc.1`) with assets like `ralph-1.0.0-rc.1-linux-amd64`, `ralph-1.0.0-rc.1-darwin-arm64`, etc.
- Install and test:
  ```bash
  ./scripts/install.sh 1.0.0-rc.1
  ralph version
  ```

## Troubleshooting

| Problem | What to check |
|--------|----------------|
| No release created | Ensure the branch name is exactly `rc`, `alpha`, or `beta` and that there is at least one commit that triggers a bump (e.g. `feat:` or `fix:`). |
| Commitlint fails | Fix commit messages to follow conventional commits; amend or add a new commit and push again. |
| Build/upload fails | Check the **Build** job logs for the matrix entry (e.g. `linux-amd64`). Fix any Go or Makefile issues and push a new commit to the same branch to re-run the workflow. |
| Want to re-run without new commits | Re-run the failed jobs from the GitHub Actions run, or push an empty commit: `git commit --allow-empty -m "ci: re-run release" && git push origin rc`. |

## Later: stable release from `main`

When you’re ready for the first **stable** release (e.g. `v1.0.0`):

1. Merge your pre-release branch (e.g. `rc`) into `main`.
2. Push `main`. Semantic-release will then publish a stable release from `main` (no `-rc`/`-alpha`/`-beta` suffix).
