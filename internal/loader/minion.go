package loader

import "path/filepath"

// Minion is a raw map of model parameters. Every key-value pair is passed
// directly to the API, so any field supported by the model can be set in
// the minion YAML without any code changes.
type Minion map[string]any

func (l *Loader) LoadMinions() (map[string]Minion, error) {
	return loadYAMLDir[Minion](filepath.Join(l.baseDir, "minions"))
}
