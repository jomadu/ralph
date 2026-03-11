package review

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DefaultReportFilename is the default report file name when --report is omitted (T5.3, O005/R005).
const DefaultReportFilename = "ralph-review-report.txt"

// ErrNotImplemented is no longer used for report generation (T5.2 implemented).
// Kept for compatibility; callers should treat as exit 2.
var ErrNotImplemented = fmt.Errorf("%w: not implemented", ErrExit2)

// ErrReportPathRequired is kept for compatibility; Run no longer returns it when ReportPath is empty (default path is used).
var ErrReportPathRequired = fmt.Errorf("%w: report path required to write report file", ErrExit2)

// ErrApplyPromptOutputRequired is returned when apply is requested but no revision output path is available (e.g. alias with no path, or stdin already validated by CLI).
var ErrApplyPromptOutputRequired = fmt.Errorf("%w: --apply requires --prompt-output when revision output path cannot be defaulted", ErrExit2)

// RunOptions holds options for a review run (report path, apply, etc.).
// Used by the CLI when wiring to the review component.
// When ReportPath is empty, report is written to WorkingDir/DefaultReportFilename (or cwd if WorkingDir is empty).
// For apply: when prompt is from file/alias, SourcePath can default the revision path; otherwise --prompt-output is required.
// NonInteractive is true when stdin is not a TTY (e.g. CI); when true and overwrite would need confirmation, Run returns ErrApplyConfirmationRequired unless Yes is set.
// Quiet minimizes output (O004/R006); Verbose adds revision path when applying.
type RunOptions struct {
	ReportPath       string
	PromptOutputPath string
	SourcePath       string // when prompt from file or alias with path; used to default revision path when Apply and PromptOutputPath empty
	WorkingDir       string // used for default report path when ReportPath is empty
	Apply            bool
	Yes              bool
	NonInteractive   bool
	Verbose          bool
	Quiet            bool
	LogLevel         string
}

// Run performs the review: produce report and optionally apply revision.
// T5.2: produces report file with narrative, machine-parseable summary line, and full suggested revision.
// T5.3: when ReportPath is empty, writes to WorkingDir/DefaultReportFilename (or current directory if WorkingDir empty).
// T5.4: when Apply is true, writes revision to chosen path; interactive confirm before overwrite unless --yes; non-interactive without --yes exits 2.
// T5.5: after report is written, parses machine-parseable summary to derive exit code 0 vs 1; returns that code on success. Report write failure or apply/precondition → error (caller exits 2).
// Report path must not be an existing directory (error exit 2). Callers should check review.IsExit2(err) and exit 2 when true; on nil error, exit with the returned code (0 or 1).
func Run(promptContent []byte, opts RunOptions) (exitCode int, err error) {
	reportPath := opts.ReportPath
	if reportPath == "" {
		dir := opts.WorkingDir
		if dir == "" {
			var e error
			dir, e = os.Getwd()
			if e != nil {
				return 0, fmt.Errorf("%w: resolving working dir for default report: %v", ErrExit2, e)
			}
		}
		reportPath = filepath.Join(dir, DefaultReportFilename)
	}
	// R005: path is a directory → error
	if fi, e := os.Stat(reportPath); e == nil && fi.IsDir() {
		return 0, fmt.Errorf("%w: report path is a directory: %s", ErrExit2, reportPath)
	}
	report := GenerateReport(promptContent)
	body := report.String()
	if e := os.WriteFile(reportPath, []byte(body), 0644); e != nil {
		return 0, fmt.Errorf("%w: writing report: %v", ErrExit2, e)
	}
	if !opts.Quiet {
		fmt.Fprintf(os.Stderr, "Report written to %s\n", reportPath)
	}
	if opts.Apply {
		if e := applyRevision(report.Revision, opts); e != nil {
			return 0, e
		}
		if opts.Verbose {
			applyPath := opts.PromptOutputPath
			if applyPath == "" {
				applyPath = opts.SourcePath
			}
			if applyPath != "" && opts.WorkingDir != "" && !filepath.IsAbs(applyPath) {
				applyPath = filepath.Join(opts.WorkingDir, applyPath)
			}
			if applyPath != "" {
				fmt.Fprintf(os.Stderr, "Revision applied to %s\n", applyPath)
			}
		}
	}
	// T5.5: parse written report body to derive exit code 0 vs 1 (O005/R008, O010/R003).
	status, errorsCount, _, ok := ParseSummaryFromReport([]byte(body))
	if !ok {
		// Missing or malformed summary → 1 (fail-safe for CI).
		return 1, nil
	}
	return ExitCodeFromSummary(status, errorsCount), nil
}

// applyRevision writes the revision to the chosen path with confirmation when overwriting (O005/R004, O009/R001, O009/R003).
func applyRevision(revision string, opts RunOptions) error {
	applyPath := opts.PromptOutputPath
	if applyPath == "" {
		applyPath = opts.SourcePath
	}
	if applyPath == "" {
		return ErrApplyPromptOutputRequired
	}
	// Resolve relative to working dir if needed
	if !filepath.IsAbs(applyPath) && opts.WorkingDir != "" {
		applyPath = filepath.Join(opts.WorkingDir, applyPath)
	}
	fi, err := os.Stat(applyPath)
	if err == nil && fi.IsDir() {
		return fmt.Errorf("%w: revision output path is a directory: %s", ErrExit2, applyPath)
	}
	exists := err == nil
	if exists {
		if !opts.Yes {
			if opts.NonInteractive {
				return ErrApplyConfirmationRequired
			}
			allowed, err := confirmOverwrite(applyPath)
			if err != nil {
				return fmt.Errorf("%w: %v", ErrExit2, err)
			}
			if !allowed {
				return nil // user declined; no write, exit 0 (review completed)
			}
		}
	}
	if err := os.WriteFile(applyPath, []byte(revision), 0644); err != nil {
		return fmt.Errorf("%w: writing revision: %v", ErrExit2, err)
	}
	return nil
}

// confirmOverwrite prompts on stdout/stderr and reads from stdin. Returns true if user confirms overwrite, false otherwise.
func confirmOverwrite(path string) (bool, error) {
	fmt.Fprintf(os.Stderr, "Overwrite %s? [y/N] ", path)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false, scanner.Err()
	}
	line := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return line == "y" || line == "yes", nil
}
