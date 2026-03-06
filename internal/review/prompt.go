// Package review provides prompt review behavior (O5). Prompt composition for review phase (R2).
package review

import (
	"bytes"
	_ "embed"
	"fmt"
)

//go:embed review_instructions.md
var reviewInstructionsContent []byte

// ComposeReviewPrompt builds the prompt sent to the AI for the review phase (R2).
// Order: (1) embedded review instructions, (2) path directive with reportPath, (3) user prompt content.
// reportPath is the R3 review output path; the AI is instructed to write the report to this path.
func ComposeReviewPrompt(reportPath string, userContent []byte) ([]byte, error) {
	if reportPath == "" {
		return nil, fmt.Errorf("review report path is required for prompt composition")
	}
	pathDirective := fmt.Sprintf("\n\n---\n\nWrite your report to the following path (you must write the full report to this file):\n\n%s\n\n---\n\nUser prompt to review:\n\n", reportPath)
	out := bytes.Buffer{}
	out.Write(reviewInstructionsContent)
	out.WriteString(pathDirective)
	out.Write(userContent)
	return out.Bytes(), nil
}
