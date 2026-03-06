# Implementation Summary: ralph-1qi

## Task: Command source precedence ai_cmd vs ai_cmd_alias (O3/R6)

### Changes Made

1. **Config Structs** (`internal/config/config.go`)
   - Added `AICmd string` field to `LoopConfig`
   - Added `AICmd ValueWithProvenance[string]` field to `LoopConfigWithProvenance`
   - Removed default value for `AICmdAlias` (was "claude", now "")
   - Added `CommandResolution` struct to hold resolution result
   - Implemented `ResolveAICommand()` function with precedence logic:
     - Direct command (`AICmd`) takes precedence over alias
     - Returns error if neither is configured
     - Returns command string and source description

2. **Config Loading** (`internal/config/loader.go`)
   - Added `AICmd` field to `CLIFlags` struct
   - Added `AICmd` overlay in `overlayLoopConfig()`
   - Added `AICmd` overlay in `overlayLoopConfigWithMap()`
   - Wired `RALPH_LOOP_AI_CMD` environment variable in `overlayEnvironment()`
   - Added `AICmd` overlay in `OverlayCLIFlags()`

3. **Validation** (`internal/config/validate.go`)
   - Updated comment in `validateSemantic()` to clarify that command resolution errors are handled at runtime by `ResolveAICommand()`

4. **Tests** (`internal/config/resolve_test.go`)
   - Created comprehensive test suite for `ResolveAICommand()`
   - Tests cover: direct command precedence, alias resolution, no command configured, unknown alias

### Acceptance Criteria Status

- [x] If --ai-cmd is specified on the CLI, it is used regardless of any other setting
- [x] If --ai-cmd-alias is specified on the CLI (and no --ai-cmd), the alias is resolved
- [x] If neither CLI flag is specified, environment variables are checked (RALPH_LOOP_AI_CMD, RALPH_LOOP_AI_CMD_ALIAS)
- [x] If no CLI flags or environment variables are set, config file values are used (loop.ai_cmd, loop.ai_cmd_alias)
- [x] There is no built-in default for ai_cmd or ai_cmd_alias — if no layer provides a value, resolution fails
- [x] At each precedence layer, a direct command (ai_cmd) takes precedence over an alias (ai_cmd_alias)
- [x] The resolved command source is visible (returned by ResolveAICommand as CommandResolution.Source)

### Next Steps

The implementation is complete and tested. The next task (ralph-qv1) will add error handling for unknown aliases and missing commands at the CLI level.

### Build Status

✅ All packages compile successfully
✅ All tests pass (4/4 test cases)
