package minion

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Minion is a raw map of model parameters. Every key-value pair is passed
// directly to the API, so any field supported by the model can be set in
// the minion YAML without any code changes.
type Minion map[string]any

func Load(dir string) map[string]Minion {
	minionsDir := filepath.Join(dir, "minions")
	entries, err := os.ReadDir(minionsDir)
	if err != nil {
		panic("reading minions directory: " + err.Error())
	}
	result := make(map[string]Minion)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(minionsDir, e.Name()))
		if err != nil {
			panic(err)
		}
		var m Minion
		if err := yaml.Unmarshal(data, &m); err != nil {
			panic(err)
		}
		result[strings.TrimSuffix(e.Name(), ".yaml")] = m
	}
	return result
}
