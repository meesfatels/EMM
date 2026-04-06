package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/meesfatels/emm/internal/runtime"
)

func (m chatModel) handleSlash(input string) (tea.Model, tea.Cmd) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return m, nil
	}

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
		m.messages = append(m.messages, message{
			role: "system",
			content: fmt.Sprintf(
				"/agent <name>   — switch agent (resets session)\n"+
					"/minion <name>  — switch minion (resets session)\n"+
					"/save <name>    — save conversation to .EMM/conversations/<name>.md\n"+
					"/load <name>    — load conversation from .EMM/conversations/<name>.md\n"+
					"/destroy <name> — delete a saved conversation\n"+
					"/help           — show this help\n\n"+
					"agents:  %s\n"+
					"minions: %s",
				strings.Join(agents, ", "),
				strings.Join(minions, ", ")),
		})

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
		m.session = runtime.NewSession(agent, m.minionName, m.rt.Minions[m.minionName], m.rt.Client, m.rt.Config.Username)
		m.messages = []message{{role: "system", content: fmt.Sprintf("switched to agent %q", name)}}
		m.historyCache = "" // Clear cache on switch

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
		m.session = runtime.NewSession(m.rt.Agents[m.agentName], name, minion, m.rt.Client, m.rt.Config.Username)
		m.messages = []message{{role: "system", content: fmt.Sprintf("switched to minion %q", name)}}
		m.historyCache = "" // Clear cache on switch

	case "/save":
		if len(parts) < 2 {
			m.messages = append(m.messages, message{role: "system", content: "usage: /save <name>"})
			break
		}
		name := parts[1]
		if err := m.session.Save(m.rt.Dir, name); err != nil {
			m.messages = append(m.messages, message{role: "system", content: fmt.Sprintf("error: %v", err)})
			break
		}
		m.messages = append(m.messages, message{role: "system", content: fmt.Sprintf("saved as %q", name)})
		// We don't clear cache on save

	case "/load":
		if len(parts) < 2 {
			m.messages = append(m.messages, message{role: "system", content: "usage: /load <name>"})
			break
		}
		name := parts[1]
		if err := m.session.Load(m.rt.Dir, name); err != nil {
			m.messages = append(m.messages, message{role: "system", content: fmt.Sprintf("error: %v", err)})
			break
		}
		// Update TUI messages
		m.messages = nil
		m.historyCache = "" // Reset cache for new conversation
		for _, msg := range m.session.Messages() {
			if msg.Role == "system" {
				continue
			}
			m.messages = append(m.messages, message{role: msg.Role, content: msg.Content})
		}
		// Rebuild history cache for all loaded messages except the (potential) system announcement
		m.messages = append(m.messages, message{role: "system", content: fmt.Sprintf("loaded %q", name)})

	case "/destroy":
		if len(parts) < 2 {
			m.messages = append(m.messages, message{role: "system", content: "usage: /destroy <name>"})
			break
		}
		name := parts[1]
		path := filepath.Join(m.rt.Dir, "conversations", name+".md")
		if err := os.Remove(path); err != nil {
			if os.IsNotExist(err) {
				m.messages = append(m.messages, message{role: "system", content: fmt.Sprintf("no conversation named %q", name)})
			} else {
				m.messages = append(m.messages, message{role: "system", content: fmt.Sprintf("error: %v", err)})
			}
			break
		}
		m.messages = append(m.messages, message{role: "system", content: fmt.Sprintf("destroyed %q", name)})

	default:
		m.messages = append(m.messages, message{role: "system", content: fmt.Sprintf("unknown command %q — try /help", parts[0])})
	}

	return m.refreshContent(), nil
}
