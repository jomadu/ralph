// Package backend invokes the AI CLI with prompt on stdin and captures stdout.
// No shell; exec-style invocation. Implements O003/R001 and backend Interfaces.
package backend

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ErrEmptyCommand is returned when the command string is empty or whitespace.
var ErrEmptyCommand = errors.New("backend: empty command")

// ErrTimeout is returned when the invocation exceeds the configured per-iteration timeout (T2.4).
var ErrTimeout = errors.New("backend: per-iteration timeout exceeded")

// Invoker runs the resolved AI command with the given prompt on stdin and
// returns full stdout, exit code, and any start/exec error.
// streamTo: when non-nil, stdout is also written to this writer in real time (O004/R006).
// timeoutSec: 0 = no timeout; >0 = kill process after that many seconds (T2.4).
type Invoker interface {
	Invoke(command string, promptBytes []byte, cwd string, env []string, timeoutSec int, streamTo io.Writer) (stdout []byte, exitCode int, err error)
}

// Invoke runs the resolved AI command (executable + args, no shell) with
// promptBytes written to stdin (stream closed after write), cwd and env
// applied. If streamTo is non-nil, child stdout is streamed to it while
// still captured and returned (O004/R006). If timeoutSec > 0, the process is
// killed after timeoutSec seconds and ErrTimeout is returned (T2.4).
// Returns full stdout, process exit code, and error if the process could not
// be started or timed out. Implements O003/R001, O004/R006.
func Invoke(command string, promptBytes []byte, cwd string, env []string, timeoutSec int, streamTo io.Writer) (stdout []byte, exitCode int, err error) {
	argv := splitCommand(command)
	if len(argv) == 0 {
		return nil, -1, ErrEmptyCommand
	}
	ctx := context.Background()
	var cancel context.CancelFunc
	if timeoutSec > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
		defer cancel()
	}
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Stdin = bytes.NewReader(promptBytes)
	// Stderr left nil so it inherits (implementation-defined per backend.md)
	if cwd != "" {
		cmd.Dir = cwd
	}
	if len(env) > 0 {
		cmd.Env = env
	} else {
		cmd.Env = os.Environ()
	}

	var out bytes.Buffer
	if streamTo != nil {
		pipe, err := cmd.StdoutPipe()
		if err != nil {
			return nil, -1, err
		}
		if err := cmd.Start(); err != nil {
			return nil, -1, err
		}
		tee := io.MultiWriter(streamTo, &out)
		_, _ = io.Copy(tee, pipe)
		runErr := cmd.Wait()
		stdout = out.Bytes()
		if runErr != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return stdout, -1, ErrTimeout
			}
			if exitErr, ok := runErr.(*exec.ExitError); ok {
				return stdout, exitErr.ExitCode(), nil
			}
			return stdout, -1, runErr
		}
		return stdout, 0, nil
	}

	cmd.Stdout = &out
	runErr := cmd.Run()
	stdout = out.Bytes()
	if runErr != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return stdout, -1, ErrTimeout
		}
		if exitErr, ok := runErr.(*exec.ExitError); ok {
			return stdout, exitErr.ExitCode(), nil
		}
		return stdout, -1, runErr
	}
	return stdout, 0, nil
}

// splitCommand splits a resolved command string into executable and arguments
// for exec-style invocation (no shell). Uses fields split; arguments
// containing spaces must be passed via a wrapper script.
func splitCommand(command string) []string {
	return strings.Fields(strings.TrimSpace(command))
}
