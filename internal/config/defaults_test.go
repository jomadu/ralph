package config

import (
	"testing"
)

func TestDefaultLoopSettings(t *testing.T) {
	d := DefaultLoopSettings()
	if d.MaxIterations <= 0 {
		t.Errorf("MaxIterations = %d, want positive", d.MaxIterations)
	}
	if d.FailureThreshold <= 0 {
		t.Errorf("FailureThreshold = %d, want positive", d.FailureThreshold)
	}
	if d.TimeoutSeconds != 0 {
		t.Errorf("TimeoutSeconds = %d, want 0 (no timeout)", d.TimeoutSeconds)
	}
	if d.SuccessSignal != DefaultSuccessSignal {
		t.Errorf("SuccessSignal = %q, want %q", d.SuccessSignal, DefaultSuccessSignal)
	}
	if d.FailureSignal != DefaultFailureSignal {
		t.Errorf("FailureSignal = %q, want %q", d.FailureSignal, DefaultFailureSignal)
	}
	if !d.Preamble {
		t.Errorf("Preamble = false, want true (included by default)")
	}
	if d.Context != "" {
		t.Errorf("Context = %q, want empty", d.Context)
	}
	if !d.Streaming {
		t.Errorf("Streaming = false, want true")
	}
	if d.LogLevel == "" {
		t.Errorf("LogLevel empty, want e.g. info")
	}
	if d.MaxOutputBuffer != DefaultMaxOutputBuffer {
		t.Errorf("MaxOutputBuffer = %d, want %d (DefaultMaxOutputBuffer)", d.MaxOutputBuffer, DefaultMaxOutputBuffer)
	}
}

func TestDefaultEffective(t *testing.T) {
	e := DefaultEffective()
	if e == nil {
		t.Fatal("DefaultEffective() = nil")
	}
	// Loop should match DefaultLoopSettings
	if e.Loop.MaxIterations != DefaultLoopSettings().MaxIterations {
		t.Errorf("Loop.MaxIterations = %d, want default", e.Loop.MaxIterations)
	}
	// Prompts empty
	if len(e.Prompts) != 0 {
		t.Errorf("Prompts has %d entries, want 0", len(e.Prompts))
	}
	// Built-in aliases present
	wantAliases := []string{"claude", "kiro", "copilot", "cursor-agent"}
	for _, name := range wantAliases {
		if _, ok := e.Aliases[name]; !ok {
			t.Errorf("Aliases[%q] missing", name)
		}
	}
}
