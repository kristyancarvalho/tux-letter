package article

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/kristyancarvalho/tux-letter/internal/urlx"
)

func Normalize(a Article) Article {
	a.Title = collapse(a.Title)
	a.Author = collapse(a.Author)
	a.Excerpt = collapse(a.Excerpt)
	a.Source = strings.TrimSpace(a.Source)
	if c, err := urlx.Canonical(a.URL); err == nil {
		a.URL = c
	} else {
		a.URL = strings.TrimSpace(a.URL)
	}
	a.Hash = Hash(a)
	return a
}

func Hash(a Article) string {
	canonical := strings.TrimSpace(a.URL)
	if c, err := urlx.Canonical(a.URL); err == nil {
		canonical = c
	}
	seed := strings.ToLower(canonical) + "\n" + strings.ToLower(collapse(a.Title))
	sum := sha256.Sum256([]byte(seed))
	return hex.EncodeToString(sum[:])
}

func CanonicalKey(rawURL string) string {
	if c, err := urlx.Canonical(rawURL); err == nil {
		return strings.ToLower(c)
	}
	return strings.ToLower(strings.TrimSpace(rawURL))
}

func collapse(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
