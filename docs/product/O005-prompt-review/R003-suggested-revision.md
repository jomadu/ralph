# R003: Suggested Revision

**Outcome:** O005 — Prompt Review

## Requirement

The system produces a suggested revision of the prompt as part of every review output.

## Detail

Every successful review run includes a revised version of the prompt reflecting the review’s feedback. That revision appears in the report (e.g. `revision.md` in the report directory) and is available for the user to read or apply. The revision is Ralph’s interpretation of how to improve the prompt along the evaluation dimensions; the user may accept, edit, or ignore it.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Prompt is minimal or already strong | Revision may be small or echo improvements; still produced |
| User supplies prompt via stdin | Same revision behavior; apply path rules apply separately (R006) |

### Examples

**Input:** User runs review on a prompt missing clear success signals.

**Expected output:** Report includes narrative feedback and a concrete revised prompt addressing signals and structure.

**Verification:** Open `revision.md` (or equivalent) after review; content reflects suggested changes.
