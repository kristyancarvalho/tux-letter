package config

import "testing"

func TestDefaultIsInvalidWithoutSources(t *testing.T) {
	cfg := Default()
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected default config to be invalid without sources")
	}
}

func TestValidateAcceptsMinimalConfig(t *testing.T) {
	cfg := Default()
	cfg.Sources.URLs = []string{"https://example.com"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid config, got: %v", err)
	}
}

func TestValidateRejectsBadSourceScheme(t *testing.T) {
	cfg := Default()
	cfg.Sources.URLs = []string{"ftp://example.com"}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for non-http source")
	}
}

func TestValidateRequiresEmailEnvsWhenEnabled(t *testing.T) {
	cfg := Default()
	cfg.Sources.URLs = []string{"https://example.com"}
	cfg.Email.Enabled = true
	cfg.Email.FromEnv = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for enabled email without env names")
	}
}
