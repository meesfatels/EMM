package tui

import "strings"

const pad = "  "

func (m chatModel) refreshContent() chatModel {
	var sb strings.Builder
	for _, msg := range m.messages {
		switch msg.role {
		case "user":
			sb.WriteString(pad + st.user.Render(m.rt.Config.Username()) + "\n")
			sb.WriteString(prefixLines(wordWrap(msg.content, m.viewport.Width-4), pad))
			sb.WriteString("\n\n")
		case "assistant":
			sb.WriteString(pad + st.assistant.Render(m.agentName) + "\n")
			sb.WriteString(prefixLines(wordWrap(msg.content, m.viewport.Width-4), pad))
			sb.WriteString("\n\n")
		case "system":
			sb.WriteString(st.system.Render(prefixLines(msg.content, pad)))
			sb.WriteString("\n\n")
		}
	}
	m.viewport.SetContent(sb.String())
	m.viewport.GotoBottom()
	return m
}

func prefixLines(text, prefix string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func wordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}
	var result strings.Builder
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if i > 0 {
			result.WriteByte('\n')
		}
		if len(line) <= width {
			result.WriteString(line)
			continue
		}
		col := 0
		for j, word := range strings.Fields(line) {
			if j == 0 {
				result.WriteString(word)
				col = len(word)
			} else if col+1+len(word) > width {
				result.WriteByte('\n')
				result.WriteString(word)
				col = len(word)
			} else {
				result.WriteByte(' ')
				result.WriteString(word)
				col += 1 + len(word)
			}
		}
	}
	return result.String()
}
