# Plan: O6 Install and Uninstall

This plan implements **O6 — Install and Uninstall** so users can install Ralph on their system (invocable from a shell) and uninstall it cleanly. Outcome and requirements are in `docs/intent/O6-install-uninstall/`.

**Verification (from O6 README):** User runs the install script; opens a new shell and runs `ralph version` (exit 0). User runs the uninstall script; `ralph` is no longer found; no broken PATH or leftover install artifacts.

**Scope:** Install and uninstall **scripts only** (no Homebrew or other package managers). Scripts live in `scripts/`: `scripts/install.sh` and `scripts/uninstall.sh`. The only artifact installed is the `ralph` binary. User config (e.g. `~/.config/ralph/ralph-config.yml`) is never removed by uninstall.

---

## Task dependency overview

```
T1 (install state and target convention)
  → T2 (install script)
  → T4 (documentation)
  → T3 (uninstall script)
```

T2 and T3 can be implemented in parallel after T1. T4 (documentation) should reflect the actual script behavior and therefore depends on T2 (and ideally T3) being done.

---

## T1: Install state and target directory convention

**Priority:** 1 (must be first)  
**Dependencies:** None  
**Spec:** [R1 — Install artifact and location](docs/intent/O6-install-uninstall/R1-install-artifact-and-location.md), [R4 — Uninstall behavior](docs/intent/O6-install-uninstall/R4-uninstall-behavior.md)

**Objective:** Define where the install script will place the binary and how the uninstall script will discover that location, so both scripts share a single contract. No code yet — this is a design decision that T2 and T3 will implement.

**Context:**

- **Artifact:** Only the `ralph` binary is installed (single executable). No user config is installed; user config paths (e.g. `~/.config/ralph/ralph-config.yml`) are out of scope for uninstall.
- **Target directory:** Choose a default install directory that is commonly on PATH. Typical options: `~/bin` (user-writable, often on PATH) or `/usr/local/bin` (may require sudo). The spec allows the script to "prompt the user for a directory" or use a default; document the default and any override (e.g. environment variable `RALPH_INSTALL_DIR` or `--dir`).
- **Install state:** The uninstall script must know where the binary was placed (so it can remove only that file and not touch other copies). Use a **state file** that the install script writes and the uninstall script reads. Suggested location: `~/.config/ralph/install-state` (or a single line in a file like `~/.config/ralph/.install-path`). Content: the directory path where `ralph` was installed (so uninstall can `rm "$INSTALL_DIR/ralph"`). If the install script supports a user-specified directory, that path must be what is recorded. Do **not** store user config in this file; this file is only for "where did we put the binary."
- **Idempotency:** Uninstall script should tolerate "binary already deleted" or "state file missing" (treat as "not installed" and exit 0 or report clearly). Install script overwriting an existing install (same or different dir) is acceptable: update state file and replace binary.

**Decisions to document in this task (for T2/T3):**

- Default install directory (e.g. `$HOME/bin` if it exists and is writable, else `/usr/local/bin` with a note about sudo; or always `$HOME/bin` and document "ensure ~/bin is on PATH").
- State file path and format (e.g. `~/.config/ralph/install-state` with one line: `INSTALL_DIR=/path/to/dir` or just `/path/to/dir`).
- Override mechanism for install directory (env var and/or flag).

**Acceptance:**

- [ ] Default install directory is chosen and documented in this plan or a short design note.
- [ ] State file path and format are defined; install script will write it; uninstall script will read it and remove only the binary at that path, then remove the state file.
- [ ] User config files are never written or deleted by these scripts; state file is separate from `ralph-config.yml`.

---

## T2: Install script (`scripts/install.sh`)

**Priority:** 2  
**Dependencies:** T1  
**Spec:** [R2 — Supported install methods](docs/intent/O6-install-uninstall/R2-supported-install-methods.md), [R3 — Post-install invokability](docs/intent/O6-install-uninstall/R3-post-install-invokability.md)

**Objective:** Implement `scripts/install.sh` that obtains the Ralph binary, places it in the target directory (per T1), and records the install path so uninstall can remove it. After running the script, the user can open a new shell and run `ralph version` (binary on PATH).

