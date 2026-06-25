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

func (g *Generator) Generate(ctx context.Context, cfg config.NewsletterConfig, articles []article.Article) (Article, error) {
	if len(articles) == 0 {
		return Article{}, fmt.Errorf("no articles to summarize")
	}
	if g.completer == nil || len(g.models) == 0 {
		return Article{}, fmt.Errorf("no ai models configured")
	}

	bundle := BuildBundle(cfg, articles)
	if len(bundle.Sources) == 0 {
		return Article{}, fmt.Errorf("no usable sources in bundle")
	}
	prompt := BuildPrompt(cfg, bundle)
	var lastErr error

	for _, model := range g.models {
		model = strings.TrimSpace(model)
		if model == "" {
			continue
		}
		art, err := g.tryModel(ctx, model, prompt, bundle)
		if err == nil {
			logx.Info("ai article generated", "model", model, "sections", len(art.Body), "sources", len(art.Sources))
			return art, nil
		}
		lastErr = err
		logx.Warn("ai model failed, trying next", "model", model, "error", err)
	}

	return Article{}, fmt.Errorf("all ai models failed: %w", lastErr)
}

func (g *Generator) tryModel(ctx context.Context, model string, prompt Prompt, bundle SourceBundle) (Article, error) {
	var lastErr error
	current := prompt
	for attempt := 0; attempt < 2; attempt++ {
		raw, err := g.completer.Complete(ctx, model, current)
		if err != nil {
			return Article{}, err
		}
		art, perr := ParseArticle(raw)
		if perr == nil {
			if verr := art.ValidateAgainst(bundle); verr != nil {
				perr = verr
			} else {
				return reconcile(art, bundle), nil
			}
		}
		lastErr = perr
		logx.Warn("ai output invalid, retrying", "model", model, "attempt", attempt+1, "error", perr)
		current = repairPrompt(prompt, perr)
	}
	return Article{}, fmt.Errorf("invalid output after retry: %w", lastErr)
}

func reconcile(a Article, bundle SourceBundle) Article {
	byID := bundle.byID()
	refs := make([]Reference, 0, len(a.Sources))
	for _, id := range CitedIDs(a) {
		if bs, ok := byID[id]; ok {
			refs = append(refs, Reference{ID: bs.ID, Title: bs.Title, Source: bs.Source, URL: bs.URL})
		}
	}
	a.Sources = refs
	return a
}
