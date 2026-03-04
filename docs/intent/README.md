# Intent

## Outcomes

| ID | Outcome | Verification |
|----|---------|--------------|
| O1 | An AI-driven task reaches verified completion through iterative execution | User runs `ralph run <alias>`, Ralph executes fresh AI processes across iterations, detects a success signal in the output, and exits 0 |
| O2 | Loop execution adapts to the user's constraints without prompt modification | User changes iteration limits, failure thresholds, timeouts, and signal strings via config or CLI flags — the same prompt file produces different loop behavior |
| O3 | Any stdin-accepting AI CLI serves as the execution backend | User runs the same prompt with different AI CLIs by changing a config value or flag, and Ralph works with each |
| O4 | The user knows why the loop stopped and how it performed | Exit code distinguishes success (0), failure threshold (1), exhaustion (2), and interruption (130); iteration statistics are reported at completion |
