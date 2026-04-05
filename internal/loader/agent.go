package loader

import (
	"fmt"
	"os"
	"path/filepath"
)

type Agent struct {
	Name     string
	Instinct *Instinct
}

func (l *Loader) LoadAgent(name string) (*Agent, error) {
	dir := filepath.Join(l.baseDir, "agents", name)
	instinct, err := LoadInstinct(
		filepath.Join(dir, "minion_instinct_guide.yaml"),
		filepath.Join(dir, "minion_instinct"),
	)
	if err != nil {
		return nil, fmt.Errorf("agent %s: %w", name, err)
	}
	return &Agent{Name: name, Instinct: instinct}, nil
}

func (l *Loader) LoadAgents() (map[string]*Agent, error) {
	dir := filepath.Join(l.baseDir, "agents")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading agents directory: %w", err)
	}
	agents := make(map[string]*Agent)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		agent, err := l.LoadAgent(e.Name())
		if err != nil {
			return nil, err
		}
		agents[e.Name()] = agent
	}
	return agents, nil
}
