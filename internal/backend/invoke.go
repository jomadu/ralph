// Package backend invokes the AI CLI with prompt on stdin and captures stdout.
// No shell; exec-style invocation. Implements O003/R001 and backend Interfaces.
package backend

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
)

// ErrEmptyCommand is returned when the command string is empty or whitespace.
var ErrEmptyCommand = errors.New("backend: empty command")

// Invoker runs the resolved AI command with the given prompt on stdin and
// returns full stdout, exit code, and any start/exec error.
type Invoker interface {
	Invoke(command string, promptBytes []byte, cwd string, env []string) (stdout []byte, exitCode int, err error)
}

// Invoke runs the resolved AI command (executable + args, no shell) with
// promptBytes written to stdin (stream closed after write), cwd and env
// applied. Returns full stdout, process exit code, and error if the process
// could not be started. Implements backend Interfaces and O003/R001.
func Invoke(command string, promptBytes []byte, cwd string, env []string) (stdout []byte, exitCode int, err error) {
	argv := splitCommand(command)
	if len(argv) == 0 {
		return nil, -1, ErrEmptyCommand
	}
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Stdin = bytes.NewReader(promptBytes)
	var out bytes.Buffer
	cmd.Stdout = &out
	// Stderr left nil so it inherits (implementation-defined per backend.md)
	if cwd != "" {
		cmd.Dir = cwd
	}
	if len(env) > 0 {
		cmd.Env = env
	} else {
		cmd.Env = os.Environ()
	}
	runErr := cmd.Run()
	stdout = out.Bytes()
	if runErr != nil {
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
