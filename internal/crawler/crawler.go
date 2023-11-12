package crawler

import (
	"fmt"
	"net/url"
	"strings"
	"web-crawler/internal/analizer"
	"web-crawler/internal/fetcher"
)

func Run(url string, depth int) {
	quit := make(chan struct{})
	visitedUrls := map[string]bool{url: false}

	go crawl(url, depth, quit, visitedUrls)

	<-quit

}

func crawl(urlToVisit string, depth int, quit chan struct{}, visitedUrls map[string]bool) {
	if depth == 0 {
		quit <- struct{}{}
		return
	}

	if _, ok := visitedUrls[urlToVisit]; ok {
		quit <- struct{}{}
		return
	}

	body, err := fetcher.Fetch(urlToVisit)
	if err != nil {
		visitedUrls[urlToVisit] = true
		fmt.Printf("%s%s[failed]\n", strings.Repeat("\t", depth), urlToVisit)
		quit <- struct{}{}
		return
	}

	urlParsed, err := url.Parse(urlToVisit)
	if err != nil {
		visitedUrls[urlToVisit] = true
		fmt.Printf("%s%s[failed]\n", strings.Repeat("\t", depth), urlToVisit)
		quit <- struct{}{}
		return
	}

	urls := analizer.GetLinks(urlParsed, body)

	visitedUrls[urlToVisit] = true
	fmt.Printf("%s%s[visited]\n", strings.Repeat("\t", depth), urlToVisit)

	childrenQuit := make(chan struct{})
	for _, childrenUrl := range urls {
		go crawl(childrenUrl, depth-1, childrenQuit, visitedUrls)
		<-childrenQuit
	}
	quit <- struct{}{}

}
