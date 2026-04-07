package main

import (
	"fmt"
	"os"

	emm "github.com/meesfatels/emm"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:     "emm",
		Short:   "Eidolon Minion Manager",
		Version: version,
	}
	rootCmd.AddCommand(
		newInitCmd(emm.TemplateFS),
		newRunCmd(),
		newValidateCmd(),
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
