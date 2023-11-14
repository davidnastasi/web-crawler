package crawler

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
	"web-crawler/internal/analyzer"
	"web-crawler/internal/fetcher"
)

const (
	retryMaxAttempt = 3
	retryDelay      = 100
)

type Crawler struct {
	visitedUrls sync.Map
}

func NewCrawler() *Crawler {
	return &Crawler{}
}

func (c *Crawler) Run(baseURL string, depth, delay int) error {

	var disallowed []*url.URL
	var rules []analyzer.RobotRule
	baseUrlParsed, err := url.Parse(strings.TrimSuffix(baseURL, "/"))
	if err != nil {
		return err
	}
	body, err := fetcher.Fetch(baseUrlParsed.Scheme+"://"+baseUrlParsed.Host+"/robot.txt", "")
	if err == nil {
		rules, err = analyzer.GetRobotsContent(body)
		if err != nil {
			return err
		}
	}

	if err != nil && !errors.Is(err, fetcher.ErrNotFound) {
		return err
	}

	if len(rules) > 0 {
		for _, rule := range rules {
			for _, disallow := range rule.Disallow {
				disallowUrlParsed, err := url.Parse(strings.TrimSuffix(disallow, "/"))
				if err != nil {
					return err
				}
				disallowed = append(disallowed, baseUrlParsed.ResolveReference(disallowUrlParsed))
			}

			quit := make(chan struct{})
			go c.crawl(baseUrlParsed, depth, rule.CrawDelay, rule.UserAgent, disallowed, quit)
			<-quit
		}

	} else {
		quit := make(chan struct{})
		go c.crawl(baseUrlParsed, depth, delay, "*", disallowed, quit)
		<-quit
	}

	return nil
}

func (c *Crawler) crawl(urlToVisit *url.URL, depth, delay int, userAgent string, disallowed []*url.URL, quit chan struct{}) {
	if depth == 0 {
		quit <- struct{}{}
		return
	}

	if urlToVisit.Fragment != "" {
		quit <- struct{}{}
		return
	}

	if !isAllowed(urlToVisit, disallowed) {
		quit <- struct{}{}
		return
	}

	if _, ok := c.visitedUrls.Load(urlToVisit.String()); ok {
		quit <- struct{}{}
		return
	}

	body, err := c.fetch(urlToVisit.String(), userAgent)
	if err != nil {
		log.Printf("%s[failed]\n", urlToVisit)
		quit <- struct{}{}
		return
	}

	urls := analyzer.GetLinks(urlToVisit, body)

	c.visitedUrls.Store(urlToVisit.String(), true)
	fmt.Printf("%s[visited]\n", urlToVisit)

	childrenQuit := make(chan struct{})
	for _, childrenUrl := range urls {
		if childrenUrl.Host != urlToVisit.Host {
			continue
		}
		<-time.After(time.Duration(delay) * time.Millisecond)
		go c.crawl(childrenUrl, depth-1, delay, userAgent, disallowed, childrenQuit)
		<-childrenQuit
	}
	quit <- struct{}{}

}

func (c *Crawler) GetVisitedURLs() []string {
	var visitedUrls []string
	c.visitedUrls.Range(func(key, value any) bool {
		if value == true {
			visitedUrls = append(visitedUrls, key.(string))
		}
		return true
	})
	return visitedUrls

}

func (c *Crawler) fetch(url, agent string) (body string, err error) {
	for attempt := 1; attempt <= retryMaxAttempt; attempt++ {
		body, err = fetcher.Fetch(url, agent)
		if err != nil && errors.Is(err, fetcher.ErrTooManyRequest) {
			log.Println(err.Error(), "retrying attempt", attempt, "for url", url)
			<-time.After(time.Duration(attempt*retryDelay) * time.Millisecond)
			continue
		}
		break
	}
	return body, err
}

func isAllowed(url *url.URL, disallowed []*url.URL) bool {
	for _, dis := range disallowed {
		if url.String() == dis.String() {
			return false
		}
	}
	return true
}