**Context:**

- **Obtaining the binary:** When run from a Ralph repo that contains Go code (e.g. `cmd/ralph/main.go`), build the binary: `go build -o ralph ./cmd/ralph` (from repo root). Optionally support an explicit path to a pre-built binary (e.g. `RALPH_BINARY=/path/to/ralph ./scripts/install.sh`) for CI or release workflows. If neither repo nor pre-built path is available, exit with a clear error (e.g. "Run from Ralph repo root or set RALPH_BINARY").
- **Target directory:** Use the convention from T1: default directory, overridable by env (e.g. `RALPH_INSTALL_DIR`) or flag (e.g. `--dir /path`). Ensure the directory exists (mkdir -p); if it is not writable (e.g. `/usr/local/bin`), print instructions for sudo or suggest `~/bin`.
- **Copy binary:** Copy (or move) the binary to `$INSTALL_DIR/ralph`. Use `cp -f` (per AGENTS.md) to overwrite if present. Set executable bit if needed (`chmod +x "$INSTALL_DIR/ralph"`).
- **State file:** Write the chosen install directory (and only that) to the state file path defined in T1. Create `~/.config/ralph` if needed for the state file; do not create or modify `ralph-config.yml`.
- **Platform:** Script should run on macOS and Linux (sh or bash). Detect OS if needed (e.g. for future Windows branch); document "tested on macOS, Linux" and "Windows not yet supported" in T4.
- **Non-interactive:** Prefer non-interactive behavior where possible. If prompting (e.g. "Install to ~/bin or /usr/local/bin?"), support an env var or flag to skip (e.g. `RALPH_INSTALL_DIR=~/bin` or `--dir ~/bin`).
- **Output:** Print where the binary was installed and remind the user to open a new terminal or ensure that directory is on PATH. Example: "Installed ralph to ~/bin. Run 'ralph version' in a new terminal to verify."

**Acceptance:**

- [ ] From repo root (with Go), running `./scripts/install.sh` builds the binary (if needed), installs it to the default directory, and writes the state file.
- [ ] With `RALPH_INSTALL_DIR` (or equivalent) set, the script installs to that directory and records it in the state file.
- [ ] After install, the binary at `$INSTALL_DIR/ralph` exists and is executable; `ralph version` works in a new shell if that dir is on PATH.
- [ ] Script uses `-f` for copy/move and avoids interactive prompts when a default or env is provided.
- [ ] Clear error message when not run from repo and no pre-built binary path is provided.

---

## T3: Uninstall script (`scripts/uninstall.sh`)

**Priority:** 3  
**Dependencies:** T1 (T2 not strictly required to implement T3, but uninstall is only useful after install exists)

**Spec:** [R4 — Uninstall behavior](docs/intent/O6-install-uninstall/R4-uninstall-behavior.md)

**Objective:** Implement `scripts/uninstall.sh` that removes the binary installed by the install script and the install state file. It must not remove user config or leave broken references (e.g. symlinks or PATH entries pointing at the removed binary). Idempotent when already uninstalled.

**Context:**

- **Discover install location:** Read the state file from the path defined in T1. If the state file does not exist, treat as "Ralph not installed by this method" — print a short message (e.g. "Ralph does not appear to be installed (no install state found).") and exit 0 (no error) or exit with a distinct code; do not remove arbitrary paths.
- **Remove binary:** Delete `$INSTALL_DIR/ralph` (from state file). Use `rm -f` so missing binary does not fail. Do not remove other files in `$INSTALL_DIR`.
- **Remove state file:** After removing the binary, delete the state file so a future install gets a clean state. If the state file is in `~/.config/ralph/`, do not remove the directory or `ralph-config.yml`; only remove the state file (e.g. `install-state` or `.install-path`).
- **Idempotency:** If state file is missing: exit 0 with "not installed" message. If state file exists but binary already missing: remove state file and exit 0 (cleanup).
- **No broken references:** The install script does not modify PATH or create symlinks elsewhere; uninstall only deletes the one binary and the state file. So no extra cleanup is required beyond that.

