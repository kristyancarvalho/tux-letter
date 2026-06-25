package ai

import (
	"strings"
	"unicode"
)

func languageName(code string) string {
	switch strings.ToLower(strings.TrimSpace(code)) {
	case "", "en", "en-us", "en-gb":
		return "English"
	case "pt", "pt-br":
		return "Brazilian Portuguese (português do Brasil)"
	case "pt-pt":
		return "European Portuguese (português de Portugal)"
	case "es", "es-es", "es-419", "es-mx":
		return "Spanish (español)"
	case "fr", "fr-fr":
		return "French (français)"
	case "de", "de-de":
		return "German (Deutsch)"
	case "it", "it-it":
		return "Italian (italiano)"
	}
	return strings.TrimSpace(code)
}

func isEnglishCode(code string) bool {
	c := strings.ToLower(strings.TrimSpace(code))
	return c == "" || c == "en" || strings.HasPrefix(c, "en-")
}

var englishMarkers = []string{
	"the", "and", "with", "that", "this", "from", "have", "which",
	"while", "their", "about", "these", "there", "would", "could",
	"into", "than", "them", "what", "when", "because", "through",
}

func looksEnglish(text string) bool {
	words := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	if len(words) < 12 {
		return false
	}
	present := make(map[string]bool, len(words))
	for _, w := range words {
		present[w] = true
	}
	hits := 0
	for _, marker := range englishMarkers {
		if present[marker] {
			hits++
		}
	}
	return hits >= 4
}
