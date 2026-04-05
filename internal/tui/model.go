package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/meesfatels/emm/internal/runtime"
)

type message struct {
	role    string
	content string
}

type chatModel struct {
	viewport   viewport.Model
	textarea   textarea.Model
	messages   []message
	session    *runtime.Session
	rt         *runtime.Runtime
	ctx        context.Context
	cancel     context.CancelFunc
	agentName  string
	minionName string
	streaming  bool
	tokenCh    chan string
	width      int
	ready      bool
}

func newChatModel(ctx context.Context, cancel context.CancelFunc, rt *runtime.Runtime, session *runtime.Session, agentName, minionName string) chatModel {
	ta := textarea.New()
	ta.Placeholder = "Type a message or /help..."
	ta.Focus()
	ta.CharLimit = 0
	ta.ShowLineNumbers = false
	ta.SetHeight(3)
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return chatModel{
		viewport:   viewport.New(80, 20),
		textarea:   ta,
		session:    session,
		rt:         rt,
		ctx:        ctx,
		cancel:     cancel,
		agentName:  agentName,
		minionName: minionName,
	}
}
