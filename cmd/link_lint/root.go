package linklint

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "link-lint",
	Short: "Link-Lint - A lightweight Go-based CLI tool  to detect dead links in a website",
	Long:  `Link-Lint is a lightweight Go-based CLI tool that scans websites and detects broken or dead links.It crawls web pages, checks the status of each URL, and reports links that return errors (such as 404 or 500), helping developers maintain healthy and reliable websites.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
