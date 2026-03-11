package runloop

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/maxdunn/ralph/internal/backend"
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
	var reported []string
	code, err := Run(RunOptions{
		Command:     "true",
		PromptBytes: []byte("p"),
		Loop:        loop,
		Invoker:     invokerAdapter(invoker),
		Reporter:    func(msg string) { reported = append(reported, msg) },
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitMaxIterations {
		t.Errorf("exit code = %d, want %d (ExitMaxIterations)", code, ExitMaxIterations)
	}
	all := strings.Join(reported, " ")
	if !strings.Contains(all, "Stopped after 2 iteration(s)") || !strings.Contains(all, "max: 2") {
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

func TestRun_DryRun_PrintsAssembledPromptAndExitsZero(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.Preamble = "You are helpful."
	callCount := 0
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int) ([]byte, int, error) {
		callCount++
		return nil, 0, nil
	}
	var reported string
	// Capture stdout for dry-run output
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	oldStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	code, err := Run(RunOptions{
		Command:     "true",
		PromptBytes: []byte("actual prompt content"),
		Loop:        loop,
		DryRun:      true,
		Invoker:     invokerAdapter(invoker),
		Reporter:    func(msg string) { reported = msg },
	})
	w.Close()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitSuccess {
		t.Errorf("exit code = %d, want %d (ExitSuccess)", code, ExitSuccess)
	}
	if callCount != 0 {
		t.Errorf("dry-run must not invoke backend; invoker was called %d times", callCount)
	}
	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "Iteration 1") || !strings.Contains(string(out), "You are helpful.") || !strings.Contains(string(out), "actual prompt content") {
		t.Errorf("dry-run stdout = %q; expected preamble + prompt content", out)
	}
	if !strings.Contains(reported, "Dry-run") || !strings.Contains(reported, "no run was performed") {
		t.Errorf("reported = %q", reported)
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
	var reported []string
	code, err := Run(RunOptions{
		Command:     "true",
		PromptBytes: []byte("p"),
		Loop:        loop,
		Invoker:     invokerAdapter(invoker),
		Reporter:    func(msg string) { reported = append(reported, msg) },
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitFailureThreshold {
		t.Errorf("exit code = %d, want %d (ExitFailureThreshold)", code, ExitFailureThreshold)
	}
	all := strings.Join(reported, " ")
	if !strings.Contains(all, "2 consecutive failure(s)") || !strings.Contains(all, "threshold: 2") {
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

// TestRun_NoSignalTreatedAsFailure verifies T3.8/O001/R009: when the process exits
// without success or failure signal (e.g. exit 0 but no signal in output),
// the iteration is treated as failure; consecutive-failure count increments and
// continue/exit follows the same threshold logic as failure-signal.
func TestRun_NoSignalTreatedAsFailure(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 5
	loop.FailureThreshold = 2
	loop.SuccessSignal = "<promise>SUCCESS</promise>"
	loop.FailureSignal = "<promise>FAILURE</promise>"
	callCount := 0
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int) ([]byte, int, error) {
		callCount++
		// Iteration 1: no signal (e.g. process exited 0 but no markers).
		if callCount == 1 {
			return []byte("output with no signal"), 0, nil
		}
		// Iteration 2: again no signal → threshold reached, exit.
		return []byte("again no signal"), 0, nil
	}
	var reported []string
	code, err := Run(RunOptions{
		Command:     "true",
		PromptBytes: []byte("p"),
		Loop:        loop,
		Invoker:     invokerAdapter(invoker),
		Reporter:    func(msg string) { reported = append(reported, msg) },
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitFailureThreshold {
		t.Errorf("exit code = %d, want %d (ExitFailureThreshold)", code, ExitFailureThreshold)
	}
	if callCount != 2 {
		t.Errorf("invoker called %d times, want 2 (two no-signal iterations then exit)", callCount)
	}
	all := strings.Join(reported, " ")
	if !strings.Contains(all, "without success or failure signal") {
		t.Errorf("reported = %q (should distinguish no-signal)", reported)
	}
	if !strings.Contains(all, "threshold: 2") {
		t.Errorf("reported = %q", reported)
	}
}

// TestRun_NoSignalBelowThreshold_Continues verifies that a single no-signal
// iteration increments consecutive failures but loop continues when below threshold.
func TestRun_NoSignalBelowThreshold_Continues(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 5
	loop.FailureThreshold = 2
	loop.SuccessSignal = "<promise>SUCCESS</promise>"
	loop.FailureSignal = "<promise>FAILURE</promise>"
	callCount := 0
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int) ([]byte, int, error) {
		callCount++
		if callCount == 1 {
			return []byte("no signal"), 0, nil
		}
		return []byte("<promise>SUCCESS</promise>"), 0, nil
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
		t.Errorf("exit code = %d, want ExitSuccess (one no-signal then success)", code)
	}
	if callCount != 2 {
		t.Errorf("invoker called %d times, want 2", callCount)
	}
	if !strings.Contains(reported, "2 iterations") {
		t.Errorf("reported = %q", reported)
	}
}

// TestRun_InvocationErrorTreatedAsNoSignal verifies that when the backend
// returns an error (e.g. timeout, crash), the iteration is treated as
// no-signal failure and threshold/continue logic applies (T3.8/O001/R009).
func TestRun_InvocationErrorTreatedAsNoSignal(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 5
	loop.FailureThreshold = 2
	loop.SuccessSignal = "<promise>SUCCESS</promise>"
	callCount := 0
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int) ([]byte, int, error) {
		callCount++
		if callCount == 1 {
			return nil, -1, backend.ErrTimeout
		}
		return nil, -1, backend.ErrTimeout
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
		t.Errorf("exit code = %d, want ExitFailureThreshold", code)
	}
	if callCount != 2 {
		t.Errorf("invoker called %d times, want 2", callCount)
	}
	if !strings.Contains(reported, "without success or failure signal") || !strings.Contains(reported, "invocation error") {
		t.Errorf("reported = %q (should mention invocation error)", reported)
	}
}

// TestRun_InterruptReturnsExitInterrupt verifies T3.9/O004/R005: when the interrupt
// context is cancelled (e.g. SIGINT/SIGTERM), Run returns ExitInterrupt (130).
func TestRun_InterruptReturnsExitInterrupt(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 5
	code, err := Run(RunOptions{
		Command:          "true",
		PromptBytes:      []byte("p"),
		Loop:             loop,
		Invoker:          invokerAdapter(func(_ string, _ []byte, _ string, _ []string, _ int) ([]byte, int, error) { return []byte("x"), 0, nil }),
		InterruptContext: ctx,
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitInterrupt {
		t.Errorf("exit code = %d, want ExitInterrupt (%d)", code, ExitInterrupt)
	}
}

// TestRun_IterationStatistics verifies T3.12/O004/R008: after a multi-iteration
// run, iteration statistics (min/max/mean duration) are reported.
func TestRun_IterationStatistics(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 5
	callCount := 0
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int) ([]byte, int, error) {
		callCount++
		if callCount == 3 {
			return []byte("<promise>SUCCESS</promise>"), 0, nil
		}
		return []byte("working"), 0, nil
	}
	var reported []string
	report := func(msg string) { reported = append(reported, msg) }
	code, err := Run(RunOptions{
		Command:     "true",
		PromptBytes: []byte("p"),
		Loop:        loop,
		Invoker:     invokerAdapter(invoker),
		Reporter:    report,
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitSuccess {
		t.Errorf("exit code = %d, want ExitSuccess", code)
	}
	if callCount != 3 {
		t.Errorf("invoker called %d times, want 3", callCount)
	}
	// Completion message and iteration stats (2+ iterations).
	var hasCompletion, hasStats bool
	for _, s := range reported {
		if strings.Contains(s, "Completed successfully") && strings.Contains(s, "3 iterations") {
			hasCompletion = true
		}
		if strings.Contains(s, "Iteration stats:") && strings.Contains(s, "min") && strings.Contains(s, "max") && strings.Contains(s, "mean") && strings.Contains(s, "3 iterations") {
			hasStats = true
		}
	}
	if !hasCompletion {
		t.Errorf("reported messages = %v; expected completion message", reported)
	}
	if !hasStats {
		t.Errorf("reported messages = %v; expected iteration stats line (min/max/mean)", reported)
	}
}

// TestRun_SingleIteration_NoIterationStats verifies that a single-iteration
// success does not emit the "Iteration stats:" line (only completion message).
func TestRun_SingleIteration_NoIterationStats(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 3
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int) ([]byte, int, error) {
		return []byte("<promise>SUCCESS</promise>"), 0, nil
	}
	var reported []string
	report := func(msg string) { reported = append(reported, msg) }
	code, err := Run(RunOptions{
		Command:     "true",
		PromptBytes: []byte("p"),
		Loop:        loop,
		Invoker:     invokerAdapter(invoker),
		Reporter:    report,
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitSuccess {
		t.Errorf("exit code = %d, want ExitSuccess", code)
	}
	for _, s := range reported {
		if strings.Contains(s, "Iteration stats:") {
			t.Errorf("single-iteration run should not report iteration stats; got %q", s)
		}
	}
}

// TestRun_InterruptBetweenIterations verifies that if the interrupt context is
// cancelled after the first iteration completes, Run returns ExitInterrupt
// (checked at start of next iteration and after Invoke returns).
func TestRun_InterruptBetweenIterations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	callCount := 0
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int) ([]byte, int, error) {
		callCount++
		if callCount == 1 {
			return []byte("no signal yet"), 0, nil
		}
		cancel() // simulate interrupt; Run will see it after this invocation returns
		return []byte("no signal"), 0, nil
	}
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 5
	code, err := Run(RunOptions{
		Command:          "true",
		PromptBytes:      []byte("p"),
		Loop:             loop,
		Invoker:          invokerAdapter(invoker),
		InterruptContext: ctx,
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitInterrupt {
		t.Errorf("exit code = %d, want ExitInterrupt (%d)", code, ExitInterrupt)
	}
	// Run checks ctx after each Invoke; so we may have run 2 iterations then seen cancel.
	if callCount < 1 || callCount > 2 {
		t.Errorf("invoker called %d times, want 1 or 2", callCount)
	}
}
