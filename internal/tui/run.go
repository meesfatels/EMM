package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/meesfatels/emm/internal/agent"
)

func Run(ctx context.Context, cancel context.CancelFunc, rt *agent.Runtime, session *agent.Session, agentName, minionName string) error {
	cfg = loadTheme(rt.Dir)
	st = buildStyles(cfg)
	p := tea.NewProgram(newChatModel(ctx, cancel, rt, session, agentName, minionName), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
