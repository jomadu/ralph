# R001: Review invocation inputs

**Outcome:** O005 — Prompt Review

## Requirement

The user can invoke the review command with the prompt supplied by alias (resolved from config), file path, or standard input.

## Detail

The review command accepts the prompt to be reviewed from exactly one of: an alias name (resolved to a prompt source per config), a user-chosen path to a prompt file, or when the prompt is supplied via standard input. The system reads the prompt content from that source and passes it to the reviewer. Resolution of alias and validation of path occur before the review runs; invalid or missing sources produce a clear error and do not produce a report or revision.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Alias not defined in config | Fail before review; clear error that alias is missing or invalid. |
| Alias points to missing file | Fail before review; clear error that the prompt source is not found. |
| File path does not exist or is not readable | Fail before review; clear error; no report or revision produced. |
| Both file path and standard input supplied (e.g. CLI and pipe) | Behavior is defined (e.g. one takes precedence or usage error); no silent ambiguity. |
| Prompt supplied via standard input but empty | System either treats as empty prompt and reviews it, or errors with clear message; behavior is defined. |
| Valid alias, file, or standard input | Prompt content is read and review runs; report and suggested revision are produced. |

### Examples

#### Invoke with file path

**Input:** User runs the review command with the prompt supplied by file path (e.g. a path to `my-prompt.md`) and that file exists and is readable.

**Expected output:** Review runs; user receives a report and a suggested revision (e.g. report to default or specified path, revision in output or buffer).

**Verification:** Report and revision reflect the content of the supplied file; exit code per R008.

#### Invoke with alias

**Input:** Config defines a prompt alias (e.g. "default") that resolves to a prompt file; user runs the review command with that alias. The resolved file exists.

**Expected output:** Review runs using the content of the resolved prompt; report and suggested revision are produced.

**Verification:** Output is based on that file's content; no prompt is read from standard input.

#### Invoke when prompt supplied via standard input

**Input:** User runs the review command with the prompt supplied via standard input (e.g. piped content).

**Expected output:** Review runs using the supplied content; report and suggested revision are produced. If user later requests that the revision be written (R004), they must specify the revision output path (R006).

**Verification:** Report and revision match the supplied content; apply without revision output path yields error per R006.

## Acceptance criteria

- [ ] User can supply the prompt for review via a config alias, a file path, or standard input; the system accepts exactly one source per invocation.
- [ ] When the source is invalid (alias missing, file not found or not readable), the system fails before running the review and emits a clear error; no report or revision is produced.
- [ ] When the source is valid, the system reads the prompt content and runs the review, producing a report and a suggested revision per R002 and R003.
- [ ] Behavior when both file path and standard input are supplied is defined and documented (precedence or usage error).

## Dependencies

None.