**Acceptance:**

- [ ] Running `./scripts/uninstall.sh` after a T2 install removes the binary at the recorded path and removes the state file.
- [ ] User config (`~/.config/ralph/ralph-config.yml` or workspace config) is never deleted.
- [ ] If state file is missing, script reports "not installed" and exits 0 (or documented non-zero); does not delete anything outside the recorded path.
- [ ] If binary was already manually deleted, uninstall still removes the state file and exits 0.
- [ ] Script uses `rm -f`; no interactive prompts.

---

## T4: Install and uninstall documentation

**Priority:** 4  
**Dependencies:** T2 (and T3 recommended so uninstall is documented accurately)  
**Spec:** [R5 — Install and uninstall documentation](docs/intent/O6-install-uninstall/R5-install-uninstall-documentation.md)

**Objective:** Add a single authoritative reference for installing and uninstalling Ralph so users can follow consistent steps. Content must cover: how to run the install script, where the binary is placed, how to ensure invokability (PATH), how to verify (`ralph version`), how to run the uninstall script, what is removed and what is not (user config), and which platforms are supported.

**Context:**

- **Location:** Either a new **Install** (and **Uninstall**) section in the main `README.md`, or a dedicated `docs/INSTALL.md` with a link from the README. The spec allows either; choose one and keep it as the single source of truth.
- **Install section:** Prerequisites (e.g. Go if building from source; or "pre-built binary"); how to run the script (from repo: `./scripts/install.sh`; or via curl if you add a stable URL later); optional override (e.g. `RALPH_INSTALL_DIR=~/bin`); where the binary goes (default per T1); note that user may need to "open a new terminal" or add the install dir to PATH; verification: run `ralph version` and expect exit 0.
- **Uninstall section:** Run `./scripts/uninstall.sh` (or from repo path); what is removed (the binary and install state file); what is not removed (user config in `~/.config/ralph/` or workspace); no broken PATH/symlinks.
- **Platform support:** State that the scripts are tested on macOS and Linux (and list arch if relevant, e.g. amd64/arm64). State if Windows is not yet supported.
- **Edge cases (optional):** If the user moved the binary after install, uninstall still cleans the recorded path and state file; the moved copy is out of scope. Multiple installs (e.g. different dirs): last install wins for state file; uninstall removes that one.

**Acceptance:**

- [ ] There is one clear place (README section or docs/INSTALL.md) that describes install and uninstall.
- [ ] Install steps include: how to run the script, default (or chosen) install directory, PATH note, and verification (`ralph version`).
- [ ] Uninstall steps include: how to run the uninstall script, what is removed (binary + state file), what is not removed (user config).
- [ ] Supported platforms (e.g. macOS, Linux) are stated; Windows gap is stated if applicable.
- [ ] Documentation matches the actual behavior of the scripts (T2 and T3).

---

## Reference: O6 requirement summary

| ID  | Requirement | Delivered by |
|-----|-------------|--------------|
| R1  | Install artifact and location definition | T1 (convention), T2 (script + state file), T4 (docs) |
| R2  | Supported install methods (install/uninstall scripts) | T2, T3 |
| R3  | Post-install invokability | T2 (binary on PATH), T4 (PATH + verify step) |
| R4  | Uninstall behavior (no broken state, no user config removed) | T1 (state), T3 (script), T4 (docs) |
| R5  | Install and uninstall documentation | T4 |

---

## Beads (bd) issue mapping (badges)

| Plan task | Bead ID    | Dependencies (blocked by) |
|-----------|-------------|---------------------------|
| T1        | ralph-7zf   | —                         |
| T2        | ralph-5vv   | ralph-7zf                 |
| T3        | ralph-ubw   | ralph-7zf                 |
| T4        | ralph-aux   | ralph-5vv, ralph-ubw      |

**Check ready work:** `bd ready` or `bd ready --json`. **Claim:** `bd update <id> --claim`. **Close:** `bd close <id> --reason "Completed"`.
