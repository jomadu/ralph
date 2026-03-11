package runloop

import (
	"errors"
	"testing"
)

func TestValidateAICommand_EmptyOrWhitespace(t *testing.T) {
	for _, input := range []string{"", "   ", "\t", "\n"} {
		err := ValidateAICommand(input)
		if err == nil {
			t.Errorf("ValidateAICommand(%q): expected error, got nil", input)
			continue
		}
		if !errors.Is(err, ErrInvalidCommand) {
			t.Errorf("ValidateAICommand(%q): expected ErrInvalidCommand, got %v", input, err)
		}
	}
}

func TestValidateAICommand_NotFound(t *testing.T) {
	// Use a name that is extremely unlikely to exist on PATH.
	err := ValidateAICommand("ralph-nonexistent-xyz-12345-binary")
	if err == nil {
		t.Fatal("ValidateAICommand(nonexistent): expected error, got nil")
	}
	if errors.Is(err, ErrInvalidCommand) {
		t.Errorf("expected not-found style error, got ErrInvalidCommand: %v", err)
	}
	if msg := err.Error(); msg == "" || len(msg) < 10 {
		t.Errorf("error message should be clear and mention command: %q", msg)
	}
}

func TestValidateAICommand_ExecutableOnPATH(t *testing.T) {
	// "go" is expected to be on PATH in a Go project.
	err := ValidateAICommand("go")
	if err != nil {
		t.Errorf("ValidateAICommand(\"go\"): expected nil (go is on PATH), got %v", err)
	}
}

func TestValidateAICommand_WithArgs(t *testing.T) {
	// Validation only checks the executable (first word); args are ignored.
	err := ValidateAICommand("go build -o bin/ralph ./cmd/ralph")
	if err != nil {
		t.Errorf("ValidateAICommand(\"go build ...\"): expected nil, got %v", err)
	}
}

func TestFirstWord(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"go", "go"},
		{"  go  ", "go"},
		{"go build", "go"},
		{"claude --non-interactive", "claude"},
		{"\tagent\t-p", "agent"},
	}
	for _, tt := range tests {
		got := firstWord(tt.in)
		if got != tt.want {
			t.Errorf("firstWord(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
