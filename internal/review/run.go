package review

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/maxdunn/ralph/internal/backend"
)

// DefaultReportDir is the default report directory name when --report is omitted (O005/R005).
const DefaultReportDir = "ralph-review"

// Report filenames written by the AI into the report directory.
const (
	ReportResultJSON = "result.json"
	ReportSummaryMD  = "summary.md"
	ReportOriginalMD = "original.md"
	ReportRevisionMD = "revision.md"
	ReportDiffMD     = "diff.md"
)

// ErrNotImplemented is no longer used for report generation (T5.2 implemented).
// Kept for compatibility; callers should treat as exit 2.
var ErrNotImplemented = fmt.Errorf("%w: not implemented", ErrExit2)

// ErrReportPathRequired is kept for compatibility; Run no longer returns it when ReportPath is empty (default path is used).
var ErrReportPathRequired = fmt.Errorf("%w: report path required to write report file", ErrExit2)

// ErrApplyPromptOutputRequired is returned when apply is requested but no revision output path is available (e.g. alias with no path, or stdin already validated by CLI).
var ErrApplyPromptOutputRequired = fmt.Errorf("%w: --apply requires --prompt-output when revision output path cannot be defaulted", ErrExit2)

// ErrAICommandRequired is returned when review is run without an AI command (empty Command).
var ErrAICommandRequired = fmt.Errorf("%w: review requires an AI command (set loop.ai_cmd_alias or --ai-cmd-alias)", ErrExit2)

// invokerAdapter adapts a function to backend.Invoker (same pattern as runloop).
type invokerAdapter func(command string, promptBytes []byte, cwd string, env []string, timeoutSec int, streamTo io.Writer) (stdout []byte, exitCode int, err error)

func (f invokerAdapter) Invoke(command string, promptBytes []byte, cwd string, env []string, timeoutSec int, streamTo io.Writer) ([]byte, int, error) {
	return f(command, promptBytes, cwd, env, timeoutSec, streamTo)
}

// RunOptions holds options for a review run (report path, apply, backend, etc.).
// Used by the CLI when wiring to the review component.
// When ReportPath is empty, report is written to WorkingDir/DefaultReportFilename (or cwd if WorkingDir is empty).
// For apply: when prompt is from file/alias, SourcePath can default the revision path; otherwise --prompt-output is required.
// NonInteractive is true when stdin is not a TTY (e.g. CI); when true and overwrite would need confirmation, Run returns ErrApplyConfirmationRequired unless Yes is set.
// Quiet minimizes output (O004/R006); Verbose adds revision path when applying.
// Review requires an AI command: Command must be non-empty; Invoker is used to run it (nil = backend.Invoke). Cwd, Env, TimeoutSec passed to the backend.
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
	// Backend (required for agent-based review)
	Command    string
	Invoker    backend.Invoker // nil = use backend.Invoke
	Cwd        string
	Env        []string
	TimeoutSec int
}

// resultJSON is the shape of result.json in the report directory (written by the AI).
type resultJSON struct {
	Status   string `json:"status"`
	Errors   int    `json:"errors"`
	Warnings int    `json:"warnings"`
}

// Run performs the review: invoke backend with review prompt (including report dir path);
// the AI creates five files in that directory. Run then reads result.json for exit code
// and optionally applies revision from revision.md. ReportPath is the report directory
// (default WorkingDir/DefaultReportDir when empty). Directory is created if missing.
func Run(promptContent []byte, opts RunOptions) (exitCode int, err error) {
	if strings.TrimSpace(opts.Command) == "" {
		return 0, ErrAICommandRequired
	}
	invoker := opts.Invoker
	if invoker == nil {
		invoker = invokerAdapter(backend.Invoke)
	}

	// Resolve report directory (absolute) and create it if needed.
	reportDir := opts.ReportPath
	if reportDir == "" {
		dir := opts.WorkingDir
		if dir == "" {
			var e error
			dir, e = os.Getwd()
			if e != nil {
				return 0, fmt.Errorf("%w: resolving working dir for default report: %v", ErrExit2, e)
			}
		}
		reportDir = filepath.Join(dir, DefaultReportDir)
	}
	if !filepath.IsAbs(reportDir) {
		base := opts.WorkingDir
		if base == "" {
			base, _ = os.Getwd()
		}
		reportDir = filepath.Join(base, reportDir)
	}
	reportDirAbs, err := filepath.Abs(reportDir)
	if err != nil {
		return 0, fmt.Errorf("%w: resolving report dir: %v", ErrExit2, err)
	}
	if fi, e := os.Stat(reportDirAbs); e == nil && !fi.IsDir() {
		return 0, fmt.Errorf("%w: report path is not a directory: %s", ErrExit2, reportDirAbs)
	}
	if e := os.MkdirAll(reportDirAbs, 0755); e != nil {
		return 0, fmt.Errorf("%w: creating report dir: %v", ErrExit2, e)
	}

	reviewPrompt := AssembleReviewPrompt(promptContent, reportDirAbs)
	stdout, _, invErr := invoker.Invoke(opts.Command, reviewPrompt, opts.Cwd, opts.Env, opts.TimeoutSec, nil)
	if invErr != nil {
		return 0, fmt.Errorf("%w: backend invocation failed: %v", ErrExit2, invErr)
	}
	_ = stdout // AI responds with short confirmation; we read result from files

	resultPath := filepath.Join(reportDirAbs, ReportResultJSON)
	data, err := os.ReadFile(resultPath)
	if err != nil {
		return 0, fmt.Errorf("%w: reading %s: %v", ErrExit2, resultPath, err)
	}
	var result resultJSON
	if err := json.Unmarshal(data, &result); err != nil {
		return 0, fmt.Errorf("%w: invalid %s: %v", ErrExit2, ReportResultJSON, err)
	}
	status := SummaryStatus(strings.TrimSpace(result.Status))
	if status != StatusOK && status != StatusErrors && status != StatusWarnings {
		return 0, fmt.Errorf("%w: invalid status in %s: %q", ErrExit2, ReportResultJSON, result.Status)
	}

	if !opts.Quiet {
		fmt.Fprintf(os.Stderr, "Report written to %s\n", reportDirAbs)
	}
	if opts.Apply {
		revisionPath := filepath.Join(reportDirAbs, ReportRevisionMD)
		revision, e := os.ReadFile(revisionPath)
		if e != nil {
			return 0, fmt.Errorf("%w: reading %s for apply: %v", ErrExit2, ReportRevisionMD, e)
		}
		if e := applyRevision(string(revision), opts); e != nil {
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
	return ExitCodeFromSummary(status, result.Errors), nil
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
