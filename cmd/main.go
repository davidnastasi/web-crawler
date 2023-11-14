package main

import (
	"flag"
	"log"
	"strings"
	"web-crawler/internal/crawler"
)

func main() {
	urlToCrawl := flag.String("url", "", "url to be crawled")
	crawlDelay := flag.Uint("delay", 0, "delay in milliseconds between pages")
	depth := flag.Uint("depth", 3, "depth of the crawler")
	flag.Parse()

	if urlToCrawl == nil || strings.TrimSpace(*urlToCrawl) == "" {
		log.Fatalln("url should be define and not be empty")
		return
	}

	if strings.HasPrefix("http://", *urlToCrawl) || strings.HasPrefix("https://", *urlToCrawl) {
		log.Fatalln("url should start with http or https")
		return
	}

	if *depth > 10 || *depth < 1 {
		log.Fatalln("depth should be between 1 and 10")
		return
	}

	if *crawlDelay > 30000 || *crawlDelay < 0 {
		log.Fatalln("depth should be between 0 and 30000")
		return
	}

	err := crawler.NewCrawler().Run(*urlToCrawl, int(*depth), int(*crawlDelay))
	if err != nil {
		return
	}

}
