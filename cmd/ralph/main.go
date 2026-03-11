package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/maxdunn/ralph/internal/config"
	"github.com/maxdunn/ralph/internal/review"
	"github.com/maxdunn/ralph/internal/runloop"
	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags (e.g. make build VERSION=1.0.0).
var Version = "dev"

func main() {
	root := &cobra.Command{
		Use:   "ralph",
		Short: "Ralph — loop runner for AI-driven tasks",
	}
	root.SetVersionTemplate("{{.Version}}\n")
	root.Version = Version
	root.PersistentFlags().String("config", "", "Use this file as the sole file-based config (must exist)")
	root.PersistentFlags().Bool("version", false, "Print version and exit 0")
	// Unknown top-level command → error to stderr, non-zero exit, suggest --help (cli.md).
	root.SilenceUsage = true
	root.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintf(os.Stderr, "Run '%s --help' for usage.\n", cmd.CommandPath())
		return err
	})
	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Global --version: print and exit 0 (cli.md Global options).
		if v, _ := cmd.Root().PersistentFlags().GetBool("version"); v {
			fmt.Println(Version)
			os.Exit(0)
		}
		// When --config is set, file must exist (R005); fail fast before subcommand.
		configPath, _ := cmd.Root().PersistentFlags().GetString("config")
		if configPath != "" {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			path := configPath
			if !filepath.IsAbs(path) {
				path = filepath.Join(cwd, path)
			}
			if _, err := os.Stat(path); err != nil {
				if os.IsNotExist(err) {
					fmt.Fprintf(os.Stderr, "ralph: config file not found: %s\n", configPath)
				} else {
					fmt.Fprintf(os.Stderr, "ralph: config file: %s\n", err)
				}
				os.Exit(1)
			}
		}
		return nil
	}
	root.AddCommand(runCmd())
	root.AddCommand(reviewCmd())
	root.AddCommand(listCmd())
	root.AddCommand(showCmd())
	root.AddCommand(versionCmd())
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintf(os.Stderr, "Run '%s --help' for usage.\n", root.Name())
		os.Exit(1)
	}
}

