package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kristyancarvalho/tux-letter/internal/article"
)

func TestStorePersistsSeenAndMetadata(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "state.json")

	s, err := Open(path, 1000)
	if err != nil {
		t.Fatalf("open: %v", err)
	}

	now := time.Date(2026, 6, 25, 20, 0, 0, 0, time.UTC)
	a := article.Normalize(article.Article{Title: "Kernel 6.20", URL: "https://x.example/kernel-6-20"})
	s.Remember(a)
	s.MarkAttempt(now)
	s.MarkSuccess(now)
	s.SetSourceStatus("x", "ok", now)
	s.IncError("fetch")
	if err := s.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("state file not written: %v", err)
	}

	reopened, err := Open(path, 1000)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	if !reopened.Seen(a) {
		t.Error("expected remembered article to be seen after reload")
	}
	snap := reopened.Snapshot()
	if !snap.LastSuccessfulRun.Equal(now) {
		t.Errorf("last successful run not persisted: %v", snap.LastSuccessfulRun)
	}
	if snap.SourceStatus["x"].Status != "ok" {
		t.Errorf("source status not persisted: %+v", snap.SourceStatus)
	}
	if snap.ErrorCounts["fetch"] != 1 {
		t.Errorf("error count not persisted: %d", snap.ErrorCounts["fetch"])
	}
}

func TestFilterNewSkipsSeen(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "state.json"), 1000)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	batch := article.NormalizeAll([]article.Article{
		{Title: "One", URL: "https://x.example/one"},
		{Title: "Two", URL: "https://x.example/two"},
	})

	first := s.FilterNew(batch)
	if len(first) != 2 {
		t.Fatalf("expected 2 new, got %d", len(first))
	}
	second := s.FilterNew(batch)
	if len(second) != 0 {
		t.Fatalf("expected 0 new on second pass, got %d", len(second))
	}
}

func TestPruneCapsSeenEntries(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "state.json"), 3)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	urls := []string{
		"https://x.example/a", "https://x.example/b",
		"https://x.example/c", "https://x.example/d",
		"https://x.example/e",
	}
	for _, u := range urls {
		s.Remember(article.Normalize(article.Article{Title: u, URL: u}))
	}
	snap := s.Snapshot()
	if len(snap.SeenURLs) != 3 {
		t.Fatalf("expected pruned to 3, got %d", len(snap.SeenURLs))
	}
	oldest := article.Normalize(article.Article{Title: urls[0], URL: urls[0]})
	if s.Seen(oldest) {
		t.Error("expected oldest entry to be pruned out")
	}
	newest := article.Normalize(article.Article{Title: urls[4], URL: urls[4]})
	if !s.Seen(newest) {
		t.Error("expected newest entry to remain")
	}
}

func TestOpenMissingFileStartsEmpty(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "does-not-exist.json"), 1000)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if len(s.Snapshot().SeenURLs) != 0 {
		t.Error("expected empty state for missing file")
	}
}
