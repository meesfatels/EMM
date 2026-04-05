package loader

import (
	"fmt"
	"path/filepath"
)

const defaultBaseURL = "https://openrouter.ai/api/v1/chat/completions"

type Config struct {
	APIKey   string `yaml:"api_key"`
	BaseURL  string `yaml:"base_url"`
	Username string `yaml:"username"`
}

func (l *Loader) LoadConfig() (Config, error) {
	var c Config
	if err := readYAML(filepath.Join(l.baseDir, "emm.yaml"), &c); err != nil {
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
	return c, nil
}
