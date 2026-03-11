package review

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_defaultReportPath(t *testing.T) {
	dir := t.TempDir()
	prompt := []byte("# prompt")
	err := Run(prompt, RunOptions{WorkingDir: dir})
	if err != nil {
		t.Fatalf("Run(empty ReportPath, WorkingDir set) err = %v", err)
	}
	defaultPath := filepath.Join(dir, DefaultReportFilename)
	data, err := os.ReadFile(defaultPath)
	if err != nil {
		t.Fatalf("ReadFile(default report) err = %v", err)
	}
	if !strings.Contains(string(data), "ralph-review: status=ok") {
		t.Errorf("default report missing summary line")
	}
}

func TestRun_reportPathIsDirectory(t *testing.T) {
	dir := t.TempDir()
	err := Run([]byte("x"), RunOptions{ReportPath: dir})
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
	err := Run(prompt, RunOptions{ReportPath: reportPath})
	if err != nil {
		t.Fatalf("Run err = %v", err)
	}
	data, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("ReadFile(report) err = %v", err)
	}
	body := string(data)
	if !strings.Contains(body, "ralph-review: status=ok") {
		t.Errorf("report missing summary line: %s", body)
	}
	if !strings.Contains(body, "# My prompt") {
		t.Errorf("report missing revision: %s", body)
	}
	if !strings.Contains(body, "---") {
		t.Errorf("report missing separator: %s", body)
	}
}
