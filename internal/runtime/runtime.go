package runtime
import (
	"fmt"
	"github.com/meesfatels/emm/internal/loader"
	"github.com/meesfatels/emm/internal/openrouter"
)
type Runtime struct {
	Config     loader.Config
	Minions    map[string]loader.Minion
	Allowlists map[string]loader.Allowlist
	Agents     map[string]*loader.Agent
	SubAgents  map[string]*loader.SubAgent
	Client     *openrouter.Client
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
	allowlists, err := l.LoadAllowlists()
	if err != nil {
		return nil, fmt.Errorf("loading allowlists: %w", err)
	}
	agents, err := l.LoadAgents()
	if err != nil {
		return nil, fmt.Errorf("loading agents: %w", err)
	}
	subAgents, err := l.LoadSubAgents()
	if err != nil {
		return nil, fmt.Errorf("loading sub-agents: %w", err)
	}
	rt := &Runtime{
		Config:     cfg,
		Minions:    minions,
		Allowlists: allowlists,
		Agents:     agents,
		SubAgents:  subAgents,
		Client:     openrouter.NewClient(apiKey, cfg.BaseURL()),
	}
	if err := rt.validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}
	return rt, nil
}
func (rt *Runtime) validate() error {
	for name, agent := range rt.Agents {
		for _, alName := range agent.Allowlists {
			if _, ok := rt.Allowlists[alName]; !ok {
				return fmt.Errorf("agent %s references unknown allowlist %q", name, alName)
			}
		}
		for _, saRef := range agent.SubAgents {
			if _, ok := rt.SubAgents[saRef.Name]; !ok {
				return fmt.Errorf("agent %s references unknown sub-agent %q", name, saRef.Name)
			}
			for _, mName := range saRef.Minions {
				if _, ok := rt.Minions[mName]; !ok {
					return fmt.Errorf("agent %s sub-agent %s references unknown minion %q", name, saRef.Name, mName)
				}
			}
		}
	}
	for name, sa := range rt.SubAgents {
		for _, alName := range sa.Allowlists {
			if _, ok := rt.Allowlists[alName]; !ok {
				return fmt.Errorf("sub-agent %s references unknown allowlist %q", name, alName)
			}
		}
	}
	return nil
}
