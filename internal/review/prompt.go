// Package review provides prompt review behavior (O5). Prompt composition for review phase (R2) and revision phase (R5).
package review

import (
	"bytes"
	_ "embed"
	"fmt"
)

//go:embed review_instructions.md
var reviewInstructionsContent []byte

//go:embed revision_instructions.md
var revisionInstructionsContent []byte

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

// ComposeRevisionPrompt builds the prompt sent to the AI for the revision phase (R5).
// Interpolates reviewOutputPath (read report from) and promptOutputPath (write revised prompt to).
func ComposeRevisionPrompt(reviewOutputPath, promptOutputPath string) ([]byte, error) {
	if reviewOutputPath == "" || promptOutputPath == "" {
		return nil, fmt.Errorf("revision phase requires both review output path and prompt output path")
	}
	pathBlock := fmt.Sprintf(
		"\n\n---\n\n- **Read the report from:** %s\n- **Write the revised prompt to:** %s\n\n---\n\n",
		reviewOutputPath, promptOutputPath,
	)
	out := bytes.Buffer{}
	out.Write(revisionInstructionsContent)
	out.WriteString(pathBlock)
	return out.Bytes(), nil
}
