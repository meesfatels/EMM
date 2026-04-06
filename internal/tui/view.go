package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m chatModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	// 1. Header
	header := st.header.Copy().Width(m.width).Render(
		fmt.Sprintf(" emm · %s · %s", m.agentName, m.minionName),
	)

	// 2. Viewport
	content := m.viewport.View()

	// 3. Input
	inputBorder := st.border.Copy().Width(m.width - 2)
	if m.streaming {
		inputBorder = inputBorder.BorderForeground(lipgloss.Color("240"))
	}
	input := inputBorder.Render(m.textarea.View())

	// Assemble
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		input,
	)
}
