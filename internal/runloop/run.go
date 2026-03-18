package runloop

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/jomadu/ralph/internal/backend"
	"github.com/jomadu/ralph/internal/config"
)

// RunOptions supplies the inputs for one run-loop execution. Caller must resolve
// config and prompt; run-loop validates command, runs the loop, and returns exit code.
type RunOptions struct {
	Command     string
	PromptBytes []byte
	Loop        config.LoopSettings
	Cwd         string
	Env         []string
	Invoker     backend.Invoker
	// DryRun when true: assemble prompt (with preamble if enabled), print to stdout, exit 0.
	// No backend invocation (T3.10, O004/R007).
	DryRun bool
	// Reporter receives completion message on success; nil = print to os.Stdout.
	Reporter func(msg string)
	// StreamWriter when non-nil and Loop.Streaming is true, AI stdout is streamed here (O004/R006).
	StreamWriter io.Writer
	// InterruptContext if non-nil is used for interrupt detection (e.g. in tests);
	// when it is cancelled, Run returns ExitInterrupt. If nil, Run uses
	// signal.NotifyContext(background, os.Interrupt).
	InterruptContext context.Context
}

// invokerAdapter adapts a function with optional streamTo to backend.Invoker.
type invokerAdapter func(command string, promptBytes []byte, cwd string, env []string, timeoutSec int, maxOutputBytes int, streamTo io.Writer) (stdout []byte, exitCode int, err error)

func (f invokerAdapter) Invoke(command string, promptBytes []byte, cwd string, env []string, timeoutSec int, maxOutputBytes int, streamTo io.Writer) ([]byte, int, error) {
	return f(command, promptBytes, cwd, env, timeoutSec, maxOutputBytes, streamTo)
}

// logLevelPriority returns a numeric priority for level comparison (higher = more verbose).
// Empty or unknown level is treated as "info". Used for O004/R006.
func logLevelPriority(level string) int {
	switch level {
	case "error":
		return 0
	case "warn":
		return 1
	case "info", "":
		return 2
	case "debug":
		return 3
	default:
		return 2
	}
}

// reportLevel emits msg only when messageLevel is at or above the configured log level (O004/R006).
func reportLevel(report func(string), configuredLevel, messageLevel, msg string) {
	if logLevelPriority(messageLevel) <= logLevelPriority(configuredLevel) {
		report(msg)
	}
}

// sectionHeader returns a single-line delimiter before each titled block (LOOP CONFIG,
// CONTEXT, INSTRUCTIONS). The line must not start with "---": some AI CLIs (e.g. Cursor
// agent) treat that as YAML frontmatter and exit without processing the prompt.
// Current format: # --- NAME --- (hash first; dashes are decorative).
func sectionHeader(name string) string {
	return "# --- " + name + " ---\n"
}

// ralphLoopDescription is the brief description of the Ralph loop technique injected when preamble is enabled.
const ralphLoopDescription = "You are in a Ralph loop. This prompt might be run multiple times until completion criteria is met or a consecutive failure limit is reached. When completion criteria are met, emit a success signal; when you cannot meet them, emit a failure signal; when more work remains, emit no signal so that the loop continues. When you do emit success or failure signals, put that signal on the last line of your output (Ralph only scans the last non-empty line)."

// unlimitedIterationsThreshold: max iterations at or above this are shown as "unlimited" in the context section.
const unlimitedIterationsThreshold = 1_000_000_000

// invokerContextLabel is the line that introduces user-provided context (-c) inside the CONTEXT section.
const invokerContextLabel = "Context provided by the invoker of this Ralph run:"

// formatLoopConfig returns a human-readable summary of loop settings for dry-run output.
func formatLoopConfig(loop config.LoopSettings) string {
	timeout := strconv.Itoa(loop.TimeoutSeconds)
	if loop.TimeoutSeconds == 0 {
		timeout = "0 (no limit)"
	}
	return fmt.Sprintf("max_iterations: %d\nfailure_threshold: %d\ntimeout_seconds: %s\nsuccess_signal: %q\nfailure_signal: %q\npreamble: %t\nlog_level: %s\nstreaming: %t\nmax_output_buffer: %d",
		loop.MaxIterations, loop.FailureThreshold, timeout,
		loop.SuccessSignal, loop.FailureSignal,
		loop.Preamble, loop.LogLevel, loop.Streaming, loop.MaxOutputBuffer)
}

