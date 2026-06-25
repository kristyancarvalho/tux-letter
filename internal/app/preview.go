package app

import (
	"context"
	"fmt"
	"os"

	"github.com/kristyancarvalho/tux-letter/internal/logx"
	"github.com/kristyancarvalho/tux-letter/internal/render"
)

func (a *App) Preview(ctx context.Context, output string) error {
	n := render.Demo()
	if brand := a.cfg.Newsletter.Title; brand != "" {
		n.Brand = brand
	}

	htmlOut, err := render.RenderHTML(n)
	if err != nil {
		return err
	}

	if output == "" {
		fmt.Print(htmlOut)
		return nil
	}
	if err := os.WriteFile(output, []byte(htmlOut), 0o644); err != nil {
		return fmt.Errorf("writing preview: %w", err)
	}
	logx.Info("wrote newsletter preview", "path", output)
	return nil
}
