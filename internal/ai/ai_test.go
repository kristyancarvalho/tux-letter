package ai

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/kristyancarvalho/tux-letter/internal/article"
	"github.com/kristyancarvalho/tux-letter/internal/config"
)

type stubCompleter struct {
	responses map[string]string
	errs      map[string]error
	calls     []string
}

func (s *stubCompleter) Complete(ctx context.Context, model string, prompt Prompt) (string, error) {
	s.calls = append(s.calls, model)
	if err, ok := s.errs[model]; ok && err != nil {
		return "", err
	}
	return s.responses[model], nil
}

func sampleArticles() []article.Article {
	return []article.Article{
		{Title: "Kernel 6.20", Source: "lwn", URL: "https://lwn.net/a"},
		{Title: "Mesa 25", Source: "phoronix", URL: "https://phoronix.com/b"},
	}
}

const validJSON = `{"title":"Tux Letter Weekly","subtitle":"a quiet but meaningful cycle","summary":"Two updates this week.","body":[{"heading":"kernel and graphics","paragraphs":["The 6.20 kernel broadens hardware support [1], while Mesa 25 modernizes the graphics stack [2]."]}],"sources":[{"id":1,"title":"Kernel 6.20","source":"lwn","url":"https://lwn.net/a"},{"id":2,"title":"Mesa 25","source":"phoronix","url":"https://phoronix.com/b"}]}`

func TestGenerateUsesFirstWorkingModel(t *testing.T) {
	stub := &stubCompleter{
		responses: map[string]string{"good": validJSON},
	}
	g := NewGenerator(stub, []string{"good", "unused"})
	a, err := g.Generate(context.Background(), config.NewsletterConfig{Title: "Tux Letter"}, sampleArticles())
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if a.Title == "" {
		t.Error("expected a title")
	}
	if len(a.Sources) != 2 {
		t.Fatalf("expected 2 reconciled sources, got %d", len(a.Sources))
	}
	if len(stub.calls) != 1 || stub.calls[0] != "good" {
		t.Errorf("expected only first model called, got %v", stub.calls)
	}
}

func TestGenerateFallsBackOnRateLimit(t *testing.T) {
	stub := &stubCompleter{
		responses: map[string]string{"second": validJSON},
		errs:      map[string]error{"first": &APIError{Status: http.StatusTooManyRequests, Message: "slow down"}},
	}
	g := NewGenerator(stub, []string{"first", "second"})
	a, err := g.Generate(context.Background(), config.NewsletterConfig{}, sampleArticles())
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if a.Title == "" {
		t.Error("expected a title")
	}
	if len(stub.calls) != 2 || stub.calls[0] != "first" || stub.calls[1] != "second" {
		t.Errorf("expected fallback first->second, got %v", stub.calls)
	}
}

func TestGenerateRetriesInvalidOutputThenNextModel(t *testing.T) {
	stub := &stubCompleter{
		responses: map[string]string{
			"bad":  "not json at all",
			"good": validJSON,
		},
	}
	g := NewGenerator(stub, []string{"bad", "good"})
	if _, err := g.Generate(context.Background(), config.NewsletterConfig{}, sampleArticles()); err != nil {
		t.Fatalf("generate: %v", err)
	}
	got := strings.Join(stub.calls, ",")
	if got != "bad,bad,good" {
		t.Errorf("expected bad retried once then good, got %q", got)
	}
}

func TestGenerateAllModelsFail(t *testing.T) {
	stub := &stubCompleter{responses: map[string]string{"a": "garbage", "b": "garbage"}}
	g := NewGenerator(stub, []string{"a", "b"})
	if _, err := g.Generate(context.Background(), config.NewsletterConfig{}, sampleArticles()); err == nil {
		t.Fatal("expected error when all models fail")
	}
}

func TestFallbackArticleIsValid(t *testing.T) {
	cfg := config.NewsletterConfig{Title: "Tux Letter"}
	a := FallbackArticle(cfg, sampleArticles())
	if err := a.ValidateAgainst(BuildBundle(cfg, sampleArticles())); err != nil {
		t.Fatalf("fallback article invalid: %v", err)
	}
	if len(a.Body) == 0 {
		t.Error("expected fallback body sections")
	}
	if len(a.Sources) != 2 {
		t.Errorf("expected 2 sources, got %d", len(a.Sources))
	}
	if len(CitedIDs(a)) == 0 {
		t.Error("expected inline citations in fallback")
	}
}

func TestAPIErrorTransient(t *testing.T) {
	cases := map[int]bool{429: true, 500: true, 503: true, 400: false, 401: false}
	for status, want := range cases {
		got := (&APIError{Status: status}).Transient()
		if got != want {
			t.Errorf("status %d: transient=%v want %v", status, got, want)
		}
	}
}

func TestParseArticleStripsCodeFence(t *testing.T) {
	raw := "```json\n" + validJSON + "\n```"
	a, err := ParseArticle(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if a.Title != "Tux Letter Weekly" {
		t.Errorf("unexpected title: %q", a.Title)
	}
}
