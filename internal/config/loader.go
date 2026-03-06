package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
	return LoadConfigWithProvenanceAndExplicit("")
}

// LoadConfigWithProvenanceAndExplicit loads configuration with provenance tracking.
// If explicitPath is non-empty, it is used as the sole file-based config (global and workspace are skipped).
// If explicitPath is empty, global and workspace configs are loaded (missing files skipped silently).
// Returns the merged config with provenance or an error.
func LoadConfigWithProvenanceAndExplicit(explicitPath string) (ConfigWithProvenance, error) {
	cfg := DefaultConfigWithProvenance()

	// Load builtin aliases with default provenance
	for name, cmd := range BuiltinAliases() {
		cfg.AICmdAliases[name] = ValueWithProvenance[string]{Value: cmd, Provenance: ProvenanceDefault}
	}

	if explicitPath != "" {
		// Explicit config: load only this file, error if missing
		if err := overlayFileWithProvenance(explicitPath, &cfg, ProvenanceFile, false); err != nil {
			return cfg, err
		}
	} else {
		// Default config: load global and workspace, skip if missing
		globalPath := GlobalConfigPath()
		if globalPath != "" {
			if err := overlayFileWithProvenance(globalPath, &cfg, ProvenanceGlobal, true); err != nil {
				return cfg, err
			}
		}

		workspacePath := WorkspaceConfigPath()
		if err := overlayFileWithProvenance(workspacePath, &cfg, ProvenanceWorkspace, true); err != nil {
			return cfg, err
		}
	}

	// Overlay environment variables
	if err := overlayEnvironment(&cfg); err != nil {
		return cfg, err
	}

	// Validate resolved configuration
	if err := Validate(cfg); err != nil {
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

	// Validate resolved configuration
	if err := Validate(cfg); err != nil {
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

	// Parse into map to detect explicitly set fields
	var rawMap map[string]interface{}
	if err := yaml.Unmarshal(data, &rawMap); err != nil {
		return fmt.Errorf("failed to parse config %s: %w", path, err)
	}

	// Overlay loop config with explicit field detection
	if loopMap, ok := rawMap["loop"].(map[string]interface{}); ok {
		overlayLoopConfigWithMap(&cfg.Loop, &raw.Loop, loopMap, prov)
	}

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

// overlayLoopConfigWithMap overlays loop config values with provenance, using a map to detect explicit fields.
func overlayLoopConfigWithMap(dst *LoopConfigWithProvenance, src *LoopConfig, rawMap map[string]interface{}, prov Provenance) {
	if _, ok := rawMap["default_max_iterations"]; ok {
		dst.DefaultMaxIterations = ValueWithProvenance[int]{Value: src.DefaultMaxIterations, Provenance: prov}
	}
	if _, ok := rawMap["failure_threshold"]; ok {
		dst.FailureThreshold = ValueWithProvenance[int]{Value: src.FailureThreshold, Provenance: prov}
	}
	if _, ok := rawMap["iteration_timeout"]; ok {
		dst.IterationTimeout = ValueWithProvenance[int]{Value: src.IterationTimeout, Provenance: prov}
	}
	if _, ok := rawMap["max_output_buffer"]; ok {
		dst.MaxOutputBuffer = ValueWithProvenance[int]{Value: src.MaxOutputBuffer, Provenance: prov}
	}
	if _, ok := rawMap["show_ai_output"]; ok {
		dst.ShowAIOutput = ValueWithProvenance[bool]{Value: src.ShowAIOutput, Provenance: prov}
	}
	if _, ok := rawMap["preamble"]; ok {
		dst.Preamble = ValueWithProvenance[bool]{Value: src.Preamble, Provenance: prov}
	}
	if _, ok := rawMap["ai_cmd_alias"]; ok {
		dst.AICmdAlias = ValueWithProvenance[string]{Value: src.AICmdAlias, Provenance: prov}
	}
	if signalsMap, ok := rawMap["signals"].(map[string]interface{}); ok {
		if _, ok := signalsMap["success"]; ok {
			dst.SignalSuccess = ValueWithProvenance[string]{Value: src.Signals.Success, Provenance: prov}
		}
		if _, ok := signalsMap["failure"]; ok {
			dst.SignalFailure = ValueWithProvenance[string]{Value: src.Signals.Failure, Provenance: prov}
		}
	}
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
	// Preamble is bool, overlay always
	dst.Preamble = ValueWithProvenance[bool]{Value: src.Preamble, Provenance: prov}
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

// overlayEnvironment reads RALPH_* environment variables and overlays them onto config.
func overlayEnvironment(cfg *ConfigWithProvenance) error {
	// RALPH_LOOP_AI_CMD
	if v := os.Getenv("RALPH_LOOP_AI_CMD"); v != "" {
		// ai_cmd not yet in config struct; will be added when O3/R6 is implemented
	}

	// RALPH_LOOP_AI_CMD_ALIAS
	if v := os.Getenv("RALPH_LOOP_AI_CMD_ALIAS"); v != "" {
		cfg.Loop.AICmdAlias = ValueWithProvenance[string]{Value: v, Provenance: ProvenanceEnv}
	}

	// RALPH_LOOP_ITERATION_MODE
	if v := os.Getenv("RALPH_LOOP_ITERATION_MODE"); v != "" {
		// iteration_mode not yet in config struct; will be added when O1/R4 is implemented
	}

	// RALPH_LOOP_DEFAULT_MAX_ITERATIONS
	if v := os.Getenv("RALPH_LOOP_DEFAULT_MAX_ITERATIONS"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid RALPH_LOOP_DEFAULT_MAX_ITERATIONS: %q: %w", v, err)
		}
		cfg.Loop.DefaultMaxIterations = ValueWithProvenance[int]{Value: val, Provenance: ProvenanceEnv}
	}

	// RALPH_LOOP_FAILURE_THRESHOLD
	if v := os.Getenv("RALPH_LOOP_FAILURE_THRESHOLD"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid RALPH_LOOP_FAILURE_THRESHOLD: %q: %w", v, err)
		}
		cfg.Loop.FailureThreshold = ValueWithProvenance[int]{Value: val, Provenance: ProvenanceEnv}
	}

	// RALPH_LOOP_ITERATION_TIMEOUT
	if v := os.Getenv("RALPH_LOOP_ITERATION_TIMEOUT"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid RALPH_LOOP_ITERATION_TIMEOUT: %q: %w", v, err)
		}
		cfg.Loop.IterationTimeout = ValueWithProvenance[int]{Value: val, Provenance: ProvenanceEnv}
	}

	// RALPH_LOOP_LOG_LEVEL
	if v := os.Getenv("RALPH_LOOP_LOG_LEVEL"); v != "" {
		// log_level not yet in config struct; will be added when O4/R5 is implemented
	}

	// RALPH_LOOP_SHOW_AI_OUTPUT
	if v := os.Getenv("RALPH_LOOP_SHOW_AI_OUTPUT"); v != "" {
		val, err := parseBool(v)
		if err != nil {
			return fmt.Errorf("invalid RALPH_LOOP_SHOW_AI_OUTPUT: %q: %w", v, err)
		}
		cfg.Loop.ShowAIOutput = ValueWithProvenance[bool]{Value: val, Provenance: ProvenanceEnv}
	}

	// RALPH_LOOP_PREAMBLE
	if v := os.Getenv("RALPH_LOOP_PREAMBLE"); v != "" {
		val, err := parseBool(v)
		if err != nil {
			return fmt.Errorf("invalid RALPH_LOOP_PREAMBLE: %q: %w", v, err)
		}
		cfg.Loop.Preamble = ValueWithProvenance[bool]{Value: val, Provenance: ProvenanceEnv}
	}

	return nil
}

// parseBool parses common boolean representations.
func parseBool(s string) (bool, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off", "":
		return false, nil
	default:
		return false, fmt.Errorf("not a boolean value")
	}
}

