package render

import "time"

type Newsletter struct {
	Brand       string
	Title       string
	Subtitle    string
	Summary     string
	Sections    []Section
	References  []Reference
	GeneratedAt time.Time
}

type Section struct {
	Heading    string
	Paragraphs []string
}

type Reference struct {
	ID     int
	Title  string
	Source string
	URL    string
}

const (
	RepositoryURL = "https://github.com/kristyancarvalho/tux-letter"
	Tagline       = "encrypted open-source dispatch"
	SubTagline    = "linux // foss // security // systems"
)
