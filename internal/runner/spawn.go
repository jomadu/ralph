package runner

import (
	"context"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// SpawnAI spawns the AI CLI process with inherited environment and working directory.
// Per O3/R4: no filtering of env, cwd is Ralph's cwd.
// Returns exit code (0 for success, non-zero for crash/error) and error.
func SpawnAI(argv []string, stdin io.Reader, stdout, stderr io.Writer) (int, error) {
	return SpawnAIWithContext(context.Background(), argv, stdin, stdout, stderr)
}

// SpawnAIWithContext spawns the AI CLI process with context for cancellation.
// On context cancellation: sends SIGTERM, waits 5s, then SIGKILL if needed.
// Returns exit code and error. Exit code 130 indicates interruption.
func SpawnAIWithContext(ctx context.Context, argv []string, stdin io.Reader, stdout, stderr io.Writer) (int, error) {
	if len(argv) == 0 {
		return 0, nil
	}

	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Env = os.Environ()
	cmd.Dir, _ = os.Getwd()
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return 1, err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// Context cancelled - graceful shutdown with 5s grace period (O1/R7)
		_ = cmd.Process.Signal(syscall.SIGTERM)

		gracePeriod := time.NewTimer(5 * time.Second)
		defer gracePeriod.Stop()

		select {
		case <-done:
			// Process exited during grace period
			return 130, nil
		case <-gracePeriod.C:
			// Grace period expired - force kill
			_ = cmd.Process.Kill()
			<-done
			return 130, nil
		}

	case err := <-done:
		// Process completed normally
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
}
