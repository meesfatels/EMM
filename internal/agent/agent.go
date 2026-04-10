package agent

import (
	"os"
	"path/filepath"

	"github.com/meesfatels/emm/internal/shell"
	"github.com/meesfatels/emm/internal/tool"
)

type InstinctFile struct {
	Name           string `yaml:"name"`
	Interpretation string `yaml:"interpretation"`
}

type instinctGuideYAML struct {
	Instinct []InstinctFile `yaml:"instinct"`
}

type allowlistYAML struct {
	Shell []shell.Rule `yaml:"shell"`
}

type Agent struct {
	Name     string
	Instinct []InstinctFile
	Content  map[string]string
	Shell    []shell.Rule
	Tools    []tool.Tool
}

func Load(dir, name string) *Agent {
	agentDir := filepath.Join(dir, "agents", name)

	var ig instinctGuideYAML
	readYAML(filepath.Join(agentDir, "instinct_guide.yaml"), &ig)

	var al allowlistYAML
	if p := filepath.Join(agentDir, "allowlist.yaml"); fileExists(p) {
		readYAML(p, &al)
	}

	content := make(map[string]string, len(ig.Instinct))
	for _, f := range ig.Instinct {
		data, err := os.ReadFile(filepath.Join(agentDir, "instinct", f.Name))
		if err != nil {
			panic("reading instinct file " + f.Name + ": " + err.Error())
		}
		content[f.Name] = string(data)
	}

	a := &Agent{
		Name:     name,
		Instinct: ig.Instinct,
		Content:  content,
		Shell:    al.Shell,
	}
	if len(al.Shell) > 0 {
		a.Tools = append(a.Tools, shell.NewExecutor(al.Shell))
	}
	return a
}

func LoadAll(dir string) map[string]*Agent {
	entries, err := os.ReadDir(filepath.Join(dir, "agents"))
	if err != nil {
		panic("reading agents directory: " + err.Error())
	}
	agents := make(map[string]*Agent)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		agents[e.Name()] = Load(dir, e.Name())
	}
	return agents
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
