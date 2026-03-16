# Plan: Remove AI-interpreted signal precedence and wire max-output-buffer

## Objectives

1. **Remove AI-based signal resolution (hard removal)**  
   Drop the optional AI-interpreted precedence (O001/R008): remove `signal_precedence: ai_interpreted` and `--signal-precedence ai_interpreted`, the interpretation prompt, and the extra AI invocation when both signals appear on the last line. Static precedence (e.g. success wins) is the only behavior. We are in release-candidate versions, so no deprecation period.

2. **Archive product requirement R008**  
   Mark O001/R008 (AI-interpreted signal precedence) as withdrawn/deprecated and move the requirement into an archived state for traceability.

3. **Wire and enforce max-output-buffer**  
   Keep the `--max-output-buffer` flag and add config/file/env support. Use it as a medium, configurable cap (bytes) so that when the last line is very long we still capture it in full up to that limit. Backend must enforce the cap (e.g. sliding window: retain only the last `max_output_buffer` bytes of stdout so the last line is preserved within the cap).

---

## Task dependency tree (high level)

- **T1** (Archive R008) has no dependencies.
- **T2** (Remove AI-interpreted implementation) depends on nothing; can run in parallel with T1.
- **T3** (Add max_output_buffer to config) has no dependencies.
- **T4** (Backend enforce max-output-buffer) depends on T3 if we pass limit from config; for a first slice we can add the parameter and implement with a default.
- **T5** (CLI: remove signal-precedence, wire max-output-buffer) depends on T2 and T3.
- **T6** (Update all documentation) depends on T1, T2, T5 (and T4 for backend/docs).

Detailed tasks below are ordered by dependency; each task includes full context to complete it.

---

## Task list

### T1. Archive product requirement R008 (AI-interpreted signal precedence)

**Priority:** 1 (no dependencies)

**Context:**  
O001/R008 is the product requirement for optional AI-interpreted signal precedence. We are withdrawing this feature. The requirement doc must be marked as withdrawn/deprecated and placed in an archived state so intent and history remain traceable.

**What to do:**

- Create an archive location for withdrawn requirements, e.g. `docs/product/O001-iterative-completion/archived/` or `docs/product/_archived/O001-R008-ai-interpreted-signal-precedence.md`.
- Move or copy `docs/product/O001-iterative-completion/R008-ai-interpreted-signal-precedence.md` into that archive.
- At the top of the archived doc, add a clear notice, e.g.:  
  **"Withdrawn / deprecated.** This requirement was withdrawn in [release/date]. AI-interpreted signal precedence is no longer supported; only static precedence (e.g. success wins) is used when both signals appear on the last line."
- Remove the original R008 doc from the active product tree (or replace it with a one-paragraph stub that points to the archived doc and states "Withdrawn").
- Update the engineering README component table: remove O001/R008 from the run-loop row (so run-loop no longer has that assignment).
- Update the O001 iterative-completion README (or requirement index) so R008 is listed as withdrawn/archived with a pointer to the archived doc.

**Deliverable:** R008 is archived and all references in the requirement map point to withdrawn/archived state; engineering README no longer assigns R008 to run-loop.

---

### T2. Remove AI-interpreted precedence from implementation

**Priority:** 2 (no dependencies; can run in parallel with T1)

**Context:**  
Signal detection already uses only the last non-empty line (`LastNonEmptyLine` in `internal/runloop/run.go`). When both success and failure signals appear on that line and `opts.Loop.SignalPrecedence == "ai_interpreted"`, the code invokes the AI once with `BuildInterpretationPrompt(stdout, ...)` and parses the response. We are removing that path entirely and always applying static precedence (success wins).

**What to do:**

- **Run-loop (`internal/runloop/run.go`):**
  - Remove the block that checks `opts.Loop.SignalPrecedence == "ai_interpreted"` and calls `BuildInterpretationPrompt` and `ParseInterpretationResponse` (lines ~167–181).
  - When both `hasSuccess` and `hasFailure` are true, treat the iteration as success (static precedence: success wins). So the condition becomes: if `hasSuccess`, report success and return; else (failure or no signal) increment consecutive failures and continue/exit.
  - Remove any imports used only for interpretation.
- **Interpret package:**  
  Remove or delete `internal/runloop/interpret.go` (BuildInterpretationPrompt, ParseInterpretationResponse, interpretation prompt constants). If any other package references it, remove those references; run-loop is the only caller.
