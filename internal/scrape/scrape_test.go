package scrape

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kristyancarvalho/tux-letter/internal/article"
)

func load(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return data
}

func byURL(items []article.Article) map[string]article.Article {
	m := make(map[string]article.Article)
	for _, a := range items {
		m[a.URL] = a
	}
	return m
}

func hasURL(items []article.Article, want string) bool {
	for _, a := range items {
		if a.URL == want {
			return true
		}
	}
	return false
}

func TestExtractJSONLDAndOpenGraph(t *testing.T) {
	body := load(t, "listing_jsonld.html")
	items := Extract("example", "https://example.com/news", body)

	if len(items) == 0 {
		t.Fatal("expected articles, got none")
	}
	for _, a := range items {
		if a.Source != "example" {
			t.Errorf("expected source to be set, got %q", a.Source)
		}
	}

	idx := byURL(items)

	wayland, ok := idx["https://example.com/2026/06/wayland-protocol-stability"]
	if !ok {
		t.Fatal("expected JSON-LD ItemList article (wayland) to be extracted")
	}
	if wayland.Title != "Wayland Protocol Reaches 1.0 Stability" {
		t.Errorf("unexpected wayland title: %q", wayland.Title)
	}
	if wayland.Author != "Alex Display" {
		t.Errorf("expected JSON-LD author, got %q", wayland.Author)
	}

	if !hasURL(items, "https://example.com/2026/06/rust-in-the-kernel-one-year") {
		t.Error("expected relative JSON-LD url to be resolved against base")
	}

	featured, ok := idx["https://example.com/2026/06/kernel-6-20-released"]
	if !ok {
		t.Fatal("expected Open Graph featured article to be extracted")
	}
	if !strings.Contains(strings.ToLower(featured.Title), "kernel") {
		t.Errorf("unexpected featured title: %q", featured.Title)
	}
}

func TestExtractFiltersJunkAndDedupes(t *testing.T) {
	body := load(t, "listing_jsonld.html")
	items := Extract("example", "https://example.com/news", body)

	for _, a := range items {
		low := strings.ToLower(a.URL)
		for _, junk := range []string{"/tag/", "/category/", "/login", "twitter.com"} {
			if strings.Contains(low, junk) {
				t.Errorf("junk link leaked into results: %s", a.URL)
			}
		}
	}

	seen := make(map[string]int)
	for _, a := range items {
		seen[a.URL]++
	}
	for u, n := range seen {
		if n > 1 {
			t.Errorf("duplicate url %s appeared %d times", u, n)
		}
	}
}

func TestExtractArticleBlocks(t *testing.T) {
	body := load(t, "listing_articles.html")
	items := Extract("distro", "https://distroweekly.example/", body)

	want := []string{
		"https://distroweekly.example/news/arch-installer-gets-graphical-mode",
		"https://distroweekly.example/news/debian-freezes-for-next-stable",
		"https://distroweekly.example/news/systemd-adds-userspace-feature",
	}
	for _, w := range want {
		if !hasURL(items, w) {
			t.Errorf("expected article %s to be extracted", w)
		}
	}

	for _, a := range items {
		low := strings.ToLower(a.URL)
		for _, junk := range []string{"/tags/", "/author/", "/subscribe", "/privacy", "/terms", "facebook.com"} {
			if strings.Contains(low, junk) {
				t.Errorf("junk link leaked: %s", a.URL)
			}
		}
	}
}

func TestExtractAnchorRanking(t *testing.T) {
	body := load(t, "listing_anchors.html")
	items := Extract("plain", "https://plainlinux.example/", body)

	want := []string{
		"https://plainlinux.example/2026/06/btrfs-gets-faster-scrubbing",
		"https://plainlinux.example/2026/06/mesa-25-ships-new-vulkan-driver",
		"https://plainlinux.example/2026/06/gnome-48-improves-power-usage",
	}
	for _, w := range want {
		if !hasURL(items, w) {
			t.Errorf("expected ranked article %s", w)
		}
	}

	for _, a := range items {
		low := strings.ToLower(a.URL)
		for _, junk := range []string{"/about", "/contact", "/tag/", "/login", "/p/2", "youtube.com", "mailto:"} {
			if strings.Contains(low, junk) {
				t.Errorf("junk link leaked: %s", a.URL)
			}
		}
	}
}

func TestExtractEmptyBody(t *testing.T) {
	if items := Extract("x", "https://x.example/", []byte("")); len(items) != 0 {
		t.Errorf("expected no items for empty body, got %d", len(items))
	}
	if items := Extract("x", "https://x.example/", []byte("<html><body><p>no links</p></body></html>")); len(items) != 0 {
		t.Errorf("expected no items for link-free page, got %d", len(items))
	}
}
