package linklint

import (
	"fmt"
	linklint "link-lint/pkg/link_lint"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var scrapeCmd = &cobra.Command{
	Use:     "scrape",
	Aliases: []string{"scrape"},
	Short:   "Scrape URL",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ch := make(chan string)
		start := time.Now()

		var wg sync.WaitGroup

		wg.Add(1)

		go linklint.ScrapeWebsite(args[0], ch, &wg)

		go func() {
			wg.Wait()
			close(ch)
		}()

		// for result := range ch {
		// 	fmt.Println(result)
		// }

		elapsed := time.Since(start)
		fmt.Printf("Scraping took %s ", elapsed)
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)
}
