package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads configuration from default locations (global and workspace).
// Missing default config files are skipped silently.
// Returns the merged config or an error if a file exists but cannot be read/parsed.
func LoadConfig() (Config, error) {
	cfg := DefaultConfig()

	// Load global config if it exists
	globalPath := GlobalConfigPath()
	if globalPath != "" {
		if err := loadFile(globalPath, &cfg, true); err != nil {
			return cfg, err
		}
	}

	// Load workspace config if it exists
	workspacePath := WorkspaceConfigPath()
	if err := loadFile(workspacePath, &cfg, true); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// LoadConfigFromFile loads configuration from an explicit file path.
// Returns an error if the file does not exist or cannot be read/parsed.
func LoadConfigFromFile(path string) (Config, error) {
	cfg := DefaultConfig()
	if err := loadFile(path, &cfg, false); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// loadFile loads a YAML config file into cfg.
// If silent is true, missing files are skipped without error.
// If silent is false, missing files produce an error.
func loadFile(path string, cfg *Config, silent bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if silent {
				return nil // Silent skip
			}
			return fmt.Errorf("config file not found: %s", path)
		}
		if os.IsPermission(err) {
			if silent {
				fmt.Fprintf(os.Stderr, "warning: cannot read config %s: permission denied\n", path)
				return nil
			}
			return fmt.Errorf("config file not readable: %s: permission denied", path)
		}
		return fmt.Errorf("failed to read config %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse config %s: %w", path, err)
	}

	return nil
}
