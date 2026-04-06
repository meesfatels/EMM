package tui

import (
	"strings"
)

// refreshContent updates the viewport's content and scrolls to the bottom.
// It caches all but the last message to keep performance stable even
// in huge conversations during streaming.
func (m chatModel) refreshContent() chatModel {
	if m.width <= 0 {
		return m
	}

	width := m.width - 4
	if width < 20 {
		width = 20
	}

	// If the window resized, clear the cache so everything is re-wrapped.
	if m.width != m.lastWidth {
		m.historyCache = ""
		m.lastWidth = m.width
	}

	// Rebuild the history cache if it's empty.
	if m.historyCache == "" && len(m.messages) > 1 {
		var sb strings.Builder
		for i := 0; i < len(m.messages)-1; i++ {
			sb.WriteString(renderMessage(m.messages[i], m.agentName, m.rt.Config.Username, width))
			sb.WriteString("\n")
		}
		m.historyCache = sb.String()
	}

	// The current content is the history cache + the last message.
	var sb strings.Builder
	sb.WriteString(m.historyCache)
	if len(m.messages) > 0 {
		sb.WriteString(renderMessage(m.messages[len(m.messages)-1], m.agentName, m.rt.Config.Username, width))
	}
	// Add a little extra space at the bottom for breathing room
	sb.WriteString("\n")

	m.viewport.SetContent(sb.String())
	if m.autoScroll {
		m.viewport.GotoBottom()
	}
	return m
}

// finalizeLastMessage moves the last message into the history cache.
func (m chatModel) finalizeLastMessage() chatModel {
	width := m.width - 4
	if width < 20 {
		width = 20
	}
	if len(m.messages) > 0 {
		m.historyCache += renderMessage(m.messages[len(m.messages)-1], m.agentName, m.rt.Config.Username, width)
		m.historyCache += "\n"
	}
	return m
}

func renderMessage(msg message, agentName, userName string, width int) string {
	switch msg.role {
	case "user":
		name := st.user.Render(userName)
		content := st.msgUser.Width(width).Render(msg.content)
		return name + "\n" + content + "\n"
	case "assistant":
		name := st.assistant.Render(agentName)
		content := st.msgAssistant.Width(width).Render(msg.content)
		return name + "\n" + content + "\n"
	case "system":
		return st.system.Width(width).Render(msg.content) + "\n"
	default:
		return ""
	}
}
