package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/meesfatels/emm/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List EMM resources",
	}
	cmd.AddCommand(newListAgentsCmd(), newListMinionsCmd())
	return cmd
}

func newListAgentsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "agents",
		Short: "List all agents",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := config.Dir()
			if err != nil {
				return err
			}
			entries, err := os.ReadDir(filepath.Join(dir, "agents"))
			if err != nil {
				return err
			}
			for _, e := range entries {
				if e.IsDir() {
					fmt.Println(e.Name())
				}
			}
			return nil
		},
	}
}

func newListMinionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "minions",
		Short: "List all minions",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := config.Dir()
			if err != nil {
				return err
			}
			entries, err := os.ReadDir(filepath.Join(dir, "minions"))
			if err != nil {
				return err
			}
			for _, e := range entries {
				if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
					continue
				}
				name := strings.TrimSuffix(e.Name(), ".yaml")
				data, err := os.ReadFile(filepath.Join(dir, "minions", e.Name()))
				if err != nil {
					return err
				}
				var m map[string]any
				if err := yaml.Unmarshal(data, &m); err != nil {
					return err
				}
				model, _ := m["model"].(string)
				fmt.Printf("%s\t%s\n", name, model)
			}
			return nil
		},
	}
}
