package app

import (
	"strings"
	"testing"
	"time"

	"github.com/kristyancarvalho/tux-letter/internal/ai"
	"github.com/kristyancarvalho/tux-letter/internal/article"
	"github.com/kristyancarvalho/tux-letter/internal/config"
	"github.com/kristyancarvalho/tux-letter/internal/render"
)

func TestFallbackArticleRendersAsOneDigest(t *testing.T) {
	cfg := config.Default()
	articles := []article.Article{
		{Title: "Kernel 6.20 Released", Source: "lwn", URL: "https://lwn.net/a", Excerpt: "Broad hardware support."},
		{Title: "Scheduler Rework Merged", Source: "lwn", URL: "https://lwn.net/b", Excerpt: "Latency improvements."},
		{Title: "Mesa 25 Ships", Source: "phoronix", URL: "https://phoronix.com/c", Excerpt: "New Vulkan driver."},
	}

	art := ai.FallbackArticle(cfg.Newsletter, articles)
	a := New(cfg)
	n := a.toNewsletter(art, time.Date(2026, 6, 25, 20, 0, 0, 0, time.UTC))

	htmlOut, err := render.RenderHTML(n)
	if err != nil {
		t.Fatalf("render html: %v", err)
	}
	textOut, err := render.RenderText(n)
	if err != nil {
		t.Fatalf("render text: %v", err)
	}

	for _, want := range []string{"// sources", "https://lwn.net/a", "https://phoronix.com/c", "[1]"} {
		if !strings.Contains(htmlOut, want) {
			t.Errorf("fallback HTML missing %q", want)
		}
	}
	if strings.Contains(htmlOut, "why_it_matters") {
		t.Error("fallback HTML must not use the old why_it_matters card layout")
	}

	for _, want := range []string{"SOURCES", "https://lwn.net/b", "[3]"} {
		if !strings.Contains(textOut, want) {
			t.Errorf("fallback text missing %q", want)
		}
	}
}
