package config

import (
	"fmt"
	"strings"
)

// ValidationError represents a single validation error.
type ValidationError struct {
	Field      string
	Value      string
	Message    string
	Provenance Provenance
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	if e.Provenance != "" {
		return fmt.Sprintf("validation error: %s (source: %s)", e.Message, e.Provenance)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

// Error implements the error interface.
func (e ValidationErrors) Error() string {
	var sb strings.Builder
	for i, err := range e {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}

// Validate validates the resolved configuration.
// Returns ValidationErrors if any validation errors are found.
func Validate(cfg ConfigWithProvenance) error {
	var errors ValidationErrors

	// Schema validation
	errors = append(errors, validateSchema(cfg)...)

	// Semantic validation
	errors = append(errors, validateSemantic(cfg)...)

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// validateSchema performs schema validation on resolved config values.
func validateSchema(cfg ConfigWithProvenance) ValidationErrors {
	var errors ValidationErrors

	// loop.default_max_iterations: minimum 1
	if cfg.Loop.DefaultMaxIterations.Value < 1 {
		errors = append(errors, ValidationError{
			Field:      "loop.default_max_iterations",
			Value:      fmt.Sprintf("%d", cfg.Loop.DefaultMaxIterations.Value),
			Message:    fmt.Sprintf("loop.default_max_iterations must be >= 1, got %d", cfg.Loop.DefaultMaxIterations.Value),
			Provenance: cfg.Loop.DefaultMaxIterations.Provenance,
		})
	}

	// loop.iteration_mode: must be "max-iterations" or "unlimited"
	if cfg.Loop.IterationMode.Value != "max-iterations" && cfg.Loop.IterationMode.Value != "unlimited" {
		errors = append(errors, ValidationError{
			Field:      "loop.iteration_mode",
			Value:      cfg.Loop.IterationMode.Value,
			Message:    fmt.Sprintf("loop.iteration_mode must be \"max-iterations\" or \"unlimited\", got %q", cfg.Loop.IterationMode.Value),
			Provenance: cfg.Loop.IterationMode.Provenance,
		})
	}

	// loop.failure_threshold: minimum 1
	if cfg.Loop.FailureThreshold.Value < 1 {
		errors = append(errors, ValidationError{
			Field:      "loop.failure_threshold",
			Value:      fmt.Sprintf("%d", cfg.Loop.FailureThreshold.Value),
			Message:    fmt.Sprintf("loop.failure_threshold must be >= 1, got %d", cfg.Loop.FailureThreshold.Value),
			Provenance: cfg.Loop.FailureThreshold.Provenance,
		})
	}

	// loop.iteration_timeout: minimum 0 (0 means no timeout)
	if cfg.Loop.IterationTimeout.Value < 0 {
		errors = append(errors, ValidationError{
			Field:      "loop.iteration_timeout",
			Value:      fmt.Sprintf("%d", cfg.Loop.IterationTimeout.Value),
			Message:    fmt.Sprintf("loop.iteration_timeout must be >= 0, got %d", cfg.Loop.IterationTimeout.Value),
			Provenance: cfg.Loop.IterationTimeout.Provenance,
		})
	}

	// loop.max_output_buffer: minimum 1
	if cfg.Loop.MaxOutputBuffer.Value < 1 {
		errors = append(errors, ValidationError{
			Field:      "loop.max_output_buffer",
			Value:      fmt.Sprintf("%d", cfg.Loop.MaxOutputBuffer.Value),
			Message:    fmt.Sprintf("loop.max_output_buffer must be >= 1, got %d", cfg.Loop.MaxOutputBuffer.Value),
			Provenance: cfg.Loop.MaxOutputBuffer.Provenance,
		})
	}

	// loop.signals.success: min length 1
	if cfg.Loop.SignalSuccess.Value == "" {
		errors = append(errors, ValidationError{
			Field:      "loop.signals.success",
			Value:      `""`,
			Message:    "loop.signals.success must not be empty",
			Provenance: cfg.Loop.SignalSuccess.Provenance,
		})
	}

	// loop.signals.failure: min length 1
	if cfg.Loop.SignalFailure.Value == "" {
		errors = append(errors, ValidationError{
			Field:      "loop.signals.failure",
			Value:      `""`,
			Message:    "loop.signals.failure must not be empty",
			Provenance: cfg.Loop.SignalFailure.Provenance,
		})
	}

	// Validate prompts
	for alias, prompt := range cfg.Prompts {
		if prompt.Path == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("prompts.%s.path", alias),
				Value:   `""`,
				Message: fmt.Sprintf("prompts.%s.path must not be empty", alias),
			})
		}
	}

	return errors
}

// validateSemantic performs semantic validation (cross-referencing).
func validateSemantic(cfg ConfigWithProvenance) ValidationErrors {
	var errors ValidationErrors

	// Validate AI command resolution
	// Note: We don't fail validation if neither ai_cmd nor ai_cmd_alias is set,
	// because ResolveAICommand will produce a clear error at runtime.
	// We only validate that if an alias is specified, it exists.
	if cfg.Loop.AICmdAlias.Value != "" {
		if _, exists := cfg.AICmdAliases[cfg.Loop.AICmdAlias.Value]; !exists {
			errors = append(errors, ValidationError{
				Field:      "loop.ai_cmd_alias",
				Value:      cfg.Loop.AICmdAlias.Value,
				Message:    fmt.Sprintf("invalid ai_cmd_alias %q: no matching alias defined", cfg.Loop.AICmdAlias.Value),
				Provenance: cfg.Loop.AICmdAlias.Provenance,
			})
		}
	}

	return errors
}
