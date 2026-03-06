# Using the Ralph Binary and Build Prompt for Prompt Review

This doc explains how to **build and run the ralph binary** and how to use the **build procedure prompt** (`prompts/build.md`) to burn down the **Prompt Review (O5)** tasks.

---

## Building the Ralph Binary

From the repository root:

```bash
go build -o ralph ./cmd/ralph
```

This produces a `ralph` executable in the current directory.

### Putting ralph on your PATH

Copy the binary into a directory that is **already on your PATH**. If you copy into a directory that isn’t on PATH, `ralph` will not be found when you run it by name.

- **Prefer a directory you already use:** Many setups add `~/.local/bin` to PATH (e.g. in `.zshrc`). If yours does, use that:
  ```bash
  mkdir -p ~/.local/bin
  cp -f ralph ~/.local/bin/ralph
  chmod +x ~/.local/bin/ralph
  ```
- **Using `~/bin`:** On macOS with zsh, `~/bin` is often **not** on PATH by default. If you copy there, you must add it yourself, e.g. in `~/.zshrc`:
  ```bash
  export PATH="$HOME/bin:$PATH"
  ```
  Then run `source ~/.zshrc` or open a new terminal.

When the install script exists (O6), you can use `./scripts/install.sh` instead; it will place the binary and (per the plan) record the install location for uninstall.

**Verify:**

```bash
ralph version
# ralph 0.1.0
```

If you get `command not found: ralph`, the directory you copied into is not on your PATH. Either copy `ralph` into a directory that is (e.g. `~/.local/bin` if that’s in your PATH), or add the directory you used to PATH in your shell config.

**Existing subcommands (O1–O4):**

- `ralph run [alias]` — run the loop with a prompt (alias, `-f <path>`, or stdin)
- `ralph run -f <path>` — run with prompt from file
- `ralph list prompts` — list configured prompt aliases
- `ralph list aliases` — list AI command aliases
- `ralph version` — print version

The **Prompt Review** feature adds `ralph review` (see O5 specs). That subcommand is not yet implemented; the first task (T1 / ralph-e8a) adds it.

---

## The Build Procedure Prompt (`prompts/build.md`)

`prompts/build.md` is an **executable procedure** for an AI agent. It defines:

1. **Scope:** One task per iteration. The agent picks the single next most important task, completes it, updates work tracking, then signals or continues.
2. **Phases:** OBSERVE → ORIENT → DECIDE → ACT.
3. **Signals:**
   - `<promise>SUCCESS</promise>` — all tasks done or no ready work.
   - `<promise>FAILURE</promise>` — blocked (missing info, tools, or permissions).
   - No signal — one task completed but more ready work remains (loop runs again).

The procedure tells the agent to:

- Read **AGENTS.md** (work tracking, build/test/lint, specs, implementation).
- Query **ready work** via `bd ready --json`.
- Study **specs** in `docs/intent/` and **implementation** in `cmd/`, `internal/`, `scripts/`.
- **Pick one task** from ready work (priority, dependencies, impact).
- **Implement** it, run tests if any, **update work tracking** (e.g. `bd close <id> --reason "Completed"`), and commit.

So you use the build prompt by **running the ralph loop** with that prompt: the loop runs the procedure, the agent does one task, then the loop runs again until SUCCESS or FAILURE.

---

## How to Use the Build Prompt to Burn Down Prompt Review Tasks

### 1. Run Ralph with the Build Prompt

Use the **build** prompt as the task. You can point ralph at it by alias or by file.

**By alias (if `build` is configured in `ralph-config.yml`):**

```bash
ralph run build
```

**By file (no config needed):**

```bash
ralph run -f prompts/build.md
```

**From repo root**, ensure a config exists that defines the `build` alias, for example in `./ralph-config.yml` or `~/.config/ralph/ralph-config.yml`:

