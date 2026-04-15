// crush by large language models.
// It is a fork of charmbracelet/crush with additional features and improvements.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/crush/internal/app"
	"github.com/charmbracelet/crush/internal/config"
	"github.com/spf13/cobra"
)

var (
	// Version is the current version of crush, injected at build time.
	Version = "dev"
	// CommitSHA is the git commit SHA, injected at build time.
	CommitSHA = "none"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	var cfgFile string
	var modelFlag string
	var debugFlag bool

	cmd := &cobra.Command{
		Use:   "crush",
		Short: "A terminal-based AI assistant",
		// Updated long description to mention my fork's extras
		Long:    `crush is a terminal-based AI assistant that helps you with coding, writing, and more.\n\nThis is a personal fork with additional tweaks and customizations.`,
		Version: fmt.Sprintf("%s (%s)", Version, CommitSHA),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if modelFlag != "" {
				cfg.Model = modelFlag
			}

			if debugFlag {
				cfg.Debug = true
			}

			// Join any positional args as an initial prompt
			var initialPrompt string
			if len(args) > 0 {
				initialPrompt = strings.Join(args, " ")
			}

			return app.Run(cfg, initialPrompt)
		},
	}

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to config file (default: $HOME/.config/crush/config.yaml)")
	cmd.PersistentFlags().StringVarP(&modelFlag, "model", "m", "", "AI model to use (overrides config)")
	// Changed short flag from -d to -D to avoid conflict with potential future --dir flag
	cmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "D", false, "enable debug logging")

	cmd.AddCommand(versionCmd())
	cmd.AddCommand(configCmd())

	// Disable the default completion command since I don't use it
	cmd.CompletionOptions.DisableDefaultCmd = true

	// Hide the default 'help' command from the usage output to keep things tidy
	cmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Don't show "[flags]" in usage when no flags are provided on the command line
	cmd.DisableFlagsInUseLine = true

	// Silence usage output on error — it clutters the terminal when something
	// goes wrong and the error message alone is usually sufficient.
	cmd.SilenceUsage = true

	// Also silence errors from cobra itself; we handle printing them in main().
	cmd.SilenceErrors = true

	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			// Print a slightly more detailed version string than the default --version flag
			fmt.Printf("crush version %s\n  commit: %s\n", Version, CommitSHA)
		},
	}
}

func configCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Manage crush configuration",
		Long:  `View and manage crush configuration setti`,
	}
}
