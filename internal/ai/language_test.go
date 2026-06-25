package ai

import (
	"strings"
	"testing"

	"github.com/kristyancarvalho/tux-letter/internal/article"
	"github.com/kristyancarvalho/tux-letter/internal/config"
)

func ptConfig() config.NewsletterConfig {
	return config.NewsletterConfig{Title: "Tux Letter", Language: "pt-BR", Tone: "tecnico, direto"}
}

func TestBuildPromptEnforcesConfiguredLanguage(t *testing.T) {
	bundle := BuildBundle(ptConfig(), sampleArticles())
	p := BuildPrompt(ptConfig(), bundle)
	for _, want := range []string{"Brazilian Portuguese", "pt-BR", "tecnico, direto", "ENTIRE newsletter"} {
		if !strings.Contains(p.System, want) {
			t.Errorf("system prompt missing %q", want)
		}
	}
}

func TestValidateRejectsEnglishWhenPortugueseConfigured(t *testing.T) {
	bundle := BuildBundle(ptConfig(), sampleArticles())
	english := Article{
		Title: "Kernel and graphics roundup",
		Body: []Section{{
			Heading: "kernel and graphics",
			Paragraphs: []string{
				"The new kernel broadens hardware support and improves power management, while the graphics stack continues to mature with drivers that benefit developers [1][2].",
			},
		}},
		Sources: []Reference{{ID: 1}, {ID: 2}},
	}
	if err := english.ValidateAgainst(bundle); err == nil {
		t.Error("expected English output to be rejected when pt-BR is configured")
	}
}

func TestValidateAcceptsPortugueseWhenPortugueseConfigured(t *testing.T) {
	bundle := BuildBundle(ptConfig(), sampleArticles())
	pt := Article{
		Title: "Panorama do kernel e da pilha gráfica",
		Body: []Section{{
			Heading: "kernel e gráficos",
			Paragraphs: []string{
				"O novo kernel amplia o suporte a hardware e melhora o gerenciamento de energia, enquanto a pilha gráfica continua amadurecendo com drivers Vulkan mais rápidos para usuários e desenvolvedores [1][2].",
			},
		}},
		Sources: []Reference{{ID: 1}, {ID: 2}},
	}
	if err := pt.ValidateAgainst(bundle); err != nil {
		t.Errorf("expected Portuguese output to pass validation, got %v", err)
	}
}

func TestEnglishOutputStillValidWhenEnglishConfigured(t *testing.T) {
	bundle := BuildBundle(config.NewsletterConfig{Language: "en"}, sampleArticles())
	a := loadFixture(t, "valid.json")
	if err := a.ValidateAgainst(bundle); err != nil {
		t.Errorf("English output must remain valid for English config: %v", err)
	}
}

func TestFallbackArticleUsesConfiguredLanguage(t *testing.T) {
	a := FallbackArticle(ptConfig(), []article.Article{
		{Title: "Kernel 6.20", Source: "lwn", URL: "https://lwn.net/a", Excerpt: "novidades do kernel"},
	})
	if !strings.Contains(a.Summary, "Compilado autom") || !strings.Contains(a.Summary, "fontes") {
		t.Errorf("expected Portuguese fallback summary, got %q", a.Summary)
	}
	if a.Subtitle != "compilado automático" {
		t.Errorf("expected Portuguese fallback subtitle, got %q", a.Subtitle)
	}
}

func TestLanguageName(t *testing.T) {
	cases := map[string]string{
		"pt-BR": "Brazilian Portuguese (português do Brasil)",
		"en":    "English",
		"":      "English",
		"xx":    "xx",
	}
	for code, want := range cases {
		if got := languageName(code); got != want {
			t.Errorf("languageName(%q) = %q, want %q", code, got, want)
		}
	}
}
