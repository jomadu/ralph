// Package runloop implements the run-loop component: validate AI command,
// load prompt once, invoke backend, detect signals, and exit with documented codes.
package runloop

// Run exit codes (run-loop is the authority; see docs/engineering/components/run-loop.md).
const (
	// ExitSuccess is returned when the success signal is detected.
	ExitSuccess = 0
	// ExitErrorPreLoop is returned before the loop starts when the AI command
	// is missing, invalid, or not executable (O001/R001, O004/R001).
	ExitErrorPreLoop = 2
	// ExitMaxIterations is returned when the loop reaches max iterations without
	// detecting the success signal (O001/R007, O004/R004). Exact value TBD in user docs.
	ExitMaxIterations = 3
)
