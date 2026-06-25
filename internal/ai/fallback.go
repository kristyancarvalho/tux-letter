package ai

import (
	"fmt"
	"strings"

	"github.com/kristyancarvalho/tux-letter/internal/article"
	"github.com/kristyancarvalho/tux-letter/internal/config"
)

func FallbackArticle(cfg config.NewsletterConfig, articles []article.Article) Article {
	title := strings.TrimSpace(cfg.Title)
	if title == "" {
		title = "Tux Letter"
	}

	bundle := BuildBundle(cfg, articles)
	if len(bundle.Sources) == 0 {
		return Article{Title: title, Subtitle: "automated fallback digest"}
	}

	grouped := make(map[string][]BundleSource)
	var order []string
	for _, s := range bundle.Sources {
		name := strings.TrimSpace(s.Source)
		if name == "" {
			name = "other sources"
		}
		if _, ok := grouped[name]; !ok {
			order = append(order, name)
		}
		grouped[name] = append(grouped[name], s)
	}

	sections := make([]Section, 0, len(order))
	refs := make([]Reference, 0, len(bundle.Sources))
	uniqueSources := make(map[string]bool)
	for _, name := range order {
		items := grouped[name]
		paragraphs := make([]string, 0, len(items))
		for _, s := range items {
			paragraphs = append(paragraphs, fallbackParagraph(s))
			refs = append(refs, Reference{ID: s.ID, Title: s.Title, Source: s.Source, URL: s.URL})
			if strings.TrimSpace(s.Source) != "" {
				uniqueSources[s.Source] = true
			}
		}
		sections = append(sections, Section{Heading: name, Paragraphs: paragraphs})
	}

	summary := fmt.Sprintf("Automated digest of %d articles from %d sources. AI synthesis was unavailable, so this edition links the collected stories directly, grouped by source.",
		len(bundle.Sources), len(uniqueSources))

	return Article{
		Title:    title,
		Subtitle: "automated fallback digest",
		Summary:  summary,
		Body:     sections,
		Sources:  refs,
	}
}

func fallbackParagraph(s BundleSource) string {
	para := s.Title
	body := limitRunes(strings.TrimSpace(s.Content), 360)
	if body != "" && body != s.Title {
		para = s.Title + " — " + body
	}
	return fmt.Sprintf("%s [%d]", para, s.ID)
}
