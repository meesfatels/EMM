package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m chatModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

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

	// Metadata label sits directly above the input box.
	meta := st.dim.Render(fmt.Sprintf(" %s  %s", m.agentName, m.minionName))

	borderColor := cfg.Colors.Accent
	if m.streaming {
		borderColor = cfg.Colors.System
	}
	inputBox := st.border.Copy().
		BorderForeground(lipgloss.Color(borderColor)).
		Width(m.width - 4).
		Render(m.textarea.View())

	parts = append(parts, meta, inputBox)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}
