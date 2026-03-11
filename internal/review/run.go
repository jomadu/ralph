package review

import (
	"fmt"
	"os"
)

// ErrNotImplemented is no longer used for report generation (T5.2 implemented).
// Kept for compatibility; callers should treat as exit 2.
var ErrNotImplemented = fmt.Errorf("%w: not implemented", ErrExit2)

// ErrReportPathRequired indicates ReportPath was empty; report file could not be written.
var ErrReportPathRequired = fmt.Errorf("%w: report path required to write report file", ErrExit2)

// RunOptions holds options for a review run (report path, apply, etc.).
// Used by the CLI when wiring to the review component.
type RunOptions struct {
	ReportPath       string
	PromptOutputPath string
	Apply            bool
	Yes              bool
	Quiet            bool
	LogLevel         string
}

// Run performs the review: produce report and optionally apply revision.
// T5.2: produces report file with narrative, machine-parseable summary line, and full suggested revision.
// Writes to opts.ReportPath when set; returns ErrReportPathRequired when ReportPath is empty.
// Apply (T5.4) and exit code derivation (T5.5) are not yet implemented.
// Callers should check review.IsExit2(err) and exit 2 when true.
func Run(promptContent []byte, opts RunOptions) error {
	if opts.ReportPath == "" {
		return ErrReportPathRequired
	}
	report := GenerateReport(promptContent)
	body := report.String()
	if err := os.WriteFile(opts.ReportPath, []byte(body), 0644); err != nil {
		return fmt.Errorf("%w: writing report: %v", ErrExit2, err)
	}
	return nil
}
