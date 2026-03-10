# O004: Observability

## Who

Users who run Ralph and need to know why a command or operation stopped, how it performed, and what to do when something fails — for diagnosis, tuning, and scripting.

## Statement

The user understands why the command or operation stopped, how it performed, and what to do when something fails.

## Why it matters

A loop runner that silently exits is a black box. The user re-runs the command, checks file diffs, reads AI output, and tries to reconstruct what happened. Did it succeed? Did it hit the failure threshold? Did it time out? Did the signal never appear? Without clear termination reporting and execution statistics, diagnosing loop behavior is guesswork and tuning config is trial and error. The same applies to the review command: the user must know whether the review succeeded, whether the prompt had errors, or whether the run failed.

## Output and log controls

Two separate controls govern what the user sees:

- **Log level** — Controls how much Ralph itself logs (e.g. error, warn, info, debug): iteration progress, timing, errors, and other operational messages. The user can set log level explicitly (e.g. via config or CLI); that value overrides any shortcut.

- **Show AI command output** — Controls whether the AI process’s stdout is streamed to the terminal in real time. Ralph always captures that output for signal scanning; this setting only determines if the user sees it. Default is true in normal runs so the user can watch the AI work. When false, the user sees only Ralph’s logs and final summary.

**Quiet** is a shortcut for minimal output: it sets log level to a minimal level (e.g. error only) and show AI command output to false, so scripts and CI get only essential messages and no streamed AI output. Explicit log level or show AI command output override the shortcut where set. There is no separate “verbose” flag; the user gets more output by raising log level and/or enabling show AI command output.

## Verification

- When the AI CLI command is missing or invalid (e.g. alias not found, binary not on PATH), Ralph reports a clear error before the loop or review starts, so the user understands why the run could not start and can fix config or install the tool.
- On success signal: Ralph prints a completion message, reports iteration count and timing, exits 0.
- On failure threshold: Ralph prints the threshold value and consecutive failure count, exits with a distinct code (e.g. 1).
- On max iterations exhausted: Ralph prints the iteration count and limit, exits with a distinct code (e.g. 2).
- On interruption (e.g. SIGINT/SIGTERM): Ralph exits with a distinct code (e.g. 130).
- User sets log level and show AI command output as desired: Ralph’s logs follow the log level; when show AI command output is true, the AI’s output is streamed to the terminal while Ralph still captures it for signal scanning.
- User uses quiet: only minimal Ralph logs (e.g. errors) and no streamed AI output, unless the user explicitly overrides log level or show AI command output.
- User runs dry-run and sees the fully assembled prompt (preamble + prompt) printed without spawning an AI process.
- After a multi-iteration run, Ralph reports iteration statistics (e.g. min/max/mean duration).
- For review: exit code and report make it clear whether the review completed, the prompt had errors, or the run failed; the user can see what to do next.

## Non-outcomes

- Operational messages and the AI command stream go to stdout (the run's log). stderr is reserved for fatal or startup errors only. Ralph does not provide persistent log files.
- Ralph does not provide per-iteration diffs or file change tracking between iterations.
- Ralph does not integrate with external monitoring, alerting, or observability systems.
- Ralph does not support replaying or debugging past executions.
