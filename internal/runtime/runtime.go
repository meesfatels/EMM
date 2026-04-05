package runtime

import (
	"fmt"

	"github.com/meesfatels/emm/internal/loader"
	"github.com/meesfatels/emm/internal/openrouter"
)

type Runtime struct {
	Dir     string
	Config  loader.Config
	Minions map[string]loader.Minion
	Agents  map[string]*loader.Agent
	Client  *openrouter.Client
}

func New(emmDir string) (*Runtime, error) {
	l := loader.NewLoader(emmDir)
	cfg, err := l.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}
	apiKey, err := cfg.APIKey()
	if err != nil {
		return nil, err
	}
	minions, err := l.LoadMinions()
	if err != nil {
		return nil, fmt.Errorf("loading minions: %w", err)
	}
	agents, err := l.LoadAgents()
	if err != nil {
		return nil, fmt.Errorf("loading agents: %w", err)
	}
	return &Runtime{
		Dir:     emmDir,
		Config:  cfg,
		Minions: minions,
		Agents:  agents,
		Client:  openrouter.NewClient(apiKey, cfg.BaseURL()),
	}, nil
}
