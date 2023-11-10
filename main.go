package main

import (
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var (
	visited = sync.Map{}
)

const startURL = "https://parserdigital.com/"

func main() {
	crawl(startURL, 0)
}

func crawl(url string, depth int) {
	if depth > 2 {
		return
	}

	if _, ok := visited.Load(url); ok {
		return
	}
	visited.Store(url, struct{}{})

	fmt.Printf("%s%s[visited]\n", strings.Repeat("\t", depth), url)

	links := getLinks(url)
	/*if len(links) == 0 {
		return
	}*/

	var wg sync.WaitGroup

	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			if !strings.HasPrefix(link, startURL) {
				//fmt.Printf("%s%s\n", strings.Repeat("\t", depth), url)
				return
			}
			crawl(link, depth+1)
		}(link)
	}

	wg.Wait()
}

func getLinks(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", urlStr, err)
		return nil
	}
	defer resp.Body.Close()

	links := extractLinks(resp)
	if err != nil {
		fmt.Printf("Error extracting links from %s: %v\n", urlStr, err)
		return nil
	}

	return links
}

func extractLinks(resp *http.Response) []string {
	var links []string

	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return links
		}
		token := tokenizer.Token()
		if token.DataAtom == atom.A {
			if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						link, err := resolveURL(resp.Request.URL, attr.Val)
						if err == nil {
							links = append(links, link)
						}
					}
				}
			}
		}
	}
}

func resolveURL(base *url.URL, href string) (string, error) {
	relURL, err := url.Parse(href)
	if err != nil {
		return "", err
	}

	resolvedURL := base.ResolveReference(relURL)
	return resolvedURL.String(), nil
}
