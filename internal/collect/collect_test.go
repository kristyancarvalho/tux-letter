package collect

import (
	"context"
	"fmt"
	"testing"

	"github.com/kristyancarvalho/tux-letter/internal/config"
	"github.com/kristyancarvalho/tux-letter/internal/fetch"
	"github.com/kristyancarvalho/tux-letter/internal/source"
)

type stubFetcher struct {
	pages map[string]*fetch.Response
}

func (s *stubFetcher) Get(ctx context.Context, url string) (*fetch.Response, error) {
	if r, ok := s.pages[url]; ok {
		return r, nil
	}
	return nil, fmt.Errorf("not found: %s", url)
}

const rssBody = `<?xml version="1.0"?>
<rss version="2.0"><channel>
<title>Feedy</title>
<item><title>Kernel 6.20 Released</title><link>https://feedy.example/kernel-6-20</link></item>
<item><title>Mesa 25 Lands</title><link>https://feedy.example/mesa-25</link></item>
</channel></rss>`

const htmlBody = `<!DOCTYPE html><html><body>
<nav><a href="/login">Login</a><a href="/tag/x">x</a></nav>
<ul>
<li><a href="/2026/06/gnome-48-improves-power-usage">GNOME 48 Improves Power Usage On Laptops</a></li>
<li><a href="/2026/06/btrfs-faster-scrubbing-lands">Btrfs Faster Scrubbing Lands In The Kernel</a></li>
</ul>
</body></html>`

func TestSourcesFeedAndScrape(t *testing.T) {
	sf := &stubFetcher{pages: map[string]*fetch.Response{
		"https://feedy.example/": {URL: "https://feedy.example/", ContentType: "application/rss+xml", Body: []byte(rssBody)},
		"https://blog.example/":  {URL: "https://blog.example/", ContentType: "text/html", Body: []byte(htmlBody)},
	}}

	sources, _ := source.Build([]string{"https://feedy.example/", "https://blog.example/"})
	if len(sources) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(sources))
	}

	res := Sources(context.Background(), sf, sources, config.Default().Fetch)

	if res.Statuses["Feedy"] != "ok" && res.Statuses["feedy"] != "ok" {
		t.Errorf("expected feed source ok, got %+v", res.Statuses)
	}
	if len(res.Articles) < 4 {
		t.Fatalf("expected at least 4 articles (2 feed + 2 scraped), got %d", len(res.Articles))
	}

	var hasFeed, hasScrape bool
	for _, a := range res.Articles {
		if a.URL == "https://feedy.example/kernel-6-20" {
			hasFeed = true
		}
		if a.URL == "https://blog.example/2026/06/gnome-48-improves-power-usage" {
			hasScrape = true
		}
	}
	if !hasFeed {
		t.Error("expected feed-parsed article present")
	}
	if !hasScrape {
		t.Error("expected scraped article present")
	}
}

func TestSourcesSkipsFailingSource(t *testing.T) {
	sf := &stubFetcher{pages: map[string]*fetch.Response{
		"https://ok.example/": {URL: "https://ok.example/", ContentType: "application/rss+xml", Body: []byte(rssBody)},
	}}
	sources, _ := source.Build([]string{"https://ok.example/", "https://dead.example/"})
	res := Sources(context.Background(), sf, sources, config.Default().Fetch)
	if len(res.Articles) == 0 {
		t.Fatal("expected articles from the working source")
	}
}