- **Config structs and defaults:**
  - In `internal/config/layer.go`: remove `SignalPrecedence` from `LoopSection` (or leave the field but stop documenting it; if we leave it, schema and merge must still be updated so we don’t accept `ai_interpreted`).
  - In `internal/config/defaults.go`: remove `SignalPrecedence: "static"` from `DefaultLoopSettings()` (or leave as the only allowed value; see schema below).
  - In `internal/config/loop_merge.go`: remove the block in `ApplyLoopSection` that sets `out.SignalPrecedence` from `section.SignalPrecedence`; remove from `ApplyEnvOverlayToLoop` any overlay of SignalPrecedence.
  - In `internal/config/env.go`: remove `SignalPrecedence` from `EnvOverlay` and remove parsing of any `RALPH_LOOP_SIGNAL_PRECEDENCE` (if present).
  - In `internal/config/schema.go`: remove `ErrInvalidSignalPrecedence` and `validSignalPrecedence`; remove validation of `loop.SignalPrecedence` in `validateLoop`. If we keep `SignalPrecedence` in the struct for backward compatibility (e.g. ignore unknown values), document that only `static` is supported and any other value is ignored; otherwise remove the field from `LoopSection` and `LoopSettings` everywhere.
- **CLI (partial; full removal of flag is T5):**  
  Do not remove the flag yet if you want T5 to own “CLI surface”; here we only remove the behavior that uses it (run-loop and config). So T2 focuses on run-loop, interpret, config, and schema. T5 will remove the `--signal-precedence` flag and overlay.

**Tests:**

- In `internal/runloop/run_test.go`: remove or rewrite tests that set `SignalPrecedence = "ai_interpreted"` and assert on interpretation (e.g. “interpreter said success”). Either remove those test cases or change them to assert static precedence (both signals on last line → success) without interpretation.
- In `internal/config/schema_test.go`: remove or adjust tests that validate `signal_precedence` (e.g. `ai_interpreted` valid). If the field is removed, remove those tests; if we only allow `static`, update expected validation accordingly.
- Run `make test` and fix any remaining references to interpretation or `ai_interpreted`.

**Deliverable:** No code path uses AI-interpreted precedence; config/schema no longer accept or apply `ai_interpreted`; tests updated and passing.

---

### T3. Add max_output_buffer to config (file, env, defaults)

**Priority:** 3 (no dependencies)

**Context:**  
`--max-output-buffer` exists in the CLI but is not in config or LoopSettings and is not passed to the run-loop or backend. We want a medium, configurable cap (bytes) so the last line of AI stdout can be captured in full up to that limit. This task adds the setting to the config model, file schema, env overlay, and defaults.

**What to do:**

- **Structs:**
  - In `internal/config/layer.go`: add to `LoopSection` a field, e.g. `MaxOutputBuffer *int` with YAML tag `max_output_buffer,omitempty`.
  - In `internal/config/loop_merge.go`: ensure `LoopSettings` has a field for effective value (e.g. `MaxOutputBuffer int`). If `LoopSettings` is in the same package, add `MaxOutputBuffer int` to it (in `internal/config`, e.g. in layer.go or a shared file).
  - In `internal/config/loop_merge.go`: in `ApplyLoopSection`, when `section.MaxOutputBuffer != nil`, set `out.MaxOutputBuffer = *section.MaxOutputBuffer`.
  - In `internal/config/env.go`: add `MaxOutputBuffer *int` to `EnvOverlay`; in `ParseEnvOverlay`, parse `RALPH_LOOP_MAX_OUTPUT_BUFFER` (integer, ≥ 0); invalid or negative → clear error.
  - In `internal/config/loop_merge.go`: in `ApplyEnvOverlayToLoop`, when `overlay.MaxOutputBuffer != nil`, set `out.MaxOutputBuffer = *overlay.MaxOutputBuffer`.
- **Defaults:**  
  In `internal/config/defaults.go`, set a default for max output buffer in `DefaultLoopSettings()`, e.g. `MaxOutputBuffer: 65536` (64 KiB) or `262144` (256 KiB). Document the default in `docs/engineering/components/config.md`.
- **Schema:**  
  In `internal/config/schema.go`, add validation for `max_output_buffer` in file layers (e.g. must be ≥ 0 if present). Add an error like `ErrInvalidMaxOutputBuffer` if needed.
- **Docs:**  
  In `docs/engineering/components/config.md`: add `max_output_buffer` to the config file structure (loop section), to the built-in defaults table, and to the environment variables table (e.g. `RALPH_LOOP_MAX_OUTPUT_BUFFER`).

**Deliverable:** Config supports `max_output_buffer` from file and env with a sensible default; schema validates it; docs updated.

---

### T4. Backend: enforce max-output-buffer when capturing stdout

**Priority:** 4 (depends on T3 for the numeric value to pass; can implement with a constant default first, then wire in T5)

