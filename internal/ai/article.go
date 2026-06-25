package ai

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Article struct {
	Title    string      `json:"title"`
	Subtitle string      `json:"subtitle"`
	Summary  string      `json:"summary"`
	Body     []Section   `json:"body"`
	Sources  []Reference `json:"sources"`
}

type Section struct {
	Heading    string   `json:"heading"`
	Paragraphs []string `json:"paragraphs"`
}

type Reference struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Source string `json:"source"`
	URL    string `json:"url"`
}

type SourceBundle struct {
	Language string         `json:"language"`
	Style    string         `json:"style"`
	Sources  []BundleSource `json:"sources"`
}

type BundleSource struct {
	ID          int    `json:"id"`
	Source      string `json:"source"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	PublishedAt string `json:"published_at,omitempty"`
	Content     string `json:"content"`
}

var citationPattern = regexp.MustCompile(`\[(\d+)\]`)

func (b SourceBundle) ids() map[int]bool {
	out := make(map[int]bool, len(b.Sources))
	for _, s := range b.Sources {
		out[s.ID] = true
	}
	return out
}

func (b SourceBundle) byID() map[int]BundleSource {
	out := make(map[int]BundleSource, len(b.Sources))
	for _, s := range b.Sources {
		out[s.ID] = s
	}
	return out
}

func (a Article) paragraphs() []string {
	var out []string
	for _, sec := range a.Body {
		for _, p := range sec.Paragraphs {
			if strings.TrimSpace(p) != "" {
				out = append(out, p)
			}
		}
	}
	return out
}

func CitedIDs(a Article) []int {
	seen := make(map[int]bool)
	var ids []int
	for _, p := range a.paragraphs() {
		for _, m := range citationPattern.FindAllStringSubmatch(p, -1) {
			n, err := strconv.Atoi(m[1])
			if err != nil {
				continue
			}
			if !seen[n] {
				seen[n] = true
				ids = append(ids, n)
			}
		}
	}
	sort.Ints(ids)
	return ids
}

func (a Article) ValidateAgainst(bundle SourceBundle) error {
	if strings.TrimSpace(a.Title) == "" {
		return fmt.Errorf("article title is empty")
	}
	if len(a.paragraphs()) == 0 {
		return fmt.Errorf("article body has no paragraphs")
	}
	if len(a.Sources) == 0 {
		return fmt.Errorf("article has no sources")
	}
	cited := CitedIDs(a)
	if len(cited) == 0 {
		return fmt.Errorf("article body has no inline citations")
	}
	allowed := bundle.ids()
	for _, id := range cited {
		if !allowed[id] {
			return fmt.Errorf("citation [%d] does not match any provided source", id)
		}
	}
	if !isEnglishCode(bundle.Language) && looksEnglish(a.languageSample()) {
		return fmt.Errorf("output language does not match configured language %q", bundle.Language)
	}
	return nil
}

func (a Article) languageSample() string {
	var b strings.Builder
	b.WriteString(a.Title + " " + a.Subtitle + " " + a.Summary + " ")
	for _, sec := range a.Body {
		b.WriteString(sec.Heading + " ")
		for _, p := range sec.Paragraphs {
			b.WriteString(p + " ")
		}
	}
	return b.String()
}
