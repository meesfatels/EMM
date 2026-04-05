package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/meesfatels/emm/internal/runtime"
)

func (m chatModel) handleSlash(input string) (tea.Model, tea.Cmd) {
	parts := strings.Fields(input)
	switch parts[0] {
	case "/help":
		agents := make([]string, 0, len(m.rt.Agents))
		for name := range m.rt.Agents {
			agents = append(agents, name)
		}
		minions := make([]string, 0, len(m.rt.Minions))
		for name := range m.rt.Minions {
			minions = append(minions, name)
		}
		m.messages = append(m.messages, message{role: "system", content: fmt.Sprintf(
			"/agent <name>  — switch agent (resets session)\n/minion <name> — switch minion (resets session)\n/help          — show this help\n\nagents:  %s\nminions: %s",
			strings.Join(agents, ", "),
			strings.Join(minions, ", "),
		)})

	case "/agent":
		if len(parts) < 2 {
			m.messages = append(m.messages, message{role: "system", content: "usage: /agent <name>"})
			break
		}
		name := parts[1]
		agent, ok := m.rt.Agents[name]
		if !ok {
			m.messages = append(m.messages, message{role: "system", content: fmt.Sprintf("unknown agent %q", name)})
			break
		}
		m.agentName = name
		m.session = runtime.NewSession(agent, m.rt.Minions[m.minionName], m.rt.Client)
		m.messages = []message{{role: "system", content: fmt.Sprintf("switched to agent %q", name)}}

	case "/minion":
		if len(parts) < 2 {
			m.messages = append(m.messages, message{role: "system", content: "usage: /minion <name>"})
			break
		}
		name := parts[1]
		minion, ok := m.rt.Minions[name]
		if !ok {
			m.messages = append(m.messages, message{role: "system", content: fmt.Sprintf("unknown minion %q", name)})
			break
		}
		m.minionName = name
		m.session = runtime.NewSession(m.rt.Agents[m.agentName], minion, m.rt.Client)
		m.messages = []message{{role: "system", content: fmt.Sprintf("switched to minion %q", name)}}

	default:
		m.messages = append(m.messages, message{role: "system", content: fmt.Sprintf("unknown command %q — try /help", parts[0])})
	}

	return m.refreshContent(), nil
}
