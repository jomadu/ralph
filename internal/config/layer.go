package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ErrExplicitConfigMissing is returned when the user-supplied config file path
// does not exist or is not a regular file (O002/R005).
var ErrExplicitConfigMissing = errors.New("config file missing or unreadable")

// FileLayer holds the parsed content of a single config file (global, workspace, or explicit).
// It mirrors the canonical config file structure; merging and defaults are applied later.
type FileLayer struct {
	Loop    *LoopSection      `yaml:"loop,omitempty"`
	Prompts map[string]Prompt `yaml:"prompts,omitempty"`
	Aliases map[string]Alias  `yaml:"aliases,omitempty"`
}

// LoopSection holds loop behavior settings from a config file.
type LoopSection struct {
	MaxIterations    *int        `yaml:"max_iterations,omitempty"`
	FailureThreshold *int        `yaml:"failure_threshold,omitempty"`
	TimeoutSeconds   *int        `yaml:"timeout_seconds,omitempty"`
	SuccessSignal    string      `yaml:"success_signal,omitempty"`
	FailureSignal    string      `yaml:"failure_signal,omitempty"`
	Preamble         interface{} `yaml:"preamble,omitempty"` // string or bool
	Streaming        *bool       `yaml:"streaming,omitempty"`
	LogLevel         string      `yaml:"log_level,omitempty"`
	MaxOutputBuffer  *int        `yaml:"max_output_buffer,omitempty"`
	AiCmd            string      `yaml:"ai_cmd,omitempty"`
	AiCmdAlias       string      `yaml:"ai_cmd_alias,omitempty"`
}

// Prompt holds a single prompt definition (path or content, optional loop overrides).
// DisplayName and Description are optional and used for list output (R006).
type Prompt struct {
	Path        string       `yaml:"path,omitempty"`
	Content     string       `yaml:"content,omitempty"`
	DisplayName string       `yaml:"display_name,omitempty"`
	Description string       `yaml:"description,omitempty"`
	Loop        *LoopSection `yaml:"loop,omitempty"`
}

// LoopSettings holds effective loop behavior (concrete values). Used by run-loop
// and resolution; built-in defaults ensure the tool works without a config file (O002/R001, R002).
type LoopSettings struct {
	MaxIterations    int
	FailureThreshold int
	TimeoutSeconds   int
	SuccessSignal    string
	FailureSignal    string
	Preamble         string // empty = no preamble injection
	Streaming        bool
	LogLevel         string
	MaxOutputBuffer  int // bytes; 0 = unlimited (for backward compat); default from DefaultLoopSettings is 65536
	AICmd            string
	AICmdAlias       string
}

// Alias holds an AI command alias. In YAML, an alias value may be a string
// (the command) or an object with a "command" key.
type Alias struct {
	Command string `yaml:"command"`
}

// UnmarshalYAML allows alias to be a string (command) or object { command: "..." }.
func (a *Alias) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		a.Command = value.Value
		return nil
	}
	var raw struct {
		Command string `yaml:"command"`
	}
	if err := value.Decode(&raw); err != nil {
		return err
	}
	a.Command = raw.Command
	return nil
}

// ReadLayer reads and parses the config file at path. If the file is missing,
// returns (nil, nil) — skip without error. Returns an error only for read or
// parse failures.
func ReadLayer(path string) (*FileLayer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return ParseLayer(data)
}

// ReadLayerRequired reads and parses the config file at path when the path is
// explicitly supplied by the user (e.g. --config). The file must exist and be
// readable; if it is missing, a directory, or unreadable, returns an error
// (O002/R005). Path may be relative to cwd or absolute.
func ReadLayerRequired(path string) (*FileLayer, error) {
	path = filepath.Clean(path)
	if path == "" || path == "." {
		return nil, fmt.Errorf("%w: path is empty or invalid", ErrExplicitConfigMissing)
	}
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrExplicitConfigMissing, path)
		}
		return nil, fmt.Errorf("config file %s: %w", path, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("%w: %s is a directory", ErrExplicitConfigMissing, path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config file %s: %w", path, err)
	}
	return ParseLayer(data)
}

// ParseLayer parses YAML bytes into a FileLayer and validates against the
// canonical schema. If data is nil or empty, returns (nil, nil). Invalid or
// out-of-range values produce a clear error (config component Config file structure).
func ParseLayer(data []byte) (*FileLayer, error) {
	if len(data) == 0 {
		return nil, nil
	}
	var layer FileLayer
	if err := yaml.Unmarshal(data, &layer); err != nil {
		return nil, err
	}
	if err := ValidateFileLayer(&layer); err != nil {
		return nil, err
	}
	return &layer, nil
}
