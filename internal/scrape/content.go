package scrape

import (
	"encoding/json"
	"strings"

	"golang.org/x/net/html"
)

const (
	maxContentRunes = 6000
	minContentRunes = 200
)

var nonContentTags = map[string]bool{
	"script": true, "style": true, "nav": true, "header": true,
	"footer": true, "aside": true, "form": true, "noscript": true,
	"button": true, "svg": true, "figure": true, "template": true,
}

func ExtractContent(body []byte) string {
	root, err := parse(body)
	if err != nil {
		return ""
	}

	if c := articleContent(root); c != "" {
		return limitContent(c)
	}
	if c := jsonLDValue(root, "articleBody"); c != "" {
		return limitContent(c)
	}
	if c := mainContent(root); c != "" {
		return limitContent(c)
	}

	meta := metaContent(root)
	if d := firstNonEmpty(meta["og:description"], meta["description"], meta["twitter:description"]); d != "" {
		return limitContent(collapse(d))
	}
	if c := jsonLDValue(root, "description", "abstract"); c != "" {
		return limitContent(c)
	}
	return ""
}

func articleContent(root *html.Node) string {
	best := ""
	for _, a := range elements(root, "article") {
		t := contentText(a)
		if runeLen(t) > runeLen(best) {
			best = t
		}
	}
	if runeLen(best) >= minContentRunes {
		return best
	}
	return ""
}

func mainContent(root *html.Node) string {
	for _, m := range elements(root, "main") {
		if t := contentText(m); runeLen(t) >= minContentRunes {
			return t
		}
	}

	best := ""
	bestScore := 0
	for _, tag := range []string{"div", "section"} {
		for _, n := range elements(root, tag) {
			ps := elements(n, "p")
			if len(ps) < 3 {
				continue
			}
			var b strings.Builder
			for _, p := range ps {
				b.WriteString(contentText(p))
				b.WriteString(" ")
			}
			t := collapse(b.String())
			if score := runeLen(t); score > bestScore {
				best = t
				bestScore = score
			}
		}
	}
	if bestScore >= minContentRunes {
		return best
	}
	return ""
}

func contentText(n *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode && nonContentTags[strings.ToLower(node.Data)] {
			return
		}
		if node.Type == html.TextNode {
			b.WriteString(node.Data)
			b.WriteString(" ")
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return collapse(b.String())
}

func jsonLDValue(root *html.Node, keys ...string) string {
	for _, s := range elements(root, "script") {
		if !strings.EqualFold(strings.TrimSpace(attr(s, "type")), "application/ld+json") {
			continue
		}
		raw := strings.TrimSpace(text(s))
		if raw == "" {
			continue
		}
		var node interface{}
		if err := json.Unmarshal([]byte(raw), &node); err != nil {
			continue
		}
		if v := walkJSONLDValue(node, keys); v != "" {
			return v
		}
	}
	return ""
}

func walkJSONLDValue(node interface{}, keys []string) string {
	switch v := node.(type) {
	case []interface{}:
		for _, item := range v {
			if s := walkJSONLDValue(item, keys); s != "" {
				return s
			}
		}
	case map[string]interface{}:
		if g, ok := v["@graph"]; ok {
			if s := walkJSONLDValue(g, keys); s != "" {
				return s
			}
		}
		for _, k := range keys {
			if s, ok := v[k].(string); ok && strings.TrimSpace(s) != "" {
				return collapse(s)
			}
		}
	}
	return ""
}

func limitContent(s string) string {
	r := []rune(s)
	if len(r) <= maxContentRunes {
		return s
	}
	return strings.TrimSpace(string(r[:maxContentRunes])) + "…"
}

func runeLen(s string) int { return len([]rune(s)) }
