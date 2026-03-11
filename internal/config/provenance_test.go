package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRootLoopWithProvenance_defaults(t *testing.T) {
	cwd := t.TempDir()
	getenv := func(string) string { return "" }
	loop, prov, err := RootLoopWithProvenance(getenv, cwd, "")
	if err != nil {
		t.Fatalf("RootLoopWithProvenance() err = %v", err)
	}
	if loop.MaxIterations != 10 || loop.LogLevel != "info" {
		t.Errorf("loop = %+v, want default values", loop)
	}
	if prov.MaxIterations != ProvenanceDefault || prov.LogLevel != ProvenanceDefault {
		t.Errorf("provenance = %+v, want all default", prov)
	}
}

func TestRootLoopWithProvenance_envOverlay(t *testing.T) {
	cwd := t.TempDir()
	getenv := func(k string) string {
		switch k {
		case "RALPH_LOOP_DEFAULT_MAX_ITERATIONS":
			return "5"
		case "RALPH_LOOP_LOG_LEVEL":
			return "debug"
		default:
			return ""
		}
	}
	loop, prov, err := RootLoopWithProvenance(getenv, cwd, "")
	if err != nil {
		t.Fatalf("RootLoopWithProvenance(env) err = %v", err)
	}
	if loop.MaxIterations != 5 || loop.LogLevel != "debug" {
		t.Errorf("loop = max_iterations=%d log_level=%q, want 5, debug", loop.MaxIterations, loop.LogLevel)
	}
	if prov.MaxIterations != ProvenanceEnv || prov.LogLevel != ProvenanceEnv {
		t.Errorf("provenance = max_iterations=%q log_level=%q, want env, env", prov.MaxIterations, prov.LogLevel)
	}
	if prov.FailureThreshold != ProvenanceDefault {
		t.Errorf("provenance.FailureThreshold = %q, want default (not overridden by env)", prov.FailureThreshold)
	}
}

func TestRootLoopWithProvenance_explicitFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "ralph-config.yml")
	if err := os.WriteFile(configPath, []byte("loop:\n  max_iterations: 2\n  log_level: warn\n"), 0644); err != nil {
		t.Fatal(err)
	}
	getenv := func(string) string { return "" }
	loop, prov, err := RootLoopWithProvenance(getenv, dir, configPath)
	if err != nil {
		t.Fatalf("RootLoopWithProvenance(explicit) err = %v", err)
	}
	if loop.MaxIterations != 2 || loop.LogLevel != "warn" {
		t.Errorf("loop = max_iterations=%d log_level=%q, want 2, warn", loop.MaxIterations, loop.LogLevel)
	}
	if prov.MaxIterations != ProvenanceExplicit || prov.LogLevel != ProvenanceExplicit {
		t.Errorf("provenance = max_iterations=%q log_level=%q, want explicit, explicit", prov.MaxIterations, prov.LogLevel)
	}
}
