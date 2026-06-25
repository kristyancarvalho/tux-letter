package feed

import (
	"os"
	"strings"
	"testing"
	"time"
)

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return data
}

func TestParseRSS(t *testing.T) {
	articles, err := Parse("Example", readFixture(t, "rss.xml"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(articles) != 2 {
		t.Fatalf("got %d articles, want 2 (empty-title item must be skipped)", len(articles))
	}
	a := articles[0]
	if a.Title != "Kernel 6.10 Released" {
		t.Errorf("title = %q", a.Title)
	}
	if a.URL != "https://example.com/kernel-6-10" {
		t.Errorf("url = %q", a.URL)
	}
	if a.Author != "Jane Doe" {
		t.Errorf("author = %q", a.Author)
	}
	if strings.Contains(a.Excerpt, "<b>") {
		t.Errorf("excerpt should strip tags: %q", a.Excerpt)
	}
	want := time.Date(2026, 6, 24, 10, 0, 0, 0, time.UTC)
	if !a.Published.Equal(want) {
		t.Errorf("published = %v, want %v", a.Published, want)
	}
}

func TestParseAtom(t *testing.T) {
	articles, err := Parse("", readFixture(t, "atom.xml"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(articles) != 2 {
		t.Fatalf("got %d articles, want 2", len(articles))
	}
	if articles[0].URL != "https://example.org/wayland" {
		t.Errorf("url = %q", articles[0].URL)
	}
	if articles[0].Source != "Example Atom Feed" {
		t.Errorf("source = %q, want feed title fallback", articles[0].Source)
	}
	if articles[0].Author != "John Smith" {
		t.Errorf("author = %q", articles[0].Author)
	}
}

func TestDiscover(t *testing.T) {
	feeds := Discover("https://example.com/", strings.NewReader(string(readFixture(t, "page.html"))))
	if len(feeds) != 2 {
		t.Fatalf("got %d feeds, want 2: %v", len(feeds), feeds)
	}
	if feeds[0] != "https://example.com/feed.xml" {
		t.Errorf("relative feed not resolved: %q", feeds[0])
	}
	if feeds[1] != "https://example.com/atom.xml" {
		t.Errorf("absolute feed wrong: %q", feeds[1])
	}
}

func TestLooksLikeFeed(t *testing.T) {
	if !LooksLikeFeed("application/rss+xml", nil) {
		t.Error("rss content-type should look like feed")
	}
	if !LooksLikeFeed("text/plain", []byte("<?xml version='1.0'?><rss>")) {
		t.Error("rss body should look like feed")
	}
	if LooksLikeFeed("text/html", []byte("<html><body>")) {
		t.Error("html should not look like feed")
	}
}
