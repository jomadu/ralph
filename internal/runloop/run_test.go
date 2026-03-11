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
