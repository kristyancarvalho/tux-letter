package app

import (
	"context"

	"github.com/kristyancarvalho/tux-letter/internal/logx"
	"github.com/kristyancarvalho/tux-letter/internal/source"
)

func (a *App) Once(ctx context.Context) error {
	sources, errs := source.Build(a.cfg.Sources.URLs)
	for _, err := range errs {
		logx.Warn("skipping invalid source", "error", err)
	}
	logx.Info("starting one-shot run", "sources", len(sources))
	logx.Warn("collection pipeline not yet wired", "stage", "foundation")
	return nil
}

func (a *App) SourcesTest(ctx context.Context) error {
	sources, errs := source.Build(a.cfg.Sources.URLs)
	for _, err := range errs {
		logx.Warn("skipping invalid source", "error", err)
	}
	logx.Info("testing sources", "count", len(sources))
	for _, s := range sources {
		logx.Info("source", "name", s.Name, "host", s.Host, "url", s.URL)
	}
	return nil
}
