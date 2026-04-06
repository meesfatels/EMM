package agent

import (
	"fmt"

	"github.com/meesfatels/emm/internal/minion"
	"github.com/meesfatels/emm/internal/openrouter"
)

type Runtime struct {
	Dir     string
	Config  Config
	Minions map[string]minion.Minion
	Agents  map[string]*Agent
	Client  *openrouter.Client
}

func NewRuntime(dir string) (*Runtime, error) {
	cfg, err := LoadConfig(dir)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}
	minions, err := minion.Load(dir)
	if err != nil {
		return nil, fmt.Errorf("loading minions: %w", err)
	}
	agents, err := LoadAll(dir)
	if err != nil {
		return nil, fmt.Errorf("loading agents: %w", err)
	}
	return &Runtime{
		Dir:     dir,
		Config:  cfg,
		Minions: minions,
		Agents:  agents,
		Client:  openrouter.NewClient(cfg.APIKey, cfg.BaseURL),
	}, nil
}
