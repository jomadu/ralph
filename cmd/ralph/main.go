package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/maxdunn/ralph/internal/config"
	"github.com/maxdunn/ralph/internal/review"
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
	root.AddCommand(versionCmd())
	root.AddCommand(reviewCmd())
	root.AddCommand(listCmd())
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
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

			var provider review.PromptProvider
			if configPath != "" {
				path := configPath
				if !filepath.IsAbs(path) {
					path = filepath.Join(cwd, path)
				}
				layer, err := config.ReadLayer(path)
				if err != nil {
					return fmt.Errorf("config file: %w", err)
				}
				if layer == nil {
					return fmt.Errorf("config file not found: %s", path)
				}
				provider = &review.FileLayerProvider{Layer: layer}
			} else {
				global, workspace, err := config.LoadGlobalAndWorkspace(os.Getenv, cwd)
				if err != nil {
					return fmt.Errorf("config: %w", err)
				}
				provider = &mergedPromptProvider{global: global, workspace: workspace}
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
	resolved, err := resolveConfigForList(cmd)
	if err != nil {
		return err
	}
	resolved = config.ResolvedWithBuiltins(resolved)
	printPrompts(resolved.Prompts)
	printAliases(resolved.Aliases)
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
			resolved, err := resolveConfigForList(cmd)
			if err != nil {
				return err
			}
			printPrompts(resolved.Prompts)
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
			resolved, err := resolveConfigForList(cmd)
			if err != nil {
				return err
			}
			resolved = config.ResolvedWithBuiltins(resolved)
			printAliases(resolved.Aliases)
			return nil
		},
	}
}

// resolveConfigForList loads config the same way as review: --config explicit path, or global + workspace.
func resolveConfigForList(cmd *cobra.Command) (*config.Resolved, error) {
	configPath, _ := cmd.Root().PersistentFlags().GetString("config")
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if configPath != "" {
		path := configPath
		if !filepath.IsAbs(path) {
			path = filepath.Join(cwd, path)
		}
		layer, err := config.ReadLayer(path)
		if err != nil {
			return nil, fmt.Errorf("config file: %w", err)
		}
		if layer == nil {
			return nil, fmt.Errorf("config file not found: %s", path)
		}
		return config.MergeLayers(nil, layer), nil
	}
	global, workspace, err := config.LoadGlobalAndWorkspace(os.Getenv, cwd)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return config.MergeLayers(global, workspace), nil
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
