package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// CommandResolution holds the result of AI command resolution.
type CommandResolution struct {
	Command string
	Source  string // "direct command" or "alias <name>"
}

// ResolveAICommand resolves the AI command from effective config.
// Direct command (AICmd) takes precedence over alias (AICmdAlias).
// Returns error if neither is configured.
func ResolveAICommand(cfg ConfigWithProvenance) (CommandResolution, error) {
	// Step 1: Check for direct command
	if cfg.Loop.AICmd.Value != "" {
		return CommandResolution{
			Command: cfg.Loop.AICmd.Value,
			Source:  "direct command",
		}, nil
	}

	// Step 2: Check for alias
	if cfg.Loop.AICmdAlias.Value != "" {
		aliasName := cfg.Loop.AICmdAlias.Value
		
		// Look up alias in merged alias map
		if aliasValue, ok := cfg.AICmdAliases[aliasName]; ok {
			return CommandResolution{
				Command: aliasValue.Value,
				Source:  fmt.Sprintf("alias %s", aliasName),
			}, nil
		}
		
		// Alias not found - build error with available aliases
		plainCfg := Config{AICmdAliases: make(map[string]string)}
		for k, v := range cfg.AICmdAliases {
			plainCfg.AICmdAliases[k] = v.Value
		}
		merged := MergedAliases(plainCfg)
		var available []string
		for k := range merged {
			available = append(available, k)
		}
		return CommandResolution{}, fmt.Errorf("unknown AI command alias %q; available aliases: %v", aliasName, available)
	}

	// Step 3: No command configured
	return CommandResolution{}, fmt.Errorf("no AI command configured: set loop.ai_cmd or loop.ai_cmd_alias via config, RALPH_LOOP_AI_CMD or RALPH_LOOP_AI_CMD_ALIAS via environment, or --ai-cmd or --ai-cmd-alias via CLI")
}


// Provenance identifies the source layer of a config value.
type Provenance string

const (
	ProvenanceDefault   Provenance = "default"
	ProvenanceGlobal    Provenance = "global"
	ProvenanceWorkspace Provenance = "workspace"
	ProvenanceFile      Provenance = "file"
	ProvenancePrompt    Provenance = "prompt"
	ProvenanceEnv       Provenance = "env"
	ProvenanceCLI       Provenance = "cli"
)

// LoopConfig holds loop execution settings.
type LoopConfig struct {
	DefaultMaxIterations int    `yaml:"default_max_iterations"`
	IterationMode        string `yaml:"iteration_mode"`
	FailureThreshold     int    `yaml:"failure_threshold"`
	IterationTimeout     int    `yaml:"iteration_timeout"`
	MaxOutputBuffer      int    `yaml:"max_output_buffer"`
	ShowAIOutput         bool   `yaml:"show_ai_output"`
	LogLevel             string `yaml:"log_level"`
	Preamble             bool   `yaml:"preamble"`
	AICmd                string `yaml:"ai_cmd"`
	AICmdAlias           string `yaml:"ai_cmd_alias"`
	Signals              struct {
		Success string `yaml:"success"`
		Failure string `yaml:"failure"`
	} `yaml:"signals"`
}

// LoopConfigWithProvenance holds loop config with provenance metadata.
type LoopConfigWithProvenance struct {
	DefaultMaxIterations ValueWithProvenance[int]
	IterationMode        ValueWithProvenance[string]
	FailureThreshold     ValueWithProvenance[int]
	IterationTimeout     ValueWithProvenance[int]
	MaxOutputBuffer      ValueWithProvenance[int]
	ShowAIOutput         ValueWithProvenance[bool]
	LogLevel             ValueWithProvenance[string]
	Preamble             ValueWithProvenance[bool]
	AICmd                ValueWithProvenance[string]
	AICmdAlias           ValueWithProvenance[string]
	SignalSuccess        ValueWithProvenance[string]
	SignalFailure        ValueWithProvenance[string]
}

// ValueWithProvenance wraps a value with its source layer.
type ValueWithProvenance[T any] struct {
	Value      T
	Provenance Provenance
}

// PromptConfig defines a prompt alias with optional loop overrides.
type PromptConfig struct {
	Path        string      `yaml:"path"`
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Loop        *LoopConfig `yaml:"loop,omitempty"`
}

