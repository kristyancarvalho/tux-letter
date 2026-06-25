package render

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var citationRe = regexp.MustCompile(`\[(\d+)\]`)

func RenderHTML(n Newsletter) (string, error) {
	return RenderHTMLWithTheme(n, Cypherpunk())
}

func RenderHTMLWithTheme(n Newsletter, t Theme) (string, error) {
	brand := strings.TrimSpace(n.Brand)
	if brand == "" {
		brand = "Tux Letter"
	}
	headline := strings.TrimSpace(n.Title)
	if headline == "" {
		headline = brand
	}
	refs := referenceMap(n.References)
	loc := localeFor(n.Language)

	var b strings.Builder
	b.WriteString("<!DOCTYPE html>\n")
	b.WriteString(`<html lang="` + esc(htmlLang(n.Language)) + `"><head><meta charset="utf-8">`)
	b.WriteString(`<meta name="viewport" content="width=device-width,initial-scale=1">`)
	b.WriteString(`<meta name="color-scheme" content="dark">`)
	b.WriteString(`<title>` + esc(headline) + `</title>`)
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

	writeHeader(&b, t, loc, brand)
	writeMeta(&b, t, loc, n)
	writeHeadline(&b, t, headline, n.Subtitle)
	if s := strings.TrimSpace(n.Summary); s != "" {
		writeSummary(&b, t, loc, s)
	}
	writeBody(&b, t, loc, n.Sections, refs)
	writeReferences(&b, t, loc, n.References)
	writeFooter(&b, t, loc, n)

	b.WriteString(`</table></td></tr></table></body></html>`)
	return b.String(), nil
}

func writeHeader(b *strings.Builder, t Theme, loc locale, brand string) {
	b.WriteString(fmt.Sprintf(`<tr><td class="tl-pad" style="padding:22px 28px;background:%s;border-bottom:1px solid %s;">`, t.Surface2, t.Border))

	b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:12px;color:%s;letter-spacing:1px;">`, fontMono, t.Muted))
	b.WriteString(`<span style="color:#ff5f56">&#9679;</span> <span style="color:` + t.Warning + `">&#9679;</span> <span style="color:` + t.Accent2 + `">&#9679;</span>`)
	b.WriteString(`&nbsp;&nbsp;tux@letter:~$ cat dispatch</div>`)

	b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:30px;font-weight:700;letter-spacing:4px;color:%s;margin-top:12px;">%s</div>`,
		fontMono, t.Accent, esc(strings.ToUpper(brand))))

	b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:13px;color:%s;margin-top:4px;">%s</div>`,
		fontMono, t.Accent2, esc(loc.tagline)))
	b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:12px;color:%s;margin-top:2px;">%s</div>`,
		fontMono, t.Muted, esc(loc.subTagline)))
	b.WriteString(`</td></tr>`)
}

func writeMeta(b *strings.Builder, t Theme, loc locale, n Newsletter) {
	stamp := "unknown"
	if !n.GeneratedAt.IsZero() {
		stamp = n.GeneratedAt.UTC().Format("2006-01-02 15:04 UTC")
	}
	b.WriteString(fmt.Sprintf(`<tr><td class="tl-pad" style="padding:16px 28px 0 28px;"><div style="font-family:%s;font-size:11px;color:%s;letter-spacing:1px;text-transform:uppercase;">issue // %s // %d %s</div></td></tr>`,
		fontMono, t.Muted, esc(stamp), len(n.References), esc(loc.sourcesWord)))
}

func writeHeadline(b *strings.Builder, t Theme, headline, subtitle string) {
	b.WriteString(`<tr><td class="tl-pad" style="padding:10px 28px 0 28px;">`)
	b.WriteString(fmt.Sprintf(`<div style="font-size:26px;font-weight:700;line-height:1.28;color:%s;">%s</div>`, t.Text, esc(headline)))
	if s := strings.TrimSpace(subtitle); s != "" {
		b.WriteString(fmt.Sprintf(`<div style="font-size:16px;line-height:1.5;color:%s;margin-top:8px;">%s</div>`, t.Muted, esc(s)))
	}
	b.WriteString(`</td></tr>`)
}

func writeSummary(b *strings.Builder, t Theme, loc locale, summary string) {
	b.WriteString(`<tr><td class="tl-pad" style="padding:20px 28px 6px 28px;">`)
	b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:11px;text-transform:uppercase;letter-spacing:2px;color:%s;margin-bottom:8px;">// %s</div>`, fontMono, t.Accent, esc(loc.briefing)))
	b.WriteString(fmt.Sprintf(`<div style="background:%s;border-left:3px solid %s;border-radius:6px;padding:16px 18px;font-size:15px;line-height:1.6;color:%s;">%s</div>`,
		t.Surface2, t.Accent, t.Text, esc(summary)))
	b.WriteString(`</td></tr>`)
}

