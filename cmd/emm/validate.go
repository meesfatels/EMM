package main

import (
	"fmt"

	"github.com/meesfatels/emm/internal/agent"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate EMM configuration",
		Run: func(cmd *cobra.Command, args []string) {
			agent.NewRuntime(agent.Dir())
			fmt.Println("Configuration valid.")
		},
	}
}
