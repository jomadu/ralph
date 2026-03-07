// Package review provides prompt review behavior (O5). Path resolution for report output (R3).
package review

import (
	"fmt"
	"os"
	"path/filepath"
)

// ResolveReportPath returns the path where the review report will be written.
// If explicitPath is non-empty, it is resolved (relative to CWD) and validated:
// parent must exist and be writable; path must be a file or non-existent (not an existing directory).
// If explicitPath is empty, a unique path in the system temp directory is chosen and validated
// (temp dir must exist and be writable). Returns (path, isTemp, nil) or ("", false, error).
// Call before spawning the review-phase AI; on error the caller should exit 2 (R8).
func ResolveReportPath(explicitPath string) (reportPath string, isTemp bool, err error) {
	if explicitPath != "" {
		path, err := resolveExplicitPath(explicitPath)
		return path, false, err
	}
	path, err := resolveTempPath()
	return path, true, err
}

// reportPathErr formats an error for report path validation (R8: "cannot write report to <path>: <reason>").
func reportPathErr(path string, reason string) error {
	return fmt.Errorf("cannot write report to %s: %s", path, reason)
}

func resolveExplicitPath(explicitPath string) (string, error) {
	path, err := filepath.Abs(explicitPath)
	if err != nil {
		return "", reportPathErr(explicitPath, err.Error())
	}

	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return "", reportPathErr(path, "path is a directory (must be a file)")
		}
		// Existing file: overwrite allowed per R3
		return path, nil
	}
	if !os.IsNotExist(err) {
		return "", reportPathErr(path, err.Error())
	}

	// Path does not exist; ensure parent exists and is writable
	parent := filepath.Dir(path)
	parentInfo, err := os.Stat(parent)
	if err != nil {
		if os.IsNotExist(err) {
			return "", reportPathErr(path, "parent directory does not exist: "+parent)
		}
		return "", reportPathErr(path, err.Error())
	}
	if !parentInfo.IsDir() {
		return "", reportPathErr(path, "parent is not a directory: "+parent)
	}

	if err := checkWritable(parent); err != nil {
		return "", reportPathErr(path, err.Error())
	}
	return path, nil
}

func resolveTempPath() (string, error) {
	dir := os.TempDir()
	if dir == "" {
		return "", fmt.Errorf("cannot write report: system temporary directory not available")
	}
	dirInfo, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("cannot write report: temp directory unavailable: %w", err)
	}
	if !dirInfo.IsDir() {
		return "", fmt.Errorf("cannot write report: temp path is not a directory: %s", dir)
	}
	if err := checkWritable(dir); err != nil {
		return "", fmt.Errorf("cannot write report: temp directory not writable: %w", err)
	}

	// Create a unique path: create then remove so we get a path the OS guarantees is unique
	f, err := os.CreateTemp(dir, "ralph-review-*.md")
	if err != nil {
		return "", fmt.Errorf("cannot write report: temp directory not writable: %w", err)
	}
	path := f.Name()
	_ = f.Close()
	if err := os.Remove(path); err != nil {
		return "", fmt.Errorf("cannot write report: failed to reserve temp path: %w", err)
	}
	return path, nil
}

// checkWritable verifies that dir is writable by creating and removing a temp file in it.
func checkWritable(dir string) error {
	f, err := os.CreateTemp(dir, ".ralph-write-check-")
	if err != nil {
		return err
	}
	name := f.Name()
	_ = f.Close()
	return os.Remove(name)
}

// InputMode represents the prompt input mode for review (alias, file, or stdin).
type InputMode int

const (
	InputModeAlias InputMode = iota
	InputModeFile
	InputModeStdin
)

// ResolvePromptOutputPathResult holds the resolved prompt output path and whether it was explicitly set.
// Used for revision-phase interpolation (R4, R5) and for writing the suggested revision when not applying.
type ResolvePromptOutputPathResult struct {
	Path     string // Resolved absolute path where revision is written; empty if no path needed
	Explicit bool   // True if user set --prompt-output
	NeedPath bool   // True when revision will be written (apply or non-apply with --prompt-output)
}

