package app

import (
	"context"
	"fmt"
	"time"

	"github.com/kristyancarvalho/tux-letter/internal/logx"
	"github.com/kristyancarvalho/tux-letter/internal/schedule"
)

func (a *App) Serve(ctx context.Context) error {
	loc, err := schedule.LoadLocation(a.cfg.App.Timezone)
	if err != nil {
		logx.Warn("invalid timezone, falling back to UTC", "timezone", a.cfg.App.Timezone, "error", err)
		loc = time.UTC
	}
	if _, _, err := schedule.ParseTimeOfDay(a.cfg.App.Schedule); err != nil {
		return fmt.Errorf("invalid schedule: %w", err)
	}

	logx.Info("service mode starting", "schedule", a.cfg.App.Schedule, "timezone", loc.String())

	if a.cfg.App.RunOnStart {
		if err := a.Once(ctx); err != nil {
			logx.Error("initial run failed", "error", err)
		}
	}

	for {
		next, err := schedule.NextRun(time.Now(), a.cfg.App.Schedule, loc)
		if err != nil {
			return fmt.Errorf("compute next run: %w", err)
		}
		wait := time.Until(next)
		logx.Info("next run scheduled", "at", next.Format(time.RFC1123), "in", wait.Round(time.Second).String())

		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			logx.Info("service mode stopped")
			return nil
		case <-timer.C:
			logx.Info("scheduled run starting")
			if err := a.Once(ctx); err != nil {
				logx.Error("scheduled run failed", "error", err)
			}
		}
	}
}
