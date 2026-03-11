package runloop

import "bytes"

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
