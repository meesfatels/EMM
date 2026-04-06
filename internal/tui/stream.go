package tui

import tea "github.com/charmbracelet/bubbletea"

type tokenMsg string
type doneMsg struct{ err error }

type streamEvent struct {
	token string
	err   error
	done  bool
}

func (m chatModel) sendMessage(input string) tea.Cmd {
	ch := m.tokenCh
	return func() tea.Msg {
		_, err := m.session.Send(m.ctx, input, func(token string) {
			ch <- streamEvent{token: token}
		})
		ch <- streamEvent{done: true, err: err}
		close(ch)
		return nil
	}
}

func (m chatModel) waitForToken() tea.Cmd {
	ch := m.tokenCh
	return func() tea.Msg {
		ev, ok := <-ch
		if !ok {
			return nil
		}
		if ev.done {
			return doneMsg{err: ev.err}
		}
		return tokenMsg(ev.token)
	}
}
