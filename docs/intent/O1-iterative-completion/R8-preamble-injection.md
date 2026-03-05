# R8: Preamble Injection

**Outcome:** O1 — Iterative Completion

## Requirement

The system wraps the user's prompt with loop state — iteration number, iteration limit, and optional user-provided context — before piping the assembled content to the AI CLI. The preamble provides dynamic metadata that the static prompt file cannot know on its own. It is generated per iteration (iteration number changes), but the underlying prompt content is immutable across iterations.

## Specification

Before each iteration, Ralph assembles the input to pipe to the AI CLI by concatenating a generated preamble with the buffered prompt content. The preamble provides dynamic loop state that the static prompt file cannot know.

**Assembled input structure:**

```
<preamble>\n\n<prompt_content>
```

Two newlines separate the preamble from the prompt content.

**Preamble format (without context):**

```
[RALPH] Iteration <N> of <limit>
```

Where `<N>` is the current iteration number (starting from 1) and `<limit>` is the value of `default_max_iterations`. In unlimited mode, `<limit>` is the string `unlimited`:

```
[RALPH] Iteration <N> of unlimited
```

**Preamble format (with context):**

When one or more `--context` flags are provided, a CONTEXT section is appended:

```
[RALPH] Iteration <N> of <limit>

CONTEXT:
<context_content>
```

Multiple `--context` flag values are concatenated with double newlines between them, in the order provided on the command line.

**Configuration:**

- Field: `preamble`
- Type: boolean
- Default: `true`
- Can be set globally (`loop.preamble: false`) or per-prompt (`prompts.<name>.loop.preamble: false`)
- When `false`, the assembled input is just `<prompt_content>` with no preamble or separator

**Context flag:**

- CLI: `--context <string>` (repeatable)
- No config file equivalent — context is inherently per-invocation
- When preamble is disabled and `--context` is provided, Ralph emits a warning that context will not be included

**Immutability:**

The prompt content portion is immutable across iterations (R9). The preamble is regenerated each iteration because the iteration number changes. The assembled input therefore changes between iterations only in the preamble's iteration number.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Preamble disabled | Prompt content piped directly; no `[RALPH]` header, no separator |
| Preamble enabled, no context | Preamble is the single `[RALPH] Iteration N of M` line, followed by two newlines, followed by prompt content |
| Preamble enabled, one `--context` flag | Preamble includes CONTEXT section with the provided string |
| Preamble enabled, multiple `--context` flags | Context values concatenated with double newlines, in command-line order |
| Preamble disabled, `--context` provided | Warning emitted; context is not included; prompt piped directly |
| Unlimited mode | Preamble reads `Iteration N of unlimited` |
| Iteration 1 of 1 | Preamble reads `Iteration 1 of 1` |

### Examples

#### Standard preamble, iteration 3 of 10

**Input:**
`default_max_iterations: 10`, `preamble: true`, no context. Prompt content is `Fix the failing tests.\n`.

**Expected output:**
Assembled input piped to AI CLI:

```
[RALPH] Iteration 3 of 10

CONTEXT:
Fix the failing tests.
```

**Verification:**
- AI CLI's stdin begins with `[RALPH] Iteration 3 of 10`
- Prompt content follows after two newlines

#### Preamble with context

**Input:**
`default_max_iterations: 5`, `preamble: true`, `--context "Focus on the auth module" --context "Ignore deprecation warnings"`. Prompt content is `Refactor for clarity.\n`.

**Expected output:**
Assembled input:

```
[RALPH] Iteration 1 of 5

CONTEXT:
Focus on the auth module

Ignore deprecation warnings

Refactor for clarity.
```

**Verification:**
- Both context strings appear in order between the iteration line and the prompt content

#### Preamble disabled

**Input:**
`preamble: false`. Prompt content is `Generate a README.\n`.

**Expected output:**
Assembled input:

```
Generate a README.
```

**Verification:**
- No `[RALPH]` header in the AI CLI's stdin
- Prompt content is piped as-is

## Acceptance criteria

- [ ] The preamble includes the current iteration number and the max iteration limit (or "unlimited" when in unlimited mode)
- [ ] When context is provided via --context flag(s), the preamble includes a CONTEXT section with the provided content
- [ ] When no context is provided, the CONTEXT section is omitted entirely
- [ ] Preamble injection is enabled by default
- [ ] Preamble can be disabled globally (loop.preamble: false) or per-prompt (prompts.<name>.loop.preamble: false)
- [ ] When preamble is disabled, the prompt file content is piped directly to the AI CLI without any wrapping

## Dependencies

_None identified._
