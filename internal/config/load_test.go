package config

import (
	"os"
	"path/filepath"
	"testing"
)

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
