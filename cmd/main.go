package main

import (
	"flag"
	"log"
	"strings"
	"sync"
)

//TODO verify mailto: and tel in href

/*
Creating a web crawler involves several considerations to ensure its effectiveness, efficiency, and ethical use. Here are some key considerations to take into account:

Respect Robots.txt: Check and respect the robots.txt file of a website. It provides guidelines on which parts of the site are open to crawling and which are off-limits. Ignoring this file may lead to legal issues or being blocked by the website.

Rate Limiting: Implement rate limiting to avoid overwhelming servers and to be respectful of a website's resources. Crawling too quickly or making too many requests in a short time can lead to being blocked.

Politeness: Be a polite crawler by adhering to the rules specified in the robots.txt file and including a meaningful user agent in your HTTP requests. This helps websites identify your crawler and contact you if there are any issues.

Crawl Depth: Decide on the depth of your crawl. You might want to focus on specific pages or go deeper into the site. This depends on your goals and the nature of the content you are looking for.

Duplicate Content Handling: Implement mechanisms to avoid crawling and storing duplicate content. Duplicate content not only wastes resources but can also lead to incorrect analysis of data.

Handling Dynamic Content: Some websites use JavaScript to load content dynamically. Ensure that your crawler can handle such dynamic content by using tools like headless browsers or by analyzing the JavaScript on the page.

User-Agent Rotation: Rotate your User-Agent header to mimic different web browsers or devices. This helps prevent being easily identified and blocked by websites.

Proxy Usage: Consider using proxies to distribute requests across multiple IP addresses. This can help avoid IP-based rate limiting and enhance your ability to crawl without interruption.

Handling Disruptions: Develop mechanisms to handle disruptions such as network errors, server timeouts, or changes in website structure. This ensures your crawler is robust and can recover from unexpected situations.

Legal and Ethical Compliance: Ensure that your web crawling activities comply with legal regulations and ethical standards. Avoid crawling sensitive information or violating privacy laws.

Monitoring and Logging: Implement monitoring and logging to keep track of your crawler's activities. This helps in identifying issues, analyzing performance, and understanding the behavior of the websites you are crawling.

Respectful Crawling Hours: Avoid crawling during peak hours or times when a website experiences heavy traffic. This ensures that your crawler doesn't adversely impact the user experience for regular visitors.

Handling Session-based Content: If a website relies on sessions or cookies for content access, ensure your crawler can manage and maintain sessions as needed.

Follow Redirects: Handle redirects appropriately to ensure that your crawler doesn't miss content due to URL redirection.

Remember to always check the terms of service of the website you are crawling and adhere to ethical standards to ensure a positive and collaborative relationship with the websites you interact with.
*/

// TODO Check protocol (https http)
// TODO regex for not https://google.com/sarasa#saras2 and query params
// TODO look redirect.
// TODO runner (optional)

var (
	visited = sync.Map{}
)

const startURL = "https://google.com/"

func main() {
	urlToCrawl := flag.String("url", "", "url to be crawled")
	crawlDelay := flag.Uint("delay", 500, "delay in milliseconds between pages")
	depth := flag.Uint("delay", 3, "depth of the crawler")
	flag.Parse()

	if urlToCrawl == nil || strings.TrimSpace(*urlToCrawl) == "" {
		log.Fatalln("url should be define and not be empty")
	}

	if strings.HasPrefix("http://", *urlToCrawl) || strings.HasPrefix("https://", *urlToCrawl) {
		log.Fatalln("url should start with http or https")
	}

	if *depth > 10 || *depth < 1 {
		log.Fatalln("depth should be between 1 and 10")
	}

	if *crawlDelay > 5000 || *crawlDelay < 0 {
		log.Fatalln("depth should be between 1 and 10")
	}

	log.Println()
}

/*
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
	}

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
*/
