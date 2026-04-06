package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/meesfatels/emm/internal/agent"
	"github.com/meesfatels/emm/internal/tui"
	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	var minionFlag string
	var agentFlag string
	cmd := &cobra.Command{
		Use:   "run [agent]",
		Short: "Run an agent",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := agent.Dir()
			if err != nil {
				return err
			}
			rt, err := agent.NewRuntime(dir)
			if err != nil {
				return err
			}

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

			return runAgent(rt, agentName, minionName)
		},
	}
	cmd.Flags().StringVarP(&agentFlag, "agent", "a", "", "agent to use")
	cmd.Flags().StringVarP(&minionFlag, "minion", "m", "", "minion to use")
	return cmd
}

func runAgent(rt *agent.Runtime, agentName, minionName string) error {
	a, ok := rt.Agents[agentName]
	if !ok {
		return fmt.Errorf("unknown agent %q", agentName)
	}
	m, ok := rt.Minions[minionName]
	if !ok {
		return fmt.Errorf("unknown minion %q", minionName)
	}
	session := agent.NewSession(a, minionName, m, rt.Client, rt.Config.Username)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	return tui.Run(ctx, cancel, rt, session, agentName, minionName)
}
