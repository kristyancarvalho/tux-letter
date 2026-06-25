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
			{ID: 2, Title: "Arch Installer Refresh", Source: "Arch Linux", URL: "https://archlinux.org/news/x"},
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
		"A Quiet Week in Kernel Land",
		"https://9to5linux.com/kernel-6-10",
		"// sources",
		"[1]",
		"#0b0f14",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("HTML missing %q", want)
		}
	}
	if strings.Contains(out, "why_it_matters") {
		t.Error("unified newsletter must not render the old why_it_matters layout")
	}
}

func TestRenderTextContainsCoreContent(t *testing.T) {
	out, err := RenderText(sampleNewsletter())
	if err != nil {
		t.Fatalf("RenderText: %v", err)
	}
	for _, want := range []string{
		"TUX LETTER",
		"A Quiet Week in Kernel Land",
		"SOURCES",
		"[1] Kernel 6.10 Released",
		"https://9to5linux.com/kernel-6-10",
		"2 sources",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("text missing %q", want)
		}
	}
	if strings.Contains(out, "<") {
		t.Error("plain text should not contain HTML tags")
	}
}
