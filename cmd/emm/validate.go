package main

import (
	"fmt"

	"github.com/meesfatels/emm/internal/config"
	"github.com/meesfatels/emm/internal/runtime"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate EMM configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := config.Dir()
			if err != nil {
				return err
			}
			_, err = runtime.New(dir)
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}
			fmt.Println("Configuration valid.")
			return nil
		},
	}
}
