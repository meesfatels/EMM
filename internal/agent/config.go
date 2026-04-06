package agent

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const defaultBaseURL = "https://openrouter.ai/api/v1/chat/completions"

type Config struct {
	APIKey        string `yaml:"api_key"`
	BaseURL       string `yaml:"base_url"`
	Username      string `yaml:"username"`
	DefaultAgent  string `yaml:"default_agent"`
	DefaultMinion string `yaml:"default_minion"`
}

func LoadConfig(dir string) (Config, error) {
	var c Config
	if err := readYAML(filepath.Join(dir, "emm.yaml"), &c); err != nil {
		return Config{}, fmt.Errorf("loading config: %w", err)
	}
	if c.APIKey == "" {
		return Config{}, fmt.Errorf("api_key not set in emm.yaml")
	}
	if c.BaseURL == "" {
		c.BaseURL = defaultBaseURL
	}
	if c.Username == "" {
		c.Username = "user"
	}
	if c.DefaultAgent == "" {
		c.DefaultAgent = "example"
	}
	if c.DefaultMinion == "" {
		c.DefaultMinion = "example"
	}
	return c, nil
}

// Dir returns the path to the user's EMM config directory (~/.emm).
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home directory: %w", err)
	}
	return filepath.Join(home, ".emm"), nil
}

// Init copies the embedded template into the EMM config directory,
// skipping files that already exist.
func Init(templateFS fs.FS) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	return fs.WalkDir(templateFS, ".EMM", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking template: %w", err)
		}
		rel, err := filepath.Rel(".EMM", path)
		if err != nil {
			return fmt.Errorf("resolving path %s: %w", path, err)
		}
		target := filepath.Join(dir, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if _, err := os.Stat(target); err == nil {
			return nil // already exists, skip
		}
		data, err := fs.ReadFile(templateFS, path)
		if err != nil {
			return fmt.Errorf("reading template %s: %w", rel, err)
		}
		return os.WriteFile(target, data, 0o644)
	})
}
