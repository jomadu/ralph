# R003: Structured Output via Signals or Wrapper

**Outcome:** O003 — Backend Agnosticism

## Requirement

The system works with AI CLIs that produce structured or non–plain-text output when the user configures success and failure signal strings to match content within that output, or uses a wrapper so Ralph sees plain text for signal scanning.

## Detail

Some AI CLIs emit JSON, tool-call traces, or other structured streams on stdout. Ralph scans stdout for configurable success and failure substrings. The user can point signals at fragments that appear in structured output, or run a wrapper script (invoked via alias or direct command) that normalizes output to plain text Ralph can scan. Ralph does not ship vendor-specific adapters beyond built-in aliases; adapting output remains the user’s responsibility when defaults are insufficient.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User configures signals matching tokens inside JSON lines | Loop stops when those tokens appear in captured stdout |
| User’s wrapper writes progress to stderr and assistant text to stdout | Ralph scans stdout only; wrapper pattern is documented for agents that mix streams |
| Output is binary or unreadable | User must wrap or reconfigure; out of scope for core normalization |

### Examples

**Input:** User runs an alias that invokes a small script filtering an agent’s JSON stream to a single line containing `DONE` on success.

**Expected output:** Ralph detects the configured success signal and stops the loop as intended.

**Verification:** Dry-run and short runs confirm signals fire when the wrapper emits the expected text.
