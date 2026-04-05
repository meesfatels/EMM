package loader
import (
	"fmt"
	"path/filepath"
)
type SubAgent struct {
	Name       string
	Allowlists []string
	Instinct   *Instinct
}
func (l *Loader) LoadSubAgent(name string) (*SubAgent, error) {
	dir := filepath.Join(l.baseDir, "sub_agents", name)
	var allowlists []string
	if err := readYAML(filepath.Join(dir, "allowlists.yaml"), &allowlists); err != nil {
		return nil, fmt.Errorf("sub-agent %s: %w", name, err)
	}
	instinct, err := LoadInstinct(
		filepath.Join(dir, "sub_minion_instinct_guide.yaml"),
		filepath.Join(dir, "sub_minion_instinct"),
	)
	if err != nil {
		return nil, fmt.Errorf("sub-agent %s: %w", name, err)
	}
	return &SubAgent{
		Name:       name,
		Allowlists: allowlists,
		Instinct:   instinct,
	}, nil
}
