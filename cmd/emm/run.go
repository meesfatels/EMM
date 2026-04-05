package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/meesfatels/emm/internal/config"
	"github.com/meesfatels/emm/internal/runtime"
	"github.com/meesfatels/emm/internal/tui"
	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	var minionFlag string
	cmd := &cobra.Command{
		Use:   "run [agent]",
		Short: "Run an agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAgent(args[0], minionFlag)
		},
	}
	cmd.Flags().StringVarP(&minionFlag, "minion", "m", "", "minion to use (required)")
	cmd.MarkFlagRequired("minion")
	return cmd
}

func runAgent(agentName string, minionName string) error {
	dir, err := config.Dir()
	if err != nil {
		return err
	}
	rt, err := runtime.New(dir)
	if err != nil {
		return err
	}
	agent, ok := rt.Agents[agentName]
	if !ok {
		return fmt.Errorf("unknown agent %q", agentName)
	}
	minion, ok := rt.Minions[minionName]
	if !ok {
		return fmt.Errorf("unknown minion %q", minionName)
	}
	session := runtime.NewSession(agent, minion, rt.Client)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	return tui.Run(ctx, cancel, rt, session, agentName, minionName)
}
