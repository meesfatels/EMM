package runtime
import (
	"fmt"
	"github.com/meesfatels/emm/internal/allowlist"
	"github.com/meesfatels/emm/internal/loader"
)
type ResolvedSubAgent struct {
	Name     string
	Minions  map[string]loader.Minion
	Instinct *loader.Instinct
	Enforcer *allowlist.Enforcer
}
type ResolvedAgent struct {
	Name      string
	Instinct  *loader.Instinct
	Enforcer  *allowlist.Enforcer
	SubAgents map[string]*ResolvedSubAgent
}
func (rt *Runtime) Resolve(agentName string) (*ResolvedAgent, error) {
	agent, ok := rt.Agents[agentName]
	if !ok {
		return nil, fmt.Errorf("unknown agent %q", agentName)
	}
	enforcer := rt.buildEnforcer(agent.Allowlists)
	subAgents := make(map[string]*ResolvedSubAgent, len(agent.SubAgents))
	for _, ref := range agent.SubAgents {
		sa, ok := rt.SubAgents[ref.Name]
		if !ok {
			return nil, fmt.Errorf("unknown sub-agent %q", ref.Name)
		}
		minions := make(map[string]loader.Minion, len(ref.Minions))
		for _, mName := range ref.Minions {
			m, ok := rt.Minions[mName]
			if !ok {
				return nil, fmt.Errorf("unknown minion %q", mName)
			}
			minions[mName] = m
		}
		subAgents[ref.Name] = &ResolvedSubAgent{
			Name:     ref.Name,
			Minions:  minions,
			Instinct: sa.Instinct,
			Enforcer: rt.buildEnforcer(sa.Allowlists),
		}
	}
	return &ResolvedAgent{
		Name:      agentName,
		Instinct:  agent.Instinct,
		Enforcer:  enforcer,
		SubAgents: subAgents,
	}, nil
}
func (rt *Runtime) buildEnforcer(names []string) *allowlist.Enforcer {
	lists := make([][]string, 0, len(names))
	for _, name := range names {
		if al, ok := rt.Allowlists[name]; ok {
			lists = append(lists, al)
		}
	}
	return allowlist.NewEnforcer(lists...)
}
