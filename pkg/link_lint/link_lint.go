package linklint

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var urls []string
var visitedPages = map[string]struct{}{}

const (
	ColorRed   = "\x1b[31m"
	ColorGreen = "\x1b[32m"
	ColorReset = "\x1b[0m"
)

func ScrapeWebsite(url string) {
	if _, ok := visitedPages[url]; ok {
		fmt.Printf("Page %s already scraped\n", url)
	}
	visitedPages[url] = struct{}{}
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(ColorRed+url, err, ColorReset)
		return
	}
	if DeadPage(url, resp.StatusCode) {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	doc, err := html.Parse(strings.NewReader(string(body)))

	if err != nil {
		panic(err)
	}
	urls = ParseHTML(doc, url)
	for _, url := range urls {
		// fmt.Println("scraping url", url)
		if _, ok := visitedPages[url]; ok {
			// fmt.Printf("Page %s already scraped\n", url)
			continue
		}
		ScrapeWebsite(url)
	}
}

func DeadPage(url string, statusCode int) bool {
	if statusCode > 400 {
		message := url + " is dead page"
		fmt.Println(ColorRed + message + ColorReset)
		urls = append(urls, message)
		return true
	}
	return false
}

func ParseHTML(n *html.Node, url string) []string {
	var urls []string
	if n.Type == html.ElementNode {
		if n.Data == "a" && n.Attr[0].Key == "href" {
			// fmt.Print("found : ", n.Attr[0].Val, "\n")

			if strings.HasPrefix(n.Attr[0].Val, "http") {
				url = n.Attr[0].Val
			} else {
				parts := strings.Split(url, "/")
				url = parts[0] + "//" + parts[2] + n.Attr[0].Val
			}
			urls = append(urls, url)
			return urls
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		urls = append(urls, ParseHTML(c, url)...)
	}
	return urls
}
