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
	_, code, err := Invoke("", []byte("hi"), "", nil, 0, 0, nil)
	if err != ErrEmptyCommand {
		t.Fatalf("Invoke(empty): err = %v, want ErrEmptyCommand", err)
	}
	if code != -1 {
		t.Errorf("exitCode = %d, want -1", code)
	}
}

func TestInvoke_WhitespaceCommand(t *testing.T) {
	_, code, err := Invoke("   \t  ", []byte("hi"), "", nil, 0, 0, nil)
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
	stdout, code, err := Invoke(cmd, input, "", nil, 0, 0, nil)
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
		stdout, code, err := Invoke("pwd", nil, dir, nil, 0, 0, nil)
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
	stdout, code, err := Invoke("pwd", nil, "", nil, 0, 0, nil)
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
	stdout, code, err := Invoke(cmd, nil, "", nil, 0, 0, nil)
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

// TestInvoke_Timeout verifies that when timeoutSec > 0, the process is killed
// after that many seconds and ErrTimeout is returned (T2.4).
func TestInvoke_Timeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		// sleep 10 may not exist or behave the same on Windows
		t.Skip("timeout test uses sleep (Unix)")
	}
	if path, _ := exec.LookPath("sleep"); path == "" {
		t.Skip("sleep not available")
	}
	// sleep 10 should be killed after 1 second
	stdout, code, err := Invoke("sleep 10", nil, "", nil, 1, 0, nil)
	if err != ErrTimeout {
		t.Fatalf("Invoke(sleep 10, timeout 1s): err = %v, want ErrTimeout", err)
	}
	if code != -1 {
		t.Errorf("exitCode = %d, want -1 on timeout", code)
	}
	if len(stdout) != 0 {
		t.Errorf("stdout = %q, want empty on timeout", stdout)
	}
}

// TestInvoke_NoTimeout verifies that timeoutSec 0 means no timeout (process runs to completion).
func TestInvoke_NoTimeout(t *testing.T) {
	cmd := "true"
	if runtime.GOOS == "windows" {
		cmd = "cmd /c exit 0"
	}
	if path, _ := exec.LookPath("true"); path == "" && runtime.GOOS != "windows" {
		cmd, _ = exec.LookPath("echo")
		if cmd == "" {
			t.Skip("no true/echo available")
		}
	}
	_, code, err := Invoke(cmd, nil, "", nil, 0, 0, nil)
	if err != nil {
		t.Fatalf("Invoke(no timeout): %v", err)
	}
	if code != 0 {
		t.Errorf("exitCode = %d, want 0", code)
	}
}

// TestInvoke_StreamTo verifies O004/R006: when streamTo is non-nil, stdout is
// both streamed to that writer and captured in the returned buffer.
func TestInvoke_StreamTo(t *testing.T) {
	cmd := "cat"
	if runtime.GOOS == "windows" {
		cmd = "type"
	}
	if path, _ := exec.LookPath("cat"); path != "" {
		cmd = "cat"
	} else if path, _ := exec.LookPath("type"); path != "" && runtime.GOOS == "windows" {
		cmd = "type"
	} else {
		t.Skip("cat/type not available")
	}
	input := []byte("streamed output")
	var streamBuf bytes.Buffer
	stdout, code, err := Invoke(cmd, input, "", nil, 0, 0, &streamBuf)
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	if code != 0 {
		t.Errorf("exitCode = %d, want 0", code)
	}
	if string(stdout) != string(input) {
		t.Errorf("returned stdout = %q, want %q", stdout, input)
	}
	if streamBuf.String() != string(input) {
		t.Errorf("streamTo buffer = %q, want %q", streamBuf.String(), input)
	}
}

// TestInvoke_MaxOutputBytes_capsReturnedStdout verifies that when maxOutputBytes > 0,
// only the last maxOutputBytes bytes are returned (sliding window).
func TestInvoke_MaxOutputBytes_capsReturnedStdout(t *testing.T) {
	cmd := "cat"
	if runtime.GOOS == "windows" {
		cmd = "type"
	}
	if path, _ := exec.LookPath("cat"); path != "" {
		cmd = "cat"
	} else if path, _ := exec.LookPath("type"); path != "" && runtime.GOOS == "windows" {
		cmd = "type"
	} else {
		t.Skip("cat/type not available")
	}
	// Output 200 bytes; cap at 10.
	input := bytes.Repeat([]byte("x"), 200)
	const cap = 10
	stdout, code, err := Invoke(cmd, input, "", nil, 0, cap, nil)
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	if code != 0 {
		t.Errorf("exitCode = %d, want 0", code)
	}
	if len(stdout) != cap {
		t.Errorf("len(stdout) = %d, want %d", len(stdout), cap)
	}
	want := input[len(input)-cap:]
	if !bytes.Equal(stdout, want) {
		t.Errorf("stdout = %q, want last %d bytes %q", stdout, cap, want)
	}
}

// TestInvoke_MaxOutputBytes_streaming_capsReturnedStdout verifies that when streaming
// and maxOutputBytes > 0, streamTo receives full output but returned buffer is capped.
func TestInvoke_MaxOutputBytes_streaming_capsReturnedStdout(t *testing.T) {
	cmd := "cat"
	if runtime.GOOS == "windows" {
		cmd = "type"
	}
	if path, _ := exec.LookPath("cat"); path != "" {
		cmd = "cat"
	} else if path, _ := exec.LookPath("type"); path != "" && runtime.GOOS == "windows" {
		cmd = "type"
	} else {
		t.Skip("cat/type not available")
	}
	input := bytes.Repeat([]byte("a"), 100)
	var streamBuf bytes.Buffer
	const cap = 20
	stdout, code, err := Invoke(cmd, input, "", nil, 0, cap, &streamBuf)
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	if code != 0 {
		t.Errorf("exitCode = %d, want 0", code)
	}
	if streamBuf.Len() != len(input) {
		t.Errorf("streamTo received %d bytes, want %d", streamBuf.Len(), len(input))
	}
	if !bytes.Equal(streamBuf.Bytes(), input) {
		t.Errorf("streamTo = %q, want full input", streamBuf.Bytes())
	}
	if len(stdout) != cap {
		t.Errorf("len(returned stdout) = %d, want %d", len(stdout), cap)
	}
	want := input[len(input)-cap:]
	if !bytes.Equal(stdout, want) {
		t.Errorf("returned stdout = %q, want last %d bytes %q", stdout, cap, want)
	}
}

// TestInvoke_MaxOutputBytes_zeroUnlimited verifies that maxOutputBytes <= 0 means no cap.
func TestInvoke_MaxOutputBytes_zeroUnlimited(t *testing.T) {
	cmd := "cat"
	if runtime.GOOS == "windows" {
		cmd = "type"
	}
	if path, _ := exec.LookPath("cat"); path != "" {
		cmd = "cat"
	} else if path, _ := exec.LookPath("type"); path != "" && runtime.GOOS == "windows" {
		cmd = "type"
	} else {
		t.Skip("cat/type not available")
	}
	input := bytes.Repeat([]byte("y"), 500)
	stdout, code, err := Invoke(cmd, input, "", nil, 0, 0, nil)
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	if code != 0 {
		t.Errorf("exitCode = %d, want 0", code)
	}
	if len(stdout) != len(input) {
		t.Errorf("len(stdout) = %d, want %d (no cap)", len(stdout), len(input))
	}
}
