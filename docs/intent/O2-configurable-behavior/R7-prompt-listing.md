# R7: Resource Listing Command

**Outcome:** O2 — Configurable Behavior

## Requirement

The system provides a command to list configured resources — prompt aliases and AI command aliases — so users can discover what is available without reading config files. Both built-in and user-defined resources are shown.

## Specification

Ralph provides `ralph list` with two subcommands for discovering configured resources. Both commands load and merge configuration (R5, R1) and run config validation (R3) before producing output, but they do not validate prompt file existence (R4) since no prompt is being executed.

### `ralph list prompts`

Lists all prompt aliases defined in the resolved configuration. In default mode, this is global + workspace merged. When `--config <path>` is provided, only the explicit file is used (R5).

**Output format:**

Output is YAML, one entry per alias, sorted alphabetically by alias key:

```yaml
bootstrap:
  path: ./prompts/bootstrap.md
build:
  name: Build
  description: Run the main build loop
  path: ./prompts/build.md
```

**Rules:**

- Each top-level key is the prompt alias name.
- `name` is included only if explicitly set in config (omitted when it would just repeat the alias key).
- `description` is included only if explicitly set in config.
- `path` is always included, shown as written in config (not resolved to absolute).
- Entries are sorted alphabetically by alias key.
- Output goes to stdout.
- If no prompts are configured, output a single line: `No prompts configured.`

### `ralph list aliases`

Lists all AI command aliases — both built-in and user-defined — with their resolved commands.

**Built-in aliases:**

| Alias | Command |
|-------|---------|
| `claude` | `claude -p --dangerously-skip-permissions` |
| `kiro` | `kiro-cli chat --no-interactive --trust-all-tools` |
| `copilot` | `copilot --yolo` |
| `cursor-agent` | `cursor-agent -p -f --stream-partial-output --output-format stream-json` |

**Merge behavior:**

User-defined aliases in `ai_cmd_aliases` are merged with built-in aliases. User-defined aliases with the same name as a built-in override the built-in command. User-defined aliases with new names are added to the set.

**Output format:**

Output is YAML, one entry per alias, sorted alphabetically:

```yaml
claude:
  command: claude -p --dangerously-skip-permissions
copilot:
  command: copilot --yolo
cursor-agent:
  command: cursor-agent -p -f --stream-partial-output --output-format stream-json
kiro:
  command: kiro-cli chat --no-interactive --trust-all-tools
my-tool:
  command: my-custom-cli --flag
```

**Rules:**

- Each top-level key is the alias name. Each entry has a `command` field with the resolved command string.
- Entries are sorted alphabetically by alias name.
- Output goes to stdout.
- Built-in aliases are always shown, even if no config files exist.

### General behavior

- Both subcommands exit with code 0 on success.
- Both subcommands exit with code 1 if config loading or validation fails (R3, R5).
- `ralph list` with no subcommand prints usage help listing available subcommands.
- These commands do not start the loop. They are informational only.

### Edge cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No prompts configured in any config file | `ralph list prompts` outputs: `No prompts configured.` Exit 0. |
| No config files exist at all | `ralph list prompts` outputs: `No prompts configured.` `ralph list aliases` shows only built-in aliases. Exit 0. |
| Prompt defined in global config but not workspace config | Prompt appears in listing — global config prompts are included |
| Same prompt alias defined in both global and workspace config (no `--config`) | Workspace definition wins (field-level merge per R6); listing shows the resolved values |
| `--config` provided | Only prompts from the explicit file appear; global and workspace configs are not consulted |
| User-defined alias overrides built-in `claude` | `ralph list aliases` shows `claude` with the user-defined command string |
| User-defined alias with a new name not in built-ins | Alias appears in the listing with its command string |
| Config validation fails (e.g., invalid `log_level`) | Both list commands exit with code 1 after printing validation errors. No listing is shown. |
| Prompt path points to a nonexistent file | Listing still shows the prompt — `ralph list prompts` does not validate file existence. The `path` field shows whatever path is configured. |

### Examples

#### Listing prompts with mixed config

**Input:**
Workspace config:
```yaml
prompts:
  build:
    path: "./prompts/build.md"
    name: "Build"
    description: "Run the main build loop"
  bootstrap:
    path: "./prompts/bootstrap.md"
```

**Expected output:**
```yaml
bootstrap:
  path: ./prompts/bootstrap.md
build:
  name: Build
  description: Run the main build loop
  path: ./prompts/build.md
```

**Verification:**
- Both prompts appear sorted alphabetically
- `bootstrap` omits `name` and `description` (neither was set in config)
- `build` shows its configured name and description

#### Listing aliases with a user override

**Input:**
Workspace config:
```yaml
ai_cmd_aliases:
  claude: "claude --no-permissions"
  my-tool: "my-custom-cli --flag"
```

**Expected output:**
```yaml
claude:
  command: claude --no-permissions
copilot:
  command: copilot --yolo
cursor-agent:
  command: cursor-agent -p -f --stream-partial-output --output-format stream-json
kiro:
  command: kiro-cli chat --no-interactive --trust-all-tools
my-tool:
  command: my-custom-cli --flag
```

**Verification:**
- `claude` shows the user-defined command (override took effect)
- `my-tool` appears with its command
- All built-in aliases are present

#### No config files, listing aliases

**Input:**
No config files exist. User runs `ralph list aliases`.

**Expected output:**
YAML output showing all four built-in aliases with their command strings.

**Verification:**
- Four aliases appear: `claude`, `copilot`, `cursor-agent`, `kiro`
- Exit code is 0

## Acceptance criteria

- [ ] `ralph list prompts` outputs all prompt aliases defined in the resolved configuration
- [ ] Each prompt entry shows the alias key, display name (if set), description (if set), and prompt file path
- [ ] If no prompts are configured, the output clearly indicates that no prompts are available
- [ ] `ralph list aliases` outputs all AI command aliases (built-in and user-defined, merged)
- [ ] Each AI command alias entry shows the alias name and the resolved command string
- [ ] Output is YAML for both subcommands

## Dependencies

_None identified._
