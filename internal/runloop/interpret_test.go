package runloop

import (
	"bytes"
	"testing"
)

func TestBuildInterpretationPrompt(t *testing.T) {
	stdout := []byte("iteration output here")
	got := BuildInterpretationPrompt(stdout, "DONE", "FAIL")
	if !bytes.Contains(got, stdout) {
		t.Error("prompt must contain iteration stdout")
	}
	if !bytes.Contains(got, []byte("DONE")) || !bytes.Contains(got, []byte("FAIL")) {
		t.Errorf("prompt must contain success and failure markers: %s", got)
	}
	if !bytes.Contains(got, []byte("--- iteration stdout ---")) {
		t.Error("prompt must contain delimiter")
	}
}

func TestParseInterpretationResponse_clearSuccess(t *testing.T) {
	outcome := ParseInterpretationResponse([]byte("some text DONE more"), "DONE", "FAIL")
	if outcome != InterpretedSuccess {
		t.Errorf("ParseInterpretationResponse(DONE only) = %v, want InterpretedSuccess", outcome)
	}
}

func TestParseInterpretationResponse_clearFailure(t *testing.T) {
	outcome := ParseInterpretationResponse([]byte("output FAIL here"), "DONE", "FAIL")
	if outcome != InterpretedFailure {
		t.Errorf("ParseInterpretationResponse(FAIL only) = %v, want InterpretedFailure", outcome)
	}
}

func TestParseInterpretationResponse_bothUnclear(t *testing.T) {
	outcome := ParseInterpretationResponse([]byte("DONE and FAIL"), "DONE", "FAIL")
	if outcome != InterpretedUnclear {
		t.Errorf("ParseInterpretationResponse(both) = %v, want InterpretedUnclear", outcome)
	}
}

func TestParseInterpretationResponse_neitherUnclear(t *testing.T) {
	outcome := ParseInterpretationResponse([]byte("no markers"), "DONE", "FAIL")
	if outcome != InterpretedUnclear {
		t.Errorf("ParseInterpretationResponse(neither) = %v, want InterpretedUnclear", outcome)
	}
}

func TestParseInterpretationResponse_defaultMarkers(t *testing.T) {
	success := ParseInterpretationResponse([]byte("<promise>SUCCESS</promise>"), "", "")
	if success != InterpretedSuccess {
		t.Errorf("default success marker: got %v, want InterpretedSuccess", success)
	}
	fail := ParseInterpretationResponse([]byte("<promise>FAILURE</promise>"), "", "")
	if fail != InterpretedFailure {
		t.Errorf("default failure marker: got %v, want InterpretedFailure", fail)
	}
}
