package runner

import (
	"bytes"
	"errors"
	"io"
	"log"

	"github.com/maxdunn123/ralph/internal/config"
)

// Exit code errors
var (
	ExitCodeExhausted         = errors.New("max iterations exhausted")
	ExitCodeFailureThreshold  = errors.New("failure threshold reached")
)

// IterationResult captures the outcome of a single iteration.
type IterationResult struct {
	Output   []byte
	ExitCode int
	Error    error
}

// RunIteration executes a single iteration: assemble prompt, spawn process, capture output.
func RunIteration(
	iteration int,
	aiCmd []string,
	promptContent []byte,
	cfg *config.ConfigWithProvenance,
	contextStrings []string,
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

	// Spawn AI process with assembled prompt as stdin
	stdin := bytes.NewReader(assembled)
	exitCode, err := SpawnAI(aiCmd, stdin, buffer, buffer)

	return IterationResult{
		Output:   buffer.Bytes(),
		ExitCode: exitCode,
		Error:    err,
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

	for i := 1; ; i++ {
		// Check iteration limit before executing iteration (O1/R4)
		if iterationMode == "max-iterations" && i > maxIterations {
			// Max iterations exhausted without success signal
			return ExitCodeExhausted
		}

		result := RunIteration(i, aiCmd, promptContent, cfg, contextStrings)

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
