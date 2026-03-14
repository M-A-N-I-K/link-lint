package linklint

import (
	linklint "link-lint/pkg/link_lint"

	"github.com/spf13/cobra"
)

var scrapeCmd = &cobra.Command{
	Use:     "scrape",
	Aliases: []string{"scrape"},
	Short:   "Scrape URL",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		linklint.ScrapeWebsite(args[0])
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)
}
