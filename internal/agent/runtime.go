package agent

import (
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

func NewRuntime(dir string) *Runtime {
	cfg := LoadConfig(dir)
	return &Runtime{
		Dir:     dir,
		Config:  cfg,
		Minions: minion.Load(dir),
		Agents:  LoadAll(dir),
		Client:  openrouter.NewClient(cfg.APIKey, cfg.BaseURL),
	}
}
