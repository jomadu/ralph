# R4: Alignment of User-facing Documentation with Behavior

**Outcome:** O7 — User-facing Documentation

## Requirement

User-facing documentation is maintained so that it matches actual product behavior or explicitly states limitations, so that users are not misled by outdated or incorrect docs and the outcome "docs enable effective use" remains true over time.

## Specification

**Alignment:** The content of user-facing docs (under R1) must describe behavior that matches the product as implemented, or it must explicitly state that a feature is not yet implemented, is limited, or differs from the description (e.g. "Currently, the install script supports macOS and Linux only; Windows is not yet supported."). Docs must not describe behavior that the product does not provide without a caveat.

**Process:** The project does not prescribe a single process (e.g. mandatory review step or automation). The requirement is that alignment is maintained: when behavior changes (e.g. a new flag, a changed default, a new workaround), the corresponding user-facing doc is updated in the same change or in a follow-up that is tracked. When docs are updated, they are checked for consistency with the intent tree (where the doc traces to an outcome/requirement) so that user docs and specifications do not contradict.

**Intent tree as source of truth:** Where a user doc traces to a requirement (e.g. O3/R1 for Cursor Agent), the requirement and its specification are authoritative for *behavior*. The user doc translates that into user-oriented language; it must not contradict the requirement. If the requirement changes, the user doc must be updated to stay aligned.

**Limitations and gaps:** If the product has a known limitation (e.g. "install script not tested on Windows"), the user doc may state it so users are not misled. If a feature is not yet implemented, the doc may say so or omit the topic until it exists; it must not claim the feature exists.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| A new flag is added to the CLI (e.g. in O2/O4) | User-facing docs that describe that area (e.g. configuration or CLI usage) are updated to include or reference the new flag, or the scope of the doc is stated so that "full reference is in CLI help" is clear. |
| Intent tree requirement is updated (e.g. default value change) | Any user doc that describes that behavior is updated to match, or the doc explicitly defers to the authoritative source (e.g. "See `ralph run --help` for current flags."). |
| Doc and implementation diverge accidentally | Treated as a bug: doc or implementation is corrected so they align. No requirement for automated checks; human process (review, smoke-test, or release checklist) can enforce alignment. |
| Product has platform-specific behavior | Docs state which platforms are supported or tested so users are not misled. |

### Examples

#### Wrapper path change

**Input:** The Cursor Agent wrapper script is moved from `scripts/cursor-wrapper.sh` to `scripts/backends/cursor-wrapper.sh`.

**Expected:** User doc `docs/user/cursor-agent-workaround.md` (and any other reference to the path) is updated to the new path. O3/R1 or related intent is updated if it references the path. User docs stay aligned with where the script actually lives.

**Verification:** User following the doc finds the script at the documented path.

#### New workaround from intent

**Input:** A new backend workaround is specified in the intent tree and implemented.

**Expected:** A new user-facing topic (or section) is added (R3 coverage) and written to match the implementation. The doc does not claim behavior that the implementation does not provide.

**Verification:** User can follow the doc and succeed; doc and behavior match.

## Acceptance criteria

- [ ] User-facing docs describe behavior that matches the product or explicitly state limitations or "not yet implemented"
- [ ] When product behavior changes, the corresponding user docs are updated (in same or follow-up change)
- [ ] User docs that trace to an intent requirement do not contradict that requirement
- [ ] Known limitations (e.g. platform support) are stated where relevant so users are not misled

## Dependencies

- R1 — Alignment applies to the docs under the location and structure of R1.
- R3 — Coverage topics are the scope for alignment; each covered topic must be accurate.
- Intent tree — For traced topics, the requirement is the source of truth for behavior; user doc aligns with it.
