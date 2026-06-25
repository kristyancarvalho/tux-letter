package app

import (
	"github.com/kristyancarvalho/tux-letter/internal/config"
	"github.com/kristyancarvalho/tux-letter/internal/logx"
)

type App struct {
	cfg config.Config
}

func New(cfg config.Config) *App {
	logx.SetLevel(cfg.App.LogLevel)
	return &App{cfg: cfg}
}

func (a *App) Config() config.Config { return a.cfg }
