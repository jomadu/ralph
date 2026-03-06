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

		// TODO: Signal scanning (ralph-3u5)
		// TODO: Failure tracking (ralph-3u5)
		// TODO: Exit code handling (ralph-3u5)

		if result.Error != nil {
			// Process crash - for now, just return the error
			// TODO: Proper crash handling per O1/R1
			return result.Error
		}

		// For now, just run one iteration and return
		// Full loop logic will be implemented when signal scanning is added
		break
	}

	return nil
}
