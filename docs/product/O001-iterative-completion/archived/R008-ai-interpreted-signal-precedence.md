> **Withdrawn.** This requirement has been archived. When both success and failure signals appear on the last line, the system treats the outcome as success (static precedence). See [R006 — Signal precedence](../R006-signal-precedence.md). Archived 2026-03-16.

---

# R008: AI-interpreted signal precedence

**Outcome:** O001 — Iterative Completion

## Requirement

The system may optionally resolve signal precedence by invoking the AI once with a built-in prompt that asks the AI to interpret the iteration output and decide success or failure; if the interpretation run does not yield a clear answer, the system applies a defined fallback (e.g. treat as failure or use static precedence).

## Detail

When both success and failure signals appear in the same output, the user may enable an option to let the AI decide the outcome. In that case, the system invokes the AI once with a built-in prompt (owned by the product, not user-editable) that provides the iteration output and asks whether the task succeeded or failed. The AI's response is parsed to determine success or failure. If the response is clear, that outcome is used. If the response is ambiguous or the interpretation run fails (e.g. crash, no parseable answer), the system applies a defined fallback: e.g. treat as failure, or apply static precedence (R006). Exactly one extra AI invocation is made per ambiguous iteration; there are no retries of the interpretation step.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Option disabled | R006 static precedence always applies. |
| Both signals present; option enabled; AI returns clear "success" | Treat iteration as success. |
| Both signals present; option enabled; AI returns clear "failure" | Treat iteration as failure. |
| Interpretation run crashes or times out | Apply fallback (e.g. failure or static precedence). |
| AI response not parseable | Apply fallback. |
| Only one signal present | No interpretation step; R004 or R005 applies directly. |
| Interpretation run counts as iteration or not | Documented (e.g. does not count toward max iterations, or does). |

### Examples

#### Interpretation yields success

**Input:** Both signals in output; option enabled. The system runs interpretation with the built-in prompt; the AI responds that the task succeeded.

**Expected output:** The system treats the iteration as success; exits with the documented success code (if this was the only iteration) or continues and may succeed later.

**Verification:** Outcome is success; only one interpretation invocation for that iteration.

#### Interpretation unclear; fallback to failure

**Input:** Both signals in output; option enabled. Interpretation run returns an ambiguous or unparseable response. Fallback = treat as failure.

**Expected output:** The system treats the iteration as failure; increments consecutive-failure count and continues or exits per R005.

**Verification:** No undefined state; fallback is applied; user may see that fallback was used (where documented).

## Acceptance criteria

- [ ] When the user enables AI-interpreted precedence and both success and failure signals appear in an iteration's output, the system may invoke the AI once with a built-in prompt to interpret the output.
- [ ] The built-in prompt is owned by the product (not user-editable).
- [ ] At most one interpretation invocation is made per ambiguous iteration; no retries of the interpretation step.
- [ ] If the interpretation run does not yield a clear success/failure answer, the system applies a defined fallback (e.g. treat as failure or use R006).
- [ ] When the option is disabled, R006 applies; no interpretation step is run.

## Dependencies

- R006 — Fallback and default behavior when interpretation is not used or fails.
