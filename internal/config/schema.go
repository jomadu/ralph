// Package config: validate canonical YAML structure per docs/engineering/components/config.md.
// Invalid or out-of-range values produce a clear error (O002, config component).

package config

import (
	"errors"
	"fmt"
)

// Schema validation errors (callers can use errors.Is).
var (
	ErrInvalidMaxIterations     = errors.New("max_iterations must be >= 0")
	ErrInvalidFailureThreshold  = errors.New("failure_threshold must be >= 0")
	ErrInvalidTimeoutSeconds    = errors.New("timeout_seconds must be >= 0")
	ErrInvalidSignalPrecedence  = errors.New("signal_precedence must be empty, \"static\", or \"ai_interpreted\"")
	ErrInvalidLogLevel          = errors.New("log_level must be empty, \"debug\", \"info\", \"warn\", or \"error\"")
	ErrPromptNeedsPathOrContent = errors.New("prompt must have path or content")
	ErrAliasCommandEmpty        = errors.New("alias command must be non-empty")
)

// Valid signal_precedence values (empty = use default).
var validSignalPrecedence = map[string]bool{
	"":               true,
	"static":         true,
	"ai_interpreted": true,
}

// Valid log_level values (empty = use default).
var validLogLevel = map[string]bool{
	"":      true,
	"debug": true,
	"info":  true,
	"warn":  true,
	"error": true,
}

// ValidateFileLayer validates a parsed file layer against the canonical config schema.
// Reject or error on invalid/out-of-range values (config component spec).
// Call after ParseLayer when the layer is non-nil.
func ValidateFileLayer(layer *FileLayer) error {
	if layer == nil {
		return nil
	}
	if err := validateLoop(layer.Loop); err != nil {
		return err
	}
	if err := validatePrompts(layer.Prompts); err != nil {
		return err
	}
	if err := validateAliases(layer.Aliases); err != nil {
		return err
	}
	return nil
}

func validateLoop(loop *LoopSection) error {
	if loop == nil {
		return nil
	}
	if loop.MaxIterations != nil && *loop.MaxIterations < 0 {
		return fmt.Errorf("%w (got %d)", ErrInvalidMaxIterations, *loop.MaxIterations)
	}
	if loop.FailureThreshold != nil && *loop.FailureThreshold < 0 {
		return fmt.Errorf("%w (got %d)", ErrInvalidFailureThreshold, *loop.FailureThreshold)
	}
	if loop.TimeoutSeconds != nil && *loop.TimeoutSeconds < 0 {
		return fmt.Errorf("%w (got %d)", ErrInvalidTimeoutSeconds, *loop.TimeoutSeconds)
	}
	if loop.SignalPrecedence != "" && !validSignalPrecedence[loop.SignalPrecedence] {
		return fmt.Errorf("%w (got %q)", ErrInvalidSignalPrecedence, loop.SignalPrecedence)
	}
	if loop.LogLevel != "" && !validLogLevel[loop.LogLevel] {
		return fmt.Errorf("%w (got %q)", ErrInvalidLogLevel, loop.LogLevel)
	}
	// preamble: string or bool — no extra validation
	// streaming: *bool — no extra validation
	return nil
}

func validatePrompts(prompts map[string]Prompt) error {
	for name, p := range prompts {
		hasPath := p.Path != ""
		hasContent := p.Content != ""
		if !hasPath && !hasContent {
			return fmt.Errorf("prompt %q: %w", name, ErrPromptNeedsPathOrContent)
		}
		if err := validateLoop(p.Loop); err != nil {
			return fmt.Errorf("prompt %q loop: %w", name, err)
		}
	}
	return nil
}

func validateAliases(aliases map[string]Alias) error {
	for name, a := range aliases {
		if a.Command == "" {
			return fmt.Errorf("alias %q: %w", name, ErrAliasCommandEmpty)
		}
	}
	return nil
}
