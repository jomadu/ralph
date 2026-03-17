package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jomadu/ralph/internal/config"
	"github.com/jomadu/ralph/internal/review"
	"github.com/jomadu/ralph/internal/runloop"
	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags (e.g. make build VERSION=1.0.0).
var Version = "dev"

// promptGuideContent is the full Writing Ralph prompts guide (build-time copy of docs/writing-ralph-prompts.md). Single source of truth: docs/writing-ralph-prompts.md; Makefile copies before build.
//
//go:embed embed/writing-ralph-prompts.md
var promptGuideContent []byte

// newRoot returns the root cobra command (used by main and tests).
func newRoot() *cobra.Command {
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
	return root
}

func main() {
	if err := newRoot().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintf(os.Stderr, "Run '%s --help' for usage.\n", "ralph")
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
		contextStrings   []string
		verbose          bool
		quiet            bool
		logLevel         string
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
			if directCmd == "" {
				directCmd = eff.Loop.AICmd
			}
			aliasName := aiCmdAlias
			if aliasName == "" && overlay != nil && overlay.AICmdAlias != nil {
				aliasName = *overlay.AICmdAlias
			}
			if aliasName == "" {
				aliasName = eff.Loop.AICmdAlias
			}
			command, ok := config.ResolveAICommand(eff, directCmd, aliasName)
			if !ok {
				fmt.Fprintln(os.Stderr, "ralph run: AI command not resolved (set --ai-cmd or --ai-cmd-alias, or config (loop.ai_cmd / loop.ai_cmd_alias) or env)")
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
				context:          contextStrings,
				verbose:          verbose,
				quiet:            quiet,
				logLevel:         logLevel,
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
	// Context / preamble
	cmd.Flags().StringArrayVarP(&contextStrings, "context", "c", nil, "Inline context injected into preamble (repeatable)")
	// Output and observability (cli.md ralph run)
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Increase verbosity (log level debug)")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Minimal output: log level error-only; do not show AI command output")
	cmd.Flags().StringVar(&logLevel, "log-level", "", "Log level: debug, info, warn, error (overrides config and shortcuts)")
	cmd.Flags().BoolVar(&noStream, "no-stream", false, "Do not show AI command output in the terminal")
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
	context          []string
	verbose          bool
	quiet            bool
	logLevel         string
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
	if o.noPreamble {
		out.Preamble = false
	}
	if len(o.context) > 0 {
		// -c/--context: invoker-provided context; stored as raw text, injected into CONTEXT section with a label in the run-loop.
		out.Context = strings.Join(o.context, "\n")
	}
	// Output and observability (cli.md: --verbose/-v, --quiet/-q, --log-level, --no-stream). Default is streaming on.
	// Apply shortcuts first; explicit --log-level overrides; only --no-stream turns off AI output (no enable flag).
	if o.quiet && !o.verbose {
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
	if o.noStream {
		out.Streaming = false
	}
	if o.maxOutputBuffer >= 0 {
		out.MaxOutputBuffer = o.maxOutputBuffer
	}
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
		return &review.FileLayerProvider{Layer: layer, ConfigPath: path}, nil
	}
	globalPath := config.GlobalPath(os.Getenv)
	workspacePath := config.WorkspacePath(cwd)
	global, workspace, err := config.LoadGlobalAndWorkspace(globalPath, workspacePath)
	if err != nil {
		return nil, err
	}
	return &mergedPromptProvider{global: global, workspace: workspace, globalPath: globalPath, workspacePath: workspacePath}, nil
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

// showCmd returns the ralph show command (T6.3 config, T6.4 prompt/alias, prompt-guide). cli.md: show config, show prompt [name], show alias [name], show prompt-guide.
func showCmd() *cobra.Command {
	showRoot := &cobra.Command{
		Use:   "show",
		Short: "Show effective config or detail for a prompt/alias",
		Long:  "Use 'show config', 'show prompt [name]', 'show alias [name]', or 'show prompt-guide'. Same config resolution as run (except prompt-guide).",
	}
	showRoot.AddCommand(showConfigCmd())
	showRoot.AddCommand(showPromptCmd())
	showRoot.AddCommand(showAliasCmd())
	showRoot.AddCommand(showPromptGuideCmd())
	return showRoot
}

func showConfigCmd() *cobra.Command {
	var (
		provenance bool
		promptName string
	)
	c := &cobra.Command{
		Use:   "config",
		Short: "Output the effective config for the current context",
		Long:  "Uses the same config resolution as run (default, global, workspace, explicit file, env). Use --provenance to show which layer supplied each loop value. Use --prompt to show effective config for a specific prompt (including prompt-level loop overrides).",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("ralph show config: unexpected argument %q", args[0])
			}
			configPath, _ := cmd.Root().PersistentFlags().GetString("config")
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			eff, ok, err := config.Resolve(os.Getenv, cwd, configPath, promptName)
			if err != nil {
				return fmt.Errorf("config: %w", err)
			}
			if promptName != "" && !ok {
				return fmt.Errorf("config: prompt %q not found", promptName)
			}
			loop := eff.Loop
			if provenance {
				rootInput, resolved, err := config.NewRootLoopInput(os.Getenv, cwd, configPath)
				if err != nil {
					return fmt.Errorf("config: %w", err)
				}
				loop, prov := config.LoopWithProvenance(config.LoopWithProvenanceInput{
					Root: rootInput, Resolved: resolved, PromptName: promptName, CLI: nil,
				})
				fmt.Printf("loop:\n  max_iterations: %d  # %s\n  failure_threshold: %d  # %s\n  timeout_seconds: %d  # %s\n  success_signal: %q  # %s\n  failure_signal: %q  # %s\n  preamble: %t  # %s\n  context: %q  # %s\n  streaming: %t  # %s\n  log_level: %q  # %s\n  max_output_buffer: %d  # %s\n  ai_cmd: %q  # %s\n  ai_cmd_alias: %q  # %s\n",
					loop.MaxIterations, prov.MaxIterations, loop.FailureThreshold, prov.FailureThreshold, loop.TimeoutSeconds, prov.TimeoutSeconds,
					loop.SuccessSignal, prov.SuccessSignal, loop.FailureSignal, prov.FailureSignal,
					loop.Preamble, prov.Preamble, loop.Context, prov.Context, loop.Streaming, prov.Streaming, loop.LogLevel, prov.LogLevel, loop.MaxOutputBuffer, prov.MaxOutputBuffer,
					loop.AICmd, prov.AICmd, loop.AICmdAlias, prov.AICmdAlias)
			} else {
				fmt.Printf("loop:\n  max_iterations: %d\n  failure_threshold: %d\n  timeout_seconds: %d\n  success_signal: %q\n  failure_signal: %q\n  preamble: %t\n  context: %q\n  streaming: %t\n  log_level: %q\n  max_output_buffer: %d\n  ai_cmd: %q\n  ai_cmd_alias: %q\n",
					loop.MaxIterations, loop.FailureThreshold, loop.TimeoutSeconds,
					loop.SuccessSignal, loop.FailureSignal, loop.Preamble, loop.Context, loop.Streaming, loop.LogLevel, loop.MaxOutputBuffer,
					loop.AICmd, loop.AICmdAlias)
			}
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
	c.Flags().BoolVar(&provenance, "provenance", false, "Show which layer supplied each loop value")
	c.Flags().StringVar(&promptName, "prompt", "", "Show effective config for this prompt (includes prompt-level loop overrides)")
	return c
}

func showPromptCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prompt [name]",
		Short: "Show detailed information for the prompt named [name]",
		Long:  "Name is required. Errors if the prompt is not defined in resolved config.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Fprintln(os.Stderr, "ralph show prompt: name required (use 'ralph show prompt <name>' or 'ralph list prompts')")
				return fmt.Errorf("ralph show prompt: name required")
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
				return fmt.Errorf("ralph show prompt: unknown prompt %q", name)
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
				return fmt.Errorf("ralph show alias: name required")
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
				return fmt.Errorf("ralph show alias: unknown alias %q", name)
			}
			fmt.Printf("name: %s\ncommand: %q\n", name, a.Command)
			return nil
		},
	}
}

