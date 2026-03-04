# R4: Process Environment Inheritance

**Outcome:** O3 — Backend Agnosticism

## Requirement

The system passes the current process's environment variables and working directory to each spawned AI CLI process. The AI CLI operates in the same environment as the user who invoked Ralph, with access to the same API keys, PATH, and filesystem context.

## Specification

_To be specified in Phase 3._

## Acceptance criteria

- [ ] The AI CLI process inherits all environment variables from Ralph's process
- [ ] The AI CLI process's working directory is the same as Ralph's working directory
- [ ] Environment variables set by the user before invoking Ralph (API keys, PATH additions, config paths) are available to the AI CLI
- [ ] Ralph does not modify, filter, or sanitize the inherited environment
