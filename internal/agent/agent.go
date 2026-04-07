package agent

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/meesfatels/emm/internal/shell"
	"github.com/meesfatels/emm/internal/tool"
)

type InstinctFile struct {
	Name           string `yaml:"name"`
	Interpretation string `yaml:"interpretation"`
}

type agentYAML struct {
	Instinct []InstinctFile `yaml:"instinct"`
	Shell    []shell.Rule   `yaml:"shell"`
}

type Agent struct {
	Name     string
	Instinct []InstinctFile
	Content  map[string]string
	Shell    []shell.Rule
	Tools    []tool.Tool
}

func Load(dir, name string) (*Agent, error) {
	agentDir := filepath.Join(dir, "agents", name)
	var cfg agentYAML
	if err := readYAML(filepath.Join(agentDir, "agent.yaml"), &cfg); err != nil {
		return nil, fmt.Errorf("agent %s: %w", name, err)
	}
	content := make(map[string]string, len(cfg.Instinct))
	for _, f := range cfg.Instinct {
		data, err := os.ReadFile(filepath.Join(agentDir, "instinct", f.Name))
		if err != nil {
			return nil, fmt.Errorf("agent %s: reading instinct file %s: %w", name, f.Name, err)
		}
		content[f.Name] = string(data)
	}
	a := &Agent{
		Name:     name,
		Instinct: cfg.Instinct,
		Content:  content,
		Shell:    cfg.Shell,
	}
	if len(cfg.Shell) > 0 {
		a.Tools = append(a.Tools, shell.NewExecutor(cfg.Shell))
	}
	return a, nil
}

func LoadAll(dir string) (map[string]*Agent, error) {
	entries, err := os.ReadDir(filepath.Join(dir, "agents"))
	if err != nil {
		return nil, fmt.Errorf("reading agents directory: %w", err)
	}
	agents := make(map[string]*Agent)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		a, err := Load(dir, e.Name())
		if err != nil {
			return nil, err
		}
		agents[e.Name()] = a
	}
	return agents, nil
}
