package loader
import (
	"fmt"
	"os"
	"path/filepath"
)
type Loader struct {
	baseDir string
}
func NewLoader(baseDir string) *Loader {
	return &Loader{baseDir: baseDir}
}
func (l *Loader) LoadConfig() (Config, error) {
	var c Config
	if err := readYAML(filepath.Join(l.baseDir, "emm.yaml"), &c); err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}
	return c, nil
}
func (l *Loader) LoadMinions() (map[string]Minion, error) {
	return loadYAMLDir[Minion](filepath.Join(l.baseDir, "minions"))
}
func (l *Loader) LoadAllowlists() (map[string]Allowlist, error) {
	return loadYAMLDir[Allowlist](filepath.Join(l.baseDir, "allowlists"))
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
func (l *Loader) LoadSubAgents() (map[string]*SubAgent, error) {
	dir := filepath.Join(l.baseDir, "sub_agents")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading sub_agents directory: %w", err)
	}
	subAgents := make(map[string]*SubAgent)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		sa, err := l.LoadSubAgent(e.Name())
		if err != nil {
			return nil, err
		}
		subAgents[e.Name()] = sa
	}
	return subAgents, nil
}