func runCmd() *cobra.Command {
	var (
		filePath         string
		aiCmd            string
		aiCmdAlias       string
		dryRun           bool
		maxIterations    int
		unlimited        bool
		failureThreshold int
		iterationTimeout int
		maxOutputBuffer  int
		noPreamble       bool
		signalSuccess    string
		signalFailure    string
		signalPrecedence string
		contextStrings   []string
		verbose          bool
		quiet            bool
		logLevel         string
		stream           bool
		noStream         bool
	)
	cmd := &cobra.Command{
		Use:   "run [alias]",
		Short: "Run the iteration loop",
		Long:  "Prompt via alias, --file/-f <path>, or stdin. Exactly one source. Runs until success signal or failure threshold / max iterations.",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, _ := cmd.Root().PersistentFlags().GetString("config")
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			// Validate loop flag values (cli.md: invalid or out-of-range → error and exit non-zero).
			// Defaults -1 mean "not set"; we only validate when user passed a value (including 0 for timeouts).
			if maxIterations < 0 {
				fmt.Fprintln(os.Stderr, "ralph run: --max-iterations must be >= 0")
				os.Exit(1)
			}
			if failureThreshold != -1 && failureThreshold < 0 {
				fmt.Fprintln(os.Stderr, "ralph run: --failure-threshold must be >= 0")
				os.Exit(1)
			}
			if iterationTimeout != -1 && iterationTimeout < 0 {
				fmt.Fprintln(os.Stderr, "ralph run: --iteration-timeout must be >= 0")
				os.Exit(1)
			}
			if maxOutputBuffer != -1 && maxOutputBuffer < 0 {
				fmt.Fprintln(os.Stderr, "ralph run: --max-output-buffer must be >= 0")
				os.Exit(1)
			}
			if logLevel != "" && !config.ValidLogLevel(logLevel) {
				fmt.Fprintf(os.Stderr, "ralph run: --log-level must be debug, info, warn, or error (got %q)\n", logLevel)
				os.Exit(1)
			}
			opts := review.ResolveOptions{}
			if len(args) > 0 {
				opts.Alias = args[0]
			}
			if filePath != "" {
				opts.FilePath = filePath
			}
			if opts.Alias != "" && opts.FilePath != "" {
				fmt.Fprintln(os.Stderr, "ralph run: exactly one of alias, --file/-f, or stdin required (cannot combine alias with --file)")
				os.Exit(1)
			}
			if opts.Alias == "" && opts.FilePath == "" {
				if stat, err := os.Stdin.Stat(); err == nil && (stat.Mode()&os.ModeCharDevice) != 0 {
					fmt.Fprintln(os.Stderr, "ralph run: no prompt source (stdin is a terminal); use an alias, -f <path>, or pipe input")
					os.Exit(1)
				}
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					opts.Stdin = append(opts.Stdin, scanner.Bytes()...)
					opts.Stdin = append(opts.Stdin, '\n')
				}
				if err := scanner.Err(); err != nil {
					return err
				}
			}
			provider, err := promptProviderForCmd(cmd, configPath, cwd)
			if err != nil {
				return fmt.Errorf("config: %w", err)
			}
			promptBytes, err := runloop.LoadPromptOnce(provider, cwd, opts)
			if err != nil {
				if review.IsExit2(err) {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(2)
				}
				return err
			}
			promptName := opts.Alias
			eff, ok, err := config.Resolve(os.Getenv, cwd, configPath, promptName)
			if err != nil {
				return fmt.Errorf("config: %w", err)
			}
			if !ok && promptName != "" {
				return fmt.Errorf("config: prompt %q not found", promptName)
			}
			overlay, _ := config.ParseEnvOverlay(os.Getenv)
			directCmd := aiCmd
			if directCmd == "" && overlay != nil && overlay.AICmd != nil {
				directCmd = *overlay.AICmd
			}
			aliasName := aiCmdAlias
			if aliasName == "" && overlay != nil && overlay.AICmdAlias != nil {
				aliasName = *overlay.AICmdAlias
			}
			if directCmd == "" && aliasName == "" {
				aliasName = "cursor-agent"
			}
			command, ok := config.ResolveAICommand(eff, directCmd, aliasName)
			if !ok {
				fmt.Fprintln(os.Stderr, "ralph run: AI command not resolved (missing or invalid --ai-cmd / --ai-cmd-alias or config)")
				os.Exit(runloop.ExitErrorPreLoop)
			}
			// Apply CLI overrides to effective loop (cli.md: flags override config for this run).
			loop := applyRunLoopOverrides(eff.Loop, runLoopOverrides{
				maxIterations:    maxIterations,
				unlimited:        unlimited,
				failureThreshold: failureThreshold,
				iterationTimeout: iterationTimeout,
				maxOutputBuffer:  maxOutputBuffer,
				noPreamble:       noPreamble,
				signalSuccess:    signalSuccess,
				signalFailure:    signalFailure,
				signalPrecedence: signalPrecedence,
				context:          contextStrings,
				verbose:          verbose,
				quiet:            quiet,
				logLevel:         logLevel,
				stream:           stream,
				noStream:         noStream,
			})
			streamWriter := io.Writer(io.Discard)
			if loop.Streaming {
				streamWriter = os.Stdout
			}
			runOpts := runloop.RunOptions{
				Command:      command,
				PromptBytes:  promptBytes,
				Loop:         loop,
				Cwd:          cwd,
				Env:          os.Environ(),
				DryRun:       dryRun,
				StreamWriter: streamWriter,
			}
			code, err := runloop.Run(runOpts)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(runloop.ExitErrorPreLoop)
			}
			os.Exit(code)
			return nil
		},
	}
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Read prompt from this file (mutually exclusive with alias and stdin)")
	// Loop control
	cmd.Flags().IntVarP(&maxIterations, "max-iterations", "n", 0, "Override max iterations for this run (0 = use config default)")
	cmd.Flags().BoolVarP(&unlimited, "unlimited", "u", false, "Run until success signal or failure threshold; no iteration cap")
	cmd.Flags().IntVar(&failureThreshold, "failure-threshold", -1, "Consecutive failures before exit; override for this run")
	cmd.Flags().IntVar(&iterationTimeout, "iteration-timeout", -1, "Per-iteration timeout in seconds (0 = no timeout)")
	cmd.Flags().IntVar(&maxOutputBuffer, "max-output-buffer", -1, "Max output buffer in bytes for capturing AI stdout")
	cmd.Flags().BoolVar(&noPreamble, "no-preamble", false, "Disable preamble injection for this run")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Print assembled prompt and exit 0; do not invoke AI")
	// AI command
	cmd.Flags().StringVar(&aiCmd, "ai-cmd", "", "Direct AI command string for this run")
	cmd.Flags().StringVar(&aiCmdAlias, "ai-cmd-alias", "", "AI command alias name from config for this run")
	// Signals
	cmd.Flags().StringVar(&signalSuccess, "signal-success", "", "Success signal string for this run")
	cmd.Flags().StringVar(&signalFailure, "signal-failure", "", "Failure signal string for this run")
	cmd.Flags().StringVar(&signalPrecedence, "signal-precedence", "", "When both signals appear: static or ai_interpreted")
	// Context / preamble
	cmd.Flags().StringArrayVarP(&contextStrings, "context", "c", nil, "Inline context injected into preamble (repeatable)")
	// Output and observability (cli.md ralph run)
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Increase verbosity (e.g. log level debug and enable streaming)")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Minimal output: log level error-only; streaming disabled")
	cmd.Flags().StringVar(&logLevel, "log-level", "", "Log level: debug, info, warn, error (overrides config and shortcuts)")
	cmd.Flags().BoolVar(&stream, "stream", false, "Enable streaming of AI output to terminal (still captured for signal scan)")
	cmd.Flags().BoolVar(&noStream, "no-stream", false, "Disable streaming of AI output to terminal")
	return cmd
}

