package tui

import tea "github.com/charmbracelet/bubbletea"

type tokenEvent string
type toolEvent struct{ name, input, output string }
type doneMsg struct{}

func (m model) sendMessage(input string) tea.Cmd {
	ch := m.eventCh
	return func() tea.Msg {
		m.session.Send(
			m.ctx,
			input,
			func(token string) { ch <- tokenEvent(token) },
			func(name, input, output string) { ch <- toolEvent{name, input, output} },
		)
		close(ch)
		return doneMsg{}
	}
}

func (m model) waitForEvent() tea.Cmd {
	ch := m.eventCh
	return func() tea.Msg {
		ev, ok := <-ch
		if !ok {
			return nil
		}
		return ev
	}
}
