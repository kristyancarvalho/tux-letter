package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kristyancarvalho/tux-letter/internal/article"
	"github.com/kristyancarvalho/tux-letter/internal/config"
)

type Prompt struct {
	System string
	User   string
}

func BuildPrompt(cfg config.NewsletterConfig, articles []article.Article) Prompt {
	title := strings.TrimSpace(cfg.Title)
	if title == "" {
		title = "Tux Letter"
	}
	language := strings.TrimSpace(cfg.Language)
	if language == "" {
		language = "en"
	}
	tone := strings.TrimSpace(cfg.Tone)
	if tone == "" {
		tone = "natural, concise, technical, friendly"
	}

	var sys strings.Builder
	sys.WriteString("You are the editor of a Linux and open-source newsletter called ")
	sys.WriteString(title + ".\n")
	sys.WriteString("Write in language code: " + language + ".\n")
	sys.WriteString("Tone: " + tone + ".\n")
	sys.WriteString("Summarize the provided articles into a concise digest.\n")
	sys.WriteString("Respond ONLY with a single valid JSON object, no markdown, matching this shape:\n")
	sys.WriteString(`{"title":string,"summary":string,"items":[{"title":string,"source":string,"url":string,"summary":string,"why_it_matters":string,"tags":[string]}]}` + "\n")
	sys.WriteString("Use only the provided article URLs. Do not invent links or facts.")

	var user strings.Builder
	user.WriteString("Newsletter title: " + title + "\n")
	user.WriteString(fmt.Sprintf("Articles (%d):\n", len(articles)))
	for i, a := range articles {
		user.WriteString(fmt.Sprintf("\n[%d]\n", i+1))
		user.WriteString("title: " + a.Title + "\n")
		if a.Source != "" {
			user.WriteString("source: " + a.Source + "\n")
		}
		user.WriteString("url: " + a.URL + "\n")
		if a.Author != "" {
			user.WriteString("author: " + a.Author + "\n")
		}
		if a.Excerpt != "" {
			user.WriteString("excerpt: " + a.Excerpt + "\n")
		}
	}

	return Prompt{System: sys.String(), User: user.String()}
}

func ParseDigest(raw string) (Digest, error) {
	jsonText := extractJSON(raw)
	if jsonText == "" {
		return Digest{}, fmt.Errorf("no json object found in model output")
	}
	var d Digest
	if err := json.Unmarshal([]byte(jsonText), &d); err != nil {
		return Digest{}, fmt.Errorf("decode digest json: %w", err)
	}
	return d, nil
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