// runLoopOverrides holds CLI flag values that override effective config for this run (cli.md ralph run Flags).
type runLoopOverrides struct {
	maxIterations    int
	unlimited        bool
	failureThreshold int
	iterationTimeout int
	maxOutputBuffer  int
	noPreamble       bool
	signalSuccess    string
	signalFailure    string
	signalPrecedence string
	context          []string
	verbose          bool
	quiet            bool
	logLevel         string
	stream           bool
	noStream         bool
}

// applyRunLoopOverrides applies CLI overrides to the effective loop. Only non-zero or set overrides apply.
// maxIterations 0 = use config; -1 sentinel for "not set" is not used (we use 0 for max-iterations default).
// For failure-threshold, iteration-timeout, max-output-buffer we use -1 as "not set".
func applyRunLoopOverrides(base config.LoopSettings, o runLoopOverrides) config.LoopSettings {
	out := base
	if o.unlimited {
		// Run-loop has a bounded for-loop; use a large cap so it effectively runs until signal or threshold.
		const unlimitedCap = 1<<31 - 1
		out.MaxIterations = unlimitedCap
	} else if o.maxIterations > 0 {
		out.MaxIterations = o.maxIterations
	}
	if o.failureThreshold >= 0 {
		out.FailureThreshold = o.failureThreshold
	}
	if o.iterationTimeout >= 0 {
		out.TimeoutSeconds = o.iterationTimeout
	}
	if o.signalSuccess != "" {
		out.SuccessSignal = o.signalSuccess
	}
	if o.signalFailure != "" {
		out.FailureSignal = o.signalFailure
	}
	if o.signalPrecedence != "" {
		out.SignalPrecedence = o.signalPrecedence
	}
	if o.noPreamble {
		out.Preamble = ""
	} else if len(o.context) > 0 {
		// -c/--context: inject as CONTEXT section into preamble (repeatable; join with newlines).
		contextBlock := "CONTEXT\n" + strings.Join(o.context, "\n")
		if out.Preamble != "" {
			out.Preamble = out.Preamble + "\n" + contextBlock
		} else {
			out.Preamble = contextBlock
		}
	}
	// Output and observability (cli.md: --verbose/-v, --quiet/-q, --log-level, --stream, --no-stream).
	// Quiet wins for minimal output unless log-level or stream is explicitly set (O004/R006).
	if o.quiet {
		out.LogLevel = "error"
		out.Streaming = false
	}
	if o.verbose {
		out.LogLevel = "debug"
		out.Streaming = true
	}
	if o.logLevel != "" {
		out.LogLevel = o.logLevel
	}
	if o.stream {
		out.Streaming = true
	}
	if o.noStream {
		out.Streaming = false
	}
	// maxOutputBuffer: parsed and validated but run-loop/backend do not yet use it; no field on LoopSettings.
	_ = o.maxOutputBuffer
	return out
}

