package render

import (
	"fmt"
	"html"
	"strings"
	"time"
)

func RenderHTML(n Newsletter) (string, error) {
	return RenderHTMLWithTheme(n, Cypherpunk())
}

func RenderHTMLWithTheme(n Newsletter, t Theme) (string, error) {
	title := n.Title
	if strings.TrimSpace(title) == "" {
		title = "Tux Letter"
	}

	var b strings.Builder
	b.WriteString("<!DOCTYPE html>\n")
	b.WriteString(`<html lang="en"><head><meta charset="utf-8">`)
	b.WriteString(`<meta name="viewport" content="width=device-width,initial-scale=1">`)
	b.WriteString(`<meta name="color-scheme" content="dark">`)
	b.WriteString(`<title>` + esc(title) + `</title>`)
	b.WriteString("<style>")
	b.WriteString(`a{text-decoration:none}`)
	b.WriteString(`a:hover{text-decoration:underline}`)
	b.WriteString(`@media (max-width:620px){.tl-container{width:100%!important}.tl-pad{padding:18px!important}}`)
	b.WriteString("</style></head>")

	b.WriteString(fmt.Sprintf(`<body style="margin:0;padding:0;background:%s;color:%s;font-family:%s;">`,
		t.Background, t.Text, fontBody))

	b.WriteString(fmt.Sprintf(`<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background:%s;padding:24px 12px;">`, t.Background))
	b.WriteString(`<tr><td align="center">`)
	b.WriteString(fmt.Sprintf(`<table role="presentation" class="tl-container" width="680" cellpadding="0" cellspacing="0" style="width:680px;max-width:680px;background:%s;border:1px solid %s;border-radius:10px;overflow:hidden;">`, t.Surface, t.Border))

	writeHeader(&b, t, title)
	if s := strings.TrimSpace(n.Summary); s != "" {
		writeSummary(&b, t, s)
	}
	writeItems(&b, t, n.Items)
	writeFooter(&b, t, n)

	b.WriteString(`</table></td></tr></table></body></html>`)
	return b.String(), nil
}

func writeHeader(b *strings.Builder, t Theme, title string) {
	b.WriteString(fmt.Sprintf(`<tr><td class="tl-pad" style="padding:22px 28px;background:%s;border-bottom:1px solid %s;">`, t.Surface2, t.Border))

	b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:12px;color:%s;letter-spacing:1px;">`, fontMono, t.Muted))
	b.WriteString(`<span style="color:#ff5f56">&#9679;</span> <span style="color:` + t.Warning + `">&#9679;</span> <span style="color:` + t.Accent2 + `">&#9679;</span>`)
	b.WriteString(`&nbsp;&nbsp;tux@letter:~$ cat dispatch</div>`)

	b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:30px;font-weight:700;letter-spacing:4px;color:%s;margin-top:12px;">%s</div>`,
		fontMono, t.Accent, esc(strings.ToUpper(title))))

	b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:13px;color:%s;margin-top:4px;">%s</div>`,
		fontMono, t.Accent2, esc(Tagline)))
	b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:12px;color:%s;margin-top:2px;">%s</div>`,
		fontMono, t.Muted, esc(SubTagline)))
	b.WriteString(`</td></tr>`)
}

func writeSummary(b *strings.Builder, t Theme, summary string) {
	b.WriteString(`<tr><td class="tl-pad" style="padding:24px 28px 6px 28px;">`)
	b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:11px;text-transform:uppercase;letter-spacing:2px;color:%s;margin-bottom:8px;">// briefing</div>`, fontMono, t.Accent))
	b.WriteString(fmt.Sprintf(`<div style="background:%s;border-left:3px solid %s;border-radius:6px;padding:16px 18px;font-size:15px;line-height:1.6;color:%s;">%s</div>`,
		t.Surface2, t.Accent, t.Text, esc(summary)))
	b.WriteString(`</td></tr>`)
}

func writeItems(b *strings.Builder, t Theme, items []Item) {
	b.WriteString(`<tr><td class="tl-pad" style="padding:18px 28px 8px 28px;">`)
	if len(items) == 0 {
		b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:14px;color:%s;padding:8px 0;">No new dispatches in this cycle.</div>`, fontMono, t.Muted))
		b.WriteString(`</td></tr>`)
		return
	}
	for i, it := range items {
		writeCard(b, t, i+1, it)
	}
	b.WriteString(`</td></tr>`)
}

