// crush is a terminal-based AI assistant powered by large language models.
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
		Use:     "crush",
		Short:   "A terminal-based AI assistant",
		Long:    `crush is a terminal-based AI assistant that helps you with coding, writing, and more.`,
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

	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("crush version %s (commit: %s)\n", Version, CommitSHA)
		},
	}
}

func configCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Manage crush configuration",
		Long:  `View and manage crush configuration settings.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use cfgFile from parent flag if provided; fall back to default discovery.
			// Accessing the flag here via cmd.Root() so the -c flag is respected
			// even when running `crush config`.
			cfgFile, _ := cmd.Root().PersistentFlags().GetString("config")
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			fmt.Printf("Config file: %s\n", cfg.Path)
			fmt.Printf("Model:       %s\n", cfg.Model)
			fmt.Printf("Debug:       %v\n", cfg.Debug)
			// Print version so it's easy to confirm which build is in use
			fmt.Printf("Version:     %s (%s)\n", Version, CommitSHA)
			return nil
		},
	}
}
