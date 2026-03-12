# R006: Log level and show AI output

**Outcome:** O004 — Observability

## Requirement

The system respects the configured log level and show-AI-output setting for what is emitted to the user.

## Detail

Two separate controls govern what the user sees. **Log level** controls how much the tool logs (e.g. error, warn, info, debug): iteration progress, timing, errors, and other operational messages. The user can set log level explicitly; that value overrides any shortcut. **Show AI command output** (streaming) controls whether the AI process's stdout is streamed to the terminal in real time. This applies to **both** `ralph run` and `ralph review`: when enabled, the user sees the AI's output as it is produced; when disabled, they see only the tool's logs and final outcome. For run, the tool also captures output for signal scanning regardless of streaming. For review, the outcome is still derived from the report directory (result.json, etc.); streaming only affects visibility. Default is typically true so the user can watch the AI work; when false, the user sees only the tool's logs and final summary.

**Default:** Show AI command output (streaming) is **on** by default so the user can watch the AI work. The only CLI flag to turn it off is **`--no-stream`**; there is no "enable streaming" flag. Config and environment can set streaming to false or true; `--no-stream` overrides to off for that run.

**Quiet** is a shortcut that sets log level to a minimal level (e.g. error only) and show AI command output to false. Explicit log level overrides the shortcut for logs; to re-enable streaming after quiet, the user uses config or env (no CLI flag to re-enable). There is no separate "verbose" flag; the user gets more output by raising log level.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Log level set to debug | The system emits debug-level logs (e.g. iteration progress, timing) in addition to info/warn/error |
| Log level set to error only | The system emits only error-level messages (and any required completion/summary per other requirements) |
| Show AI output false | AI process stdout is not streamed to the terminal; user sees the tool's logs and final summary (run: still captured for signal scanning; review: outcome from report directory) |
| Show AI output true (default) | AI stdout is streamed to the terminal in real time (run: also captured for signal scanning; review: outcome still from report directory) |
| Quiet shortcut with no overrides | Minimal logs (e.g. errors only) and no streamed AI output |
| Quiet plus explicit log level | Log level override wins; show AI output remains off unless config/env set streaming true |

### Examples

#### Quiet mode in CI (run)

**Input:** User runs `ralph run` with quiet (minimal log level and show AI output off). Run completes successfully.

**Expected output:** Only minimal logs (e.g. errors) and no streamed AI output; user still sees completion message and the process exits with the documented success code per R002.

**Verification:** Script/CI gets essential outcome without noisy AI stream; completion message still present.

#### Quiet mode in CI (review)

**Input:** User runs `ralph review` with quiet (minimal log level and show AI output off). Review completes successfully.

**Expected output:** Only minimal logs (e.g. errors) and no streamed AI output; user still sees report path and the process exits with the documented review exit code.

**Verification:** Script/CI gets essential outcome without noisy AI stream; report directory is still produced and exit code is correct.

#### Explicit log level override

**Input:** User sets quiet but overrides log level to info.

**Expected output:** Logs at info level (iteration progress, timing, etc.); show AI output remains off unless also overridden.

**Verification:** User sees operational messages but not the AI stream, as configured.

## Acceptance criteria

- [ ] The system honors the configured log level for the tool's own logs (error, warn, info, debug or equivalent).
- [ ] The system honors the show-AI-output setting for both run and review: when true, AI stdout is streamed to the terminal; when false, it is not streamed (run: still captured for signal scanning; review: outcome from report directory).
- [ ] When the user sets log level explicitly, that overrides the quiet shortcut for logs; streaming has no CLI "on" flag (default is on; only `--no-stream` and quiet turn it off).
- [ ] Quiet shortcut results in minimal log level and no streamed AI output when not overridden.
- [ ] Default behavior (e.g. show AI output true in normal runs) is as specified so the user can watch the AI work when not using quiet.
