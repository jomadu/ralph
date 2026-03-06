package runner

import (
	"fmt"
	"strings"
)

// PreambleConfig holds settings for preamble generation.
type PreambleConfig struct {
	Enabled        bool
	Iteration      int
	MaxIterations  int
	Unlimited      bool
	ContextStrings []string
}

// GeneratePreamble creates the preamble text per O1/R8.
func GeneratePreamble(cfg PreambleConfig) string {
	if !cfg.Enabled {
		return ""
	}

	var sb strings.Builder

	// Iteration line
	if cfg.Unlimited {
		fmt.Fprintf(&sb, "[RALPH] Iteration %d of unlimited", cfg.Iteration)
	} else {
		fmt.Fprintf(&sb, "[RALPH] Iteration %d of %d", cfg.Iteration, cfg.MaxIterations)
	}

	// Context section if provided
	if len(cfg.ContextStrings) > 0 {
		sb.WriteString("\n\nCONTEXT:\n")
		sb.WriteString(strings.Join(cfg.ContextStrings, "\n\n"))
	}

	return sb.String()
}

// AssemblePrompt combines preamble and prompt content per O1/R8.
func AssemblePrompt(preamble string, promptContent []byte) []byte {
	if preamble == "" {
		return promptContent
	}
	return []byte(preamble + "\n\n" + string(promptContent))
}
