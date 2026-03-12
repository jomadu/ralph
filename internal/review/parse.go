// Package review: parse AI stdout into Report (narrative, summary line, revision).
// See docs/engineering/components/review.md "Expected AI output format".
package review

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
)

// ErrParseAIOutput is returned when AI stdout cannot be parsed into a valid report.
var ErrParseAIOutput = errors.New("review: could not parse AI output into report (missing or malformed summary line or revision separator)")

// ParseAIOutput parses backend stdout into a Report.
// Expects: optional narrative, one line matching ralph-review: status=..., then a line
// containing "---", then the full suggested revision. Returns ErrParseAIOutput if
// the summary line or separator is missing.
func ParseAIOutput(stdout []byte) (*Report, error) {
	lines := splitLines(stdout)
	var summaryIdx int = -1
	var separatorIdx int = -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if _, _, _, ok := ParseSummaryLine(trimmed); ok {
			summaryIdx = i
			break
		}
	}
	if summaryIdx < 0 {
		return nil, ErrParseAIOutput
	}
	for i := summaryIdx + 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" || strings.HasPrefix(strings.TrimSpace(lines[i]), "---") {
			separatorIdx = i
			break
		}
	}
	if separatorIdx < 0 {
		return nil, ErrParseAIOutput
	}

	narrative := strings.TrimSpace(strings.Join(lines[0:summaryIdx], "\n"))
	summaryLine := strings.TrimSpace(lines[summaryIdx])
	revision := strings.TrimSpace(strings.Join(lines[separatorIdx+1:], "\n"))
	// Normalize: ensure revision ends with newline when non-empty for Report.String()
	if revision != "" && !strings.HasSuffix(revision, "\n") {
		revision += "\n"
	}

	return &Report{
		Narrative:   narrative,
		SummaryLine: summaryLine,
		Revision:    revision,
	}, nil
}

func splitLines(b []byte) []string {
	var out []string
	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		out = append(out, scanner.Text())
	}
	return out
}
