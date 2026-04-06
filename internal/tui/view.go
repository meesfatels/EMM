package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m chatModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	sep := st.dim.Render(strings.Repeat("─", m.width))

	inputStyle := lipgloss.NewStyle().Padding(0, 1)
	if m.streaming {
		inputStyle = inputStyle.Foreground(lipgloss.Color("240"))
	}
	input := inputStyle.Render(m.textarea.View())

	parts := []string{}

	if cfg.Layout.ShowHeader {
		label := fmt.Sprintf("%s  %s", m.agentName, m.minionName)
		parts = append(parts, st.header.Copy().Width(m.width).Render(label))
	}

	parts = append(parts, lipgloss.NewStyle().Padding(0, 1).Render(m.viewport.View()))

	if cfg.Layout.ShowStatus {
		status := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
		if m.viewport.AtBottom() {
			status = ""
		}
		if !m.autoScroll && !m.viewport.AtBottom() {
			status += "  scroll locked"
		}
		parts = append(parts, st.dim.Copy().Width(m.width).Render(" "+status))
	}

	parts = append(parts, sep, input)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}
