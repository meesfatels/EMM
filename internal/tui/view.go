package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m chatModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	// 1. Header (Fixed 1 line)
	header := st.header.Copy().Width(m.width).Render(
		fmt.Sprintf(" emm · %s · %s", m.agentName, m.minionName),
	)

	// 2. Viewport (Flex height)
	// We wrap the viewport in a style that provides some side padding
	viewportStyle := lipgloss.NewStyle().Padding(0, 1)
	viewport := viewportStyle.Width(m.width).Render(m.viewport.View())

	// 3. Status Bar (1 line)
	scrollStatus := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
	if m.viewport.AtBottom() {
		scrollStatus = "BOT"
	}
	autoStatus := ""
	if !m.autoScroll && !m.viewport.AtBottom() {
		autoStatus = " PAUSED (Down/End to resume) "
	}
	statusLine := st.dim.Copy().Width(m.width).Render(
		fmt.Sprintf(" %s %s", scrollStatus, autoStatus),
	)

	// 4. Input (Fixed 5 lines)
	inputBorder := st.border.Copy().Width(m.width - 2)
	if m.streaming {
		inputBorder = inputBorder.BorderForeground(lipgloss.Color("240"))
	}
	input := inputBorder.Render(m.textarea.View())

	// Vertical assembly: ensure exactly the terminal height is used
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		viewport,
		statusLine,
		input,
	)
}
