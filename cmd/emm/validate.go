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
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := agent.Dir()
			if err != nil {
				return err
			}
			if _, err = agent.NewRuntime(dir); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}
			fmt.Println("Configuration valid.")
			return nil
		},
	}
}
