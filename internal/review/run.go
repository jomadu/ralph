package review

import "fmt"

// ErrNotImplemented indicates report generation is not yet implemented (T5.2).
// Callers should treat as exit 2.
var ErrNotImplemented = fmt.Errorf("%w: report generation not yet implemented", ErrExit2)

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
// Currently returns ErrNotImplemented until report content and apply are implemented (T5.2, T5.3, T5.4).
// Callers should check review.IsExit2(err) and exit 2 when true.
func Run(promptContent []byte, opts RunOptions) error {
	_ = promptContent
	_ = opts
	return ErrNotImplemented
}
