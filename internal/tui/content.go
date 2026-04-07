package tui

import (
	"strings"
)

func (m chatModel) contentWidth() int {
	w := m.width - 4
	if w < 20 {
		return 20
	}
	return w
}

// refreshContent updates the viewport's content and scrolls to the bottom.
// All but the last message are cached to keep streaming fast in long conversations.
func (m chatModel) refreshContent() chatModel {
	if m.width <= 0 {
		return m
	}

	width := m.contentWidth()

	if m.width != m.lastWidth {
		m.historyCache = ""
		m.lastWidth = m.width
	}

	if m.historyCache == "" && len(m.messages) > 1 {
		var sb strings.Builder
		for i := 0; i < len(m.messages)-1; i++ {
			sb.WriteString(renderMessage(m.messages[i], m.agentName, m.rt.Config.Username, width))
			sb.WriteString("\n")
		}
		m.historyCache = sb.String()
	}

	var sb strings.Builder
	sb.WriteString(m.historyCache)
	if len(m.messages) > 0 {
		sb.WriteString(renderMessage(m.messages[len(m.messages)-1], m.agentName, m.rt.Config.Username, width))
	}
	sb.WriteString("\n")

	m.viewport.SetContent(sb.String())
	if m.autoScroll {
		m.viewport.GotoBottom()
	}
	return m
}

// finalizeLastMessage moves the last message into the history cache.
func (m chatModel) finalizeLastMessage() chatModel {
	if len(m.messages) > 0 {
		m.historyCache += renderMessage(m.messages[len(m.messages)-1], m.agentName, m.rt.Config.Username, m.contentWidth())
		m.historyCache += "\n"
	}
	return m
}

func renderMessage(msg message, agentName, userName string, width int) string {
	switch msg.role {
	case "user":
		return st.user.Render(userName) + "\n" + st.msg.Width(width).Render(msg.content) + "\n"
	case "assistant":
		return st.assistant.Render(agentName) + "\n" + st.msg.Width(width).Render(msg.content) + "\n"
	case "tool":
		lines := strings.SplitN(msg.content, "\n", 2)
		header := lines[0]
		output := ""
		if len(lines) > 1 {
			output = lines[1]
		}
		rendered := st.toolHeader.Render("🔧 "+header) + "\n" + st.toolBody.Width(width).Render(output)
		return rendered + "\n"
	case "system":
		return st.system.Width(width).Render(msg.content) + "\n"
	default:
		return ""
	}
}
