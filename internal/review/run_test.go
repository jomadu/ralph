package review

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maxdunn/ralph/internal/backend"
)

// mockInvoker simulates the AI creating the five files in the report dir (embedded in the prompt).
// It writes result.json and, when revisionContent is non-empty, revision.md into the directory
// found after "**Directory where you must create the files:**" in the prompt, then returns a short confirmation.
func mockInvoker(revisionContent string) backend.Invoker {
	return invokerAdapter(func(_ string, prompt []byte, _ string, _ []string, _ int, _ io.Writer) ([]byte, int, error) {
		reportDir := extractReportDirFromPrompt(prompt)
		if reportDir == "" {
			return []byte("error: no report dir in prompt"), 1, nil
		}
		if err := os.MkdirAll(reportDir, 0755); err != nil {
			return nil, 0, err
		}
		result := resultJSON{Status: "ok", Errors: 0, Warnings: 0}
		b, _ := json.Marshal(result)
		if err := os.WriteFile(filepath.Join(reportDir, ReportResultJSON), b, 0644); err != nil {
			return nil, 0, err
		}
		if revisionContent != "" {
			if err := os.WriteFile(filepath.Join(reportDir, ReportRevisionMD), []byte(revisionContent), 0644); err != nil {
				return nil, 0, err
			}
		}
		return []byte("Created the review report at " + reportDir + ". Files: result.json, summary.md, original.md, revision.md, diff.md."), 0, nil
	})
}

// extractReportDirFromPrompt finds the report directory path in the assembled prompt.
// The path appears on its own line between "**Directory where you must create the files:**" and "**Files to create".
func extractReportDirFromPrompt(prompt []byte) string {
	s := string(prompt)
	dirHeader := "**Directory where you must create the files:**"
	filesHeader := "**Files to create"
	idx := strings.Index(s, dirHeader)
	if idx < 0 {
		return ""
	}
	block := s[idx+len(dirHeader):]
	if end := strings.Index(block, filesHeader); end >= 0 {
		block = block[:end]
	}
	for _, line := range strings.Split(block, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// The interpolated path is always absolute; it's the only such line in this block.
		if filepath.IsAbs(line) {
			return line
		}
	}
	return ""
}

func TestRun_emptyCommand_exit2(t *testing.T) {
	_, err := Run([]byte("x"), RunOptions{ReportPath: "/tmp/report.txt"})
	if err == nil {
		t.Fatal("Run(empty Command) err = nil, want ErrAICommandRequired")
	}
	if !IsExit2(err) {
		t.Error("IsExit2(err) = false, want true")
	}
	if !errors.Is(err, ErrAICommandRequired) {
		t.Errorf("err = %v, want ErrAICommandRequired", err)
	}
}

func TestRun_defaultReportPath(t *testing.T) {
	dir := t.TempDir()
	prompt := []byte("# prompt")
	opts := RunOptions{WorkingDir: dir, Command: "echo", Invoker: mockInvoker("")}
	_, err := Run(prompt, opts)
	if err != nil {
		t.Fatalf("Run(empty ReportPath, WorkingDir set) err = %v", err)
	}
	defaultDir := filepath.Join(dir, DefaultReportDir)
	resultPath := filepath.Join(defaultDir, ReportResultJSON)
	data, err := os.ReadFile(resultPath)
	if err != nil {
		t.Fatalf("ReadFile(default result.json) err = %v", err)
	}
	var result resultJSON
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Unmarshal result.json: %v", err)
	}
	if result.Status != "ok" {
		t.Errorf("result.status = %q, want ok", result.Status)
	}
}

func TestRun_reportPathIsFile_exit2(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "existing-file.txt")
	if err := os.WriteFile(filePath, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Run([]byte("x"), RunOptions{ReportPath: filePath, Command: "echo", Invoker: mockInvoker("x")})
	if err == nil {
		t.Fatal("Run(report path = existing file) err = nil, want error")
	}
	if !IsExit2(err) {
		t.Error("IsExit2(err) = false, want true")
	}
}

