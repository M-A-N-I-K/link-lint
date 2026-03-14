package linklint

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

var urls []string

func ScrapeWebsite(url string) {
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

	ParseHTML(doc, url)
}

func CheckIfDeadPage(url string, statusCode int) bool {
	if statusCode > 400 {
		message := url + " is dead page"
		urls = append(urls, message)
		return true
	}
	return false
}

func ParseHTML(n *html.Node, url string) {
	if n.Type == html.ElementNode {
		if n.Data == "a" && n.Attr[0].Key == "href" {
			fmt.Println("checking", n.Attr[0].Val)
			ScrapeWebsite(url + n.Attr[0].Val)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ParseHTML(c, url)
	}
}