// showPromptGuideCmd returns the ralph show prompt-guide command (PLAN T1–T3, cli.md). Outputs full guide verbatim; --markdown is optional (same output; for saving or piping to a pager).
func showPromptGuideCmd() *cobra.Command {
	var markdown bool
	c := &cobra.Command{
		Use:   "prompt-guide",
		Short: "Output the full Writing Ralph prompts guide",
		Long:  "Output the full prompt-writing guide verbatim (same content as docs/writing-ralph-prompts.md). No config required; exit 0.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("ralph show prompt-guide: unexpected argument %q", args[0])
			}
			_ = markdown // output is same with or without flag; flag exists for CLI consistency and scripts (saving/piping to pager)
			_, err := cmd.OutOrStdout().Write(promptGuideContent)
			return err
		},
	}
	c.Flags().BoolVar(&markdown, "markdown", false, "Output the full guide as markdown (e.g. for saving or piping to a pager)")
	return c
}

// mergedPromptProvider implements review.PromptProvider by checking workspace then global.
// Prompt paths are resolved relative to the config file that defined each prompt.
type mergedPromptProvider struct {
	global, workspace         *config.FileLayer
	globalPath, workspacePath string
}

func (m *mergedPromptProvider) PromptByName(name string) (path, content string, ok bool) {
	if m.workspace != nil && m.workspace.Prompts != nil {
		if p, ok := m.workspace.Prompts[name]; ok {
			path := p.Path
			if path != "" && m.workspacePath != "" {
				path = config.ResolvePromptPath(m.workspacePath, path)
			}
			return path, p.Content, true
		}
	}
	if m.global != nil && m.global.Prompts != nil {
		if p, ok := m.global.Prompts[name]; ok {
			path := p.Path
			if path != "" && m.globalPath != "" {
				path = config.ResolvePromptPath(m.globalPath, path)
			}
			return path, p.Content, true
		}
	}
	return "", "", false
}

