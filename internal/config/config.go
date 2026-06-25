package config

import "time"

type Config struct {
	App        AppConfig        `toml:"app" json:"app"`
	State      StateConfig      `toml:"state" json:"state"`
	Fetch      FetchConfig      `toml:"fetch" json:"fetch"`
	OpenRouter OpenRouterConfig `toml:"openrouter" json:"openrouter"`
	Newsletter NewsletterConfig `toml:"newsletter" json:"newsletter"`
	Email      EmailConfig      `toml:"email" json:"email"`
	Sources    SourcesConfig    `toml:"sources" json:"sources"`
}

type AppConfig struct {
	Timezone   string `toml:"timezone" json:"timezone"`
	Schedule   string `toml:"schedule" json:"schedule"`
	RunOnStart bool   `toml:"run_on_start" json:"run_on_start"`
	LogLevel   string `toml:"log_level" json:"log_level"`
}

type StateConfig struct {
	Dir             string `toml:"dir" json:"dir"`
	MaxSeenArticles int    `toml:"max_seen_articles" json:"max_seen_articles"`
}

type FetchConfig struct {
	Timeout              Duration `toml:"timeout" json:"timeout"`
	UserAgent            string   `toml:"user_agent" json:"user_agent"`
	MaxArticlesPerSource int      `toml:"max_articles_per_source" json:"max_articles_per_source"`
	MaxTotalArticles     int      `toml:"max_total_articles" json:"max_total_articles"`
}

type OpenRouterConfig struct {
	APIKeyEnv string   `toml:"api_key_env" json:"api_key_env"`
	BaseURL   string   `toml:"base_url" json:"base_url"`
	Models    []string `toml:"models" json:"models"`
}

type NewsletterConfig struct {
	Language       string `toml:"language" json:"language"`
	Tone           string `toml:"tone" json:"tone"`
	Title          string `toml:"title" json:"title"`
	IncludeLinks   bool   `toml:"include_links" json:"include_links"`
	IncludeSources bool   `toml:"include_sources" json:"include_sources"`
}

type EmailConfig struct {
	Enabled     bool   `toml:"enabled" json:"enabled"`
	FromEnv     string `toml:"from_env" json:"from_env"`
	ToEnv       string `toml:"to_env" json:"to_env"`
	SMTPHostEnv string `toml:"smtp_host_env" json:"smtp_host_env"`
	SMTPPortEnv string `toml:"smtp_port_env" json:"smtp_port_env"`
	SMTPUserEnv string `toml:"smtp_user_env" json:"smtp_user_env"`
	SMTPPassEnv string `toml:"smtp_pass_env" json:"smtp_pass_env"`
}

type SourcesConfig struct {
	URLs []string `toml:"urls" json:"urls"`
}

func Default() Config {
	return Config{
		App: AppConfig{
			Timezone:   "UTC",
			Schedule:   "20:00",
			RunOnStart: false,
			LogLevel:   "info",
		},
		State: StateConfig{
			Dir:             "~/.local/state/tux-letter",
			MaxSeenArticles: 1000,
		},
		Fetch: FetchConfig{
			Timeout:              Duration(15 * time.Second),
			UserAgent:            "tux-letter/3.0 (+https://github.com/kristyancarvalho/tux-letter)",
			MaxArticlesPerSource: 8,
			MaxTotalArticles:     25,
		},
		OpenRouter: OpenRouterConfig{
			APIKeyEnv: "OPENROUTER_API_KEY",
			BaseURL:   "https://openrouter.ai/api/v1",
			Models: []string{
				"openai/gpt-4o-mini",
				"anthropic/claude-3.5-haiku",
				"google/gemini-flash-1.5",
			},
		},
		Newsletter: NewsletterConfig{
			Language:       "en",
			Tone:           "natural, concise, technical, friendly",
			Title:          "Tux Letter",
			IncludeLinks:   true,
			IncludeSources: true,
		},
		Email: EmailConfig{
			Enabled:     false,
			FromEnv:     "TUX_LETTER_EMAIL_FROM",
			ToEnv:       "TUX_LETTER_EMAIL_TO",
			SMTPHostEnv: "TUX_LETTER_SMTP_HOST",
			SMTPPortEnv: "TUX_LETTER_SMTP_PORT",
			SMTPUserEnv: "TUX_LETTER_SMTP_USER",
			SMTPPassEnv: "TUX_LETTER_SMTP_PASS",
		},
		Sources: SourcesConfig{URLs: nil},
	}
}
