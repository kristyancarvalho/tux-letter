package ai

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kristyancarvalho/tux-letter/internal/config"
)

func loadFixture(t *testing.T, name string) Article {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	a, err := ParseArticle(string(data))
	if err != nil {
		t.Fatalf("parse fixture %s: %v", name, err)
	}
	return a
}

func testBundle() SourceBundle {
	return BuildBundle(config.NewsletterConfig{}, sampleArticles())
}

func TestValidateAcceptsCitedArticle(t *testing.T) {
	a := loadFixture(t, "valid.json")
	if err := a.ValidateAgainst(testBundle()); err != nil {
		t.Fatalf("expected valid article, got error: %v", err)
	}
}

func TestValidateRejectsInvalidOutput(t *testing.T) {
	cases := []string{
		"no_citations.json",
		"bad_citation.json",
		"missing_title.json",
		"empty_body.json",
		"no_sources.json",
	}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			a := loadFixture(t, name)
			if err := a.ValidateAgainst(testBundle()); err == nil {
				t.Errorf("expected validation error for %s, got nil", name)
			}
		})
	}
}

func TestReconcileFooterMatchesCitations(t *testing.T) {
	bundle := testBundle()
	a := Article{
		Title: "T",
		Body: []Section{
			{Heading: "k", Paragraphs: []string{"Only the second source is cited here [2]."}},
		},
		Sources: []Reference{
			{ID: 1, Title: "wrong", Source: "wrong", URL: "https://wrong.example/x"},
			{ID: 2, Title: "wrong too", Source: "wrong", URL: "https://wrong.example/y"},
		},
	}
	out := reconcile(a, bundle)
	if len(out.Sources) != 1 {
		t.Fatalf("expected footer to contain only cited sources, got %d", len(out.Sources))
	}
	if out.Sources[0].ID != 2 {
		t.Errorf("expected footer id 2, got %d", out.Sources[0].ID)
	}
	if out.Sources[0].URL != "https://phoronix.com/b" {
		t.Errorf("expected reconciled url from bundle, got %q", out.Sources[0].URL)
	}
}

type seqCompleter struct {
	responses []string
	calls     int
}

func (s *seqCompleter) Complete(ctx context.Context, model string, prompt Prompt) (string, error) {
	i := s.calls
	s.calls++
	if i < len(s.responses) {
		return s.responses[i], nil
	}
	return "", fmt.Errorf("no more responses")
}

func TestGeneratorRepairsInvalidThenSucceeds(t *testing.T) {
	seq := &seqCompleter{responses: []string{"not json at all", validJSON}}
	g := NewGenerator(seq, []string{"only"})
	a, err := g.Generate(context.Background(), config.NewsletterConfig{}, sampleArticles())
	if err != nil {
		t.Fatalf("expected repair to recover, got error: %v", err)
	}
	if seq.calls != 2 {
		t.Errorf("expected exactly one repair attempt (2 calls), got %d", seq.calls)
	}
	if len(a.Sources) == 0 {
		t.Error("expected reconciled sources after repair")
	}
}

func TestRepairPromptIncludesReason(t *testing.T) {
	base := Prompt{System: "base", User: "u"}
	repaired := repairPrompt(base, fmt.Errorf("article body has no inline citations"))
	if repaired.User != base.User {
		t.Error("repair prompt must preserve the user message")
	}
	if repaired.System == base.System {
		t.Error("repair prompt must strengthen the system message")
	}
	if !strings.Contains(repaired.System, "no inline citations") {
		t.Error("repair prompt should include the rejection reason")
	}
}
