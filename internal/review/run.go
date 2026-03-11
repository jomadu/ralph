package review

import (
	"fmt"
	"os"
	"path/filepath"
)

// DefaultReportFilename is the default report file name when --report is omitted (T5.3, O005/R005).
const DefaultReportFilename = "ralph-review-report.txt"

// ErrNotImplemented is no longer used for report generation (T5.2 implemented).
// Kept for compatibility; callers should treat as exit 2.
var ErrNotImplemented = fmt.Errorf("%w: not implemented", ErrExit2)

// ErrReportPathRequired is kept for compatibility; Run no longer returns it when ReportPath is empty (default path is used).
var ErrReportPathRequired = fmt.Errorf("%w: report path required to write report file", ErrExit2)

// RunOptions holds options for a review run (report path, apply, etc.).
// Used by the CLI when wiring to the review component.
// When ReportPath is empty, report is written to WorkingDir/DefaultReportFilename (or cwd if WorkingDir is empty).
type RunOptions struct {
	ReportPath       string
	PromptOutputPath string
	WorkingDir       string // used for default report path when ReportPath is empty
	Apply            bool
	Yes              bool
	Quiet            bool
	LogLevel         string
}

// Run performs the review: produce report and optionally apply revision.
// T5.2: produces report file with narrative, machine-parseable summary line, and full suggested revision.
// T5.3: when ReportPath is empty, writes to WorkingDir/DefaultReportFilename (or current directory if WorkingDir empty).
// Report path must not be an existing directory (error exit 2). Apply (T5.4) and exit code derivation (T5.5) are not yet implemented.
// Callers should check review.IsExit2(err) and exit 2 when true.
func Run(promptContent []byte, opts RunOptions) error {
	reportPath := opts.ReportPath
	if reportPath == "" {
		dir := opts.WorkingDir
		if dir == "" {
			var err error
			dir, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("%w: resolving working dir for default report: %v", ErrExit2, err)
			}
		}
		reportPath = filepath.Join(dir, DefaultReportFilename)
	}
	// R005: path is a directory → error
	if fi, err := os.Stat(reportPath); err == nil && fi.IsDir() {
		return fmt.Errorf("%w: report path is a directory: %s", ErrExit2, reportPath)
	}
	report := GenerateReport(promptContent)
	body := report.String()
	if err := os.WriteFile(reportPath, []byte(body), 0644); err != nil {
		return fmt.Errorf("%w: writing report: %v", ErrExit2, err)
	}
	return nil
}
