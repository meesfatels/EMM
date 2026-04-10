package main

import (
	"io/fs"

	"github.com/meesfatels/emm/internal/agent"
	"github.com/spf13/cobra"
)

func newInitCmd(templateFS fs.FS) *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize EMM configuration in ~/.emm/",
		Run: func(cmd *cobra.Command, args []string) {
			agent.Init(templateFS, force)
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "force overwrite existing files")
	return cmd
}
