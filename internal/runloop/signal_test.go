package runloop

import "testing"

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
