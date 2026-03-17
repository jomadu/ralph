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
	if prov.AICmd != ProvenanceDefault || prov.AICmdAlias != ProvenanceDefault {
		t.Errorf("provenance AICmd=%q AICmdAlias=%q, want both default", prov.AICmd, prov.AICmdAlias)
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

func TestRootLoopWithProvenance_explicitFile_aiCmd(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "ralph-config.yml")
	cfg := "loop:\n  ai_cmd: \"claude --non-interactive\"\n  ai_cmd_alias: cursor-agent\n"
	if err := os.WriteFile(configPath, []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}
	getenv := func(string) string { return "" }
	loop, prov, err := RootLoopWithProvenance(getenv, dir, configPath)
	if err != nil {
		t.Fatalf("RootLoopWithProvenance(explicit ai_cmd) err = %v", err)
	}
	if loop.AICmd != "claude --non-interactive" || loop.AICmdAlias != "cursor-agent" {
		t.Errorf("loop = AICmd=%q AICmdAlias=%q, want explicit values", loop.AICmd, loop.AICmdAlias)
	}
	if prov.AICmd != ProvenanceExplicit || prov.AICmdAlias != ProvenanceExplicit {
		t.Errorf("provenance AICmd=%q AICmdAlias=%q, want explicit, explicit", prov.AICmd, prov.AICmdAlias)
	}
}

func TestRootLoopWithProvenance_envOverlay_aiCmd(t *testing.T) {
	cwd := t.TempDir()
	getenv := func(k string) string {
		switch k {
		case "RALPH_LOOP_AI_CMD":
			return "custom-ai --batch"
		case "RALPH_LOOP_AI_CMD_ALIAS":
			return "claude"
		default:
			return ""
		}
	}
	loop, prov, err := RootLoopWithProvenance(getenv, cwd, "")
	if err != nil {
		t.Fatalf("RootLoopWithProvenance(env ai_cmd) err = %v", err)
	}
	if loop.AICmd != "custom-ai --batch" || loop.AICmdAlias != "claude" {
		t.Errorf("loop = AICmd=%q AICmdAlias=%q, want env values", loop.AICmd, loop.AICmdAlias)
	}
	if prov.AICmd != ProvenanceEnv || prov.AICmdAlias != ProvenanceEnv {
		t.Errorf("provenance AICmd=%q AICmdAlias=%q, want env, env", prov.AICmd, prov.AICmdAlias)
	}
}

func TestRootLoopWithProvenance_workspaceOnly_aiCmdAlias(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "ralph-config.yml")
	cfg := "loop:\n  ai_cmd_alias: cursor-agent\n"
	if err := os.WriteFile(configPath, []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}
	getenv := func(string) string { return "" }
	loop, prov, err := RootLoopWithProvenance(getenv, dir, "")
	if err != nil {
		t.Fatalf("RootLoopWithProvenance(workspace) err = %v", err)
	}
	if loop.AICmdAlias != "cursor-agent" {
		t.Errorf("loop.AICmdAlias = %q, want cursor-agent", loop.AICmdAlias)
	}
	if prov.AICmdAlias != ProvenanceWorkspace {
		t.Errorf("provenance.AICmdAlias = %q, want workspace", prov.AICmdAlias)
	}
}

func TestRootLoopWithProvenance_globalOnly_aiCmdAlias(t *testing.T) {
	globalDir := t.TempDir()
	configPath := filepath.Join(globalDir, ConfigFileName)
	cfg := "loop:\n  ai_cmd_alias: kiro\n"
	if err := os.WriteFile(configPath, []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}
	cwd := t.TempDir()
	getenv := func(k string) string {
		if k == "RALPH_CONFIG_HOME" {
			return globalDir
		}
		return ""
	}
	loop, prov, err := RootLoopWithProvenance(getenv, cwd, "")
	if err != nil {
		t.Fatalf("RootLoopWithProvenance(global) err = %v", err)
	}
	if loop.AICmdAlias != "kiro" {
		t.Errorf("loop.AICmdAlias = %q, want kiro", loop.AICmdAlias)
	}
	if prov.AICmdAlias != ProvenanceGlobal {
		t.Errorf("provenance.AICmdAlias = %q, want global", prov.AICmdAlias)
	}
}

func TestLoopWithProvenance_cliOverlay(t *testing.T) {
	cwd := t.TempDir()
	getenv := func(string) string { return "" }
	// FailureThreshold/IterationTimeout -1 = not set (don't override)
	cli := &CLIOverlay{MaxIterations: 7, LogLevel: "warn", FailureThreshold: -1, IterationTimeout: -1}
	loop, prov, err := LoopWithProvenance(getenv, cwd, "", "", cli)
	if err != nil {
		t.Fatalf("LoopWithProvenance(cli) err = %v", err)
	}
	if loop.MaxIterations != 7 || loop.LogLevel != "warn" {
		t.Errorf("loop = max_iterations=%d log_level=%q, want 7, warn", loop.MaxIterations, loop.LogLevel)
	}
	if prov.MaxIterations != ProvenanceCLI || prov.LogLevel != ProvenanceCLI {
		t.Errorf("provenance = max_iterations=%q log_level=%q, want cli, cli", prov.MaxIterations, prov.LogLevel)
	}
	if prov.FailureThreshold != ProvenanceDefault {
		t.Errorf("provenance.FailureThreshold = %q, want default (not overridden by cli)", prov.FailureThreshold)
	}
}

func TestLoopWithProvenance_promptOverride(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "ralph-config.yml")
	cfg := "loop:\n  max_iterations: 5\nprompts:\n  p1:\n    path: \"x\"\n    loop:\n      max_iterations: 2\n      log_level: debug\n"
	if err := os.WriteFile(configPath, []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}
	getenv := func(string) string { return "" }
	loop, prov, err := LoopWithProvenance(getenv, dir, configPath, "p1", nil)
	if err != nil {
		t.Fatalf("LoopWithProvenance(prompt p1) err = %v", err)
	}
	if loop.MaxIterations != 2 || loop.LogLevel != "debug" {
		t.Errorf("loop = max_iterations=%d log_level=%q, want 2, debug (prompt overrides)", loop.MaxIterations, loop.LogLevel)
	}
	if prov.MaxIterations != ProvenancePrompt || prov.LogLevel != ProvenancePrompt {
		t.Errorf("provenance = max_iterations=%q log_level=%q, want prompt, prompt", prov.MaxIterations, prov.LogLevel)
	}
	// Root had max_iterations 5 from explicit file; failure_threshold from default
	if prov.FailureThreshold != ProvenanceDefault {
		t.Errorf("provenance.FailureThreshold = %q, want default", prov.FailureThreshold)
	}
}

func TestLoopWithProvenance_promptOverride_aiCmdAlias(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "ralph-config.yml")
	cfg := "loop:\n  ai_cmd_alias: claude\nprompts:\n  p1:\n    path: \"x\"\n    loop:\n      ai_cmd_alias: cursor-agent\n"
	if err := os.WriteFile(configPath, []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}
	getenv := func(string) string { return "" }
	loop, prov, err := LoopWithProvenance(getenv, dir, configPath, "p1", nil)
	if err != nil {
		t.Fatalf("LoopWithProvenance(prompt p1 ai_cmd_alias) err = %v", err)
	}
	if loop.AICmdAlias != "cursor-agent" {
		t.Errorf("loop.AICmdAlias = %q, want cursor-agent (prompt overrides root)", loop.AICmdAlias)
	}
	if prov.AICmdAlias != ProvenancePrompt {
		t.Errorf("provenance.AICmdAlias = %q, want prompt", prov.AICmdAlias)
	}
}
