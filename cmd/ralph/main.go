package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/maxdunn/ralph/internal/cmdparse"
	"github.com/maxdunn/ralph/internal/config"
	"github.com/maxdunn/ralph/internal/logger"
	"github.com/maxdunn/ralph/internal/prompt"
	"github.com/maxdunn/ralph/internal/review"
	"github.com/maxdunn/ralph/internal/runner"
)

var rootCmd = &cobra.Command{
	Use:   "ralph",
	Short: "Ralph - a dumb loop that pipes prompts to AI CLIs",
}

var (
	// Configuration
	configFlag string

	// Loop control
	maxIterationsFlag   int
	unlimitedFlag       bool
	failureThresholdFlag int
	iterationTimeoutFlag int
	maxOutputBufferFlag int
	preambleFlag        bool
	noPreambleFlag      bool
	dryRunFlag          bool

	// AI command
	aiCmdFlag      string
	aiCmdAliasFlag string

	// Signals
	signalSuccessFlag string
	signalFailureFlag string

	// Context
	contextFlags []string

	// Output control
	verboseFlag  bool
	quietFlag    bool
	logLevelFlag string

	// Prompt input
	fileFlag string

	// Review-specific flags (T1; T2 adds --review-output; T7 adds --prompt-output, --apply for validation)
	reviewFileFlag         string
	reviewOutputFlag       string
	reviewPromptOutputFlag string
	reviewApplyFlag        bool
)

