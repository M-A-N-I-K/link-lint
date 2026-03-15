package linklint

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

var urls []string
var total_pages_visited int

const (
	ColorRed   = "\x1b[31m"
	ColorGreen = "\x1b[32m"
	ColorReset = "\x1b[0m"
)

func ScrapeWebsite(url string) {
	total_pages_visited = total_pages_visited + 1
	client := &http.Client{}

	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	CheckIfDeadPage(url, resp.StatusCode)
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
		fmt.Println("scrpaing url", url)
		ScrapeWebsite(url)
	}
}

func CheckIfDeadPage(url string, statusCode int) bool {
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
			url = url + n.Attr[0].Val
			urls = append(urls, url)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		urls = append(urls, ParseHTML(c, url)...)
	}
	return urls
}
