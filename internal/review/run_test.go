package review

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_reportPathRequired(t *testing.T) {
	err := Run([]byte("prompt"), RunOptions{})
	if err == nil {
		t.Fatal("Run(empty ReportPath) err = nil, want ErrReportPathRequired")
	}
	if !errors.Is(err, ErrReportPathRequired) {
		t.Errorf("err = %v, want ErrReportPathRequired", err)
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
