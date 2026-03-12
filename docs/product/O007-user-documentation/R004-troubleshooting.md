# R004: Troubleshooting

**Outcome:** O007 — User Documentation

## Requirement

User can resolve common problems using documented troubleshooting guidance.

## Detail

User-facing documentation includes guidance for common problems—for example: prompt or alias not found, config file not found, wrong or unexpected exit code, AI command not found, or the `ralph` command not found after install. For each, the docs explain what causes the situation and how to fix it (e.g. check config location, PATH, or exit code semantics). Troubleshooting content may live in a dedicated section (e.g. README Troubleshooting) or be linked from a single place; it should reference the canonical CLI, config, and exit-code docs so users can resolve issues without leaving the doc set.

### Examples

#### User gets "prompt not found" or unknown alias

**Input:** User runs `ralph run build` and sees an error that the prompt or alias is not found.

**Expected output:** Documentation explains that the alias must be defined in resolved config (global, workspace, or `--config`) and suggests using `ralph list prompts` / `ralph list aliases` to verify. It clarifies the effect of `--config` (no fallback) and where config files are read.

**Verification:** User can follow the guidance to define or fix the alias or config path and get a successful run.

#### User gets an unexpected exit code

**Input:** User ran `ralph run` or `ralph review` and got a non-zero exit code; they are unsure what it means.

**Expected output:** Documentation lists exit codes (run: 0, 2, 3, 4, 130; review: 0, 1, 2) with clear meanings and common causes (e.g. exit 2 = error before loop or invocation error; 3/4 = loop ended without success). For review, the report is a directory; result.json holds the machine outcome. It points to the stable exit-code contract (e.g. docs/exit-codes.md) for scripts and CI.

**Verification:** User can interpret the exit code and take the right action (fix config, check signals, add `--yes` for apply, etc.).

#### User gets "ralph: command not found"

**Input:** User installed Ralph but the shell cannot find `ralph`.

**Expected output:** Documentation states that the install directory must be on PATH and suggests verifying with `ralph version` and checking `~/.config/ralph/install-state` for the install path. It reminds the user to add that directory to PATH (e.g. in shell profile).

**Verification:** User can add the install directory to PATH and run `ralph version` successfully.

## Acceptance criteria

- [ ] Documentation covers at least: prompt/alias not found, config not found, wrong exit code, AI command not found, and `ralph` not found.
- [ ] For each covered problem, the doc explains cause and resolution (or links to the right config/CLI/exit-code docs).
- [ ] Troubleshooting is findable (e.g. README section or linked from a single place) and references canonical specs where relevant.

## Dependencies

None.
