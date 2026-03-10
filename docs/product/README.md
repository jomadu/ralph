# Product

Ralph is a loop runner for AI-driven tasks. The user supplies a prompt and chooses an AI CLI (e.g. Claude, Cursor); Ralph re-invokes it until a success or failure signal appears in the output or a limit is reached. Multi-step work that would otherwise require manual read–judge–re-run cycles can reach verified completion without hand-holding. The outcomes below are the measurable results of that product.

## Outcomes

| ID | Outcome |
|----|---------|
| [O001](./O001-iterative-completion/README.md) | An AI-driven task reaches verified completion through iterative execution |
| [O002](./O002-configurable-behavior/README.md) | Loop execution adapts to the user's constraints without changing the prompt file |
| [O003](./O003-backend-agnosticism/README.md) | Any stdin-accepting AI CLI serves as the execution backend |
| [O004](./O004-observability/README.md) | The user understands why the loop or review stopped, how it performed, and what to do when something fails |
| [O005](./O005-prompt-review/README.md) | Prompts can be reviewed for quality and structure before or without running the loop; user gets a report and optional revision; `--apply` writes revision to the prompt file |
| [O006](./O006-install-uninstall/README.md) | Users can install Ralph on their system and uninstall it cleanly |
| [O007](./O007-user-documentation/README.md) | Users have documentation that explains how to use the product |
| [O008](./O008-discoverability/README.md) | A new user can discover what Ralph does and get to a first successful run |
| [O009](./O009-predictability/README.md) | Ralph only changes user content when the user explicitly requests it (e.g. `--apply`) |
| [O010](./O010-automation/README.md) | Users can run Ralph from scripts and CI |
| [O011](./O011-upgrade-path/README.md) | Users can upgrade Ralph without breaking existing config or workflows |
