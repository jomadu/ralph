package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadLayer_missing(t *testing.T) {
	layer, err := ReadLayer("/nonexistent/path/ralph-config.yml")
	if err != nil {
		t.Fatalf("ReadLayer(missing) err = %v, want nil (skip without error)", err)
	}
	if layer != nil {
		t.Errorf("ReadLayer(missing) layer = %+v, want nil", layer)
	}
}

func TestParseLayer_empty(t *testing.T) {
	layer, err := ParseLayer(nil)
	if err != nil {
		t.Fatalf("ParseLayer(nil) err = %v", err)
	}
	if layer != nil {
		t.Errorf("ParseLayer(nil) layer = %+v, want nil", layer)
	}
	layer, err = ParseLayer([]byte{})
	if err != nil {
		t.Fatalf("ParseLayer([]) err = %v", err)
	}
	if layer != nil {
		t.Errorf("ParseLayer([]) layer = %+v, want nil", layer)
	}
}

func TestParseLayer_valid(t *testing.T) {
	yaml := `
loop:
  max_iterations: 5
  failure_threshold: 2
prompts:
  default:
    path: prompts/main.md
aliases:
  claude: "claude --non-interactive"
  kiro:
    command: "kiro run"
`
	layer, err := ParseLayer([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseLayer() err = %v", err)
	}
	if layer == nil {
		t.Fatal("ParseLayer() layer = nil")
	}
	if layer.Loop == nil || *layer.Loop.MaxIterations != 5 || *layer.Loop.FailureThreshold != 2 {
		t.Errorf("loop: got %+v", layer.Loop)
	}
	if layer.Prompts["default"].Path != "prompts/main.md" {
		t.Errorf("prompts[default].Path = %q", layer.Prompts["default"].Path)
	}
	if layer.Aliases["claude"].Command != "claude --non-interactive" {
		t.Errorf("aliases[claude].Command = %q", layer.Aliases["claude"].Command)
	}
	if layer.Aliases["kiro"].Command != "kiro run" {
		t.Errorf("aliases[kiro].Command = %q", layer.Aliases["kiro"].Command)
	}
}

func TestLoadGlobalAndWorkspace_skip_missing(t *testing.T) {
	dir := t.TempDir()
	getenv := func(string) string { return "" } // no RALPH_CONFIG_HOME or XDG
	global, workspace, err := LoadGlobalAndWorkspace(getenv, dir)
	if err != nil {
		t.Fatalf("LoadGlobalAndWorkspace() err = %v", err)
	}
	if global != nil {
		t.Errorf("global = %+v, want nil (missing)", global)
	}
	if workspace != nil {
		t.Errorf("workspace = %+v, want nil (missing)", workspace)
	}
}

func TestLoadGlobalAndWorkspace_read_workspace(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ConfigFileName)
	if err := os.WriteFile(path, []byte("loop:\n  max_iterations: 3\n"), 0644); err != nil {
		t.Fatal(err)
	}
	// Use a nonexistent global config dir so only workspace is loaded
	globalDir := t.TempDir()
	getenv := func(k string) string {
		if k == "RALPH_CONFIG_HOME" {
			return globalDir
		}
		return ""
	}
	global, workspace, err := LoadGlobalAndWorkspace(getenv, dir)
	if err != nil {
		t.Fatalf("LoadGlobalAndWorkspace() err = %v", err)
	}
	if global != nil {
		t.Errorf("global = %+v, want nil", global)
	}
	if workspace == nil || workspace.Loop == nil || *workspace.Loop.MaxIterations != 3 {
		t.Errorf("workspace = %+v", workspace)
	}
}
