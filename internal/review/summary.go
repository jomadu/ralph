// Package review: parser for the machine-parseable summary line (T5.5, T5.7).
// See docs/engineering/components/review.md and O005/R008, O010/R003.
package review

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"
)

// Summary line format: ralph-review: status=(ok|errors|warnings)(\s+errors=(\d+))?(\s+warnings=(\d+))?
var summaryLineRE = regexp.MustCompile(`ralph-review:\s*status=(ok|errors|warnings)(\s+errors=(\d+))?(\s+warnings=(\d+))?`)

// ParseSummaryFromReport reads the report content and finds the first line matching
// the canonical summary format. Returns status, errors count, warnings count, and
// whether a valid line was found. If no line matches or the line is malformed,
// ok is false (caller should treat as exit 1 per fail-safe).
func ParseSummaryFromReport(reportContent []byte) (status SummaryStatus, errorsCount, warningsCount int, ok bool) {
	scanner := bufio.NewScanner(strings.NewReader(string(reportContent)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		status, errorsCount, warningsCount, ok = ParseSummaryLine(line)
		if ok {
			return status, errorsCount, warningsCount, true
		}
	}
	return "", 0, 0, false
}

// ParseSummaryLine parses a single line in the format:
//
//	ralph-review: status=(ok|errors|warnings)(\s+errors=N)?(\s+warnings=N)?
//
// Returns status, errors count (default 0 if absent), warnings count (default 0 if absent), and ok.
func ParseSummaryLine(line string) (status SummaryStatus, errorsCount, warningsCount int, ok bool) {
	matches := summaryLineRE.FindStringSubmatch(line)
	if matches == nil {
		return "", 0, 0, false
	}
	status = SummaryStatus(matches[1])
	if matches[3] != "" {
		errorsCount, _ = strconv.Atoi(matches[3])
	}
	if matches[5] != "" {
		warningsCount, _ = strconv.Atoi(matches[5])
	}
	return status, errorsCount, warningsCount, true
}

// ExitCodeFromSummary derives process exit code from parsed summary (T5.5, O005/R008).
// 0 = review completed, no errors; 1 = review completed, prompt has errors or missing/malformed summary; 2 is never returned here (caller uses ErrExit2).
// Missing or malformed summary → 1 (fail-safe for CI).
func ExitCodeFromSummary(status SummaryStatus, errorsCount int) int {
	if status == StatusOK && errorsCount == 0 {
		return 0
	}
	// status=errors or errors>=1 → 1; status=warnings with errors=0 → 0 per spec; missing/malformed → 1
	if status == StatusErrors || errorsCount >= 1 {
		return 1
	}
	if status == StatusWarnings {
		return 0
	}
	// Unknown status or empty → fail-safe 1
	return 1
}
