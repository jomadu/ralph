// Package review provides prompt review behavior (O5). Report summary parsing for exit code derivation (R6).
package review

import (
	"os"
	"regexp"
)

// SummaryLinePrefix is the exact prefix the AI must emit for the machine-parseable summary (R6).
// Full format: ralph-review: status=ok|errors|warnings [errors=N] [warnings=N]
const SummaryLinePrefix = "ralph-review:"

// summaryLineRegex matches the canonical summary line (R6).
// Example: ralph-review: status=ok
// Example: ralph-review: status=errors errors=2
// Example: ralph-review: status=warnings warnings=1 errors=0
var summaryLineRegex = regexp.MustCompile(
	`ralph-review:\s*status=(ok|errors|warnings)(\s+errors=(\d+))?(\s+warnings=(\d+))?`,
)

// ParseReportSummary reads the report file at reportPath, finds the machine-parseable summary line (R6),
// and returns the exit code Ralph should use: 0 (no errors / only warnings) or 1 (one or more errors).
// Call only after VerifyReportExists has succeeded. Missing or malformed summary → 1 (fail-safe for CI).
func ParseReportSummary(reportPath string) (exitCode int) {
	data, err := os.ReadFile(reportPath)
	if err != nil {
		return 1
	}
	return deriveExitCodeFromSummary(string(data))
}

// deriveExitCodeFromSummary scans content for a line matching the summary format and returns 0 or 1.
func deriveExitCodeFromSummary(content string) int {
	// Search line by line for ralph-review: ...
	for _, line := range lineSlice(content) {
		loc := summaryLineRegex.FindStringSubmatch(line)
		if loc == nil {
			continue
		}
		status := loc[1]
		errorsVal := loc[3] // capture group for errors=N
		_ = loc[5]          // warnings=N unused for exit code policy
		switch status {
		case "ok":
			if errorsVal == "" || errorsVal == "0" {
				return 0
			}
			return 1
		case "errors":
			return 1
		case "warnings":
			if errorsVal == "" || errorsVal == "0" {
				return 0
			}
			return 1
		default:
			return 1
		}
	}
	return 1
}

func lineSlice(s string) []string {
	var lines []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	return lines
}
