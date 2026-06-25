package feed

import (
	"encoding/xml"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/kristyancarvalho/tux-letter/internal/article"
)

type xmlFeed struct {
	XMLName xml.Name
	Title   string `xml:"title"`

	Channel struct {
		Title string    `xml:"title"`
		Items []rssItem `xml:"item"`
	} `xml:"channel"`

	RDFItems []rssItem `xml:"item"`

	Entries []atomEntry `xml:"entry"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Creator     string `xml:"creator"`
	Author      string `xml:"author"`
	PubDate     string `xml:"pubDate"`
	Date        string `xml:"date"`
	GUID        string `xml:"guid"`
}

type atomEntry struct {
	Title   string     `xml:"title"`
	Links   []atomLink `xml:"link"`
	Summary string     `xml:"summary"`
	Content string     `xml:"content"`
	Updated string     `xml:"updated"`
	Publish string     `xml:"published"`
	Author  struct {
		Name string `xml:"name"`
	} `xml:"author"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

var dateLayouts = []string{
	time.RFC1123Z,
	time.RFC1123,
	time.RFC822Z,
	time.RFC822,
	time.RFC3339,
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02 15:04:05",
	"2006-01-02",
	"Mon, 2 Jan 2006 15:04:05 -0700",
}

func Parse(sourceName string, data []byte) ([]article.Article, error) {
	var f xmlFeed
	if err := xml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parse feed: %w", err)
	}

	feedTitle := strings.TrimSpace(f.Channel.Title)
	if feedTitle == "" {
		feedTitle = strings.TrimSpace(f.Title)
	}
	if sourceName == "" {
		sourceName = feedTitle
	}

	var articles []article.Article

	items := f.Channel.Items
	if len(items) == 0 {
		items = f.RDFItems
	}
	for _, it := range items {
		title := clean(it.Title)
		if title == "" || strings.TrimSpace(it.Link) == "" {
			continue
		}
		author := firstNonEmpty(it.Creator, it.Author)
		articles = append(articles, article.Article{
			Title:     title,
			URL:       strings.TrimSpace(it.Link),
			Source:    sourceName,
			Author:    clean(author),
			Excerpt:   excerpt(it.Description),
			Published: parseDate(firstNonEmpty(it.PubDate, it.Date)),
		})
	}

	for _, e := range f.Entries {
		title := clean(e.Title)
		link := atomBestLink(e.Links)
		if title == "" || link == "" {
			continue
		}
		articles = append(articles, article.Article{
			Title:     title,
			URL:       link,
			Source:    sourceName,
			Author:    clean(e.Author.Name),
			Excerpt:   excerpt(firstNonEmpty(e.Summary, e.Content)),
			Published: parseDate(firstNonEmpty(e.Publish, e.Updated)),
		})
	}

	return articles, nil
}

func atomBestLink(links []atomLink) string {
	var fallback string
	for _, l := range links {
		if strings.TrimSpace(l.Href) == "" {
			continue
		}
		if l.Rel == "alternate" || l.Rel == "" {
			return strings.TrimSpace(l.Href)
		}
		if fallback == "" {
			fallback = strings.TrimSpace(l.Href)
		}
	}
	return fallback
}

func parseDate(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	for _, layout := range dateLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}

func clean(s string) string {
	return strings.TrimSpace(html.UnescapeString(s))
}

func excerpt(s string) string {
	s = clean(stripTags(s))
	const max = 400
	if len(s) > max {
		return strings.TrimSpace(s[:max]) + "…"
	}
	return s
}

func stripTags(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
