package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) handleSlash(input string) (tea.Model, tea.Cmd) {
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
		m = m.sys(
			"/agent <name>   — switch agent\n" +
				"/minion <name>  — switch minion\n" +
				"/save <name>    — save conversation to .emm/conversations/<name>.md\n" +
				"/load <name>    — load conversation from .emm/conversations/<name>.md\n" +
				"/destroy <name> — delete a saved conversation\n" +
				"/help           — show this help\n\n" +
				"agents:  " + strings.Join(agents, ", ") + "\n" +
				"minions: " + strings.Join(minions, ", "),
		)

	case "/agent":
		if len(parts) < 2 {
			m = m.sys("usage: /agent <name>")
			break
		}
		name := parts[1]
		a, ok := m.rt.Agents[name]
		if !ok {
			m = m.sys(fmt.Sprintf("unknown agent %q", name))
			break
		}
		m.agentName = name
		m.session.SwitchAgent(a)
		m.historyCache = ""
		m = m.sys(fmt.Sprintf("switched to agent %q", name))

	case "/minion":
		if len(parts) < 2 {
			m = m.sys("usage: /minion <name>")
			break
		}
		name := parts[1]
		mn, ok := m.rt.Minions[name]
		if !ok {
			m = m.sys(fmt.Sprintf("unknown minion %q", name))
			break
		}
		m.minionName = name
		m.session.SwitchMinion(mn, name)
		m.historyCache = ""
		m = m.sys(fmt.Sprintf("switched to minion %q", name))

	case "/save":
		if len(parts) < 2 {
			m = m.sys("usage: /save <name>")
			break
		}
		name := parts[1]
		m.session.Save(m.rt.Dir, name)
		m = m.sys(fmt.Sprintf("saved as %q", name))

	case "/load":
		if len(parts) < 2 {
			m = m.sys("usage: /load <name>")
			break
		}
		name := parts[1]
		if !m.session.Load(m.rt.Dir, name) {
			m = m.sys(fmt.Sprintf("no conversation named %q", name))
			break
		}
		m.messages = nil
		m.historyCache = ""
		for _, msg := range m.session.Messages() {
			if msg.Role == "system" || msg.Role == "tool" || msg.Content == "" {
				continue
			}
			m.messages = append(m.messages, message{role: msg.Role, content: msg.Content})
		}
		m = m.sys(fmt.Sprintf("loaded %q", name))

	case "/destroy":
		if len(parts) < 2 {
			m = m.sys("usage: /destroy <name>")
			break
		}
		name := parts[1]
		path := filepath.Join(m.rt.Dir, "conversations", name+".md")
		if err := os.Remove(path); err != nil {
			if os.IsNotExist(err) {
				m = m.sys(fmt.Sprintf("no conversation named %q", name))
			} else {
				m = m.sys(fmt.Sprintf("error: %v", err))
			}
			break
		}
		m = m.sys(fmt.Sprintf("destroyed %q", name))

	default:
		m = m.sys(fmt.Sprintf("unknown command %q — try /help", parts[0]))
	}

	return m.refreshContent(), nil
}

func (m model) sys(text string) model {
	m.messages = append(m.messages, message{role: "system", content: text})
	return m
}
