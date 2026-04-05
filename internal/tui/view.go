package tui

import "fmt"

func (m chatModel) View() string {
	if !m.ready {
		return "\n  Loading..."
	}
	header := st.header.Width(m.width).Render(
		fmt.Sprintf("  emm  ·  %s  ·  %s", m.agentName, m.minionName),
	)
	input := st.border.Width(m.width - 2).Render(m.textarea.View())
	return header + "\n" + m.viewport.View() + "\n" + input
}