**Context:**  
The backend currently captures full stdout in a `bytes.Buffer` with no limit. We want to cap the retained stdout at `max_output_buffer` bytes so that (1) we don’t grow unbounded, and (2) we still capture the entirety of the last line up to that cap. When streaming, we still stream all output to the user; only the buffer used for signal detection (and returned to the run-loop) is capped. A sliding-window approach: once we have read more than `max_output_buffer` bytes, keep only the last `max_output_buffer` bytes so the last line is as complete as possible within the cap.

**What to do:**

- **Interface:**  
  Extend the `Invoker` interface and `Invoke` function to accept a max-output size. For example: add a parameter `maxOutputBytes int` to `Invoke(...) (stdout []byte, exitCode int, err error)`. Semantics: `maxOutputBytes <= 0` means unlimited (current behavior). When `maxOutputBytes > 0`, the returned `stdout` must contain at most the last `maxOutputBytes` bytes of the process stdout (sliding window so the last line is preserved within the cap).
- **Implementation in `internal/backend/invoke.go`:**
  - When `maxOutputBytes <= 0`, keep current behavior (single buffer, no cap).
  - When `maxOutputBytes > 0`, use a writer that maintains a sliding window of the last `maxOutputBytes` bytes (e.g. ring buffer or discard older bytes when exceeding cap). If streaming, tee to both the stream writer and the capped buffer; the capped buffer is what is returned.
  - Return the contents of the capped buffer as `stdout`.
- **Tests:**  
  Add tests that invoke with `maxOutputBytes > 0` and large output; assert that returned stdout is at most `maxOutputBytes` long and that the last line (or last N bytes) matches the tail of the expected output.
- **Callers:**  
  Run-loop will pass `opts.Loop.MaxOutputBuffer` in T5. Review uses the same `Invoker`; pass `0` (unlimited) for review so behavior is unchanged. Update `internal/runloop/run.go` and `internal/review/run.go` (and any test adapters) to pass the new parameter; run-loop passes the configured value, review passes 0.

**Deliverable:** Backend respects `maxOutputBytes`; run-loop and review call sites updated; tests pass.

---

### T5. CLI: remove signal-precedence, wire max-output-buffer into run

**Priority:** 5 (depends on T2 and T3; and on T4 if run-loop is to pass the buffer size)

**Context:**  
CLI currently defines `--signal-precedence` and `--max-output-buffer` for `ralph run`. We are removing `--signal-precedence` entirely (hard removal). We are keeping `--max-output-buffer` and wiring it into effective config and run options so the run-loop (and thus backend) receive the configured max output buffer.

**What to do:**

- **Remove signal-precedence:**
  - In `cmd/ralph/main.go`: remove the `--signal-precedence` flag definition and the `signalPrecedence` variable (and `signalPrecedence` from the run command’s local vars and from `runLoopOverrides`).
  - In `applyRunLoopOverrides`: remove the branch that sets `out.SignalPrecedence` from `o.signalPrecedence`.
  - Remove any `ralph show config` or help output that prints `signal_precedence` (e.g. in the loop section).
- **Wire max-output-buffer:**
  - Ensure `runLoopOverrides` still has `maxOutputBuffer int` (e.g. -1 for “not set”).
  - In `applyRunLoopOverrides`: when `o.maxOutputBuffer >= 0`, set `out.MaxOutputBuffer = o.maxOutputBuffer`. When not set (-1), leave the base config value (from T3) unchanged.
  - When building `RunOptions` for the run-loop, pass `opts.Loop.MaxOutputBuffer` so the run-loop can pass it to the backend (T4). Ensure the effective loop used for the run includes the overlay from `applyRunLoopOverrides`.
  - Keep validation that `--max-output-buffer` must be >= 0 when provided.
- **Invoker signature:**  
  After T4, `Invoker.Invoke` has an extra parameter. In `cmd/ralph/main.go`, the run path does not construct the Invoker; it uses `runloop.Run` with `opts.Invoker == nil` (so `backend.Invoke` is used). The run-loop’s `Run` must pass `opts.Loop.MaxOutputBuffer` into the Invoker. So the wiring is: CLI overlay → effective `LoopSettings.MaxOutputBuffer` → `RunOptions.Loop` → run-loop → `Invoker.Invoke(..., maxOutputBytes)`.

**Deliverable:** `--signal-precedence` is gone; `--max-output-buffer` overrides config and is passed through to the run-loop/backend; `ralph show config` and help no longer mention signal_precedence.

---

### T6. Update all documentation (engineering, product, user-facing)

**Priority:** 6 (depends on T1, T2, T5; and T4 for backend/docs)

