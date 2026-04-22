package parser

import (
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func ExtractLinks(body io.ReadCloser, baseURL string) []string {
	var links []string
	tokenizer := html.NewTokenizer(body)

	// URL parsing, converting absolute links to relative
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return links
	}

	for {
		tokenType := tokenizer.Next()

		// If EOF or Read error
		if tokenType == html.ErrorToken {
			return links
		}

		//
		if tokenType == html.StartTagToken {
			token := tokenizer.Token()

			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						link := cleanURL(attr.Val, parsedBaseURL)
						if link != "" {
							links = append(links, link)
						}
					}
				}
			}
		}
	}
}

// cleanURL converting link: deleting anchor(#), absolute link -> relative
func cleanURL(href string, base *url.URL) string {
	// Trimming space from the edges
	href = strings.TrimSpace(href)

	// Deleting anchors, emails and empty links
	if href == "" || strings.HasPrefix(href, "#") {
		return ""
	}

	// Trying to convert struct to url.URL
	parsedHref, err := url.Parse(href)
	if err != nil {
		return ""
	}

	// Split relative or absolute link to normal
	resolvedURL := base.ResolveReference(parsedHref)

	// Delete all != http/https
	if resolvedURL.Scheme != "http" && resolvedURL.Scheme != "https" {
		return ""
	}

	// Cleaning all anchors at the end "#"
	resolvedURL.Fragment = ""

	// Return clean string
	return resolvedURL.String()
}
