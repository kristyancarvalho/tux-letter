package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

const EnvConfigPath = "TUX_LETTER_CONFIG"

func Load(explicitPath string) (Config, string, error) {
	cfg := Default()

	path, err := Discover(explicitPath)
	if err != nil {
		return cfg, "", err
	}
	if path == "" {
		return cfg, "", nil
	}

	if err := decodeFile(path, &cfg); err != nil {
		return cfg, path, err
	}
	return cfg, path, nil
}

func Discover(explicitPath string) (string, error) {
	if explicitPath != "" {
		if !fileExists(explicitPath) {
			return "", fmt.Errorf("config file not found: %s", explicitPath)
		}
		return explicitPath, nil
	}

	if env := strings.TrimSpace(os.Getenv(EnvConfigPath)); env != "" {
		if !fileExists(env) {
			return "", fmt.Errorf("config file from %s not found: %s", EnvConfigPath, env)
		}
		return env, nil
	}

	for _, candidate := range candidatePaths() {
		if candidate != "" && fileExists(candidate) {
			return candidate, nil
		}
	}
	return "", nil
}

func candidatePaths() []string {
	paths := []string{"tux-letter.toml"}

	if xdg := strings.TrimSpace(os.Getenv("XDG_CONFIG_HOME")); xdg != "" {
		paths = append(paths, filepath.Join(xdg, "tux-letter", "config.toml"))
	}
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".config", "tux-letter", "config.toml"))
	}
	return paths
}

func decodeFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		if err := json.Unmarshal(data, cfg); err != nil {
			return fmt.Errorf("parsing JSON config %s: %w", path, err)
		}
	default:
		if err := toml.Unmarshal(data, cfg); err != nil {
			return fmt.Errorf("parsing TOML config %s: %w", path, err)
		}
	}
	return nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
