package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "frontend-init",
	Short: "Scaffold and configure frontend projects",
	Long:  "An interactive TUI wizard to scaffold React, Vue, Svelte, Angular, and Astro projects with your preferred tooling.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
