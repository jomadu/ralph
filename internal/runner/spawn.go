package runner

import (
	"io"
	"os"
	"os/exec"
)

// SpawnAI spawns the AI CLI process with inherited environment and working directory.
// Per O3/R4: no filtering of env, cwd is Ralph's cwd.
func SpawnAI(argv []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(argv) == 0 {
		return nil
	}

	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Env = os.Environ()
	cmd.Dir, _ = os.Getwd()
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}
