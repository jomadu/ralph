package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadGlobalAndWorkspace_bothMissing(t *testing.T) {
	// No files exist at global or workspace paths; skip without error (O002/R001).
	dir := t.TempDir()
	getenv := func(k string) string {
		if k == "RALPH_CONFIG_HOME" {
			return dir
		}
		return ""
	}
	global, workspace, err := LoadGlobalAndWorkspace(getenv, dir)
	if err != nil {
		t.Fatalf("LoadGlobalAndWorkspace(both missing) err = %v, want nil", err)
	}
	if global != nil || workspace != nil {
		t.Errorf("LoadGlobalAndWorkspace(both missing) = global=%v, workspace=%v, want nil, nil", global, workspace)
	}
}

func TestLoadGlobalAndWorkspace_globalPresent(t *testing.T) {
	dir := t.TempDir()
	globalPath := filepath.Join(dir, ConfigFileName)
	if err := os.WriteFile(globalPath, []byte("loop:\n  max_iterations: 7\n"), 0644); err != nil {
		t.Fatal(err)
	}
	getenv := func(k string) string {
		if k == "RALPH_CONFIG_HOME" {
			return dir
		}
		return ""
	}
	cwd := t.TempDir() // workspace path is cwd/ralph-config.yml which does not exist
	global, workspace, err := LoadGlobalAndWorkspace(getenv, cwd)
	if err != nil {
		t.Fatalf("LoadGlobalAndWorkspace(global present) err = %v", err)
	}
	if global == nil || global.Loop == nil || global.Loop.MaxIterations == nil || *global.Loop.MaxIterations != 7 {
		t.Errorf("LoadGlobalAndWorkspace(global present) global = %+v, want loop.max_iterations=7", global)
	}
	if workspace != nil {
		t.Errorf("LoadGlobalAndWorkspace(global present) workspace = %v, want nil", workspace)
	}
}

func TestLoadGlobalAndWorkspace_workspacePresent(t *testing.T) {
	globalDir := t.TempDir() // no file there
	workspaceDir := t.TempDir()
	workspacePath := filepath.Join(workspaceDir, ConfigFileName)
	if err := os.WriteFile(workspacePath, []byte("prompts:\n  p1:\n    path: p1.md\n"), 0644); err != nil {
		t.Fatal(err)
	}
	getenv := func(k string) string {
		if k == "RALPH_CONFIG_HOME" {
			return globalDir
		}
		return ""
	}
	global, workspace, err := LoadGlobalAndWorkspace(getenv, workspaceDir)
	if err != nil {
		t.Fatalf("LoadGlobalAndWorkspace(workspace present) err = %v", err)
	}
	if global != nil {
		t.Errorf("LoadGlobalAndWorkspace(workspace present) global = %v, want nil", global)
	}
	if workspace == nil || workspace.Prompts["p1"].Path != "p1.md" {
		t.Errorf("LoadGlobalAndWorkspace(workspace present) workspace = %+v", workspace)
	}
}

func TestLoadGlobalAndWorkspace_bothPresent(t *testing.T) {
	dir := t.TempDir()
	globalPath := filepath.Join(dir, ConfigFileName)
	if err := os.WriteFile(globalPath, []byte("loop:\n  max_iterations: 3\n"), 0644); err != nil {
		t.Fatal(err)
	}
	workspaceDir := t.TempDir()
	workspacePath := filepath.Join(workspaceDir, ConfigFileName)
	if err := os.WriteFile(workspacePath, []byte("loop:\n  max_iterations: 5\nprompts:\n  w: { path: w.md }\n"), 0644); err != nil {
		t.Fatal(err)
	}
	getenv := func(k string) string {
		if k == "RALPH_CONFIG_HOME" {
			return dir
		}
		return ""
	}
	global, workspace, err := LoadGlobalAndWorkspace(getenv, workspaceDir)
	if err != nil {
		t.Fatalf("LoadGlobalAndWorkspace(both present) err = %v", err)
	}
	if global == nil || global.Loop == nil || global.Loop.MaxIterations == nil || *global.Loop.MaxIterations != 3 {
		t.Errorf("LoadGlobalAndWorkspace(both present) global = %+v", global)
	}
	if workspace == nil || workspace.Loop == nil || workspace.Loop.MaxIterations == nil || *workspace.Loop.MaxIterations != 5 ||
		workspace.Prompts["w"].Path != "w.md" {
		t.Errorf("LoadGlobalAndWorkspace(both present) workspace = %+v", workspace)
	}
}

