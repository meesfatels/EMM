package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/meesfatels/emm/internal/agent"
)

type message struct {
	role    string
	content string
}

type model struct {
	viewport     viewport.Model
	textarea     textarea.Model
	messages     []message
	session      *agent.Session
	rt           *agent.Runtime
	ctx          context.Context
	cancel       context.CancelFunc
	agentName    string
	minionName   string
	streaming    bool
	eventCh      chan any
	width        int
	height       int
	ready        bool
	historyCache string
	lastWidth    int
	autoScroll   bool
}

func newModel(ctx context.Context, cancel context.CancelFunc, rt *agent.Runtime, session *agent.Session, agentName, minionName string) model {
	ta := textarea.New()
	ta.Placeholder = cfg.Input.Placeholder
	ta.Focus()
	ta.CharLimit = 0
	ta.ShowLineNumbers = false
	ta.SetHeight(1)
	ta.KeyMap.InsertNewline.SetEnabled(false)
	plain := lipgloss.NewStyle()
	ta.FocusedStyle.Base = plain
	ta.BlurredStyle.Base = plain
	ta.FocusedStyle.CursorLine = plain
	ta.BlurredStyle.CursorLine = plain
	ta.FocusedStyle.Text = plain
	ta.BlurredStyle.Text = plain
	ta.FocusedStyle.Prompt = plain
	ta.BlurredStyle.Prompt = plain
	ta.SetPromptFunc(0, func(int) string { return "" })

	vp := viewport.New(0, 0)
	vp.KeyMap.Up.SetEnabled(false)
	vp.KeyMap.Down.SetEnabled(false)

	return model{
		viewport:   vp,
		textarea:   ta,
		session:    session,
		rt:         rt,
		ctx:        ctx,
		cancel:     cancel,
		agentName:  agentName,
		minionName: minionName,
		autoScroll: true,
	}
}

// inputLines returns the number of visual rows the textarea content occupies,
// capped at 6.
func (m model) inputLines() int {
	w := m.width - 4 // matches textarea.SetWidth argument
	if w < 1 {
		return 1
	}
	lines := len([]rune(m.textarea.Value()))/w + 1
	if lines > 6 {
		lines = 6
	}
	return lines
}

// applyLayout resizes the textarea and viewport to match the current input size.
func (m model) applyLayout() model {
	lines := m.inputLines()
	m.textarea.SetHeight(lines)
	fixed := lines + 2 + 1 // textarea rows + border (top+bottom) + meta label
	if cfg.Layout.ShowHeader {
		fixed++
	}
	if cfg.Layout.ShowStatus {
		fixed++
	}
	vph := m.height - fixed
	if vph < 1 {
		vph = 1
	}
	m.viewport.Height = vph
	return m
}
