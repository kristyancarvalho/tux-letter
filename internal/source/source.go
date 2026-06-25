package source

import (
	"strings"
	"unicode"

	"github.com/kristyancarvalho/tux-letter/internal/urlx"
)

type Source struct {
	URL  string
	Name string
	Host string
}

func Build(urls []string) ([]Source, []error) {
	var sources []Source
	var errs []error
	seen := make(map[string]bool)

	for _, raw := range urls {
		canonical, err := urlx.Canonical(raw)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		host, err := urlx.Host(raw)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if seen[canonical] {
			continue
		}
		seen[canonical] = true
		sources = append(sources, Source{
			URL:  canonical,
			Name: NameFromHost(host),
			Host: host,
		})
	}
	return sources, errs
}

func NameFromHost(host string) string {
	host = strings.TrimPrefix(strings.ToLower(host), "www.")
	if host == "" {
		return "Unknown"
	}
	labels := strings.Split(host, ".")
	label := labels[0]
	if label == "" {
		return host
	}
	return capitalize(label)
}

func capitalize(s string) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return s
	}
	if unicode.IsLetter(runes[0]) {
		runes[0] = unicode.ToUpper(runes[0])
	}
	return string(runes)
}
