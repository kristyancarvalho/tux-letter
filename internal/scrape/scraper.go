package scrape

import (
	"net/url"
	"strings"

	"github.com/kristyancarvalho/tux-letter/internal/article"
	"golang.org/x/net/html"
)

func Extract(source, baseURL string, body []byte) []article.Article {
	root, err := parse(body)
	if err != nil {
		return nil
	}

	var candidates []article.Article
	candidates = append(candidates, fromJSONLD(root, baseURL)...)
	if a, ok := fromOpenGraph(root, baseURL); ok {
		candidates = append(candidates, a)
	}
	candidates = append(candidates, fromArticleBlocks(root, baseURL)...)
	candidates = append(candidates, fromHeadingLinks(root, baseURL)...)
	candidates = append(candidates, rankAnchors(root, baseURL)...)

	return finalize(source, baseURL, candidates)
}

func fromOpenGraph(root *html.Node, baseURL string) (article.Article, bool) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return article.Article{}, false
	}
	meta := metaContent(root)

	title := firstNonEmpty(meta["og:title"], meta["twitter:title"])
	link := firstNonEmpty(meta["og:url"], canonicalLink(root))
	desc := firstNonEmpty(meta["og:description"], meta["twitter:description"], meta["description"])
	author := firstNonEmpty(meta["article:author"], meta["author"])

	if title == "" || link == "" {
		return article.Article{}, false
	}
	ref, err := url.Parse(strings.TrimSpace(link))
	if err != nil {
		return article.Article{}, false
	}
	resolved := base.ResolveReference(ref)
	return article.Article{
		Title:   collapse(title),
		URL:     resolved.String(),
		Excerpt: collapse(desc),
		Author:  collapse(author),
	}, true
}

func canonicalLink(root *html.Node) string {
	for _, l := range elements(root, "link") {
		if strings.EqualFold(strings.TrimSpace(attr(l, "rel")), "canonical") {
			return strings.TrimSpace(attr(l, "href"))
		}
	}
	return ""
}

func fromArticleBlocks(root *html.Node, baseURL string) []article.Article {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil
	}
	var out []article.Article
	for _, block := range elements(root, "article") {
		title, href := firstHeadingLink(block)
		if href == "" {
			if a := firstAnchor(block); a != nil {
				href = strings.TrimSpace(attr(a, "href"))
				if title == "" {
					title = collapse(text(a))
				}
			}
		}
		if a, ok := buildArticle(title, href, "", "", base); ok {
			out = append(out, a)
		}
	}
	return out
}

func fromHeadingLinks(root *html.Node, baseURL string) []article.Article {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil
	}
	var out []article.Article
	for _, tag := range []string{"h1", "h2", "h3"} {
		for _, h := range elements(root, tag) {
			a := firstAnchor(h)
			if a == nil {
				continue
			}
			href := strings.TrimSpace(attr(a, "href"))
			title := collapse(text(a))
			if title == "" {
				title = collapse(text(h))
			}
			if art, ok := buildArticle(title, href, "", "", base); ok {
				out = append(out, art)
			}
		}
	}
	return out
}

func finalize(source, baseURL string, candidates []article.Article) []article.Article {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil
	}

	var out []article.Article
	index := make(map[string]int)

	for _, c := range candidates {
		ref, err := url.Parse(strings.TrimSpace(c.URL))
		if err != nil {
			continue
		}
		resolved := base.ResolveReference(ref)
		if !acceptableLink(base, resolved) {
			continue
		}
		if collapse(c.Title) == "" {
			continue
		}
		key := normalizeKey(resolved)
		c.URL = resolved.String()
		c.Title = collapse(c.Title)
		if c.Source == "" {
			c.Source = source
		}

		if i, ok := index[key]; ok {
			out[i] = merge(out[i], c)
			continue
		}
		index[key] = len(out)
		out = append(out, c)
	}
	return out
}

func normalizeKey(u *url.URL) string {
	u2 := *u
	u2.Fragment = ""
	u2.Host = strings.TrimPrefix(strings.ToLower(u2.Host), "www.")
	u2.Path = strings.TrimRight(u2.Path, "/")
	if u2.Path == "" {
		u2.Path = "/"
	}
	return u2.String()
}

func merge(a, b article.Article) article.Article {
	if len([]rune(b.Title)) > len([]rune(a.Title)) {
		a.Title = b.Title
	}
	if a.Excerpt == "" {
		a.Excerpt = b.Excerpt
	}
	if a.Author == "" {
		a.Author = b.Author
	}
	if a.Source == "" {
		a.Source = b.Source
	}
	return a
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
