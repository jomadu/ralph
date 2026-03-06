# Cursor Agent: optional wrapper for plain-text stdout

**Intent:** [O3 — Backend agnosticism](../intent/O3-backend-agnosticism/README.md), [R1 — Built-in command aliases](../intent/O3-backend-agnosticism/R1-builtin-aliases.md)

The built-in `cursor-agent` alias resolves to the raw command:

```text
agent -p --force --output-format stream-json --stream-partial-output
```

That command works from any working directory (no script path dependency). The Cursor Agent writes **newline-delimited JSON** to stdout, not plain text. You can either:

- Configure **signal strings** to match content inside those JSON lines (so the loop can detect SUCCESS/FAILURE in the stream), or  
- Use an **optional wrapper script** that invokes the same command, parses the JSON, and emits plain text to stdout so Ralph’s normal signal scanning works.

This page describes the optional wrapper.

## What the wrapper does

The wrapper script (in the Ralph repo at `scripts/cursor-wrapper.sh`):

1. Reads the prompt from stdin and pipes it to `agent -p --force --output-format stream-json --stream-partial-output`.
2. Parses each JSON line (e.g. with `jq`), extracts text from `assistant` message content, and writes that text to **stdout**.
3. Sends tool-call progress and other metadata to **stderr**.
4. Ensures Ralph’s signal scanning sees plain text on stdout, so default signal strings work without JSON-aware patterns.

**Dependencies:** The script requires `jq` and the `agent` binary on your PATH. If either is missing, the script exits non-zero and prints an error to stderr.

## How to use the wrapper

The wrapper is **not** the default for the `cursor-agent` alias (the default is the raw command). To use the wrapper:

1. Ensure `scripts/cursor-wrapper.sh` is available (e.g. clone the Ralph repo or copy the script to a directory you control).
2. Point Ralph at the script via **direct command** or a **user-defined alias**:

   **Option A — config, user-defined alias (override built-in):**

   ```yaml
   ai_cmd_aliases:
     cursor-agent: "/path/to/ralph/scripts/cursor-wrapper.sh"
   ```

   Use a path that works from the directory where you run `ralph` (absolute path is most reliable).

   **Option B — CLI:**

   ```bash
   ralph run build --ai-cmd "/path/to/ralph/scripts/cursor-wrapper.sh"
   ```

3. Run as usual, e.g. `ralph run build --ai-cmd-alias cursor-agent` (if you overrode the alias in config) or with `--ai-cmd` as above.

**Note:** The script is not an embedded resource; it must exist on disk at the path you give. If you run Ralph from a different project, use an absolute path or a path relative to your current working directory that reaches the script.