// ResolvePromptOutputPath resolves where the suggested revised prompt is written (R4).
// - Without apply: if promptOutputFlag is set, that path is used; otherwise no path need be interpolated for writing.
// - With apply + alias/file: default is sourcePath; promptOutputFlag can override to a different path.
// - With apply + stdin: promptOutputFlag is required; absence returns ErrStdinApplyRequiresPromptOutput.
// When a path is determined and will be written to (NeedPath true), the path is validated for writability; invalid returns error (exit 2 per R8).
func ResolvePromptOutputPath(inputMode InputMode, applyRequested bool, promptOutputFlag string, sourcePath string) (*ResolvePromptOutputPathResult, error) {
	out := &ResolvePromptOutputPathResult{}

	// Stdin + apply without --prompt-output is invalid (R4, R5, R8)
	if inputMode == InputModeStdin && applyRequested && promptOutputFlag == "" {
		return nil, ErrStdinApplyRequiresPromptOutput
	}

	needPath := false
	var path string

	if applyRequested {
		if inputMode == InputModeStdin {
			path = promptOutputFlag
			needPath = true
		} else {
			// Alias or file: default to source; --prompt-output overrides
			if promptOutputFlag != "" {
				path = promptOutputFlag
			} else {
				path = sourcePath
			}
			needPath = true
		}
		out.Explicit = promptOutputFlag != ""
	} else {
		// No apply: only need path if user set --prompt-output
		if promptOutputFlag != "" {
			path = promptOutputFlag
			needPath = true
			out.Explicit = true
		}
	}

	out.NeedPath = needPath
	if !needPath {
		return out, nil
	}

	resolved, err := resolvePromptOutputPathWritable(path)
	if err != nil {
		return nil, err
	}
	out.Path = resolved
	return out, nil
}

// ErrStdinApplyRequiresPromptOutput is returned when stdin + --apply is used without --prompt-output (R4, R8).
var ErrStdinApplyRequiresPromptOutput = fmt.Errorf("stdin input with --apply requires --prompt-output <path>")

func resolvePromptOutputPathWritable(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("prompt output path invalid: %w", err)
	}

	info, err := os.Stat(abs)
	if err == nil {
		if info.IsDir() {
			return "", fmt.Errorf("prompt output path is a directory (must be a file): %s", abs)
		}
		// Existing file: must be writable
		if err := checkFileWritable(abs); err != nil {
			return "", fmt.Errorf("prompt output path unwritable: %w", err)
		}
		return abs, nil
	}
	if !os.IsNotExist(err) {
		return "", fmt.Errorf("prompt output path: %w", err)
	}

	parent := filepath.Dir(abs)
	parentInfo, err := os.Stat(parent)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("prompt output path invalid: parent directory does not exist: %s", parent)
		}
		return "", fmt.Errorf("prompt output path: %w", err)
	}
	if !parentInfo.IsDir() {
		return "", fmt.Errorf("prompt output path invalid: parent is not a directory: %s", parent)
	}
	if err := checkWritable(parent); err != nil {
		return "", fmt.Errorf("prompt output path unwritable: %w", err)
	}
	return abs, nil
}

func checkFileWritable(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	return f.Close()
}

// VerifyReportExists checks that a regular file exists at reportPath after the review-phase AI exits (R9).
// If the file is missing, a directory, or a symlink to a missing file, returns an error (caller should exit 2 per R8).
func VerifyReportExists(reportPath string) error {
	return verifyFileExists(reportPath, "report")
}

// VerifyRevisionExists checks that a regular file exists at revisionPath after the revision-phase AI exits (R5, R8).
// When apply was requested, if the AI did not write to the prompt output path, returns error (caller should exit 2).
func VerifyRevisionExists(revisionPath string) error {
	return verifyFileExists(revisionPath, "revision")
}

func verifyFileExists(path string, label string) error {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s file not found at %s", label, path)
		}
		return fmt.Errorf("%s path: %w", label, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%s file not found at %s: path is a directory", label, path)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("%s file not found at %s (symlink target missing)", label, path)
			}
			return fmt.Errorf("%s path: %w", label, err)
		}
		if target.IsDir() {
			return fmt.Errorf("%s file not found at %s: symlink target is a directory", label, path)
		}
	}
	return nil
}
