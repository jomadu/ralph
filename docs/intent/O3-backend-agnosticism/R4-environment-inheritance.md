# R4: Process Environment Inheritance

**Outcome:** O3 — Backend Agnosticism

## Requirement

The system passes the current process's environment variables and working directory to each spawned AI CLI process. The AI CLI operates in the same environment as the user who invoked Ralph, with access to the same API keys, PATH, and filesystem context.

## Specification

When Ralph spawns the AI CLI process for an iteration, it uses **exec-family** semantics (e.g. `execve`): the child process is created with an explicit argv (from R2 parsing) and an explicit environment. This requirement specifies that environment and working directory.

**Environment:** The child process's environment must be a copy of Ralph's process environment. Every environment variable present in Ralph's process (e.g. `PATH`, `HOME`, `OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, user-defined vars) is present in the child with the same name and value. Ralph must not modify, filter, or sanitize the environment — no removal of "dangerous" variables, no override of `PATH` or `HOME`, no injection of Ralph-specific vars unless required for correct operation (and none are specified here). The intent is that the AI CLI sees exactly what the user's shell would pass to it if the user ran the same command manually.

**Working directory:** The child process's current working directory must be the same as Ralph's current working directory at the time of the spawn. So if the user ran `ralph run build` from `/home/user/project`, each AI CLI process is also started with cwd `/home/user/project`. Relative paths in the prompt or in the AI's tool use resolve relative to that directory.

**When it applies:** This applies to every AI CLI process Ralph starts for the loop (each iteration). It does not apply to any helper processes Ralph might start (e.g. for validation); only to the process that receives the prompt on stdin and whose stdout is scanned for signals.

**No shell:** Because Ralph does not invoke a shell (R2, O3 non-outcomes), the environment is not "whatever the shell would set" — it is Ralph's own environment. Typically Ralph is invoked from a shell, so Ralph's environment is the shell's environment at exec time. If Ralph is invoked by a process that does not pass the user's environment, the AI CLI will not see it either; that is outside Ralph's control.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User sets `ANTHROPIC_API_KEY` in shell, then runs `ralph run build --ai-cmd-alias claude` | AI CLI process has `ANTHROPIC_API_KEY` in its environment with the same value. |
| User runs Ralph from `/project`; prompt says "edit ./README.md" | AI CLI runs with cwd `/project`; `./README.md` resolves to `/project/README.md`. |
| Ralph's process has 100 env vars | Child gets the same 100 vars (or equivalent). No filtering. |
| User runs `env -i ralph run build` (empty env) | Child gets empty (or minimal) environment. Ralph does not add vars. |
| Ralph changes cwd internally between iterations | Each spawn uses Ralph's cwd at the time of that spawn. Ralph must not change cwd between iterations unless specified elsewhere; if cwd is unchanged, every iteration sees the same cwd. |

### Examples

#### API key available to CLI

**Input:**
```bash
export ANTHROPIC_API_KEY=sk-...
ralph run build --ai-cmd-alias claude
```

**Expected output:**
The spawned `claude` process has `ANTHROPIC_API_KEY` in its environment. The CLI can authenticate without the user embedding the key in the command or config.

**Verification:**
- Run with a valid key; loop proceeds. Or inspect process (e.g. `/proc/<pid>/environ` on Linux) to confirm var is present.

#### Working directory

**Input:**
User in `/home/user/myapp` runs `ralph run build`. Prompt says "list files in the current directory."

**Expected output:**
AI CLI's cwd is `/home/user/myapp`. Listing "current directory" shows contents of `myapp`.

**Verification:**
- AI output or tool calls reference paths under `/home/user/myapp`. No confusion with Ralph's binary directory or temp dirs.

## Acceptance criteria

- [ ] The AI CLI process inherits all environment variables from Ralph's process
- [ ] The AI CLI process's working directory is the same as Ralph's working directory
- [ ] Environment variables set by the user before invoking Ralph (API keys, PATH additions, config paths) are available to the AI CLI
- [ ] Ralph does not modify, filter, or sanitize the inherited environment

## Dependencies

_None identified._
