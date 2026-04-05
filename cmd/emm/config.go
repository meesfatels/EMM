package main
import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"github.com/meesfatels/emm/internal/config"
)
func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Get or set configuration values",
	}
	cmd.AddCommand(newConfigGetCmd(), newConfigSetCmd())
	return cmd
}
func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfigMap()
			if err != nil {
				return err
			}
			val, ok := cfg[args[0]]
			if !ok {
				return fmt.Errorf("key %q not found in emm.yaml", args[0])
			}
			fmt.Println(val)
			return nil
		},
	}
}
func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfigMap()
			if err != nil {
				return err
			}
			cfg[args[0]] = args[1]
			dir, err := config.Dir()
			if err != nil {
				return err
			}
			data, err := yaml.Marshal(cfg)
			if err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(dir, "emm.yaml"), data, 0644)
		},
	}
}
func loadConfigMap() (map[string]any, error) {
	dir, err := config.Dir()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(dir, "emm.yaml"))
	if err != nil {
		return nil, err
	}
	var cfg map[string]any
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = make(map[string]any)
	}
	return cfg, nil
}
