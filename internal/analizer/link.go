package analizer

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"log"
	"net/url"
	"strings"
)

func GetLinks(url *url.URL, value string) []string {
	var links []string
	tokenizer := html.NewTokenizer(strings.NewReader(value))
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
						link, err := resolveURL(url, attr.Val)
						if err != nil {
							log.Println("resolve url failed", err.Error())
							continue
						}
						if link.Fragment != "" {
							log.Println("url contains fragment")
							continue
						}
						if link.Scheme != "http" && link.Scheme != "https" {
							log.Println("url link is not http or https")
							continue
						}

						links = append(links, link.String())

					}
				}
			}
		}
	}
}

func resolveURL(base *url.URL, href string) (*url.URL, error) {
	relURL, err := url.Parse(href)
	if err != nil {
		return nil, err
	}

	resolvedURL := base.ResolveReference(relURL)
	return resolvedURL, nil
}
