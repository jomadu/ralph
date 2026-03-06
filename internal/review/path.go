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

func resolveExplicitPath(explicitPath string) (string, error) {
	path, err := filepath.Abs(explicitPath)
	if err != nil {
		return "", fmt.Errorf("review output path invalid: %w", err)
	}

	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return "", fmt.Errorf("review output path is a directory (must be a file): %s", path)
		}
		// Existing file: overwrite allowed per R3
		return path, nil
	}
	if !os.IsNotExist(err) {
		return "", fmt.Errorf("review output path: %w", err)
	}

	// Path does not exist; ensure parent exists and is writable
	parent := filepath.Dir(path)
	parentInfo, err := os.Stat(parent)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("review output path invalid: parent directory does not exist: %s", parent)
		}
		return "", fmt.Errorf("review output path: %w", err)
	}
	if !parentInfo.IsDir() {
		return "", fmt.Errorf("review output path invalid: parent is not a directory: %s", parent)
	}

	if err := checkWritable(parent); err != nil {
		return "", fmt.Errorf("review output path unwritable: %w", err)
	}
	return path, nil
}

func resolveTempPath() (string, error) {
	dir := os.TempDir()
	if dir == "" {
		return "", fmt.Errorf("system temporary directory not available")
	}
	dirInfo, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("temp directory unavailable: %w", err)
	}
	if !dirInfo.IsDir() {
		return "", fmt.Errorf("temp path is not a directory: %s", dir)
	}
	if err := checkWritable(dir); err != nil {
		return "", fmt.Errorf("temp directory not writable: %w", err)
	}

	// Create a unique path: create then remove so we get a path the OS guarantees is unique
	f, err := os.CreateTemp(dir, "ralph-review-*.md")
	if err != nil {
		return "", fmt.Errorf("temp directory not writable: %w", err)
	}
	path := f.Name()
	_ = f.Close()
	if err := os.Remove(path); err != nil {
		return "", fmt.Errorf("failed to reserve temp path: %w", err)
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
