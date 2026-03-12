// Package review: assemble the review prompt sent to the backend (instructions + user prompt).
// See docs/engineering/components/review.md "Review invocation" and O005/R007.
package review

import "strings"

// PlaceholderReportDir is replaced by the actual report directory path when assembling the prompt.
const PlaceholderReportDir = "{{REPORT_DIR}}"

// Embedded review instructions (Ralph-owned). Evaluates the user prompt along the four
// dimensions (O005/R007). The AI must create five files in the given directory and
// respond only with a short confirmation. reportDir is interpolated from run options.
const reviewInstructionsTemplate = `You are a prompt reviewer for Ralph, an iterative AI execution loop. Evaluate the user prompt below and create five files in the directory we specify. Do not put the file contents in your response. Create the files on disk, then respond briefly that you have created the files and where they are.

**Dimensions to evaluate:**
1. **Signal and state** — Does the prompt define clear success and failure signals Ralph can detect? Is statefulness compatible with a fresh process per iteration?
2. **Iteration awareness** — Does the prompt acknowledge multi-iteration execution with a fresh process each time, so the AI re-reads state and emits signals each run?
3. **Scope and convergence** — Does the task have defined scope and checkable completion criteria so the loop can converge?
4. **Subjective completion criteria** — If "done" is subjective (e.g. "good enough"), does the prompt include escape techniques (variation, step back, challenge assumptions) to avoid getting stuck?

**Directory where you must create the files:**  
This is the report output directory (by default a directory named ralph-review in the project). It is a real directory in the project, not a temporary location. Create the files there.

` + PlaceholderReportDir + `

**Files to create in that directory (create these files; do not output their contents in your reply):**

1. **result.json** — JSON object with the review outcome only. Keys: "status" (string: "ok" | "errors" | "warnings"), optional "errors" (number), optional "warnings" (number). No narrative. Example: {"status":"ok","errors":0,"warnings":0}

2. **summary.md** — Human-readable narrative feedback (e.g. by dimension). Markdown. This is the main feedback the user reads.

3. **original.md** — The exact prompt you were given to review, unchanged. Markdown.

4. **revision.md** — Your full suggested revised prompt: the complete prompt as you would have the user use it. Markdown. This is what gets applied if the user runs apply.

5. **diff.md** — The diff between the original and the revision: unified diff in a code block, or prose describing changes. Markdown.

**What to do:** Create the five files above in the directory path shown. Then respond with a short message only, for example: "Created the review report at <path>. Files: result.json, summary.md, original.md, revision.md, diff.md." Do not include the contents of the files in your response.

Here is the user prompt to review:

---
`

// AssembleReviewPrompt returns the bytes to send to the backend: instructions with
// reportDir interpolated, plus the user prompt content. reportDir should be the
// absolute path to the report directory (from run options) so the AI knows where
// to create the five files.
func AssembleReviewPrompt(userPrompt []byte, reportDir string) []byte {
	instructions := strings.ReplaceAll(reviewInstructionsTemplate, PlaceholderReportDir, reportDir)
	return []byte(instructions + string(userPrompt))
}
