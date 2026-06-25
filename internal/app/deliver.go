package app

import (
	"fmt"
	"time"

	"github.com/kristyancarvalho/tux-letter/internal/ai"
	"github.com/kristyancarvalho/tux-letter/internal/logx"
	"github.com/kristyancarvalho/tux-letter/internal/mail"
	"github.com/kristyancarvalho/tux-letter/internal/render"
)

func (a *App) toNewsletter(d ai.Digest, sources []string, at time.Time) render.Newsletter {
	title := d.Title
	if title == "" {
		title = a.cfg.Newsletter.Title
	}
	items := make([]render.Item, 0, len(d.Items))
	for _, it := range d.Items {
		items = append(items, render.Item{
			Title:        it.Title,
			Source:       it.Source,
			URL:          it.URL,
			Summary:      it.Summary,
			WhyItMatters: it.WhyItMatters,
			Tags:         it.Tags,
		})
	}
	return render.Newsletter{
		Title:       title,
		Summary:     d.Summary,
		Items:       items,
		GeneratedAt: at,
		Sources:     sources,
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
		logx.Info("email disabled, wrote newsletter to stdout", "items", len(n.Items))
		return nil
	}

	settings, err := mail.Resolve(a.cfg.Email)
	if err != nil {
		return err
	}
	subject := n.Title
	if subject == "" {
		subject = "Tux Letter"
	}
	if err := mail.Send(settings, mail.Message{Subject: subject, Text: textOut, HTML: htmlOut}); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	logx.Info("newsletter delivered by email", "recipients", len(settings.To), "items", len(n.Items))
	return nil
}
