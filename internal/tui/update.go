package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

func (m chatModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.cancel()
			return m, tea.Quit
		case tea.KeyEnter:
			if m.streaming {
				break
			}
			input := strings.TrimSpace(m.textarea.Value())
			if input == "" {
				break
			}
			m.textarea.Reset()
			if strings.HasPrefix(input, "/") {
				return m.handleSlash(input)
			}
			m.messages = append(m.messages, message{role: "user", content: input})
			m.messages = append(m.messages, message{role: "assistant", content: ""})
			m.streaming = true
			m.tokenCh = make(chan string, 256)
			m = m.refreshContent()
			return m, tea.Batch(m.sendMessage(input), m.waitForToken())
		}

	case tokenMsg:
		if len(m.messages) > 0 {
			m.messages[len(m.messages)-1].content += string(msg)
		}
		m = m.refreshContent()
		return m, m.waitForToken()

	case doneMsg:
		m.streaming = false
		if msg.err != nil && len(m.messages) > 0 {
			m.messages[len(m.messages)-1].content += fmt.Sprintf("\n[error: %v]", msg.err)
		}
		return m.refreshContent(), nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.ready = true
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 1 - 5
		m.textarea.SetWidth(msg.Width - 2)
		return m.refreshContent(), nil
	}

	var cmd tea.Cmd
	if !m.streaming {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
