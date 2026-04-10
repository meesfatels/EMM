package agent

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func readYAML(path string, out any) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic("reading " + filepath.Base(path) + ": " + err.Error())
	}
	if err := yaml.Unmarshal(data, out); err != nil {
		panic("parsing " + filepath.Base(path) + ": " + err.Error())
	}
}
