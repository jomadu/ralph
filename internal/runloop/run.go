package runloop

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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
	// DryRun when true: assemble prompt (with preamble if enabled), print to stdout, exit 0.
	// No backend invocation (T3.10, O004/R007).
	DryRun bool
	// Reporter receives completion message on success; nil = print to os.Stdout.
	Reporter func(msg string)
	// InterruptContext if non-nil is used for interrupt detection (e.g. in tests);
	// when it is cancelled, Run returns ExitInterrupt. If nil, Run uses
	// signal.NotifyContext(background, os.Interrupt).
	InterruptContext context.Context
}

// invokerAdapter adapts a package-level Invoke function to backend.Invoker.
type invokerAdapter func(command string, promptBytes []byte, cwd string, env []string, timeoutSec int) (stdout []byte, exitCode int, err error)

func (f invokerAdapter) Invoke(command string, promptBytes []byte, cwd string, env []string, timeoutSec int) ([]byte, int, error) {
	return f(command, promptBytes, cwd, env, timeoutSec)
}

// Run validates the AI command, then runs the loop: for each iteration invokes
// the backend with the assembled prompt, captures stdout, and scans for the
// configured success and failure signals. On success: reports completion and
// returns ExitSuccess. On failure signal or process exit without signal (T3.8,
// O001/R009): increments consecutive-failure count; if count >= failure
// threshold, reports and returns ExitFailureThreshold. When max iterations is
// reached without success, returns ExitMaxIterations. On SIGINT or SIGTERM
// (T3.9, O004/R005): returns ExitInterrupt (130). Static precedence
// (T3.6, O001/R006): success is checked before failure; when both signals
// appear in the same output, the iteration is treated as success.
// Implements T3.4, T3.5, T3.6, T3.7, T3.8, T3.9, O001/R004, O001/R005, O001/R006, O001/R007,
// O001/R009, O004/R002, O004/R003, O004/R004, O004/R005.
func Run(opts RunOptions) (exitCode int, err error) {
	if opts.Invoker == nil {
		opts.Invoker = invokerAdapter(backend.Invoke)
	}
	report := opts.Reporter
	if report == nil {
		report = func(msg string) { fmt.Fprintln(os.Stdout, msg) }
	}

	// Dry-run: assemble prompt (with preamble if enabled), print to stdout, exit 0. No backend (T3.10, O004/R007).
	if opts.DryRun {
		preamble := buildPreamble(opts.Loop.Preamble, 1)
		assembled := AssemblePrompt(preamble, opts.PromptBytes)
		os.Stdout.Write(assembled)
		report("Dry-run: assembled prompt printed; no run was performed.")
		return ExitSuccess, nil
	}

	if err := ValidateAICommand(opts.Command); err != nil {
		return ExitErrorPreLoop, err
	}

	// Interrupt: use optional context (e.g. for tests) or os.Interrupt (O004/R005).
	var ctx context.Context
	var stop context.CancelFunc
	if opts.InterruptContext != nil {
		ctx = opts.InterruptContext
		stop = func() {}
	} else {
		ctx, stop = signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()
	}

	start := time.Now()
	consecutiveFailures := 0
	threshold := opts.Loop.FailureThreshold
	if threshold <= 0 {
		threshold = 1
	}
	for i := 1; i <= opts.Loop.MaxIterations; i++ {
		select {
		case <-ctx.Done():
			return ExitInterrupt, nil
		default:
		}
		preamble := buildPreamble(opts.Loop.Preamble, i)
		assembled := AssemblePrompt(preamble, opts.PromptBytes)
		stdout, _, invErr := opts.Invoker.Invoke(opts.Command, assembled, opts.Cwd, opts.Env, opts.Loop.TimeoutSeconds)
		if ctx.Err() != nil {
			return ExitInterrupt, nil
		}
		// Invocation error (e.g. timeout, crash, exec failure): treat as no-signal failure per T3.8/O001/R009.
		if invErr != nil {
			consecutiveFailures++
			if consecutiveFailures >= threshold {
				report(fmt.Sprintf("Stopped after %d consecutive iteration(s) without success or failure signal (invocation error: %v; threshold: %d).", consecutiveFailures, invErr, opts.Loop.FailureThreshold))
				return ExitFailureThreshold, nil
			}
			continue
		}
		if ContainsSuccessSignal(stdout, opts.Loop.SuccessSignal) {
			elapsed := time.Since(start)
			report(completionMessage(i, elapsed))
			return ExitSuccess, nil
		}
		// Failure: either failure signal present or process exited without signal (T3.8, O001/R009).
		hasFailureSignal := ContainsFailureSignal(stdout, opts.Loop.FailureSignal)
		consecutiveFailures++
		if consecutiveFailures >= threshold {
			if hasFailureSignal {
				report(fmt.Sprintf("Stopped after %d consecutive failure(s) (threshold: %d).", consecutiveFailures, opts.Loop.FailureThreshold))
			} else {
				report(fmt.Sprintf("Stopped after %d consecutive iteration(s) without success or failure signal (threshold: %d).", consecutiveFailures, opts.Loop.FailureThreshold))
			}
			return ExitFailureThreshold, nil
		}
	}
	report(fmt.Sprintf("Stopped after %d iteration(s) without success signal (max: %d).", opts.Loop.MaxIterations, opts.Loop.MaxIterations))
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