var runCmd = &cobra.Command{
	Use:   "run [alias]",
	Short: "Run the loop with a prompt",
	Run: func(cmd *cobra.Command, args []string) {
		// Load config with explicit path if provided
		cfg, err := config.LoadConfigWithProvenanceAndExplicit(configFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		// Build CLI flags struct
		cliFlags := config.CLIFlags{}
		if cmd.Flags().Changed("max-iterations") {
			cliFlags.MaxIterations = &maxIterationsFlag
		}
		if cmd.Flags().Changed("unlimited") {
			unlimited := "unlimited"
			cliFlags.IterationMode = &unlimited
		}
		if cmd.Flags().Changed("failure-threshold") {
			cliFlags.FailureThreshold = &failureThresholdFlag
		}
		if cmd.Flags().Changed("iteration-timeout") {
			cliFlags.IterationTimeout = &iterationTimeoutFlag
		}
		if cmd.Flags().Changed("max-output-buffer") {
			cliFlags.MaxOutputBuffer = &maxOutputBufferFlag
		}
		if cmd.Flags().Changed("ai-cmd") {
			cliFlags.AICmd = &aiCmdFlag
		}
		if cmd.Flags().Changed("ai-cmd-alias") {
			cliFlags.AICmdAlias = &aiCmdAliasFlag
		}
		if cmd.Flags().Changed("signal-success") {
			cliFlags.SignalSuccess = &signalSuccessFlag
		}
		if cmd.Flags().Changed("signal-failure") {
			cliFlags.SignalFailure = &signalFailureFlag
		}
		if cmd.Flags().Changed("verbose") {
			cliFlags.ShowAIOutput = &verboseFlag
		}
		if cmd.Flags().Changed("quiet") {
			quietLevel := "error"
			cliFlags.LogLevel = &quietLevel
		}
		if cmd.Flags().Changed("log-level") {
			cliFlags.LogLevel = &logLevelFlag
		}

		// Handle preamble/no-preamble flags
		if cmd.Flags().Changed("preamble") {
			cliFlags.Preamble = &preambleFlag
		} else if cmd.Flags().Changed("no-preamble") {
			noPreamble := !noPreambleFlag
			cliFlags.Preamble = &noPreamble
		}

		// Resolve prompt input mode
		var alias string
		if len(args) > 0 {
			alias = args[0]
		}

		mode, err := prompt.ResolveMode(alias, fileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		// Apply prompt-level loop overrides if running from alias
		if mode == prompt.ModeAlias {
			cfg, err = config.ResolveEffectiveConfigForPrompt(cfg, alias)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		}

		// Overlay CLI flags (after prompt overrides, so CLI takes precedence)
		config.OverlayCLIFlags(&cfg, cliFlags)

		// Initialize logger with effective log level (O4/R5)
		// Precedence: --log-level > --quiet > --verbose > config > default (info)
		logLevel := cfg.Loop.LogLevel.Value
		if cmd.Flags().Changed("verbose") && !cmd.Flags().Changed("log-level") && !cmd.Flags().Changed("quiet") {
			logLevel = "debug"
		}
		if err := logger.SetLevel(logLevel); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		// Log config provenance at debug level (O2/R1)
		logger.Debug("config: loop.iteration_mode = %s (source: %s)", cfg.Loop.IterationMode.Value, cfg.Loop.IterationMode.Provenance)
		logger.Debug("config: loop.default_max_iterations = %d (source: %s)", cfg.Loop.DefaultMaxIterations.Value, cfg.Loop.DefaultMaxIterations.Provenance)
		logger.Debug("config: loop.failure_threshold = %d (source: %s)", cfg.Loop.FailureThreshold.Value, cfg.Loop.FailureThreshold.Provenance)
		logger.Debug("config: loop.iteration_timeout = %d (source: %s)", cfg.Loop.IterationTimeout.Value, cfg.Loop.IterationTimeout.Provenance)
		logger.Debug("config: loop.max_output_buffer = %d (source: %s)", cfg.Loop.MaxOutputBuffer.Value, cfg.Loop.MaxOutputBuffer.Provenance)
		logger.Debug("config: loop.log_level = %s (source: %s)", cfg.Loop.LogLevel.Value, cfg.Loop.LogLevel.Provenance)
		logger.Debug("config: loop.show_ai_output = %v (source: %s)", cfg.Loop.ShowAIOutput.Value, cfg.Loop.ShowAIOutput.Provenance)
		logger.Debug("config: loop.preamble = %v (source: %s)", cfg.Loop.Preamble.Value, cfg.Loop.Preamble.Provenance)
		logger.Debug("config: loop.ai_cmd = %s (source: %s)", cfg.Loop.AICmd.Value, cfg.Loop.AICmd.Provenance)
		logger.Debug("config: loop.ai_cmd_alias = %s (source: %s)", cfg.Loop.AICmdAlias.Value, cfg.Loop.AICmdAlias.Provenance)
		logger.Debug("config: loop.signals.success = %s (source: %s)", cfg.Loop.SignalSuccess.Value, cfg.Loop.SignalSuccess.Provenance)
		logger.Debug("config: loop.signals.failure = %s (source: %s)", cfg.Loop.SignalFailure.Value, cfg.Loop.SignalFailure.Provenance)

		// Load prompt content
		src, err := prompt.LoadPrompt(mode, alias, fileFlag, &cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		// Dry-run mode: print assembled prompt and exit (O4/R4)
		if dryRunFlag {
			// Print config with provenance (O2/R1)
			fmt.Fprintf(os.Stderr, "=== Configuration ===\n")
			fmt.Fprintf(os.Stderr, "loop.iteration_mode = %s (source: %s)\n", cfg.Loop.IterationMode.Value, cfg.Loop.IterationMode.Provenance)
			fmt.Fprintf(os.Stderr, "loop.default_max_iterations = %d (source: %s)\n", cfg.Loop.DefaultMaxIterations.Value, cfg.Loop.DefaultMaxIterations.Provenance)
			fmt.Fprintf(os.Stderr, "loop.failure_threshold = %d (source: %s)\n", cfg.Loop.FailureThreshold.Value, cfg.Loop.FailureThreshold.Provenance)
			fmt.Fprintf(os.Stderr, "loop.iteration_timeout = %d (source: %s)\n", cfg.Loop.IterationTimeout.Value, cfg.Loop.IterationTimeout.Provenance)
			fmt.Fprintf(os.Stderr, "loop.max_output_buffer = %d (source: %s)\n", cfg.Loop.MaxOutputBuffer.Value, cfg.Loop.MaxOutputBuffer.Provenance)
			fmt.Fprintf(os.Stderr, "loop.log_level = %s (source: %s)\n", cfg.Loop.LogLevel.Value, cfg.Loop.LogLevel.Provenance)
			fmt.Fprintf(os.Stderr, "loop.show_ai_output = %v (source: %s)\n", cfg.Loop.ShowAIOutput.Value, cfg.Loop.ShowAIOutput.Provenance)
			fmt.Fprintf(os.Stderr, "loop.preamble = %v (source: %s)\n", cfg.Loop.Preamble.Value, cfg.Loop.Preamble.Provenance)
			fmt.Fprintf(os.Stderr, "loop.ai_cmd = %s (source: %s)\n", cfg.Loop.AICmd.Value, cfg.Loop.AICmd.Provenance)
			fmt.Fprintf(os.Stderr, "loop.ai_cmd_alias = %s (source: %s)\n", cfg.Loop.AICmdAlias.Value, cfg.Loop.AICmdAlias.Provenance)
			fmt.Fprintf(os.Stderr, "loop.signals.success = %s (source: %s)\n", cfg.Loop.SignalSuccess.Value, cfg.Loop.SignalSuccess.Provenance)
			fmt.Fprintf(os.Stderr, "loop.signals.failure = %s (source: %s)\n", cfg.Loop.SignalFailure.Value, cfg.Loop.SignalFailure.Provenance)
			fmt.Fprintf(os.Stderr, "\n=== Assembled Prompt ===\n")

			maxIterations := cfg.Loop.DefaultMaxIterations.Value
			iterationMode := cfg.Loop.IterationMode.Value
			preambleEnabled := cfg.Loop.Preamble.Value

			preambleCfg := runner.PreambleConfig{
				Enabled:        preambleEnabled,
				Iteration:      1,
				MaxIterations:  maxIterations,
				Unlimited:      iterationMode == "unlimited",
				ContextStrings: contextFlags,
			}

			preamble := runner.GeneratePreamble(preambleCfg)
			assembled := runner.AssemblePrompt(preamble, src.Content)
			fmt.Print(string(assembled))
			os.Exit(0)
		}

		// Resolve AI command (alias or direct command)
		resolution, err := config.ResolveAICommand(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		// Parse AI command string into argv
		aiCmd, err := cmdparse.Parse(resolution.Command)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to parse AI command: %v\n", err)
			os.Exit(1)
		}

		// Execute loop
		loopErr := runner.Loop(aiCmd, src.Content, &cfg, contextFlags)

		// Map loop error to exit code per O4/R1
		if loopErr == nil {
			// Success signal received
			os.Exit(0)
		} else if errors.Is(loopErr, runner.ExitCodeFailureThreshold) {
			// Failure threshold reached
			os.Exit(1)
		} else if errors.Is(loopErr, runner.ExitCodeExhausted) {
			// Max iterations exhausted
			os.Exit(2)
		} else if errors.Is(loopErr, runner.ExitCodeInterrupted) {
			// Interrupted by SIGINT/SIGTERM
			os.Exit(130)
		} else {
			// Unexpected error
			fmt.Fprintf(os.Stderr, "error: %v\n", loopErr)
			os.Exit(1)
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured resources",
}

var listPromptsCmd = &cobra.Command{
	Use:   "prompts",
	Short: "List prompt aliases",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfigWithProvenanceAndExplicit(configFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		if err := config.Validate(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		if len(cfg.Prompts) == 0 {
			fmt.Println("No prompts configured.")
			return
		}

		var keys []string
		for k := range cfg.Prompts {
			keys = append(keys, k)
		}

		// Sort alphabetically
		for i := 0; i < len(keys); i++ {
			for j := i + 1; j < len(keys); j++ {
				if keys[i] > keys[j] {
					keys[i], keys[j] = keys[j], keys[i]
				}
			}
		}

		for _, k := range keys {
			p := cfg.Prompts[k]
			fmt.Printf("%s:\n", k)
			if p.Name.Value != "" {
				fmt.Printf("  name: %s\n", p.Name.Value)
			}
			if p.Description.Value != "" {
				fmt.Printf("  description: %s\n", p.Description.Value)
			}
			fmt.Printf("  path: %s\n", p.Path.Value)
		}
	},
}

var listAliasesCmd = &cobra.Command{
	Use:   "aliases",
	Short: "List AI command aliases",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfigWithProvenanceAndExplicit(configFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		plainCfg := config.Config{AICmdAliases: make(map[string]string)}
		for k, v := range cfg.AICmdAliases {
			plainCfg.AICmdAliases[k] = v.Value
		}

		merged := config.MergedAliases(plainCfg)
		
		var keys []string
		for k := range merged {
			keys = append(keys, k)
		}
		
		// Sort alphabetically
		for i := 0; i < len(keys); i++ {
			for j := i + 1; j < len(keys); j++ {
				if keys[i] > keys[j] {
					keys[i], keys[j] = keys[j], keys[i]
				}
			}
		}

		for _, k := range keys {
			fmt.Printf("%s:\n  command: %s\n", k, merged[k])
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ralph 0.1.0")
	},
}

// reviewCmd implements the ralph review subcommand (O5). T1: input modes and validation; rest stubbed.
var reviewCmd = &cobra.Command{
	Use:   "review [alias]",
	Short: "Review a prompt for quality and structure (O5)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfigWithProvenanceAndExplicit(configFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(2)
		}
		if err := config.Validate(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(2)
		}

		var alias string
		if len(args) > 0 {
			alias = args[0]
		}
		// R1: resolution order file > alias > stdin (same ResolveMode as run)
		mode, err := prompt.ResolveMode(alias, reviewFileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(2)
		}

		if mode == prompt.ModeAlias {
			cfg, err = config.ResolveEffectiveConfigForPrompt(cfg, alias)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(2)
			}
		}

		src, err := prompt.LoadPrompt(mode, alias, reviewFileFlag, &cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(2)
		}

		// R4: Resolve prompt output path; validate stdin+apply requires --prompt-output (T7)
		reviewInputMode := promptModeToReviewInputMode(mode)
		promptOutputResult, err := review.ResolvePromptOutputPath(reviewInputMode, reviewApplyFlag, reviewPromptOutputFlag, src.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(2)
		}
		_ = promptOutputResult // T8 will use .Path for revision-phase interpolation

		// R3: Resolve report path (--review-output or temp); validate before spawning AI (T2)
		reportPath, isTemp, err := review.ResolveReportPath(reviewOutputFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(2)
		}
		if isTemp {
			fmt.Fprintf(os.Stderr, "Report will be written to: %s\n", reportPath)
		}

		// R2: Compose review prompt (embedded instructions + path directive + user prompt)
		composed, err := review.ComposeReviewPrompt(reportPath, src.Content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(2)
		}

		// Resolve and parse AI command (same mechanism as runner)
		resolution, err := config.ResolveAICommand(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(2)
		}
		aiArgv, err := cmdparse.Parse(resolution.Command)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to parse AI command: %v\n", err)
			os.Exit(2)
		}

		// Single AI invocation; report is file at reportPath (R2), not parsed from stdout
		_, err = runner.SpawnAI(aiArgv, bytes.NewReader(composed), os.Stdout, os.Stderr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: AI spawn failed: %v\n", err)
			os.Exit(2)
		}

		// R9: Verify report file exists at reportPath after review-phase AI exits (R8 condition 7)
		if err := review.VerifyReportExists(reportPath); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(2)
		}

		// R6: Parse report for machine-parseable summary; derive exit 0 (no errors) or 1 (errors in prompt).
		code := review.ParseReportSummary(reportPath)
		os.Exit(code)
	},
}

func init() {
	// Configuration
	runCmd.Flags().StringVar(&configFlag, "config", "", "Explicit config file path")
	listPromptsCmd.Flags().StringVar(&configFlag, "config", "", "Explicit config file path")
	listAliasesCmd.Flags().StringVar(&configFlag, "config", "", "Explicit config file path")
	reviewCmd.Flags().StringVar(&configFlag, "config", "", "Explicit config file path")

	// Prompt input
	runCmd.Flags().StringVarP(&fileFlag, "file", "f", "", "Read prompt from file")
	reviewCmd.Flags().StringVarP(&reviewFileFlag, "file", "f", "", "Read prompt from file (R1: wins over alias)")
	reviewCmd.Flags().StringVar(&reviewOutputFlag, "review-output", "", "Write review report to this path (default: temp file; path communicated to user)")
	reviewCmd.Flags().StringVar(&reviewPromptOutputFlag, "prompt-output", "", "Write suggested revised prompt to this path; required with --apply when input is stdin (R4)")
	reviewCmd.Flags().BoolVar(&reviewApplyFlag, "apply", false, "Apply suggested revision to prompt file (or --prompt-output when set); with stdin, --prompt-output required (R5)")

	// Loop control
	runCmd.Flags().IntVarP(&maxIterationsFlag, "max-iterations", "n", 0, "Override max iterations")
	runCmd.Flags().BoolVarP(&unlimitedFlag, "unlimited", "u", false, "Run until signal or failure threshold")
	runCmd.Flags().IntVar(&failureThresholdFlag, "failure-threshold", 0, "Consecutive failures before abort")
	runCmd.Flags().IntVar(&iterationTimeoutFlag, "iteration-timeout", 0, "Per-iteration timeout in seconds")
	runCmd.Flags().IntVar(&maxOutputBufferFlag, "max-output-buffer", 0, "Max output buffer in bytes")
	runCmd.Flags().BoolVar(&preambleFlag, "preamble", false, "Enable preamble injection")
	runCmd.Flags().BoolVar(&noPreambleFlag, "no-preamble", false, "Disable preamble injection")
	runCmd.Flags().BoolVarP(&dryRunFlag, "dry-run", "d", false, "Validate and show assembled prompt")

	// AI command
	runCmd.Flags().StringVar(&aiCmdFlag, "ai-cmd", "", "Direct AI command string")
	runCmd.Flags().StringVar(&aiCmdAliasFlag, "ai-cmd-alias", "", "AI command alias")

	// Signals
	runCmd.Flags().StringVar(&signalSuccessFlag, "signal-success", "", "Success signal string")
	runCmd.Flags().StringVar(&signalFailureFlag, "signal-failure", "", "Failure signal string")

	// Context
	runCmd.Flags().StringArrayVarP(&contextFlags, "context", "c", nil, "Inject context into preamble (repeatable)")

	// Output control
	runCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Stream AI output to terminal")
	runCmd.Flags().BoolVarP(&quietFlag, "quiet", "q", false, "Suppress non-error output")
	runCmd.Flags().StringVar(&logLevelFlag, "log-level", "", "Set log level (debug, info, warn, error)")

	listCmd.AddCommand(listPromptsCmd)
	listCmd.AddCommand(listAliasesCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(reviewCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(versionCmd)
}

// promptModeToReviewInputMode maps prompt.Mode to review.InputMode for R4 path resolution.
func promptModeToReviewInputMode(mode prompt.Mode) review.InputMode {
	switch mode {
	case prompt.ModeAlias:
		return review.InputModeAlias
	case prompt.ModeFile:
		return review.InputModeFile
	case prompt.ModeStdin:
		return review.InputModeStdin
	default:
		return review.InputModeStdin
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