// promptProviderForCmd builds a review.PromptProvider from root --config and cwd (explicit file or global+workspace).
func promptProviderForCmd(cmd *cobra.Command, configPath, cwd string) (review.PromptProvider, error) {
	if configPath != "" {
		path := configPath
		if !filepath.IsAbs(path) {
			path = filepath.Join(cwd, path)
		}
		layer, err := config.LoadExplicit(path)
		if err != nil {
			return nil, err
		}
		return &review.FileLayerProvider{Layer: layer}, nil
	}
	global, workspace, err := config.LoadGlobalAndWorkspace(os.Getenv, cwd)
	if err != nil {
		return nil, err
	}
	return &mergedPromptProvider{global: global, workspace: workspace}, nil
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version string and exit 0",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(Version)
			return nil
		},
	}
}

func showCmd() *cobra.Command {
	showRoot := &cobra.Command{
		Use:   "show",
		Short: "Show effective config or detail for a prompt/alias",
		Long:  "Use 'show config', 'show prompt [name]', or 'show alias [name]'. Same config resolution as run.",
	}
	showRoot.AddCommand(showConfigCmd())
	showRoot.AddCommand(showPromptCmd())
	showRoot.AddCommand(showAliasCmd())
	return showRoot
}

func showConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Output the effective config for the current context",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("ralph show config: unexpected argument %q", args[0])
			}
			configPath, _ := cmd.Root().PersistentFlags().GetString("config")
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			eff, _, err := config.Resolve(os.Getenv, cwd, configPath, "")
			if err != nil {
				return fmt.Errorf("config: %w", err)
			}
			// Simple YAML-like output for effective config (loop + prompts + aliases).
			fmt.Printf("loop:\n  max_iterations: %d\n  failure_threshold: %d\n  timeout_seconds: %d\n  success_signal: %q\n  failure_signal: %q\n  log_level: %q\n",
				eff.Loop.MaxIterations, eff.Loop.FailureThreshold, eff.Loop.TimeoutSeconds,
				eff.Loop.SuccessSignal, eff.Loop.FailureSignal, eff.Loop.LogLevel)
			if len(eff.Prompts) > 0 {
				fmt.Println("prompts:")
				for name, p := range eff.Prompts {
					fmt.Printf("  %s:\n", name)
					if p.Path != "" {
						fmt.Printf("    path: %q\n", p.Path)
					}
					if p.Content != "" {
						fmt.Printf("    content: \"\"\"\n%s\"\"\"\n", p.Content)
					}
				}
			}
			if len(eff.Aliases) > 0 {
				fmt.Println("aliases:")
				for name, a := range eff.Aliases {
					fmt.Printf("  %s: %q\n", name, a.Command)
				}
			}
			return nil
		},
	}
}

func showPromptCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prompt [name]",
		Short: "Show detailed information for the prompt named [name]",
		Long:  "Name is required. Errors if the prompt is not defined in resolved config.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Fprintln(os.Stderr, "ralph show prompt: name required (use 'ralph show prompt <name>' or 'ralph list prompts')")
				os.Exit(1)
			}
			name := args[0]
			configPath, _ := cmd.Root().PersistentFlags().GetString("config")
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			eff, _, err := config.Resolve(os.Getenv, cwd, configPath, "")
			if err != nil {
				return fmt.Errorf("config: %w", err)
			}
			p, ok := eff.Prompts[name]
			if !ok {
				fmt.Fprintf(os.Stderr, "ralph show prompt: unknown prompt %q\n", name)
				os.Exit(1)
			}
			fmt.Printf("name: %s\n", name)
			if p.DisplayName != "" {
				fmt.Printf("display_name: %q\n", p.DisplayName)
			}
			if p.Description != "" {
				fmt.Printf("description: %q\n", p.Description)
			}
			if p.Path != "" {
				fmt.Printf("path: %q\n", p.Path)
			}
			if p.Content != "" {
				fmt.Printf("content: \"\"\"\n%s\"\"\"\n", p.Content)
			}
			return nil
		},
	}
}

func showAliasCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "alias [name]",
		Short: "Show detailed information for the alias named [name]",
		Long:  "Name is required. Errors if the alias is not defined in resolved config.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Fprintln(os.Stderr, "ralph show alias: name required (use 'ralph show alias <name>' or 'ralph list aliases')")
				os.Exit(1)
			}
			name := args[0]
			configPath, _ := cmd.Root().PersistentFlags().GetString("config")
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			eff, _, err := config.Resolve(os.Getenv, cwd, configPath, "")
			if err != nil {
				return fmt.Errorf("config: %w", err)
			}
			a, ok := eff.Aliases[name]
			if !ok {
				fmt.Fprintf(os.Stderr, "ralph show alias: unknown alias %q\n", name)
				os.Exit(1)
			}
			fmt.Printf("name: %s\ncommand: %q\n", name, a.Command)
			return nil
		},
	}
}

// mergedPromptProvider implements review.PromptProvider by checking workspace then global.
type mergedPromptProvider struct {
	global, workspace *config.FileLayer
}

func (m *mergedPromptProvider) PromptByName(name string) (path, content string, ok bool) {
	if m.workspace != nil && m.workspace.Prompts != nil {
		if p, ok := m.workspace.Prompts[name]; ok {
			return p.Path, p.Content, true
		}
	}
	if m.global != nil && m.global.Prompts != nil {
		if p, ok := m.global.Prompts[name]; ok {
			return p.Path, p.Content, true
		}
	}
	return "", "", false
}

func reviewCmd() *cobra.Command {
	var (
		filePath     string
		reportPath   string
		promptOutput string
		apply        bool
		yes          bool
		quiet        bool
		logLevel     string
	)
	cmd := &cobra.Command{
		Use:   "review [alias]",
		Short: "Review a prompt (alias, file, or stdin); report and optional apply",
		Long:  "Exactly one of: positional alias, --file/-f <path>, or stdin. Produces a report and optional suggested revision.",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, _ := cmd.Root().PersistentFlags().GetString("config")
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			provider, err := promptProviderForCmd(cmd, configPath, cwd)
			if err != nil {
				return fmt.Errorf("config: %w", err)
			}

			opts := review.ResolveOptions{}
			if len(args) > 0 {
				opts.Alias = args[0]
			}
			if filePath != "" {
				opts.FilePath = filePath
			}
			if opts.Alias != "" && opts.FilePath != "" {
				fmt.Fprintln(os.Stderr, "ralph review: exactly one of alias, --file/-f, or stdin required (cannot combine alias with --file)")
				os.Exit(1)
			}
			if opts.Alias == "" && opts.FilePath == "" {
				// Stdin: error if TTY (no prompt source)
				if stat, err := os.Stdin.Stat(); err == nil && (stat.Mode()&os.ModeCharDevice) != 0 {
					fmt.Fprintln(os.Stderr, "ralph review: no prompt source (stdin is a terminal); use an alias, -f <path>, or pipe input")
					os.Exit(2)
				}
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					opts.Stdin = append(opts.Stdin, scanner.Bytes()...)
					opts.Stdin = append(opts.Stdin, '\n')
				}
				if err := scanner.Err(); err != nil {
					return err
				}
				// Stdin + apply requires revision output path (cli.md, review error handling).
				if apply && promptOutput == "" {
					fmt.Fprintln(os.Stderr, "ralph review: --apply requires --prompt-output when prompt is from stdin")
					os.Exit(2)
				}
			}

			content, err := review.ResolvePromptSource(provider, cwd, opts)
			if err != nil {
				if review.IsExit2(err) {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(2)
				}
				return err
			}

			runOpts := review.RunOptions{
				ReportPath:       reportPath,
				PromptOutputPath: promptOutput,
				Apply:            apply,
				Yes:              yes,
				Quiet:            quiet,
				LogLevel:         logLevel,
			}
			if err := review.Run(content, runOpts); err != nil {
				if review.IsExit2(err) {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(2)
				}
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Read prompt from this file (mutually exclusive with alias and stdin)")
	cmd.Flags().StringVar(&reportPath, "report", "", "Write report file to this path")
	cmd.Flags().StringVar(&promptOutput, "prompt-output", "", "When using --apply, write revision to this path (required when prompt is from stdin)")
	cmd.Flags().BoolVar(&apply, "apply", false, "Write suggested revision to a file")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Non-interactive apply: do not prompt for confirmation")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Minimize output")
	cmd.Flags().StringVar(&logLevel, "log-level", "", "Log level: debug, info, warn, error")
	return cmd
}