// PromptConfigWithProvenance tracks provenance for prompt configuration.
type PromptConfigWithProvenance struct {
	Path        ValueWithProvenance[string]
	Name        ValueWithProvenance[string]
	Description ValueWithProvenance[string]
	Loop        *LoopConfig // Raw loop config from YAML
	LoopRawMap  map[string]interface{} // Raw map for explicit field detection
	Provenance  Provenance // Provenance of this prompt definition
}

// Config is the root configuration structure.
type Config struct {
	Loop          LoopConfig              `yaml:"loop"`
	Prompts       map[string]PromptConfig `yaml:"prompts"`
	AICmdAliases  map[string]string       `yaml:"ai_cmd_aliases"`
}

// ConfigWithProvenance is the resolved configuration with provenance metadata.
type ConfigWithProvenance struct {
	Loop         LoopConfigWithProvenance
	Prompts      map[string]PromptConfigWithProvenance
	AICmdAliases map[string]ValueWithProvenance[string]
}

// DefaultConfig returns a Config with built-in defaults.
func DefaultConfig() Config {
	return Config{
		Loop: LoopConfig{
			DefaultMaxIterations: 5,
			IterationMode:        "max-iterations",
			FailureThreshold:     3,
			IterationTimeout:     300,
			MaxOutputBuffer:      10485760, // 10 MB
			ShowAIOutput:         false,
			LogLevel:             "info",
			Preamble:             true,
			AICmd:                "",
			AICmdAlias:           "",
			Signals: struct {
				Success string `yaml:"success"`
				Failure string `yaml:"failure"`
			}{
				Success: "<promise>SUCCESS</promise>",
				Failure: "<promise>FAILURE</promise>",
			},
		},
		Prompts:      make(map[string]PromptConfig),
		AICmdAliases: BuiltinAliases(),
	}
}

// DefaultConfigWithProvenance returns a ConfigWithProvenance with built-in defaults tagged.
func DefaultConfigWithProvenance() ConfigWithProvenance {
	return ConfigWithProvenance{
		Loop: LoopConfigWithProvenance{
			DefaultMaxIterations: ValueWithProvenance[int]{Value: 5, Provenance: ProvenanceDefault},
			IterationMode:        ValueWithProvenance[string]{Value: "max-iterations", Provenance: ProvenanceDefault},
			FailureThreshold:     ValueWithProvenance[int]{Value: 3, Provenance: ProvenanceDefault},
			IterationTimeout:     ValueWithProvenance[int]{Value: 300, Provenance: ProvenanceDefault},
			MaxOutputBuffer:      ValueWithProvenance[int]{Value: 10485760, Provenance: ProvenanceDefault},
			ShowAIOutput:         ValueWithProvenance[bool]{Value: false, Provenance: ProvenanceDefault},
			LogLevel:             ValueWithProvenance[string]{Value: "info", Provenance: ProvenanceDefault},
			Preamble:             ValueWithProvenance[bool]{Value: true, Provenance: ProvenanceDefault},
			AICmd:                ValueWithProvenance[string]{Value: "", Provenance: ProvenanceDefault},
			AICmdAlias:           ValueWithProvenance[string]{Value: "", Provenance: ProvenanceDefault},
			SignalSuccess:        ValueWithProvenance[string]{Value: "<promise>SUCCESS</promise>", Provenance: ProvenanceDefault},
			SignalFailure:        ValueWithProvenance[string]{Value: "<promise>FAILURE</promise>", Provenance: ProvenanceDefault},
		},
		Prompts:      make(map[string]PromptConfigWithProvenance),
		AICmdAliases: make(map[string]ValueWithProvenance[string]),
	}
}

// GlobalConfigPath resolves the global config file path using the fallback chain:
// RALPH_CONFIG_HOME → XDG_CONFIG_HOME/ralph → ~/.config/ralph
func GlobalConfigPath() string {
	if dir := os.Getenv("RALPH_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, "ralph-config.yml")
	}
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, "ralph", "ralph-config.yml")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "ralph", "ralph-config.yml")
}

// WorkspaceConfigPath returns the workspace config file path relative to cwd.
func WorkspaceConfigPath() string {
	return "ralph-config.yml"
}
