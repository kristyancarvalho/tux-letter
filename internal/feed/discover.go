package feed

import (
	"io"
	"strings"

	"github.com/kristyancarvalho/tux-letter/internal/urlx"
	"golang.org/x/net/html"
)

var feedTypes = map[string]bool{
	"application/rss+xml":   true,
	"application/atom+xml":  true,
	"application/feed+json": true,
}

func Discover(baseURL string, body io.Reader) []string {
	doc, err := html.Parse(body)
	if err != nil {
		return nil
	}

	var found []string
	seen := make(map[string]bool)

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			var rel, typ, href string
			for _, attr := range n.Attr {
				switch strings.ToLower(attr.Key) {
				case "rel":
					rel = strings.ToLower(attr.Val)
				case "type":
					typ = strings.ToLower(strings.TrimSpace(attr.Val))
				case "href":
					href = attr.Val
				}
			}
			if strings.Contains(rel, "alternate") && feedTypes[typ] && href != "" {
				if resolved, err := urlx.Resolve(baseURL, href); err == nil && !seen[resolved] {
					seen[resolved] = true
					found = append(found, resolved)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return found
}

func LooksLikeFeed(contentType string, body []byte) bool {
	ct := strings.ToLower(contentType)
	if strings.Contains(ct, "xml") || strings.Contains(ct, "rss") || strings.Contains(ct, "atom") {
		return true
	}
	prefix := strings.TrimSpace(strings.ToLower(string(body)))
	if len(prefix) > 512 {
		prefix = prefix[:512]
	}
	return strings.Contains(prefix, "<rss") || strings.Contains(prefix, "<feed") || strings.Contains(prefix, "<rdf")
}
