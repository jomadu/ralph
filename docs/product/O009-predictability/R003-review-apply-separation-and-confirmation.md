# R003: Review–apply separation and confirmation

**Outcome:** O009 — Predictability

## Requirement

Ralph separates review (report and suggestions only, no writes) from apply (write revision) and requires confirmation for apply in interactive use unless a non-interactive option is used.

## Detail

Review produces a report and a suggested revision; it does not modify the user's prompt file or write the revision anywhere unless the user explicitly requests apply. Apply is a distinct step: the user must request that the suggested revision be written (e.g. to the source file or a user-specified path). In interactive use, Ralph requires confirmation before performing the write (e.g. "Apply revision? [y/N]" or equivalent). If the user declines or does not confirm, no write occurs. When the user (or a script) uses a documented non-interactive option for apply (e.g. documented flag), confirmation is not required and Ralph may write without prompting. This separation prevents the user from believing they are only reviewing when a revision is actually applied.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User runs review only | Report and suggested revision are produced; no write to prompt file or revision path. No apply step is performed. |
| User runs review and then requests apply (interactive) | Ralph prompts for confirmation (e.g. "Apply? [y/N]"). If user confirms, Ralph writes; if user declines or aborts, no write. |
| User runs review with a documented non-interactive apply option (e.g. documented flag) | Ralph writes the revision without prompting; the flow is explicitly apply and non-interactive. |
| User runs review, sees "Apply?" prompt, and answers "N" or Ctrl+C | No write occurs; prompt file and any revision output path are unchanged. |
| Review output is to file (e.g. report path, revision to file); user has not requested apply | Writing the report or suggested revision to user-specified output paths (per O005) is allowed; writing the revision *into the source prompt file* or overwriting the source is apply and requires explicit request and confirmation (or non-interactive option). |

### Examples

#### Review without apply: no confirmation, no write

**Input:** User runs the review command with a prompt file (e.g. ./prompt.md). Report and suggested revision are shown or written to user-chosen report/revision paths.

**Expected output:** No prompt to "Apply?"; no write to the source prompt file. User can inspect report and revision without any change to their prompt file.

**Verification:** Source prompt file unchanged; no apply confirmation was shown.

#### Apply in interactive use: confirmation required

**Input:** User runs review, then requests that the revision be applied (e.g. selects apply in UI or runs a follow-up apply step). In an interactive session, Ralph shows a confirmation (e.g. "Apply revision to ./prompt.md? [y/N]").

**Expected output:** If user answers "y" (or equivalent), Ralph writes the revision. If user answers "n" or aborts, Ralph does not write.

**Verification:** After "n", file unchanged; after "y", file contains revision. UX clearly separates "review" from "apply" and "confirm."

#### Non-interactive apply: no confirmation prompt

**Input:** User runs the review command with the documented apply option and a non-interactive option (e.g. documented flag) in a non-interactive session.

**Expected output:** Ralph writes the revision to the appropriate path without showing an interactive confirmation prompt.

**Verification:** Revision is written; no "[y/N]" prompt; behavior is documented for scripts and automation.

#### Ambiguity avoided: user cannot accidentally apply

**Input:** User runs review and sees the suggested revision. They intend only to read the report and close.

**Expected output:** Closing or exiting without explicitly requesting apply and confirming does not result in a write. No auto-apply on exit or on "view revision."

**Verification:** Only an explicit apply request plus confirmation (or non-interactive option) triggers a write.

## Acceptance criteria

- [ ] Review (report and suggested revision) is a distinct phase from apply (writing the revision); running review alone never writes the revision to the source prompt file or user path.
- [ ] In interactive use, when the user requests apply, Ralph prompts for confirmation before writing (e.g. "Apply? [y/N]" or equivalent). If the user does not confirm, Ralph does not write.
- [ ] When a documented non-interactive apply option is used (e.g. documented flag), Ralph may write without prompting; the option is documented.
- [ ] The UX does not imply "only reviewing" while silently applying; apply is clearly a separate, confirmable step (or explicitly bypassed by a non-interactive option).

## Dependencies

- O005 (Prompt Review) defines the review command, report, and suggested revision. R003 defines the separation of review and apply and the confirmation requirement for apply. R001 defines that writes occur only when apply is requested and confirmed or non-interactive option is used.
