package runner

import (
	"io"
	"os"
	"os/exec"
)

// SpawnAI spawns the AI CLI process with inherited environment and working directory.
// Per O3/R4: no filtering of env, cwd is Ralph's cwd.
// Returns exit code (0 for success, non-zero for crash/error) and error.
func SpawnAI(argv []string, stdin io.Reader, stdout, stderr io.Writer) (int, error) {
	if len(argv) == 0 {
		return 0, nil
	}

	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Env = os.Environ()
	cmd.Dir, _ = os.Getwd()
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return exitCode, err
}
