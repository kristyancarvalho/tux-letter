package render

import "strings"

type locale struct {
	tagline     string
	subTagline  string
	briefing    string
	sources     string
	sourcesWord string
	emptyBody   string
	generated   string
	generatedBy string
}

func localeFor(language string) locale {
	switch primaryLang(language) {
	case "pt":
		return locale{
			tagline:     "boletim open-source cifrado",
			subTagline:  "linux // foss // segurança // sistemas",
			briefing:    "resumo",
			sources:     "fontes",
			sourcesWord: "fontes",
			emptyBody:   "Nenhum conteúdo de dispatch neste ciclo.",
			generated:   "gerado",
			generatedBy: "gerado localmente pelo tux-letter",
		}
	default:
		return locale{
			tagline:     Tagline,
			subTagline:  SubTagline,
			briefing:    "briefing",
			sources:     "sources",
			sourcesWord: "sources",
			emptyBody:   "No dispatch content in this cycle.",
			generated:   "generated",
			generatedBy: "generated locally by tux-letter",
		}
	}
}

func primaryLang(language string) string {
	lang := strings.ToLower(strings.TrimSpace(language))
	if lang == "" {
		return "en"
	}
	if i := strings.IndexAny(lang, "-_"); i > 0 {
		lang = lang[:i]
	}
	return lang
}

func htmlLang(language string) string {
	lang := strings.TrimSpace(language)
	if lang == "" {
		return "en"
	}
	return lang
}
