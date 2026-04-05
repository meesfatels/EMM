package tui

import tea "github.com/charmbracelet/bubbletea"

type tokenMsg string
type doneMsg struct{ err error }

func (m chatModel) sendMessage(input string) tea.Cmd {
	ch := m.tokenCh
	return func() tea.Msg {
		_, err := m.session.Send(m.ctx, input, func(token string) {
			ch <- token
		})
		close(ch)
		return doneMsg{err: err}
	}
}

func (m chatModel) waitForToken() tea.Cmd {
	ch := m.tokenCh
	return func() tea.Msg {
		token, ok := <-ch
		if !ok {
			return nil
		}
		return tokenMsg(token)
	}
}
