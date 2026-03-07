package runner

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maxdunn/ralph/internal/config"
	"github.com/maxdunn/ralph/internal/logger"
)

// Exit code errors
var (
	ExitCodeExhausted        = errors.New("max iterations exhausted")
	ExitCodeFailureThreshold = errors.New("failure threshold reached")
	ExitCodeInterrupted      = errors.New("interrupted")
)

// IterationResult captures the outcome of a single iteration.
type IterationResult struct {
	Output      []byte
	ExitCode    int
	Error       error
	Interrupted bool
	Duration    time.Duration
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
	start := time.Now()

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

	// Setup output writers: buffer always captures, optionally mirror to stdout (run log)
	var stdout, stderr io.Writer
	if verbose {
		stdout = io.MultiWriter(buffer, os.Stdout)
		stderr = io.MultiWriter(buffer, os.Stdout)
	} else {
		stdout = buffer
		stderr = buffer
	}

	// Apply per-iteration timeout if configured (O1/R3)
	iterCtx := ctx
	var cancel context.CancelFunc
	if cfg.Loop.IterationTimeout.Value > 0 {
		iterCtx, cancel = context.WithTimeout(ctx, time.Duration(cfg.Loop.IterationTimeout.Value)*time.Second)
		defer cancel()
	}

	// Spawn AI process with assembled prompt as stdin
	stdin := bytes.NewReader(assembled)
	exitCode, err := SpawnAIWithContext(iterCtx, aiCmd, stdin, stdout, stderr)

	duration := time.Since(start)
	interrupted := exitCode == 130

	return IterationResult{
		Output:      buffer.Bytes(),
		ExitCode:    exitCode,
		Error:       err,
		Interrupted: interrupted,
		Duration:    duration,
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

	// Initialize statistics tracker (O4/R2)
	stats := NewIterationStats()

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
			// Interrupted - do not report statistics (O4/R2)
			return ExitCodeInterrupted
		default:
		}

		// Check iteration limit before executing iteration (O1/R4)
		if iterationMode == "max-iterations" && i > maxIterations {
			// Max iterations exhausted without success signal
			// Report statistics before exit (O4/R2)
			logger.Info(stats.Report())
			return ExitCodeExhausted
		}

		// Emit iteration progress message to stdout (O4/R6)
		if iterationMode == "unlimited" {
			logger.Info("Iteration %d (unlimited)", i)
		} else {
			logger.Info("Iteration %d/%d", i, maxIterations)
		}

		result := RunIteration(ctx, i, aiCmd, promptContent, cfg, contextStrings, verbose)

		// If interrupted, discard output and exit 130 (O1/R7)
		// Do not report statistics on interruption (O4/R2)
		if result.Interrupted {
			return ExitCodeInterrupted
		}

		// Record iteration duration (O4/R2)
		stats.Add(result.Duration)

		// Log crashes (non-zero exit) at warn level (O1/R1)
		if result.ExitCode != 0 {
			logger.Warn("AI process exited with code %d (crash)", result.ExitCode)
		}

		// Scan for signals after process exit
		outcome := ScanForSignals(result.Output, cfg.Loop.SignalSuccess.Value, cfg.Loop.SignalFailure.Value)

		// Handle iteration outcome (O1/R5)
		switch outcome {
		case OutcomeSuccess:
			// Success signal found - report statistics and exit with success (O4/R2)
			logger.Info(stats.Report())
			return nil
		case OutcomeFailure:
			// Failure signal found - increment consecutive failure counter
			consecutiveFailures++
			if consecutiveFailures >= failureThreshold {
				// Failure threshold reached - report statistics and abort loop (O4/R2)
				logger.Info("Failure threshold (%d) reached after %d consecutive failure(s); stopping.", failureThreshold, consecutiveFailures)
				logger.Info(stats.Report())
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
