package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func readYAML(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filepath.Base(path), err)
	}
	if err := yaml.Unmarshal(data, out); err != nil {
		return fmt.Errorf("parsing %s: %w", filepath.Base(path), err)
	}
	return nil
}

func loadYAMLDir[T any](dir string) (map[string]T, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", filepath.Base(dir), err)
	}
	result := make(map[string]T)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		var v T
		if err := readYAML(filepath.Join(dir, e.Name()), &v); err != nil {
			return nil, fmt.Errorf("loading %s: %w", e.Name(), err)
		}
		name := strings.TrimSuffix(e.Name(), ".yaml")
		result[name] = v
	}
	return result, nil
}
