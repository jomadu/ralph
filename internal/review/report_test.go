package review

import (
	"strings"
	"testing"
)

func TestFormatSummaryLine(t *testing.T) {
	tests := []struct {
		status   SummaryStatus
		errors   int
		warnings int
		want     string
	}{
		{StatusOK, 0, 0, "ralph-review: status=ok errors=0 warnings=0"},
		{StatusErrors, 2, 0, "ralph-review: status=errors errors=2 warnings=0"},
		{StatusWarnings, 0, 1, "ralph-review: status=warnings errors=0 warnings=1"},
	}
	for _, tt := range tests {
		got := FormatSummaryLine(tt.status, tt.errors, tt.warnings)
		if got != tt.want {
			t.Errorf("FormatSummaryLine(%q, %d, %d) = %q, want %q", tt.status, tt.errors, tt.warnings, got, tt.want)
		}
	}
}

func TestReport_String(t *testing.T) {
	r := &Report{
		Narrative:   "Some feedback.",
		SummaryLine: "ralph-review: status=ok errors=0 warnings=0",
		Revision:    "revised prompt",
	}
	s := r.String()
	if !strings.Contains(s, "Some feedback.") {
		t.Error("Report.String() missing narrative")
	}
	if !strings.Contains(s, "ralph-review: status=ok") {
		t.Error("Report.String() missing summary line")
	}
	if !strings.Contains(s, "---") {
		t.Error("Report.String() missing separator")
	}
	if !strings.Contains(s, "revised prompt") {
		t.Error("Report.String() missing revision")
	}
}

func TestGenerateReport(t *testing.T) {
	prompt := []byte("# Task\nDo it.")
	report := GenerateReport(prompt)
	if report == nil {
		t.Fatal("GenerateReport returned nil")
	}
	if report.Narrative == "" {
		t.Error("Narrative empty")
	}
	// T5.6: narrative is structured by dimension
	if !strings.Contains(report.Narrative, "Signal and state") || !strings.Contains(report.Narrative, "Scope and convergence") {
		t.Errorf("Narrative missing dimension labels: %s", report.Narrative)
	}
	if !strings.HasPrefix(report.SummaryLine, "ralph-review:") {
		t.Errorf("SummaryLine = %q, want ralph-review:...", report.SummaryLine)
	}
	// Revision includes prompt and may include suggestions (T5.6)
	if !strings.Contains(report.Revision, "# Task") || !strings.Contains(report.Revision, "Do it.") {
		t.Errorf("Revision missing prompt content: %q", report.Revision)
	}
	body := report.String()
	if !strings.Contains(body, report.SummaryLine) {
		t.Error("Report body missing summary line")
	}
	if !strings.Contains(body, "# Task") {
		t.Error("Report body missing revision content")
	}
}
