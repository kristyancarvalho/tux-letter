package article

import "time"

type Article struct {
	Title     string
	URL       string
	Source    string
	Author    string
	Excerpt   string
	Published time.Time
	Hash      string
}

func (a Article) IsZero() bool {
	return a.Title == "" && a.URL == ""
}
