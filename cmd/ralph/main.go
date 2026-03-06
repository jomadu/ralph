package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/maxdunn/ralph/internal/config"
	"github.com/maxdunn/ralph/internal/prompt"
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
		if cmd.Flags().Changed("failure-threshold") {
			cliFlags.FailureThreshold = &failureThresholdFlag
		}
		if cmd.Flags().Changed("iteration-timeout") {
			cliFlags.IterationTimeout = &iterationTimeoutFlag
		}
		if cmd.Flags().Changed("max-output-buffer") {
			cliFlags.MaxOutputBuffer = &maxOutputBufferFlag
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

		// Handle preamble/no-preamble flags
		if cmd.Flags().Changed("preamble") {
			cliFlags.Preamble = &preambleFlag
		} else if cmd.Flags().Changed("no-preamble") {
			noPreamble := !noPreambleFlag
			cliFlags.Preamble = &noPreamble
		}

		// Overlay CLI flags
		config.OverlayCLIFlags(&cfg, cliFlags)

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

		// Load prompt content
		src, err := prompt.LoadPrompt(mode, alias, fileFlag, &cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Prompt loaded: mode=%v, size=%d bytes\n", src.Mode, len(src.Content))
		fmt.Println("run: loop execution not yet implemented")
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
		fmt.Println("list prompts: not implemented")
	},
}

var listAliasesCmd = &cobra.Command{
	Use:   "aliases",
	Short: "List AI command aliases",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list aliases: not implemented")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ralph 0.1.0")
	},
}

func init() {
	// Configuration
	runCmd.Flags().StringVar(&configFlag, "config", "", "Explicit config file path")

	// Prompt input
	runCmd.Flags().StringVarP(&fileFlag, "file", "f", "", "Read prompt from file")

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
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
