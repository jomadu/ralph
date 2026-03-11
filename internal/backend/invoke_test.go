package backend

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestInvoke_EmptyCommand(t *testing.T) {
	_, code, err := Invoke("", []byte("hi"), "", nil)
	if err != ErrEmptyCommand {
		t.Fatalf("Invoke(empty): err = %v, want ErrEmptyCommand", err)
	}
	if code != -1 {
		t.Errorf("exitCode = %d, want -1", code)
	}
}

func TestInvoke_WhitespaceCommand(t *testing.T) {
	_, code, err := Invoke("   \t  ", []byte("hi"), "", nil)
	if err != ErrEmptyCommand {
		t.Fatalf("Invoke(whitespace): err = %v, want ErrEmptyCommand", err)
	}
	if code != -1 {
		t.Errorf("exitCode = %d, want -1", code)
	}
}

func TestInvoke_EchoStdin(t *testing.T) {
	// Use a command that reads stdin and exits 0 (e.g. cat or true).
	// On Unix: "cat" echoes stdin; "true" ignores it and exits 0.
	cmd := "cat"
	if runtime.GOOS == "windows" {
		cmd = "type"
	}
	// Prefer cat if available so we can assert on stdout.
	if path, _ := exec.LookPath("cat"); path != "" {
		cmd = "cat"
	} else if path, _ := exec.LookPath("type"); path != "" && runtime.GOOS == "windows" {
		cmd = "type"
	} else {
		// Just run something that exits 0 with stdin closed
		cmd, _ = exec.LookPath("true")
		if cmd == "" {
			cmd, _ = exec.LookPath("echo")
		}
		if cmd == "" {
			t.Skip("no cat/true/echo available")
		}
	}
	input := []byte("hello stdin")
	stdout, code, err := Invoke(cmd, input, "", nil)
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	if code != 0 {
		t.Errorf("exitCode = %d, want 0", code)
	}
	if path, _ := exec.LookPath("cat"); path != "" && cmd == "cat" {
		if string(stdout) != string(input) {
			t.Errorf("stdout = %q, want %q", stdout, input)
		}
	}
}

func TestInvoke_Cwd(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("cwd test uses pwd (Unix)")
	}
	if path, _ := exec.LookPath("pwd"); path != "" {
		dir, _ := os.Getwd()
		stdout, code, err := Invoke("pwd", nil, dir, nil)
		if err != nil {
			t.Fatalf("Invoke: %v", err)
		}
		if code != 0 {
			t.Errorf("exitCode = %d, want 0", code)
		}
		got := string(bytes.TrimSpace(stdout))
		abs, _ := filepath.Abs(dir)
		if got != abs && got != dir {
			t.Errorf("cwd: got %q, want %q or %q", got, dir, abs)
		}
	} else {
		t.Skip("pwd not available")
	}
}

// TestInvoke_CwdInherit verifies that when cwd is empty, the child inherits the parent's working directory (O003/R002).
func TestInvoke_CwdInherit(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("cwd test uses pwd (Unix)")
	}
	if path, _ := exec.LookPath("pwd"); path == "" {
		t.Skip("pwd not available")
	}
	want, _ := os.Getwd()
	stdout, code, err := Invoke("pwd", nil, "", nil)
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	if code != 0 {
		t.Errorf("exitCode = %d, want 0", code)
	}
	got := string(bytes.TrimSpace(stdout))
	absWant, _ := filepath.Abs(want)
	if got != absWant && got != want {
		t.Errorf("cwd inherit: got %q, want %q or %q", got, want, absWant)
	}
}

// TestInvoke_EnvInherit verifies that when env is nil, the child inherits the parent's environment (O003/R002).
func TestInvoke_EnvInherit(t *testing.T) {
	const testVar = "RALPH_TEST_ENV_INHERIT"
	const testVal = "inherited"
	t.Setenv(testVar, testVal)
	// printenv VAR prints the value and exits 0; not on Windows, so skip there.
	var cmd string
	if runtime.GOOS != "windows" {
		if path, _ := exec.LookPath("printenv"); path != "" {
			cmd = "printenv " + testVar
		}
	}
	if cmd == "" {
		t.Skip("printenv not available (Unix only)")
	}
	stdout, code, err := Invoke(cmd, nil, "", nil)
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	if code != 0 {
		t.Errorf("exitCode = %d, want 0", code)
	}
	got := string(bytes.TrimSpace(stdout))
	if got != testVal {
		t.Errorf("env inherit: got %q, want %q", got, testVal)
	}
}
