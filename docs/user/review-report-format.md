# Review report summary format (O5 R6)

When you run `ralph review`, the AI writes a report file. The report must include a **machine-parseable summary line** so scripts and CI can gate on the result without parsing narrative text.

## Canonical format

A single line in the report matching:

```
ralph-review: status=<status> [errors=N] [warnings=M]
```

- **status** — One of `ok`, `errors`, or `warnings`.
- **errors** — Optional. Integer count of errors found in the prompt.
- **warnings** — Optional. Integer count of warnings.

**Regex (for parsers):**  
`ralph-review:\s*status=(ok|errors|warnings)(\s+errors=(\d+))?(\s+warnings=(\d+))?`

**Examples:**

- `ralph-review: status=ok` → exit 0
- `ralph-review: status=ok errors=0` → exit 0
- `ralph-review: status=errors` → exit 1
- `ralph-review: status=errors errors=2` → exit 1
- `ralph-review: status=warnings warnings=1 errors=0` → exit 0

## Exit code derivation

Ralph reads the report file after the review-phase AI exits and looks for this line. Exit code is derived as follows:

| Summary content              | Ralph exit code |
|-----------------------------|-----------------|
| `status=ok` (and errors absent or 0) | 0 |
| `status=warnings` and errors=0 or absent | 0 |
| `status=errors` or `errors>=1`       | 1 |
| Line missing or unparseable         | 1 (fail-safe for CI) |

Exit code **2** is used when the review did not complete (e.g. config invalid, prompt load failure, report file missing). It is never set from the summary line; see R8/R9.

## Reference

- Spec: [R6 — Report format and exit codes](../intent/O5-prompt-review/R6-report-format-exit-codes.md)
- Implementation: `internal/review/summary.go` (ParseReportSummary), embedded instructions in `internal/review/review_instructions.md`
