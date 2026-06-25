package config

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return "invalid configuration:\n  - " + strings.Join(e.Issues, "\n  - ")
}

func (c *Config) Validate() error {
	var issues []string

	if len(c.Sources.URLs) == 0 {
		issues = append(issues, "sources.urls must contain at least one URL")
	}
	for i, u := range c.Sources.URLs {
		if strings.TrimSpace(u) == "" {
			issues = append(issues, fmt.Sprintf("sources.urls[%d] is empty", i))
			continue
		}
		if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
			issues = append(issues, fmt.Sprintf("sources.urls[%d] %q must start with http:// or https://", i, u))
		}
	}

	if c.Fetch.Timeout.Duration() <= 0 {
		issues = append(issues, "fetch.timeout must be greater than zero")
	}
	if c.Fetch.MaxArticlesPerSource <= 0 {
		issues = append(issues, "fetch.max_articles_per_source must be greater than zero")
	}
	if c.Fetch.MaxTotalArticles <= 0 {
		issues = append(issues, "fetch.max_total_articles must be greater than zero")
	}

	if strings.TrimSpace(c.OpenRouter.BaseURL) == "" {
		issues = append(issues, "openrouter.base_url must not be empty")
	}
	if len(c.OpenRouter.Models) == 0 {
		issues = append(issues, "openrouter.models must contain at least one model")
	}
	if strings.TrimSpace(c.OpenRouter.APIKeyEnv) == "" {
		issues = append(issues, "openrouter.api_key_env must name an environment variable")
	}

	if strings.TrimSpace(c.App.Schedule) == "" {
		issues = append(issues, "app.schedule must not be empty")
	}
	if c.State.MaxSeenArticles <= 0 {
		issues = append(issues, "state.max_seen_articles must be greater than zero")
	}

	if c.Email.Enabled {
		for _, pair := range [][2]string{
			{"email.from_env", c.Email.FromEnv},
			{"email.to_env", c.Email.ToEnv},
			{"email.smtp_host_env", c.Email.SMTPHostEnv},
			{"email.smtp_port_env", c.Email.SMTPPortEnv},
		} {
			if strings.TrimSpace(pair[1]) == "" {
				issues = append(issues, pair[0]+" must be set when email is enabled")
			}
		}
	}

	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}
