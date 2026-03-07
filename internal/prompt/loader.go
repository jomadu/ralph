package prompt

import (
	"fmt"
	"io"
	"os"

	"github.com/maxdunn/ralph/internal/config"
)

// Mode represents the prompt input mode.
type Mode int

const (
	ModeAlias Mode = iota
	ModeFile
	ModeStdin
)

// Source holds the resolved prompt source and content.
type Source struct {
	Mode    Mode
	Alias   string // Set for ModeAlias
	Path    string // Set for ModeAlias and ModeFile
	Content []byte // Prompt content, read once
}

// ResolveMode determines the prompt input mode from CLI arguments and stdin state.
// Resolution order (O5 R1, shared by run and review): (1) if -f present use file,
// (2) else positional alias, (3) else stdin. When both alias and -f are present,
// file wins and alias is ignored.
// Returns an error only when no source is identified (no alias, no file, no piped stdin).
func ResolveMode(alias string, filePath string) (Mode, error) {
	hasAlias := alias != ""
	hasFile := filePath != ""

	// Mode resolution order: file > alias > stdin (O5 R1)
	if hasFile {
		return ModeFile, nil
	}
	if hasAlias {
		return ModeAlias, nil
	}

	// Check if stdin is piped
	stat, err := os.Stdin.Stat()
	if err != nil {
		return 0, fmt.Errorf("failed to check stdin: %w", err)
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// stdin is not a TTY (piped input)
		return ModeStdin, nil
	}

	// No source identified
	return 0, fmt.Errorf("no prompt source: provide an alias, --file, or pipe input via stdin")
}

// LoadPrompt loads the prompt content based on the resolved mode.
// Validates the source and reads the content once into memory.
func LoadPrompt(mode Mode, alias string, filePath string, cfg *config.ConfigWithProvenance) (*Source, error) {
	switch mode {
	case ModeAlias:
		return loadFromAlias(alias, cfg)
	case ModeFile:
		return loadFromFile(filePath)
	case ModeStdin:
		return loadFromStdin()
	default:
		return nil, fmt.Errorf("unknown prompt mode: %d", mode)
	}
}

// loadFromAlias loads prompt from a configured alias.
func loadFromAlias(alias string, cfg *config.ConfigWithProvenance) (*Source, error) {
	// Look up alias in config
	promptCfg, ok := cfg.Prompts[alias]
	if !ok {
		return nil, fmt.Errorf("unknown prompt alias %q", alias)
	}

	path := promptCfg.Path.Value
	if path == "" {
		return nil, fmt.Errorf("prompt alias %q has no path configured", alias)
	}

	// Read and validate file
	content, err := readFile(path)
	if err != nil {
		return nil, fmt.Errorf("prompt file error (alias: %s): %w", alias, err)
	}

	if len(content) == 0 {
		return nil, fmt.Errorf("prompt file is empty: %s (alias: %s)", path, alias)
	}

	return &Source{
		Mode:    ModeAlias,
		Alias:   alias,
		Path:    path,
		Content: content,
	}, nil
}

// loadFromFile loads prompt from a file path.
func loadFromFile(path string) (*Source, error) {
	content, err := readFile(path)
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return nil, fmt.Errorf("prompt file is empty: %s", path)
	}

	return &Source{
		Mode:    ModeFile,
		Path:    path,
		Content: content,
	}, nil
}

// loadFromStdin loads prompt from stdin.
func loadFromStdin() (*Source, error) {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read stdin: %w", err)
	}

	if len(content) == 0 {
		return nil, fmt.Errorf("stdin is empty: no prompt content provided")
	}

	return &Source{
		Mode:    ModeStdin,
		Content: content,
	}, nil
}

// readFile reads a file and returns its content with validation.
func readFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("prompt file not found: %s", path)
		}
		if os.IsPermission(err) {
			return nil, fmt.Errorf("prompt file not readable: %s: permission denied", path)
		}
		return nil, fmt.Errorf("failed to read prompt file %s: %w", path, err)
	}
	return content, nil
}
