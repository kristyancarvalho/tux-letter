package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return path
}

func TestLoadTOML(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "c.toml", `
[fetch]
timeout = "30s"
max_articles_per_source = 3

[sources]
urls = ["https://example.com", "https://example.org"]
`)

	cfg, used, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if used != path {
		t.Errorf("used = %q, want %q", used, path)
	}
	if cfg.Fetch.Timeout.Duration() != 30*time.Second {
		t.Errorf("timeout = %v, want 30s", cfg.Fetch.Timeout.Duration())
	}
	if len(cfg.Sources.URLs) != 2 {
		t.Errorf("sources = %d, want 2", len(cfg.Sources.URLs))
	}
	if cfg.OpenRouter.BaseURL == "" {
		t.Error("defaults should be preserved for unset fields")
	}
}

func TestLoadJSON(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "c.json", `{"sources":{"urls":["https://example.com"]},"fetch":{"timeout":"5s"}}`)

	cfg, _, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Fetch.Timeout.Duration() != 5*time.Second {
		t.Errorf("timeout = %v, want 5s", cfg.Fetch.Timeout.Duration())
	}
	if len(cfg.Sources.URLs) != 1 {
		t.Errorf("sources = %d, want 1", len(cfg.Sources.URLs))
	}
}

func TestLoadMissingExplicitPath(t *testing.T) {
	if _, _, err := Load(filepath.Join(t.TempDir(), "nope.toml")); err == nil {
		t.Fatal("expected error for missing explicit path")
	}
}

func TestDiscoverPrefersEnvVar(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "env.toml", "[sources]\nurls=[\"https://x.com\"]\n")
	t.Setenv(EnvConfigPath, path)

	got, err := Discover("")
	if err != nil {
		t.Fatalf("Discover: %v", err)
	}
	if got != path {
		t.Errorf("Discover = %q, want %q", got, path)
	}
}

func TestDiscoverExplicitBeatsEnv(t *testing.T) {
	dir := t.TempDir()
	envPath := writeFile(t, dir, "env.toml", "x=1\n")
	explicit := writeFile(t, dir, "explicit.toml", "x=1\n")
	t.Setenv(EnvConfigPath, envPath)

	got, err := Discover(explicit)
	if err != nil {
		t.Fatalf("Discover: %v", err)
	}
	if got != explicit {
		t.Errorf("Discover = %q, want %q", got, explicit)
	}
}

func TestDiscoverNoneFound(t *testing.T) {
	t.Setenv(EnvConfigPath, "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	got, err := Discover("")
	if err != nil {
		t.Fatalf("Discover: %v", err)
	}
	if got != "" {
		t.Errorf("Discover = %q, want empty", got)
	}
}

func TestExpandPath(t *testing.T) {
	t.Setenv("HOME", "/home/tester")
	if got := ExpandPath("~/.local/state"); got != "/home/tester/.local/state" {
		t.Errorf("ExpandPath = %q", got)
	}
	if got := ExpandPath("/abs/path"); got != "/abs/path" {
		t.Errorf("ExpandPath absolute = %q", got)
	}
}
