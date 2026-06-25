package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kristyancarvalho/tux-letter/internal/ai"
	"github.com/kristyancarvalho/tux-letter/internal/article"
	"github.com/kristyancarvalho/tux-letter/internal/collect"
	"github.com/kristyancarvalho/tux-letter/internal/fetch"
	"github.com/kristyancarvalho/tux-letter/internal/logx"
	"github.com/kristyancarvalho/tux-letter/internal/scrape"
	"github.com/kristyancarvalho/tux-letter/internal/source"
	"github.com/kristyancarvalho/tux-letter/internal/state"
)

func (a *App) Once(ctx context.Context) error {
	sources, errs := source.Build(a.cfg.Sources.URLs)
	for _, err := range errs {
		logx.Warn("skipping invalid source", "error", err)
	}
	if len(sources) == 0 {
		return fmt.Errorf("no valid sources configured")
	}
	logx.Info("starting one-shot run", "sources", len(sources))

	now := time.Now()
	st := a.openState()
	if st != nil {
		st.MarkAttempt(now)
	}

	fetcher := fetch.New(a.cfg.Fetch.Timeout.Duration(), a.cfg.Fetch.UserAgent)
	result := collect.Sources(ctx, fetcher, sources, a.cfg.Fetch)

	if st != nil {
		for name, status := range result.Statuses {
			st.SetSourceStatus(name, status, now)
		}
	}

	articles := result.Articles
	if st != nil {
		articles = st.FilterNew(articles)
	}
	articles = capTotal(articles, a.cfg.Fetch.MaxTotalArticles)
	logx.Info("collected articles", "new", len(articles))

	if len(articles) == 0 {
		logx.Info("no new articles, nothing to deliver")
		a.saveState(st, now)
		return nil
	}

	articles = a.enrichContent(ctx, fetcher, articles)

	art := a.summarize(ctx, articles)
	newsletter := a.toNewsletter(art, now)
	if err := a.deliver(newsletter); err != nil {
		return err
	}

	a.saveState(st, now)
	return nil
}

func (a *App) summarize(ctx context.Context, articles []article.Article) ai.Article {
	apiKey := os.Getenv(a.cfg.OpenRouter.APIKeyEnv)
	if apiKey == "" || len(a.cfg.OpenRouter.Models) == 0 {
		logx.Warn("openrouter not configured, using non-ai fallback digest")
		return ai.FallbackArticle(a.cfg.Newsletter, articles)
	}
	client := ai.NewOpenRouter(a.cfg.OpenRouter.BaseURL, apiKey, a.cfg.Fetch.Timeout.Duration())
	gen := ai.NewGenerator(client, a.cfg.OpenRouter.Models)
	art, err := gen.Generate(ctx, a.cfg.Newsletter, articles)
	if err != nil {
		logx.Warn("ai generation failed, using non-ai fallback digest", "error", err)
		return ai.FallbackArticle(a.cfg.Newsletter, articles)
	}
	return art
}

func (a *App) enrichContent(ctx context.Context, f fetch.Fetcher, articles []article.Article) []article.Article {
	for i := range articles {
		url := strings.TrimSpace(articles[i].URL)
		if url == "" {
			continue
		}
		resp, err := f.Get(ctx, url)
		if err != nil {
			logx.Warn("article fetch failed, using excerpt", "url", url, "error", err)
			continue
		}
		if content := scrape.ExtractContent(resp.Body); content != "" {
			articles[i].Content = content
		}
	}
	return articles
}

func (a *App) openState() *state.Store {
	path := filepath.Join(a.cfg.StateDir(), "state.json")
	st, err := state.Open(path, a.cfg.State.MaxSeenArticles)
	if err != nil {
		logx.Warn("could not open state, continuing without persistence", "path", path, "error", err)
		return nil
	}
	return st
}

func (a *App) saveState(st *state.Store, at time.Time) {
	if st == nil {
		return
	}
	st.MarkSuccess(at)
	if err := st.Save(); err != nil {
		logx.Warn("could not save state", "error", err)
	}
}

func capTotal(items []article.Article, max int) []article.Article {
	if max > 0 && len(items) > max {
		return items[:max]
	}
	return items
}

func (a *App) SourcesTest(ctx context.Context) error {
	sources, errs := source.Build(a.cfg.Sources.URLs)
	for _, err := range errs {
		logx.Warn("skipping invalid source", "error", err)
	}
	if len(sources) == 0 {
		return fmt.Errorf("no valid sources configured")
	}
	logx.Info("testing sources", "count", len(sources))

	fetcher := fetch.New(a.cfg.Fetch.Timeout.Duration(), a.cfg.Fetch.UserAgent)
	result := collect.Sources(ctx, fetcher, sources, a.cfg.Fetch)

	for _, s := range sources {
		logx.Info("source", "name", s.Name, "host", s.Host, "url", s.URL, "status", result.Statuses[s.Name])
	}
	logx.Info("extraction complete", "articles", len(result.Articles))
	for _, art := range result.Articles {
		logx.Info("article", "source", art.Source, "title", art.Title, "url", art.URL)
	}
	return nil
}
