package review

import (
	"fmt"
	"strings"
)

// Report holds the three parts of a review report per docs/engineering/components/review.md.
type Report struct {
	// Narrative is human-readable feedback (e.g. signal discipline, statefulness, scope).
	Narrative string
	// SummaryLine is the single machine-parseable line: ralph-review: status=(ok|errors|warnings) [errors=N] [warnings=N]
	SummaryLine string
	// Revision is the full suggested prompt text.
	Revision string
}

// SummaryStatus is the status value for the machine-parseable summary line.
type SummaryStatus string

const (
	StatusOK       SummaryStatus = "ok"
	StatusErrors   SummaryStatus = "errors"
	StatusWarnings SummaryStatus = "warnings"
)

// FormatSummaryLine builds the canonical summary line.
// Format: ralph-review: status=(ok|errors|warnings)(\s+errors=(\d+))?(\s+warnings=(\d+))?
func FormatSummaryLine(status SummaryStatus, errors, warnings int) string {
	s := fmt.Sprintf("ralph-review: status=%s", status)
	if errors >= 0 {
		s += fmt.Sprintf(" errors=%d", errors)
	}
	if warnings >= 0 {
		s += fmt.Sprintf(" warnings=%d", warnings)
	}
	return s
}

// String returns the report as the report file contents: narrative, summary line, then full revision.
// Structure: narrative block, blank line, summary line, separator, revision.
func (r *Report) String() string {
	var b strings.Builder
	if r.Narrative != "" {
		b.WriteString(r.Narrative)
		b.WriteString("\n\n")
	}
	b.WriteString(r.SummaryLine)
	b.WriteString("\n\n---\n\n")
	b.WriteString(r.Revision)
	if r.Revision != "" && !strings.HasSuffix(r.Revision, "\n") {
		b.WriteString("\n")
	}
	return b.String()
}

// GenerateReport produces a report from prompt content.
// T5.2: produces narrative, machine-parseable summary line, and full suggested revision.
// T5.6 (O005/R007): evaluates prompt on four dimensions (signal/state, iteration awareness, scope/convergence, subjective completion) and structures narrative by dimension.
func GenerateReport(promptContent []byte) *Report {
	dimensions := evaluateDimensions(promptContent)
	narrative := narrativeFromDimensions(dimensions)
	revision := suggestedRevisionFromDimensions(promptContent, dimensions)

	errorsCount := 0
	for _, d := range dimensions {
		if !d.OK {
			errorsCount++
		}
	}
	status := StatusOK
	if errorsCount >= 1 {
		status = StatusErrors
	}
	warningsCount := 0
	summaryLine := FormatSummaryLine(status, errorsCount, warningsCount)

	return &Report{
		Narrative:   narrative,
		SummaryLine: summaryLine,
		Revision:    revision,
	}
}
