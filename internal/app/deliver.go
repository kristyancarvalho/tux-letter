package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/kristyancarvalho/tux-letter/internal/ai"
	"github.com/kristyancarvalho/tux-letter/internal/logx"
	"github.com/kristyancarvalho/tux-letter/internal/mail"
	"github.com/kristyancarvalho/tux-letter/internal/render"
)

func (a *App) toNewsletter(art ai.Article, at time.Time) render.Newsletter {
	brand := strings.TrimSpace(a.cfg.Newsletter.Title)
	if brand == "" {
		brand = "Tux Letter"
	}

	sections := make([]render.Section, 0, len(art.Body))
	for _, s := range art.Body {
		sections = append(sections, render.Section{Heading: s.Heading, Paragraphs: s.Paragraphs})
	}

	refs := make([]render.Reference, 0, len(art.Sources))
	for _, r := range art.Sources {
		refs = append(refs, render.Reference{ID: r.ID, Title: r.Title, Source: r.Source, URL: r.URL})
	}

	return render.Newsletter{
		Brand:       brand,
		Title:       art.Title,
		Subtitle:    art.Subtitle,
		Summary:     art.Summary,
		Sections:    sections,
		References:  refs,
		GeneratedAt: at,
	}
}

func (a *App) deliver(n render.Newsletter) error {
	htmlOut, err := render.RenderHTML(n)
	if err != nil {
		return fmt.Errorf("render html: %w", err)
	}
	textOut, err := render.RenderText(n)
	if err != nil {
		return fmt.Errorf("render text: %w", err)
	}

	if !a.cfg.Email.Enabled {
		fmt.Println(textOut)
		logx.Info("email disabled, wrote newsletter to stdout", "sources", len(n.References))
		return nil
	}

	settings, err := mail.Resolve(a.cfg.Email)
	if err != nil {
		return err
	}
	subject := strings.TrimSpace(n.Title)
	if subject == "" {
		subject = strings.TrimSpace(n.Brand)
	}
	if subject == "" {
		subject = "Tux Letter"
	}
	if err := mail.Send(settings, mail.Message{Subject: subject, Text: textOut, HTML: htmlOut}); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	logx.Info("newsletter delivered by email", "recipients", len(settings.To), "sources", len(n.References))
	return nil
}
