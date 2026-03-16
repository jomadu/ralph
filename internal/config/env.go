// Package config: environment variable overlay per docs/engineering/components/config.md.
// RALPH_CONFIG_HOME is handled in paths.go. This file implements RALPH_LOOP_* overlay.

package config

import (
	"fmt"
	"strconv"
	"strings"
)

// EnvOverlay holds values from RALPH_LOOP_* environment variables.
// Only set (non-nil) fields should override file-based config; unset vars leave fields nil.
type EnvOverlay struct {
	AICmd            *string
	AICmdAlias       *string
	IterationMode    *string
	MaxIterations    *int
	FailureThreshold *int
	IterationTimeout *int
	LogLevel         *string
	Streaming        *bool
	Preamble         *bool
	MaxOutputBuffer  *int
}

// ParseEnvOverlay reads RALPH_LOOP_* from getenv and returns an overlay.
// Unset variables do not appear in the overlay. Invalid values produce a clear error.
// getenv is typically os.Getenv.
func ParseEnvOverlay(getenv func(string) string) (*EnvOverlay, error) {
	o := &EnvOverlay{}

	if v := getenv("RALPH_LOOP_AI_CMD"); v != "" {
		s := v
		o.AICmd = &s
	}
	if v := getenv("RALPH_LOOP_AI_CMD_ALIAS"); v != "" {
		s := v
		o.AICmdAlias = &s
	}
	if v := getenv("RALPH_LOOP_ITERATION_MODE"); v != "" {
		s := v
		o.IterationMode = &s
	}

	if v := getenv("RALPH_LOOP_DEFAULT_MAX_ITERATIONS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("RALPH_LOOP_DEFAULT_MAX_ITERATIONS: invalid integer %q: %w", v, err)
		}
		if n < 0 {
			return nil, fmt.Errorf("RALPH_LOOP_DEFAULT_MAX_ITERATIONS: must be >= 0, got %d", n)
		}
		o.MaxIterations = &n
	}
	if v := getenv("RALPH_LOOP_FAILURE_THRESHOLD"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("RALPH_LOOP_FAILURE_THRESHOLD: invalid integer %q: %w", v, err)
		}
		if n < 0 {
			return nil, fmt.Errorf("RALPH_LOOP_FAILURE_THRESHOLD: must be >= 0, got %d", n)
		}
		o.FailureThreshold = &n
	}
	if v := getenv("RALPH_LOOP_ITERATION_TIMEOUT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("RALPH_LOOP_ITERATION_TIMEOUT: invalid integer %q: %w", v, err)
		}
		if n < 0 {
			return nil, fmt.Errorf("RALPH_LOOP_ITERATION_TIMEOUT: must be >= 0, got %d", n)
		}
		o.IterationTimeout = &n
	}

	if v := getenv("RALPH_LOOP_LOG_LEVEL"); v != "" {
		s := v
		o.LogLevel = &s
	}

	if v := getenv("RALPH_LOOP_STREAMING"); v != "" {
		b, err := parseBool(v)
		if err != nil {
			return nil, fmt.Errorf("RALPH_LOOP_STREAMING: %w", err)
		}
		o.Streaming = &b
	}
	if v := getenv("RALPH_LOOP_PREAMBLE"); v != "" {
		b, err := parseBool(v)
		if err != nil {
			return nil, fmt.Errorf("RALPH_LOOP_PREAMBLE: %w", err)
		}
		o.Preamble = &b
	}

	if v := getenv("RALPH_LOOP_MAX_OUTPUT_BUFFER"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("RALPH_LOOP_MAX_OUTPUT_BUFFER: invalid integer %q: %w", v, err)
		}
		if n < 0 {
			return nil, fmt.Errorf("RALPH_LOOP_MAX_OUTPUT_BUFFER: must be >= 0, got %d", n)
		}
		o.MaxOutputBuffer = &n
	}

	return o, nil
}

// parseBool parses env-style boolean: true/1/yes/on → true, false/0/no/off → false.
// Empty string is invalid (caller should skip unset). Other values return error.
func parseBool(s string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean %q (use true/1/yes/on or false/0/no/off)", s)
	}
}
