package loader

import (
	"fmt"
	"os"
	"path/filepath"
)

type InstinctFile struct {
	Name           string `yaml:"name"`
	Interpretation string `yaml:"interpretation"`
}

type Instinct struct {
	Files   []InstinctFile    `yaml:"files"`
	Content map[string]string `yaml:"-"`
}

func LoadInstinct(guideFile, contentDir string) (*Instinct, error) {
	instinct := &Instinct{}
	if err := readYAML(guideFile, instinct); err != nil {
		return nil, fmt.Errorf("loading instinct guide: %w", err)
	}
	instinct.Content = make(map[string]string, len(instinct.Files))
	for _, f := range instinct.Files {
		data, err := os.ReadFile(filepath.Join(contentDir, f.Name))
		if err != nil {
			return nil, fmt.Errorf("loading instinct file %s: %w", f.Name, err)
		}
		instinct.Content[f.Name] = string(data)
	}
	return instinct, nil
}