func TestLoadExplicit_missing(t *testing.T) {
	_, err := LoadExplicit("/nonexistent/config.yml")
	if err == nil {
		t.Fatal("LoadExplicit(missing) err = nil, want error")
	}
}

func TestLoadExplicit_valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "explicit.yml")
	if err := os.WriteFile(path, []byte("prompts:\n  p1:\n    path: p1.md\n"), 0644); err != nil {
		t.Fatal(err)
	}
	layer, err := LoadExplicit(path)
	if err != nil {
		t.Fatalf("LoadExplicit(valid) err = %v", err)
	}
	if layer == nil || layer.Prompts["p1"].Path != "p1.md" {
		t.Errorf("LoadExplicit(valid) layer = %+v", layer)
	}
}

func TestResolveEffective_noConfig_envOverlay(t *testing.T) {
	dir := t.TempDir()
	getenv := func(k string) string {
		switch k {
		case "RALPH_CONFIG_HOME":
			return dir
		case "RALPH_LOOP_DEFAULT_MAX_ITERATIONS":
			return "5"
		case "RALPH_LOOP_LOG_LEVEL":
			return "warn"
		default:
			return ""
		}
	}
	eff, err := ResolveEffective(getenv, dir, "")
	if err != nil {
		t.Fatalf("ResolveEffective() err = %v", err)
	}
	if eff.Loop.MaxIterations != 5 {
		t.Errorf("MaxIterations = %d, want 5 (from env)", eff.Loop.MaxIterations)
	}
	if eff.Loop.LogLevel != "warn" {
		t.Errorf("LogLevel = %q, want warn (from env)", eff.Loop.LogLevel)
	}
}

func TestResolveEffective_invalidEnv_clearError(t *testing.T) {
	dir := t.TempDir()
	getenv := func(k string) string {
		if k == "RALPH_LOOP_DEFAULT_MAX_ITERATIONS" {
			return "not-a-number"
		}
		return ""
	}
	_, err := ResolveEffective(getenv, dir, "")
	if err == nil {
		t.Fatal("ResolveEffective(invalid env) err = nil, want error")
	}
	if err.Error() == "" {
		t.Error("error message should be non-empty and mention the variable")
	}
}

// TestResolveEffectiveForPrompt_promptOverride verifies that when resolving for a
// named prompt, prompt-level loop overrides are applied (T1.6: defaults → … → env → prompt overrides).
func TestResolveEffectiveForPrompt_promptOverride(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "ralph-config.yml")
	cfg := `
loop:
  max_iterations: 10
  log_level: info
prompts:
  myprompt:
    path: p.md
    loop:
      max_iterations: 2
      log_level: debug
`
	if err := os.WriteFile(configPath, []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}
	getenv := func(string) string { return "" }

	// Resolve for the named prompt: effective loop must have prompt overrides (2, debug).
	eff, ok, err := ResolveEffectiveForPrompt(getenv, dir, configPath, "myprompt")
	if err != nil {
		t.Fatalf("ResolveEffectiveForPrompt err = %v", err)
	}
	if !ok || eff == nil {
		t.Fatalf("ResolveEffectiveForPrompt(_, _, _, %q) = %v, %v; want effective, true", "myprompt", eff, ok)
	}
	if eff.Loop.MaxIterations != 2 {
		t.Errorf("Loop.MaxIterations = %d, want 2 (prompt override)", eff.Loop.MaxIterations)
	}
	if eff.Loop.LogLevel != "debug" {
		t.Errorf("Loop.LogLevel = %q, want debug (prompt override)", eff.Loop.LogLevel)
	}

	// Root (no prompt name) or empty prompt name: root loop only (10, info).
	root, okRoot, err := ResolveEffectiveForPrompt(getenv, dir, configPath, "")
	if err != nil || !okRoot || root == nil {
		t.Fatalf("ResolveEffectiveForPrompt(empty name) err=%v ok=%v root=%v", err, okRoot, root)
	}
	if root.Loop.MaxIterations != 10 || root.Loop.LogLevel != "info" {
		t.Errorf("root Loop = max_iter=%d log_level=%q, want 10, info", root.Loop.MaxIterations, root.Loop.LogLevel)
	}

	// Unknown prompt: not found.
	_, okUnknown, err := ResolveEffectiveForPrompt(getenv, dir, configPath, "nonexistent")
	if err != nil {
		t.Fatalf("ResolveEffectiveForPrompt(nonexistent) err = %v", err)
	}
	if okUnknown {
		t.Error("ResolveEffectiveForPrompt(nonexistent) ok = true, want false")
	}
}
