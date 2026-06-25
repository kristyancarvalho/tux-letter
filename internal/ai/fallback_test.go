package ai

import (
	"strings"
	"testing"

	"github.com/kristyancarvalho/tux-letter/internal/article"
	"github.com/kristyancarvalho/tux-letter/internal/config"
)

func fallbackArticles() []article.Article {
	return []article.Article{
		{Title: "Kernel 6.20 Released", Source: "lwn", URL: "https://lwn.net/a", Excerpt: "Broad hardware support lands."},
		{Title: "Scheduler Rework Merged", Source: "lwn", URL: "https://lwn.net/b", Excerpt: "Latency improvements."},
		{Title: "Mesa 25 Ships", Source: "phoronix", URL: "https://phoronix.com/c", Excerpt: "New Vulkan driver."},
	}
}

func TestFallbackGroupsArticlesBySource(t *testing.T) {
	a := FallbackArticle(config.NewsletterConfig{Title: "Tux Letter"}, fallbackArticles())

	if len(a.Body) != 2 {
		t.Fatalf("expected 2 sections (one per source), got %d", len(a.Body))
	}

	bySource := map[string]int{}
	for _, sec := range a.Body {
		bySource[sec.Heading] = len(sec.Paragraphs)
	}
	if bySource["lwn"] != 2 {
		t.Errorf("expected lwn section to group 2 articles, got %d", bySource["lwn"])
	}
	if bySource["phoronix"] != 1 {
		t.Errorf("expected phoronix section to have 1 article, got %d", bySource["phoronix"])
	}
}

func TestFallbackHasInlineCitationsAndFooter(t *testing.T) {
	a := FallbackArticle(config.NewsletterConfig{}, fallbackArticles())

	cited := CitedIDs(a)
	if len(cited) != 3 {
		t.Fatalf("expected every article cited inline, got cited=%v", cited)
	}
	if len(a.Sources) != 3 {
		t.Fatalf("expected 3 footer references, got %d", len(a.Sources))
	}

	footerIDs := map[int]bool{}
	for _, r := range a.Sources {
		footerIDs[r.ID] = true
	}
	for _, id := range cited {
		if !footerIDs[id] {
			t.Errorf("cited id %d missing from footer references", id)
		}
	}

	if strings.TrimSpace(a.Summary) == "" {
		t.Error("expected a fallback summary paragraph")
	}
}

func TestFallbackIsNotPerCardLayout(t *testing.T) {
	a := FallbackArticle(config.NewsletterConfig{}, fallbackArticles())
	if len(a.Body) >= len(fallbackArticles()) {
		t.Errorf("fallback should group sources, not render one section per article (sections=%d, articles=%d)",
			len(a.Body), len(fallbackArticles()))
	}
	for _, sec := range a.Body {
		for _, p := range sec.Paragraphs {
			if strings.Contains(strings.ToLower(p), "why it matters") || strings.Contains(p, "why_it_matters") {
				t.Errorf("fallback paragraph must not use the old why_it_matters card field: %q", p)
			}
		}
	}
}

func TestFallbackEmptyArticles(t *testing.T) {
	a := FallbackArticle(config.NewsletterConfig{Title: "Tux Letter"}, nil)
	if a.Title != "Tux Letter" {
		t.Errorf("expected title preserved, got %q", a.Title)
	}
	if len(a.Body) != 0 || len(a.Sources) != 0 {
		t.Error("expected empty body and sources for no articles")
	}
}