// Run validates the AI command, then runs the loop: for each iteration invokes
// the backend with the assembled prompt, captures stdout, and scans for the
// configured success and failure signals. On success: reports completion and
// returns ExitSuccess. On failure signal, non-zero process exit, or invocation
// error: increments consecutive-failure count; if count >= failure threshold,
// reports and returns ExitFailureThreshold. When the process exits 0 and the
// last non-empty line has neither success nor failure signal (“null signal”),
// the iteration does not count toward the threshold; the failure streak resets
// and the loop continues (aligned with the preamble: more work remains).
// When max iterations is reached without success, returns ExitMaxIterations.
// On SIGINT or SIGTERM (T3.9, O004/R005): returns ExitInterrupt (130). Static
// precedence (T3.6, O001/R006): success is checked before failure. Log level
// and streaming are respected (T3.13, O004/R006).
func Run(opts RunOptions) (exitCode int, err error) {
	if opts.Invoker == nil {
		opts.Invoker = invokerAdapter(backend.Invoke)
	}
	report := opts.Reporter
	if report == nil {
		report = func(msg string) { fmt.Fprintln(os.Stdout, msg) }
	}
	logLevel := opts.Loop.LogLevel
	if logLevel == "" {
		logLevel = "info"
	}
	// streamTo: when Streaming and StreamWriter set, backend will tee stdout here (O004/R006).
	var streamTo io.Writer
	if opts.Loop.Streaming && opts.StreamWriter != nil {
		streamTo = opts.StreamWriter
	}

	// Dry-run: print LOOP CONFIG, then assembled prompt (CONTEXT + INSTRUCTIONS sections), exit 0. No backend (T3.10, O004/R007).
	if opts.DryRun {
		configSection := sectionHeader("LOOP CONFIG") + "\n\n" + formatLoopConfig(opts.Loop)
		contextBody := buildContextBody(opts.Loop.Preamble, 1, opts.Loop.MaxIterations, opts.Loop.Context)
		assembled := assembleWithSectionHeaders(contextBody, opts.PromptBytes)
		os.Stdout.Write([]byte(configSection + "\n\n"))
		os.Stdout.Write(assembled)
		os.Stdout.Write([]byte("\n"))
		reportLevel(report, logLevel, "info", "Dry-run: loop config and assembled prompt printed; no run was performed.")
		return ExitSuccess, nil
	}

	if err := ValidateAICommand(opts.Command); err != nil {
		return ExitErrorPreLoop, err
	}

	// Interrupt: use optional context (e.g. for tests) or os.Interrupt (O004/R005).
	var ctx context.Context
	var stop context.CancelFunc
	if opts.InterruptContext != nil {
		ctx = opts.InterruptContext
		stop = func() {}
	} else {
		ctx, stop = signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()
	}

	start := time.Now()
	consecutiveFailures := 0
	threshold := opts.Loop.FailureThreshold
	if threshold <= 0 {
		threshold = 1
	}
	var iterationDurations []time.Duration
	for i := 1; i <= opts.Loop.MaxIterations; i++ {
		select {
		case <-ctx.Done():
			return ExitInterrupt, nil
		default:
		}
		reportLevel(report, logLevel, "debug", fmt.Sprintf("Starting iteration %d.", i))
		contextBody := buildContextBody(opts.Loop.Preamble, i, opts.Loop.MaxIterations, opts.Loop.Context)
		assembled := assembleWithSectionHeaders(contextBody, opts.PromptBytes)
		iterStart := time.Now()
		stdout, procExit, invErr := opts.Invoker.Invoke(opts.Command, assembled, opts.Cwd, opts.Env, opts.Loop.TimeoutSeconds, opts.Loop.MaxOutputBuffer, streamTo)
		iterationDurations = append(iterationDurations, time.Since(iterStart))
		if ctx.Err() != nil {
			return ExitInterrupt, nil
		}
		// Invocation error (e.g. timeout, crash, exec failure): counts toward failure threshold.
		if invErr != nil {
			consecutiveFailures++
			if consecutiveFailures >= threshold {
				reportLevel(report, logLevel, "error", fmt.Sprintf("Stopped after %d consecutive iteration(s) with invocation error (last: %v; threshold: %d).", consecutiveFailures, invErr, opts.Loop.FailureThreshold))
				return ExitFailureThreshold, nil
			}
			continue
		}
		lastLine := LastNonEmptyLine(stdout)
		hasSuccess := ContainsSuccessSignal(lastLine, opts.Loop.SuccessSignal)
		hasFailure := ContainsFailureSignal(lastLine, opts.Loop.FailureSignal)
		// Static precedence (O001/R006): when both signals present on last line, success wins.
		if hasSuccess {
			elapsed := time.Since(start)
			report(completionMessage(i, elapsed))
			reportIterationStatsLevel(report, logLevel, iterationDurations)
			return ExitSuccess, nil
		}
		if hasFailure {
			consecutiveFailures++
			if consecutiveFailures >= threshold {
				reportLevel(report, logLevel, "error", fmt.Sprintf("Stopped after %d consecutive failure(s) (threshold: %d).", consecutiveFailures, opts.Loop.FailureThreshold))
				reportIterationStatsLevel(report, logLevel, iterationDurations)
				return ExitFailureThreshold, nil
			}
			continue
		}
		if procExit != 0 {
			consecutiveFailures++
			if consecutiveFailures >= threshold {
				reportLevel(report, logLevel, "error", fmt.Sprintf("Stopped after %d consecutive iteration(s) where the AI process exited with code %d without success signal (threshold: %d).", consecutiveFailures, procExit, opts.Loop.FailureThreshold))
				reportIterationStatsLevel(report, logLevel, iterationDurations)
				return ExitFailureThreshold, nil
			}
			continue
		}
		// Exit 0, no success/failure on last line: neutral iteration (more work remains).
		consecutiveFailures = 0
		reportLevel(report, logLevel, "debug", fmt.Sprintf("Iteration %d: no success or failure signal on last line; continuing (failure streak reset).", i))
	}
	reportLevel(report, logLevel, "error", fmt.Sprintf("Stopped after %d iteration(s) without success signal (max: %d).", opts.Loop.MaxIterations, opts.Loop.MaxIterations))
	reportIterationStatsLevel(report, logLevel, iterationDurations)
	return ExitMaxIterations, nil
}

