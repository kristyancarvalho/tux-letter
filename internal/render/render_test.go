package render

import (
	"strings"
	"testing"
	"time"
)

func sampleNewsletter() Newsletter {
	return Newsletter{
		Brand:       "Tux Letter",
		Title:       "A Quiet Week in Kernel Land",
		Subtitle:    "a few notable releases and one security scare",
		Summary:     "A calm but meaningful cycle for the open-source desktop.",
		GeneratedAt: time.Date(2026, 6, 25, 20, 0, 0, 0, time.UTC),
		Sections: []Section{
			{
				Heading: "kernel",
				Paragraphs: []string{
					"The new kernel ships broader hardware support and better power management [1].",
				},
			},
			{
				Heading: "distro news",
				Paragraphs: []string{
					"Arch refreshed its guided installer, lowering the barrier for new users [2].",
				},
			},
		},
		References: []Reference{
			{ID: 1, Title: "Kernel 6.10 Released", Source: "9to5Linux", URL: "https://9to5linux.com/kernel-6-10"},
			{ID: 2, Title: "Arch Installer Refresh", Source: "Arch Linux", URL: "https://archlinux.org/news/installer"},
		},
	}
}

func TestRenderHTMLUnifiedArticle(t *testing.T) {
	out, err := RenderHTML(sampleNewsletter())
	if err != nil {
		t.Fatalf("RenderHTML: %v", err)
	}
	for _, want := range []string{
		"<!DOCTYPE html>",
		"TUX LETTER",
		"A Quiet Week in Kernel Land",
		"a few notable releases and one security scare",
		"// briefing",
		"// kernel",
		"// distro news",
		"broader hardware support",
		"// sources",
		"Kernel 6.10 Released",
		"Arch Installer Refresh",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("HTML missing %q", want)
		}
	}
}

func TestRenderHTMLIncludesSourceURLs(t *testing.T) {
	out, _ := RenderHTML(sampleNewsletter())
	for _, u := range []string{
		"https://9to5linux.com/kernel-6-10",
		"https://archlinux.org/news/installer",
	} {
		if !strings.Contains(out, u) {
			t.Errorf("HTML missing source URL %q", u)
		}
	}
}

func TestRenderHTMLInlineCitationsLinkToSources(t *testing.T) {
	out, _ := RenderHTML(sampleNewsletter())
	if !strings.Contains(out, `<sup`) {
		t.Error("expected inline citations rendered as superscripts")
	}
	if !strings.Contains(out, `href="https://9to5linux.com/kernel-6-10"`) {
		t.Error("expected inline citation [1] to link to its source URL")
	}
	if !strings.Contains(out, `id="ref-1"`) || !strings.Contains(out, `id="ref-2"`) {
		t.Error("expected footer reference anchors for each cited source")
	}
}

func TestRenderHTMLFooterMatchesInlineCitations(t *testing.T) {
	n := sampleNewsletter()
	out, _ := RenderHTML(n)
	for _, r := range n.References {
		if !strings.Contains(out, r.Title) {
			t.Errorf("footer missing referenced title %q", r.Title)
		}
		if !strings.Contains(out, r.URL) {
			t.Errorf("footer missing referenced url %q", r.URL)
		}
	}
}

func TestRenderHTMLDoesNotUseCardLayout(t *testing.T) {
	out, _ := RenderHTML(sampleNewsletter())
	for _, banned := range []string{"why_it_matters", "why it matters", "read &rarr;"} {
		if strings.Contains(out, banned) {
			t.Errorf("unified email must not render old card element %q", banned)
		}
	}
}

func TestRenderHTMLEscapesUnsafeContent(t *testing.T) {
	n := Newsletter{
		Brand:    "Tux Letter",
		Title:    `<script>alert('t')</script>`,
		Subtitle: `<img src=x onerror=alert(2)>`,
		Sections: []Section{{
			Heading:    `<b>head</b>`,
			Paragraphs: []string{`A paragraph with <script>alert('p')</script> and a cite [1].`},
		}},
		References: []Reference{{ID: 1, Title: `<script>alert('s')</script>`, Source: `<i>evil</i>`, URL: "https://example.com/a"}},
	}
	out, err := RenderHTML(n)
	if err != nil {
		t.Fatal(err)
	}
	for _, raw := range []string{"<script>alert", "<img src=x onerror", "<b>head</b>", "<i>evil</i>"} {
		if strings.Contains(out, raw) {
			t.Errorf("raw unsafe content not escaped: %q", raw)
		}
	}
	if !strings.Contains(out, "&lt;script&gt;") {
		t.Error("expected escaped script entity")
	}
}

func TestRenderHTMLRejectsUnsafeURLScheme(t *testing.T) {
	n := Newsletter{
		Brand:    "Tux Letter",
		Title:    "Bad links",
		Sections: []Section{{Heading: "x", Paragraphs: []string{"A cite [1]."}}},
		References: []Reference{
			{ID: 1, Title: "Bad", Source: "x", URL: "javascript:alert(1)"},
		},
	}
	out, err := RenderHTML(n)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "javascript:alert") {
		t.Error("javascript: URL must not be emitted as a link")
	}
}

func TestRenderHTMLReadableCypherpunkIdentity(t *testing.T) {
	out, _ := RenderHTML(sampleNewsletter())
	for _, want := range []string{
		"#0b0f14",
		"#00e5ff",
		"font-size:15px;line-height:1.75",
		"@media (max-width:620px)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected readable cypherpunk marker %q", want)
		}
	}
}

func TestRenderHTMLNoSections(t *testing.T) {
	out, err := RenderHTML(Newsletter{Brand: "Tux Letter", Title: "Empty"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "No dispatch content") {
		t.Error("empty newsletter should show a placeholder")
	}
	if strings.Contains(out, "// sources") {
		t.Error("no references means no sources footer")
	}
}

func TestRenderTextUnifiedArticle(t *testing.T) {
	out, err := RenderText(sampleNewsletter())
	if err != nil {
		t.Fatalf("RenderText: %v", err)
	}
	for _, want := range []string{
		"TUX LETTER",
		"A Quiet Week in Kernel Land",
		"// KERNEL",
		"broader hardware support",
		"management [1].",
		"SOURCES",
		"[1] Kernel 6.10 Released — 9to5Linux",
		"https://9to5linux.com/kernel-6-10",
		"https://archlinux.org/news/installer",
		"2 sources",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("text missing %q", want)
		}
	}
	if strings.Contains(out, "<") {
		t.Error("plain text should not contain HTML tags")
	}
	if strings.Contains(out, "why_it_matters") {
		t.Error("plain text must not render the old why_it_matters layout")
	}
}