```yaml
prompts:
  build:
    path: prompts/build.md
loop:
  ai_cmd_alias: claude   # or kiro, cursor-agent, etc.
```

Then `ralph run build` will run the loop with `prompts/build.md` as the task prompt.

### 2. What the Loop Does

Each iteration:

1. Ralph assembles the prompt (with optional preamble) and sends it to the AI (e.g. Claude, Kiro).
2. The AI follows the procedure: observes (AGENTS.md, `bd ready`, specs, implementation), orients, decides, picks **one** task, then acts (code, tests, `bd close`, commit).
3. Ralph scans the AI output for `<promise>SUCCESS</promise>` or `<promise>FAILURE</promise>`.
4. If SUCCESS or FAILURE, the loop stops. Otherwise it runs another iteration with the same build prompt (so the agent can pick the next task).

So **burning down** Prompt Review tasks means: run `ralph run build` (or `ralph run -f prompts/build.md`) and let the loop repeatedly execute the procedure until all ready work is done or the agent signals failure.

### 3. Prompt Review Tasks and Order

Ready work is determined by **bd (beads)**. For O5 Prompt Review, the plan is in `PLAN_PROMPT_REVIEW.md` and the first **unblocked** task is:

| Order | Task | Bead ID   | Description |
|-------|------|-----------|-------------|
| T1    | First | **ralph-e8a** | Review subcommand and R1 input modes (alias / file / stdin) |
| T2    |      | ralph-wrd | R3 review output path |
| T3    |      | ralph-orr | R2 review prompt composition |
| T4    |      | ralph-5vr | R8 failure handling |
| T5    |      | ralph-5b8 | R9 report file verification |
| T6    |      | ralph-xkg | R6 report format and exit codes |
| T7    |      | ralph-bgo | R4 prompt output path |
| T8    |      | ralph-1a3 | R5 apply and revision phase |
| T9    |      | ralph-b1e | R7 configurable review stdout |

T1 (ralph-e8a) has no dependencies, so it appears in `bd ready`. After T1 is closed, T2, T3, T4, T7 become ready, and so on.

**Check ready work anytime:**

```bash
bd ready --json
```

**Claim and close (when doing work manually):**

```bash
bd update ralph-e8a --claim --json
# ... do the work ...
bd close ralph-e8a --reason "Completed" --json
```

When using the **build prompt**, the agent is expected to claim/close the task it picks (per the procedure).

### 4. One-Shot vs Loop

- **Loop (recommended for burning down many tasks):**  
  `ralph run build` (or `ralph run -f prompts/build.md`). Each iteration = one task; repeat until SUCCESS or FAILURE.

- **One-shot (single task):**  
  Run the AI once with the contents of `prompts/build.md` and instruct it to do exactly one iteration (pick one task from `bd ready`, complete it, then output SUCCESS or continue). No need to run the full ralph loop if you only want one task done.

### 5. After Prompt Review Is Implemented

Once `ralph review` exists, you can also **review** the build prompt itself:

```bash
ralph review build
# or
ralph review -f prompts/build.md
```

That will run the O5 reviewer on `prompts/build.md` and produce a report (and optional apply).

---

## Summary

| Goal | Command / Action |
|------|-------------------|
| Build ralph | `go build -o ralph ./cmd/ralph` |
| Put on PATH | Copy to a directory already on PATH (e.g. `~/.local/bin`); avoid `~/bin` unless you’ve added it to PATH in your shell config. Verify with `ralph version`. |
| Run build procedure (burn down tasks) | `ralph run build` or `ralph run -f prompts/build.md` |
| Check ready O5 tasks | `bd ready --json` |
| Close a task | `bd close <id> --reason "Completed"` |

The **build prompt** is the procedure; the **ralph binary** runs the loop that executes it; **bd** holds the Prompt Review (and other) issues. Use the loop with the build prompt to burn down Prompt Review tasks one at a time until no ready work remains or the agent signals failure.
