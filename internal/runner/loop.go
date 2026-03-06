package runner

import (
	"bytes"
	"io"

	"github.com/maxdunn123/ralph/internal/config"
)

// IterationResult captures the outcome of a single iteration.
type IterationResult struct {
	Output []byte
	Error  error
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
		Unlimited:      false, // TODO: support unlimited mode when implemented
		ContextStrings: contextStrings,
	})

	// Assemble prompt
	assembled := AssemblePrompt(preamble, promptContent)

	// Create bounded buffer for output capture
	buffer := NewBoundedBuffer(cfg.Loop.MaxOutputBuffer.Value)

	// Spawn AI process with assembled prompt as stdin
	stdin := bytes.NewReader(assembled)
	err := SpawnAI(aiCmd, stdin, buffer, buffer)

	return IterationResult{
		Output: buffer.Bytes(),
		Error:  err,
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

	for i := 1; i <= maxIterations; i++ {
		result := RunIteration(i, aiCmd, promptContent, cfg, contextStrings)

		// Scan for signals after process exit
		outcome := ScanForSignals(result.Output, cfg.Loop.SignalSuccess.Value, cfg.Loop.SignalFailure.Value)

		// Handle iteration outcome
		switch outcome {
		case OutcomeSuccess:
			// Success signal found - exit loop with success
			return nil
		case OutcomeFailure:
			// Failure signal found - for now, continue loop
			// TODO: Consecutive failure tracking (ralph-wnp)
			continue
		case OutcomeNoSignal:
			// No signal found - continue to next iteration
			// TODO: Reset failure counter (ralph-wnp)
			continue
		}
	}

	// Max iterations exhausted
	// TODO: Exit code 2 (ralph-qim)
	return nil
}
