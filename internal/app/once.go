package app

import (
	"context"

	"github.com/kristyancarvalho/tux-letter/internal/logx"
)

func (a *App) Once(ctx context.Context) error {
	logx.Info("starting one-shot run", "sources", len(a.cfg.Sources.URLs))
	logx.Warn("collection pipeline not yet wired", "stage", "foundation")
	return nil
}

func (a *App) SourcesTest(ctx context.Context) error {
	logx.Info("testing sources", "count", len(a.cfg.Sources.URLs))
	for _, u := range a.cfg.Sources.URLs {
		logx.Info("source", "url", u)
	}
	return nil
}
