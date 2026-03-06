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

// LoadConfigWithProvenance loads configuration with provenance tracking.
// Missing default config files are skipped silently.
// Returns the merged config with provenance or an error.
func LoadConfigWithProvenance() (ConfigWithProvenance, error) {
	cfg := DefaultConfigWithProvenance()

	// Load builtin aliases with default provenance
	for name, cmd := range BuiltinAliases() {
		cfg.AICmdAliases[name] = ValueWithProvenance[string]{Value: cmd, Provenance: ProvenanceDefault}
	}

	// Load global config if it exists
	globalPath := GlobalConfigPath()
	if globalPath != "" {
		if err := overlayFileWithProvenance(globalPath, &cfg, ProvenanceGlobal, true); err != nil {
			return cfg, err
		}
	}

	// Load workspace config if it exists
	workspacePath := WorkspaceConfigPath()
	if err := overlayFileWithProvenance(workspacePath, &cfg, ProvenanceWorkspace, true); err != nil {
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

// LoadConfigFromFileWithProvenance loads configuration from an explicit file with provenance.
// Returns an error if the file does not exist or cannot be read/parsed.
func LoadConfigFromFileWithProvenance(path string) (ConfigWithProvenance, error) {
	cfg := DefaultConfigWithProvenance()

	// Load builtin aliases with default provenance
	for name, cmd := range BuiltinAliases() {
		cfg.AICmdAliases[name] = ValueWithProvenance[string]{Value: cmd, Provenance: ProvenanceDefault}
	}

	if err := overlayFileWithProvenance(path, &cfg, ProvenanceFile, false); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// overlayFileWithProvenance loads a YAML config file and overlays values with provenance.
func overlayFileWithProvenance(path string, cfg *ConfigWithProvenance, prov Provenance, silent bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if silent {
				return nil
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

	var raw Config
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to parse config %s: %w", path, err)
	}

	// Overlay loop config
	overlayLoopConfig(&cfg.Loop, &raw.Loop, prov)

	// Overlay prompts (no provenance for prompts yet)
	for k, v := range raw.Prompts {
		cfg.Prompts[k] = v
	}

	// Overlay AI command aliases
	for k, v := range raw.AICmdAliases {
		cfg.AICmdAliases[k] = ValueWithProvenance[string]{Value: v, Provenance: prov}
	}

	return nil
}

// overlayLoopConfig overlays loop config values with provenance.
func overlayLoopConfig(dst *LoopConfigWithProvenance, src *LoopConfig, prov Provenance) {
	if src.DefaultMaxIterations != 0 {
		dst.DefaultMaxIterations = ValueWithProvenance[int]{Value: src.DefaultMaxIterations, Provenance: prov}
	}
	if src.FailureThreshold != 0 {
		dst.FailureThreshold = ValueWithProvenance[int]{Value: src.FailureThreshold, Provenance: prov}
	}
	if src.IterationTimeout != 0 {
		dst.IterationTimeout = ValueWithProvenance[int]{Value: src.IterationTimeout, Provenance: prov}
	}
	if src.MaxOutputBuffer != 0 {
		dst.MaxOutputBuffer = ValueWithProvenance[int]{Value: src.MaxOutputBuffer, Provenance: prov}
	}
	// ShowAIOutput is bool, check if explicitly set (requires more sophisticated detection; for now overlay always)
	dst.ShowAIOutput = ValueWithProvenance[bool]{Value: src.ShowAIOutput, Provenance: prov}
	if src.AICmdAlias != "" {
		dst.AICmdAlias = ValueWithProvenance[string]{Value: src.AICmdAlias, Provenance: prov}
	}
	if src.Signals.Success != "" {
		dst.SignalSuccess = ValueWithProvenance[string]{Value: src.Signals.Success, Provenance: prov}
	}
	if src.Signals.Failure != "" {
		dst.SignalFailure = ValueWithProvenance[string]{Value: src.Signals.Failure, Provenance: prov}
	}
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
