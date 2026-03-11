package runloop

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/maxdunn/ralph/internal/backend"
	"github.com/maxdunn/ralph/internal/config"
)

// RunOptions supplies the inputs for one run-loop execution. Caller must resolve
// config and prompt; run-loop validates command, runs the loop, and returns exit code.
type RunOptions struct {
	Command     string
	PromptBytes []byte
	Loop        config.LoopSettings
	Cwd         string
	Env         []string
	Invoker     backend.Invoker
	// Reporter receives completion message on success; nil = print to os.Stdout.
	Reporter func(msg string)
}

// invokerAdapter adapts a package-level Invoke function to backend.Invoker.
type invokerAdapter func(command string, promptBytes []byte, cwd string, env []string, timeoutSec int) (stdout []byte, exitCode int, err error)

func (f invokerAdapter) Invoke(command string, promptBytes []byte, cwd string, env []string, timeoutSec int) ([]byte, int, error) {
	return f(command, promptBytes, cwd, env, timeoutSec)
}

// Run validates the AI command, then runs the loop: for each iteration invokes
// the backend with the assembled prompt, captures stdout, and scans for the
// configured success signal. On match: reports completion (message, iteration
// count, timing) and returns ExitSuccess. When max iterations is reached
// without success, returns ExitMaxIterations. Implements T3.4, O001/R004, O004/R002.
func Run(opts RunOptions) (exitCode int, err error) {
	if opts.Invoker == nil {
		opts.Invoker = invokerAdapter(backend.Invoke)
	}
	if err := ValidateAICommand(opts.Command); err != nil {
		return ExitErrorPreLoop, err
	}
	report := opts.Reporter
	if report == nil {
		report = func(msg string) { fmt.Fprintln(os.Stdout, msg) }
	}

	start := time.Now()
	for i := 1; i <= opts.Loop.MaxIterations; i++ {
		preamble := buildPreamble(opts.Loop.Preamble, i)
		assembled := AssemblePrompt(preamble, opts.PromptBytes)
		stdout, _, invErr := opts.Invoker.Invoke(opts.Command, assembled, opts.Cwd, opts.Env, opts.Loop.TimeoutSeconds)
		if invErr != nil {
			return ExitErrorPreLoop, invErr
		}
		if ContainsSuccessSignal(stdout, opts.Loop.SuccessSignal) {
			elapsed := time.Since(start)
			report(completionMessage(i, elapsed))
			return ExitSuccess, nil
		}
	}
	report(fmt.Sprintf("Stopped after %d iteration(s) without success signal.", opts.Loop.MaxIterations))
	return ExitMaxIterations, nil
}

func buildPreamble(preamble string, iteration int) string {
	if preamble == "" {
		return ""
	}
	return "Iteration " + strconv.Itoa(iteration) + "\n" + preamble
}

func completionMessage(iterations int, elapsed time.Duration) string {
	sec := elapsed.Seconds()
	if iterations == 1 {
		return fmt.Sprintf("Completed successfully in 1 iteration (%.2fs).", sec)
	}
	return fmt.Sprintf("Completed successfully in %d iterations (%.2fs).", iterations, sec)
}
