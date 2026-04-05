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
type InstinctGuide struct {
	Files []InstinctFile `yaml:"files"`
}
type Instinct struct {
	Guide   InstinctGuide
	Content map[string]string
}

func LoadInstinct(guideFile string, contentDir string) (*Instinct, error) {
	var guide InstinctGuide
	if err := readYAML(guideFile, &guide); err != nil {
		return nil, fmt.Errorf("loading instinct guide: %w", err)
	}
	content := make(map[string]string, len(guide.Files))
	for _, f := range guide.Files {
		path := filepath.Join(contentDir, f.Name)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("loading instinct file %s: %w", f.Name, err)
		}
		content[f.Name] = string(data)
	}
	return &Instinct{Guide: guide, Content: content}, nil
}
