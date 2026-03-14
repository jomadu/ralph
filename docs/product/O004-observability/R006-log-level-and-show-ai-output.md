# R006: Log level and show AI output

**Outcome:** O004 — Observability

## Requirement

The system respects the configured log level and show-AI-output setting for what is emitted to the user.

## Detail

Two separate controls govern what the user sees. **Log level** controls how much the tool logs (e.g. error, warn, info, debug): iteration progress, timing, errors, and other operational messages. The user can set log level explicitly; that value overrides any shortcut. **Show AI command output** controls whether the AI process's stdout is streamed to the terminal in real time. The tool always captures that output for signal scanning; this setting only determines visibility to the user. Default is typically true so the user can watch the AI work; when false, the user sees only the tool's logs and final summary.

**Quiet** is a shortcut that sets log level to a minimal level (e.g. error only) and show AI command output to false. Explicit log level or show AI command output override the shortcut where set. There is no separate "verbose" flag; the user gets more output by raising log level and/or enabling show AI command output.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Log level set to debug | The system emits debug-level logs (e.g. iteration progress, timing) in addition to info/warn/error |
| Log level set to error only | The system emits only error-level messages (and any required completion/summary per other requirements) |
| Show AI output false | AI process stdout is not streamed to the terminal; the system still captures it for signal scanning; user sees the tool's logs and final summary |
| Show AI output true (default) | AI stdout is streamed to the terminal in real time while the system captures it |
| Quiet shortcut with no overrides | Minimal logs (e.g. errors only) and no streamed AI output |
| Quiet plus explicit log level | Log level override wins; show AI output remains off unless overridden |
| Quiet plus explicit show AI output | Show AI output override wins; user sees streamed AI output |

### Examples

#### Quiet mode in CI

**Input:** User runs with quiet (minimal log level and show AI output off). Run completes successfully.

**Expected output:** Only minimal logs (e.g. errors) and no streamed AI output; user still sees completion message and the process exits with the documented success code per R002.

**Verification:** Script/CI gets essential outcome without noisy AI stream; completion message still present.

#### Explicit log level override

**Input:** User sets quiet but overrides log level to info.

**Expected output:** Logs at info level (iteration progress, timing, etc.); show AI output remains off unless also overridden.

**Verification:** User sees operational messages but not the AI stream, as configured.

## Acceptance criteria

- [ ] The system honors the configured log level for the tool's own logs (error, warn, info, debug or equivalent).
- [ ] The system honors the show-AI-output setting: when true, AI stdout is streamed to the terminal; when false, it is not streamed (but still captured for signal scanning).
- [ ] When the user sets log level or show AI output explicitly, that setting overrides the quiet shortcut for that dimension.
- [ ] Quiet shortcut results in minimal log level and no streamed AI output when not overridden.
- [ ] Default behavior (e.g. show AI output true in normal runs) is as specified so the user can watch the AI work when not using quiet.
