package render

import "time"

type Newsletter struct {
	Title       string
	Summary     string
	Items       []Item
	GeneratedAt time.Time
	Sources     []string
}

type Item struct {
	Title        string
	Source       string
	URL          string
	Summary      string
	WhyItMatters string
	Tags         []string
}

const (
	RepositoryURL = "https://github.com/kristyancarvalho/tux-letter"
	Tagline       = "encrypted open-source dispatch"
	SubTagline    = "linux // foss // security // systems"
)
