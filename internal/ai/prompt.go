package ai

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kristyancarvalho/tux-letter/internal/article"
	"github.com/kristyancarvalho/tux-letter/internal/config"
)

type Prompt struct {
	System string
	User   string
}

const maxBundleContentRunes = 6000

func BuildBundle(cfg config.NewsletterConfig, articles []article.Article) SourceBundle {
	language := strings.TrimSpace(cfg.Language)
	if language == "" {
		language = "en"
	}
	style := strings.TrimSpace(cfg.Tone)
	if style == "" {
		style = "natural, technical, concise, editorial"
	}

	sources := make([]BundleSource, 0, len(articles))
	id := 0
	for _, a := range articles {
		title := strings.TrimSpace(a.Title)
		url := strings.TrimSpace(a.URL)
		if title == "" || url == "" {
			continue
		}
		id++
		content := strings.TrimSpace(a.Content)
		if content == "" {
			content = strings.TrimSpace(a.Excerpt)
		}
		published := ""
		if !a.Published.IsZero() {
			published = a.Published.UTC().Format(time.RFC3339)
		}
		sources = append(sources, BundleSource{
			ID:          id,
			Source:      strings.TrimSpace(a.Source),
			Title:       title,
			URL:         url,
			PublishedAt: published,
			Content:     limitRunes(content, maxBundleContentRunes),
		})
	}
	return SourceBundle{Language: language, Style: style, Sources: sources}
}

func BuildPrompt(cfg config.NewsletterConfig, bundle SourceBundle) Prompt {
	title := strings.TrimSpace(cfg.Title)
	if title == "" {
		title = "Tux Letter"
	}

	payload, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		payload = []byte("{}")
	}

	var user strings.Builder
	user.WriteString("Newsletter title: " + title + "\n")
	user.WriteString("Source bundle (JSON):\n")
	user.Write(payload)

	return Prompt{System: systemPrompt(title, bundle), User: user.String()}
}

func systemPrompt(title string, bundle SourceBundle) string {
	var s strings.Builder
	s.WriteString("You are the editor of a Linux and open-source intelligence newsletter called " + title + ".\n")
	s.WriteString("Write in language code: " + bundle.Language + ".\n")
	s.WriteString("Tone: " + bundle.Style + ".\n")
	s.WriteString("Read ALL provided source articles in the bundle and synthesize them into ONE cohesive newsletter article.\n")
	s.WriteString("Group related developments together and explain the broader context for the reader.\n")
	s.WriteString("Do NOT write one mini-summary per source. Do NOT produce a list of cards. Write a flowing editorial article.\n")
	s.WriteString("Cite sources inline using numeric markers like [1], [2], [3] wherever you use information from a source.\n")
	s.WriteString("Use only the source ids present in the bundle. Never invent sources or cite ids that are not in the bundle.\n")
	s.WriteString("Preserve factual uncertainty; do not overstate or fabricate details.\n")
	s.WriteString("Respond ONLY with a single valid JSON object, no markdown, matching this shape:\n")
	s.WriteString(`{"title":string,"subtitle":string,"summary":string,"body":[{"heading":string,"paragraphs":[string]}],"sources":[{"id":number,"title":string,"source":string,"url":string}]}` + "\n")
	s.WriteString("Every output must contain at least one inline [n] citation in the body. The sources array must list the sources you cited, using their bundle id, title, source and url.\n")
	return s.String()
}

func repairPrompt(p Prompt, reason error) Prompt {
	var s strings.Builder
	s.WriteString(p.System)
	s.WriteString("\nThe previous response was rejected: " + reason.Error() + ".\n")
	s.WriteString("Return ONLY a valid JSON object in the required shape, with at least one inline [n] citation that matches a provided source id, and a sources array listing the cited sources.\n")
	return Prompt{System: s.String(), User: p.User}
}

func ParseArticle(raw string) (Article, error) {
	jsonText := extractJSON(raw)
	if jsonText == "" {
		return Article{}, fmt.Errorf("no json object found in model output")
	}
	var a Article
	if err := json.Unmarshal([]byte(jsonText), &a); err != nil {
		return Article{}, fmt.Errorf("decode article json: %w", err)
	}
	return a, nil
}

func extractJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "```") {
		raw = strings.TrimPrefix(raw, "```")
		if i := strings.IndexByte(raw, '\n'); i >= 0 {
			raw = raw[i+1:]
		}
		if i := strings.LastIndex(raw, "```"); i >= 0 {
			raw = raw[:i]
		}
		raw = strings.TrimSpace(raw)
	}
	start := strings.IndexByte(raw, '{')
	end := strings.LastIndexByte(raw, '}')
	if start < 0 || end < 0 || end < start {
		return ""
	}
	return raw[start : end+1]
}

func limitRunes(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return strings.TrimSpace(string(r[:max])) + "…"
}