func writeCard(b *strings.Builder, t Theme, index int, it Item) {
	b.WriteString(fmt.Sprintf(`<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="margin:0 0 16px 0;background:%s;border:1px solid %s;border-radius:8px;"><tr><td style="padding:16px 18px;">`,
		t.Surface2, t.Border))

	source := strings.TrimSpace(it.Source)
	if source == "" {
		source = "unknown"
	}
	b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:11px;color:%s;letter-spacing:1px;">[%02d] %s</div>`,
		fontMono, t.Accent2, index, esc(strings.ToLower(source))))

	titleHTML := esc(it.Title)
	if u := safeURL(it.URL); u != "" {
		titleHTML = fmt.Sprintf(`<a href="%s" style="color:%s;">%s</a>`, u, t.Accent, esc(it.Title))
	}
	b.WriteString(fmt.Sprintf(`<div style="font-size:17px;font-weight:600;line-height:1.4;margin:6px 0 8px 0;color:%s;">%s</div>`,
		t.Text, titleHTML))

	if s := strings.TrimSpace(it.Summary); s != "" {
		b.WriteString(fmt.Sprintf(`<div style="font-size:15px;line-height:1.6;color:%s;">%s</div>`, t.Text, esc(s)))
	}

	if w := strings.TrimSpace(it.WhyItMatters); w != "" {
		b.WriteString(fmt.Sprintf(`<div style="margin-top:10px;font-size:14px;line-height:1.55;color:%s;">`, t.Muted))
		b.WriteString(fmt.Sprintf(`<span style="font-family:%s;color:%s;">why_it_matters &gt; </span>`, fontMono, t.Warning))
		b.WriteString(esc(w) + `</div>`)
	}

	writeTags(b, t, it.Tags)

	if u := safeURL(it.URL); u != "" {
		b.WriteString(fmt.Sprintf(`<div style="margin-top:12px;font-family:%s;font-size:12px;"><a href="%s" style="color:%s;">read &rarr; %s</a></div>`,
			fontMono, u, t.Accent, esc(displayURL(it.URL))))
	}

	b.WriteString(`</td></tr></table>`)
}

func writeTags(b *strings.Builder, t Theme, tags []string) {
	var clean []string
	for _, tag := range tags {
		if s := strings.TrimSpace(tag); s != "" {
			clean = append(clean, s)
		}
	}
	if len(clean) == 0 {
		return
	}
	b.WriteString(`<div style="margin-top:12px;line-height:2;">`)
	for _, tag := range clean {
		b.WriteString(fmt.Sprintf(`<span style="font-family:%s;font-size:11px;color:%s;background:%s;border:1px solid %s;border-radius:4px;padding:3px 8px;margin-right:6px;">%s</span>`,
			fontMono, t.Accent, t.Surface, t.Border, esc(strings.ToLower(tag))))
	}
	b.WriteString(`</div>`)
}

func writeFooter(b *strings.Builder, t Theme, n Newsletter) {
	b.WriteString(fmt.Sprintf(`<tr><td class="tl-pad" style="padding:20px 28px 26px 28px;border-top:1px solid %s;background:%s;">`, t.Border, t.Surface2))

	stamp := "unknown"
	if !n.GeneratedAt.IsZero() {
		stamp = n.GeneratedAt.UTC().Format(time.RFC1123)
	}
	b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:12px;color:%s;line-height:1.7;">`, fontMono, t.Muted))
	b.WriteString(fmt.Sprintf(`tux-letter // %d articles // %d sources<br>`, len(n.Items), len(n.Sources)))
	b.WriteString(`generated ` + esc(stamp) + `<br>`)
	b.WriteString(fmt.Sprintf(`<a href="%s" style="color:%s;">%s</a><br>`, RepositoryURL, t.Accent, esc(RepositoryURL)))
	b.WriteString(`<span style="color:` + t.Muted + `">generated locally by tux-letter</span>`)
	b.WriteString(`</div></td></tr>`)
}

func esc(s string) string { return html.EscapeString(s) }

func safeURL(raw string) string {
	raw = strings.TrimSpace(raw)
	low := strings.ToLower(raw)
	if strings.HasPrefix(low, "http://") || strings.HasPrefix(low, "https://") {
		return html.EscapeString(raw)
	}
	return ""
}

func displayURL(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "https://")
	raw = strings.TrimPrefix(raw, "http://")
	if len(raw) > 60 {
		raw = raw[:60] + "…"
	}
	return raw
}
