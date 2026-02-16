package scraper

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kristyancarvalho/tux-letter/internal/database"

	"github.com/gocolly/colly/v2"
)

type Site struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Selector string `json:"selector"`
}

type NewsItem struct {
	Title string
	Link  string
	Site  string
}

func LoadSites(path string) ([]Site, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var sites []Site
	err = json.Unmarshal(file, &sites)
	return sites, err
}

func ScrapeAll(sitesPath string) ([]NewsItem, error) {
	sites, err := LoadSites(sitesPath)
	if err != nil {
		return nil, err
	}

	var allNews []NewsItem

	for _, site := range sites {
		news, err := scrapeSite(site)
		if err != nil {
			fmt.Printf("Error scraping %s: %v\n", site.Name, err)
			continue
		}
		allNews = append(allNews, news...)
	}

	return allNews, nil
}

func scrapeSite(site Site) ([]NewsItem, error) {
	var news []NewsItem

	c := colly.NewCollector(
		colly.AllowedDomains(getDomain(site.URL)),
	)

	c.OnHTML(site.Selector, func(e *colly.HTMLElement) {
		title := e.Text
		link := e.Attr("href")

		if link == "" {
			return
		}

		if link[0] == '/' {
			link = site.URL + link
		}

		if database.IsURLSeen(link) {
			return
		}

		news = append(news, NewsItem{
			Title: title,
			Link:  link,
			Site:  site.Name,
		})

		database.SaveArticle(link, title)
	})

	err := c.Visit(site.URL)
	return news, err
}

func getDomain(url string) string {
	if len(url) > 8 && url[:8] == "https://" {
		url = url[8:]
	} else if len(url) > 7 && url[:7] == "http://" {
		url = url[7:]
	}

	for i, char := range url {
		if char == '/' {
			return url[:i]
		}
	}
	return url
}
