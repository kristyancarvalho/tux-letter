package scrape

import (
	"strings"
	"testing"
)

func TestExtractContentFromArticle(t *testing.T) {
	body := load(t, "article_page.html")
	content := ExtractContent(body)

	if content == "" {
		t.Fatal("expected extracted content, got empty")
	}
	for _, want := range []string{
		"broad hardware enablement",
		"long-awaited filesystem fix",
		"reduce idle power draw",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("content missing %q", want)
		}
	}
	for _, junk := range []string{
		"Login",
		"Subscribe to our newsletter",
		"Advertisement",
		"Privacy policy",
		"Cookie settings",
	} {
		if strings.Contains(content, junk) {
			t.Errorf("content leaked non-article text %q", junk)
		}
	}
}

func TestExtractContentEmptyBody(t *testing.T) {
	if c := ExtractContent([]byte("")); c != "" {
		t.Errorf("expected empty content for empty body, got %q", c)
	}
}

func TestExtractContentFallsBackToMeta(t *testing.T) {
	page := `<!DOCTYPE html><html><head>` +
		`<meta name="description" content="A concise meta description used when no article body is present on the page.">` +
		`</head><body><div><p>too short</p></div></body></html>`
	c := ExtractContent([]byte(page))
	if !strings.Contains(c, "concise meta description") {
		t.Errorf("expected meta description fallback, got %q", c)
	}
}
