package runner

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/maxdunn/ralph/internal/config"
)

// Exit code errors
var (
	ExitCodeExhausted         = errors.New("max iterations exhausted")
	ExitCodeFailureThreshold  = errors.New("failure threshold reached")
	ExitCodeInterrupted       = errors.New("interrupted")
)

// IterationResult captures the outcome of a single iteration.
type IterationResult struct {
	Output      []byte
	ExitCode    int
	Error       error
	Interrupted bool
}

// RunIteration executes a single iteration: assemble prompt, spawn process, capture output.
func RunIteration(
	ctx context.Context,
	iteration int,
	aiCmd []string,
	promptContent []byte,
	cfg *config.ConfigWithProvenance,
	contextStrings []string,
	verbose bool,
) IterationResult {
	// Generate preamble
	preamble := GeneratePreamble(PreambleConfig{
		Enabled:        cfg.Loop.Preamble.Value,
		Iteration:      iteration,
		MaxIterations:  cfg.Loop.DefaultMaxIterations.Value,
		Unlimited:      cfg.Loop.IterationMode.Value == "unlimited",
		ContextStrings: contextStrings,
	})

	// Assemble prompt
	assembled := AssemblePrompt(preamble, promptContent)

	// Create bounded buffer for output capture
	buffer := NewBoundedBuffer(cfg.Loop.MaxOutputBuffer.Value)

	// Setup output writers: buffer always captures, optionally mirror to terminal
	var stdout, stderr io.Writer
	if verbose {
		stdout = io.MultiWriter(buffer, os.Stderr)
		stderr = io.MultiWriter(buffer, os.Stderr)
	} else {
		stdout = buffer
		stderr = buffer
	}

	// Spawn AI process with assembled prompt as stdin
	stdin := bytes.NewReader(assembled)
	exitCode, err := SpawnAIWithContext(ctx, aiCmd, stdin, stdout, stderr)

	interrupted := exitCode == 130

	return IterationResult{
		Output:      buffer.Bytes(),
		ExitCode:    exitCode,
		Error:       err,
		Interrupted: interrupted,
	}
}

// Loop executes the main iteration loop.
func Loop(
	aiCmd []string,
	promptContent []byte,
	cfg *config.ConfigWithProvenance,
	contextStrings []string,
) error {
	maxIterations := cfg.Loop.DefaultMaxIterations.Value
	iterationMode := cfg.Loop.IterationMode.Value
	failureThreshold := cfg.Loop.FailureThreshold.Value
	consecutiveFailures := 0
	verbose := cfg.Loop.ShowAIOutput.Value

	// Setup signal handling for SIGINT/SIGTERM (O1/R7)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	go func() {
		<-sigChan
		cancel()
	}()

	for i := 1; ; i++ {
		// Check for interruption between iterations (O1/R7)
		select {
		case <-ctx.Done():
			return ExitCodeInterrupted
		default:
		}

		// Check iteration limit before executing iteration (O1/R4)
		if iterationMode == "max-iterations" && i > maxIterations {
			// Max iterations exhausted without success signal
			return ExitCodeExhausted
		}

		result := RunIteration(ctx, i, aiCmd, promptContent, cfg, contextStrings, verbose)

		// If interrupted, discard output and exit 130 (O1/R7)
		if result.Interrupted {
			return ExitCodeInterrupted
		}

		// Log crashes (non-zero exit) at warn level (O1/R1)
		if result.ExitCode != 0 {
			log.Printf("WARN: AI process exited with code %d (crash)", result.ExitCode)
		}

		// Scan for signals after process exit
		outcome := ScanForSignals(result.Output, cfg.Loop.SignalSuccess.Value, cfg.Loop.SignalFailure.Value)

		// Handle iteration outcome (O1/R5)
		switch outcome {
		case OutcomeSuccess:
			// Success signal found - exit loop with success
			return nil
		case OutcomeFailure:
			// Failure signal found - increment consecutive failure counter
			consecutiveFailures++
			if consecutiveFailures >= failureThreshold {
				// Failure threshold reached - abort loop
				return ExitCodeFailureThreshold
			}
			// Continue to next iteration
			continue
		case OutcomeNoSignal:
			// No signal found - reset consecutive failure counter
			consecutiveFailures = 0
			// Continue to next iteration
			continue
		}
	}
}
