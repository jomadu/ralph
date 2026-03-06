package cmdparse

import (
	"fmt"
	"strings"
)

// Parse parses a command string using shell-style quoting rules.
// Returns argv slice (program + arguments) or error for invalid syntax.
func Parse(cmd string) ([]string, error) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return nil, fmt.Errorf("empty command string")
	}

	var argv []string
	var current strings.Builder
	var inDoubleQuote, inSingleQuote bool
	var escaped bool

	for i := 0; i < len(cmd); i++ {
		ch := cmd[i]

		if escaped {
			if inDoubleQuote && (ch == '"' || ch == '\\') {
				current.WriteByte(ch)
			} else {
				current.WriteByte(ch)
			}
			escaped = false
			continue
		}

		if ch == '\\' && inDoubleQuote {
			escaped = true
			continue
		}

		if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}

		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}

		if (ch == ' ' || ch == '\t') && !inDoubleQuote && !inSingleQuote {
			if current.Len() > 0 {
				argv = append(argv, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteByte(ch)
	}

	if inDoubleQuote {
		return nil, fmt.Errorf("unclosed double quote")
	}
	if inSingleQuote {
		return nil, fmt.Errorf("unclosed single quote")
	}

	if current.Len() > 0 {
		argv = append(argv, current.String())
	}

	if len(argv) == 0 {
		return nil, fmt.Errorf("empty command string")
	}

	return argv, nil
}