// reportIterationStats reports min/max/mean duration per iteration when there are
// two or more iterations (T3.12, O004/R008). Single-iteration timing is already
// in the completion message; multi-iteration runs get statistics for tuning.
func reportIterationStats(report func(string), durations []time.Duration) {
	if len(durations) < 2 {
		return
	}
	var total time.Duration
	minD, maxD := durations[0], durations[0]
	for _, d := range durations {
		total += d
		if d < minD {
			minD = d
		}
		if d > maxD {
			maxD = d
		}
	}
	meanD := total / time.Duration(len(durations))
	report(fmt.Sprintf("Iteration stats: min %.2fs, max %.2fs, mean %.2fs (%d iterations).", minD.Seconds(), maxD.Seconds(), meanD.Seconds(), len(durations)))
}

// reportIterationStatsLevel emits iteration stats only when log level allows info (O004/R006).
func reportIterationStatsLevel(report func(string), logLevel string, durations []time.Duration) {
	if logLevelPriority("info") > logLevelPriority(logLevel) {
		return
	}
	reportIterationStats(report, durations)
}

// buildContextBody returns the body of the single CONTEXT section: when preamble is enabled,
// the Ralph loop description and iteration line; optionally, invoker-provided context (-c)
// with an explicit label. invokerContext is the raw text from the invoker (no "CONTEXT" prefix).
func buildContextBody(injectPreamble bool, iteration, maxIterations int, invokerContext string) string {
	var parts []string
	if injectPreamble {
		iterLine := "Iteration " + strconv.Itoa(iteration)
		if maxIterations >= unlimitedIterationsThreshold {
			iterLine += " (unlimited)"
		} else {
			iterLine += " of max " + strconv.Itoa(maxIterations)
		}
		parts = append(parts, ralphLoopDescription+"\n"+iterLine)
	}
	if invokerContext != "" {
		parts = append(parts, invokerContextLabel+"\n"+invokerContext)
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, "\n\n")
}

// assembleWithSectionHeaders builds the full prompt with titled section separators (# --- NAME ---).
// Single CONTEXT section (when non-empty) then INSTRUCTIONS. Sections are separated by a blank line.
func assembleWithSectionHeaders(contextBody string, promptBytes []byte) []byte {
	var parts []string
	if contextBody != "" {
		parts = append(parts, sectionHeader("CONTEXT")+"\n\n"+contextBody)
	}
	parts = append(parts, sectionHeader("INSTRUCTIONS")+"\n\n"+string(promptBytes))
	return []byte(strings.Join(parts, "\n\n"))
}

func completionMessage(iterations int, elapsed time.Duration) string {
	sec := elapsed.Seconds()
	if iterations == 1 {
		return fmt.Sprintf("Completed successfully in 1 iteration (%.2fs).", sec)
	}
	return fmt.Sprintf("Completed successfully in %d iterations (%.2fs).", iterations, sec)
}
