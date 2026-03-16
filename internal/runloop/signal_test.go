package runloop

import (
	"bytes"
	"testing"
)

func TestLastNonEmptyLine(t *testing.T) {
	tests := []struct {
		name   string
		stdout []byte
		want   []byte
	}{
		{"empty", []byte(""), nil},
		{"single line no newline", []byte("only"), []byte("only")},
		{"single line with newline", []byte("only\n"), []byte("only")},
		{"multi-line last", []byte("a\nb\nc"), []byte("c")},
		{"trailing newlines", []byte("last\n\n\n"), []byte("last")},
		{"all empty lines", []byte("\n\n\n"), nil},
		{"trim", []byte("  \n  last  \n"), []byte("last")},
		{"only whitespace", []byte("   \n  \t  \n"), nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LastNonEmptyLine(tt.stdout)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("LastNonEmptyLine(%q) = %q, want %q", tt.stdout, got, tt.want)
			}
		})
	}
}

func TestContainsSuccessSignal(t *testing.T) {
	tests := []struct {
		name   string
		stdout []byte
		signal string
		want   bool
	}{
		{"match substring", []byte("foo <promise>SUCCESS</promise> bar"), "<promise>SUCCESS</promise>", true},
		{"no match", []byte("Still working..."), "<promise>SUCCESS</promise>", false},
		{"empty output", []byte(""), "<promise>SUCCESS</promise>", false},
		{"empty signal never matches", []byte("anything"), "", false},
		{"signal equals output", []byte("DONE"), "DONE", true},
		{"match at start", []byte("<promise>SUCCESS</promise> tail"), "<promise>SUCCESS</promise>", true},
		{"match at end", []byte("head <promise>SUCCESS</promise>"), "<promise>SUCCESS</promise>", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsSuccessSignal(tt.stdout, tt.signal)
			if got != tt.want {
				t.Errorf("ContainsSuccessSignal(%q, %q) = %v, want %v", tt.stdout, tt.signal, got, tt.want)
			}
		})
	}
}

// TestContainsSuccessSignal_LastLineOnly verifies that success is detected only when
// the signal appears on the last non-empty line (run-loop passes LastNonEmptyLine(stdout)).
func TestContainsSuccessSignal_LastLineOnly(t *testing.T) {
	// Full stdout with signal only on an earlier line -> last line has no signal -> no success.
	full := []byte("Status: DONE\nStill working...")
	lastLine := LastNonEmptyLine(full)
	if bytes.Equal(lastLine, []byte("Still working...")) {
		// good: last line is "Still working..."
	}
	got := ContainsSuccessSignal(lastLine, "DONE")
	if got {
		t.Errorf("ContainsSuccessSignal(lastLine, \"DONE\") = true with last line %q (signal only on earlier line); want false", lastLine)
	}
	// Signal on last non-empty line -> success.
	full2 := []byte("Still working...\nStatus: DONE")
	lastLine2 := LastNonEmptyLine(full2)
	got2 := ContainsSuccessSignal(lastLine2, "DONE")
	if !got2 {
		t.Errorf("ContainsSuccessSignal(lastLine, \"DONE\") = false with last line %q; want true", lastLine2)
	}
}

// TestContainsFailureSignal_LastLineOnly verifies that failure is detected only when
// the signal appears on the last non-empty line.
func TestContainsFailureSignal_LastLineOnly(t *testing.T) {
	// Failure only on earlier line -> last line has no failure signal -> no failure.
	full := []byte("FAIL\nStill working...")
	lastLine := LastNonEmptyLine(full)
	got := ContainsFailureSignal(lastLine, "FAIL")
	if got {
		t.Errorf("ContainsFailureSignal(lastLine, \"FAIL\") = true with last line %q; want false", lastLine)
	}
	// Failure on last non-empty line -> failure.
	full2 := []byte("Still working...\nFAIL")
	lastLine2 := LastNonEmptyLine(full2)
	got2 := ContainsFailureSignal(lastLine2, "FAIL")
	if !got2 {
		t.Errorf("ContainsFailureSignal(lastLine, \"FAIL\") = false with last line %q; want true", lastLine2)
	}
}

func TestContainsFailureSignal(t *testing.T) {
	tests := []struct {
		name   string
		stdout []byte
		signal string
		want   bool
	}{
		{"match substring", []byte("err <promise>FAILURE</promise> more"), "<promise>FAILURE</promise>", true},
		{"no match", []byte("Still working..."), "<promise>FAILURE</promise>", false},
		{"empty signal never matches", []byte("anything"), "", false},
		{"signal equals output", []byte("FAIL"), "FAIL", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsFailureSignal(tt.stdout, tt.signal)
			if got != tt.want {
				t.Errorf("ContainsFailureSignal(%q, %q) = %v, want %v", tt.stdout, tt.signal, got, tt.want)
			}
		})
	}
}
