package config

// LoopConfig holds loop execution settings.
type LoopConfig struct {
	DefaultMaxIterations int    `yaml:"default_max_iterations"`
	FailureThreshold     int    `yaml:"failure_threshold"`
	IterationTimeout     int    `yaml:"iteration_timeout"`
	MaxOutputBuffer      int    `yaml:"max_output_buffer"`
	ShowAIOutput         bool   `yaml:"show_ai_output"`
	AICmdAlias           string `yaml:"ai_cmd_alias"`
	Signals              struct {
		Success string `yaml:"success"`
		Failure string `yaml:"failure"`
	} `yaml:"signals"`
}

// PromptConfig defines a prompt alias with optional loop overrides.
type PromptConfig struct {
	Path        string      `yaml:"path"`
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Loop        *LoopConfig `yaml:"loop,omitempty"`
}

// Config is the root configuration structure.
type Config struct {
	Loop          LoopConfig              `yaml:"loop"`
	Prompts       map[string]PromptConfig `yaml:"prompts"`
	AICmdAliases  map[string]string       `yaml:"ai_cmd_aliases"`
}

// DefaultConfig returns a Config with built-in defaults.
func DefaultConfig() Config {
	return Config{
		Loop: LoopConfig{
			DefaultMaxIterations: 5,
			FailureThreshold:     3,
			IterationTimeout:     300,
			MaxOutputBuffer:      10485760, // 10 MB
			ShowAIOutput:         false,
			AICmdAlias:           "claude",
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