// listCmd returns the ralph list command and its subcommands (list, list prompts, list aliases).
// Uses the same config resolution as run/review: explicit file only, or global + workspace.
func listCmd() *cobra.Command {
	listRoot := &cobra.Command{
		Use:   "list [prompts|aliases]",
		Short: "List prompts and/or AI command aliases from resolved config",
		Long:  "With no subcommand: list all. Use 'list prompts' or 'list aliases' to list only one. Same config resolution as run (global, workspace, or --config). Read-only.",
		RunE:  runListAll,
	}
	listRoot.AddCommand(listPromptsCmd())
	listRoot.AddCommand(listAliasesCmd())
	return listRoot
}

func runListAll(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("ralph list: invalid subcommand %q (use 'ralph list prompts', 'ralph list aliases', or 'ralph list --help')", args[0])
	}
	eff, _, err := resolveConfigForList(cmd)
	if err != nil {
		return err
	}
	printPrompts(eff.Prompts)
	printAliases(eff.Aliases)
	return nil
}

func listPromptsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prompts",
		Short: "List only prompts (names and optional display name, description, path)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("ralph list prompts: no positional arguments allowed")
			}
			eff, _, err := resolveConfigForList(cmd)
			if err != nil {
				return err
			}
			printPrompts(eff.Prompts)
			return nil
		},
	}
}

func listAliasesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "aliases",
		Short: "List only AI command aliases (names and expansion)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("ralph list aliases: no positional arguments allowed")
			}
			eff, _, err := resolveConfigForList(cmd)
			if err != nil {
				return err
			}
			printAliases(eff.Aliases)
			return nil
		},
	}
}

// resolveConfigForList resolves effective config via the single Resolve entrypoint
// (cwd, --config, env; no prompt). Used by list and list prompts/aliases.
func resolveConfigForList(cmd *cobra.Command) (*config.Effective, bool, error) {
	configPath, _ := cmd.Root().PersistentFlags().GetString("config")
	cwd, err := os.Getwd()
	if err != nil {
		return nil, false, err
	}
	eff, ok, err := config.Resolve(os.Getenv, cwd, configPath, "")
	if err != nil {
		return nil, false, fmt.Errorf("config: %w", err)
	}
	return eff, ok, nil
}

func printPrompts(prompts map[string]config.Prompt) {
	if len(prompts) == 0 {
		fmt.Println("(no prompts defined)")
		return
	}
	for name, p := range prompts {
		display := name
		if p.DisplayName != "" {
			display = p.DisplayName + " (" + name + ")"
		}
		fmt.Println(display)
		if p.Description != "" {
			fmt.Println("  description:", p.Description)
		}
		source := p.Path
		if source == "" && p.Content != "" {
			source = "(inline)"
		}
		if source != "" {
			fmt.Println("  path:", source)
		}
	}
}

func printAliases(aliases map[string]config.Alias) {
	if len(aliases) == 0 {
		fmt.Println("(no aliases defined)")
		return
	}
	for name, a := range aliases {
		fmt.Printf("%s\t%s\n", name, a.Command)
	}
}
