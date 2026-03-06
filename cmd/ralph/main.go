package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ralph",
	Short: "Ralph - a dumb loop that pipes prompts to AI CLIs",
}

var (
	configFlag string
)

var runCmd = &cobra.Command{
	Use:   "run [alias]",
	Short: "Run the loop with a prompt",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run: not implemented")
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
	runCmd.Flags().StringVar(&configFlag, "config", "", "Explicit config file path")
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
