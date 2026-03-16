package runloop

import (
	"bytes"
	"strings"
)

// LastNonEmptyLine returns the last non-empty line of stdout. Stdout is split
// on newline (\n), each line is trimmed (leading and trailing whitespace), and
// the last line that is non-empty after trim is returned. If there is no
// non-empty line, returns nil. Callers use this for signal detection so only
// the last line is scanned (run-loop spec).
func LastNonEmptyLine(stdout []byte) []byte {
	lines := strings.Split(string(stdout), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		s := strings.TrimSpace(lines[i])
		if s != "" {
			return []byte(s)
		}
	}
	return nil
}

// ContainsSuccessSignal reports whether the configured success signal appears
// in the captured AI output. Uses substring match; empty signal never matches.
// Implements O001/R004: scan captured output for configured success signal.
func ContainsSuccessSignal(stdout []byte, signal string) bool {
	if signal == "" {
		return false
	}
	return bytes.Contains(stdout, []byte(signal))
}

// ContainsFailureSignal reports whether the configured failure signal appears
// in the captured AI output. Uses substring match; empty signal never matches.
// Implements O001/R005: detect failure signal for consecutive-failure count.
func ContainsFailureSignal(stdout []byte, signal string) bool {
	if signal == "" {
		return false
	}
	return bytes.Contains(stdout, []byte(signal))
}
