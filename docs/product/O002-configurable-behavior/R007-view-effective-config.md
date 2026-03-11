# R007: View Effective Config

**Outcome:** O002 — Configurable Behavior

## Requirement

The user can view the effective (resolved) configuration the tool would use for the current context — current working directory, chosen prompt (if any), explicit config file, environment variables, and command-line overrides — and optionally see which layer supplied each value (defaults, global file, workspace file, explicit file, environment, prompt-level override, CLI).

## Detail

Users need to answer "what config would Ralph actually use if I ran it here?" without running the loop. With multiple layers (defaults, global, workspace, explicit file, env, prompt overrides, CLI), it can be unclear which value applies for a given setting. A read-only view of the resolved config lets users verify behavior before running, debug config issues, and script or document the effective settings.

**Context:** The view uses the same resolution as the run command (R001): current working directory (determines whether workspace config is loaded), any explicit config file path supplied, environment variables, and command-line options. When a prompt is selected (e.g. by name), the effective config includes that prompt's overrides merged with the rest. So the output reflects exactly what would be used for a run with the same invocation (same cwd, same config option, same env, same flags, same prompt choice).

**Output:** The system exposes the effective values for all configurable settings (loop behavior, prompts, AI commands as applicable). Optionally, the system can show provenance per setting (which layer supplied the value — e.g. "default", "global", "workspace", "env", "cli"). Exact format (e.g. YAML dump, key-value list, or structured output) and whether provenance is default or opt-in are implementation details; the requirement is that the user can see the resolved config and, optionally, where each value came from.

**Mechanism:** Implemented via a documented entry point (e.g. a show command or equivalent). Product does not prescribe the exact CLI shape; engineering defines the command and output format.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No config files; only defaults | View shows default values for all settings. |
| Explicit config file specified | Resolved config is from that file only (plus defaults, env, CLI); global and workspace not loaded. |
| Prompt name supplied with view | Effective config includes that prompt's overrides; root settings not overridden by that prompt appear as resolved for the run. |
| Prompt name not supplied | Effective config shows root/default loop settings; prompt-level overrides do not apply (no specific prompt chosen). |
| Missing global or workspace file | That layer is skipped; view shows resolved config as if that layer were absent (same as run). |
| Explicit config file missing | Same as run: error (R005); view does not produce resolved config when explicit file is required and missing. |
| Environment variable overrides a setting | Resolved value shown is the env value; optional provenance shows "env". |
| CLI flag overrides a setting | Resolved value shown is the CLI value; optional provenance shows "cli". |

### Examples

#### View effective config in project directory

**Input:** User is in a directory with a workspace config that sets max iterations to 5 and failure threshold to 2. No explicit config file; no env or CLI overrides. User invokes the command or option that views effective config.

**Expected output:** Output includes resolved max iterations 5, failure threshold 2, and other loop settings (from defaults or workspace). Optionally, provenance indicates "workspace" for those two and "default" for others.

**Verification:** Invoke the view-effective-config capability; confirm values match workspace config for set keys and defaults elsewhere. Run the loop; loop behavior matches the viewed config.

#### View with CLI override

**Input:** Same workspace as above. User invokes the view-effective-config capability with the documented CLI option to set max iterations to 1 for this invocation.

**Expected output:** Resolved max iterations is 1; failure threshold remains 2. Optionally, provenance shows "cli" for max iterations and "workspace" for failure threshold.

**Verification:** View output shows 1 and 2; running the loop with the same CLI flag stops after 1 iteration.

#### View with explicit config file

**Input:** User invokes the view-effective-config capability with the documented config file option pointing to a specific file (e.g. testdata/ci.yml). Global and workspace config files exist.

**Expected output:** Resolved config is from the specified file only (plus defaults, env, CLI). Values from global or workspace do not appear unless also in the explicit file.

**Verification:** Change global or workspace; view again with same explicit file; output unchanged.

## Acceptance criteria

- [ ] The user can invoke a documented command or option that outputs the effective (resolved) configuration the tool would use for the current context.
- [ ] Resolution context matches the run command: current working directory, explicit config file (if specified), environment variables, command-line options, and optional prompt selection (R001).
- [ ] Output includes effective values for configurable loop behavior (and for the chosen prompt's overrides when a prompt is specified). Prompt definitions and AI command aliases from resolved config may be included or referenced as designed.
- [ ] When an explicit config file is specified and that file is missing, the system reports an error and does not output resolved config (same as R005).
- [ ] Optionally, the user can see which layer supplied each value (provenance). If provenance is opt-in, it is documented.
- [ ] View is read-only; it does not modify config or run the loop.

## Dependencies

- R001 — Config layer resolution (view uses the same resolution as run).
- R005 — Explicit config file only (when explicit file is used, same load and error semantics).
