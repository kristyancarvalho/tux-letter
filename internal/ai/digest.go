package ai

import (
	"fmt"
	"strings"
)

type Digest struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Items   []Item `json:"items"`
}

type Item struct {
	Title        string   `json:"title"`
	Source       string   `json:"source"`
	URL          string   `json:"url"`
	Summary      string   `json:"summary"`
	WhyItMatters string   `json:"why_it_matters"`
	Tags         []string `json:"tags"`
}

func (d Digest) Validate() error {
	if strings.TrimSpace(d.Title) == "" {
		return fmt.Errorf("digest title is empty")
	}
	if len(d.Items) == 0 {
		return fmt.Errorf("digest has no items")
	}
	for i, it := range d.Items {
		if strings.TrimSpace(it.Title) == "" {
			return fmt.Errorf("item %d has empty title", i)
		}
		if strings.TrimSpace(it.URL) == "" {
			return fmt.Errorf("item %d (%q) has empty url", i, it.Title)
		}
	}
	return nil
}
