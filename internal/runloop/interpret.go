package runloop

import "bytes"

// InterpretationPrompt is the built-in prompt used when signal_precedence is
// ai_interpreted and both success and failure signals appear in the same output.
// Product-owned, not user-editable (O001/R008). The AI is asked to interpret
// the iteration output and respond with exactly one of the two signal markers.
const interpretationPromptPrefix = `You are a classifier. Below is the full stdout from a single iteration of an automated task. The iteration output may contain both a success marker and a failure marker; your job is to decide whether the task for that iteration ultimately succeeded or failed.

Reply with exactly one line containing only one of these two markers (copy it literally):
- SUCCESS_MARKER
- FAILURE_MARKER

Do not add any other text before or after the marker. If the overall outcome is success, output SUCCESS_MARKER. If the overall outcome is failure, output FAILURE_MARKER.

--- iteration stdout ---
`

const interpretationPromptSuffix = `
--- end iteration stdout ---
`

// BuildInterpretationPrompt returns the prompt bytes for the one-off AI
// invocation that interprets iteration output when both signals are present.
// successMarker and failureMarker are the configured signal strings; they
// are substituted into the prompt so the AI knows which literal to output.
// Implements O001/R008: at most one interpretation invocation per ambiguous iteration.
func BuildInterpretationPrompt(iterationStdout []byte, successMarker, failureMarker string) []byte {
	if successMarker == "" {
		successMarker = "<promise>SUCCESS</promise>"
	}
	if failureMarker == "" {
		failureMarker = "<promise>FAILURE</promise>"
	}
	prefix := bytes.ReplaceAll([]byte(interpretationPromptPrefix), []byte("SUCCESS_MARKER"), []byte(successMarker))
	prefix = bytes.ReplaceAll(prefix, []byte("FAILURE_MARKER"), []byte(failureMarker))
	out := make([]byte, 0, len(prefix)+len(iterationStdout)+len(interpretationPromptSuffix))
	out = append(out, prefix...)
	out = append(out, iterationStdout...)
	out = append(out, interpretationPromptSuffix...)
	return out
}

// InterpretedOutcome is the result of parsing the AI's interpretation response.
type InterpretedOutcome int

const (
	InterpretedUnclear InterpretedOutcome = iota
	InterpretedSuccess
	InterpretedFailure
)

// ParseInterpretationResponse parses the stdout from the interpretation
// invocation. Returns InterpretedSuccess or InterpretedFailure if the
// response clearly contains the corresponding marker (and not the other);
// returns InterpretedUnclear if both, neither, or ambiguous. When unclear,
// the run-loop applies the defined fallback (e.g. treat as failure).
func ParseInterpretationResponse(stdout []byte, successMarker, failureMarker string) InterpretedOutcome {
	if successMarker == "" {
		successMarker = "<promise>SUCCESS</promise>"
	}
	if failureMarker == "" {
		failureMarker = "<promise>FAILURE</promise>"
	}
	hasSuccess := ContainsSuccessSignal(stdout, successMarker)
	hasFailure := ContainsFailureSignal(stdout, failureMarker)
	if hasSuccess && !hasFailure {
		return InterpretedSuccess
	}
	if hasFailure && !hasSuccess {
		return InterpretedFailure
	}
	return InterpretedUnclear
}
