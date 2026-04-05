package loader

import "path/filepath"

type Minion map[string]any

func (l *Loader) LoadMinions() (map[string]Minion, error) {
	return loadYAMLDir[Minion](filepath.Join(l.baseDir, "minions"))
}