func TestRun_writesReportDir(t *testing.T) {
	dir := t.TempDir()
	reportDir := filepath.Join(dir, "report")
	prompt := []byte("# My prompt\nDo the thing.")
	opts := RunOptions{ReportPath: reportDir, WorkingDir: dir, Command: "echo", Invoker: mockInvoker("# My prompt\nDo the thing.")}
	_, err := Run(prompt, opts)
	if err != nil {
		t.Fatalf("Run err = %v", err)
	}
	data, err := os.ReadFile(filepath.Join(reportDir, ReportResultJSON))
	if err != nil {
		t.Fatalf("ReadFile(result.json) err = %v", err)
	}
	var result resultJSON
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Unmarshal result.json: %v", err)
	}
	if result.Status != "ok" {
		t.Errorf("result.status = %q, want ok", result.Status)
	}
	revisionData, err := os.ReadFile(filepath.Join(reportDir, ReportRevisionMD))
	if err != nil {
		t.Fatalf("ReadFile(revision.md) err = %v", err)
	}
	if !strings.Contains(string(revisionData), "# My prompt") {
		t.Errorf("revision.md missing prompt content: %s", revisionData)
	}
}

func TestRun_apply_writesNewFile(t *testing.T) {
	dir := t.TempDir()
	reportDir := filepath.Join(dir, "report")
	applyPath := filepath.Join(dir, "revision.md")
	opts := RunOptions{
		ReportPath:       reportDir,
		PromptOutputPath: applyPath,
		WorkingDir:       dir,
		Apply:            true,
		Yes:              true,
		Command:          "echo",
		Invoker:          mockInvoker("original"),
	}
	_, err := Run([]byte("original"), opts)
	if err != nil {
		t.Fatalf("Run(apply to new file) err = %v", err)
	}
	data, err := os.ReadFile(applyPath)
	if err != nil {
		t.Fatalf("ReadFile(revision) err = %v", err)
	}
	if !strings.Contains(string(data), "original") {
		t.Errorf("revision file content unexpected: %s", data)
	}
}

func TestRun_apply_overwriteWithYes(t *testing.T) {
	dir := t.TempDir()
	reportDir := filepath.Join(dir, "report")
	applyPath := filepath.Join(dir, "revision.md")
	if err := os.WriteFile(applyPath, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	opts := RunOptions{
		ReportPath:       reportDir,
		PromptOutputPath: applyPath,
		WorkingDir:       dir,
		Apply:            true,
		Yes:              true,
		Command:          "echo",
		Invoker:          mockInvoker("new content"),
	}
	_, err := Run([]byte("new content"), opts)
	if err != nil {
		t.Fatalf("Run(apply overwrite with Yes) err = %v", err)
	}
	data, err := os.ReadFile(applyPath)
	if err != nil {
		t.Fatalf("ReadFile(revision) err = %v", err)
	}
	if !strings.Contains(string(data), "new content") {
		t.Errorf("revision file missing user content: %q", data)
	}
}

func TestRun_apply_overwriteNonInteractiveWithoutYes_exit2(t *testing.T) {
	dir := t.TempDir()
	reportDir := filepath.Join(dir, "report")
	applyPath := filepath.Join(dir, "revision.md")
	if err := os.WriteFile(applyPath, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	opts := RunOptions{
		ReportPath:       reportDir,
		PromptOutputPath: applyPath,
		WorkingDir:       dir,
		Apply:            true,
		Yes:              false,
		NonInteractive:   true,
		Command:          "echo",
		Invoker:          mockInvoker("new"),
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
		ReportPath:       filepath.Join(dir, "report"),
		PromptOutputPath: "",
		SourcePath:       "",
		WorkingDir:       dir,
		Apply:            true,
		Command:          "echo",
		Invoker:          mockInvoker("x"),
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
	reportDir := filepath.Join(dir, "report")
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w
	defer func() { os.Stderr = oldStderr }()

	_, runErr := Run([]byte("# p"), RunOptions{ReportPath: reportDir, WorkingDir: dir, Quiet: true, Command: "echo", Invoker: mockInvoker("# p")})
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
	reportDir := filepath.Join(dir, "report")
	applyPath := filepath.Join(dir, "rev.md")
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w
	defer func() { os.Stderr = oldStderr }()

	_, runErr := Run([]byte("content"), RunOptions{
		ReportPath:       reportDir,
		PromptOutputPath: applyPath,
		WorkingDir:       dir,
		Apply:            true,
		Yes:              true,
		Verbose:          true,
		Command:          "echo",
		Invoker:          mockInvoker("content"),
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
