package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tms",
	Short: "TMS - Tmux Session Manager",
	Long: `TMS is a powerful tmux session manager that helps you manage complex tmux sessions with ease.

It supports:
- Multiple sessions and windows with panes
- YAML and JSON configuration files
- Named and local workspaces
- Templates for reusable configurations`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
