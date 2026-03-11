package review

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_defaultReportPath(t *testing.T) {
	dir := t.TempDir()
	prompt := []byte("# prompt")
	_, err := Run(prompt, RunOptions{WorkingDir: dir})
	if err != nil {
		t.Fatalf("Run(empty ReportPath, WorkingDir set) err = %v", err)
	}
	defaultPath := filepath.Join(dir, DefaultReportFilename)
	data, err := os.ReadFile(defaultPath)
	if err != nil {
		t.Fatalf("ReadFile(default report) err = %v", err)
	}
	if !strings.Contains(string(data), "ralph-review:") {
		t.Errorf("default report missing summary line")
	}
}

func TestRun_reportPathIsDirectory(t *testing.T) {
	dir := t.TempDir()
	_, err := Run([]byte("x"), RunOptions{ReportPath: dir})
	if err == nil {
		t.Fatal("Run(report path = dir) err = nil, want error")
	}
	if !IsExit2(err) {
		t.Error("IsExit2(err) = false, want true")
	}
}

func TestRun_writesReportFile(t *testing.T) {
	dir := t.TempDir()
	reportPath := filepath.Join(dir, "report.txt")
	prompt := []byte("# My prompt\nDo the thing.")
	_, err := Run(prompt, RunOptions{ReportPath: reportPath})
	if err != nil {
		t.Fatalf("Run err = %v", err)
	}
	data, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("ReadFile(report) err = %v", err)
	}
	body := string(data)
	if !strings.Contains(body, "ralph-review:") {
		t.Errorf("report missing summary line: %s", body)
	}
	if !strings.Contains(body, "# My prompt") {
		t.Errorf("report missing revision: %s", body)
	}
	if !strings.Contains(body, "---") {
		t.Errorf("report missing separator: %s", body)
	}
}

func TestRun_apply_writesNewFile(t *testing.T) {
	dir := t.TempDir()
	reportPath := filepath.Join(dir, "report.txt")
	applyPath := filepath.Join(dir, "revision.md")
	opts := RunOptions{
		ReportPath:       reportPath,
		PromptOutputPath: applyPath,
		WorkingDir:       dir,
		Apply:            true,
		Yes:              true,
	}
	_, err := Run([]byte("original"), opts)
	if err != nil {
		t.Fatalf("Run(apply to new file) err = %v", err)
	}
	data, err := os.ReadFile(applyPath)
	if err != nil {
		t.Fatalf("ReadFile(revision) err = %v", err)
	}
	// GenerateReport currently returns input as revision; we wrote that
	if !strings.Contains(string(data), "original") {
		t.Errorf("revision file content unexpected: %s", data)
	}
}

func TestRun_apply_overwriteWithYes(t *testing.T) {
	dir := t.TempDir()
	reportPath := filepath.Join(dir, "report.txt")
	applyPath := filepath.Join(dir, "revision.md")
	if err := os.WriteFile(applyPath, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	opts := RunOptions{
		ReportPath:       reportPath,
		PromptOutputPath: applyPath,
		WorkingDir:       dir,
		Apply:            true,
		Yes:              true,
	}
	_, err := Run([]byte("new content"), opts)
	if err != nil {
		t.Fatalf("Run(apply overwrite with Yes) err = %v", err)
	}
	data, err := os.ReadFile(applyPath)
	if err != nil {
		t.Fatalf("ReadFile(revision) err = %v", err)
	}
	// Revision may include review suggestions block (T5.6); must contain the user content
	if !strings.Contains(string(data), "new content") {
		t.Errorf("revision file missing user content: %q", data)
	}
}

func TestRun_apply_overwriteNonInteractiveWithoutYes_exit2(t *testing.T) {
	dir := t.TempDir()
	reportPath := filepath.Join(dir, "report.txt")
	applyPath := filepath.Join(dir, "revision.md")
	if err := os.WriteFile(applyPath, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	opts := RunOptions{
		ReportPath:       reportPath,
		PromptOutputPath: applyPath,
		WorkingDir:       dir,
		Apply:            true,
		Yes:              false,
		NonInteractive:   true,
	}
	_, err := Run([]byte("new"), opts)
	if err == nil {
		t.Fatal("Run(apply overwrite, non-interactive, no Yes) err = nil, want ErrApplyConfirmationRequired")
	}
	if !IsExit2(err) {
		t.Error("IsExit2(err) = false, want true")
	}
	data, _ := os.ReadFile(applyPath)
	if string(data) != "old" {
		t.Errorf("file was overwritten: %q", data)
	}
}

func TestRun_apply_requiresPromptOutputWhenNoSource(t *testing.T) {
	dir := t.TempDir()
	opts := RunOptions{
		ReportPath:       filepath.Join(dir, "report.txt"),
		PromptOutputPath: "",
		SourcePath:       "",
		WorkingDir:       dir,
		Apply:            true,
	}
	_, err := Run([]byte("x"), opts)
	if err == nil {
		t.Fatal("Run(apply, no prompt-output, no source) err = nil, want error")
	}
	if !errors.Is(err, ErrApplyPromptOutputRequired) && !IsExit2(err) {
		t.Errorf("err = %v, want ErrApplyPromptOutputRequired or IsExit2", err)
	}
}

// TestRun_Quiet_suppressesReportPath ensures --quiet suppresses "Report written to" (T7.2, O004/R006).
func TestRun_Quiet_suppressesReportPath(t *testing.T) {
	dir := t.TempDir()
	reportPath := filepath.Join(dir, "report.txt")
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w
	defer func() { os.Stderr = oldStderr }()

	_, runErr := Run([]byte("# p"), RunOptions{ReportPath: reportPath, WorkingDir: dir, Quiet: true})
	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	serr := buf.String()

	if runErr != nil {
		t.Fatalf("Run: %v", runErr)
	}
	if strings.Contains(serr, "Report written to") {
		t.Errorf("Quiet should suppress report path; stderr contained: %q", serr)
	}
}

// TestRun_Verbose_apply_printsRevisionPath ensures --verbose with --apply prints revision path (T7.2).
func TestRun_Verbose_apply_printsRevisionPath(t *testing.T) {
	dir := t.TempDir()
	reportPath := filepath.Join(dir, "report.txt")
	applyPath := filepath.Join(dir, "rev.md")
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w
	defer func() { os.Stderr = oldStderr }()

	_, runErr := Run([]byte("content"), RunOptions{
		ReportPath:       reportPath,
		PromptOutputPath: applyPath,
		WorkingDir:       dir,
		Apply:            true,
		Yes:              true,
		Verbose:          true,
	})
	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	serr := buf.String()

	if runErr != nil {
		t.Fatalf("Run: %v", runErr)
	}
	if !strings.Contains(serr, "Report written to") {
		t.Errorf("expected 'Report written to' in stderr; got: %q", serr)
	}
	if !strings.Contains(serr, "Revision applied to") {
		t.Errorf("Verbose+apply should print revision path; got: %q", serr)
	}
}
