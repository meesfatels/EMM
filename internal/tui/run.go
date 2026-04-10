package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/meesfatels/emm/internal/agent"
)

func Run(ctx context.Context, cancel context.CancelFunc, rt *agent.Runtime, session *agent.Session, agentName, minionName string) {
	cfg = loadTheme(rt.Dir)
	st = buildStyles(cfg)
	tea.NewProgram(newModel(ctx, cancel, rt, session, agentName, minionName), tea.WithAltScreen()).Run()
}
