package collect

import (
	"bytes"
	"context"

	"github.com/kristyancarvalho/tux-letter/internal/article"
	"github.com/kristyancarvalho/tux-letter/internal/config"
	"github.com/kristyancarvalho/tux-letter/internal/feed"
	"github.com/kristyancarvalho/tux-letter/internal/fetch"
	"github.com/kristyancarvalho/tux-letter/internal/logx"
	"github.com/kristyancarvalho/tux-letter/internal/scrape"
	"github.com/kristyancarvalho/tux-letter/internal/source"
)

type Result struct {
	Articles []article.Article
	Statuses map[string]string
}

func Sources(ctx context.Context, f fetch.Fetcher, sources []source.Source, cfg config.FetchConfig) Result {
	var all []article.Article
	statuses := make(map[string]string)

	for _, s := range sources {
		arts, err := collectOne(ctx, f, s, cfg)
		if err != nil {
			logx.Warn("source failed", "source", s.Name, "url", s.URL, "error", err)
			statuses[s.Name] = "error"
			continue
		}
		logx.Info("source collected", "source", s.Name, "articles", len(arts))
		statuses[s.Name] = "ok"
		all = append(all, arts...)
	}

	deduped := article.Dedupe(article.NormalizeAll(all))
	return Result{Articles: deduped, Statuses: statuses}
}

func collectOne(ctx context.Context, f fetch.Fetcher, s source.Source, cfg config.FetchConfig) ([]article.Article, error) {
	resp, err := f.Get(ctx, s.URL)
	if err != nil {
		return nil, err
	}

	if feed.LooksLikeFeed(resp.ContentType, resp.Body) {
		if arts, perr := feed.Parse(s.Name, resp.Body); perr == nil && len(arts) > 0 {
			return capArticles(arts, cfg.MaxArticlesPerSource), nil
		}
	}

	for _, feedURL := range feed.Discover(s.URL, bytes.NewReader(resp.Body)) {
		fr, ferr := f.Get(ctx, feedURL)
		if ferr != nil {
			logx.Warn("feed fetch failed", "source", s.Name, "feed", feedURL, "error", ferr)
			continue
		}
		if arts, perr := feed.Parse(s.Name, fr.Body); perr == nil && len(arts) > 0 {
			return capArticles(arts, cfg.MaxArticlesPerSource), nil
		}
	}

	base := resp.URL
	if base == "" {
		base = s.URL
	}
	arts := scrape.Extract(s.Name, base, resp.Body)
	return capArticles(arts, cfg.MaxArticlesPerSource), nil
}

func capArticles(items []article.Article, max int) []article.Article {
	if max > 0 && len(items) > max {
		return items[:max]
	}
	return items
}