// CLIFlags holds CLI flag values for config overlay.
type CLIFlags struct {
	MaxIterations    *int
	FailureThreshold *int
	IterationTimeout *int
	MaxOutputBuffer  *int
	Preamble         *bool
	AICmdAlias       *string
	SignalSuccess    *string
	SignalFailure    *string
	ShowAIOutput     *bool
}

// OverlayCLIFlags applies CLI flag values to config with ProvenanceCLI.
// Only non-nil flag values are applied.
func OverlayCLIFlags(cfg *ConfigWithProvenance, flags CLIFlags) {
	if flags.MaxIterations != nil {
		cfg.Loop.DefaultMaxIterations = ValueWithProvenance[int]{Value: *flags.MaxIterations, Provenance: ProvenanceCLI}
	}
	if flags.FailureThreshold != nil {
		cfg.Loop.FailureThreshold = ValueWithProvenance[int]{Value: *flags.FailureThreshold, Provenance: ProvenanceCLI}
	}
	if flags.IterationTimeout != nil {
		cfg.Loop.IterationTimeout = ValueWithProvenance[int]{Value: *flags.IterationTimeout, Provenance: ProvenanceCLI}
	}
	if flags.MaxOutputBuffer != nil {
		cfg.Loop.MaxOutputBuffer = ValueWithProvenance[int]{Value: *flags.MaxOutputBuffer, Provenance: ProvenanceCLI}
	}
	if flags.Preamble != nil {
		cfg.Loop.Preamble = ValueWithProvenance[bool]{Value: *flags.Preamble, Provenance: ProvenanceCLI}
	}
	if flags.AICmdAlias != nil {
		cfg.Loop.AICmdAlias = ValueWithProvenance[string]{Value: *flags.AICmdAlias, Provenance: ProvenanceCLI}
	}
	if flags.SignalSuccess != nil {
		cfg.Loop.SignalSuccess = ValueWithProvenance[string]{Value: *flags.SignalSuccess, Provenance: ProvenanceCLI}
	}
	if flags.SignalFailure != nil {
		cfg.Loop.SignalFailure = ValueWithProvenance[string]{Value: *flags.SignalFailure, Provenance: ProvenanceCLI}
	}
	if flags.ShowAIOutput != nil {
		cfg.Loop.ShowAIOutput = ValueWithProvenance[bool]{Value: *flags.ShowAIOutput, Provenance: ProvenanceCLI}
	}
}
