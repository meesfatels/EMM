package loader
import (
	"fmt"
	"path/filepath"
)
type SubAgentRef struct {
	Name    string   `yaml:"name"`
	Minions []string `yaml:"minions"`
}
type Agent struct {
	Name       string
	SubAgents  []SubAgentRef
	Allowlists []string
	Instinct   *Instinct
}
func (l *Loader) LoadAgent(name string) (*Agent, error) {
	dir := filepath.Join(l.baseDir, "agents", name)
	var saConfig struct {
		SubAgents []SubAgentRef `yaml:"sub_agents"`
	}
	if err := readYAML(filepath.Join(dir, "sub_agents.yaml"), &saConfig); err != nil {
		return nil, fmt.Errorf("agent %s: %w", name, err)
	}
	var allowlists []string
	if err := readYAML(filepath.Join(dir, "allowlists.yaml"), &allowlists); err != nil {
		return nil, fmt.Errorf("agent %s: %w", name, err)
	}
	instinct, err := LoadInstinct(
		filepath.Join(dir, "minion_instinct_guide.yaml"),
		filepath.Join(dir, "minion_instinct"),
	)
	if err != nil {
		return nil, fmt.Errorf("agent %s: %w", name, err)
	}
	return &Agent{
		Name:       name,
		SubAgents:  saConfig.SubAgents,
		Allowlists: allowlists,
		Instinct:   instinct,
	}, nil
}
