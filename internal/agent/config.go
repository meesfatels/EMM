package agent

import (
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

func LoadConfig(dir string) Config {
	var c Config
	readYAML(filepath.Join(dir, "emm.yaml"), &c)
	if c.APIKey == "" {
		panic("api_key not set in emm.yaml")
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
	return c
}

func Dir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("home directory: " + err.Error())
	}
	return filepath.Join(home, ".emm")
}

func Init(templateFS fs.FS, force bool) {
	dir := Dir()
	fs.WalkDir(templateFS, ".EMM", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		rel, _ := filepath.Rel(".EMM", path)
		target := filepath.Join(dir, rel)
		if d.IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				panic(err)
			}
			return nil
		}
		if !force {
			if fileExists(target) {
				return nil
			}
		}
		data, err := fs.ReadFile(templateFS, path)
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile(target, data, 0o644); err != nil {
			panic(err)
		}
		return nil
	})
}
