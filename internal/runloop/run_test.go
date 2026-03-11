package runloop

import (
	"strings"
	"testing"

	"github.com/maxdunn/ralph/internal/config"
)

func TestRun_SuccessOnFirstIteration(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 3
	var reported string
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int) ([]byte, int, error) {
		return []byte("output with <promise>SUCCESS</promise> in it"), 0, nil
	}

	code, err := Run(RunOptions{
		Command:     "true",
		PromptBytes: []byte("prompt"),
		Loop:        loop,
		Invoker:     invokerAdapter(invoker),
		Reporter:    func(msg string) { reported = msg },
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitSuccess {
		t.Errorf("exit code = %d, want %d (ExitSuccess)", code, ExitSuccess)
	}
	if reported == "" {
		t.Error("expected completion message to be reported")
	}
	if !strings.Contains(reported, "Completed successfully") || !strings.Contains(reported, "1 iteration") {
		t.Errorf("reported message = %q", reported)
	}
}

func TestRun_SuccessOnSecondIteration(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 5
	callCount := 0
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int) ([]byte, int, error) {
		callCount++
		if callCount == 2 {
			return []byte("<promise>SUCCESS</promise>"), 0, nil
		}
		return []byte("still working"), 0, nil
	}
	var reported string
	code, err := Run(RunOptions{
		Command:     "true",
		PromptBytes: []byte("p"),
		Loop:        loop,
		Invoker:     invokerAdapter(invoker),
		Reporter:    func(msg string) { reported = msg },
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitSuccess {
		t.Errorf("exit code = %d, want %d", code, ExitSuccess)
	}
	if callCount != 2 {
		t.Errorf("invoker called %d times, want 2", callCount)
	}
	if !strings.Contains(reported, "2 iterations") {
		t.Errorf("reported = %q", reported)
	}
}

func TestRun_MaxIterationsWithoutSuccess(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 2
	loop.SuccessSignal = "<promise>SUCCESS</promise>"
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int) ([]byte, int, error) {
		return []byte("no signal here"), 0, nil
	}
	var reported string
	code, err := Run(RunOptions{
		Command:     "true",
		PromptBytes: []byte("p"),
		Loop:        loop,
		Invoker:     invokerAdapter(invoker),
		Reporter:    func(msg string) { reported = msg },
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitMaxIterations {
		t.Errorf("exit code = %d, want %d (ExitMaxIterations)", code, ExitMaxIterations)
	}
	if !strings.Contains(reported, "Stopped after 2 iteration(s)") {
		t.Errorf("reported = %q", reported)
	}
}

func TestRun_InvalidCommand(t *testing.T) {
	loop := config.DefaultLoopSettings()
	code, err := Run(RunOptions{
		Command:     "nonexistent-ralph-command-xyz",
		PromptBytes: []byte("p"),
		Loop:        loop,
	})
	if err == nil {
		t.Fatal("expected error for invalid command")
	}
	if code != ExitErrorPreLoop {
		t.Errorf("exit code = %d, want %d (ExitErrorPreLoop)", code, ExitErrorPreLoop)
	}
}

func TestRun_FailureSignalBelowThreshold_Continues(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 5
	loop.FailureThreshold = 3
	loop.SuccessSignal = "<promise>SUCCESS</promise>"
	loop.FailureSignal = "<promise>FAILURE</promise>"
	callCount := 0
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int) ([]byte, int, error) {
		callCount++
		if callCount == 1 {
			return []byte("<promise>FAILURE</promise>"), 0, nil
		}
		if callCount == 2 {
			return []byte("<promise>SUCCESS</promise>"), 0, nil
		}
		return []byte(""), 0, nil
	}
	var reported string
	code, err := Run(RunOptions{
		Command:     "true",
		PromptBytes: []byte("p"),
		Loop:        loop,
		Invoker:     invokerAdapter(invoker),
		Reporter:    func(msg string) { reported = msg },
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitSuccess {
		t.Errorf("exit code = %d, want ExitSuccess (success on iteration 2 after one failure)", code)
	}
	if callCount != 2 {
		t.Errorf("invoker called %d times, want 2", callCount)
	}
	if !strings.Contains(reported, "2 iterations") {
		t.Errorf("reported = %q", reported)
	}
}

func TestRun_FailureThresholdReached_ExitsWithCode(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 5
	loop.FailureThreshold = 2
	loop.SuccessSignal = "<promise>SUCCESS</promise>"
	loop.FailureSignal = "<promise>FAILURE</promise>"
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int) ([]byte, int, error) {
		return []byte("<promise>FAILURE</promise>"), 0, nil
	}
	var reported string
	code, err := Run(RunOptions{
		Command:     "true",
		PromptBytes: []byte("p"),
		Loop:        loop,
		Invoker:     invokerAdapter(invoker),
		Reporter:    func(msg string) { reported = msg },
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitFailureThreshold {
		t.Errorf("exit code = %d, want %d (ExitFailureThreshold)", code, ExitFailureThreshold)
	}
	if !strings.Contains(reported, "2 consecutive failure(s)") || !strings.Contains(reported, "threshold: 2") {
		t.Errorf("reported = %q", reported)
	}
}

// TestRun_StaticPrecedence_BothSignalsPresent verifies T3.6/O001/R006: when both
// success and failure signals appear in the same output, static precedence applies
// (success checked first), so the iteration is treated as success.
func TestRun_StaticPrecedence_BothSignalsPresent(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 3
	loop.SuccessSignal = "DONE"
	loop.FailureSignal = "FAIL"
	callCount := 0
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int) ([]byte, int, error) {
		callCount++
		// Single iteration output contains both signals; success must win.
		return []byte("output FAIL and DONE together"), 0, nil
	}
	var reported string
	code, err := Run(RunOptions{
		Command:     "true",
		PromptBytes: []byte("p"),
		Loop:        loop,
		Invoker:     invokerAdapter(invoker),
		Reporter:    func(msg string) { reported = msg },
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitSuccess {
		t.Errorf("exit code = %d, want ExitSuccess (static precedence: success wins when both present)", code)
	}
	if callCount != 1 {
		t.Errorf("invoker called %d times, want 1 (success on first iteration)", callCount)
	}
	if !strings.Contains(reported, "Completed successfully") || !strings.Contains(reported, "1 iteration") {
		t.Errorf("reported = %q", reported)
	}
}

func TestRun_SuccessResetsConsecutiveFailures(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 5
	loop.FailureThreshold = 2
	loop.SuccessSignal = "DONE"
	loop.FailureSignal = "FAIL"
	callCount := 0
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int) ([]byte, int, error) {
		callCount++
		switch callCount {
		case 1:
			return []byte("FAIL"), 0, nil
		case 2:
			return []byte("DONE"), 0, nil
		case 3:
			return []byte("FAIL"), 0, nil
		case 4:
			return []byte("DONE"), 0, nil
		}
		return []byte(""), 0, nil
	}
	var reported string
	code, err := Run(RunOptions{
		Command:     "true",
		PromptBytes: []byte("p"),
		Loop:        loop,
		Invoker:     invokerAdapter(invoker),
		Reporter:    func(msg string) { reported = msg },
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitSuccess {
		t.Errorf("exit code = %d, want ExitSuccess", code)
	}
	if callCount != 2 {
		t.Errorf("invoker called %d times, want 2 (success on 2nd iteration)", callCount)
	}
	if !strings.Contains(reported, "2 iterations") {
		t.Errorf("reported = %q", reported)
	}
}
