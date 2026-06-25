package render

import (
	"strings"
	"testing"
	"time"
)

func sampleNewsletter() Newsletter {
	return Newsletter{
		Title:       "Tux Letter",
		Summary:     "A quiet week in kernel land with a few notable releases.",
		GeneratedAt: time.Date(2026, 6, 25, 20, 0, 0, 0, time.UTC),
		Sources:     []string{"9to5linux.com", "archlinux.org"},
		Items: []Item{
			{
				Title:        "Kernel 6.10 Released",
				Source:       "9to5Linux",
				URL:          "https://9to5linux.com/kernel-6-10",
				Summary:      "The new kernel ships better hardware support.",
				WhyItMatters: "Improves laptop battery life.",
				Tags:         []string{"linux", "kernel"},
			},
			{
				Title:  "Arch News Update",
				Source: "Archlinux",
				URL:    "https://archlinux.org/news/x",
			},
		},
	}
}

func TestRenderHTMLContainsCoreContent(t *testing.T) {
	out, err := RenderHTML(sampleNewsletter())
	if err != nil {
		t.Fatalf("RenderHTML: %v", err)
	}
	for _, want := range []string{
		"<!DOCTYPE html>",
		"TUX LETTER",
		"Kernel 6.10 Released",
		"https://9to5linux.com/kernel-6-10",
		"why_it_matters",
		"#0b0f14",
		"kernel",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("HTML missing %q", want)
		}
	}
}

func TestRenderHTMLEscapesUnsafeContent(t *testing.T) {
	n := Newsletter{
		Title: "Tux Letter",
		Items: []Item{{
			Title:   `<script>alert('x')</script>`,
			Summary: `<img src=x onerror=alert(1)>`,
			Source:  "evil",
			URL:     "https://example.com/a",
		}},
	}
	out, err := RenderHTML(n)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "<script>alert") {
		t.Error("raw <script> must be escaped")
	}
	if strings.Contains(out, "<img src=x onerror") {
		t.Error("raw <img> must be escaped")
	}
	if !strings.Contains(out, "&lt;script&gt;") {
		t.Error("expected escaped script entity")
	}
}

func TestRenderHTMLRejectsUnsafeURLScheme(t *testing.T) {
	n := Newsletter{
		Title: "Tux Letter",
		Items: []Item{{Title: "Bad", URL: "javascript:alert(1)", Source: "x"}},
	}
	out, err := RenderHTML(n)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "javascript:alert") {
		t.Error("javascript: URL must not be emitted as a link")
	}
}

func TestRenderHTMLEmptyOptionalFields(t *testing.T) {
	n := Newsletter{
		Title: "Tux Letter",
		Items: []Item{{Title: "Just a title", Source: "", URL: ""}},
	}
	out, err := RenderHTML(n)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Just a title") {
		t.Error("title should render even with empty optional fields")
	}
	if !strings.Contains(out, "unknown") {
		t.Error("missing source should fall back to 'unknown'")
	}
}

func TestRenderHTMLNoItems(t *testing.T) {
	out, err := RenderHTML(Newsletter{Title: "Tux Letter"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "No new dispatches") {
		t.Error("empty newsletter should show a placeholder")
	}
}

func TestRenderText(t *testing.T) {
	out, err := RenderText(sampleNewsletter())
	if err != nil {
		t.Fatalf("RenderText: %v", err)
	}
	for _, want := range []string{
		"TUX LETTER",
		"Kernel 6.10 Released",
		"https://9to5linux.com/kernel-6-10",
		"#linux #kernel",
		"why it matters",
		"2 articles // 2 sources",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("text missing %q", want)
		}
	}
	if strings.Contains(out, "<") {
		t.Error("plain text should not contain HTML tags")
	}
}
