package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/kristyancarvalho/tux-letter/internal/article"
	"github.com/kristyancarvalho/tux-letter/internal/config"
	"github.com/kristyancarvalho/tux-letter/internal/logx"
)

type Generator struct {
	completer Completer
	models    []string
}

func NewGenerator(completer Completer, models []string) *Generator {
	return &Generator{completer: completer, models: models}
}

func (g *Generator) Generate(ctx context.Context, cfg config.NewsletterConfig, articles []article.Article) (Digest, error) {
	if len(articles) == 0 {
		return Digest{}, fmt.Errorf("no articles to summarize")
	}
	if g.completer == nil || len(g.models) == 0 {
		return Digest{}, fmt.Errorf("no ai models configured")
	}

	prompt := BuildPrompt(cfg, articles)
	var lastErr error

	for _, model := range g.models {
		model = strings.TrimSpace(model)
		if model == "" {
			continue
		}
		digest, err := g.tryModel(ctx, model, prompt)
		if err == nil {
			logx.Info("ai digest generated", "model", model, "items", len(digest.Items))
			return digest, nil
		}
		lastErr = err
		logx.Warn("ai model failed, trying next", "model", model, "error", err)
	}

	return Digest{}, fmt.Errorf("all ai models failed: %w", lastErr)
}

func (g *Generator) tryModel(ctx context.Context, model string, prompt Prompt) (Digest, error) {
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		raw, err := g.completer.Complete(ctx, model, prompt)
		if err != nil {
			return Digest{}, err
		}
		digest, err := ParseDigest(raw)
		if err == nil {
			if err = digest.Validate(); err == nil {
				return digest, nil
			}
		}
		lastErr = err
		logx.Warn("ai output invalid, retrying", "model", model, "attempt", attempt+1, "error", err)
	}
	return Digest{}, fmt.Errorf("invalid output after retry: %w", lastErr)
}
