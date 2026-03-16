package config

import (
	"errors"
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
	// Pass paths that don't exist; no getenv/cwd, so test is independent of dev machine state.
	globalPath := filepath.Join(t.TempDir(), ConfigFileName)
	workspacePath := filepath.Join(t.TempDir(), ConfigFileName)
	global, workspace, err := LoadGlobalAndWorkspace(globalPath, workspacePath)
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
	workspaceDir := t.TempDir()
	workspacePath := filepath.Join(workspaceDir, ConfigFileName)
	if err := os.WriteFile(workspacePath, []byte("loop:\n  max_iterations: 3\n"), 0644); err != nil {
		t.Fatal(err)
	}
	globalPath := filepath.Join(t.TempDir(), ConfigFileName) // missing
	global, workspace, err := LoadGlobalAndWorkspace(globalPath, workspacePath)
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

func TestReadLayerRequired_missing(t *testing.T) {
	_, err := ReadLayerRequired("/nonexistent/path/ralph-config.yml")
	if err == nil {
		t.Fatal("ReadLayerRequired(missing) err = nil, want error")
	}
	if !errors.Is(err, ErrExplicitConfigMissing) {
		t.Errorf("ReadLayerRequired(missing) err = %v, want ErrExplicitConfigMissing", err)
	}
}

func TestReadLayerRequired_directory(t *testing.T) {
	dir := t.TempDir()
	_, err := ReadLayerRequired(dir)
	if err == nil {
		t.Fatal("ReadLayerRequired(directory) err = nil, want error")
	}
	if !errors.Is(err, ErrExplicitConfigMissing) {
		t.Errorf("ReadLayerRequired(directory) err = %v, want ErrExplicitConfigMissing", err)
	}
}

func TestReadLayerRequired_emptyPath(t *testing.T) {
	_, err := ReadLayerRequired("")
	if err == nil {
		t.Fatal("ReadLayerRequired(empty) err = nil, want error")
	}
	if !errors.Is(err, ErrExplicitConfigMissing) {
		t.Errorf("ReadLayerRequired(empty) err = %v, want ErrExplicitConfigMissing", err)
	}
}

func TestReadLayerRequired_valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "custom.yml")
	if err := os.WriteFile(path, []byte("loop:\n  max_iterations: 5\n"), 0644); err != nil {
		t.Fatal(err)
	}
	layer, err := ReadLayerRequired(path)
	if err != nil {
		t.Fatalf("ReadLayerRequired(valid) err = %v", err)
	}
	if layer == nil || layer.Loop == nil || *layer.Loop.MaxIterations != 5 {
		t.Errorf("ReadLayerRequired(valid) layer = %+v", layer)
	}
}
