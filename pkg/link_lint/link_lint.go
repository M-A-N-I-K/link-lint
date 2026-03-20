package linklint

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var urls []string
var visitedPages sync.Map

const (
	ColorRed   = "\x1b[31m"
	ColorGreen = "\x1b[32m"
	ColorReset = "\x1b[0m"
)

func ScrapeWebsite(url string, ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	if _, ok := visitedPages.Load(url); ok {
		fmt.Printf("Page %s already scraped\n", url)
	}
	visitedPages.Store("url", true)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(ColorRed+url, err, ColorReset)
		return
	}
	dead_page, message := DeadPage(url, resp.StatusCode)
	if dead_page {
		ch <- message
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
		wg.Add(1)
		if _, ok := visitedPages.Load(url); ok {
			continue
		}
		go ScrapeWebsite(url, ch, wg)
	}
}

func DeadPage(url string, statusCode int) (bool, string) {
	if statusCode > 400 {
		message := url + " is dead page"
		fmt.Println(ColorRed + message + ColorReset)
		urls = append(urls, message)
		return true, message
	}
	if statusCode > 300 && statusCode < 400 {
		message := url + " is a redirected page"
		fmt.Println(message)
		urls = append(urls, message)
		return true, message
	}
	return false, ""
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
