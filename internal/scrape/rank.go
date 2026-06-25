package scrape

import (
	"net/url"
	"sort"
	"strings"

	"github.com/kristyancarvalho/tux-letter/internal/article"
	"golang.org/x/net/html"
)

var junkSegments = map[string]bool{
	"tag": true, "tags": true, "category": true, "categories": true,
	"topic": true, "topics": true, "author": true, "authors": true,
	"page": true, "login": true, "signin": true, "sign-in": true,
	"signup": true, "sign-up": true, "register": true, "subscribe": true,
	"account": true, "about": true, "contact": true, "privacy": true,
	"terms": true, "cookie": true, "cookies": true, "search": true,
	"feed": true, "rss": true, "atom": true, "newsletter": true,
	"advertise": true, "ads": true, "sponsor": true, "donate": true,
	"jobs": true, "careers": true, "legal": true, "faq": true,
	"help": true, "support": true, "shop": true, "store": true,
}

var socialHosts = map[string]bool{
	"twitter.com": true, "x.com": true, "facebook.com": true,
	"linkedin.com": true, "youtube.com": true, "youtu.be": true,
	"instagram.com": true, "reddit.com": true, "mastodon.social": true,
	"t.me": true, "telegram.org": true, "discord.com": true,
	"discord.gg": true, "patreon.com": true, "paypal.com": true,
	"tiktok.com": true, "threads.net": true, "bsky.app": true,
}

func acceptableLink(base *url.URL, target *url.URL) bool {
	if target.Scheme != "http" && target.Scheme != "https" {
		return false
	}
	host := strings.TrimPrefix(strings.ToLower(target.Hostname()), "www.")
	if host == "" {
		return false
	}
	if socialHosts[host] {
		return false
	}
	baseHost := strings.TrimPrefix(strings.ToLower(base.Hostname()), "www.")
	if host != baseHost {
		return false
	}
	if sameDocument(base, target) {
		return false
	}
	segments := pathSegments(target.Path)
	if len(segments) == 0 {
		return false
	}
	for _, s := range segments {
		if junkSegments[s] {
			return false
		}
	}
	return true
}

func sameDocument(base, target *url.URL) bool {
	bp := strings.TrimRight(base.Path, "/")
	tp := strings.TrimRight(target.Path, "/")
	return bp == tp && base.Hostname() == target.Hostname()
}

func pathSegments(p string) []string {
	var out []string
	for _, s := range strings.Split(p, "/") {
		s = strings.TrimSpace(strings.ToLower(s))
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

func looksLikeArticlePath(p string) int {
	segments := pathSegments(p)
	score := 0
	if len(segments) >= 2 {
		score += 2
	} else if len(segments) == 1 {
		score++
	}
	last := ""
	if len(segments) > 0 {
		last = segments[len(segments)-1]
	}
	if strings.Count(last, "-") >= 2 {
		score += 3
	} else if strings.Contains(last, "-") {
		score++
	}
	if strings.HasSuffix(last, ".html") || strings.HasSuffix(last, ".htm") {
		score++
	}
	if hasDateSegment(segments) {
		score += 2
	}
	if len(last) > 12 {
		score++
	}
	return score
}

func hasDateSegment(segments []string) bool {
	for _, s := range segments {
		if len(s) == 4 && isAllDigits(s) && (strings.HasPrefix(s, "19") || strings.HasPrefix(s, "20")) {
			return true
		}
	}
	return false
}

func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return s != ""
}

func rankAnchors(root *html.Node, baseURL string) []article.Article {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil
	}

	type scored struct {
		art   article.Article
		score int
		order int
	}
	var ranked []scored
	seen := make(map[string]bool)

	order := 0
	for _, a := range elements(root, "a") {
		href := strings.TrimSpace(attr(a, "href"))
		if href == "" {
			continue
		}
		ref, err := url.Parse(href)
		if err != nil {
			continue
		}
		resolved := base.ResolveReference(ref)
		if !acceptableLink(base, resolved) {
			continue
		}
		title := collapse(text(a))
		if len([]rune(title)) < 15 || len([]rune(title)) > 200 {
			continue
		}
		key := strings.TrimRight(resolved.String(), "/")
		if seen[key] {
			continue
		}
		seen[key] = true

		score := looksLikeArticlePath(resolved.Path)
		tl := len([]rune(title))
		if tl >= 25 && tl <= 130 {
			score += 2
		}
		ranked = append(ranked, scored{
			art:   article.Article{Title: title, URL: resolved.String()},
			score: score,
			order: order,
		})
		order++
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].score != ranked[j].score {
			return ranked[i].score > ranked[j].score
		}
		return ranked[i].order < ranked[j].order
	})

	var out []article.Article
	for _, r := range ranked {
		if r.score <= 0 {
			continue
		}
		out = append(out, r.art)
	}
	return out
}
