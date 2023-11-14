package analyzer

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"log"
	"net/url"
	"strings"
)

func GetLinks(baseUrl *url.URL, value string) []*url.URL {
	var links []*url.URL
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
						if strings.HasPrefix(strings.ToLower(attr.Val), "mailto") || strings.HasPrefix(strings.ToLower(attr.Val), "tel") {
							log.Println("baseUrl link is not http or https")
							continue
						}

						link, err := resolveURL(baseUrl, attr.Val)
						if err != nil {
							log.Println("resolve baseUrl failed", err.Error())
							continue
						}

						links = append(links, link)

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
