package app

import (
	"context"

	"github.com/kristyancarvalho/tux-letter/internal/logx"
)

func (a *App) Serve(ctx context.Context) error {
	logx.Info("service mode starting", "schedule", a.cfg.App.Schedule, "timezone", a.cfg.App.Timezone)
	if a.cfg.App.RunOnStart {
		if err := a.Once(ctx); err != nil {
			logx.Error("initial run failed", "error", err)
		}
	}
	logx.Warn("scheduler not yet wired", "stage", "foundation")
	<-ctx.Done()
	logx.Info("service mode stopped")
	return nil
}
