package render

import (
	"fmt"
	"strings"
	"time"
)

func RenderText(n Newsletter) (string, error) {
	title := strings.TrimSpace(n.Title)
	if title == "" {
		title = "Tux Letter"
	}

	var b strings.Builder
	rule := strings.Repeat("=", 60)

	b.WriteString(rule + "\n")
	b.WriteString("  " + strings.ToUpper(title) + "\n")
	b.WriteString("  " + Tagline + "\n")
	b.WriteString("  " + SubTagline + "\n")
	b.WriteString(rule + "\n\n")

	if s := strings.TrimSpace(n.Summary); s != "" {
		b.WriteString("// BRIEFING\n")
		b.WriteString(wrap(s, 72) + "\n\n")
	}

	if len(n.Items) == 0 {
		b.WriteString("No new dispatches in this cycle.\n\n")
	}
	for i, it := range n.Items {
		source := strings.TrimSpace(it.Source)
		if source == "" {
			source = "unknown"
		}
		b.WriteString(fmt.Sprintf("[%02d] %s\n", i+1, strings.ToLower(source)))
		b.WriteString("     " + strings.TrimSpace(it.Title) + "\n")
		if s := strings.TrimSpace(it.Summary); s != "" {
			b.WriteString(indent(wrap(s, 68), "     ") + "\n")
		}
		if w := strings.TrimSpace(it.WhyItMatters); w != "" {
			b.WriteString(indent(wrap("why it matters > "+w, 68), "     ") + "\n")
		}
		if tags := cleanTags(it.Tags); len(tags) > 0 {
			b.WriteString("     #" + strings.Join(tags, " #") + "\n")
		}
		if u := strings.TrimSpace(it.URL); u != "" {
			b.WriteString("     " + u + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(rule + "\n")
	stamp := "unknown"
	if !n.GeneratedAt.IsZero() {
		stamp = n.GeneratedAt.UTC().Format(time.RFC1123)
	}
	b.WriteString(fmt.Sprintf("tux-letter // %d articles // %d sources\n", len(n.Items), len(n.Sources)))
	b.WriteString("generated " + stamp + "\n")
	b.WriteString(RepositoryURL + "\n")
	b.WriteString("generated locally by tux-letter\n")

	return b.String(), nil
}

func cleanTags(tags []string) []string {
	var out []string
	for _, tag := range tags {
		if s := strings.TrimSpace(tag); s != "" {
			out = append(out, strings.ToLower(s))
		}
	}
	return out
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

func indent(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = prefix + l
	}
	return strings.Join(lines, "\n")
}
