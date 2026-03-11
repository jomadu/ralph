package config

import (
	"path/filepath"
	"testing"
)

func TestGlobalPath(t *testing.T) {
	t.Run("RALPH_CONFIG_HOME set", func(t *testing.T) {
		getenv := func(k string) string {
			if k == "RALPH_CONFIG_HOME" {
				return "/custom/ralph"
			}
			return ""
		}
		got := GlobalPath(getenv)
		want := filepath.Join("/custom/ralph", ConfigFileName)
		if got != want {
			t.Errorf("GlobalPath() = %q, want %q", got, want)
		}
	})

	t.Run("XDG_CONFIG_HOME fallback", func(t *testing.T) {
		getenv := func(k string) string {
			if k == "XDG_CONFIG_HOME" {
				return "/xdg/config"
			}
			return ""
		}
		got := GlobalPath(getenv)
		want := filepath.Join("/xdg/config", "ralph", ConfigFileName)
		if got != want {
			t.Errorf("GlobalPath() = %q, want %q", got, want)
		}
	})

	t.Run("RALPH_CONFIG_HOME overrides XDG", func(t *testing.T) {
		getenv := func(k string) string {
			if k == "RALPH_CONFIG_HOME" {
				return "/ralph-home"
			}
			if k == "XDG_CONFIG_HOME" {
				return "/xdg"
			}
			return ""
		}
		got := GlobalPath(getenv)
		want := filepath.Join("/ralph-home", ConfigFileName)
		if got != want {
			t.Errorf("GlobalPath() = %q, want %q", got, want)
		}
	})
}

func TestWorkspacePath(t *testing.T) {
	cwd := "/project/root"
	got := WorkspacePath(cwd)
	want := filepath.Join(cwd, ConfigFileName)
	if got != want {
		t.Errorf("WorkspacePath(%q) = %q, want %q", cwd, got, want)
	}
}
