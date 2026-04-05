package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/meesfatels/emm/internal/runtime"
)

var (
	userStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	assistantStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
	headerStyle    = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Bold(true).Padding(0, 1)
	inputBorder    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62")).Padding(0, 0)
)

type tokenMsg string
type doneMsg struct{ err error }

type message struct {
	role    string
	content string
}

type chatModel struct {
	viewport   viewport.Model
	textarea   textarea.Model
	messages   []message
	session    *runtime.Session
	ctx        context.Context
	cancel     context.CancelFunc
	agentName  string
	minionName string
	streaming  bool
	tokenCh    chan string
	width      int
	height     int
	ready      bool
}

// Run starts the interactive chat TUI for the given session.
func Run(ctx context.Context, cancel context.CancelFunc, session *runtime.Session, agentName, minionName string) error {
	m := newChatModel(ctx, cancel, session, agentName, minionName)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func newChatModel(ctx context.Context, cancel context.CancelFunc, session *runtime.Session, agentName, minionName string) chatModel {
	ta := textarea.New()
	ta.Placeholder = "Type a message..."
	ta.Focus()
	ta.CharLimit = 0
	ta.ShowLineNumbers = false
	ta.SetHeight(3)
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(80, 20)

	return chatModel{
		viewport:   vp,
		textarea:   ta,
		session:    session,
		ctx:        ctx,
		cancel:     cancel,
		agentName:  agentName,
		minionName: minionName,
	}
}

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
			m.messages = append(m.messages, message{role: "user", content: input})
			m.messages = append(m.messages, message{role: "assistant", content: ""})
			m.streaming = true
			m.tokenCh = make(chan string, 256)
			m.updateViewportContent()
			cmds = append(cmds, m.sendMessage(input), m.waitForToken())
			return m, tea.Batch(cmds...)
		}

	case tokenMsg:
		if len(m.messages) > 0 {
			last := &m.messages[len(m.messages)-1]
			last.content += string(msg)
		}
		m.updateViewportContent()
		cmds = append(cmds, m.waitForToken())
		return m, tea.Batch(cmds...)

	case doneMsg:
		m.streaming = false
		if msg.err != nil {
			if len(m.messages) > 0 {
				last := &m.messages[len(m.messages)-1]
				last.content += fmt.Sprintf("\n[error: %v]", msg.err)
			}
		}
		m.updateViewportContent()
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		headerHeight := 1
		inputHeight := 5
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - headerHeight - inputHeight
		m.textarea.SetWidth(msg.Width - 2)
		m.updateViewportContent()
		return m, nil
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

func (m *chatModel) updateViewportContent() {
	var sb strings.Builder
	for _, msg := range m.messages {
		switch msg.role {
		case "user":
			sb.WriteString(userStyle.Render("You") + "\n")
		case "assistant":
			sb.WriteString(assistantStyle.Render("Assistant") + "\n")
		}
		sb.WriteString(msg.content)
		sb.WriteString("\n\n")
	}
	m.viewport.SetContent(sb.String())
	m.viewport.GotoBottom()
}

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

func (m chatModel) View() string {
	if !m.ready {
		return "Loading..."
	}
	header := headerStyle.Render(fmt.Sprintf(" emm — %s (%s) ", m.agentName, m.minionName))
	header = lipgloss.PlaceHorizontal(m.width, lipgloss.Left, header)
	input := inputBorder.Width(m.width - 2).Render(m.textarea.View())
	return header + "\n" + m.viewport.View() + "\n" + input
}
