package config

import (
	"errors"
	"testing"
)

func TestValidateFileLayer_nil(t *testing.T) {
	if err := ValidateFileLayer(nil); err != nil {
		t.Errorf("ValidateFileLayer(nil) err = %v, want nil", err)
	}
}

func TestValidateFileLayer_negativeMaxIterations(t *testing.T) {
	layer := &FileLayer{
		Loop: &LoopSection{MaxIterations: intPtr(-1)},
	}
	err := ValidateFileLayer(layer)
	if err == nil {
		t.Fatal("ValidateFileLayer(negative max_iterations) err = nil, want error")
	}
	if !errors.Is(err, ErrInvalidMaxIterations) {
		t.Errorf("err = %v, want ErrInvalidMaxIterations", err)
	}
}

func TestValidateFileLayer_negativeFailureThreshold(t *testing.T) {
	layer := &FileLayer{
		Loop: &LoopSection{FailureThreshold: intPtr(-2)},
	}
	err := ValidateFileLayer(layer)
	if err == nil {
		t.Fatal("ValidateFileLayer(negative failure_threshold) err = nil, want error")
	}
	if !errors.Is(err, ErrInvalidFailureThreshold) {
		t.Errorf("err = %v, want ErrInvalidFailureThreshold", err)
	}
}

func TestValidateFileLayer_negativeTimeoutSeconds(t *testing.T) {
	layer := &FileLayer{
		Loop: &LoopSection{TimeoutSeconds: intPtr(-5)},
	}
	err := ValidateFileLayer(layer)
	if err == nil {
		t.Fatal("ValidateFileLayer(negative timeout_seconds) err = nil, want error")
	}
	if !errors.Is(err, ErrInvalidTimeoutSeconds) {
		t.Errorf("err = %v, want ErrInvalidTimeoutSeconds", err)
	}
}

func TestValidateFileLayer_invalidSignalPrecedence(t *testing.T) {
	layer := &FileLayer{
		Loop: &LoopSection{SignalPrecedence: "unknown"},
	}
	err := ValidateFileLayer(layer)
	if err == nil {
		t.Fatal("ValidateFileLayer(invalid signal_precedence) err = nil, want error")
	}
	if !errors.Is(err, ErrInvalidSignalPrecedence) {
		t.Errorf("err = %v, want ErrInvalidSignalPrecedence", err)
	}
}

func TestValidateFileLayer_validSignalPrecedence(t *testing.T) {
	for _, v := range []string{"", "static", "ai_interpreted"} {
		layer := &FileLayer{Loop: &LoopSection{SignalPrecedence: v}}
		if err := ValidateFileLayer(layer); err != nil {
			t.Errorf("ValidateFileLayer(signal_precedence=%q) err = %v", v, err)
		}
	}
}

func TestValidateFileLayer_invalidLogLevel(t *testing.T) {
	layer := &FileLayer{
		Loop: &LoopSection{LogLevel: "trace"},
	}
	err := ValidateFileLayer(layer)
	if err == nil {
		t.Fatal("ValidateFileLayer(invalid log_level) err = nil, want error")
	}
	if !errors.Is(err, ErrInvalidLogLevel) {
		t.Errorf("err = %v, want ErrInvalidLogLevel", err)
	}
}

func TestValidateFileLayer_validLogLevel(t *testing.T) {
	for _, v := range []string{"", "debug", "info", "warn", "error"} {
		layer := &FileLayer{Loop: &LoopSection{LogLevel: v}}
		if err := ValidateFileLayer(layer); err != nil {
			t.Errorf("ValidateFileLayer(log_level=%q) err = %v", v, err)
		}
	}
}

func TestValidateFileLayer_promptNeedsPathOrContent(t *testing.T) {
	layer := &FileLayer{
		Prompts: map[string]Prompt{
			"empty": {}, // neither path nor content
		},
	}
	err := ValidateFileLayer(layer)
	if err == nil {
		t.Fatal("ValidateFileLayer(prompt no path/content) err = nil, want error")
	}
	if !errors.Is(err, ErrPromptNeedsPathOrContent) {
		t.Errorf("err = %v, want ErrPromptNeedsPathOrContent", err)
	}
}

func TestValidateFileLayer_promptWithPathOK(t *testing.T) {
	layer := &FileLayer{
		Prompts: map[string]Prompt{
			"p1": {Path: "main.md"},
		},
	}
	if err := ValidateFileLayer(layer); err != nil {
		t.Errorf("ValidateFileLayer(prompt path) err = %v", err)
	}
}

func TestValidateFileLayer_promptWithContentOK(t *testing.T) {
	layer := &FileLayer{
		Prompts: map[string]Prompt{
			"p1": {Content: "inline prompt"},
		},
	}
	if err := ValidateFileLayer(layer); err != nil {
		t.Errorf("ValidateFileLayer(prompt content) err = %v", err)
	}
}

func TestValidateFileLayer_aliasCommandEmpty(t *testing.T) {
	layer := &FileLayer{
		Aliases: map[string]Alias{
			"bad": {Command: ""},
		},
	}
	err := ValidateFileLayer(layer)
	if err == nil {
		t.Fatal("ValidateFileLayer(empty alias command) err = nil, want error")
	}
	if !errors.Is(err, ErrAliasCommandEmpty) {
		t.Errorf("err = %v, want ErrAliasCommandEmpty", err)
	}
}

func TestValidateFileLayer_aliasCommandOK(t *testing.T) {
	layer := &FileLayer{
		Aliases: map[string]Alias{
			"good": {Command: "claude --non-interactive"},
		},
	}
	if err := ValidateFileLayer(layer); err != nil {
		t.Errorf("ValidateFileLayer(alias command) err = %v", err)
	}
}

func TestParseLayer_rejectsInvalid(t *testing.T) {
	yaml := "loop:\n  max_iterations: -1\n"
	_, err := ParseLayer([]byte(yaml))
	if err == nil {
		t.Fatal("ParseLayer(negative max_iterations) err = nil, want error")
	}
	if !errors.Is(err, ErrInvalidMaxIterations) {
		t.Errorf("err = %v, want ErrInvalidMaxIterations", err)
	}
}

func intPtr(n int) *int { return &n }
