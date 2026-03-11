// Package config resolves configuration from defined layers (defaults, global file,
// workspace file, explicit file, environment, CLI). See docs/engineering/components/config.md.
package config

import (
	"os"
	"path/filepath"
)

// ConfigFileName is the fixed config filename in any config directory.
const ConfigFileName = "ralph-config.yml"

// GlobalPath returns the path to the global config file.
// Order: $RALPH_CONFIG_HOME/ralph-config.yml, then $XDG_CONFIG_HOME/ralph/ralph-config.yml,
// then ~/.config/ralph/ralph-config.yml. getenv is typically os.Getenv.
func GlobalPath(getenv func(string) string) string {
	if d := getenv("RALPH_CONFIG_HOME"); d != "" {
		return filepath.Join(d, ConfigFileName)
	}
	if d := getenv("XDG_CONFIG_HOME"); d != "" {
		return filepath.Join(d, "ralph", ConfigFileName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "ralph", ConfigFileName)
}

// WorkspacePath returns the path to the workspace config file in the given
// current working directory (project-level config). Uses ConfigFileName in cwd.
func WorkspacePath(cwd string) string {
	return filepath.Join(cwd, ConfigFileName)
}
