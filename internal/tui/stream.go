package tui

import tea "github.com/charmbracelet/bubbletea"

// sessionEvent is the discriminated union of events from a running session.
type sessionEvent interface{ isSessionEvent() }

type tokenEvent string
type toolEvent struct{ name, input, output string }
type doneMsg struct{ err error }

func (tokenEvent) isSessionEvent() {}
func (toolEvent) isSessionEvent() {}

func (m chatModel) sendMessage(input string) tea.Cmd {
	ch := m.eventCh
	return func() tea.Msg {
		_, err := m.session.Send(
			m.ctx,
			input,
			func(token string) { ch <- tokenEvent(token) },
			func(name, input, output string) { ch <- toolEvent{name, input, output} },
		)
		close(ch)
		return doneMsg{err: err}
	}
}

func (m chatModel) waitForEvent() tea.Cmd {
	ch := m.eventCh
	return func() tea.Msg {
		ev, ok := <-ch
		if !ok {
			return nil
		}
		return ev
	}
}
