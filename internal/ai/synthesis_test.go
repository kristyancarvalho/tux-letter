package ai

import (
	"strings"
	"testing"

	"github.com/kristyancarvalho/tux-letter/internal/article"
	"github.com/kristyancarvalho/tux-letter/internal/config"
)

func TestBuildBundleCarriesArticleContent(t *testing.T) {
	articles := []article.Article{
		{Title: "Kernel 6.20", Source: "lwn", URL: "https://lwn.net/a", Excerpt: "short excerpt", Content: "Full extracted article body describing the scheduler rework and new power management defaults."},
		{Title: "Mesa 25", Source: "phoronix", URL: "https://phoronix.com/b", Excerpt: "fallback excerpt"},
	}
	bundle := BuildBundle(config.NewsletterConfig{}, articles)
	if len(bundle.Sources) != 2 {
		t.Fatalf("expected 2 bundle sources, got %d", len(bundle.Sources))
	}
	if !strings.Contains(bundle.Sources[0].Content, "scheduler rework") {
		t.Errorf("expected extracted Content to be sent in the bundle, got %q", bundle.Sources[0].Content)
	}
	if bundle.Sources[1].Content != "fallback excerpt" {
		t.Errorf("expected excerpt fallback when Content is empty, got %q", bundle.Sources[1].Content)
	}
}

func TestBuildBundleCapsContentSize(t *testing.T) {
	huge := strings.Repeat("a", maxBundleContentRunes+500)
	bundle := BuildBundle(config.NewsletterConfig{}, []article.Article{
		{Title: "t", Source: "s", URL: "https://x/y", Content: huge},
	})
	if got := len([]rune(bundle.Sources[0].Content)); got > maxBundleContentRunes+1 {
		t.Errorf("expected content capped near %d runes, got %d", maxBundleContentRunes, got)
	}
}

func TestSystemPromptRequestsDepth(t *testing.T) {
	bundle := BuildBundle(config.NewsletterConfig{}, sampleArticles())
	p := BuildPrompt(config.NewsletterConfig{}, bundle)
	for _, want := range []string{
		"concrete technical details",
		"why each development matters",
		"Connect related stories",
		"at least two substantial paragraphs",
	} {
		if !strings.Contains(p.System, want) {
			t.Errorf("system prompt missing depth instruction %q", want)
		}
	}
}
