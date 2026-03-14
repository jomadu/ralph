# R002: Inherit Env and Cwd

**Outcome:** O003 — Backend Agnosticism

## Requirement

The system inherits the user's environment and working directory when invoking the AI CLI.

## Detail

When Ralph spawns the AI CLI process, that process must run with the same environment variables and current working directory as the Ralph process (and thus the user's shell). This ensures the AI CLI can access the user's API keys, config paths, and project context without Ralph having to manage or pass through credentials. No special stripping or overriding of env or cwd is required unless specified elsewhere; the default is inheritance.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User runs Ralph from a project directory | AI CLI process has the same cwd so relative paths in prompts or tooling resolve correctly |
| User has API key or auth in environment | AI CLI process sees the same env; Ralph does not inject or remove vars |
| Ralph is invoked from a script with a modified env | AI CLI receives whatever env Ralph received |
| Working directory is removed or inaccessible after Ralph starts | Out of scope; runtime failure is reported when it occurs |

### Examples

#### Same cwd as invoker

**Input:** User runs `ralph run` from `/home/user/project` with a prompt that references `./src/file.go`.

**Expected output:** The AI CLI process is started with cwd `/home/user/project`, so `./src/file.go` resolves as the user expects.

**Verification:** Run from a known directory; inspect process or CLI behavior to confirm cwd matches.

#### Environment passed through

**Input:** User has `OPENAI_API_KEY` set in the environment and runs Ralph with an AI CLI that uses it.

**Expected output:** The AI CLI process inherits the same environment and can use `OPENAI_API_KEY`.

**Verification:** Run with a required env var unset in Ralph's env; AI CLI fails. Set it and re-run; AI CLI succeeds (subject to its own runtime behavior).

## Acceptance criteria

- [ ] When the system invokes the AI CLI process, the process's current working directory is the same as the process that started Ralph (e.g. the user's shell or script).
- [ ] When the system invokes the AI CLI process, the process's environment variables are the same as the process that started Ralph (no stripping or overriding unless specified by another requirement).

## Dependencies

- O003/R001 — The system must invoke the AI CLI; inheritance applies to that invocation.
