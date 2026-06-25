package article

import "testing"

func TestNormalizeCanonicalizesURLAndWhitespace(t *testing.T) {
	a := Normalize(Article{
		Title:  "  Kernel   6.20   Released ",
		URL:    "HTTP://Example.com/News/?utm_source=rss&id=9#frag",
		Author: "  Jane  Maintainer ",
	})
	if a.Title != "Kernel 6.20 Released" {
		t.Errorf("title not collapsed: %q", a.Title)
	}
	if a.Author != "Jane Maintainer" {
		t.Errorf("author not collapsed: %q", a.Author)
	}
	if a.URL != "http://example.com/News?id=9" {
		t.Errorf("url not canonicalized: %q", a.URL)
	}
	if a.Hash == "" {
		t.Error("expected hash to be computed")
	}
}

func TestHashStableAcrossTrackingParams(t *testing.T) {
	h1 := Hash(Article{Title: "Same Title", URL: "https://x.example/post?utm_source=a"})
	h2 := Hash(Article{Title: "Same Title", URL: "https://x.example/post/"})
	if h1 != h2 {
		t.Errorf("expected equal hashes after canonicalization, got %s vs %s", h1, h2)
	}
	h3 := Hash(Article{Title: "Other", URL: "https://x.example/post"})
	if h1 == h3 {
		t.Error("expected different hash for different title")
	}
}

func TestDedupeAcrossSources(t *testing.T) {
	items := NormalizeAll([]Article{
		{Title: "A", URL: "https://a.example/one", Source: "a"},
		{Title: "A", URL: "https://a.example/one/", Source: "b"},
		{Title: "B", URL: "https://a.example/two?utm_medium=x", Source: "a"},
		{Title: "B again", URL: "https://a.example/two", Source: "c"},
		{Title: "C", URL: "https://a.example/three", Source: "a"},
	})
	out := Dedupe(items)
	if len(out) != 3 {
		t.Fatalf("expected 3 unique articles, got %d", len(out))
	}
}

func TestNormalizeAllDropsEmpty(t *testing.T) {
	out := NormalizeAll([]Article{
		{Title: "", URL: ""},
		{Title: "Valid", URL: "https://x.example/p"},
	})
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
}
