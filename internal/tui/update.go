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

		case tea.KeyUp:
			m.viewport.LineUp(1)
			m.autoScroll = false
			return m, nil
		case tea.KeyDown:
			m.viewport.LineDown(1)
			if m.viewport.AtBottom() {
				m.autoScroll = true
			}
			return m, nil
		case tea.KeyPgUp, tea.KeyCtrlU:
			m.viewport.HalfPageUp()
			m.autoScroll = false
			return m, nil
		case tea.KeyPgDown, tea.KeyCtrlD:
			m.viewport.HalfPageDown()
			if m.viewport.AtBottom() {
				m.autoScroll = true
			}
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

			m.messages = append(m.messages, message{role: "user", content: input})
			m = m.finalizeLastMessage()

			m.messages = append(m.messages, message{role: "assistant", content: ""})
			m.streaming = true
			m.eventCh = make(chan sessionEvent, 256)
			m.autoScroll = true
			m = m.refreshContent()
			return m, tea.Batch(m.sendMessage(input), m.waitForEvent())
		}

	case tokenEvent:
		if len(m.messages) > 0 {
			m.messages[len(m.messages)-1].content += string(msg)
		}
		m = m.refreshContent()
		return m, m.waitForEvent()

	case shellEvent:
		// Finalize or remove the current assistant bubble.
		if len(m.messages) > 0 && m.messages[len(m.messages)-1].role == "assistant" {
			last := &m.messages[len(m.messages)-1]
			if last.content != "" {
				m = m.finalizeLastMessage()
			} else {
				m.messages = m.messages[:len(m.messages)-1]
			}
		}
		// Add the shell execution record and open a new assistant bubble.
		m.messages = append(m.messages, message{role: "shell", content: msg.cmd + "\n" + msg.output})
		m = m.finalizeLastMessage()
		m.messages = append(m.messages, message{role: "assistant", content: ""})
		m = m.refreshContent()
		return m, m.waitForEvent()

	case doneMsg:
		m.streaming = false
		if msg.err != nil && len(m.messages) > 0 {
			m.messages[len(m.messages)-1].content += fmt.Sprintf("\n[error: %v]", msg.err)
		}
		m = m.refreshContent()
		m = m.finalizeLastMessage()
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.ready = true

		fixed := cfg.Layout.InputHeight + 4 + 1 // box (height + 2 border + 2 padding) + meta label
		if cfg.Layout.ShowHeader {
			fixed++
		}
		if cfg.Layout.ShowStatus {
			fixed++
		}
		m.viewport.Height = msg.Height - fixed
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width - 4)
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