func writeBody(b *strings.Builder, t Theme, loc locale, sections []Section, refs map[int]Reference) {
	if len(sections) == 0 {
		b.WriteString(`<tr><td class="tl-pad" style="padding:18px 28px 8px 28px;">`)
		b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:14px;color:%s;padding:8px 0;">%s</div>`, fontMono, t.Muted, esc(loc.emptyBody)))
		b.WriteString(`</td></tr>`)
		return
	}

	b.WriteString(`<tr><td class="tl-pad" style="padding:8px 28px 8px 28px;">`)
	for _, sec := range sections {
		if h := strings.TrimSpace(sec.Heading); h != "" {
			b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:13px;text-transform:uppercase;letter-spacing:2px;color:%s;margin:22px 0 10px 0;">// %s</div>`,
				fontMono, t.Accent2, esc(strings.ToLower(h))))
		}
		for _, p := range sec.Paragraphs {
			if strings.TrimSpace(p) == "" {
				continue
			}
			b.WriteString(fmt.Sprintf(`<p style="margin:0 0 14px 0;font-size:15px;line-height:1.75;color:%s;">%s</p>`,
				t.Text, renderParagraph(t, p, refs)))
		}
	}
	b.WriteString(`</td></tr>`)
}

func writeReferences(b *strings.Builder, t Theme, loc locale, refs []Reference) {
	if len(refs) == 0 {
		return
	}
	b.WriteString(fmt.Sprintf(`<tr><td class="tl-pad" style="padding:20px 28px 6px 28px;border-top:1px solid %s;">`, t.Border))
	b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:11px;text-transform:uppercase;letter-spacing:2px;color:%s;margin-bottom:12px;">// %s</div>`, fontMono, t.Accent, esc(loc.sources)))
	for _, r := range refs {
		b.WriteString(fmt.Sprintf(`<div id="ref-%d" style="margin:0 0 12px 0;font-size:13px;line-height:1.5;color:%s;">`, r.ID, t.Muted))
		b.WriteString(fmt.Sprintf(`<span style="font-family:%s;color:%s;">[%d]</span> `, fontMono, t.Accent, r.ID))

		label := esc(strings.TrimSpace(r.Title))
		if src := strings.TrimSpace(r.Source); src != "" {
			label += ` &mdash; <span style="color:` + t.Accent2 + `;">` + esc(src) + `</span>`
		}
		if u := safeURL(r.URL); u != "" {
			b.WriteString(fmt.Sprintf(`<a href="%s" style="color:%s;">%s</a>`, u, t.Text, label))
			b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:12px;color:%s;margin-top:2px;">%s</div>`, fontMono, t.Muted, esc(r.URL)))
		} else {
			b.WriteString(label)
		}
		b.WriteString(`</div>`)
	}
	b.WriteString(`</td></tr>`)
}

func writeFooter(b *strings.Builder, t Theme, loc locale, n Newsletter) {
	b.WriteString(fmt.Sprintf(`<tr><td class="tl-pad" style="padding:18px 28px 26px 28px;border-top:1px solid %s;background:%s;">`, t.Border, t.Surface2))

	stamp := "unknown"
	if !n.GeneratedAt.IsZero() {
		stamp = n.GeneratedAt.UTC().Format(time.RFC1123)
	}
	b.WriteString(fmt.Sprintf(`<div style="font-family:%s;font-size:12px;color:%s;line-height:1.7;">`, fontMono, t.Muted))
	b.WriteString(fmt.Sprintf(`tux-letter // %d %s<br>`, len(n.References), esc(loc.sourcesWord)))
	b.WriteString(esc(loc.generated) + ` ` + esc(stamp) + `<br>`)
	b.WriteString(fmt.Sprintf(`<a href="%s" style="color:%s;">%s</a><br>`, RepositoryURL, t.Accent, esc(RepositoryURL)))
	b.WriteString(`<span style="color:` + t.Muted + `">` + esc(loc.generatedBy) + `</span>`)
	b.WriteString(`</div></td></tr>`)
}

func renderParagraph(t Theme, paragraph string, refs map[int]Reference) string {
	escaped := esc(paragraph)
	return citationRe.ReplaceAllStringFunc(escaped, func(m string) string {
		sub := citationRe.FindStringSubmatch(m)
		id, err := strconv.Atoi(sub[1])
		if err != nil {
			return m
		}
		if ref, ok := refs[id]; ok {
			if u := safeURL(ref.URL); u != "" {
				return fmt.Sprintf(`<sup style="font-size:11px;"><a href="%s" style="color:%s;font-family:%s;">[%d]</a></sup>`, u, t.Accent, fontMono, id)
			}
		}
		return fmt.Sprintf(`<sup style="font-size:11px;color:%s;font-family:%s;">[%d]</sup>`, t.Accent, fontMono, id)
	})
}

func referenceMap(refs []Reference) map[int]Reference {
	m := make(map[int]Reference, len(refs))
	for _, r := range refs {
		m[r.ID] = r
	}
	return m
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
