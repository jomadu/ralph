package runloop

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ErrInvalidCommand is returned when the resolved command is empty or whitespace.
var ErrInvalidCommand = errors.New("invalid AI command: empty or whitespace")

// ValidateAICommand checks that the resolved AI command (executable name or path,
// optionally with args) is available and executable. It does not resolve aliases;
// the caller must pass the final command string (e.g. after alias expansion).
// Implements O001/R001 and O004/R001: fail before loop with clear error if
// command is missing, invalid, or not executable.
func ValidateAICommand(resolvedCommand string) error {
	cmd := strings.TrimSpace(resolvedCommand)
	if cmd == "" {
		return ErrInvalidCommand
	}
	exe := firstWord(cmd)
	path, err := exec.LookPath(exe)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return fmt.Errorf("AI command not found: %q (not on PATH or invalid)", exe)
		}
		return fmt.Errorf("AI command lookup failed: %w", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("AI command path not accessible: %w", err)
	}
	if info.Mode()&0111 == 0 {
		return fmt.Errorf("AI command not executable: %s", path)
	}
	return nil
}

func firstWord(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexAny(s, " \t"); i >= 0 {
		return s[:i]
	}
	return s
}
