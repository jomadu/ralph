# Ralph prompt review instructions

You are reviewing a prompt that will be used with Ralph: a loop runner that pipes the prompt to an AI CLI once per iteration and scans output for success/failure signals.

## Your task

1. **Evaluate** the prompt for qualities that support Ralph's execution model:
   - **Signal discipline**: Does the prompt tell the AI what to emit so Ralph can detect success or failure (e.g. `<promise>SUCCESS</promise>` / `<promise>FAILURE</promise>`)?
   - **Statefulness**: Does it acknowledge that each iteration is a fresh process and that state lives on the filesystem?
   - **Scope and convergence**: Is the task scoped so the AI can complete it in one or a few iterations? Are "done" criteria clear?
   - **Iteration awareness**: Does the prompt account for iteration count, failure threshold, or loop context where relevant?

2. **Produce** a single report file containing:
   - **Narrative feedback** (what works, what to improve, risks).
   - **Machine-parseable summary**: Include a line of the form `ralph-review: status=ok` or `ralph-review: status=errors` (and optionally `errors=N`, `warnings=N`) so scripts can gate on the result.
   - **Full suggested revision**: The complete revised prompt text (so the user or an apply step can use it).

You must write the entire report to the exact path you are given below. Do not write it to stdout only; Ralph reads the report from that file.
