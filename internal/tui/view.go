package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	userStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	assistantStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
	systemStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
	headerStyle    = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Bold(true).Padding(0, 1)
	inputBorder    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
)

func (m chatModel) View() string {
	if !m.ready {
		return "Loading..."
	}
	header := lipgloss.PlaceHorizontal(m.width, lipgloss.Left, headerStyle.Render(fmt.Sprintf("emm — %s (%s)", m.agentName, m.minionName)))
	input := inputBorder.Width(m.width - 2).Render(m.textarea.View())
	return header + "\n" + m.viewport.View() + "\n" + input
}
