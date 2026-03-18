package runloop

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/jomadu/ralph/internal/backend"
	"github.com/jomadu/ralph/internal/config"
)

func TestRun_SuccessOnFirstIteration(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 3
	var reported string
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
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
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
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
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
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

// TestRun_QuietLogLevel_stillShowsCompletionMessage ensures O004/R006: completion message is shown even when log level is error (quiet).
func TestRun_QuietLogLevel_stillShowsCompletionMessage(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 3
	loop.LogLevel = "error"
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
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
		t.Errorf("exit code = %d, want ExitSuccess", code)
	}
	if reported == "" || !strings.Contains(reported, "Completed successfully") {
		t.Errorf("quiet mode must still show completion message (R006); got %q", reported)
	}
}

// TestRun_Streaming_passesStreamWriterToInvoker ensures O004/R006: when Streaming is true
// and StreamWriter is set, the invoker is called with that writer so AI stdout is streamed.
func TestRun_Streaming_passesStreamWriterToInvoker(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 2
	loop.Streaming = true
	loop.SuccessSignal = "<promise>SUCCESS</promise>"
	var receivedStreamTo io.Writer
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, streamTo io.Writer) ([]byte, int, error) {
		receivedStreamTo = streamTo
		return []byte("<promise>SUCCESS</promise>"), 0, nil
	}
	streamWriter := &strings.Builder{}
	code, err := Run(RunOptions{
		Command:      "true",
		PromptBytes:  []byte("p"),
		Loop:         loop,
		Invoker:      invokerAdapter(invoker),
		StreamWriter: streamWriter,
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitSuccess {
		t.Errorf("exit code = %d, want ExitSuccess", code)
	}
	if receivedStreamTo != streamWriter {
		t.Errorf("invoker streamTo = %v, want StreamWriter %v (streaming must pass writer for O004/R006)", receivedStreamTo, streamWriter)
	}
}

// TestRun_NoStreaming_invokerReceivesNilStreamTo ensures O004/R006: when Streaming is false,
// the invoker is called with nil streamTo so AI stdout is not streamed (still captured).
func TestRun_NoStreaming_invokerReceivesNilStreamTo(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 2
	loop.Streaming = false
	loop.SuccessSignal = "<promise>SUCCESS</promise>"
	var receivedStreamTo io.Writer
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, streamTo io.Writer) ([]byte, int, error) {
		receivedStreamTo = streamTo
		return []byte("<promise>SUCCESS</promise>"), 0, nil
	}
	code, err := Run(RunOptions{
		Command:      "true",
		PromptBytes:  []byte("p"),
		Loop:         loop,
		Invoker:      invokerAdapter(invoker),
		StreamWriter: os.Stdout, // even if set, run-loop should pass nil when !Streaming
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitSuccess {
		t.Errorf("exit code = %d, want ExitSuccess", code)
	}
	if receivedStreamTo != nil {
		t.Errorf("invoker streamTo = %v, want nil when Streaming is false (O004/R006)", receivedStreamTo)
	}
}

// TestRun_LogLevelError_suppressesDebugAndInfo ensures O004/R006: at log level "error",
// debug and info messages (e.g. "Starting iteration") are not reported; completion message still is.
func TestRun_LogLevelError_suppressesDebugAndInfo(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 3
	loop.LogLevel = "error"
	loop.SuccessSignal = "<promise>SUCCESS</promise>"
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
		return []byte("<promise>SUCCESS</promise>"), 0, nil
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
	if code != ExitSuccess {
		t.Errorf("exit code = %d, want ExitSuccess", code)
	}
	for _, msg := range reported {
		if strings.Contains(msg, "Starting iteration") {
			t.Errorf("log level error must not emit debug messages (R006); got %q", msg)
		}
	}
	hasCompletion := false
	for _, msg := range reported {
		if strings.Contains(msg, "Completed successfully") {
			hasCompletion = true
			break
		}
	}
	if !hasCompletion {
		t.Errorf("completion message must still be shown at log level error (R006); reported = %v", reported)
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
	loop.Preamble = true
	loop.MaxIterations = 10
	callCount := 0
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
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
	outStr := string(out)
	if !strings.Contains(outStr, "# --- LOOP CONFIG ---\n") || !strings.Contains(outStr, "max_iterations: 10") || !strings.Contains(outStr, "failure_threshold: 3") {
		t.Errorf("dry-run stdout missing LOOP CONFIG section or content; got %q", out)
	}
	if !strings.Contains(outStr, "# --- CONTEXT ---\n") || !strings.Contains(outStr, ralphLoopDescription) || !strings.Contains(outStr, "Iteration 1") || !strings.Contains(outStr, "of max 10") {
		t.Errorf("dry-run stdout missing CONTEXT section or content; got %q", out)
	}
	if !strings.Contains(outStr, "# --- INSTRUCTIONS ---\n") || !strings.Contains(outStr, "actual prompt content") {
		t.Errorf("dry-run stdout missing INSTRUCTIONS section or content; got %q", out)
	}
	if !strings.Contains(reported, "Dry-run") || !strings.Contains(reported, "loop config") || !strings.Contains(reported, "no run was performed") {
		t.Errorf("reported = %q", reported)
	}
}

func TestRun_DryRun_WithInvokerContext_ShowsLabelNoDuplicate(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.Preamble = true
	loop.MaxIterations = 10
	loop.Context = "Hello world"
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	code, err := Run(RunOptions{
		Command:     "true",
		PromptBytes: []byte("prompt"),
		Loop:        loop,
		DryRun:      true,
	})
	w.Close()
	if err != nil || code != ExitSuccess {
		t.Fatalf("Run: err=%v code=%d", err, code)
	}
	outStr := string(mustReadAll(r))
	if !strings.Contains(outStr, invokerContextLabel) || !strings.Contains(outStr, "Hello world") {
		t.Errorf("dry-run with -c should show invoker label and content; got %q", outStr)
	}
	// No duplicate: body must not contain "CONTEXT\nHello" (section title is "# --- CONTEXT ---", not "CONTEXT" as a line in the body).
	if strings.Contains(outStr, "CONTEXT\nHello world") {
		t.Errorf("dry-run should not duplicate CONTEXT in body; use invoker label only; got %q", outStr)
	}
}

func mustReadAll(r *os.File) []byte {
	b, _ := io.ReadAll(r)
	return b
}

func TestRun_FailureSignalBelowThreshold_Continues(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 5
	loop.FailureThreshold = 3
	loop.SuccessSignal = "<promise>SUCCESS</promise>"
	loop.FailureSignal = "<promise>FAILURE</promise>"
	callCount := 0
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
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
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
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
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
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
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
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

// TestRun_NullSignalDoesNotCountTowardThreshold verifies that exit 0 with neither
// success nor failure on the last line does not increment the failure counter;
// the loop runs until max iterations.
func TestRun_NullSignalDoesNotCountTowardThreshold(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 4
	loop.FailureThreshold = 2
	loop.SuccessSignal = "<promise>SUCCESS</promise>"
	loop.FailureSignal = "<promise>FAILURE</promise>"
	callCount := 0
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
		callCount++
		return []byte("still working, no signal on last line"), 0, nil
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
	if callCount != 4 {
		t.Errorf("invoker called %d times, want 4", callCount)
	}
}

// TestRun_NullSignalResetsFailureStreak verifies that a neutral iteration after
// a failure signal resets the consecutive-failure count.
func TestRun_NullSignalResetsFailureStreak(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 10
	loop.FailureThreshold = 2
	loop.SuccessSignal = "<promise>SUCCESS</promise>"
	loop.FailureSignal = "<promise>FAILURE</promise>"
	callCount := 0
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
		callCount++
		switch callCount {
		case 1:
			return []byte("<promise>FAILURE</promise>"), 0, nil
		case 2:
			return []byte("progress, no explicit signal"), 0, nil
		case 3:
			return []byte("<promise>FAILURE</promise>"), 0, nil
		default:
			return []byte("<promise>SUCCESS</promise>"), 0, nil
		}
	}
	code, err := Run(RunOptions{
		Command:     "true",
		PromptBytes: []byte("p"),
		Loop:        loop,
		Invoker:     invokerAdapter(invoker),
		Reporter:    func(msg string) {},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if code != ExitSuccess {
		t.Errorf("exit code = %d, want ExitSuccess (FAIL, null, FAIL should not hit threshold 2)", code)
	}
	if callCount != 4 {
		t.Errorf("invoker called %d times, want 4", callCount)
	}
}

// TestRun_NoSignalBelowThreshold_Continues verifies null then success completes.
func TestRun_NoSignalBelowThreshold_Continues(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 5
	loop.FailureThreshold = 2
	loop.SuccessSignal = "<promise>SUCCESS</promise>"
	loop.FailureSignal = "<promise>FAILURE</promise>"
	callCount := 0
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
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
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
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
	if !strings.Contains(reported, "invocation error") {
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
		Command:     "true",
		PromptBytes: []byte("p"),
		Loop:        loop,
		Invoker: invokerAdapter(func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
			return []byte("x"), 0, nil
		}),
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
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
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
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
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

// TestRun_SuccessSignalOnlyOnLastLine verifies R004: success is detected only when
// the signal appears on the last non-empty line. When the signal is only on an
// earlier line, the iteration is not successful.
func TestRun_SuccessSignalOnlyOnLastLine(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 2
	loop.SuccessSignal = "DONE"
	// Output: "DONE" on first line, "Still working..." on last line -> last non-empty line has no signal.
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
		return []byte("DONE\nStill working..."), 0, nil
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
	// Should exit with max iterations (no success detected; signal was not on last line).
	if code != ExitMaxIterations {
		t.Errorf("exit code = %d, want %d (ExitMaxIterations); success signal only on earlier line must not count", code, ExitMaxIterations)
	}
	all := strings.Join(reported, " ")
	if !strings.Contains(all, "without success signal") || !strings.Contains(all, "max: 2") {
		t.Errorf("reported = %q", reported)
	}
}

// TestRun_SuccessSignalOnLastLine_detected verifies R004: when the success signal
// is on the last non-empty line, the run completes successfully.
func TestRun_SuccessSignalOnLastLine_detected(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 2
	loop.SuccessSignal = "DONE"
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
		return []byte("Still working...\nStatus: DONE"), 0, nil
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
		t.Errorf("exit code = %d, want ExitSuccess (signal on last line)", code)
	}
	if !strings.Contains(reported, "Completed successfully") || !strings.Contains(reported, "1 iteration") {
		t.Errorf("reported = %q", reported)
	}
}

// TestRun_FailureSignalOnLastLine_detected verifies R005: when the failure signal
// appears on the last non-empty line, the iteration is treated as failure.
func TestRun_FailureSignalOnLastLine_detected(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 2
	loop.FailureThreshold = 1
	loop.SuccessSignal = "DONE"
	loop.FailureSignal = "FAIL"
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
		return []byte("Working...\nFAIL"), 0, nil
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
		t.Errorf("exit code = %d, want ExitFailureThreshold (failure on last line)", code)
	}
	all := strings.Join(reported, " ")
	if !strings.Contains(all, "consecutive failure(s)") {
		t.Errorf("reported = %q", reported)
	}
}

// TestRun_StaticPrecedence_BothSignalsOnLastLine verifies R006: when both success
// and failure signals appear on the same (last non-empty) line, static precedence
// applies (success wins).
func TestRun_StaticPrecedence_BothSignalsOnLastLine(t *testing.T) {
	loop := config.DefaultLoopSettings()
	loop.MaxIterations = 3
	loop.SuccessSignal = "DONE"
	loop.FailureSignal = "FAIL"
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
		// Both on last non-empty line; success must win.
		return []byte("line one\nFAIL and DONE together"), 0, nil
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
		t.Errorf("exit code = %d, want ExitSuccess (static precedence when both on last line)", code)
	}
	if !strings.Contains(reported, "Completed successfully") || !strings.Contains(reported, "1 iteration") {
		t.Errorf("reported = %q", reported)
	}
}

// TestRun_InterruptBetweenIterations verifies that if the interrupt context is
// cancelled after the first iteration completes, Run returns ExitInterrupt
// (checked at start of next iteration and after Invoke returns).
func TestRun_InterruptBetweenIterations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	callCount := 0
	invoker := func(_ string, _ []byte, _ string, _ []string, _ int, _ int, _ io.Writer) ([]byte, int, error) {
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
