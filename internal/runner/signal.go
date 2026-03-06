package runner

import "bytes"

// IterationOutcome represents the result of signal scanning.
type IterationOutcome int

const (
	OutcomeNoSignal IterationOutcome = iota
	OutcomeSuccess
	OutcomeFailure
)

// ScanForSignals scans output buffer for success/failure signals.
// Returns the iteration outcome based on signal precedence: failure wins over success.
func ScanForSignals(output []byte, successSignal, failureSignal string) IterationOutcome {
	hasFailure := bytes.Contains(output, []byte(failureSignal))
	hasSuccess := bytes.Contains(output, []byte(successSignal))

	if hasFailure {
		return OutcomeFailure
	}
	if hasSuccess {
		return OutcomeSuccess
	}
	return OutcomeNoSignal
}
