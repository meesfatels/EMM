package main
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"github.com/meesfatels/emm/internal/config"
)
func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List EMM resources",
	}
	cmd.AddCommand(
		newListAgentsCmd(),
		newListMinionsCmd(),
		newListAllowlistsCmd(),
		newListSubAgentsCmd(),
	)
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
func newListAllowlistsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "allowlists",
		Short: "List all allowlists",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := config.Dir()
			if err != nil {
				return err
			}
			entries, err := os.ReadDir(filepath.Join(dir, "allowlists"))
			if err != nil {
				return err
			}
			for _, e := range entries {
				if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
					continue
				}
				name := strings.TrimSuffix(e.Name(), ".yaml")
				data, err := os.ReadFile(filepath.Join(dir, "allowlists", e.Name()))
				if err != nil {
					return err
				}
				var items []string
				if err := yaml.Unmarshal(data, &items); err != nil {
					return err
				}
				fmt.Printf("%s\t%d entries\n", name, len(items))
			}
			return nil
		},
	}
}
func newListSubAgentsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sub-agents",
		Short: "List all sub-agents",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := config.Dir()
			if err != nil {
				return err
			}
			entries, err := os.ReadDir(filepath.Join(dir, "sub_agents"))
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
