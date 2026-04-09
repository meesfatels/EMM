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

type chatModel struct {
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
	eventCh      chan sessionEvent
	width        int
	ready        bool
	historyCache string
	lastWidth    int
	autoScroll   bool
}

func newChatModel(ctx context.Context, cancel context.CancelFunc, rt *agent.Runtime, session *agent.Session, agentName, minionName string) chatModel {
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

	return chatModel{
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
