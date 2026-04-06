package main

import (
	"io/fs"

	"github.com/meesfatels/emm/internal/agent"
	"github.com/spf13/cobra"
)

func newInitCmd(templateFS fs.FS) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize EMM configuration in ~/.emm/",
		RunE: func(cmd *cobra.Command, args []string) error {
			return agent.Init(templateFS)
		},
	}
}
