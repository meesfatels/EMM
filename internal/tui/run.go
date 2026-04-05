package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/meesfatels/emm/internal/runtime"
)

func Run(ctx context.Context, cancel context.CancelFunc, rt *runtime.Runtime, session *runtime.Session, agentName, minionName string) error {
	st = buildStyles(loadTheme(rt.Dir))
	p := tea.NewProgram(newChatModel(ctx, cancel, rt, session, agentName, minionName), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
