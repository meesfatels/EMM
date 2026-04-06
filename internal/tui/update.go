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
		case tea.KeyCtrlU:
			m.viewport.HalfPageUp()
			return m, nil
		case tea.KeyCtrlD:
			m.viewport.HalfPageDown()
			return m, nil
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
				nm, cmd := m.handleSlash(input)
				m = nm.(chatModel)
				return m, cmd
			}

			// Add user message, then immediately finalize it into cache
			m.messages = append(m.messages, message{role: "user", content: input})
			m = m.finalizeLastMessage()

			// Start assistant message
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
		m = m.refreshContent()
		m = m.finalizeLastMessage() // Move finished assistant msg to cache
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		// Header (1) + Border Padding (2) + Textarea (3) + extra spacing (2)
		m.viewport.Height = msg.Height - 1 - 3 - 2
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width - 4)
		m.ready = true
		return m.refreshContent(), nil
	}

	var cmd tea.Cmd
	if !m.streaming {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Always update viewport for mouse/keys
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
