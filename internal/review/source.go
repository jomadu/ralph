// Package review implements the review component: prompt source resolution,
// report format, and exit code derivation. See docs/engineering/components/review.md.
package review

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrExit2 is the sentinel for "review or apply did not complete" (exit code 2).
// Callers should check with errors.Is and exit 2 when true.
var ErrExit2 = errors.New("review invocation error (exit 2)")

// ErrMissingSource indicates no prompt source was supplied (alias, file, or stdin).
var ErrMissingSource = fmt.Errorf("%w: exactly one of alias, file path, or stdin required", ErrExit2)

// ErrInvalidAlias indicates the prompt alias is not defined in config.
var ErrInvalidAlias = fmt.Errorf("%w: prompt alias not defined", ErrExit2)

// ErrAliasSourceMissing indicates the alias resolved to a path that does not exist or is unreadable.
var ErrAliasSourceMissing = fmt.Errorf("%w: prompt source file not found or unreadable", ErrExit2)

// ErrFileNotFound indicates the supplied file path does not exist or is not readable.
var ErrFileNotFound = fmt.Errorf("%w: prompt file not found or unreadable", ErrExit2)

// ErrMultipleSources indicates more than one prompt source was supplied.
var ErrMultipleSources = fmt.Errorf("%w: exactly one of alias, file path, or stdin allowed", ErrExit2)

// ErrApplyConfirmationRequired indicates apply would overwrite but confirmation is required and session is non-interactive (use --yes).
var ErrApplyConfirmationRequired = fmt.Errorf("%w: confirmation required to overwrite in non-interactive mode; use --yes to apply without confirmation", ErrExit2)

// PromptProvider supplies prompt definitions by name (for alias resolution).
// Implemented by effective config; used by review to resolve alias to path or content.
type PromptProvider interface {
	// PromptByName returns the path (file) or content (inline) for the named prompt.
	// If the prompt is file-based, path is set and content may be empty.
	// If the prompt is inline, content is set and path is empty.
	// ok is false if the name is not defined.
	PromptByName(name string) (path, content string, ok bool)
}

// ResolveOptions specifies exactly one prompt source for a review invocation.
// Exactly one of Alias, FilePath, or Stdin must be set; otherwise resolution returns an error.
type ResolveOptions struct {
	// Alias is the config prompt name (from prompts map). Mutually exclusive with FilePath and Stdin.
	Alias string
	// FilePath is the path to a prompt file. Mutually exclusive with Alias and Stdin.
	FilePath string
	// Stdin is the prompt content from standard input. Mutually exclusive with Alias and FilePath.
	Stdin []byte
}

// ResolvePromptSource resolves the single prompt source to prompt content.
// Exactly one of opts.Alias, opts.FilePath, or opts.Stdin must be set.
// cwd is used to resolve relative paths for alias-resolved paths and FilePath.
// Returns content and nil, or nil and an error that may be ErrExit2 (exit 2) or another error.
func ResolvePromptSource(provider PromptProvider, cwd string, opts ResolveOptions) ([]byte, error) {
	hasAlias := opts.Alias != ""
	hasFile := opts.FilePath != ""
	hasStdin := len(opts.Stdin) > 0
	n := 0
	if hasAlias {
		n++
	}
	if hasFile {
		n++
	}
	if hasStdin {
		n++
	}
	if n == 0 {
		return nil, ErrMissingSource
	}
	if n > 1 {
		return nil, ErrMultipleSources
	}

	if hasStdin {
		return opts.Stdin, nil
	}

	if hasFile {
		path := opts.FilePath
		if !filepath.IsAbs(path) {
			path = filepath.Join(cwd, path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("%w: %s", ErrFileNotFound, path)
			}
			return nil, fmt.Errorf("%w: %s: %v", ErrExit2, path, err)
		}
		return data, nil
	}

	// Alias
	path, content, ok := provider.PromptByName(opts.Alias)
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrInvalidAlias, opts.Alias)
	}
	if content != "" {
		return []byte(content), nil
	}
	if path == "" {
		return nil, fmt.Errorf("%w: alias %q has no path or content", ErrInvalidAlias, opts.Alias)
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s (alias %q)", ErrAliasSourceMissing, path, opts.Alias)
		}
		return nil, fmt.Errorf("%w: %s: %v", ErrExit2, path, err)
	}
	return data, nil
}

// IsExit2 reports whether err is or wraps ErrExit2 (caller should exit 2).
func IsExit2(err error) bool {
	return errors.Is(err, ErrExit2)
}
