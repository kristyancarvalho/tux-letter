package ai

import (
	"fmt"
	"strings"

	"github.com/kristyancarvalho/tux-letter/internal/article"
	"github.com/kristyancarvalho/tux-letter/internal/config"
)

func FallbackDigest(cfg config.NewsletterConfig, articles []article.Article) Digest {
	title := strings.TrimSpace(cfg.Title)
	if title == "" {
		title = "Tux Letter"
	}

	sources := make(map[string]bool)
	items := make([]Item, 0, len(articles))
	for _, a := range articles {
		if strings.TrimSpace(a.Title) == "" || strings.TrimSpace(a.URL) == "" {
			continue
		}
		if a.Source != "" {
			sources[a.Source] = true
		}
		items = append(items, Item{
			Title:   a.Title,
			Source:  a.Source,
			URL:     a.URL,
			Summary: a.Excerpt,
			Tags:    fallbackTags(a),
		})
	}

	summary := fmt.Sprintf("Automated digest of %d articles from %d sources. AI summarization was unavailable, so entries are listed with their original metadata.",
		len(items), len(sources))

	return Digest{Title: title, Summary: summary, Items: items}
}

func fallbackTags(a article.Article) []string {
	if strings.TrimSpace(a.Source) == "" {
		return nil
	}
	return []string{strings.ToLower(a.Source)}
}
