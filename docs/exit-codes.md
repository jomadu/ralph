# Ralph exit codes

This document is the **stable contract** for Ralph’s exit codes. Scripts and CI can rely on these values within the compatibility contract (e.g. same major version). Exact numeric values are defined by the run-loop and review components; see `docs/engineering/components/run-loop.md` and `docs/engineering/components/review.md`.

---

## ralph run

| Exit code | Meaning |
|-----------|---------|
| **0** | Success — success signal detected; loop completed. |
| **2** | Error before loop — invalid or missing AI command, invalid config, or prompt source error. Clear error is printed; loop did not start. |
| **3** | Max iterations — iteration limit reached without detecting the success signal. |
| **4** | Failure threshold — consecutive failure count reached the configured threshold. |
| **130** | Interrupted — process received SIGINT (e.g. Ctrl+C) or SIGTERM. Conventional code for SIGINT is 128 + 2. |

### For scripts and CI

- **Success:** Exit 0 → task completed successfully.
- **Retry or alert:** Exit 3 (max iterations) or 4 (failure threshold) → run did not complete successfully; you may retry or fail the job.
- **Configuration/input error:** Exit 2 → fix config or prompt source before re-running.
- **Interrupted:** Exit 130 → user or system interrupted; do not treat as success or normal failure.

---

## ralph review

| Exit code | Meaning |
|-----------|---------|
| **0** | Review completed; report directory written (result.json, summary.md, original.md, revision.md, diff.md); no prompt errors (result.json indicates OK). |
| **1** | Review completed; report directory written; prompt has one or more errors (result.json indicates errors). |
| **2** | Review or apply did not complete — invalid prompt source, report write failure, stdin + apply without `--prompt-output`, confirmation required in non-interactive mode without `--yes`, or internal error. |

### For scripts and CI

- **Gate on prompt quality:** Exit 0 = pass, exit 1 = prompt has errors (fail the gate), exit 2 = invocation error (fail the job). For outcome details, read result.json in the report directory.
- **Apply flow:** Use `--yes` when applying in non-interactive mode to avoid exit 2 for missing confirmation.

---

## Other commands

- **ralph version**, **ralph list**, **ralph show** — Exit 0 on success. Parse/usage errors (unknown command, bad args, missing required argument) use a consistent non-zero code (e.g. 1); see CLI documentation.

---

## Stability

Within the compatibility contract (e.g. same major version), these exit codes do not change meaning. New codes may be added in a backward-compatible way. Release notes will describe any change to the contract.