// reviewCmd builds the ralph review subcommand (T6.1). Syntax: review [alias], review -f <path>, or stdin.
// Flags: --file/-f, --report, --prompt-output, --apply, --yes/-y, --verbose/-v, --quiet/-q, --log-level, --ai-cmd, --ai-cmd-alias; --config from root.
func reviewCmd() *cobra.Command {
	var (
		filePath     string
		reportPath   string
		promptOutput string
		apply        bool
		yes          bool
		verbose      bool
		quiet        bool
		logLevel     string
		noStream     bool
		aiCmd        string
		aiCmdAlias   string
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

			// Resolve AI command for review (same precedence as run: config, env overlay; no default).
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
			if directCmd == "" {
				directCmd = eff.Loop.AICmd
			}
			aliasName := aiCmdAlias
			if aliasName == "" && overlay != nil && overlay.AICmdAlias != nil {
				aliasName = *overlay.AICmdAlias
			}
			if aliasName == "" {
				aliasName = eff.Loop.AICmdAlias
			}
			command, ok := config.ResolveAICommand(eff, directCmd, aliasName)
			if !ok {
				fmt.Fprintln(os.Stderr, "ralph review: AI command not resolved (set --ai-cmd or --ai-cmd-alias, or config (loop.ai_cmd / loop.ai_cmd_alias) or env)")
				os.Exit(2)
			}

			// Source path for defaulting revision output when --apply without --prompt-output (file or alias with path).
			var sourcePath string
			if filePath != "" {
				sourcePath = filePath
			} else if opts.Alias != "" {
				if p, _, ok := provider.PromptByName(opts.Alias); ok && p != "" {
					if !filepath.IsAbs(p) {
						sourcePath = filepath.Join(cwd, p)
					} else {
						sourcePath = p
					}
				}
			}
			// Non-interactive when stdin is not a TTY (O009/R003, O010; exit 2 if confirmation would be needed and --yes not set).
			stdinStat, _ := os.Stdin.Stat()
			nonInteractive := (stdinStat.Mode() & os.ModeCharDevice) == 0

			// Resolve effective log level and streaming: same precedence as run (quiet→error, no-stream→no AI output).
			effLogLevel := ""
			if quiet && !verbose {
				effLogLevel = "error"
			}
			if verbose {
				effLogLevel = "debug"
			}
			if logLevel != "" {
				effLogLevel = logLevel
			}
			streaming := eff.Loop.Streaming
			if quiet && !verbose {
				streaming = false
			}
			if noStream {
				streaming = false
			}
			var streamWriter io.Writer
			if streaming {
				streamWriter = os.Stdout
			}

			runOpts := review.RunOptions{
				ReportPath:       reportPath,
				PromptOutputPath: promptOutput,
				SourcePath:       sourcePath,
				WorkingDir:       cwd,
				Apply:            apply,
				Yes:              yes,
				NonInteractive:   nonInteractive,
				Verbose:          verbose,
				Quiet:            quiet,
				LogLevel:         effLogLevel,
				StreamWriter:     streamWriter,
				Command:          command,
				Cwd:              cwd,
				Env:              os.Environ(),
				TimeoutSec:       eff.Loop.TimeoutSeconds,
			}
			code, err := review.Run(content, runOpts)
			if err != nil {
				if review.IsExit2(err) {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(2)
				}
				return err
			}
			os.Exit(code)
			return nil
		},
	}
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Read prompt from this file (mutually exclusive with alias and stdin)")
	cmd.Flags().StringVar(&reportPath, "report", "", "Write report directory to this path (default: ./ralph-review, a directory in the current working directory); AI creates result.json, summary.md, original.md, revision.md, diff.md inside it")
	cmd.Flags().StringVar(&promptOutput, "prompt-output", "", "When using --apply, write revision to this path (required when prompt is from stdin)")
	cmd.Flags().BoolVar(&apply, "apply", false, "Write suggested revision to a file")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Non-interactive apply: do not prompt for confirmation")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Increase verbosity (log level debug)")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Minimal output: log level error-only; do not show AI command output")
	cmd.Flags().StringVar(&logLevel, "log-level", "", "Log level: debug, info, warn, error (overrides shortcuts)")
	cmd.Flags().BoolVar(&noStream, "no-stream", false, "Do not show AI command output in the terminal")
	cmd.Flags().StringVar(&aiCmd, "ai-cmd", "", "Direct AI command string for this review")
	cmd.Flags().StringVar(&aiCmdAlias, "ai-cmd-alias", "", "AI command alias name from config for this review")
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
