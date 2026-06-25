package state

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"time"

	"github.com/kristyancarvalho/tux-letter/internal/article"
)

type SourceStatus struct {
	Status string    `json:"status"`
	At     time.Time `json:"at"`
}

type Data struct {
	SeenURLs          []string                `json:"seen_urls"`
	SeenHashes        []string                `json:"seen_hashes"`
	LastSuccessfulRun time.Time               `json:"last_successful_run"`
	LastAttemptedRun  time.Time               `json:"last_attempted_run"`
	SourceStatus      map[string]SourceStatus `json:"source_status"`
	ErrorCounts       map[string]int          `json:"error_counts"`
}

type Store struct {
	path    string
	maxSeen int
	data    Data
	urls    map[string]bool
	hashes  map[string]bool
}

func Open(path string, maxSeen int) (*Store, error) {
	s := &Store{
		path:    path,
		maxSeen: maxSeen,
		urls:    make(map[string]bool),
		hashes:  make(map[string]bool),
		data: Data{
			SourceStatus: make(map[string]SourceStatus),
			ErrorCounts:  make(map[string]int),
		},
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return s, nil
		}
		return nil, err
	}
	if len(raw) == 0 {
		return s, nil
	}
	if err := json.Unmarshal(raw, &s.data); err != nil {
		return nil, err
	}
	if s.data.SourceStatus == nil {
		s.data.SourceStatus = make(map[string]SourceStatus)
	}
	if s.data.ErrorCounts == nil {
		s.data.ErrorCounts = make(map[string]int)
	}
	for _, u := range s.data.SeenURLs {
		s.urls[u] = true
	}
	for _, h := range s.data.SeenHashes {
		s.hashes[h] = true
	}
	return s, nil
}

func (s *Store) Seen(a article.Article) bool {
	if key := article.CanonicalKey(a.URL); key != "" && s.urls[key] {
		return true
	}
	hash := a.Hash
	if hash == "" {
		hash = article.Hash(a)
	}
	return hash != "" && s.hashes[hash]
}

func (s *Store) Remember(a article.Article) {
	key := article.CanonicalKey(a.URL)
	hash := a.Hash
	if hash == "" {
		hash = article.Hash(a)
	}
	if key != "" && !s.urls[key] {
		s.urls[key] = true
		s.data.SeenURLs = append(s.data.SeenURLs, key)
	}
	if hash != "" && !s.hashes[hash] {
		s.hashes[hash] = true
		s.data.SeenHashes = append(s.data.SeenHashes, hash)
	}
	s.prune()
}

func (s *Store) FilterNew(items []article.Article) []article.Article {
	var out []article.Article
	for _, a := range items {
		if s.Seen(a) {
			continue
		}
		out = append(out, a)
		s.Remember(a)
	}
	return out
}

func (s *Store) MarkAttempt(t time.Time) {
	s.data.LastAttemptedRun = t.UTC()
}

func (s *Store) MarkSuccess(t time.Time) {
	s.data.LastSuccessfulRun = t.UTC()
}

func (s *Store) SetSourceStatus(name, status string, t time.Time) {
	s.data.SourceStatus[name] = SourceStatus{Status: status, At: t.UTC()}
}

func (s *Store) IncError(key string) {
	s.data.ErrorCounts[key]++
}

func (s *Store) Snapshot() Data {
	return s.data
}

func (s *Store) prune() {
	if s.maxSeen <= 0 {
		return
	}
	if n := len(s.data.SeenURLs); n > s.maxSeen {
		drop := s.data.SeenURLs[:n-s.maxSeen]
		for _, u := range drop {
			delete(s.urls, u)
		}
		s.data.SeenURLs = append([]string(nil), s.data.SeenURLs[n-s.maxSeen:]...)
	}
	if n := len(s.data.SeenHashes); n > s.maxSeen {
		drop := s.data.SeenHashes[:n-s.maxSeen]
		for _, h := range drop {
			delete(s.hashes, h)
		}
		s.data.SeenHashes = append([]string(nil), s.data.SeenHashes[n-s.maxSeen:]...)
	}
}

func (s *Store) Save() error {
	raw, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return atomicWrite(s.path, raw)
}