**Context:**  
Docs must reflect: (1) R008 is withdrawn/archived (T1), (2) only static precedence exists and there is no AI-interpreted option (T2, T5), (3) max_output_buffer is configurable and enforced (T3, T4, T5). Update every doc that mentions AI-interpreted precedence, signal_precedence, or max-output-buffer.

**What to do:**

- **Engineering**
  - **docs/engineering/README.md:** Run path: remove “or applies AI-interpreted precedence when configured.” Components table: run-loop row already updated in T1 (R008 removed).
  - **docs/engineering/components/run-loop.md:** One-line description: remove “or applies AI-interpreted precedence when configured.” Loop algorithm step 3: remove “or AI-interpreted when configured.” Signal detection: remove the entire “AI-interpreted precedence (O001/R008)” bullet; state that only the last non-empty line is scanned and static precedence always applies when both signals appear on that line. Add a short note that stdout capture may be capped by `max_output_buffer` (see config) so the last line is preserved up to that limit.
  - **docs/engineering/components/config.md:** Remove `signal_precedence` from the built-in defaults table and from the config file structure (loop keys). Add `max_output_buffer` to the loop section (integer, optional, bytes; default e.g. 65536). Add env var `RALPH_LOOP_MAX_OUTPUT_BUFFER` to the environment variables table.
  - **docs/engineering/components/cli.md:** Remove `--signal-precedence` from the Signals table. Keep `--max-output-buffer` in the Loop control table; clarify that it sets the max bytes retained for AI stdout (last line is preserved within this cap).
  - **docs/engineering/components/backend.md:** If it describes stdout capture, add that when `max_output_buffer` is set, only the last N bytes are retained (sliding window).
  - **docs/engineering/signal-detection-design-options.md:** Update “Current behavior” and Option F / Recommendation to state that AI-interpreted precedence is no longer offered (withdrawn). Optionally add a one-line “Withdrawn” note for Option F.

- **Product**
  - **docs/product/O001-iterative-completion/R006-signal-precedence.md:** Remove all references to “Optional AI-interpreted precedence (R008)” and “when AI-interpreted precedence (R008) is not used”; state that static precedence is the only behavior.
  - **docs/product/O001-iterative-completion/R004** and **R005:** Remove “or R008” and “(or R008 if AI-interpreted)”; reference only R006 for precedence.
  - **docs/product/O001-iterative-completion/README.md:** If there is a requirements table, remove R008 from active list and add a line that R008 is withdrawn/archived with pointer to archived doc.
  - **docs/product/O002-configurable-behavior/README.md** (and R002 if applicable): Remove “signal precedence mode (static default vs optional AI-interpreted)” and the “User enables AI-interpreted…” scenario.
  - **docs/product/O004-observability/README.md:** Remove or reword the bullet about “When AI-interpreted signal precedence is used…”.

- **User-facing**
  - **README.md:** Remove `signal_precedence: static` from the defaults list if it’s redundant (static is now the only behavior). Remove the sentence that mentions `signal_precedence: ai_interpreted` and `--signal-precedence`. Keep or add a brief note that `max_output_buffer` (and `--max-output-buffer`) cap the retained stdout for signal detection; default value as in config docs.

- **Exit codes / automation:**  
  If `docs/exit-codes.md` or automation docs mention signal precedence, update to describe only static precedence.

**Deliverable:** All listed docs updated; no references to AI-interpreted precedence as a supported option; max_output_buffer documented where relevant.

---

## Verification

- `make build`, `make test`, `make lint` pass.
- Run a quick manual test: `ralph run` with a prompt that outputs both signals on the last line → outcome is success (static precedence). No interpretation run.
- Run with `--max-output-buffer 1000` and a prompt that produces a long last line → last line is truncated to 1000 bytes in captured output (or full if under 1000); signal detection uses that captured tail.
- Config with `max_output_buffer: 65536` and no flag → run uses 65536; flag overrides config.

---

## Summary dependency order

1. **T1** – Archive R008 (product + engineering README).
2. **T2** – Remove AI-interpreted code (run-loop, interpret, config, schema, tests).
3. **T3** – Add max_output_buffer to config (file, env, defaults, schema, docs).
4. **T4** – Backend: enforce max-output cap; extend Invoker and callers.
5. **T5** – CLI: remove --signal-precedence, wire max-output-buffer into run.
6. **T6** – Update all documentation.

T1 and T2 can be done in parallel. T3 can be done in parallel with T1/T2. T4 should use the config type from T3 (or a constant default) and T5 wires CLI to config and run-loop. T6 should be last so it reflects the final behavior and flags.
