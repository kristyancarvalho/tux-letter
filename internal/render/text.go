package render

import (
	"fmt"
	"strings"
	"time"
)

func RenderText(n Newsletter) (string, error) {
	brand := strings.TrimSpace(n.Brand)
	if brand == "" {
		brand = "Tux Letter"
	}
	headline := strings.TrimSpace(n.Title)
	if headline == "" {
		headline = brand
	}
	loc := localeFor(n.Language)

	var b strings.Builder
	rule := strings.Repeat("=", 60)

	b.WriteString(rule + "\n")
	b.WriteString("  " + strings.ToUpper(brand) + "\n")
	b.WriteString("  " + loc.tagline + "\n")
	b.WriteString("  " + loc.subTagline + "\n")
	b.WriteString(rule + "\n\n")

	b.WriteString(headline + "\n")
	if s := strings.TrimSpace(n.Subtitle); s != "" {
		b.WriteString(wrap(s, 72) + "\n")
	}
	b.WriteString("\n")

	if s := strings.TrimSpace(n.Summary); s != "" {
		b.WriteString("// " + strings.ToUpper(loc.briefing) + "\n")
		b.WriteString(wrap(s, 72) + "\n\n")
	}

	if len(n.Sections) == 0 {
		b.WriteString(loc.emptyBody + "\n\n")
	}
	for _, sec := range n.Sections {
		if h := strings.TrimSpace(sec.Heading); h != "" {
			b.WriteString("// " + strings.ToUpper(h) + "\n")
		}
		for _, p := range sec.Paragraphs {
			if strings.TrimSpace(p) == "" {
				continue
			}
			b.WriteString(wrap(p, 72) + "\n\n")
		}
	}

	if len(n.References) > 0 {
		b.WriteString(strings.ToUpper(loc.sources) + "\n")
		for _, r := range n.References {
			line := fmt.Sprintf("[%d] %s", r.ID, strings.TrimSpace(r.Title))
			if src := strings.TrimSpace(r.Source); src != "" {
				line += " — " + src
			}
			b.WriteString(line + "\n")
			if u := strings.TrimSpace(r.URL); u != "" {
				b.WriteString("    " + u + "\n")
			}
		}
		b.WriteString("\n")
	}

	b.WriteString(rule + "\n")
	stamp := "unknown"
	if !n.GeneratedAt.IsZero() {
		stamp = n.GeneratedAt.UTC().Format(time.RFC1123)
	}
	b.WriteString(fmt.Sprintf("tux-letter // %d %s\n", len(n.References), loc.sourcesWord))
	b.WriteString(loc.generated + " " + stamp + "\n")
	b.WriteString(RepositoryURL + "\n")
	b.WriteString(loc.generatedBy + "\n")

	return b.String(), nil
}

func wrap(s string, width int) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}
	var b strings.Builder
	lineLen := 0
	for i, w := range words {
		if lineLen > 0 && lineLen+1+len(w) > width {
			b.WriteString("\n")
			lineLen = 0
		} else if i > 0 {
			b.WriteString(" ")
			lineLen++
		}
		b.WriteString(w)
		lineLen += len(w)
	}
	return b.String()
}
