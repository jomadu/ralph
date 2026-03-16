package config

import (
	"testing"
)

func TestParseEnvOverlay_empty(t *testing.T) {
	getenv := func(string) string { return "" }
	o, err := ParseEnvOverlay(getenv)
	if err != nil {
		t.Fatalf("ParseEnvOverlay() err = %v", err)
	}
	if o == nil {
		t.Fatal("ParseEnvOverlay() returned nil overlay")
	}
	if o.AICmd != nil || o.MaxIterations != nil || o.Streaming != nil || o.MaxOutputBuffer != nil {
		t.Errorf("expected all nil when no env set, got AICmd=%v MaxIterations=%v Streaming=%v MaxOutputBuffer=%v", o.AICmd, o.MaxIterations, o.Streaming, o.MaxOutputBuffer)
	}
}

func TestParseEnvOverlay_strings(t *testing.T) {
	getenv := func(k string) string {
		switch k {
		case "RALPH_LOOP_AI_CMD":
			return "claude --non-interactive"
		case "RALPH_LOOP_AI_CMD_ALIAS":
			return "my-alias"
		case "RALPH_LOOP_LOG_LEVEL":
			return "debug"
		default:
			return ""
		}
	}
	o, err := ParseEnvOverlay(getenv)
	if err != nil {
		t.Fatalf("ParseEnvOverlay() err = %v", err)
	}
	if o.AICmd == nil || *o.AICmd != "claude --non-interactive" {
		t.Errorf("AICmd = %v", o.AICmd)
	}
	if o.AICmdAlias == nil || *o.AICmdAlias != "my-alias" {
		t.Errorf("AICmdAlias = %v", o.AICmdAlias)
	}
	if o.LogLevel == nil || *o.LogLevel != "debug" {
		t.Errorf("LogLevel = %v", o.LogLevel)
	}
}

func TestParseEnvOverlay_ints(t *testing.T) {
	getenv := func(k string) string {
		switch k {
		case "RALPH_LOOP_DEFAULT_MAX_ITERATIONS":
			return "10"
		case "RALPH_LOOP_FAILURE_THRESHOLD":
			return "3"
		case "RALPH_LOOP_ITERATION_TIMEOUT":
			return "60"
		default:
			return ""
		}
	}
	o, err := ParseEnvOverlay(getenv)
	if err != nil {
		t.Fatalf("ParseEnvOverlay() err = %v", err)
	}
	if o.MaxIterations == nil || *o.MaxIterations != 10 {
		t.Errorf("MaxIterations = %v", o.MaxIterations)
	}
	if o.FailureThreshold == nil || *o.FailureThreshold != 3 {
		t.Errorf("FailureThreshold = %v", o.FailureThreshold)
	}
	if o.IterationTimeout == nil || *o.IterationTimeout != 60 {
		t.Errorf("IterationTimeout = %v", o.IterationTimeout)
	}
}

func TestParseEnvOverlay_int_invalid(t *testing.T) {
	getenv := func(k string) string {
		if k == "RALPH_LOOP_DEFAULT_MAX_ITERATIONS" {
			return "not-a-number"
		}
		return ""
	}
	_, err := ParseEnvOverlay(getenv)
	if err == nil {
		t.Fatal("expected error for invalid integer")
	}
	if err.Error() == "" {
		t.Error("error message should be non-empty")
	}
}

func TestParseEnvOverlay_int_negative(t *testing.T) {
	getenv := func(k string) string {
		if k == "RALPH_LOOP_DEFAULT_MAX_ITERATIONS" {
			return "-1"
		}
		return ""
	}
	_, err := ParseEnvOverlay(getenv)
	if err == nil {
		t.Fatal("expected error for negative max iterations")
	}
}

func TestParseEnvOverlay_bool_true(t *testing.T) {
	for _, v := range []string{"true", "1", "yes", "on"} {
		getenv := func(k string) string {
			if k == "RALPH_LOOP_STREAMING" {
				return v
			}
			return ""
		}
		o, err := ParseEnvOverlay(getenv)
		if err != nil {
			t.Fatalf("ParseEnvOverlay(%q) err = %v", v, err)
		}
		if o.Streaming == nil || !*o.Streaming {
			t.Errorf("Streaming for %q = %v", v, o.Streaming)
		}
	}
}

func TestParseEnvOverlay_bool_false(t *testing.T) {
	for _, v := range []string{"false", "0", "no", "off"} {
		getenv := func(k string) string {
			if k == "RALPH_LOOP_STREAMING" {
				return v
			}
			return ""
		}
		o, err := ParseEnvOverlay(getenv)
		if err != nil {
			t.Fatalf("ParseEnvOverlay(%q) err = %v", v, err)
		}
		if o.Streaming == nil || *o.Streaming {
			t.Errorf("Streaming for %q = %v", v, o.Streaming)
		}
	}
}

func TestParseEnvOverlay_bool_invalid(t *testing.T) {
	getenv := func(k string) string {
		if k == "RALPH_LOOP_STREAMING" {
			return "invalid"
		}
		return ""
	}
	_, err := ParseEnvOverlay(getenv)
	if err == nil {
		t.Fatal("expected error for invalid boolean")
	}
}

func TestParseEnvOverlay_iteration_timeout_zero(t *testing.T) {
	getenv := func(k string) string {
		if k == "RALPH_LOOP_ITERATION_TIMEOUT" {
			return "0"
		}
		return ""
	}
	o, err := ParseEnvOverlay(getenv)
	if err != nil {
		t.Fatalf("ParseEnvOverlay() err = %v", err)
	}
	if o.IterationTimeout == nil || *o.IterationTimeout != 0 {
		t.Errorf("IterationTimeout = %v (0 = no timeout)", o.IterationTimeout)
	}
}

func TestParseEnvOverlay_max_output_buffer(t *testing.T) {
	getenv := func(k string) string {
		if k == "RALPH_LOOP_MAX_OUTPUT_BUFFER" {
			return "65536"
		}
		return ""
	}
	o, err := ParseEnvOverlay(getenv)
	if err != nil {
		t.Fatalf("ParseEnvOverlay() err = %v", err)
	}
	if o.MaxOutputBuffer == nil || *o.MaxOutputBuffer != 65536 {
		t.Errorf("MaxOutputBuffer = %v, want 65536", o.MaxOutputBuffer)
	}
}

func TestParseEnvOverlay_max_output_buffer_negative(t *testing.T) {
	getenv := func(k string) string {
		if k == "RALPH_LOOP_MAX_OUTPUT_BUFFER" {
			return "-1"
		}
		return ""
	}
	_, err := ParseEnvOverlay(getenv)
	if err == nil {
		t.Fatal("expected error for negative RALPH_LOOP_MAX_OUTPUT_BUFFER")
	}
}

func TestParseEnvOverlay_max_output_buffer_invalid(t *testing.T) {
	getenv := func(k string) string {
		if k == "RALPH_LOOP_MAX_OUTPUT_BUFFER" {
			return "not-a-number"
		}
		return ""
	}
	_, err := ParseEnvOverlay(getenv)
	if err == nil {
		t.Fatal("expected error for invalid RALPH_LOOP_MAX_OUTPUT_BUFFER")
	}
}
