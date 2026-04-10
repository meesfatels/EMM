package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/meesfatels/emm/internal/agent"
	"github.com/meesfatels/emm/internal/tui"
	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	var agentFlag, minionFlag string
	cmd := &cobra.Command{
		Use:   "run [agent]",
		Short: "Run an agent",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rt := agent.NewRuntime(agent.Dir())

			agentName := agentFlag
			if len(args) > 0 {
				agentName = args[0]
			}
			if agentName == "" {
				agentName = rt.Config.DefaultAgent
			}
			minionName := minionFlag
			if minionName == "" {
				minionName = rt.Config.DefaultMinion
			}

			a, ok := rt.Agents[agentName]
			if !ok {
				panic("unknown agent: " + agentName)
			}
			m, ok := rt.Minions[minionName]
			if !ok {
				panic("unknown minion: " + minionName)
			}

			session := agent.NewSession(a, minionName, m, rt.Client, rt.Config.Username)
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()
			tui.Run(ctx, cancel, rt, session, agentName, minionName)
		},
	}
	cmd.Flags().StringVarP(&agentFlag, "agent", "a", "", "agent to use")
	cmd.Flags().StringVarP(&minionFlag, "minion", "m", "", "minion to use")
	return cmd
}
