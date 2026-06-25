package scrape

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/kristyancarvalho/tux-letter/internal/article"
	"golang.org/x/net/html"
)

var articleTypes = map[string]bool{
	"article":     true,
	"newsarticle": true,
	"blogposting": true,
	"techarticle": true,
	"report":      true,
	"webpage":     true,
}

func fromJSONLD(root *html.Node, baseURL string) []article.Article {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil
	}

	var out []article.Article
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
		out = append(out, walkJSONLD(node, base)...)
	}
	return out
}

func walkJSONLD(node interface{}, base *url.URL) []article.Article {
	var out []article.Article
	switch v := node.(type) {
	case []interface{}:
		for _, item := range v {
			out = append(out, walkJSONLD(item, base)...)
		}
	case map[string]interface{}:
		if graph, ok := v["@graph"]; ok {
			out = append(out, walkJSONLD(graph, base)...)
		}
		if list, ok := v["itemListElement"]; ok {
			out = append(out, walkItemList(list, base)...)
		}
		if a, ok := articleFromMap(v, base); ok {
			out = append(out, a)
		}
	}
	return out
}

func walkItemList(node interface{}, base *url.URL) []article.Article {
	var out []article.Article
	items, ok := node.([]interface{})
	if !ok {
		return nil
	}
	for _, raw := range items {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if inner, ok := m["item"].(map[string]interface{}); ok {
			if a, ok := articleFromMap(inner, base); ok {
				out = append(out, a)
				continue
			}
		}
		title := firstString(m, "name", "headline")
		link := firstString(m, "url", "@id")
		if a, ok := buildArticle(title, link, "", "", base); ok {
			out = append(out, a)
		}
	}
	return out
}

func articleFromMap(m map[string]interface{}, base *url.URL) (article.Article, bool) {
	if !isArticleType(m["@type"]) {
		return article.Article{}, false
	}
	title := firstString(m, "headline", "name", "title")
	link := firstString(m, "url", "@id")
	if link == "" {
		if main, ok := m["mainEntityOfPage"].(map[string]interface{}); ok {
			link = firstString(main, "@id", "url")
		}
	}
	desc := firstString(m, "description", "abstract")
	author := authorName(m["author"])
	return buildArticle(title, link, desc, author, base)
}

func buildArticle(title, link, desc, author string, base *url.URL) (article.Article, bool) {
	title = collapse(title)
	link = strings.TrimSpace(link)
	if title == "" || link == "" {
		return article.Article{}, false
	}
	ref, err := url.Parse(link)
	if err != nil {
		return article.Article{}, false
	}
	resolved := base.ResolveReference(ref)
	return article.Article{
		Title:   title,
		URL:     resolved.String(),
		Excerpt: collapse(desc),
		Author:  collapse(author),
	}, true
}

func isArticleType(v interface{}) bool {
	switch t := v.(type) {
	case string:
		return articleTypes[strings.ToLower(strings.TrimSpace(t))]
	case []interface{}:
		for _, item := range t {
			if s, ok := item.(string); ok && articleTypes[strings.ToLower(strings.TrimSpace(s))] {
				return true
			}
		}
	}
	return false
}

func authorName(v interface{}) string {
	switch a := v.(type) {
	case string:
		return a
	case map[string]interface{}:
		return firstString(a, "name")
	case []interface{}:
		if len(a) > 0 {
			return authorName(a[0])
		}
	}
	return ""
}

func firstString(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if s, ok := m[k].(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	return ""
}
