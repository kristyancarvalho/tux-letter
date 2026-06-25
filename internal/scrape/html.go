package scrape

import (
	"strings"

	"golang.org/x/net/html"
)

func parse(body []byte) (*html.Node, error) {
	return html.Parse(strings.NewReader(string(body)))
}

func attr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if strings.EqualFold(a.Key, key) {
			return a.Val
		}
	}
	return ""
}

func hasAttr(n *html.Node, key string) bool {
	for _, a := range n.Attr {
		if strings.EqualFold(a.Key, key) {
			return true
		}
	}
	return false
}

func forEach(root *html.Node, fn func(*html.Node)) {
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		fn(n)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(root)
}

func elements(root *html.Node, tag string) []*html.Node {
	var out []*html.Node
	forEach(root, func(n *html.Node) {
		if n.Type == html.ElementNode && strings.EqualFold(n.Data, tag) {
			out = append(out, n)
		}
	})
	return out
}

func text(n *html.Node) string {
	var b strings.Builder
	forEach(n, func(c *html.Node) {
		if c.Type == html.TextNode {
			b.WriteString(c.Data)
		}
	})
	return collapse(b.String())
}

func collapse(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func metaContent(root *html.Node) map[string]string {
	out := make(map[string]string)
	for _, m := range elements(root, "meta") {
		content := strings.TrimSpace(attr(m, "content"))
		if content == "" {
			continue
		}
		if p := strings.ToLower(strings.TrimSpace(attr(m, "property"))); p != "" {
			if _, ok := out[p]; !ok {
				out[p] = content
			}
		}
		if name := strings.ToLower(strings.TrimSpace(attr(m, "name"))); name != "" {
			if _, ok := out[name]; !ok {
				out[name] = content
			}
		}
	}
	return out
}

func firstHeadingLink(n *html.Node) (string, string) {
	for _, tag := range []string{"h1", "h2", "h3", "h4"} {
		for _, h := range elements(n, tag) {
			for _, a := range elements(h, "a") {
				href := strings.TrimSpace(attr(a, "href"))
				title := text(a)
				if title == "" {
					title = text(h)
				}
				if href != "" {
					return title, href
				}
			}
			if t := text(h); t != "" {
				if a := firstAnchor(n); a != nil {
					return t, strings.TrimSpace(attr(a, "href"))
				}
			}
		}
	}
	return "", ""
}

func firstAnchor(n *html.Node) *html.Node {
	for _, a := range elements(n, "a") {
		if strings.TrimSpace(attr(a, "href")) != "" {
			return a
		}
	}
	return nil
}
