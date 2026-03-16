package review

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestParseSummaryLine(t *testing.T) {
	tests := []struct {
		line         string
		wantStatus   SummaryStatus
		wantErrors   int
		wantWarnings int
		wantOK       bool
	}{
		{"ralph-review: status=ok", StatusOK, 0, 0, true},
		{"ralph-review: status=ok errors=0", StatusOK, 0, 0, true},
		{"ralph-review: status=ok errors=0 warnings=0", StatusOK, 0, 0, true},
		{"ralph-review: status=errors errors=2", StatusErrors, 2, 0, true},
		{"ralph-review: status=errors errors=1 warnings=0", StatusErrors, 1, 0, true},
		{"ralph-review: status=warnings warnings=1 errors=0", StatusWarnings, 0, 1, true},
		{"  ralph-review: status=ok  ", StatusOK, 0, 0, true},
		{"ralph-review: status=ok errors=10", StatusOK, 10, 0, true},
		{"not-a-summary", "", 0, 0, false},
		{"ralph-review: status=unknown", "", 0, 0, false},
		{"", "", 0, 0, false},
	}
	for _, tt := range tests {
		gotStatus, gotErrors, gotWarnings, gotOK := ParseSummaryLine(tt.line)
		if gotOK != tt.wantOK || gotStatus != tt.wantStatus || gotErrors != tt.wantErrors || gotWarnings != tt.wantWarnings {
			t.Errorf("ParseSummaryLine(%q) = status=%q errors=%d warnings=%d ok=%v, want status=%q errors=%d warnings=%d ok=%v",
				tt.line, gotStatus, gotErrors, gotWarnings, gotOK, tt.wantStatus, tt.wantErrors, tt.wantWarnings, tt.wantOK)
		}
	}
}

func TestParseSummaryFromReport(t *testing.T) {
	report := "Some narrative.\n\nralph-review: status=errors errors=2\n\n---\n\nrevision"
	status, errs, warns, ok := ParseSummaryFromReport([]byte(report))
	if !ok {
		t.Fatal("ParseSummaryFromReport ok = false, want true")
	}
	if status != StatusErrors || errs != 2 || warns != 0 {
		t.Errorf("ParseSummaryFromReport = status=%q errors=%d warnings=%d, want status=errors errors=2 warnings=0", status, errs, warns)
	}

	// No summary line in content
	_, _, _, ok = ParseSummaryFromReport([]byte("only narrative\nno summary"))
	if ok {
		t.Error("ParseSummaryFromReport(no summary line) ok = true, want false")
	}
}

func TestExitCodeFromSummary(t *testing.T) {
	tests := []struct {
		status SummaryStatus
		errors int
		want   int
	}{
		{StatusOK, 0, 0},
		{StatusOK, 1, 1},
		{StatusErrors, 0, 1},
		{StatusErrors, 2, 1},
		{StatusWarnings, 0, 0},
		{StatusWarnings, 1, 1},
		{"", 0, 1},
	}
	for _, tt := range tests {
		got := ExitCodeFromSummary(tt.status, tt.errors)
		if got != tt.want {
			t.Errorf("ExitCodeFromSummary(%q, %d) = %d, want %d", tt.status, tt.errors, got, tt.want)
		}
	}
}

func TestRun_exitCodeMatchesResultJSON(t *testing.T) {
	dir := t.TempDir()
	reportDir := filepath.Join(dir, "report")
	// Mock invoker writes result.json (status=ok) so exit code is 0.
	mock := invokerAdapter(func(_ string, prompt []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
		reportDir := extractReportDirFromPrompt(prompt)
		if reportDir != "" {
			_ = os.MkdirAll(reportDir, 0755)
			b, _ := json.Marshal(resultJSON{Status: "ok", Errors: 0, Warnings: 0})
			_ = os.WriteFile(filepath.Join(reportDir, ReportResultJSON), b, 0644)
		}
		return []byte("Created the review report."), 0, nil
	})
	opts := RunOptions{ReportPath: reportDir, WorkingDir: dir, Command: "echo", Invoker: mock}
	code, err := Run([]byte("# prompt"), opts)
	if err != nil {
		t.Fatalf("Run err = %v", err)
	}
	data, err := os.ReadFile(filepath.Join(reportDir, ReportResultJSON))
	if err != nil {
		t.Fatalf("ReadFile(result.json) err = %v", err)
	}
	var result resultJSON
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal("result.json invalid:", err)
	}
	wantCode := ExitCodeFromSummary(SummaryStatus(result.Status), result.Errors)
	if code != wantCode {
		t.Errorf("Run exit code = %d, want %d (from result status=%q errors=%d)", code, wantCode, result.Status, result.Errors)
	}
	if code != 0 && code != 1 {
		t.Errorf("Run exit code = %d, want 0 or 1", code)
	}
}
